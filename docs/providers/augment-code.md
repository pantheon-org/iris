# Augment Code

**iris provider name:** not yet implemented

## Config file

Augment Code is a VS Code (and JetBrains) extension. MCP server configuration is managed through the extension's settings panel and stored in VS Code extension global storage at a version-dependent path:

```text
~/.vscode/extensions/augment.vscode-augment-<version>/globalStorage/
```

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
