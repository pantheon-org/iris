Feature: Sync MCP servers to all providers

  Background:
    Given a clean workspace

  Scenario: Full pipeline syncs servers to all providers
    Given an MCP server "filesystem" with command "npx" and args "-y,@modelcontextprotocol/server-filesystem,/tmp"
    And an MCP server "fetch" with command "uvx" and args "mcp-server-fetch"
    When I sync to all providers
    Then the provider config file ".mcp.json" exists
    And the provider config file "gemini-settings.json" exists
    And the provider config file "opencode.json" exists
    And the provider config file "codex-config.toml" exists
    And the JSON provider file ".mcp.json" contains servers "filesystem,fetch" under key "mcpServers"
    And the JSON provider file "gemini-settings.json" contains servers "filesystem,fetch" under key "mcpServers"
    And the opencode provider file "opencode.json" contains servers "filesystem,fetch"
    And the TOML provider file "codex-config.toml" contains servers "filesystem,fetch"
    And the JSON provider file "claude-desktop-config.json" contains servers "filesystem,fetch" under key "mcpServers"
    And the JSON provider file "windsurf-config.json" contains servers "filesystem,fetch" under key "mcpServers"
    And the JSON provider file "warp-mcp.json" contains servers "filesystem,fetch" under key "mcpServers"
    And the JSON provider file "kimi-settings.json" contains servers "filesystem,fetch" under key "mcpServers"
    And the JSON provider file ".idea/mcp.json" contains servers "filesystem,fetch" under key "mcpServers"
    And the JSON provider file ".qwen/settings.json" contains servers "filesystem,fetch" under key "mcpServers"

  Scenario: Sync is idempotent
    Given an MCP server "filesystem" with command "npx" and args "-y,@modelcontextprotocol/server-filesystem,/tmp"
    And an MCP server "fetch" with command "uvx" and args "mcp-server-fetch"
    When I sync to all providers
    Then all providers report status "created"
    When I sync to all providers again
    Then all providers report status "unchanged"

  Scenario: All supported providers write provider-specific formats
    Given an MCP server "tool" with command "uvx" and args "some-tool"
    When I sync to all providers
    Then the JSON provider file ".mcp.json" contains servers "tool" under key "mcpServers"
    And the JSON provider file "gemini-settings.json" contains servers "tool" under key "mcpServers"
    And the opencode provider file "opencode.json" contains servers "tool"
    And the TOML provider file "codex-config.toml" contains servers "tool"
    And the JSON provider file ".cursor/mcp.json" contains servers "tool" under key "mcpServers"
    And the JSON provider file ".vscode/mcp.json" contains servers "tool" under key "servers"
    And the JSON provider file ".qwen/settings.json" contains servers "tool" under key "mcpServers"
    And the JSON provider file ".idea/mcp.json" contains servers "tool" under key "mcpServers"
    And the zed provider file "zed-settings.json" contains servers "tool"
    And the TOML mistral provider file "mistral-vibe-config.toml" contains servers "tool"
    And the JSON provider file "warp-mcp.json" contains servers "tool" under key "mcpServers"
    And the JSON provider file "kimi-settings.json" contains servers "tool" under key "mcpServers"
    And the JSON provider file "windsurf-config.json" contains servers "tool" under key "mcpServers"

  Scenario: Sync with --json flag emits JSON results
    Given an MCP server "fetch" with command "uvx" and args "mcp-server-fetch"
    When I sync to all providers with JSON output
    Then the JSON sync output has a "results" array
    And the JSON sync results contain an entry for provider "claude" with status "created"

  Scenario: Re-sync after config exists reports updated status
    Given an MCP server "tool" with command "uvx" and args "some-tool"
    When I sync to all providers
    And I add an MCP server "tool2" with command "uvx" and args "other-tool"
    And I sync to all providers
    Then the provider config file ".mcp.json" exists
    And the JSON provider file ".mcp.json" contains servers "tool,tool2" under key "mcpServers"

  Scenario: OpenCode writes correct field formats
    Given an MCP server "myserver" with command "node" and args "server.js"
    And an MCP server "myserver" with command "node" and env "MY_KEY=MY_VAL"
    When I sync to all providers
    Then the opencode provider file "opencode.json" contains servers "myserver"
    And the opencode server "myserver" in file "opencode.json" has correct field format

  Scenario: Copilot writes servers under "servers" key
    Given an MCP server "tool" with command "uvx" and args "some-tool"
    When I sync to all providers
    Then the JSON provider file ".vscode/mcp.json" contains servers "tool" under key "servers"

  Scenario: Copilot does not emit unsupported fields
    Given an MCP server "tool" with command "uvx" and args "some-tool"
    When I sync to all providers
    Then the copilot server "tool" in file ".vscode/mcp.json" does not have field "headers"
    And the copilot server "tool" in file ".vscode/mcp.json" does not have field "cwd"
    And the copilot server "tool" in file ".vscode/mcp.json" does not have field "enabled"

  Scenario: SSE servers sync url field to JSON providers
    Given an SSE server "remote" with URL "https://example.com/sse"
    When I sync to all providers
    Then the JSON provider file ".mcp.json" server "remote" under key "mcpServers" has field "url"
    And the JSON provider file "gemini-settings.json" server "remote" under key "mcpServers" has field "url"

  Scenario: Sync preserves extra keys in existing Gemini config
    Given a provider file "gemini-settings.json" exists with extra key "theme" set to "dark"
    And an MCP server "tool" with command "uvx" and args "some-tool"
    When I sync to all providers
    Then the JSON provider file "gemini-settings.json" still has key "theme"
    And the JSON provider file "gemini-settings.json" contains servers "tool" under key "mcpServers"

  Scenario: Sync preserves extra keys in existing Zed config
    Given a provider file "zed-settings.json" exists with extra key "vim_mode" set to "true"
    And an MCP server "tool" with command "uvx" and args "some-tool"
    When I sync to all providers
    Then the JSON provider file "zed-settings.json" still has key "vim_mode"
    And the zed provider file "zed-settings.json" contains servers "tool"

  Scenario: Sync uses providers from .iris.json when no --provider flag is given
    Given an MCP server "tool" with command "uvx" and args "some-tool"
    And the iris config providers list is set to "cursor"
    When I sync to all providers
    Then the provider config file ".cursor/mcp.json" exists
    And the provider config file ".mcp.json" does not exist

  Scenario: --global and --provider can be combined to sync global config of one provider
    Given an MCP server "tool" with command "uvx" and args "some-tool"
    When I sync to provider "claude" with "global" scope
    Then the provider config file "claude-global.json" exists
    And the provider config file ".mcp.json" does not exist

  Scenario: --local and --provider can be combined to sync local config of one provider
    Given an MCP server "tool" with command "uvx" and args "some-tool"
    When I sync to provider "claude" with "local" scope
    Then the provider config file ".mcp.json" exists
    And the provider config file "claude-global.json" does not exist

  Scenario: Env vars are written to JSON providers
    Given an MCP server "envtool" with command "node" and args "server.js"
    And the server "envtool" has env var "MY_KEY" set to "MY_VAL"
    When I sync to all providers
    Then the JSON provider file ".mcp.json" server "envtool" under key "mcpServers" has env var "MY_KEY"
    And the JSON provider file "gemini-settings.json" server "envtool" under key "mcpServers" has env var "MY_KEY"
