// Package commit provides shared data types, a large built-in pool of ~200
// realistic commit messages and ~200 branch names, custom message/branch
// loading with append/replace modes, and the multi-file modification logic
// that writes to the bot-data/ directory.
package commit

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Data represents the JSON structure written to the tracking data file.
type Data struct {
	Date  string `json:"date"`
	Index int    `json:"index,omitempty"`
}

// WriteDataFile writes commit tracking data as formatted JSON.
func WriteDataFile(filePath string, data Data) error {
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return fmt.Errorf("mkdir %s: %w", filepath.Dir(filePath), err)
	}
	content, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal commit data: %w", err)
	}
	if err := os.WriteFile(filePath, content, 0644); err != nil {
		return fmt.Errorf("write data file %s: %w", filePath, err)
	}
	return nil
}

// ModifyRandomFile picks a realistic dummy file inside dataDir, creates its
// parent directories if needed, and appends realistic content to it.
// Returns the path of the modified file so it can be staged with git add.
func ModifyRandomFile(dataDir, dateStr string, index int) (string, error) {
	dummyFiles := []string{
		"data.json",
		"docs/api_notes.md",
		"logs/debug.log",
		"config/settings.json",
		"scripts/build.sh",
	}

	fileName := dummyFiles[rand.IntN(len(dummyFiles))]
	file := filepath.Join(dataDir, fileName)

	dir := filepath.Dir(file)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("mkdir %s: %w", dir, err)
	}

	ext := filepath.Ext(file)
	var content string
	switch ext {
	case ".json":
		data := Data{Date: dateStr, Index: index}
		bytes, _ := json.MarshalIndent(data, "", "  ")
		content = string(bytes) + "\n"
	case ".md":
		content = fmt.Sprintf("## Update Log\n- Last updated: %s\n- Revision index: %d\n\n", dateStr, index)
	case ".log":
		content = fmt.Sprintf("[%s] DEBUG: Heartbeat check ok (iter %d)\n", dateStr, index)
	case ".sh":
		content = fmt.Sprintf("#!/bin/bash\n# Build script generated %s\necho 'Building project (rev %d)...'\n", dateStr, index)
	default:
		content = fmt.Sprintf("Updated: %s\n", dateStr)
	}

	f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return "", fmt.Errorf("open file %s: %w", file, err)
	}
	defer f.Close()

	if _, err := f.WriteString(content); err != nil {
		return "", fmt.Errorf("write file %s: %w", file, err)
	}

	return file, nil
}

// FormatSpamMessage formats a backdated commit message using the date string.
func FormatSpamMessage(t time.Time) string {
	return t.Format(time.RFC3339)
}

// ────────────────────────────────────────────────────────────────────────────
// Built-in commit message pool (~200 messages across 16 categories)
// ────────────────────────────────────────────────────────────────────────────

