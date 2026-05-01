# IntelliJ IDEA

**iris provider name:** `intellij`

## Config file

| Scope   | Path              |
|---------|-------------------|
| Project | `.idea/mcp.json`  |
| Global  | —                 |

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

- [JetBrains MCP documentation](https://www.jetbrains.com/help/idea/model-context-protocol.html)
