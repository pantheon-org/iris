# VS Code GitHub Copilot

**iris provider name:** `copilot`

## Config file

| Scope   | Path                                        |
|---------|---------------------------------------------|
| Project | `.vscode/mcp.json`                          |
| Global  | VS Code user `settings.json` → `mcp.servers` |

iris targets the project-level file. The global form embeds server definitions inside VS Code's main user settings file rather than a dedicated MCP config.

## Format

```json
{
  "servers": {
    "server-name": {
      "type": "stdio",
      "command": "npx",
      "args": ["-y", "@modelcontextprotocol/server-filesystem", "/tmp"],
      "env": { "KEY": "value" }
    }
  }
}
```

Root key: `servers` (not `mcpServers`). Each entry has a sibling `"type"` field (`"stdio"`, `"sse"`, or `"http"`). iris preserves all existing top-level keys (e.g. `"inputs"`).

## Supported server fields

| Field     | Type               | Notes                                                                                    |
|-----------|--------------------|------------------------------------------------------------------------------------------|
| `type`    | string             | `"stdio"`, `"sse"`, or `"http"` (Streamable HTTP; tries HTTP, falls back to SSE)        |
| `command` | string             | Executable path (stdio only)                                                             |
| `args`    | array of strings   | CLI arguments (stdio only)                                                               |
| `env`     | map of strings     | Environment variables                                                                    |
| `url`     | string             | Server URL (sse/http only)                                                               |

## Fields NOT supported

The following iris `MCPServer` fields have no equivalent in the VS Code Copilot format and are **silently omitted** when iris generates `.vscode/mcp.json`:

- `headers` — Copilot does not accept per-server HTTP headers
- `cwd` — working-directory override is not supported
- `enabled` — Copilot has no per-server enable/disable flag

This is intentional and correct behaviour. The omission is achieved via struct mapping (`vscodeServer` in `internal/providers/vscode_copilot.go`) — no explicit stripping code is needed. Confirmed by source inspection (2026-05-05, wave 0b audit).

## References

- [VS Code MCP documentation](https://code.visualstudio.com/docs/copilot/chat/mcp-servers)