var builtinMessages = []string{
	// 🐛 Bug fixes (15)
	"🐛 fix null pointer exception in user service",
	"🐛 resolve race condition in cache layer",
	"🐛 fix memory leak in connection pool",
	"🐛 correct off-by-one error in pagination",
	"🐛 fix broken redirect after login",
	"🐛 resolve 404 on nested routes",
	"🐛 fix timezone conversion bug",
	"🐛 correct validation logic for email fields",
	"🐛 fix crash on empty array input",
	"🐛 resolve deadlock in database transactions",
	"🐛 fix incorrect response status codes",
	"🐛 correct sorting behavior in list view",
	"🐛 fix broken image upload on mobile",
	"🐛 resolve session expiry not refreshing",
	"🐛 fix decimal rounding error in totals",

	// 🚑 Hotfixes (5)
	"🚑 hotfix: patch critical auth bypass vulnerability",
	"🚑 hotfix: fix data corruption on concurrent writes",
	"🚑 hotfix: resolve production outage in payment flow",
	"🚑 hotfix: correct broken deployment configuration",
	"🚑 hotfix: fix critical regression in search",

	// ✨ Features (25)
	"✨ add user profile page",
	"✨ implement dark mode toggle",
	"✨ add CSV export functionality",
	"✨ implement real-time notifications",
	"✨ add multi-language support",
	"✨ implement OAuth2 login with Google",
	"✨ add file upload with drag and drop",
	"✨ implement advanced search filters",
	"✨ add two-factor authentication",
	"✨ implement infinite scroll pagination",
	"✨ add keyboard shortcuts for power users",
	"✨ implement rate limiting middleware",
	"✨ add webhook support for integrations",
	"✨ implement cache warming strategy",
	"✨ add batch processing for large datasets",
	"✨ implement audit logging",
	"✨ add GraphQL endpoint",
	"✨ implement retry mechanism for failed jobs",
	"✨ add support for custom themes",
	"✨ implement role-based access control",
	"✨ add analytics dashboard",
	"✨ implement email digest feature",
	"✨ add API versioning support",
	"✨ implement content delivery network integration",
	"✨ add progress tracking for long-running tasks",

	// ♻️ Refactoring (20)
	"♻️ refactor authentication middleware",
	"♻️ extract payment service into separate module",
	"♻️ simplify database query builder",
	"♻️ decouple notification system from core logic",
	"♻️ refactor user repository pattern",
	"♻️ extract common validation utilities",
	"♻️ simplify error handling across services",
	"♻️ refactor config loading to use environment variables",
	"♻️ extract shared constants to dedicated file",
	"♻️ refactor API response formatting",
	"♻️ simplify routing configuration",
	"♻️ reduce cyclomatic complexity in parser module",
	"♻️ extract business logic from controllers",
	"♻️ refactor database connection pooling",
	"♻️ simplify retry logic implementation",
	"♻️ extract logging to centralized service",
	"♻️ decouple rendering engine from data layer",
	"♻️ refactor test helpers for reusability",
	"♻️ simplify state management logic",
	"♻️ extract shared types to common package",

	// 📝 Documentation (15)
	"📝 update API documentation",
	"📝 add inline code comments for complex logic",
	"📝 update installation guide",
	"📝 document environment variable configuration",
	"📝 add architecture decision records",
	"📝 update changelog",
	"📝 add contributing guidelines",
	"📝 document rate limiting behavior",
	"📝 add code examples to README",
	"📝 update deployment documentation",
	"📝 document database schema changes",
	"📝 add troubleshooting section to docs",
	"📝 update API endpoint references",
	"📝 document caching strategy",
	"📝 add performance benchmarks to docs",

	// 🔧 Config / Tooling / CI (15)
	"🔧 update ESLint configuration",
	"🔧 configure pre-commit hooks",
	"🔧 update CI pipeline for faster builds",
	"🔧 configure code coverage reporting",
	"🔧 update Docker base image",
	"🔧 configure automated dependency updates",
	"🔧 update build scripts for monorepo",
	"🔧 configure staging environment variables",
	"🔧 update Makefile targets",
	"🔧 configure load balancer health checks",
	"🔧 update deployment configuration",
	"🔧 configure automated database migrations",
	"🔧 update GitHub Actions workflows",
	"🔧 configure log rotation settings",
	"🔧 update type generation scripts",

	// ✅ Tests (15)
	"✅ add unit tests for auth service",
	"✅ improve integration test coverage",
	"✅ add snapshot tests for UI components",
	"✅ implement end-to-end test suite",
	"✅ add property-based testing",
	"✅ improve test isolation with mock factories",
	"✅ add load testing scenarios",
	"✅ fix flaky tests in CI environment",
	"✅ add contract tests for external APIs",
	"✅ improve test data management",
	"✅ add regression tests for fixed bugs",
	"✅ implement mutation testing",
	"✅ add performance benchmarks",
	"✅ improve test documentation",
	"✅ add accessibility tests for UI",

	// 🚀 Performance (10)
	"🚀 optimize database query performance",
	"🚀 implement query result caching",
	"🚀 reduce bundle size by lazy loading",
	"🚀 optimize image compression pipeline",
	"🚀 implement database connection pooling",
	"🚀 reduce API response time with indexing",
	"🚀 optimize memory usage in data processing",
	"🚀 implement CDN caching for static assets",
	"🚀 reduce cold start time for serverless functions",
	"🚀 optimize search indexing performance",

	// 📦 Dependencies (10)
	"📦 upgrade dependencies to latest versions",
	"📦 remove unused dependencies",
	"📦 replace deprecated library with modern alternative",
	"📦 pin dependency versions for reproducible builds",
	"📦 add missing peer dependencies",
	"📦 update lock file after dependency changes",
	"📦 migrate from deprecated package",
	"📦 add security patches for known vulnerabilities",
	"📦 consolidate duplicate dependencies",
	"📦 upgrade runtime version",

	// 🎨 Styling / Formatting (10)
	"🎨 apply consistent code formatting",
	"🎨 update color scheme to match design system",
	"🎨 improve responsive layout for mobile",
	"🎨 refactor CSS to use design tokens",
	"🎨 update typography scale",
	"🎨 improve accessibility with ARIA labels",
	"🎨 standardize button component styles",
	"🎨 fix alignment issues in grid layout",
	"🎨 update icon set to latest version",
	"🎨 improve loading state animations",

	// 🗄️ Database / Migrations (10)
	"🗄️ add database index for performance",
	"🗄️ create migration for new schema",
	"🗄️ optimize slow query with better index strategy",
	"🗄️ add foreign key constraints",
	"🗄️ create seed data for development",
	"🗄️ implement soft delete for user records",
	"🗄️ add database backup procedures",
	"🗄️ optimize join queries",
	"🗄️ implement data archiving strategy",
	"🗄️ add audit trail to sensitive tables",

	// 🔒 Security (10)
	"🔒 implement input sanitization",
	"🔒 add CSRF protection",
	"🔒 update authentication token expiry",
	"🔒 implement content security policy",
	"🔒 add rate limiting to auth endpoints",
	"🔒 encrypt sensitive data at rest",
	"🔒 implement API key rotation",
	"🔒 add security headers to responses",
	"🔒 implement OAuth token refresh flow",
	"🔒 add brute force protection",

	// 🌐 API / Networking (10)
	"🌐 add request timeout handling",
	"🌐 implement circuit breaker pattern",
	"🌐 add API response compression",
	"🌐 implement request deduplication",
	"🌐 add retry logic for external API calls",
	"🌐 implement request batching",
	"🌐 add API gateway configuration",
	"🌐 implement webhook retry mechanism",
	"🌐 add request tracing with correlation IDs",
	"🌐 implement graceful shutdown handling",

	// 🛠️ DevOps / Infrastructure (10)
	"🛠️ update Kubernetes deployment manifests",
	"🛠️ configure horizontal pod autoscaling",
	"🛠️ implement blue-green deployment strategy",
	"🛠️ update infrastructure as code",
	"🛠️ configure centralized logging",
	"🛠️ implement service mesh configuration",
	"🛠️ add monitoring dashboards",
	"🛠️ configure auto-scaling policies",
	"🛠️ update container resource limits",
	"🛠️ implement disaster recovery procedures",

	// 🚧 WIP / Drafts (10)
	"🚧 WIP: implement new checkout flow",
	"🚧 WIP: draft search algorithm improvements",
	"🚧 WIP: prototype real-time collaboration feature",
	"🚧 WIP: experimenting with new architecture",
	"🚧 WIP: draft implementation of background jobs",
	"🚧 initial implementation of recommendation engine",
	"🚧 draft: explore new caching strategies",
	"🚧 WIP: implement streaming response support",
	"🚧 initial scaffolding for plugin system",
	"🚧 draft: new onboarding flow design",

	// 📊 Data / Analytics (10)
	"📊 add metrics collection for user events",
	"📊 implement A/B testing framework",
	"📊 create analytics reporting pipeline",
	"📊 add conversion funnel tracking",
	"📊 implement data export to warehouse",
	"📊 add real-time dashboard metrics",
	"📊 implement cohort analysis",
	"📊 add session recording integration",
	"📊 create automated reporting jobs",
	"📊 implement data retention policies",
}

