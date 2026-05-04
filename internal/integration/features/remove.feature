Feature: Remove MCP servers

  Background:
    Given a clean workspace

  Scenario: Remove an existing server persists the change
    Given an MCP server "alpha" with command "cmd-alpha" and no args
    And an MCP server "beta" with command "cmd-beta" and no args
    When I remove the server "alpha"
    And I reload the config from disk
    Then the config does not contain server "alpha"
    And the config contains server "beta" with command "cmd-beta"

  Scenario: Removing a non-existent server returns ErrServerNotFound
    When I try to remove the server "ghost"
    Then the last error wraps "server not found"
