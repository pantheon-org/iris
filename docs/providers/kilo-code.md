# Kilo Code

**iris provider name:** not yet implemented

## Config file

Kilo Code is a VS Code extension (fork of Cline). MCP configuration is stored in VS Code's extension global storage at a version-dependent path, similar to Cline:

```text
~/.vscode/extensions/kilocode.kilo-code-<version>/globalStorage/settings/kilo_mcp_settings.json
```

Kilo Code also supports managing MCP servers through its in-app UI.

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
