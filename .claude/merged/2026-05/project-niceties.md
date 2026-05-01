# Branch ADR: feat/project-niceties

## Meta
- **Branch**: `feat/project-niceties`
- **Type**: feat
- **Status**: merged
- **Created**: 2026-05-01
- **Merged**: 2026-05-01
- **Author**: Claude
- **PR**: #1

## Problem Statement
### Context
A new Go project needs baseline tooling: pre-commit quality gates, markdown linting, and a unified AI agent instruction file. Without these, code quality checks only run in CI (too late) and multiple AI assistants require separate instruction files.

### Goals
- Pre-commit hooks that catch issues before they reach CI
- Markdown linting with project-appropriate rule configuration
- Single source of truth for AI agent instructions (Claude, Gemini, others)

### Non-Goals
- Full CI pipeline configuration (handled separately)
- IDE-specific configuration

## Decision Record
### Options Considered

**`pre-commit` (Python-based hook framework)**
- Widely used; large ecosystem of hooks.
- Cons: requires Python runtime; slower; less idiomatic in Go projects.

**`lefthook` (Go-based hook framework)**
- Fast, single binary, no runtime dependency beyond the binary itself.
- Supports parallel hook execution; mirrors what skill-quality-auditor uses.
- Pros: consistent with existing tooling in the team's ecosystem.

**Separate CLAUDE.md + GEMINI.md**
- Each AI assistant gets its own file.
- Cons: content duplication; any update requires editing multiple files; divergence risk.

**Single AGENTS.md with symlinks**
- One canonical file; CLAUDE.md and GEMINI.md become symlinks.
- Pros: single edit point; no divergence; extensible to additional AI tools.

### Chosen Solution
- **lefthook** for git hooks
- **AGENTS.md** as single source of truth; `CLAUDE.md` and `GEMINI.md` are symlinks

### Rationale
`lefthook` is the natural fit for a Go project — no Python dependency, fast parallel execution. The AGENTS.md symlink pattern is a pragmatic solution to multi-AI-assistant maintenance: add a new assistant by adding one symlink, not duplicating content.

## Implementation
### Key Changes
- `lefthook.yml` — pre-commit: gofmt, go vet, golangci-lint, mdlint, shellcheck; pre-push: go test -race, go build
- `mdlint.toml` — disables MD013 (line length), MD032, MD051, MD055, MD058 (false positives for this project)
- `CLAUDE.md` → symlink to `AGENTS.md`
- `GEMINI.md` → symlink to `AGENTS.md`
- `.gitignore` — excludes `.context/` from tracking
- `README.md`, `CONTRIBUTING.md` — initial documentation

### Testing Strategy
- `lefthook install` run locally to verify hooks fire correctly
- mdlint violations caught and fixed as part of this branch

## Challenges & Solutions
- **mdlint false positives**: MD013 (line length in code blocks), MD055 (pipes in backticks), and others fired on valid markdown. Disabled in `mdlint.toml` with targeted rule suppression.
- **golangci-lint config**: linter settings needed to move out of `formatters` block (fixed in same branch).

## Impact Assessment
- **Performance**: Positive — pre-commit hooks catch issues in <5s locally.
- **Security**: None.
- **Maintenance**: Positive — AGENTS.md symlink pattern reduces drift between AI tool configs.

## Outcome & Lessons
- The AGENTS.md symlink pattern works cleanly — both Claude and Gemini pick up the file correctly.
- mdlint rule tuning was necessary from day one; `mdlint.toml` should be considered a living config.
