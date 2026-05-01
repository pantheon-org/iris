# Kilo Code

**iris provider name:** not yet implemented

## Config file

| Scope   | Path                                                                                                            |
|---------|-----------------------------------------------------------------------------------------------------------------|
| Project | —                                                                                                               |
| Global  | `~/.vscode/extensions/kilocode.kilo-code-<version>/globalStorage/settings/kilo_mcp_settings.json`              |

The `<version>` suffix makes the global path non-deterministic across updates. Kilo Code also supports managing MCP servers through its in-app UI.

## Format

```json
{
  "mcpServers": {
    "server-name": {
      "command": "npx",
      "args": ["-y", "@modelcontextprotocol/server-filesystem", "/tmp"],
      "env": { "KEY": "value" },
      "disabled": false,
      "autoApprove": []
    }
  }
}
```

Root key: `mcpServers`. Format follows the Cline convention.

## Notes

Path is version-dependent (same limitation as Cline). A Go implementation would need to resolve the path dynamically.

## References

- [Kilo Code MCP documentation](https://kilocode.ai/docs/features/mcp)
