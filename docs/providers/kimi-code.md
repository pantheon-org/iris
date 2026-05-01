# Kimi Code

**iris provider name:** `kimi`

## Config file

| Scope  | Path               |
|--------|--------------------|
| Global | `~/.kimi/mcp.json` |

Project-level config is not supported. Use `--mcp-config-file <path>` at runtime to load an alternate file.

## Format

```json
{
  "mcpServers": {
    "context7": {
      "url": "https://mcp.context7.com/mcp",
      "headers": {
        "CONTEXT7_API_KEY": "your-key"
      }
    },
    "chrome-devtools": {
      "command": "npx",
      "args": ["chrome-devtools-mcp@latest"],
      "env": {
        "SOME_VAR": "value"
      }
    }
  }
}
```

Root key: `mcpServers` (standard MCP JSON format). Both stdio (`command`/`args`/`env`) and HTTP (`url`/`headers`) transports are supported.

## References

- [Kimi Code MCP documentation](https://www.kimi.com/code/docs/en/kimi-code-cli/customization/mcp.html)
- [Kimi Code CLI reference: kimi mcp](https://www.kimi.com/code/docs/en/kimi-code-cli/reference/kimi-mcp.html)
