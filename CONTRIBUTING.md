# Contributing to GitHub Contribution Bot

First off, thanks for taking the time to contribute! 🎉

The following is a set of guidelines for contributing to this repository. These are mostly guidelines, not rules. Use your best judgment, and feel free to propose changes to this document in a pull request.

## Development Setup

This project is written entirely in **Go** and has **zero external dependencies**.

1. Install [Go (1.22+)](https://go.dev/)
2. Clone the repo:
   ```bash
   git clone https://github.com/shvmpk/github-contributor-bot.git
   cd github-contribution-bot
   ```
3. Make your changes in `cmd/` or `internal/` packages.

## Build & Verify

After making changes, always verify the code compiles and passes static analysis:

```bash
# Build all packages
go build ./... && echo "✅ Build OK"

# Run Go's static analyzer
go vet ./... && echo "✅ Vet OK"

# Format code (required before committing)
go fmt ./...
```

## Testing Changes Safely

Use dry-run mode to test the bot logic without touching git at all:

```bash
# Preview daily bot behavior (zero git calls)
go run ./cmd/daily/ -dry-run

# Preview spam bot behavior
go run ./cmd/spam/ -count 10 -dry-run

# Check all available flags
go run ./cmd/daily/ -help
go run ./cmd/spam/ -help
```

## Pull Requests

- Keep PRs small and focused on a single issue/feature.
- Ensure your code passes `go build ./...` and `go vet ./...`.
- Follow standard Go formatting (`go fmt ./...`).
- Write descriptive commit messages with emojis (we like consistency 😄).

## Bug Reports and Feature Requests

Please use the [GitHub Issue tracker](https://github.com/shvmpk/github-contributor-bot/issues) to report bugs or request features. Include:
- Your OS and Go version (`go version`)
- Steps to reproduce the issue
- Expected vs actual behavior

## 🏷️ Creating a Release (Maintainers Only)

The release workflow automatically cross-compiles and attaches binaries for 5 platforms whenever a version tag is pushed. **Do not build or upload manually.**

```bash
# 1. Commit all your changes
git add .
git commit -m "🔖 bump version to v1.0.0"
git push

# 2. Push a version tag to trigger the GitHub Action
git tag v1.0.0
git push origin v1.0.0
```

Thanks for contributing! 🚀
