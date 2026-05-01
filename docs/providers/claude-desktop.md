# Claude Desktop App

**iris provider name:** `claude-desktop`

## Config file

| OS      | Path                                                                      |
|---------|---------------------------------------------------------------------------|
| macOS   | `~/Library/Application Support/Claude/claude_desktop_config.json`         |
| Windows | `%APPDATA%\Claude\claude_desktop_config.json`                             |
| Linux   | `~/.config/Claude/claude_desktop_config.json`                             |

Global only — no per-project config.

## Format

```json
{
  "mcpServers": {
    "server-name": {
      "command": "npx",
      "args": ["-y", "@modelcontextprotocol/server-filesystem", "/tmp"],
      "env": { "KEY": "value" }
    }
  }
}
```

Root key: `mcpServers` (same format as Claude Code).

## References

- [Claude Desktop MCP quickstart](https://modelcontextprotocol.io/quickstart/user)
