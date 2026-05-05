package integration_test

import (
	"fmt"

	"github.com/pantheon-org/iris/internal/cli"
	"github.com/pantheon-org/iris/internal/types"
)

// ── add steps ─────────────────────────────────────────────────────────────────

func (s *scenarioCtx) anMCPServerWithCommandAndArgs(name, command, rawArgs string) error {
	args := splitArgs(rawArgs)
	if err := cli.RunAdd(s.cfg, s.store, name, types.MCPServer{
		Transport: types.TransportStdio,
		Command:   command,
		Args:      args,
	}, nil, nil); err != nil {
		return fmt.Errorf("add server %s: %w", name, err)
	}
	return nil
}

func (s *scenarioCtx) anMCPServerWithCommandAndNoArgs(name, command string) error {
	if err := cli.RunAdd(s.cfg, s.store, name, types.MCPServer{
		Transport: types.TransportStdio,
		Command:   command,
	}, nil, nil); err != nil {
		return fmt.Errorf("add server %s: %w", name, err)
	}
	return nil
}

func (s *scenarioCtx) anMCPServerWithCommandAndEnv(name, command, rawEnv string) error {
	env := parseEnvPairs(rawEnv)
	if err := cli.RunAdd(s.cfg, s.store, name, types.MCPServer{
		Transport: types.TransportStdio,
		Command:   command,
		Env:       env,
	}, nil, nil); err != nil {
		return fmt.Errorf("add server %s: %w", name, err)
	}
	return nil
}

func (s *scenarioCtx) iAddAStdioServerWithCommandAndArgs(name, command, rawArgs string) error {
	args := splitArgs(rawArgs)
	s.lastErr = cli.RunAdd(s.cfg, s.store, name, types.MCPServer{
		Transport: types.TransportStdio,
		Command:   command,
		Args:      args,
	}, nil, nil)
	return nil
}

func (s *scenarioCtx) iAddAStdioServerWithCommandAndNoArgs(name, command string) error {
	s.lastErr = cli.RunAdd(s.cfg, s.store, name, types.MCPServer{
		Transport: types.TransportStdio,
		Command:   command,
	}, nil, nil)
	return nil
}

func (s *scenarioCtx) iAddAnSSEServerWithURL(name, url string) error {
	s.lastErr = cli.RunAdd(s.cfg, s.store, name, types.MCPServer{
		Transport: types.TransportSSE,
		URL:       url,
	}, nil, nil)
	return nil
}

func (s *scenarioCtx) iAddAStdioServerWithCommandAndEnv(name, command, rawEnv string) error {
	env := parseEnvPairs(rawEnv)
	s.lastErr = cli.RunAdd(s.cfg, s.store, name, types.MCPServer{
		Transport: types.TransportStdio,
		Command:   command,
		Env:       env,
	}, nil, nil)
	return nil
}

func (s *scenarioCtx) iTryToAddAStdioServerWithNoCommand(name string) error {
	s.lastErr = cli.RunAdd(s.cfg, s.store, name, types.MCPServer{
		Transport: types.TransportStdio,
		// no Command
	}, nil, nil)
	return nil
}

// anSSEServerWithURL is a Given-form step (no error capture) for SSE servers.
func (s *scenarioCtx) anSSEServerWithURL(name, url string) error {
	if err := cli.RunAdd(s.cfg, s.store, name, types.MCPServer{
		Transport: types.TransportSSE,
		URL:       url,
	}, nil, nil); err != nil {
		return fmt.Errorf("add SSE server %s: %w", name, err)
	}
	return nil
}

// iAddAnMCPServerWithCommandAndArgs is a When-clause alias used in sync scenarios.
func (s *scenarioCtx) iAddAnMCPServerWithCommandAndArgs(name, command, rawArgs string) error {
	args := splitArgs(rawArgs)
	if err := cli.RunAdd(s.cfg, s.store, name, types.MCPServer{
		Transport: types.TransportStdio,
		Command:   command,
		Args:      args,
	}, nil, nil); err != nil {
		return fmt.Errorf("add server %s: %w", name, err)
	}
	return nil
}

// ── remove steps ──────────────────────────────────────────────────────────────

func (s *scenarioCtx) iRemoveTheServer(name string) error {
	if err := cli.RunRemove(s.cfg, s.store, name, nil, nil); err != nil {
		return fmt.Errorf("remove server %s: %w", name, err)
	}
	return nil
}

func (s *scenarioCtx) iTryToRemoveTheServer(name string) error {
	s.lastErr = cli.RunRemove(s.cfg, s.store, name, nil, nil)
	return nil
}
