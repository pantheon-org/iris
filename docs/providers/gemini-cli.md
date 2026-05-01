# Gemini CLI

**iris provider name:** `gemini`

## Config file

| Scope   | Path                              |
|---------|-----------------------------------|
| Project | —                                 |
| Global  | `~/.config/gemini/settings.json`  |

Project-level config is not supported.

## Format

```json
{
  "mcpServers": {
    "server-name": {
      "command": "uvx",
      "args": ["mcp-server-fetch"],
      "env": { "KEY": "value" },
      "type": "stdio"
    }
  }
}
```

Root key: `mcpServers`. The `settings.json` file may contain other Gemini CLI settings; iris preserves all non-`mcpServers` keys.

## References

- [Gemini CLI MCP documentation](https://github.com/google-gemini/gemini-cli/blob/main/docs/mcp-server.md)
