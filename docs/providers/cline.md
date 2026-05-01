# Cline

**iris provider name:** not yet implemented

## Config file

Cline is a VS Code extension. MCP server configuration is stored in VS Code's extension global storage, not in a user-accessible fixed path. The file is typically at:

```text
~/.vscode/extensions/saoudrizwan.claude-dev-<version>/globalStorage/settings/cline_mcp_settings.json
```

The version suffix makes the path non-deterministic across updates. Cline also supports managing MCP servers through its in-app UI.

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

Root key: `mcpServers`. Additional Cline-specific fields: `disabled` (bool) and `autoApprove` (tool name array).

## Notes

Because the config path is version-dependent, automated syncing requires either VS Code API access or locating the path at runtime via `find`. A Go implementation would need to resolve the path dynamically.

## References

- [Cline MCP servers documentation](https://docs.cline.bot/mcp-servers/configuring-mcp-servers)
