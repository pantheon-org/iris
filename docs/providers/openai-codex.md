# OpenAI Codex

**iris provider name:** `codex`

## Config file

| Scope   | Path                     |
|---------|--------------------------|
| Project | `.codex/config.toml`     |
| Global  | `~/.codex/config.toml`   |

## Format

```toml
[mcp_servers.server-name]
command = "npx"
args    = ["-y", "@modelcontextprotocol/server-filesystem", "/tmp"]

[mcp_servers.server-name.env]
KEY = "value"
```

Root key: `mcp_servers` (TOML table keyed by server name). Codex infers the transport from the fields present:

- stdio servers use `command`, `args`, and optional `env`
- streamable HTTP servers use `url` and optional `http_headers`

## References

- [Codex CLI GitHub](https://github.com/openai/codex)
- [Codex CLI configuration reference](https://github.com/openai/codex#configuration)
- [OpenAI Codex MCP](https://developers.openai.com/codex/mcp)
