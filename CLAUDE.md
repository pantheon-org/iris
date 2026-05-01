# Iris

Go rewrite of gustavodiasdev/mcpx-cli — CLI to sync MCP server configs across AI providers.

## Commands

- `mise run build` — build binary to `dist/iris`
- `mise run test`  — run tests with race detector
- `mise run lint`  — run golangci-lint

## Rules

- Always work on a feature branch; never commit directly to main.
- Write tests before implementation (TDD).
- Run `go mod tidy` after adding/removing dependencies.
- Use `errors.Is` / `errors.As` for error checks — never string match on error messages.
- All logic in `internal/`; `cmd/iris/main.go` only wires cobra commands.
