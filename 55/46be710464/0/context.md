# Session Context

## User Prompts

### Prompt 1

Implement the following plan:

# Plan: Selected Row Style Fix

## Context
The selected checkpoint row currently has a thick left border "tab" indicator. User wants to remove that border and use a background color instead.

## Changes — `pkg/ui/dashboard/styles.go`

**`selectedRowStyle`** (line 88): Remove `BorderLeft`, `BorderStyle`, `BorderForeground`. Add `Background(lipgloss.Color("236"))` (dark gray). Set `PaddingLeft(3)` to match `normalRowStyle`.

## Verification
- `make build`
- Visual...

