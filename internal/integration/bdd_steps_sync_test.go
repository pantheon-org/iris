package integration_test

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"

	"github.com/pantheon-org/iris/internal/cli"
	irisync "github.com/pantheon-org/iris/internal/sync"
)

// ── sync steps ────────────────────────────────────────────────────────────────

func (s *scenarioCtx) iSyncToAllProviders() error {
	s.reg = buildReg(s.root)
	if len(s.cfg.Providers) > 0 {
		filtered, err := s.reg.Filter(s.cfg.Providers)
		if err != nil {
			return fmt.Errorf("filter providers: %w", err)
		}
		s.reg = filtered
	}
	s.syncResults = irisync.SyncAllProviders(s.root, s.reg, s.cfg.Servers)
	return nil
}

func (s *scenarioCtx) iSyncToAllProvidersAgain() error {
	s.reg = buildReg(s.root)
	if err := cli.RunSync(s.root, s.cfg, s.reg, io.Discard, false); err != nil {
		return fmt.Errorf("sync: %w", err)
	}
	s.syncResults = irisync.SyncAllProviders(s.root, s.reg, s.cfg.Servers)
	return nil
}

func (s *scenarioCtx) iSyncToAllProvidersWithJSONOutput() error {
	s.output.Reset()
	if err := cli.RunSync(s.root, s.cfg, s.reg, s.output, true); err != nil {
		return fmt.Errorf("sync json: %w", err)
	}
	return nil
}

// ── provider file assertions ───────────────────────────────────────────────────

func (s *scenarioCtx) theProviderConfigFileExists(filename string) error {
	path := filepath.Join(s.root, filename)
	if _, err := os.Stat(path); err != nil {
		return fmt.Errorf("expected file %s to exist: %w", path, err)
	}
	return nil
}

func (s *scenarioCtx) theProviderConfigFileDoesNotExist(filename string) error {
	path := filepath.Join(s.root, filename)
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("expected file %s to not exist, but it does", path)
	}
	return nil
}

func (s *scenarioCtx) theJSONProviderFileContainsServersUnderKey(filename, rawServers, key string) error {
	path := filepath.Join(s.root, filename)
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}
	var doc map[string]json.RawMessage
	if err := json.Unmarshal(data, &doc); err != nil {
		return fmt.Errorf("parse %s: %w", path, err)
	}
	raw, ok := doc[key]
	if !ok {
		return fmt.Errorf("%s: missing key %q", filename, key)
	}
	var section map[string]json.RawMessage
	if err := json.Unmarshal(raw, &section); err != nil {
		return fmt.Errorf("%s: key %q is not an object: %w", filename, key, err)
	}
	for _, name := range splitArgs(rawServers) {
		if _, ok := section[name]; !ok {
			return fmt.Errorf("%s: missing server %q under %q", filename, name, key)
		}
	}
	return nil
}

func (s *scenarioCtx) theOpencodeProviderFileContainsServers(filename, rawServers string) error {
	path := filepath.Join(s.root, filename)
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}
	var doc struct {
		MCP map[string]json.RawMessage `json:"mcp"`
	}
	if err := json.Unmarshal(data, &doc); err != nil {
		return fmt.Errorf("parse %s: %w", path, err)
	}
	for _, name := range splitArgs(rawServers) {
		if _, ok := doc.MCP[name]; !ok {
			return fmt.Errorf("opencode.json: missing server %q", name)
		}
	}
	return nil
}

