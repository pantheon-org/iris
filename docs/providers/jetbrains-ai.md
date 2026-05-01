# JetBrains AI Assistant

**iris provider name:** not yet implemented

## Config file

| Scope   | Path                                                                        |
|---------|-----------------------------------------------------------------------------|
| Project | —                                                                           |
| Global  | IDE settings directory (OS and IDE version dependent — see below)           |

Global paths vary by OS and IDE:

| OS    | Approximate path                                                    |
|-------|---------------------------------------------------------------------|
| macOS | `~/Library/Application Support/JetBrains/<IDE><version>/options/`  |
| Linux | `~/.config/JetBrains/<IDE><version>/options/`                       |

The `<IDE>` and `<version>` components vary (e.g. `IdeaIC2024.3`, `PyCharm2024.3`), making the path non-deterministic.

## Format

MCP server configuration is stored in IDE XML options files (JetBrains settings format), not a standard JSON or TOML file. Programmatic editing requires parsing JetBrains XML settings.

## Notes

Due to IDE-specific paths and XML-based storage, automated syncing via iris is not currently feasible. Use the JetBrains AI Assistant settings UI to configure MCP servers.

## References

- [JetBrains AI Assistant MCP documentation](https://www.jetbrains.com/help/ai/mcp-support.html)
