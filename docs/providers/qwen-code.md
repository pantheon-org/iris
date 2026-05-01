# Qwen Code

**iris provider name:** `qwen`

## Config file

| Scope   | Path                    |
|---------|-------------------------|
| Project | —                       |
| Global  | `~/.qwen/settings.json` |

Project-level config is not supported.

## Format

```json
{
  "mcpServers": {
    "server-name": {
      "command": "npx",
      "args": ["-y", "@modelcontextprotocol/server-filesystem", "/tmp"],
      "env": { "KEY": "value" },
      "cwd": "./server-directory",
      "timeout": 30000,
      "trust": false
    }
  }
}
```

Root key: `mcpServers`. Qwen Code supports additional optional fields `cwd`, `timeout`, and `trust` that iris does not currently emit but preserves when merging.

## References

- [Qwen Code MCP server documentation](https://qwenlm.github.io/qwen-code-docs/en/developers/tools/mcp-server/)