func (s *scenarioCtx) theOpencodeServerHasCorrectFieldFormat(serverName, filename string) error {
	path := filepath.Join(s.root, filename)
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}
	var doc struct {
		MCP map[string]struct {
			Command     json.RawMessage   `json:"command"`
			Type        string            `json:"type"`
			Environment map[string]string `json:"environment"`
			Env         map[string]string `json:"env"`
		} `json:"mcp"`
	}
	if err := json.Unmarshal(data, &doc); err != nil {
		return fmt.Errorf("parse %s: %w", path, err)
	}
	entry, ok := doc.MCP[serverName]
	if !ok {
		return fmt.Errorf("opencode: missing server %q", serverName)
	}
	// command must be a JSON array, not a string
	var cmdSlice []string
	if err := json.Unmarshal(entry.Command, &cmdSlice); err != nil {
		return fmt.Errorf("opencode server %q: command must be a JSON array, got: %s", serverName, entry.Command)
	}
	// type must be "local" (not "stdio")
	if entry.Type != "local" {
		return fmt.Errorf("opencode server %q: type must be %q, got %q", serverName, "local", entry.Type)
	}
	// env key must not exist; environment key used instead
	if entry.Env != nil {
		return fmt.Errorf("opencode server %q: must use \"environment\" key, not \"env\"", serverName)
	}
	return nil
}

func (s *scenarioCtx) theTOMLProviderFileContainsServers(filename, rawServers string) error {
	path := filepath.Join(s.root, filename)
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}
	var doc struct {
		MCPServers map[string]map[string]any `toml:"mcp_servers"`
	}
	if _, err := toml.Decode(string(data), &doc); err != nil {
		return fmt.Errorf("parse %s: %w", path, err)
	}
	for _, name := range splitArgs(rawServers) {
		if _, ok := doc.MCPServers[name]; !ok {
			return fmt.Errorf("codex config: missing server %q", name)
		}
	}
	return nil
}

func (s *scenarioCtx) theZedProviderFileContainsServers(filename, rawServers string) error {
	path := filepath.Join(s.root, filename)
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}
	var doc struct {
		ContextServers map[string]json.RawMessage `json:"context_servers"`
	}
	if err := json.Unmarshal(data, &doc); err != nil {
		return fmt.Errorf("parse zed %s: %w", path, err)
	}
	for _, name := range splitArgs(rawServers) {
		if _, ok := doc.ContextServers[name]; !ok {
			return fmt.Errorf("zed config: missing server %q", name)
		}
	}
	return nil
}

func (s *scenarioCtx) theTOMLMistralProviderFileContainsServers(filename, rawServers string) error {
	path := filepath.Join(s.root, filename)
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}
	var doc struct {
		MCPServers []struct {
			Name string `toml:"name"`
		} `toml:"mcp_servers"`
	}
	if _, err := toml.Decode(string(data), &doc); err != nil {
		return fmt.Errorf("parse mistral %s: %w", path, err)
	}
	names := make(map[string]bool, len(doc.MCPServers))
	for _, e := range doc.MCPServers {
		names[e.Name] = true
	}
	for _, name := range splitArgs(rawServers) {
		if !names[name] {
			return fmt.Errorf("mistral config: missing server %q", name)
		}
	}
	return nil
}

func (s *scenarioCtx) theCopilotServerDoesNotHaveField(serverName, filename, field string) error {
	path := filepath.Join(s.root, filename)
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}
	var root struct {
		Servers map[string]map[string]json.RawMessage `json:"servers"`
	}
	if err := json.Unmarshal(data, &root); err != nil {
		return fmt.Errorf("parse copilot %s: %w", path, err)
	}
	srv, ok := root.Servers[serverName]
	if !ok {
		return fmt.Errorf("copilot %s: server %q not found", filename, serverName)
	}
	if _, exists := srv[field]; exists {
		return fmt.Errorf("copilot %s: server %q unexpectedly has field %q", filename, serverName, field)
	}
	return nil
}

func (s *scenarioCtx) theJSONProviderServerHasField(filename, serverName, key, field string) error {
	path := filepath.Join(s.root, filename)
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}
	var doc map[string]json.RawMessage
	if err := json.Unmarshal(data, &doc); err != nil {
		return fmt.Errorf("parse %s: %w", path, err)
	}
	rawSection, ok := doc[key]
	if !ok {
		return fmt.Errorf("%s: missing key %q", filename, key)
	}
	var section map[string]json.RawMessage
	if err := json.Unmarshal(rawSection, &section); err != nil {
		return fmt.Errorf("%s: key %q is not an object: %w", filename, key, err)
	}
	raw, ok := section[serverName]
	if !ok {
		return fmt.Errorf("%s: missing server %q under key %q", filename, serverName, key)
	}
	var entry map[string]json.RawMessage
	if err := json.Unmarshal(raw, &entry); err != nil {
		return fmt.Errorf("%s: parse server %q: %w", filename, serverName, err)
	}
	if _, ok := entry[field]; !ok {
		return fmt.Errorf("%s: server %q under %q has no field %q (got: %v)", filename, serverName, key, field, entry)
	}
	return nil
}

