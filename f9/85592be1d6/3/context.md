# Session Context

## User Prompts

### Prompt 1

Implement the following plan:

# Plan: Fix scrolling in detail view tabs

## Context
After adding full diff rendering to the Files tab, j/k scrolling doesn't work. Root cause: `m.maxScroll` (the model struct field) is **never written to** — it stays 0. The Down handler in `handleDetailKeys` (`dashboard.go:203`) gates on `m.scrollOffset < m.maxScroll`, so scrollOffset can never increase past 0.

The View functions compute a local `maxScroll` variable for clamping display, but that never flows ...

