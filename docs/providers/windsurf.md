# Windsurf

**iris provider name:** `windsurf`

## Config file

| Scope   | Path                                   |
|---------|----------------------------------------|
| Project | —                                      |
| Global  | `~/.codeium/windsurf/mcp_config.json`  |

Project-level config is not supported.

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

Root key: `mcpServers` (standard MCP JSON format).

## References

- [Windsurf MCP documentation](https://docs.windsurf.com/windsurf/mcp)
