Feature: Status of provider configs

  Background:
    Given a clean workspace

  Scenario: Status before sync shows missing for all providers
    Given an MCP server "tool" with command "uvx" and args "some-tool"
    When I run status
    Then the status output contains provider "claude" with status "missing"

  Scenario: Status after sync shows synced for all providers
    Given an MCP server "tool" with command "uvx" and args "some-tool"
    When I sync to all providers
    And I run status
    Then the status output contains provider "claude" with status "synced"
    And the status output contains provider "gemini" with status "synced"
    And the status output contains provider "opencode" with status "synced"

  Scenario: Status detects desync when config is manually modified
    Given an MCP server "tool" with command "uvx" and args "some-tool"
    When I sync to all providers
    And I corrupt the provider config file ".mcp.json"
    And I run status
    Then the status output contains provider "claude" with status "desync"

  Scenario: Status with --json flag emits JSON providers array
    Given an MCP server "tool" with command "uvx" and args "some-tool"
    When I sync to all providers
    And I run status with JSON output
    Then the JSON status output has a "providers" array
    And the JSON status providers contain an entry for provider "claude" with status "synced"