// ────────────────────────────────────────────────────────────────────────────
// Built-in branch name pool (~200 branch names across 8 prefixes)
// ────────────────────────────────────────────────────────────────────────────

var builtinBranchNames = []string{
	// feature/ (50)
	"feature/add-user-authentication",
	"feature/implement-search-functionality",
	"feature/dark-mode-toggle",
	"feature/payment-integration",
	"feature/real-time-notifications",
	"feature/export-to-csv",
	"feature/multi-language-support",
	"feature/two-factor-auth",
	"feature/api-versioning",
	"feature/file-upload",
	"feature/advanced-filters",
	"feature/keyboard-shortcuts",
	"feature/webhook-support",
	"feature/rate-limiting",
	"feature/batch-processing",
	"feature/audit-logging",
	"feature/graphql-endpoint",
	"feature/custom-themes",
	"feature/role-based-access",
	"feature/analytics-dashboard",
	"feature/email-digest",
	"feature/infinite-scroll",
	"feature/drag-and-drop",
	"feature/oauth-google-login",
	"feature/content-delivery-network",
	"feature/progress-tracking",
	"feature/retry-mechanism",
	"feature/cache-warming",
	"feature/recommendation-engine",
	"feature/onboarding-flow",
	"feature/plugin-system",
	"feature/streaming-responses",
	"feature/data-export",
	"feature/ab-testing",
	"feature/session-management",
	"feature/notification-preferences",
	"feature/user-profile-page",
	"feature/activity-feed",
	"feature/team-collaboration",
	"feature/smart-search",
	"feature/scheduled-tasks",
	"feature/api-key-management",
	"feature/custom-webhooks",
	"feature/data-visualization",
	"feature/mobile-responsiveness",
	"feature/accessibility-improvements",
	"feature/performance-dashboard",
	"feature/event-tracking",
	"feature/integration-marketplace",
	"feature/developer-portal",

	// bugfix/ (40)
	"bugfix/fix-login-redirect",
	"bugfix/resolve-memory-leak",
	"bugfix/fix-pagination-offset",
	"bugfix/correct-date-formatting",
	"bugfix/fix-null-check",
	"bugfix/resolve-race-condition",
	"bugfix/fix-broken-links",
	"bugfix/correct-timezone-handling",
	"bugfix/fix-form-validation",
	"bugfix/resolve-cors-issue",
	"bugfix/fix-session-expiry",
	"bugfix/correct-calculation-error",
	"bugfix/fix-image-upload",
	"bugfix/resolve-deadlock",
	"bugfix/fix-search-results",
	"bugfix/correct-sorting-order",
	"bugfix/fix-email-templates",
	"bugfix/resolve-api-timeout",
	"bugfix/fix-duplicate-entries",
	"bugfix/correct-error-messages",
	"bugfix/fix-mobile-layout",
	"bugfix/resolve-cache-invalidation",
	"bugfix/fix-csv-export",
	"bugfix/correct-permission-check",
	"bugfix/fix-notification-delivery",
	"bugfix/resolve-websocket-disconnect",
	"bugfix/fix-token-refresh",
	"bugfix/correct-decimal-rounding",
	"bugfix/fix-filter-logic",
	"bugfix/resolve-import-errors",
	"bugfix/fix-broken-tests",
	"bugfix/correct-sql-query",
	"bugfix/fix-infinite-loop",
	"bugfix/resolve-encoding-issue",
	"bugfix/fix-dependency-conflict",
	"bugfix/correct-routing-config",
	"bugfix/fix-svg-rendering",
	"bugfix/resolve-port-conflict",
	"bugfix/fix-scroll-behavior",
	"bugfix/correct-cache-headers",

	// hotfix/ (20)
	"hotfix/patch-auth-bypass",
	"hotfix/fix-data-corruption",
	"hotfix/resolve-production-crash",
	"hotfix/patch-sql-injection",
	"hotfix/fix-payment-failure",
	"hotfix/resolve-outage",
	"hotfix/patch-xss-vulnerability",
	"hotfix/fix-broken-deployment",
	"hotfix/resolve-database-connection",
	"hotfix/patch-csrf-vulnerability",
	"hotfix/fix-critical-api-error",
	"hotfix/resolve-service-down",
	"hotfix/patch-memory-overflow",
	"hotfix/fix-ssl-certificate",
	"hotfix/resolve-rate-limit-bypass",
	"hotfix/patch-data-leak",
	"hotfix/fix-authentication-loop",
	"hotfix/resolve-session-hijack",
	"hotfix/patch-dependency-vulnerability",
	"hotfix/fix-critical-regression",

	// chore/ (30)
	"chore/update-dependencies",
	"chore/cleanup-unused-imports",
	"chore/remove-dead-code",
	"chore/update-gitignore",
	"chore/organize-project-structure",
	"chore/update-environment-variables",
	"chore/cleanup-test-fixtures",
	"chore/remove-deprecated-methods",
	"chore/update-package-versions",
	"chore/cleanup-log-files",
	"chore/organize-assets",
	"chore/update-ci-config",
	"chore/cleanup-database-seeds",
	"chore/remove-console-logs",
	"chore/update-docker-config",
	"chore/cleanup-migration-files",
	"chore/remove-feature-flags",
	"chore/update-scripts",
	"chore/cleanup-temp-files",
	"chore/organize-documentation",
	"chore/update-eslint-config",
	"chore/cleanup-api-keys",
	"chore/remove-experimental-code",
	"chore/update-build-tools",
	"chore/cleanup-test-coverage",
	"chore/organize-routes",
	"chore/update-type-definitions",
	"chore/cleanup-unused-styles",
	"chore/remove-legacy-endpoints",
	"chore/update-deployment-scripts",

	// refactor/ (20)
	"refactor/extract-auth-service",
	"refactor/simplify-routing",
	"refactor/decouple-services",
	"refactor/improve-error-handling",
	"refactor/optimize-queries",
	"refactor/extract-shared-utilities",
	"refactor/simplify-state-management",
	"refactor/restructure-api-layer",
	"refactor/improve-code-readability",
	"refactor/extract-common-components",
	"refactor/simplify-config-loading",
	"refactor/decouple-ui-from-logic",
	"refactor/improve-test-structure",
	"refactor/extract-validation-rules",
	"refactor/simplify-caching-logic",
	"refactor/restructure-database-layer",
	"refactor/improve-logging-strategy",
	"refactor/extract-payment-module",
	"refactor/simplify-authentication-flow",
	"refactor/restructure-notification-service",

	// docs/ (15)
	"docs/update-api-reference",
	"docs/add-getting-started-guide",
	"docs/improve-readme",
	"docs/add-architecture-diagrams",
	"docs/update-changelog",
	"docs/add-deployment-guide",
	"docs/document-environment-setup",
	"docs/update-contributing-guide",
	"docs/add-code-examples",
	"docs/document-api-authentication",
	"docs/update-troubleshooting-guide",
	"docs/add-faq-section",
	"docs/document-release-process",
	"docs/update-security-policy",
	"docs/add-performance-guide",

	// test/ (15)
	"test/add-unit-tests-auth",
	"test/improve-integration-coverage",
	"test/add-e2e-tests",
	"test/fix-flaky-tests",
	"test/add-performance-benchmarks",
	"test/improve-mock-factories",
	"test/add-contract-tests",
	"test/setup-test-database",
	"test/add-load-tests",
	"test/improve-test-documentation",
	"test/add-mutation-tests",
	"test/configure-test-reporting",
	"test/add-snapshot-tests",
	"test/improve-test-isolation",
	"test/add-regression-tests",

	// release/ (10)
	"release/v1-0-0",
	"release/v1-1-0",
	"release/v1-2-0",
	"release/v2-0-0-beta",
	"release/v2-0-0",
	"release/v2-1-0",
	"release/v3-0-0-alpha",
	"release/v1-0-1",
	"release/v1-1-1",
	"release/v2-0-1",
}

