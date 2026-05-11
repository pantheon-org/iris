# Cline

**iris provider name:** `cline`

## Config file

| Scope   | Path |
|---------|------|
| Project | — |
| Global  | `~/Library/Application Support/Code/User/globalStorage/saoudrizwan.claude-dev/settings/cline_mcp_settings.json` |

Project-level config is not supported.

## Format

```json
{
  "mcpServers": {
    "server-name": {
      "command": "npx",
      "args": ["-y", "@modelcontextprotocol/server-filesystem"],
      "env": { "KEY": "value" }
    }
  }
}
```

Root key: `mcpServers` (standard MCP JSON format). Both stdio (`command`/`args`/`env`) and HTTP (`url`/`headers`) transports are supported.

## References
