# Iris

Go rewrite of gustavodiasdev/mcpx-cli — CLI to sync MCP server configs across AI providers.

## Commands

- `mise run build` — build binary to `dist/iris`
- `mise run test`  — run tests with race detector
- `mise run lint`  — run golangci-lint

## Package layout

```text
cmd/iris/main.go          # cobra root + subcommand wiring only
internal/
  ierrors/                # sentinel errors (ErrServerNotFound, ErrMalformedConfig, etc.)
  types/                  # canonical MCPServer + IrisConfig structs
  config/                 # Codec interface (json/yaml/toml) + Store (load/save .iris.*)
  providers/              # Provider interface + Registry + Claude/Gemini/OpenCode/Codex impls
  detector/               # Detect() — scans project root for present provider configs
  merger/                 # SyncProvider / SyncAllProviders — thin orchestrator
  version/                # Version string injected at build time via ldflags
  wizard/                 # Runner interface + ScriptedRunner + BubbleteaRunner + RunInit
  cli/                    # RunList, RunStatus, RunAdd, RunRemove, RunSync, RunInitNonInteractive
  integration/            # end-to-end tests (full pipeline, no mocks)
```

## Rules

- Always work on a feature branch; never commit directly to main.
- Write tests before implementation (TDD). Test naming: `TestXxx_<scenario>_<expected>`.
- Run `go mod tidy` after adding/removing dependencies.
- Use `errors.Is` / `errors.As` for error checks — never string match on error messages.
- All logic in `internal/`; `cmd/iris/main.go` only wires cobra commands.
- Wrap all errors from external packages: `fmt.Errorf("context: %w", err)`.
- Global providers (Gemini, Codex) expose `NewXxxProviderWithPath(path string)` constructors for test isolation — use them instead of mutating `HOME`.
- For `gh` CLI commands, always prefix with `dotenvx run --` to load `GH_TOKEN`.

<!-- code-review-graph MCP tools -->
## MCP Tools: code-review-graph

**IMPORTANT: This project has a knowledge graph. ALWAYS use the
code-review-graph MCP tools BEFORE using Grep/Glob/Read to explore
the codebase.** The graph is faster, cheaper (fewer tokens), and gives
you structural context (callers, dependents, test coverage) that file
scanning cannot.

### When to use graph tools FIRST

- **Exploring code**: `semantic_search_nodes` or `query_graph` instead of Grep
- **Understanding impact**: `get_impact_radius` instead of manually tracing imports
- **Code review**: `detect_changes` + `get_review_context` instead of reading entire files
- **Finding relationships**: `query_graph` with callers_of/callees_of/imports_of/tests_for
- **Architecture questions**: `get_architecture_overview` + `list_communities`

Fall back to Grep/Glob/Read **only** when the graph doesn't cover what you need.

### Key Tools

| Tool | Use when |
|------|----------|
| `detect_changes` | Reviewing code changes — gives risk-scored analysis |
| `get_review_context` | Need source snippets for review — token-efficient |
| `get_impact_radius` | Understanding blast radius of a change |
| `get_affected_flows` | Finding which execution paths are impacted |
| `query_graph` | Tracing callers, callees, imports, tests, dependencies |
| `semantic_search_nodes` | Finding functions/classes by name or keyword |
| `get_architecture_overview` | Understanding high-level codebase structure |
| `refactor_tool` | Planning renames, finding dead code |

### Workflow

1. The graph auto-updates on file changes (via hooks).
2. Use `detect_changes` for code review.
3. Use `get_affected_flows` to understand impact.
4. Use `query_graph` pattern="tests_for" to check coverage.

@RTK.md
