# Branch ADR: main

## Meta
- **Branch**: `main`
- **Type**: chore
- **Status**: accepted
- **Created**: 2026-05-01
- **Author**: Claude
- **PR**: <!-- Add PR number when opened -->

## Problem Statement
### Context
Iris is a CLI tool to sync MCP server configurations across AI providers (a Go rewrite of gustavodiasdev/mcpx-cli). At project inception, the implementation language needed to be chosen. A key constraint surfaced: users may want to install the tool via `uvx` (UV, the Python package manager) or `npx` (NPM), which are the primary ways developers install ad-hoc CLI tools in Python and Node ecosystems.

### Goals
- Choose a language that produces a reliable, portable CLI binary
- Keep the distribution story as simple as possible
- Minimise cross-compilation complexity for CI/CD releases

### Non-Goals
- Python or Node.js bindings / library usage
- WASM target support
- Embedding inside another language runtime

## Decision Record
### Options Considered

**Rust**
- Pros: Strong npm/PyPI binary packaging story (`wasm-pack`, `cargo-dist`, platform-specific wheel packages via `maturin`); excellent binary size; memory safety without GC pauses.
- Cons: Slower compile times; steeper learning curve; more complex toolchain for cross-compilation; UV/NPM distribution requires extra wrapper tooling (`maturin`, `cargo-dist`) and non-trivial CI setup.

**Go**
- Pros: Simple, fast cross-compilation (`GOOS`/`GOARCH`); small single-binary output; straightforward CI release via `goreleaser`; less toolchain overhead.
- Cons: No native UV/NPM distribution path — requires the same wrapper layer as Rust (a PyPI or npm shim that downloads the platform binary at install time); GC (minimal impact for a CLI); slightly less ergonomic binary size optimisation.

### Chosen Solution
**Go**, with Cobra (`spf13/cobra`) as the CLI framework.

### Rationale
Neither language integrates natively with UV/NPM distribution without a wrapper layer — both require publishing platform-specific binaries and a thin install shim. Given that the distribution overhead is identical, Go's simpler cross-compilation and faster iteration cycle were the deciding factors. The codebase is a rewrite of an existing tool, so velocity matters.

Distribution via `uvx` / `npx` remains a future goal; it will be addressed separately via a packaging ADR once the core CLI is stable.

## Implementation
### Key Changes
- `cmd/iris/main.go` — Cobra root command wiring only
- `internal/version/version.go` — version injected at build time via ldflags
- `go.mod` — module `github.com/pantheon-org/iris`, Go 1.24.3, `spf13/cobra` v1.9.1

### Testing Strategy
- `mise run test` runs tests with `-race` detector
- `mise run build` verifies cross-compilation to `dist/iris`

## Challenges & Solutions
- **UV/NPM distribution**: Neither Go nor Rust solves this for free. Decision deferred — a future ADR will cover the packaging/shim strategy once the CLI feature set is stable enough to publish.

## Impact Assessment
- **Performance**: None — Go's GC has negligible impact for a short-lived CLI process.
- **Security**: None at this stage.
- **Maintenance**: Positive — Go's simpler toolchain reduces CI/CD maintenance burden.

## Outcome & Lessons
<!-- Fill in post-merge: results, metrics, lessons learned -->
