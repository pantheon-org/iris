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
