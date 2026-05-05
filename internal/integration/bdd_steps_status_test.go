package integration_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pantheon-org/iris/internal/cli"
)

// ── status steps ──────────────────────────────────────────────────────────────

func (s *scenarioCtx) iRunStatus() error {
	s.output.Reset()
	if err := cli.RunStatus(s.root, s.cfg, s.reg, s.output, false); err != nil {
		return fmt.Errorf("status: %w", err)
	}
	return nil
}

func (s *scenarioCtx) iRunStatusWithJSONOutput() error {
	s.output.Reset()
	if err := cli.RunStatus(s.root, s.cfg, s.reg, s.output, true); err != nil {
		return fmt.Errorf("status json: %w", err)
	}
	return nil
}

func (s *scenarioCtx) iCorruptTheProviderConfigFile(filename string) error {
	path := filepath.Join(s.root, filename)
	if err := os.WriteFile(path, []byte(`{"mcpServers":{}}`), 0644); err != nil {
		return fmt.Errorf("corrupt config: %w", err)
	}
	return nil
}

// ── JSON status assertions ─────────────────────────────────────────────────────

func (s *scenarioCtx) theJSONStatusOutputHasAProvidersArray() error {
	var out cli.StatusOutput
	if err := json.Unmarshal(s.output.Bytes(), &out); err != nil {
		return fmt.Errorf("parse status JSON: %w\nraw: %s", err, s.output.String())
	}
	if out.Providers == nil {
		return fmt.Errorf("JSON status output missing \"providers\" key")
	}
	return nil
}

func (s *scenarioCtx) theJSONStatusProvidersContainEntryForProviderWithStatus(provider, status string) error {
	var out cli.StatusOutput
	if err := json.Unmarshal(s.output.Bytes(), &out); err != nil {
		return fmt.Errorf("parse status JSON: %w", err)
	}
	for _, e := range out.Providers {
		if e.Provider == provider && e.Status == status {
			return nil
		}
	}
	return fmt.Errorf("no entry with provider=%q status=%q in: %s", provider, status, s.output.String())
}

// theStatusOutputContainsProviderWithStatus scans text output for provider+status.
// It performs an exact field match on the provider name (first whitespace-delimited
// token on each line) to avoid "claude" matching "claude-desktop" lines.
func (s *scenarioCtx) theStatusOutputContainsProviderWithStatus(provider, status string) error {
	out := s.output.String()
	for _, line := range strings.Split(out, "\n") {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		if fields[0] == provider && strings.Contains(line, status) {
			return nil
		}
	}
	return fmt.Errorf("no line with provider=%q and status=%q in:\n%s", provider, status, out)
}
