Feature: Add MCP servers

  Background:
    Given a clean workspace

  Scenario: Add a stdio server persists to disk
    When I add a stdio server "myserver" with command "npx" and args "-y,some-package"
    Then the config contains server "myserver" with command "npx"
    And the iris config file exists on disk

  Scenario: Add an SSE server persists to disk
    When I add an SSE server "remote" with url "https://example.com/mcp"
    Then the config contains server "remote" with transport "sse"
    And the iris config file exists on disk

  Scenario: Add a server with env vars persists them
    When I add a stdio server "envserver" with command "node" and env "TOKEN=abc123,DEBUG=true"
    Then the config contains server "envserver" with env var "TOKEN" equal to "abc123"
    And the config contains server "envserver" with env var "DEBUG" equal to "true"

  Scenario: Adding a duplicate name overwrites the existing server
    When I add a stdio server "dup" with command "cmd-v1" and no args
    And I add a stdio server "dup" with command "cmd-v2" and no args
    Then the config contains server "dup" with command "cmd-v2"
    And the config has exactly 1 server

  Scenario: Adding a stdio server without a command returns an error
    When I try to add a stdio server "bad" with no command
    Then the last error wraps "malformed config"
