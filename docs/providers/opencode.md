# OpenCode

**iris provider name:** `opencode`

## Config file

| Scope   | Path                                    |
|---------|-----------------------------------------|
| Project | `opencode.json`                         |
| Global  | `~/.config/opencode/opencode.json`      |

iris writes to the project-level file.

## Format

```json
{
  "$schema": "https://opencode.ai/config.json",
  "mcp": {
    "server-name": {
      "type": "local",
      "command": ["npx", "-y", "@modelcontextprotocol/server-filesystem", "/tmp"],
      "environment": { "KEY": "value" },
      "enabled": true
    }
  }
}
```

Root key: `mcp`. The command is an array (binary + args merged). Transport type is `"local"` for stdio and `"remote"` for HTTP.

## References

- [OpenCode MCP documentation](https://opencode.ai/docs/mcp)
