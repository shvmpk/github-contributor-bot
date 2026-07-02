// Package config provides CLI flag parsing, config file loading, and
// validation for both the daily and spam commit modes. CLI flags always
// take precedence over values in bot.config.json.
package config

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
)

const defaultConfigFile = "bot.config.json"

// fileConfig mirrors the structure of bot.config.json for unmarshalling.
// All fields are optional; zero values mean "not set in file".
type fileConfig struct {
	MinCommits      *int    `json:"min_commits"`
	MaxCommits      *int    `json:"max_commits"`
	SkipWeekends    *bool   `json:"skip_weekends"`
	DryRun          *bool   `json:"dry_run"`
	DataDir         *string `json:"data_dir"`
	MessagesFile    *string `json:"messages_file"`
	MessagesMode    *string `json:"messages_mode"`
	BranchNamesFile *string `json:"branch_names_file"`
	BranchMode      *string `json:"branch_mode"`
}

// DailyConfig holds the fully resolved configuration for daily commit mode.
type DailyConfig struct {
	// DataDir is the directory where all auto-generated files are written (default: "bot-data").
	DataDir string
	// MinCommits is the minimum number of commits on a normal day (default: 1).
	MinCommits int
	// MaxCommits is the maximum number of commits on a normal day (default: 3).
	MaxCommits int
	// SkipWeekends skips committing on Saturday and Sunday (default: false).
	SkipWeekends bool
	// DryRun prints what would happen without making any git calls (default: false).
	DryRun bool
	// MessagesFile is the path to a custom commit messages file (default: "messages.txt").
	MessagesFile string
	// MessagesMode controls how the custom file interacts with built-in messages.
	// "append" = merge; "replace" = override (default: "append").
	MessagesMode string
	// BranchNamesFile is the path to a custom branch names file (default: "branch_names.txt").
	BranchNamesFile string
	// BranchMode controls how the custom file interacts with built-in branch names.
	// "append" = merge; "replace" = override (default: "append").
	BranchMode string
	// StatsFile is the path to the persistent stats JSON file (default: "stats.json").
	StatsFile string
}

// SpamConfig holds the fully resolved configuration for spam commit mode.
type SpamConfig struct {
	// DataDir is the directory where all auto-generated files are written (default: "bot-data").
	DataDir string
	// Count is the total number of backdated commits to create (default: 100).
	Count int
	// WeeksBack is how many weeks into the past to spread commits (default: 54).
	WeeksBack int
	// DryRun prints what would happen without making any git calls (default: false).
	DryRun bool
	// MessagesFile is the path to a custom commit messages file (default: "messages.txt").
	MessagesFile string
	// MessagesMode controls how the custom file interacts with built-in messages (default: "append").
	MessagesMode string
}

// ParseDailyFlags parses CLI flags for the daily commit mode, then merges
// values from bot.config.json (CLI flags take precedence).
func ParseDailyFlags(args []string) (*DailyConfig, error) {
	fs := flag.NewFlagSet("daily", flag.ContinueOnError)

	// Defaults
	cfg := &DailyConfig{}
	var configFile string

	fs.StringVar(&configFile, "config", defaultConfigFile, "Path to bot.config.json")
	fs.StringVar(&cfg.DataDir, "data-dir", "", "Directory for auto-generated bot files (default: bot-data)")
	fs.IntVar(&cfg.MinCommits, "min-commits", 0, "Minimum commits per run (default: 1)")
	fs.IntVar(&cfg.MaxCommits, "max-commits", 0, "Maximum commits per run (default: 3)")
	fs.BoolVar(&cfg.SkipWeekends, "skip-weekends", false, "Skip committing on weekends")
	fs.BoolVar(&cfg.DryRun, "dry-run", false, "Preview without making any git calls")
	fs.StringVar(&cfg.MessagesFile, "messages-file", "", "Path to custom commit messages file (default: messages.txt)")
	fs.StringVar(&cfg.MessagesMode, "messages-mode", "", "How to use custom messages: append or replace (default: append)")
	fs.StringVar(&cfg.BranchNamesFile, "branch-names-file", "", "Path to custom branch names file (default: branch_names.txt)")
	fs.StringVar(&cfg.BranchMode, "branch-mode", "", "How to use custom branch names: append or replace (default: append)")
	fs.StringVar(&cfg.StatsFile, "stats-file", "stats.json", "Path to the persistent stats file")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	// Load config file and apply defaults from it (CLI flags override).
	fc, _ := loadFile(configFile) // Ignore error — missing file is fine.
	applyDailyFileDefaults(cfg, fc)

	if err := cfg.validate(); err != nil {
		fs.Usage()
		return nil, err
	}

	return cfg, nil
}

