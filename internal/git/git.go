// Package git provides a thin wrapper around the git CLI with retry logic,
// dry-run support, and proper error handling for all common git operations.
//
// Set DryRun = true to print what would be executed without touching git.
// The Retry* variants of Push and Commit will retry up to 3 times on failure.
package git

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// DryRun controls whether git commands are actually executed.
// When true, all functions print the command and return nil.
var DryRun bool

// run executes a git command, streaming stdout/stderr to the parent process.
func run(args ...string) error {
	if DryRun {
		fmt.Printf("[DRY RUN] git %s\n", strings.Join(args, " "))
		return nil
	}

	cmd := exec.Command("git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git %s: %w", strings.Join(args, " "), err)
	}
	return nil
}

// withRetry executes fn up to maxAttempts times, waiting delay between each
// attempt. Returns the last error if all attempts fail.
func withRetry(fn func() error, maxAttempts int, delay time.Duration) error {
	var lastErr error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		if err := fn(); err != nil {
			lastErr = err
			if attempt < maxAttempts {
				fmt.Printf("⚠️  Attempt %d/%d failed: %v. Retrying in %s...\n", attempt, maxAttempts, err, delay)
				time.Sleep(delay)
			}
			continue
		}
		return nil
	}
	return fmt.Errorf("all %d attempts failed: %w", maxAttempts, lastErr)
}

// Pull fetches and merges changes from the specified remote and branch.
func Pull(remote, branch string) error {
	return run("pull", remote, branch)
}

// Add stages the specified files for the next commit.
func Add(files ...string) error {
	args := append([]string{"add"}, files...)
	return run(args...)
}

// Commit creates a new commit with the given message.
func Commit(message string) error {
	return run("commit", "-m", message)
}

// CommitWithDate creates a commit with a backdated author and committer date.
// Both GIT_AUTHOR_DATE and GIT_COMMITTER_DATE are set so the commit appears
// on the correct date in the GitHub contribution graph.
func CommitWithDate(message, date string) error {
	if DryRun {
		fmt.Printf("[DRY RUN] git commit -m %q --date %s\n", message, date)
		return nil
	}

	cmd := exec.Command("git", "commit", "-m", message, "--date", date)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(),
		"GIT_AUTHOR_DATE="+date,
		"GIT_COMMITTER_DATE="+date,
	)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git commit --date %s: %w", date, err)
	}
	return nil
}

// CommitWithDateRetry wraps CommitWithDate with up to 3 retries.
func CommitWithDateRetry(message, date string) error {
	return withRetry(func() error {
		return CommitWithDate(message, date)
	}, 3, 5*time.Second)
}

// Push pushes committed changes to the default remote.
func Push() error {
	return run("push")
}

// PushRetry wraps Push with up to 3 retries and exponential-ish backoff.
func PushRetry() error {
	return withRetry(Push, 3, 5*time.Second)
}

// CheckoutBranch switches to an existing branch.
func CheckoutBranch(name string) error {
	return run("checkout", name)
}

// CheckoutNewBranch creates a new branch and switches to it.
func CheckoutNewBranch(name string) error {
	return run("checkout", "-b", name)
}

// Merge merges the specified branch into the current branch without fast-forwarding,
// creating an explicit merge commit that is visible in the git graph.
func Merge(branch string) error {
	return run("merge", "--no-ff", branch, "-m", fmt.Sprintf("Merge branch '%s'", branch))
}

// DeleteBranch deletes a local branch (force delete).
func DeleteBranch(branch string) error {
	return run("branch", "-D", branch)
}
