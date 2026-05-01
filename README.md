# iris

> Sync MCP server configs across all your AI providers.

`iris` is a CLI tool that keeps MCP server configurations in sync across Claude Code, Gemini CLI, OpenCode, and Codex. Define your servers once in `.iris.json` and run `iris sync` to propagate them everywhere.

## Installation

Build from source:

```sh
mise run build   # produces dist/iris
```

Then add `dist/iris` to your `$PATH`, or run it directly.

## Commands

| Command | Description |
|---|---|
| `iris init` | Scaffold `.iris.json` in the current project |
| `iris init -I` | Interactive wizard to scaffold config |
| `iris add <name> --command <cmd> [--args ...] [--env KEY=VAL ...] [--transport stdio\|sse]` | Add or update a server |
| `iris remove <name>` | Remove a server |
| `iris list` | List all configured servers |
| `iris sync` | Sync all provider config files |
| `iris status` | Show per-provider sync status |

## Canonical config format

`.iris.json` is the single source of truth for your MCP server definitions:

```json
{
  "version": 1,
  "providers": ["claude", "gemini", "opencode", "codex"],
  "servers": {
    "filesystem": {
      "transport": "stdio",
      "command": "npx",
      "args": ["-y", "@modelcontextprotocol/server-filesystem", "/tmp"]
    },
    "fetch": {
      "transport": "stdio",
      "command": "uvx",
      "args": ["mcp-server-fetch"]
    }
  }
}
```

## Supported providers

| Provider | Config file | Scope |
|---|---|---|
| Claude Code | `.mcp.json` | Project |
| Gemini CLI | `~/.config/gemini/settings.json` | Global |
| OpenCode | `opencode.json` | Project |
| Codex | `~/.codex/config.toml` | Global |

## Development

```sh
mise run build   # build binary to dist/iris
mise run test    # run tests with race detector
mise run lint    # run golangci-lint
```

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

MIT — see [LICENSE](LICENSE).