// theJSONProviderServerHasEnvVar asserts that a named server in a JSON provider
// file (under sectionKey → serverName → env) contains the given env key.
// Regex capture order: filename, serverName, sectionKey, envKey.
func (s *scenarioCtx) theJSONProviderServerHasEnvVar(filename, serverName, sectionKey, envKey string) error {
	path := filepath.Join(s.root, filename)
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}
	var doc map[string]json.RawMessage
	if err := json.Unmarshal(data, &doc); err != nil {
		return fmt.Errorf("parse %s: %w", path, err)
	}
	rawSection, ok := doc[sectionKey]
	if !ok {
		return fmt.Errorf("%s: missing key %q", filename, sectionKey)
	}
	var section map[string]struct {
		Env map[string]string `json:"env"`
	}
	if err := json.Unmarshal(rawSection, &section); err != nil {
		return fmt.Errorf("%s: key %q is not an object: %w", filename, sectionKey, err)
	}
	srv, ok := section[serverName]
	if !ok {
		return fmt.Errorf("%s: missing server %q under %q", filename, serverName, sectionKey)
	}
	if _, exists := srv.Env[envKey]; !exists {
		return fmt.Errorf("%s: server %q env missing key %q", filename, serverName, envKey)
	}
	return nil
}

func (s *scenarioCtx) theJSONProviderFileStillHasKey(filename, key string) error {
	path := filepath.Join(s.root, filename)
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}
	var doc map[string]json.RawMessage
	if err := json.Unmarshal(data, &doc); err != nil {
		return fmt.Errorf("parse %s: %w", path, err)
	}
	if _, ok := doc[key]; !ok {
		return fmt.Errorf("%s: key %q was not preserved", filename, key)
	}
	return nil
}

func (s *scenarioCtx) aProviderFileExistsWithExtraKey(filename, key, value string) error {
	path := filepath.Join(s.root, filename)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("mkdir: %w", err)
	}
	content := map[string]interface{}{key: value}
	data, err := json.MarshalIndent(content, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}

// ── sync result assertions ─────────────────────────────────────────────────────

func (s *scenarioCtx) allProvidersReportStatus(expectedStatus string) error {
	for _, r := range s.syncResults {
		if r.Err != nil {
			return fmt.Errorf("provider %s: unexpected error: %w", r.ProviderName, r.Err)
		}
		if string(r.Status) != expectedStatus {
			return fmt.Errorf("provider %s: expected status %q, got %q", r.ProviderName, expectedStatus, r.Status)
		}
	}
	return nil
}

// ── JSON sync assertions ───────────────────────────────────────────────────────

func (s *scenarioCtx) theJSONSyncOutputHasAResultsArray() error {
	var out cli.SyncOutput
	if err := json.Unmarshal(s.output.Bytes(), &out); err != nil {
		return fmt.Errorf("parse sync JSON: %w\nraw: %s", err, s.output.String())
	}
	if out.Results == nil {
		return fmt.Errorf("JSON sync output missing \"results\" key")
	}
	return nil
}

func (s *scenarioCtx) theJSONSyncResultsContainEntryForProviderWithStatus(provider, status string) error {
	var out cli.SyncOutput
	if err := json.Unmarshal(s.output.Bytes(), &out); err != nil {
		return fmt.Errorf("parse sync JSON: %w", err)
	}
	for _, e := range out.Results {
		if e.Provider == provider && e.Status == status {
			return nil
		}
	}
	return fmt.Errorf("no entry with provider=%q status=%q in: %s", provider, status, s.output.String())
}
