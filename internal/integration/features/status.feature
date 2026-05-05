Feature: Status of provider configs

  Background:
    Given a clean workspace

  Scenario: Status before sync shows missing for all providers
    Given an MCP server "tool" with command "uvx" and args "some-tool"
    When I run status
    Then the status output contains provider "claude" with status "missing"
    And the status output contains provider "claude-desktop" with status "missing"
    And the status output contains provider "gemini" with status "missing"
    And the status output contains provider "opencode" with status "missing"
    And the status output contains provider "codex" with status "missing"
    And the status output contains provider "cursor" with status "missing"
    And the status output contains provider "windsurf" with status "missing"
    And the status output contains provider "copilot" with status "missing"
    And the status output contains provider "zed" with status "missing"
    And the status output contains provider "qwen" with status "missing"
    And the status output contains provider "warp" with status "missing"
    And the status output contains provider "kimi" with status "missing"
    And the status output contains provider "mistral-vibe" with status "missing"
    And the status output contains provider "intellij" with status "missing"

  Scenario: Status after sync shows synced for all providers
    Given an MCP server "tool" with command "uvx" and args "some-tool"
    When I sync to all providers
    And I run status
    Then the status output contains provider "claude" with status "synced"
    And the status output contains provider "claude-desktop" with status "synced"
    And the status output contains provider "gemini" with status "synced"
    And the status output contains provider "opencode" with status "synced"
    And the status output contains provider "codex" with status "synced"
    And the status output contains provider "cursor" with status "synced"
    And the status output contains provider "windsurf" with status "synced"
    And the status output contains provider "copilot" with status "synced"
    And the status output contains provider "zed" with status "synced"
    And the status output contains provider "qwen" with status "synced"
    And the status output contains provider "warp" with status "synced"
    And the status output contains provider "kimi" with status "synced"
    And the status output contains provider "mistral-vibe" with status "synced"
    And the status output contains provider "intellij" with status "synced"

  Scenario: Status detects desync when config is manually modified
    Given an MCP server "tool" with command "uvx" and args "some-tool"
    When I sync to all providers
    And I corrupt the provider config file ".mcp.json"
    And I run status
    Then the status output contains provider "claude" with status "desync"

  Scenario: JSON status output contains all 14 providers
    Given an MCP server "tool" with command "uvx" and args "some-tool"
    When I sync to all providers
    And I run status with JSON output
    Then the JSON status output has a "providers" array
    And the JSON status providers contain an entry for provider "claude" with status "synced"
    And the JSON status providers contain an entry for provider "claude-desktop" with status "synced"
    And the JSON status providers contain an entry for provider "gemini" with status "synced"
    And the JSON status providers contain an entry for provider "opencode" with status "synced"
    And the JSON status providers contain an entry for provider "codex" with status "synced"
    And the JSON status providers contain an entry for provider "cursor" with status "synced"
    And the JSON status providers contain an entry for provider "windsurf" with status "synced"
    And the JSON status providers contain an entry for provider "copilot" with status "synced"
    And the JSON status providers contain an entry for provider "zed" with status "synced"
    And the JSON status providers contain an entry for provider "qwen" with status "synced"
    And the JSON status providers contain an entry for provider "warp" with status "synced"
    And the JSON status providers contain an entry for provider "kimi" with status "synced"
    And the JSON status providers contain an entry for provider "mistral-vibe" with status "synced"
    And the JSON status providers contain an entry for provider "intellij" with status "synced"

  Scenario: Status distinguishes "claude" from "claude-desktop"
    Given an MCP server "tool" with command "uvx" and args "some-tool"
    When I sync to all providers
    And I run status
    Then the status output contains provider "claude" with status "synced"
    And the status output contains provider "claude-desktop" with status "synced"
