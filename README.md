# 🤖 Github Contribution Bot

![Open Source Love svg3](https://badges.frapsoft.com/os/v3/open-source.svg?v=103) [![GPL-3.0 license](https://img.shields.io/badge/License-GPL--3.0-blue.svg)](https://github.com/shvmpk/github-contributor-bot/blob/main/LICENSE) ![Maintenance](https://img.shields.io/badge/Maintained%3F-yes-green.svg) [![Go](https://img.shields.io/badge/Language-Go-00ADD8.svg)](https://go.dev/) [![Version](https://img.shields.io/badge/Version-1.0.0-brightgreen.svg)](https://github.com/shvmpk/github-contributor-bot/releases)

## 🌍 Overview

😓🚶‍♂️💨 Maintaining a consistent GitHub commit streak can be a real challenge for developers! Busy schedules, unforeseen events, or simply forgetting to commit daily can break that coveted streak.

💡🔍 This can be particularly frustrating for those who want to showcase their dedication and progress. The pressure to commit daily often takes away from focusing on meaningful, quality work.

🤖✨ To solve this, I've created a **Go**-based bot that automates the process of making commits with human-like behavior — completely indistinguishable from a real developer.

🙏🔍 **Note:** Using this bot to gain an unfair advantage over others is not good practice. This project is for educational purposes only. Please use it responsibly and at your own risk.

---

## 🛠️ Why Go? — Language Comparison

This project was rewritten from Node.js to Go for extreme performance, zero dependencies, and simplicity.

| Criteria | JavaScript (Old & Most Competitors) | **Go (Current)** 🏆 |
|---|---|---|
| **Startup time** | ~150–300ms (Node.js runtime) | **~5–10ms (native binary)** |
| **Dependencies** | Massive `node_modules` | **0 external deps** |
| **Binary distribution** | Requires Node.js + npm install | **Single static binary** |
| **Memory usage** | ~30–50MB (V8 heap) | **~5–10MB** |
| **Error handling** | try/catch (silent failures) | **Explicit, forced error returns** |

---

## 🔥 Why This Bot is Better Than The Rest?

If you've looked around GitHub, you've probably seen other popular activity generators. Most of them suffer from the same problem: **they look like bots.** They commit at exactly midnight, never take days off, and use generic commit messages like *"Commit 123"*.

This bot is designed to be **indistinguishable from a real human developer**:

- 💤 **Rest Days & Work-Life Balance:** 15% chance to completely skip a day. You can also pass `--skip-weekends` to ignore Saturdays and Sundays entirely (the "Corporate Developer" look).
- 🏃‍♂️ **Sprint Days:** 10% chance to go into overdrive. On sprint days the commit count is always `max-commits + 1 to 5 extra`, mathematically guaranteeing it exceeds any normal day.
- 🕒 **Randomized Timestamps:** The bot backdates commits to random times between 9 AM and 10 PM, even though the Action runs at midnight. Your graph fills at human hours.
- 💬 **~200 Realistic Commit Messages:** Across 16 categories — bug fixes, features, refactoring, DevOps, security, database, analytics, WIP, and more. No more *"Commit 1 on July 2"*.
- 🔀 **~200 Realistic Branch Names:** 20% chance the bot creates a real-looking branch (`bugfix/resolve-cache-invalidation`, `feature/add-user-authentication`, etc.), commits to it, merges it, and deletes it — exactly like a real developer would.
- 🗂️ **Multi-File Modifications:** Instead of touching just one file, the bot randomly modifies different realistic files (`bot-data/logs/debug.log`, `bot-data/docs/api_notes.md`, `bot-data/scripts/build.sh`, etc.).
- 🔄 **Retry Logic:** If `git push` fails due to a network hiccup, the bot automatically retries up to 3 times with a 5-second backoff. Your streak is protected even on flaky connections.
- ⚡ **Zero Setup & Dual Mode:** Unlike competitors that require Node.js/Python, this is a single native Go binary. It supports both *Daily Maintenance* (via Actions) and *Bulk Backdating* (via CLI).

### Before 😓
<img width="1040" height="240" alt="before" src="https://github.com/user-attachments/assets/6743e3c5-623d-4c58-b23d-8e09e7af731b" />


### After 💪🔥
<img width="1036" height="236" alt="after" src="https://github.com/user-attachments/assets/ea90e2d0-6da3-419c-8594-1c110639f4eb" />

---

## 📁 Project Structure

```
github-contribution-bot/
├── .github/
│   └── workflows/
│       ├── daily-commit.yml        # Daily cron job (GitHub Actions)
│       └── release.yml             # Cross-platform binary builder
├── bot-data/                       # All bot-generated files (committed to git)
│   ├── data.json
│   ├── docs/api_notes.md
│   ├── logs/debug.log
│   ├── config/settings.json
│   └── scripts/build.sh
├── cmd/
│   ├── daily/main.go               # Daily commit CLI entrypoint
│   └── spam/main.go                # Spam commit CLI entrypoint
├── internal/
│   ├── commit/commit.go            # ~200 messages, ~200 branch names, file I/O
│   ├── config/config.go            # CLI flags + bot.config.json loader
│   ├── git/git.go                  # Git ops with retry logic and dry-run
│   └── stats/stats.go              # Persistent run stats tracking
├── bot.config.json                 # User configuration file
├── messages.txt                    # Custom commit messages (optional)
├── branch_names.txt                # Custom branch names (optional)
├── stats.json                      # Auto-generated run stats (gitignored)
├── go.mod
├── LICENSE
├── README.md
└── CONTRIBUTING.md
```

---

## 🚀 How it Works & Usage

This project provides two distinct CLI tools. Ensure you have [Go 1.22+](https://go.dev/) installed, **or** download a pre-built binary from the [Releases page](https://github.com/shvmpk/github-contributor-bot/releases).

### 📦 Download Pre-Built Binaries (No Go Required)

Go to the [Releases page](https://github.com/shvmpk/github-contributor-bot/releases) and download the binary for your platform:

| Platform | Binary |
|---|---|
| Linux (x64) | `spam-linux-amd64`, `daily-linux-amd64` |
| Linux (ARM) | `spam-linux-arm64`, `daily-linux-arm64` |
| macOS (Intel) | `spam-darwin-amd64`, `daily-darwin-amd64` |
| macOS (Apple Silicon) | `spam-darwin-arm64`, `daily-darwin-arm64` |
| Windows | `spam-windows-amd64.exe`, `daily-windows-amd64.exe` |

---

### 1. Daily Commit (GitHub Actions)

Designed to commit 1–3 times per day at a randomized hour via GitHub Actions. Fully automatic after one-time setup.

**Setup:**
1. Fork this repository.
2. Navigate to GitHub Settings → Developer settings → Personal access tokens.
3. Generate a new token (classic), set expiry to never, add **`repo`** and **`workflow`** scopes.
4. Add it to your repo secrets as `GH_TOKEN`.

> ✅ **That's it!** The bot runs **fully automatically every day at 11:55 PM UTC** — no computer needs to be on. You can also trigger it manually from the **Actions** tab using the "Run workflow" button.

You can also run it locally:
```bash
go run ./cmd/daily/ -min-commits 1 -max-commits 3

# Optional: skip weekends for the corporate dev look
go run ./cmd/daily/ -skip-weekends

# Preview what the bot WOULD do without committing anything
go run ./cmd/daily/ -dry-run
```

> **💡 How the logic works:** `-min-commits` and `-max-commits` set your normal baseline. The bot automatically overrides this 25% of the time:
> - **Rest Days (15%):** 0 commits.
> - **Sprint Days (10%):** `max-commits + 1 to 5 extra commits`. A sprint day always exceeds your highest normal day, regardless of your flags.

---

### 2. Spam Commit (Local Execution)

Designed to create a series of backdated commits to fill gaps in your contribution graph.

```bash
# Clone the repository
git clone https://github.com/shvmpk/github-contributor-bot.git
cd github-contribution-bot

# Run with defaults (100 commits over 54 weeks)
go run ./cmd/spam/ -count 100 -weeks-back 54

# Preview without touching git
go run ./cmd/spam/ -count 50 -dry-run
```

**⚠️ Warning:** This alters your git history. Use responsibly!

---

## ⚙️ Configuration (`bot.config.json`)

Instead of passing CLI flags every time, you can configure the bot once in `bot.config.json`. CLI flags always override file values.

```json
{
  "min_commits": 1,
  "max_commits": 3,
  "skip_weekends": false,
  "dry_run": false,
  "data_dir": "bot-data",
  "messages_file": "messages.txt",
  "messages_mode": "append",
  "branch_names_file": "branch_names.txt",
  "branch_mode": "append"
}
```

| Key | Type | Default | Description |
|---|---|---|---|
| `min_commits` | int | `1` | Minimum commits on a normal day |
| `max_commits` | int | `3` | Maximum commits on a normal day |
| `skip_weekends` | bool | `false` | Skip Saturday and Sunday |
| `dry_run` | bool | `false` | Preview without touching git |
| `data_dir` | string | `"bot-data"` | Directory for all auto-generated files |
| `messages_file` | string | `"messages.txt"` | Path to custom commit messages |
| `messages_mode` | string | `"append"` | `"append"` or `"replace"` |
| `branch_names_file` | string | `"branch_names.txt"` | Path to custom branch names |
| `branch_mode` | string | `"append"` | `"append"` or `"replace"` |

---

## 💬 Custom Messages & Branch Names

### `messages.txt` — Custom Commit Messages

Create or edit `messages.txt` with one message per line. Lines starting with `#` are ignored.

```
# My custom messages
🎯 implement core feature logic
🔍 investigate performance bottleneck
💡 prototype new approach for data pipeline
```

Set `"messages_mode"` in `bot.config.json`:
- `"append"` *(default)* — your messages are **added** to the ~200 built-in ones (maximum variety).
- `"replace"` — **only** your messages are used (full control).

### `branch_names.txt` — Custom Branch Names

Same format as `messages.txt`. Controls the branch names used during the 20% branching runs.

```
# My project-specific branches
feature/my-auth-module
bugfix/fix-dashboard-crash
```

Set `"branch_mode"` in `bot.config.json` to `"append"` or `"replace"`.

---

## 🧪 Dry Run Mode

Before running the bot for real, you can preview exactly what it would do:

```bash
go run ./cmd/daily/ -dry-run
```

Example output:
```
🤖 Daily Commit Bot v1.0.0
🧪 DRY RUN mode — no git commands will be executed
📝 Making 2 commit(s) today

── Commit 1/2 ─────────────────────────
   Time: 2026-07-02T14:23:00+05:30
[DRY RUN] git add bot-data/logs/debug.log
[DRY RUN] git commit -m "🔧 update CI pipeline for faster builds" --date 2026-07-02T14:23:00+05:30
   ✅ 🔧 update CI pipeline for faster builds
...
```

No files are modified and no git commands are executed.

---

## 📊 Run Statistics

After every run, the bot prints a summary and saves cumulative stats to `stats.json` (gitignored):

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📊 Run Summary
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📝 Today:           2 commit(s)
🔀 Branch used:     bugfix/resolve-cache-invalidation
🗂️  Files modified:  [bot-data/logs/debug.log bot-data/docs/api_notes.md]

📈 All-Time Stats
   Total runs:     47
   Total commits:  143
   Rest days:      8
   Sprint days:    5
   Branches used:  11
   This week:      9 commits
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

---

---

## 📜 Contribution

Contributions are welcome! Please refer to [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## ⚖️ License

This project is licensed under the **GPL-3.0 License** — see the [LICENSE](LICENSE) file for details.
