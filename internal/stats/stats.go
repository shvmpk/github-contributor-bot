// Package stats provides persistent run statistics tracking for the
// GitHub contribution bot. Stats are stored in a JSON file and accumulate
// across runs, providing a weekly and all-time summary of bot activity.
package stats

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

const defaultStatsFile = "stats.json"

// Stats holds all cumulative run statistics for the bot.
type Stats struct {
	// All-time totals
	TotalRuns       int `json:"total_runs"`
	TotalCommits    int `json:"total_commits"`
	TotalRestDays   int `json:"total_rest_days"`
	TotalSprintDays int `json:"total_sprint_days"`
	TotalBranches   int `json:"total_branches"`

	// Weekly tracking (resets every Monday)
	WeeklyCommits    int    `json:"weekly_commits"`
	WeekStartDate    string `json:"week_start_date"`

	// Session (current run only, not persisted)
	sessionCommits  int
	sessionFiles    []string
	sessionBranch   string
	sessionDryRun   bool
	sessionRestDay  bool
	sessionSprint   bool
}

// Load reads stats from the given file path. If the file doesn't exist,
// a fresh Stats struct is returned (no error).
func Load(path string) (*Stats, error) {
	if path == "" {
		path = defaultStatsFile
	}

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &Stats{WeekStartDate: currentWeekStart()}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read stats file: %w", err)
	}

	var s Stats
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("parse stats file: %w", err)
	}

	// Reset weekly counter if we're in a new week.
	if s.WeekStartDate != currentWeekStart() {
		s.WeeklyCommits = 0
		s.WeekStartDate = currentWeekStart()
	}

	return &s, nil
}

// Save writes the stats to the given file path.
func (s *Stats) Save(path string) error {
	if path == "" {
		path = defaultStatsFile
	}

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal stats: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("write stats file: %w", err)
	}
	return nil
}

// RecordRun increments the run counter.
func (s *Stats) RecordRun() {
	s.TotalRuns++
}

// RecordCommit records a commit for the current session and all-time totals.
func (s *Stats) RecordCommit(file string) {
	s.sessionCommits++
	s.TotalCommits++
	s.WeeklyCommits++
	s.sessionFiles = append(s.sessionFiles, file)
}

// RecordRestDay marks the current session as a rest day.
func (s *Stats) RecordRestDay() {
	s.sessionRestDay = true
	s.TotalRestDays++
}

// RecordSprintDay marks the current session as a sprint day.
func (s *Stats) RecordSprintDay() {
	s.sessionSprint = true
	s.TotalSprintDays++
}

// RecordBranch records a branch name used in the current session.
func (s *Stats) RecordBranch(branch string) {
	s.sessionBranch = branch
	s.TotalBranches++
}

// SetDryRun marks this session as a dry run.
func (s *Stats) SetDryRun() {
	s.sessionDryRun = true
}

// PrintSummary prints a formatted summary of the current session and all-time stats.
func (s *Stats) PrintSummary() {
	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("📊 Run Summary")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	if s.sessionDryRun {
		fmt.Println("🧪 Mode:            DRY RUN (no changes made)")
	}

	if s.sessionRestDay {
		fmt.Println("💤 Today:           Rest Day")
	} else if s.sessionSprint {
		fmt.Printf("🏃 Today:           Sprint Day (%d commits)\n", s.sessionCommits)
	} else {
		fmt.Printf("📝 Today:           %d commit(s)\n", s.sessionCommits)
	}

	if s.sessionBranch != "" {
		fmt.Printf("🔀 Branch used:     %s\n", s.sessionBranch)
	}

	if len(s.sessionFiles) > 0 {
		fmt.Printf("🗂️  Files modified:  %v\n", unique(s.sessionFiles))
	}

	fmt.Println()
	fmt.Println("📈 All-Time Stats")
	fmt.Printf("   Total runs:     %d\n", s.TotalRuns)
	fmt.Printf("   Total commits:  %d\n", s.TotalCommits)
	fmt.Printf("   Rest days:      %d\n", s.TotalRestDays)
	fmt.Printf("   Sprint days:    %d\n", s.TotalSprintDays)
	fmt.Printf("   Branches used:  %d\n", s.TotalBranches)
	fmt.Printf("   This week:      %d commits\n", s.WeeklyCommits)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}

// currentWeekStart returns the ISO week start (Monday) as a YYYY-MM-DD string.
func currentWeekStart() string {
	now := time.Now()
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7 // Sunday = 7 in ISO
	}
	monday := now.AddDate(0, 0, -(weekday - 1))
	return monday.Format("2006-01-02")
}

// unique returns deduplicated slice of strings.
func unique(ss []string) []string {
	seen := make(map[string]struct{})
	result := []string{}
	for _, s := range ss {
		if _, ok := seen[s]; !ok {
			seen[s] = struct{}{}
			result = append(result, s)
		}
	}
	return result
}
