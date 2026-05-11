# Kiro

**iris provider name:** `kiro`

## Config file

| Scope   | Path |
|---------|------|
| Project | `.kiro/settings/mcp.json` |
| Global  | `~/.kiro/settings/mcp.json` |

Both project-level and global configs are supported.

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
