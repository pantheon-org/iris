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
| `iris init --interactive/-I` | Interactive wizard — scans all provider configs (global + project), shows a TUI multi-select to pick servers to import (grouped by server name), then lets you add more manually. Saves detected providers to `.iris.json` so `iris sync` knows where to sync by default |
| `iris init --provider/-p <name> ...` | Limit init to specific providers |
| `iris add <name> --command/-c <cmd> [--args/-a ...] [--env/-e KEY=VAL ...] [--transport/-t stdio\|sse] [--url/-u <url>]` | Add or update a server |
| `iris remove <name>` | Remove a server |
| `iris list` | List all configured servers |
| `iris sync` | Sync both global and local configs for all providers (or those listed in `.iris.json` `providers` field) |
| `iris sync --provider/-p <name> ...` | Override the provider list for this sync |
| `iris sync --global/-g` | Sync only home-directory (global) configs |
| `iris sync --local/-l` | Sync only project-local configs |
| `iris sync --interactive/-I` | Interactive wizard — select which providers and scope (global/local/both) to sync |
| `iris status` | Show per-provider sync status |

## Canonical config format

`.iris.json` is the single source of truth for your MCP server definitions:

```json
{
  "version": 1,
  "providers": ["claude", "cursor"],
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

The `providers` field is written automatically by `iris init -I` and tells `iris sync` which providers to target by default. You can edit it by hand to add or remove providers between init runs. Omitting it (or running non-interactive init) falls back to syncing all providers.

## Supported providers

| Provider | Project config | Global config |
|---|---|---|
| Claude Code (`claude`) | `.mcp.json` | `~/.mcp.json` |
| Claude Desktop (`claude-desktop`) | — | `~/Library/Application Support/Claude/claude_desktop_config.json` (macOS) |
| Gemini CLI (`gemini`) | `.gemini/settings.json` | `~/.gemini/settings.json` |
| OpenCode (`opencode`) | `opencode.json` | `~/.config/opencode/opencode.json` |
| OpenAI Codex (`codex`) | `.codex/config.toml` | `~/.codex/config.toml` |
| Cursor (`cursor`) | `.cursor/mcp.json` | `~/.cursor/mcp.json` |
| Windsurf (`windsurf`) | — | `~/.codeium/windsurf/mcp_config.json` |
| VS Code Copilot (`copilot`) | `.vscode/mcp.json` | — |
| Zed (`zed`) | — | `~/.config/zed/settings.json` |
| Qwen Code (`qwen`) | `.qwen/settings.json` | `~/.qwen/settings.json` |
| Warp (`warp`) | — | `~/.warp/mcp.json` |
| Kimi Code (`kimi`) | — | `~/.kimi/mcp.json` |
| Mistral Vibe (`mistral-vibe`) | `.vibe/config.toml` | `~/.vibe/config.toml` |
| IntelliJ IDEA (`intellij`) | `.idea/mcp.json` | — |

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
