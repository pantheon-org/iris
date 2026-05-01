# Claude Desktop App

**iris provider name:** `claude-desktop`

## Config file

| Scope   | Path                                      |
|---------|-------------------------------------------|
| Project | —                                         |
| Global  | OS-dependent (see below)                  |

Global paths by OS:

| OS      | Path                                                                      |
|---------|---------------------------------------------------------------------------|
| macOS   | `~/Library/Application Support/Claude/claude_desktop_config.json`         |
| Windows | `%APPDATA%\Claude\claude_desktop_config.json`                             |
| Linux   | `~/.config/Claude/claude_desktop_config.json`                             |

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
- [Anthropic Claude Desktop configuration docs](https://docs.anthropic.com/en/docs/claude-code/mcp)
