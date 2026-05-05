package integration_test

import (
	"encoding/json"
	"fmt"

	"github.com/pantheon-org/iris/internal/cli"
)

// ── list steps ────────────────────────────────────────────────────────────────

func (s *scenarioCtx) iRunList() error {
	s.output.Reset()
	if err := cli.RunList(s.cfg, s.output, false); err != nil {
		return fmt.Errorf("list: %w", err)
	}
	return nil
}

func (s *scenarioCtx) iRunListWithJSONOutput() error {
	s.output.Reset()
	if err := cli.RunList(s.cfg, s.output, true); err != nil {
		return fmt.Errorf("list json: %w", err)
	}
	return nil
}

// ── JSON list assertions ───────────────────────────────────────────────────────

func (s *scenarioCtx) theJSONOutputHasAServersArray() error {
	var out cli.ListOutput
	if err := json.Unmarshal(s.output.Bytes(), &out); err != nil {
		return fmt.Errorf("parse list JSON: %w\nraw: %s", err, s.output.String())
	}
	if out.Servers == nil {
		return fmt.Errorf("JSON output missing \"servers\" key")
	}
	return nil
}

func (s *scenarioCtx) theJSONServersArrayContainsEntryWithNameAndCommand(name, command string) error {
	var out cli.ListOutput
	if err := json.Unmarshal(s.output.Bytes(), &out); err != nil {
		return fmt.Errorf("parse list JSON: %w", err)
	}
	for _, e := range out.Servers {
		if e.Name == name && e.Command == command {
			return nil
		}
	}
	return fmt.Errorf("no entry with name=%q command=%q in: %s", name, command, s.output.String())
}
