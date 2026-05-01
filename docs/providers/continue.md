# Continue

**iris provider name:** not yet implemented

## Config file

| Scope   | Path                       |
|---------|----------------------------|
| Project | —                          |
| Global  | `~/.continue/config.yaml`  |

Project-level config is not supported. Older versions used `~/.continue/config.json`. Continue 3.x uses YAML.

## Format (config.yaml)

```yaml
mcpServers:
  - name: server-name
    command: npx
    args:
      - -y
      - "@modelcontextprotocol/server-filesystem"
      - /tmp
    env:
      KEY: value
```

Root key: `mcpServers` (YAML array, not an object map). Each entry has an explicit `name` field.

## Format (config.json, deprecated)

```json
{
  "mcpServers": [
    {
      "name": "server-name",
      "command": "npx",
      "args": ["-y", "@modelcontextprotocol/server-filesystem", "/tmp"]
    }
  ]
}
```

## Notes

The YAML array format differs from the JSON object map used by most other providers. A Go implementation would require a YAML dependency and a separate provider type.

## References

- [Continue MCP documentation](https://docs.continue.dev/customize/deep-dives/mcp)
- [Continue config.yaml reference](https://docs.continue.dev/reference)
