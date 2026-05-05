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
  types/                  # canonical MCPServer + IrisConfig structs (IrisConfig.Providers records which providers are installed)
  config/                 # Codec interface (json/yaml/toml) + Store (load/save .iris.*)
  providers/              # Provider interface + per-provider impls + name constants
  registry/               # Registry — builds and filters the provider registry
  detector/               # Detect() — scans project root for present provider configs
  sync/                   # SyncProvider / SyncAllProviders — thin orchestrator
  i18n/                   # Internationalisation — locale loading and T() helper
  io/                     # OS helpers (UserHomeDir, etc.)
  version/                # Version string injected at build time via ldflags
  wizard/                 # Runner interface + ScriptedRunner + TerminalRunner (charmbracelet/huh) + RunInit + CollectImportCandidates + GroupImportCandidates
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

@RTK.md

@CODE_REVIEW_GRAPH.md
