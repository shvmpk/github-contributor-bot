// daily is a CLI tool that makes configurable commits daily to maintain
// your GitHub contribution streak. Designed to run via GitHub Actions cron.
//
// Usage:
//
//	daily [flags]
//	  -config string         Path to bot.config.json (default "bot.config.json")
//	  -data-dir string       Directory for bot-generated files (default "bot-data")
//	  -min-commits int       Minimum commits on a normal day (default 1)
//	  -max-commits int       Maximum commits on a normal day (default 3)
//	  -skip-weekends         Skip Saturdays and Sundays (default false)
//	  -dry-run               Preview without making any git calls (default false)
//	  -messages-file string  Path to custom messages file (default "messages.txt")
//	  -messages-mode string  "append" or "replace" (default "append")
//	  -branch-names-file     Path to custom branch names file (default "branch_names.txt")
//	  -branch-mode string    "append" or "replace" (default "append")
//	  -stats-file string     Path to persistent stats file (default "stats.json")
//	  -version               Print version and exit
package main

import (
	"fmt"
	"math/rand/v2"
	"os"
	"time"

	"github.com/shvmpk/github-contribution-bot/internal/commit"
	"github.com/shvmpk/github-contribution-bot/internal/config"
	"github.com/shvmpk/github-contribution-bot/internal/git"
	"github.com/shvmpk/github-contribution-bot/internal/stats"
)

const version = "2.1.0"

func main() {
	for _, arg := range os.Args[1:] {
		if arg == "-version" || arg == "--version" {
			config.PrintVersion("daily", version)
			os.Exit(0)
		}
	}

	cfg, err := config.ParseDailyFlags(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if err := run(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run(cfg *config.DailyConfig) error {
	// Activate dry-run mode in the git package.
	git.DryRun = cfg.DryRun

	// Load stats.
	s, err := stats.Load(cfg.StatsFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: could not load stats: %v\n", err)
		s = &stats.Stats{}
	}
	s.RecordRun()
	if cfg.DryRun {
		s.SetDryRun()
	}

	// Load the effective message and branch pools.
	messages := commit.LoadMessages(cfg.MessagesFile, cfg.MessagesMode)
	branches := commit.LoadBranchNames(cfg.BranchNamesFile, cfg.BranchMode)

	fmt.Printf("🤖 Daily Commit Bot v%s\n", version)
	if cfg.DryRun {
		fmt.Println("🧪 DRY RUN mode — no git commands will be executed")
	}

	now := time.Now()

	// Work-Life Balance: Skip weekends if enabled.
	if cfg.SkipWeekends && (now.Weekday() == time.Saturday || now.Weekday() == time.Sunday) {
		fmt.Println("🏖️  Weekend detected! Work-Life Balance is enabled. Skipping commits.")
		s.RecordRestDay()
		_ = s.Save(cfg.StatsFile)
		s.PrintSummary()
		return nil
	}

	// 15% chance for a Rest Day.
	if rand.IntN(100) < 15 {
		fmt.Println("💤 Rest day! Taking a break. No commits today.")
		s.RecordRestDay()
		_ = s.Save(cfg.StatsFile)
		s.PrintSummary()
		return nil
	}

	// Determine commit count: 10% sprint, otherwise normal.
	commitCount := cfg.MinCommits
	if rand.IntN(100) < 10 {
		fmt.Println("🏃 Sprint day! Going into overdrive.")
		commitCount = cfg.MaxCommits + 1 + rand.IntN(5)
		s.RecordSprintDay()
	} else if cfg.MaxCommits > cfg.MinCommits {
		commitCount = cfg.MinCommits + rand.IntN(cfg.MaxCommits-cfg.MinCommits+1)
	}

	fmt.Printf("📝 Making %d commit(s) today\n\n", commitCount)

	// Initial pull to avoid conflicts.
	if err := git.Pull("origin", "main"); err != nil {
		return fmt.Errorf("initial pull failed: %w", err)
	}

	// 20% chance to work on a separate branch.
	useBranch := rand.IntN(100) < 20
	var branchName string
	if useBranch {
		branchName = commit.PickBranchName(branches)
		fmt.Printf("🔀 Creating branch: %s\n\n", branchName)
		if err := git.CheckoutNewBranch(branchName); err != nil {
			return fmt.Errorf("create branch %s: %w", branchName, err)
		}
		s.RecordBranch(branchName)
	}

	// Spread commits across a random window during the day (9 AM to 10 PM).
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	for i := 0; i < commitCount; i++ {
		randomMinutes := rand.IntN(780) // 13 hours * 60 minutes
		commitTime := startOfDay.Add(9 * time.Hour).Add(time.Duration(randomMinutes) * time.Minute)
		dateStr := commitTime.Format(time.RFC3339)

		fmt.Printf("── Commit %d/%d ─────────────────────────\n", i+1, commitCount)
		fmt.Printf("   Time: %s\n", dateStr)

		modifiedFile, err := commit.ModifyRandomFile(cfg.DataDir, dateStr, i)
		if err != nil {
			return fmt.Errorf("commit %d: modify file: %w", i+1, err)
		}

		if err := git.Add(modifiedFile); err != nil {
			return fmt.Errorf("commit %d: git add: %w", i+1, err)
		}

		message := commit.PickMessage(messages)
		if err := git.CommitWithDateRetry(message, dateStr); err != nil {
			return fmt.Errorf("commit %d: git commit: %w", i+1, err)
		}

		s.RecordCommit(modifiedFile)
		fmt.Printf("   ✅ %s\n   📄 %s\n\n", message, modifiedFile)
	}

	// Merge branch back into main if we created one.
	if useBranch {
		fmt.Printf("🔀 Merging %s → main\n", branchName)
		if err := git.CheckoutBranch("main"); err != nil {
			return fmt.Errorf("checkout main: %w", err)
		}
		if err := git.Merge(branchName); err != nil {
			return fmt.Errorf("merge %s: %w", branchName, err)
		}
		if err := git.DeleteBranch(branchName); err != nil {
			return fmt.Errorf("delete branch %s: %w", branchName, err)
		}
	}

	// Push everything once at the end with retry.
	fmt.Println("📤 Pushing commits...")
	if err := git.PushRetry(); err != nil {
		return fmt.Errorf("push failed: %w", err)
	}

	// Save stats and print summary.
	if err := s.Save(cfg.StatsFile); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: could not save stats: %v\n", err)
	}
	s.PrintSummary()

	fmt.Println("\n🎉 All commits completed successfully!")
	return nil
}
