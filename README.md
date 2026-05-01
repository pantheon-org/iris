# iris

> Sync MCP server configs across all your AI providers.

`iris` is a CLI tool that keeps MCP server configurations in sync across all your AI providers. Define your servers once in `.iris.json` and run `iris sync` to propagate them everywhere.

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
| `iris add <name> --command <cmd> [--args ...] [--env KEY=VAL ...] [--transport stdio\|sse] [--url <url>]` | Add or update a server |
| `iris remove <name>` | Remove a server |
| `iris list` | List all configured servers |
| `iris sync` | Sync all provider config files |
| `iris status` | Show per-provider sync status |

## Canonical config format

`.iris.json` is the single source of truth for your MCP server definitions:

```json
{
  "version": 1,
  "providers": ["claude", "gemini", "opencode", "codex", "cursor", "windsurf"],
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

| Provider | Project config | Global config |
|---|---|---|
| Claude Code (`claude`) | `.mcp.json` | `~/.mcp.json` |
| Claude Desktop (`claude-desktop`) | — | `~/Library/Application Support/Claude/claude_desktop_config.json` (macOS) |
| Gemini CLI (`gemini`) | — | `~/.config/gemini/settings.json` |
| OpenCode (`opencode`) | `opencode.json` | `~/.config/opencode/opencode.json` |
| OpenAI Codex (`codex`) | — | `~/.codex/config.toml` |
| Cursor (`cursor`) | `.cursor/mcp.json` | `~/.cursor/mcp.json` |
| Windsurf (`windsurf`) | — | `~/.codeium/windsurf/mcp_config.json` |
| VS Code Copilot (`vscode-copilot`) | `.vscode/mcp.json` | — |
| Zed (`zed`) | — | `~/.config/zed/settings.json` |
| Qwen Code (`qwen`) | — | `~/.qwen/settings.json` |
| Warp (`warp`) | — | `~/.warp/mcp.json` |
| Kimi Code (`kimi`) | — | `~/.kimi/mcp.json` |
| Mistral Vibe (`mistral-vibe`) | — | `~/.vibe/config.toml` |

Providers not yet implemented (config references available in [`docs/providers/`](docs/providers/)): Continue, Cline, Kilo Code, JetBrains AI, Augment Code, Amp.

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
