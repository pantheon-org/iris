# Branch ADR: internal-layout

## Meta
- **Branch**: `main`
- **Type**: chore
- **Status**: accepted
- **Created**: 2026-05-01
- **Author**: Claude
- **PR**: <!-- N/A — decided at project initialisation -->

## Problem Statement
### Context
Go projects can place logic anywhere, but the visibility rules of the `internal/` package enforce that nothing outside the module can import it. The entry point location also determines how the binary is named and how cross-compilation targets are expressed.

### Goals
- All business logic isolated and unexportable to external modules
- Thin entry point that is trivially testable in isolation
- Clear convention for adding future packages

### Non-Goals
- Public library surface (iris is a CLI-only binary)
- Plugin / extension API

## Decision Record
### Options Considered

**`internal/` + thin `cmd/iris/main.go`**
- Standard Go project layout for CLI tools.
- `internal/` prevents accidental external imports.
- `cmd/iris/` makes the binary name explicit and supports multiple binaries if needed later.

**Flat `main.go` at repo root**
- Simpler for very small tools.
- Cons: logic and wiring mixed; harder to test; breaks convention once the project grows.

**`pkg/` for shared logic**
- Common in libraries.
- Cons: `pkg/` implies exported API — wrong signal for a CLI tool.

### Chosen Solution
- All logic in `internal/`
- `cmd/iris/main.go` only wires Cobra commands and calls `os.Exit`

### Rationale
This is the idiomatic Go layout for CLI tools with no public API. The CLAUDE.md project rule codifies it: "All logic in `internal/`; `cmd/iris/main.go` only wires cobra commands." Using `cmd/iris/` leaves room for a future `cmd/iris-server/` or similar without restructuring.

## Implementation
### Key Changes
- `cmd/iris/main.go` — entry point, Cobra wiring only
- `internal/version/version.go` — first internal package; version injected via ldflags at build time

### Testing Strategy
- `internal/` packages tested directly; `cmd/` not unit-tested (it has no logic)
- `mise run test` with `-race` flag

## Challenges & Solutions
None — enforced by Go's compiler via `internal/` visibility rules.

## Impact Assessment
- **Performance**: None.
- **Security**: Positive — `internal/` prevents inadvertent public API surface.
- **Maintenance**: Positive — clear convention reduces cognitive load when adding packages.

## Outcome & Lessons
<!-- Fill in post-merge -->
