# Mistral Vibe

**iris provider name:** `mistral-vibe`

## Config file

| Scope  | Path                   |
|--------|------------------------|
| Global | `~/.vibe/config.toml`  |

Project-level config is not supported. The home directory can be overridden via `VIBE_HOME`.

## Format

```toml
model = "codestral-latest"
api_key_env = "MISTRAL_API_KEY"

[[mcp_servers]]
name    = "fetch"
transport = "stdio"
command = "uvx"
args    = ["mcp-server-fetch"]

[mcp_servers.env]
DEBUG = "1"

[[mcp_servers]]
name      = "context7"
transport = "http"
url       = "https://mcp.context7.com/mcp"

[mcp_servers.headers]
CONTEXT7_API_KEY = "your-key"
```

Root key: `mcp_servers` (TOML array of tables). Each entry has an explicit `name` field. Supported transports: `stdio`, `http`, `streamable-http`. iris preserves all existing top-level keys (e.g. `model`, `api_key_env`).

## References

- [Mistral Vibe GitHub](https://github.com/mistralai/mistral-vibe)
- [Mistral Vibe MCP configuration docs](https://github.com/mistralai/mistral-vibe/blob/main/docs/mcp.md)
