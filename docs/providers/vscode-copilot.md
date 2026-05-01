# VS Code GitHub Copilot

**iris provider name:** `vscode-copilot`

## Config file

| Scope   | Path                |
|---------|---------------------|
| Project | `.vscode/mcp.json`  |

Global/user-scoped MCP servers can also be configured via VS Code user `settings.json` under `mcp.servers`, but iris targets the project-level file.

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

Root key: `servers` (not `mcpServers`). Each entry has a sibling `"type"` field (`"stdio"` or `"sse"`). iris preserves all existing top-level keys (e.g. `"inputs"`).

## References

- [VS Code MCP documentation](https://code.visualstudio.com/docs/copilot/chat/mcp-servers)
