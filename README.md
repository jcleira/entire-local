# entire-local

[![CI](https://github.com/jcleira/entire-local/actions/workflows/ci.yml/badge.svg)](https://github.com/jcleira/entire-local/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/jcleira/entire-local)](https://goreportcard.com/report/github.com/jcleira/entire-local)
[![Go Reference](https://pkg.go.dev/badge/github.com/jcleira/entire-local.svg)](https://pkg.go.dev/github.com/jcleira/entire-local)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

A terminal UI for browsing [entire.io](https://entire.io) checkpoints locally. Zero network, zero auth -- just git.

![demo](https://github.com/user-attachments/assets/0332285d-f32f-4ea9-afcb-ede757b905d9)

## Why

entire.io stores checkpoint data in a local git orphan branch. The web dashboard requires granting read access to your GitHub repos. **entire-local** is a local-first alternative that reads that data directly from your repo, keeping everything offline and private.

## Prerequisites

- **Go 1.24+**
- **git**
- A repository with [entire.io](https://entire.io) enabled (the orphan branch must exist locally)

## Installation

```bash
go install github.com/jcleira/entire-local@latest
```

Or build from source:

```bash
git clone https://github.com/jcleira/entire-local.git
cd entire-local
make install
```

## Usage

Navigate to any git repository with entire.io enabled and run:

```bash
entire-local
```

### Screens

- **Overview** -- Aggregated stats across all checkpoints: total sessions, token usage, and duration breakdown.
- **Checkpoint List** -- Browse and filter all checkpoints by prompt text. Press `/` to search.
- **Checkpoint Detail** -- Drill into a single checkpoint with three tabs:
  - **Transcript** -- Full conversation with syntax-highlighted code blocks.
  - **Files** -- Diff view showing all file changes with syntax highlighting.
  - **Plan** -- Rendered markdown plan for the session.

### Keyboard Shortcuts

| Key | Action |
|---|---|
| `j/k` | Navigate / scroll |
| `Enter/l` | Select checkpoint |
| `Esc/h` | Go back |
| `/` | Filter checkpoints |
| `Tab` | Next tab (detail view) |
| `Shift+Tab` | Previous tab |
| `1-3` | Jump to tab |
| `r` | Refresh data |
| `?` | Toggle help |
| `q` | Quit |

## Project Layout

```
cmd/                          Cobra command (root)
pkg/git/
  git.go                      ExecGit, RepoRoot, BranchExists
  reader.go                   Reader (ls-tree, show)
pkg/checkpoint/
  types.go                    Checkpoint & session types
  loader.go                   Loader (walks orphan branch tree)
  jsonl.go                    JSONL parser, flexTimestamp
  stats.go                    Token/duration stats
pkg/ui/
  commands/
    messages.go               PrintError/PrintInfo helpers
  dashboard/
    dashboard.go              Bubble Tea model, Init, Update, View
    keys.go                   Key bindings
    messages.go               Msg types
    styles.go                 Lipgloss styles
    overview.go               Overview screen
    checkpoint_list.go        Checkpoint list screen
    checkpoint_detail.go      Checkpoint detail (Files, Plan tabs)
    transcript.go             Transcript tab
    markdown.go               Markdown renderer
    help.go                   Help screen
```

## Development

```bash
make build     # Build binary
make test      # Run tests
make lint      # Run golangci-lint
make check     # fmt + vet + lint
make deps      # Download & tidy dependencies
```

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

[MIT](LICENSE)
