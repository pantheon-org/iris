# iris

> Manage MCP server configs across AI providers.

A CLI tool that syncs MCP server configurations across Claude Code, Gemini CLI, OpenCode, and Codex — so you define your servers once and iris keeps all providers in sync.

## Install

### Homebrew

```sh
brew install pantheon-org/tap/iris
```

### Go

```sh
go install github.com/pantheon-org/iris/cmd/iris@latest
```

### Binary

Download from [Releases](https://github.com/pantheon-org/iris/releases).

## Quickstart

```sh
iris sync          # sync MCP servers to all detected providers
iris list          # list configured MCP servers
iris add <name>    # add a new MCP server
```

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

MIT — see [LICENSE](LICENSE).