// ────────────────────────────────────────────────────────────────────────────
// Public API: message and branch name selection
// ────────────────────────────────────────────────────────────────────────────

// LoadMessages returns the effective message pool based on the custom file and mode.
//   - If the file doesn't exist: returns the built-in pool.
//   - mode == "append": returns built-in + custom messages merged.
//   - mode == "replace": returns only custom messages (falls back to built-in if file is empty).
func LoadMessages(filePath, mode string) []string {
	custom := readLines(filePath)
	if len(custom) == 0 {
		return builtinMessages
	}
	if mode == "replace" {
		return custom
	}
	// Default: append
	merged := make([]string, len(builtinMessages)+len(custom))
	copy(merged, builtinMessages)
	copy(merged[len(builtinMessages):], custom)
	return merged
}

// LoadBranchNames returns the effective branch name pool based on the custom file and mode.
//   - If the file doesn't exist: returns the built-in pool.
//   - mode == "append": returns built-in + custom branch names merged.
//   - mode == "replace": returns only custom names (falls back to built-in if file is empty).
func LoadBranchNames(filePath, mode string) []string {
	custom := readLines(filePath)
	if len(custom) == 0 {
		return builtinBranchNames
	}
	if mode == "replace" {
		return custom
	}
	merged := make([]string, len(builtinBranchNames)+len(custom))
	copy(merged, builtinBranchNames)
	copy(merged[len(builtinBranchNames):], custom)
	return merged
}

// PickMessage returns a random message from the given pool.
func PickMessage(pool []string) string {
	if len(pool) == 0 {
		return builtinMessages[rand.IntN(len(builtinMessages))]
	}
	return pool[rand.IntN(len(pool))]
}

// PickBranchName returns a random branch name from the given pool.
func PickBranchName(pool []string) string {
	if len(pool) == 0 {
		return builtinBranchNames[rand.IntN(len(builtinBranchNames))]
	}
	return pool[rand.IntN(len(pool))]
}

// readLines reads a text file and returns non-empty, non-comment lines.
// Lines starting with '#' are treated as comments and skipped.
func readLines(path string) []string {
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		lines = append(lines, line)
	}
	return lines
}
