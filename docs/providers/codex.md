# OpenAI Codex

**iris provider name:** `codex`

## Config file

| Scope   | Path                     |
|---------|--------------------------|
| Project | `.codex/config.toml`     |
| Global  | `~/.codex/config.toml`   |

## Format

```toml
[[mcp_servers]]
name    = "server-name"
command = "npx"
args    = ["-y", "@modelcontextprotocol/server-filesystem", "/tmp"]
type    = "stdio"

[mcp_servers.env]
KEY = "value"
```

Root key: `mcp_servers` (TOML array of tables). Unlike other providers, each server has an explicit `name` field inside the table entry.

## References

- [Codex CLI GitHub](https://github.com/openai/codex)
- [Codex CLI configuration reference](https://github.com/openai/codex#configuration)
