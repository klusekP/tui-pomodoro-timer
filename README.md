# ⏱ time-tracker

> A small, friendly **terminal time tracker** with a Stopwatch, Pomodoro and a
> searchable session history. Built in Go with [Bubble Tea].

[![CI](https://github.com/klusekP/time_pr_tracking/actions/workflows/ci.yml/badge.svg)](https://github.com/klusekP/time_pr_tracking/actions/workflows/ci.yml)
[![Release](https://github.com/klusekP/time_pr_tracking/actions/workflows/release.yml/badge.svg)](https://github.com/klusekP/time_pr_tracking/actions/workflows/release.yml)
[![Latest release](https://img.shields.io/github/v/release/klusekP/time_pr_tracking?display_name=tag&sort=semver)](https://github.com/klusekP/time_pr_tracking/releases/latest)
[![Go](https://img.shields.io/badge/Go-1.25%2B-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![Made with Charm](https://img.shields.io/badge/made%20with-charm-ff69b4)](https://charm.sh/)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

`time-tracker` is a single-binary, full-screen TUI for tracking the time you
spend on tasks. It pairs a no-nonsense **stopwatch** that persists sessions
locally with a focused **Pomodoro** timer and a **history** view that lets you
filter by task name to see, for example, _"how many hours did I spend on the
frontend this month?"_.

---

## Table of contents

- [Features](#features)
- [Installation](#installation)
  - [Pre-built binaries](#pre-built-binaries)
  - [`go install`](#go-install)
  - [From source](#from-source)
- [Usage](#usage)
  - [Tabs and global keys](#tabs-and-global-keys)
  - [Stopwatch](#stopwatch)
  - [Pomodoro](#pomodoro)
  - [Sessions (history)](#sessions-history)
- [Data & configuration](#data--configuration)
  - [Database location](#database-location)
  - [Schema](#schema)
- [Architecture](#architecture)
  - [Project layout](#project-layout)
  - [Dependency graph](#dependency-graph)
- [Development](#development)
  - [Build & run](#build--run)
  - [Testing](#testing)
  - [Formatting & linting](#formatting--linting)
- [Roadmap](#roadmap)
- [Contributing](#contributing)
- [License](#license)
- [Acknowledgements](#acknowledgements)

---

## Features

- 🎯 **Three modes, one binary** — Stopwatch, Pomodoro and a Sessions
  browser, switchable with `tab` or `1` / `2` / `3`.
- ⏱ **Stopwatch** with a large `HH:MM:SS` clock, start/pause/reset and
  one-key **save** to a local SQLite database with the task name attached.
- 🍅 **Pomodoro** with a progress bar, presets (25 / 5 / 15 min), `±5 min`
  fine-tuning and an OS-aware **alarm** when the countdown finishes.
- 📚 **Sessions** view — full history with totals (today, all loaded, per
  mode) and a **substring filter** that recomputes totals live, so you can
  answer _"how much did I spend on `<project>`?"_ in one keystroke.
- 💾 **Local-first**, single-file SQLite database — no servers, no accounts,
  no telemetry.
- 🚫 **No CGO** — uses [`modernc.org/sqlite`][sqlite-pure-go]; builds with
  plain `go build` on every platform.
- 🧪 **Solid test coverage** for the core domain (`domain` 100 %, `timer`
  100 %, `repository` ~70 %, mode logic ~50 %).

[sqlite-pure-go]: https://pkg.go.dev/modernc.org/sqlite

## Installation

### Pre-built binaries

Each tagged release ships ready-to-run binaries on the
[**Releases**](https://github.com/klusekP/time_pr_tracking/releases) page for:

- Linux — `amd64`, `arm64`
- macOS — `amd64` (Intel), `arm64` (Apple Silicon)
- Windows — `amd64`

Download the archive matching your platform, extract it and put the binary on
your `PATH`. Each release also includes a `SHA256SUMS` file you can use to
verify the download.

```bash
# Linux / macOS example
curl -sSLO https://github.com/klusekP/time_pr_tracking/releases/latest/download/time-tracker-<version>-<os>-<arch>.tar.gz
tar -xzf time-tracker-*.tar.gz
sudo mv time-tracker-*/time-tracker /usr/local/bin/
time-tracker
```

### `go install`

If you have a Go toolchain handy:

```bash
go install github.com/klusekP/time_pr_tracking/cmd/time-tracker@latest
time-tracker
```

This drops a `time-tracker` binary into `$(go env GOBIN)` (or `$GOPATH/bin`).
Add that directory to your `PATH` if it isn't already.

### From source

```bash
git clone https://github.com/klusekP/time_pr_tracking.git
cd time_pr_tracking
go run ./cmd/time-tracker
```

…or build a binary you can move around:

```bash
go build -o time-tracker ./cmd/time-tracker
./time-tracker
```

**Requirements**

- Go **1.25 or newer** (the project tracks the latest Go release).
- A modern terminal with True Color support — iTerm2, WezTerm, Alacritty,
  Kitty, Windows Terminal, GNOME Terminal, etc.

The first run creates the SQLite database under `~/.time_tracker/`
(see [Data & configuration](#data--configuration)).

## Usage

### Tabs and global keys

| Key      | Action                                          |
| -------- | ----------------------------------------------- |
| `tab`    | Cycle through tabs (Stopwatch → Pomodoro → Sessions) |
| `1` / `2` / `3` | Jump straight to a tab                   |
| `e` / `i` | Edit the **Task** field (Stopwatch only) — `Enter` / `Esc` to confirm |
| `q` / `Ctrl+C` | Quit                                       |

The status pill in the header (`● RUNNING`, `❚❚ PAUSED`, `⏰ DONE!`,
`○ IDLE`) reflects the active tab, and the right side of the header shows
the total time saved **today**.

### Stopwatch

A simple count-up timer that **persists** every saved session.

| Key     | Action                                       |
| ------- | -------------------------------------------- |
| `space` | Start / pause                                |
| `s`     | Stop and **save** the session                |
| `r`     | Reset (does **not** save)                    |

The current `Task:` value is attached to the saved session. If left empty,
the session is stored as `Untitled`.

### Pomodoro

A countdown with an alarm. Pomodoro is **not persisted** — it's purely a
focus aid.

| Key       | Action                                          |
| --------- | ----------------------------------------------- |
| `space`   | Start / pause (after the alarm: start a new one) |
| `s`       | Stop early                                      |
| `r`       | Reset                                           |
| `w`       | Preset: 25 min (work)                           |
| `b`       | Preset: 5 min (short break)                     |
| `l`       | Preset: 15 min (long break)                     |
| `+` / `-` | Adjust target by ± 5 min (clamped to 1m–8h)    |

When the countdown reaches zero the app:

1. Rings the terminal bell (`BEL`) and plays an OS-appropriate sound:
   - **macOS** — `afplay /System/Library/Sounds/Glass.aiff`
   - **Linux** — `paplay` / `aplay` (when available)
   - **Windows** — `powershell [console]::beep`
2. Shows a blinking `⏰  POMODORO FINISHED!  ⏰` banner.

### Sessions (history)

A read-only browser over your saved Stopwatch sessions, with a built-in
filter and live summaries.

| Key        | Action                                        |
| ---------- | --------------------------------------------- |
| `f` / `/`  | Edit the filter                               |
| `Enter` / `Esc` | Apply filter and exit edit mode          |
| `c`        | Clear the filter                              |

The summary updates as you change the filter and shows:

- with **no filter**: today's total, total of loaded sessions, breakdown per
  mode (Stopwatch / Pomodoro);
- with a **filter**: number of matches and total hours, plus per-mode split
  — useful for answering _"how many hours did I spend on the **frontend**?"_.

Matching is **case-insensitive substring** on the task name.

## Data & configuration

### Database location

The SQLite database lives at:

```
~/.time_tracker/sessions.db
```

The directory is created on first run. You can inspect the file with any
SQLite client:

```bash
sqlite3 ~/.time_tracker/sessions.db "SELECT * FROM sessions ORDER BY ended_at DESC LIMIT 10;"
```

### Schema

```sql
CREATE TABLE sessions (
    id           INTEGER  PRIMARY KEY AUTOINCREMENT,
    kind         TEXT     NOT NULL DEFAULT 'stopwatch',  -- 'stopwatch' | 'pomodoro'
    task         TEXT     NOT NULL,
    duration_ms  INTEGER  NOT NULL,
    started_at   DATETIME NOT NULL,
    ended_at     DATETIME NOT NULL
);

CREATE INDEX idx_sessions_ended_at ON sessions(ended_at DESC);
```

Older databases (created before the `kind` column existed) are migrated
automatically on startup.

## Architecture

`time-tracker` is structured around small, single-purpose packages and a few
well-defined interfaces. Concrete implementations are wired together in
exactly one place — the `cmd/time-tracker` composition root — so swapping
out the storage backend, the alarm strategy or even the clock requires no
changes elsewhere.

### Project layout

```
time_pr_tracking/
├── cmd/
│   └── time-tracker/
│       └── main.go            # composition root – wires everything together
└── internal/
    ├── domain/
    │   └── session.go         # Session value object + SessionKind
    ├── clock/
    │   └── clock.go           # Clock interface + System (time.Now)
    ├── timer/
    │   └── timer.go           # Pausable counter (depends on clock.Clock)
    ├── alarm/
    │   └── alarm.go           # Alarm interface + per-OS impls + Composite
    ├── repository/
    │   └── sqlite.go          # Sessions interface + SQLite implementation
    ├── ui/
    │   ├── state.go           # RunState (Idle/Running/Paused/Finished)
    │   ├── theme.go           # Lipgloss colors and styles
    │   ├── bigfont.go         # Big-font digit renderer
    │   ├── format.go          # FormatHMS / FormatHours / DisplayTask / Clamp…
    │   ├── progress.go        # RenderProgressBar
    │   └── messages.go        # tea.Msg types + factories (Tick, Saved, …)
    ├── mode/
    │   ├── mode.go            # Mode interface + SessionsAware (ISP)
    │   ├── stopwatch.go       # Stopwatch mode (Timer + Sessions)
    │   ├── pomodoro.go        # Pomodoro mode (Timer + Alarm)
    │   └── history.go         # Read-only history with substring filter
    └── app/
        ├── model.go           # Bubble Tea Model — thin controller
        └── renderer.go        # Renderer + ViewState (pure rendering)
```

### Dependency graph

```
domain  ←  repository
domain  ←  ui  ←  mode  ←  app
clock   ←  timer ─┘        ↑
alarm  ────────────────────┘
```

- `domain`, `clock`, `alarm` — leaf packages, free of any TUI imports.
- `timer` depends only on `clock`; `repository` depends only on `domain`.
- `ui` is the presentation toolkit (theme, big font, message types, format
  helpers).
- `mode` composes whatever a given mode needs (`Timer`, `Sessions` or
  `Alarm`, `ui`) behind the `Mode` interface.
- `app` (Bubble Tea Model + Renderer) only sees `mode`, `ui` and
  `repository`.
- `cmd/time-tracker` is the only place where concrete implementations
  (`repository.NewSQLite`, `alarm.NewSystem`, the mode list) appear.

Adding, say, a “HIIT” mode is a new file implementing `mode.Mode` plus one
line in `cmd/time-tracker/main.go`. Swapping SQLite for Postgres is a new
type implementing `repository.Sessions`. Nothing else changes.

## Development

### Build & run

```bash
go run ./cmd/time-tracker            # run from source
go build ./cmd/time-tracker          # produce ./time-tracker
go build -o ./bin/time-tracker ./cmd/time-tracker
```

### Testing

```bash
go test ./...                        # run the full suite
go test -cover ./...                 # with per-package coverage
go test -run TestPomodoro ./...      # focus on a subset
```

The suite uses fake clocks and an in-memory `Sessions` fake to stay fast
and deterministic — no `time.Sleep`, no real audio playback. Repository
tests run against an isolated SQLite file inside `t.TempDir()`.

### Formatting & linting

```bash
gofmt -l .                           # show files that need formatting
gofmt -w .                           # rewrite them
go vet ./...                         # static checks
```

## Roadmap

A short, non-binding wish list:

- [ ] Daily / weekly summaries grouped by date.
- [ ] Export sessions to CSV / JSON.
- [ ] Configurable database path (env var or CLI flag).
- [ ] Scrolling in the Sessions table for very long histories.
- [ ] Persistent settings (last filter, default Pomodoro target).

## Contributing

Bug reports, ideas and pull requests are very welcome.

1. Fork the repo and create a feature branch.
2. Run `gofmt -w .`, `go vet ./...` and `go test ./...` — they should all
   pass.
3. Keep new code aligned with the existing style: small packages,
   single-purpose types, dependencies injected through interfaces, no
   comments that merely restate what the code does.
4. Open a pull request describing the change and the motivation.

For larger changes please open an issue first to discuss the design — it
saves everyone time.

## License

Released under the [MIT License](https://opensource.org/licenses/MIT).

## Acknowledgements

Built on the shoulders of the wonderful [Charm] ecosystem:

- [Bubble Tea] — the Elm-inspired TUI runtime.
- [Bubbles](https://github.com/charmbracelet/bubbles) — ready-made
  components (text input, etc.).
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) — declarative
  styling for terminal output.
- [`modernc.org/sqlite`][sqlite-pure-go] — pure-Go SQLite, so the binary
  builds without CGO on every platform.

[Charm]: https://charm.sh/
[Bubble Tea]: https://github.com/charmbracelet/bubbletea
