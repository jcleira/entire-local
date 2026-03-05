# CLAUDE.md -- entire-local

## Project Overview

Local TUI viewer for entire.io checkpoints. Reads from the entire.io checkpoint orphan branch using git plumbing commands. Read-only, no network, no auth.

**Module:** `github.com/jcleira/entire-local`
**Go version:** 1.24.0
**CLI framework:** Cobra + Bubble Tea

## Project Structure

```
cmd/                          Cobra root command -> launches TUI
pkg/git/
  git.go                      ExecGit, RepoRoot, BranchExists
  reader.go                   Reader (ls-tree, show)
pkg/checkpoint/
  types.go                    Checkpoint & session types
  loader.go                   Loader (walks orphan branch tree)
  jsonl.go                    JSONL parser, flexTimestamp
  stats.go                    Token/duration stats
pkg/ui/commands/
  messages.go                 PrintError/PrintInfo helpers
pkg/ui/dashboard/
  dashboard.go                Bubble Tea model, Init, Update, View
  keys.go                     Key bindings
  messages.go                 Msg types
  styles.go                   Lipgloss styles
  overview.go                 Overview screen
  checkpoint_list.go          Checkpoint list screen
  checkpoint_detail.go        Checkpoint detail (Files, Plan tabs)
  transcript.go               Transcript tab
  markdown.go                 Markdown renderer
  help.go                     Help screen
```

## Build & Test

```bash
make build     # Compile binary
make test      # Run tests with race detector
make lint      # golangci-lint
make check     # fmt + vet + lint
```

## Architecture

- **Git plumbing reads**: All data comes from `git ls-tree` and `git show` on the orphan branch
- **Lazy transcript loading**: full.jsonl only loaded when entering detail view
- **No config needed**: discovers repo from cwd via `git rev-parse --show-toplevel`
- **Read-only**: never writes to the repository
- **Detail view tabs**: Transcript, Files (diff with syntax highlighting), Plan (markdown rendered)
- **Scroll via estimate**: Detail tab scrolling uses `contentLineEstimate()` as an upper bound; View clamps to precise max during rendering

## Key Patterns

- No comments inside function bodies
- Standard library `testing` only
- Table-driven tests with `t.TempDir()`
- Warm orange theme (lipgloss color 208)
- Vim-style keybindings (j/k/h/l)

## Data Layout on Orphan Branch

```
<shard>/<rest>/metadata.json          Checkpoint metadata
<shard>/<rest>/0/metadata.json        Session metadata (tokens, duration)
<shard>/<rest>/0/prompt.txt           User prompt
<shard>/<rest>/0/context.md           Session context
<shard>/<rest>/0/diff.patch           Code diff
<shard>/<rest>/0/plan.md              Plan
<shard>/<rest>/0/full.jsonl           Full transcript
```
