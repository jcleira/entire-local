# Contributing to entire-local

Thanks for your interest in contributing! Here's how to get started.

## Getting Started

1. Fork the repository
2. Clone your fork and create a feature branch:
   ```bash
   git checkout -b feature/my-feature
   ```
3. Make your changes
4. Verify everything passes:
   ```bash
   make check && make test
   ```
5. Commit with a clear message and open a Pull Request

## Code Style

- Format with `gofmt` / `goimports`
- Pass `go vet` and `golangci-lint` (`make check` runs all three)
- Use standard library `testing` only -- no testify or other test frameworks
- Table-driven tests with `t.TempDir()` for any file-system work

## Guidelines

- No comments inside function bodies
- Keep functions focused and small
- Prefer returning errors over panicking
- Follow existing patterns in the codebase

## Pull Requests

- One concern per PR -- don't mix features with refactors
- Include tests for new logic
- Update documentation if behavior changes
- Ensure CI is green before requesting review
