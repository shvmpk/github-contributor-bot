// spam is a CLI tool that creates a series of backdated commits to fill
// gaps in your GitHub contribution graph. All commits are pushed at once
// at the end for efficiency.
//
// Usage:
//
//	spam [flags]
//	  -config string         Path to bot.config.json (default "bot.config.json")
//	  -data-dir string       Directory for bot-generated files (default "bot-data")
//	  -count int             Number of backdated commits to create (default 100)
//	  -weeks-back int        Weeks into the past to spread commits (default 54)
//	  -dry-run               Preview without making any git calls (default false)
//	  -messages-file string  Path to custom messages file (default "messages.txt")
//	  -messages-mode string  "append" or "replace" (default "append")
//	  -version               Print version and exit
//
// ⚠️  WARNING: This tool modifies your git history. Use responsibly.
package main

import (
	"fmt"
	"math/rand/v2"
	"os"
	"time"

	"github.com/shvmpk/github-contribution-bot/internal/commit"
	"github.com/shvmpk/github-contribution-bot/internal/config"
	"github.com/shvmpk/github-contribution-bot/internal/git"
)

const version = "2.1.0"

func main() {
	for _, arg := range os.Args[1:] {
		if arg == "-version" || arg == "--version" {
			config.PrintVersion("spam", version)
			os.Exit(0)
		}
	}

	cfg, err := config.ParseSpamFlags(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if err := run(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run(cfg *config.SpamConfig) error {
	git.DryRun = cfg.DryRun

	messages := commit.LoadMessages(cfg.MessagesFile, cfg.MessagesMode)

	fmt.Printf("🤖 Spam Commit Bot v%s\n", version)
	if cfg.DryRun {
		fmt.Println("🧪 DRY RUN mode — no git commands will be executed")
	}
	fmt.Printf("📝 Creating %d backdated commits\n", cfg.Count)
	fmt.Printf("📅 Spread over %d weeks into the past\n", cfg.WeeksBack)
	fmt.Printf("📁 Data directory: %s\n\n", cfg.DataDir)

	now := time.Now()
	baseDate := now.AddDate(-1, 0, 1)

	for i := 0; i < cfg.Count; i++ {
		weekOffset := rand.IntN(cfg.WeeksBack)
		dayOffset := rand.IntN(7)

		commitDate := baseDate.
			AddDate(0, 0, weekOffset*7).
			AddDate(0, 0, dayOffset)

		dateStr := commitDate.Format(time.RFC3339)
		fmt.Printf("[%d/%d] %s\n", i+1, cfg.Count, dateStr)

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
	}

	fmt.Println("\n📤 Pushing all commits...")
	if err := git.PushRetry(); err != nil {
		return fmt.Errorf("push failed: %w", err)
	}

	fmt.Printf("\n🎉 Successfully created and pushed %d backdated commits!\n", cfg.Count)
	return nil
}
