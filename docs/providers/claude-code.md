# Claude Code

**iris provider name:** `claude`

## Config file

| Scope   | Path              |
|---------|-------------------|
| Project | `.mcp.json`        |
| Global  | `~/.claude.json`  |

## Format

```json
{
  "mcpServers": {
    "server-name": {
      "command": "npx",
      "args": ["-y", "@modelcontextprotocol/server-filesystem", "/tmp"],
      "env": { "KEY": "value" },
      "type": "stdio"
    }
  }
}
```

Root key: `mcpServers` (object map, keyed by server name).

## References

- [Claude Code MCP docs](https://docs.anthropic.com/en/docs/claude-code/mcp)
- [MCP configuration reference](https://modelcontextprotocol.io/docs/develop/connect-local-servers)
