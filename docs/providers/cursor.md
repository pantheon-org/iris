# Cursor

**iris provider name:** `cursor`

## Config file

| Scope   | Path                   |
|---------|------------------------|
| Project | `.cursor/mcp.json`     |
| Global  | `~/.cursor/mcp.json`   |

iris writes to the project-level file.

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

Root key: `mcpServers` (standard MCP JSON format).

## References

- [Cursor MCP documentation](https://docs.cursor.com/context/model-context-protocol)
- [Cursor MCP configuration reference](https://docs.cursor.com/context/model-context-protocol#configuration)
