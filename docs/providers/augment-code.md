# Augment Code

**iris provider name:** not yet implemented

## Config file

| Scope   | Path                                                                              |
|---------|-----------------------------------------------------------------------------------|
| Project | —                                                                                 |
| Global  | `~/.vscode/extensions/augment.vscode-augment-<version>/globalStorage/`           |

The `<version>` suffix makes the global path non-deterministic across updates.

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

Root key: `mcpServers` (standard format).

## Notes

Like Cline and Kilo Code, the config path is version-dependent. MCP servers can be configured through Augment Code's UI (Settings → MCP Servers).

## References

- [Augment Code MCP documentation](https://docs.augmentcode.com/setup-augment/mcp)
