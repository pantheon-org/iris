# Amp

**iris provider name:** not yet implemented

## Config file

| Scope   | Path                      |
|---------|---------------------------|
| Project | —                         |
| Global  | `~/.config/amp/mcp.json`  |

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

## Notes

Path and format are based on Amp's public documentation. Verify against `ampcode.com/manual/mcp` before implementing.

## References

- [Amp MCP documentation](https://ampcode.com/manual/mcp)
