# Zed

**iris provider name:** `zed`

## Config file

| Scope  | Path                            |
|--------|---------------------------------|
| Global | `~/.config/zed/settings.json`   |

Project-level config is not supported for MCP servers directly; servers are defined in the global settings file.

## Format

```json
{
  "context_servers": {
    "server-name": {
      "command": "npx",
      "args": ["-y", "@modelcontextprotocol/server-filesystem", "/tmp"],
      "env": { "KEY": "value" }
    }
  }
}
```

Root key: `context_servers`. For remote/HTTP servers, use `"url"` and `"headers"` instead of `"command"`. iris merges only the `context_servers` block and preserves all other Zed settings (e.g. `"theme"`, `"buffer_font_size"`).

## References

- [Zed MCP documentation](https://zed.dev/docs/ai/mcp.html)
