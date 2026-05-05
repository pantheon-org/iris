package integration_test

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/pantheon-org/iris/internal/config"
	"github.com/pantheon-org/iris/internal/ierrors"
)

// ── reload steps ──────────────────────────────────────────────────────────────

func (s *scenarioCtx) iReloadTheConfigFromDisk() error {
	store2, err := config.NewStore(s.storePath)
	if err != nil {
		return fmt.Errorf("NewStore reload: %w", err)
	}
	loaded, err := store2.Load()
	if err != nil {
		return fmt.Errorf("Load: %w", err)
	}
	s.reloadedCfg = loaded
	return nil
}

// ── error assertions ──────────────────────────────────────────────────────────

func (s *scenarioCtx) theLastErrorWraps(substr string) error {
	if s.lastErr == nil {
		return fmt.Errorf("expected an error wrapping %q, but got nil", substr)
	}
	if !strings.Contains(s.lastErr.Error(), substr) {
		return fmt.Errorf("expected error to contain %q, got: %v", substr, s.lastErr)
	}
	return nil
}

func (s *scenarioCtx) theLastErrorIsErrServerNotFound() error {
	if s.lastErr == nil {
		return fmt.Errorf("expected ErrServerNotFound, got nil")
	}
	if !errors.Is(s.lastErr, ierrors.ErrServerNotFound) {
		return fmt.Errorf("expected ErrServerNotFound, got: %v", s.lastErr)
	}
	return nil
}

// ── in-memory config assertions ───────────────────────────────────────────────

func (s *scenarioCtx) theConfigContainsServerWithCommand(name, command string) error {
	srv, ok := s.cfg.Servers[name]
	if !ok {
		return fmt.Errorf("missing server %q", name)
	}
	if srv.Command != command {
		return fmt.Errorf("server %q: expected command %q, got %q", name, command, srv.Command)
	}
	return nil
}

func (s *scenarioCtx) theConfigContainsServerWithTransport(name, transport string) error {
	srv, ok := s.cfg.Servers[name]
	if !ok {
		return fmt.Errorf("missing server %q", name)
	}
	if string(srv.Transport) != transport {
		return fmt.Errorf("server %q: expected transport %q, got %q", name, transport, srv.Transport)
	}
	return nil
}

func (s *scenarioCtx) theConfigContainsServerWithEnvVar(name, key, value string) error {
	srv, ok := s.cfg.Servers[name]
	if !ok {
		return fmt.Errorf("missing server %q", name)
	}
	got, exists := srv.Env[key]
	if !exists {
		return fmt.Errorf("server %q: env var %q not set", name, key)
	}
	if got != value {
		return fmt.Errorf("server %q: env[%q] = %q, want %q", name, key, got, value)
	}
	return nil
}

func (s *scenarioCtx) theConfigHasExactlyNServers(n int) error {
	if len(s.cfg.Servers) != n {
		return fmt.Errorf("expected %d servers in cfg, got %d", n, len(s.cfg.Servers))
	}
	return nil
}

func (s *scenarioCtx) theIrisConfigFileExistsOnDisk() error {
	if _, err := os.Stat(s.storePath); err != nil {
		return fmt.Errorf("iris config %s: %w", s.storePath, err)
	}
	return nil
}

// ── reloaded config assertions ────────────────────────────────────────────────

func (s *scenarioCtx) theConfigContainsNServers(n int) error {
	if len(s.reloadedCfg.Servers) != n {
		return fmt.Errorf("expected %d servers, got %d", n, len(s.reloadedCfg.Servers))
	}
	return nil
}

func (s *scenarioCtx) theConfigDoesNotContainServer(name string) error {
	if _, ok := s.reloadedCfg.Servers[name]; ok {
		return fmt.Errorf("server %q should have been removed", name)
	}
	return nil
}

func (s *scenarioCtx) theReloadedConfigContainsServerWithCommand(name, command string) error {
	srv, ok := s.reloadedCfg.Servers[name]
	if !ok {
		return fmt.Errorf("missing server %q", name)
	}
	if srv.Command != command {
		return fmt.Errorf("server %q: expected command %q, got %q", name, command, srv.Command)
	}
	return nil
}

// ── env var mutation ──────────────────────────────────────────────────────────

func (s *scenarioCtx) theServerHasEnvVarSetTo(name, key, value string) error {
	srv, ok := s.cfg.Servers[name]
	if !ok {
		return fmt.Errorf("server %q not found", name)
	}
	if srv.Env == nil {
		srv.Env = make(map[string]string)
	}
	srv.Env[key] = value
	s.cfg.Servers[name] = srv
	if err := s.store.Save(s.cfg); err != nil {
		return fmt.Errorf("save config: %w", err)
	}
	return nil
}

// ── output text assertions ────────────────────────────────────────────────────

func (s *scenarioCtx) theOutputContains(substr string) error {
	if !strings.Contains(s.output.String(), substr) {
		return fmt.Errorf("expected output to contain %q, got:\n%s", substr, s.output.String())
	}
	return nil
}

func (s *scenarioCtx) theOutputLinesAppearInOrder(rawNames string) error {
	names := strings.Split(rawNames, ",")
	out := s.output.String()
	lastIdx := -1
	for _, name := range names {
		idx := strings.Index(out, name)
		if idx == -1 {
			return fmt.Errorf("name %q not found in output:\n%s", name, out)
		}
		if idx <= lastIdx {
			return fmt.Errorf("name %q appears before %q in output, expected reverse order:\n%s", names[lastIdx], name, out)
		}
		lastIdx = idx
	}
	return nil
}
