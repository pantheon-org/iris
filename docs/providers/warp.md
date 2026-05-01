# Warp

**iris provider name:** `warp`

## Config file

| Scope  | Path              |
|--------|-------------------|
| Global | `~/.warp/mcp.json` |

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

- [Warp MCP documentation](https://docs.warp.dev/ai/mcp)
