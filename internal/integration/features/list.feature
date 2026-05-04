Feature: List MCP servers

  Background:
    Given a clean workspace

  Scenario: List with no servers prints empty message
    When I run list
    Then the output contains "No servers configured."

  Scenario: List with multiple servers prints them alphabetically
    Given an MCP server "zebra" with command "cmd-z" and no args
    And an MCP server "alpha" with command "cmd-a" and no args
    And an MCP server "mango" with command "cmd-m" and no args
    When I run list
    Then the output lines appear in order "alpha,mango,zebra"

  Scenario: List with --json flag emits a JSON servers array
    Given an MCP server "fetch" with command "uvx" and args "mcp-server-fetch"
    When I run list with JSON output
    Then the JSON output has a "servers" array
    And the JSON servers array contains an entry with name "fetch" and command "uvx"
