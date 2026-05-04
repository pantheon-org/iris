Feature: Manage MCP server configs

  Background:
    Given a clean workspace

  Scenario: Add and remove servers persists correctly
    Given an MCP server "alpha" with command "cmd-alpha" and no args
    And an MCP server "beta" with command "cmd-beta" and no args
    And an MCP server "gamma" with command "cmd-gamma" and no args
    When I remove the server "gamma"
    And I reload the config from disk
    Then the config contains 2 servers
    And the config does not contain server "gamma"
    And the config contains server "alpha" with command "cmd-alpha"
    And the config contains server "beta" with command "cmd-beta"