// applyDailyFileDefaults fills zero-value fields from the config file.
func applyDailyFileDefaults(cfg *DailyConfig, fc *fileConfig) {
	if fc == nil {
		return
	}
	if cfg.DataDir == "" && fc.DataDir != nil {
		cfg.DataDir = *fc.DataDir
	}
	if cfg.MinCommits == 0 && fc.MinCommits != nil {
		cfg.MinCommits = *fc.MinCommits
	}
	if cfg.MaxCommits == 0 && fc.MaxCommits != nil {
		cfg.MaxCommits = *fc.MaxCommits
	}
	if !cfg.SkipWeekends && fc.SkipWeekends != nil {
		cfg.SkipWeekends = *fc.SkipWeekends
	}
	if !cfg.DryRun && fc.DryRun != nil {
		cfg.DryRun = *fc.DryRun
	}
	if cfg.MessagesFile == "" && fc.MessagesFile != nil {
		cfg.MessagesFile = *fc.MessagesFile
	}
	if cfg.MessagesMode == "" && fc.MessagesMode != nil {
		cfg.MessagesMode = *fc.MessagesMode
	}
	if cfg.BranchNamesFile == "" && fc.BranchNamesFile != nil {
		cfg.BranchNamesFile = *fc.BranchNamesFile
	}
	if cfg.BranchMode == "" && fc.BranchMode != nil {
		cfg.BranchMode = *fc.BranchMode
	}
	// Apply hardcoded defaults for anything still zero.
	if cfg.DataDir == "" {
		cfg.DataDir = "bot-data"
	}
	if cfg.MinCommits == 0 {
		cfg.MinCommits = 1
	}
	if cfg.MaxCommits == 0 {
		cfg.MaxCommits = 3
	}
	if cfg.MessagesFile == "" {
		cfg.MessagesFile = "messages.txt"
	}
	if cfg.MessagesMode == "" {
		cfg.MessagesMode = "append"
	}
	if cfg.BranchNamesFile == "" {
		cfg.BranchNamesFile = "branch_names.txt"
	}
	if cfg.BranchMode == "" {
		cfg.BranchMode = "append"
	}
}

// validate checks DailyConfig values are sensible.
func (c *DailyConfig) validate() error {
	if c.MinCommits < 1 {
		return errors.New("min-commits must be at least 1")
	}
	if c.MaxCommits < c.MinCommits {
		return fmt.Errorf("max-commits (%d) must be >= min-commits (%d)", c.MaxCommits, c.MinCommits)
	}
	if c.MaxCommits > 20 {
		return errors.New("max-commits cannot exceed 20 (safety limit)")
	}
	if c.MessagesMode != "append" && c.MessagesMode != "replace" {
		return fmt.Errorf("messages-mode must be 'append' or 'replace', got: %q", c.MessagesMode)
	}
	if c.BranchMode != "append" && c.BranchMode != "replace" {
		return fmt.Errorf("branch-mode must be 'append' or 'replace', got: %q", c.BranchMode)
	}
	return nil
}

// ParseSpamFlags parses CLI flags for the spam commit mode.
func ParseSpamFlags(args []string) (*SpamConfig, error) {
	fs := flag.NewFlagSet("spam", flag.ContinueOnError)
	cfg := &SpamConfig{}
	var configFile string

	fs.StringVar(&configFile, "config", defaultConfigFile, "Path to bot.config.json")
	fs.StringVar(&cfg.DataDir, "data-dir", "", "Directory for auto-generated bot files (default: bot-data)")
	fs.IntVar(&cfg.Count, "count", 0, "Number of backdated commits to create (default: 100)")
	fs.IntVar(&cfg.WeeksBack, "weeks-back", 0, "Weeks into the past to spread commits (default: 54)")
	fs.BoolVar(&cfg.DryRun, "dry-run", false, "Preview without making any git calls")
	fs.StringVar(&cfg.MessagesFile, "messages-file", "", "Path to custom commit messages file")
	fs.StringVar(&cfg.MessagesMode, "messages-mode", "", "How to use custom messages: append or replace")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	fc, _ := loadFile(configFile)
	applySpamFileDefaults(cfg, fc)

	if err := cfg.validate(); err != nil {
		fs.Usage()
		return nil, err
	}

	return cfg, nil
}

// applySpamFileDefaults fills zero-value fields from the config file.
func applySpamFileDefaults(cfg *SpamConfig, fc *fileConfig) {
	if cfg.DataDir == "" {
		if fc != nil && fc.DataDir != nil {
			cfg.DataDir = *fc.DataDir
		} else {
			cfg.DataDir = "bot-data"
		}
	}
	if cfg.Count == 0 {
		cfg.Count = 100
	}
	if cfg.WeeksBack == 0 {
		cfg.WeeksBack = 54
	}
	if cfg.MessagesFile == "" {
		if fc != nil && fc.MessagesFile != nil {
			cfg.MessagesFile = *fc.MessagesFile
		} else {
			cfg.MessagesFile = "messages.txt"
		}
	}
	if cfg.MessagesMode == "" {
		if fc != nil && fc.MessagesMode != nil {
			cfg.MessagesMode = *fc.MessagesMode
		} else {
			cfg.MessagesMode = "append"
		}
	}
}

// validate checks SpamConfig values are sensible.
func (c *SpamConfig) validate() error {
	if c.Count < 1 {
		return errors.New("count must be at least 1")
	}
	if c.Count > 1000 {
		return errors.New("count cannot exceed 1000 (safety limit)")
	}
	if c.WeeksBack < 1 {
		return errors.New("weeks-back must be at least 1")
	}
	if c.WeeksBack > 104 {
		return errors.New("weeks-back cannot exceed 104 (2 years)")
	}
	return nil
}

// loadFile reads and parses bot.config.json. Returns nil, nil if file not found.
func loadFile(path string) (*fileConfig, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read config file %s: %w", path, err)
	}

	var fc fileConfig
	if err := json.Unmarshal(data, &fc); err != nil {
		return nil, fmt.Errorf("parse config file %s: %w", path, err)
	}
	return &fc, nil
}

// PrintVersion prints the version string.
func PrintVersion(name, version string) {
	fmt.Fprintf(os.Stdout, "%s version %s\n", name, version)
}
