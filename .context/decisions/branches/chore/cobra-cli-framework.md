# Branch ADR: cobra-cli-framework

## Meta
- **Branch**: `main`
- **Type**: chore
- **Status**: accepted
- **Created**: 2026-05-01
- **Author**: Claude
- **PR**: <!-- N/A — decided at project initialisation -->

## Problem Statement
### Context
Iris is a CLI tool; a framework is needed to handle command routing, flag parsing, help text, and shell completion generation. Go's standard `flag` package is minimal and would require significant boilerplate.

### Goals
- Standard flag/subcommand parsing with minimal boilerplate
- Auto-generated `--help` and shell completions
- Well-maintained, widely adopted in the Go ecosystem

### Non-Goals
- TUI / interactive prompts
- Plugin architecture for commands

## Decision Record
### Options Considered

**`spf13/cobra`**
- The de-facto standard Go CLI framework; used by kubectl, Hugo, GitHub CLI.
- Pros: subcommand tree, persistent flags, shell completion, active maintenance.
- Cons: heavier than `flag`; pulls in `spf13/pflag` as indirect dep.

**`urfave/cli`**
- Lighter alternative, good ergonomics.
- Cons: smaller community; less alignment with existing Go CLI conventions the team knows.

**`flag` (stdlib)**
- Zero dependencies.
- Cons: no subcommand support, no completion, significant boilerplate for anything beyond trivial.

### Chosen Solution
**`spf13/cobra` v1.9.1**

### Rationale
Cobra is the industry standard for Go CLIs. The project will grow to have multiple subcommands (one per provider sync operation), making Cobra's command-tree model directly applicable. The indirect dependency cost (`pflag`) is negligible.

## Implementation
### Key Changes
- `go.mod`: `require github.com/spf13/cobra v1.9.1`
- `cmd/iris/main.go`: wires the Cobra root command only — no business logic

### Testing Strategy
- Command wiring tested via `cobra.Command.Execute()` in integration tests
- `mise run test` covers the full suite

## Challenges & Solutions
None at this stage — Cobra initialisation is straightforward.

## Impact Assessment
- **Performance**: None.
- **Security**: None.
- **Maintenance**: Positive — Cobra is actively maintained and well-documented.

## Outcome & Lessons
<!-- Fill in post-merge -->
