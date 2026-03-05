# Session Context

## User Prompts

### Prompt 1

Implement the following plan:

# Plan: Chat-Style Transcript + Markdown Rendering

## Context
The transcript tab currently shows user/assistant messages as plain text with basic left-border styling. The user wants it to feel like a real chat conversation and have markdown properly rendered (headers, code blocks, bold, lists, etc.) — especially for assistant responses which are typically markdown.

## Approach
Add `github.com/charmbracelet/glamour` for terminal markdown rendering. Redesign the...

