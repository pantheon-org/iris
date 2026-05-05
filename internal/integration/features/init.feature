Feature: Init command

  Background:
    Given a clean workspace

  Scenario: Init creates the iris config file when it does not exist
    When I run init
    Then the iris config file exists on disk
    And the iris config file is valid JSON with version 1

  Scenario: Init is idempotent when config already exists
    Given the iris config already exists with one server
    When I run init
    Then the iris config file exists on disk
    And the output contains "already"

  # ── Interactive import flow ────────────────────────────────────────────────

  Scenario: Interactive init with no provider configs shows empty import list
    Given no provider config files exist
    When I run interactive init and select no servers
    Then the iris config file exists on disk
    And the iris config contains 0 servers

  Scenario: Interactive init shows detected servers with provider and scope info
    Given a Claude Code project config exists with server "fmt" command "npx" args "-y @modelcontextprotocol/server-filesystem /tmp"
    When I run interactive init and collect the import candidates
    Then the import candidates include an entry for server "fmt" from provider "claude" with scope "project"

  Scenario: Interactive init imports selected servers from a single provider
    Given a Claude Code project config exists with server "fmt" command "npx" args "-y @modelcontextprotocol/server-filesystem /tmp"
    When I run interactive init and select server "fmt"
    Then the iris config file exists on disk
    And the iris config contains 1 servers
    And the iris config contains server "fmt" with command "npx"

  Scenario: Interactive init imports servers from multiple providers
    Given a Claude Code project config exists with server "fmt" command "npx" args "-y @modelcontextprotocol/server-filesystem /tmp"
    And a Cursor project config exists with server "github" command "uvx" args "mcp-server-github"
    When I run interactive init and select all discovered servers
    Then the iris config contains 2 servers
    And the iris config contains server "fmt" with command "npx"
    And the iris config contains server "github" with command "uvx"

  Scenario: Interactive init deduplicates servers with the same name across providers
    Given a Claude Code project config exists with server "shared" command "npx" args "-y foo"
    And a Cursor project config exists with server "shared" command "npx" args "-y foo"
    When I run interactive init and collect the import candidates
    Then the import candidates contain exactly 2 entries for server "shared"

  Scenario: Interactive init shows global scope for providers with global configs
    Given a global Google Gemini config exists with server "gemini-srv" command "python" args "-m gemini_mcp"
    When I run interactive init and collect the import candidates
    Then the import candidates include an entry for server "gemini-srv" from provider "gemini" with scope "global"

  Scenario: Interactive init skips import when user selects none
    Given a Claude Code project config exists with server "fmt" command "npx" args "-y @modelcontextprotocol/server-filesystem /tmp"
    When I run interactive init and select no servers
    Then the iris config contains 0 servers

  Scenario: Interactive init allows manual server entry after import step
    Given no provider config files exist
    When I run interactive init, skip import, and manually add server "manual-srv" command "node" args "server.js"
    Then the iris config contains 1 servers
    And the iris config contains server "manual-srv" with command "node"

  Scenario: Interactive init gracefully skips a malformed provider config
    Given a malformed Claude Code project config exists
    When I run interactive init and select no servers
    Then the iris config file exists on disk
    And the iris config contains 0 servers

  Scenario: Interactive init imports from valid providers even when one config is malformed
    Given a malformed Claude Code project config exists
    And a Cursor project config exists with server "github" command "uvx" args "mcp-server-github"
    When I run interactive init and select all discovered servers
    Then the iris config contains 1 servers
    And the iris config contains server "github" with command "uvx"

  Scenario: Interactive init combines imported and manually added servers
    Given a Claude Code project config exists with server "fmt" command "npx" args "-y @modelcontextprotocol/server-filesystem /tmp"
    When I run interactive init, import server "fmt", and manually add server "extra" command "uvx" args "extra-tool"
    Then the iris config contains 2 servers
    And the iris config contains server "fmt" with command "npx"
    And the iris config contains server "extra" with command "uvx"
