package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/cucumber/godog"

	"github.com/pantheon-org/iris/internal/cli"
	"github.com/pantheon-org/iris/internal/config"
	"github.com/pantheon-org/iris/internal/ierrors"
	"github.com/pantheon-org/iris/internal/providers"
	"github.com/pantheon-org/iris/internal/registry"
	irisync "github.com/pantheon-org/iris/internal/sync"
	"github.com/pantheon-org/iris/internal/types"
)

// scenarioCtx holds per-scenario mutable state.
type scenarioCtx struct {
	root      string
	storePath string
	store     *config.Store
	cfg       *types.IrisConfig
	reg       *registry.Registry

	// captured output / results
	lastErr     error
	output      *bytes.Buffer
	syncResults []irisync.SyncResult
	reloadedCfg *types.IrisConfig
}

func newScenarioCtx(root string) *scenarioCtx {
	return &scenarioCtx{
		root:   root,
		output: &bytes.Buffer{},
	}
}

// buildReg constructs a full registry with all 14 providers, paths pinned under root.
func buildReg(root string) *registry.Registry {
	googleGeminiPath := filepath.Join(root, "gemini-settings.json")
	codexPath := filepath.Join(root, "codex-config.toml")
	claudeDesktopPath := filepath.Join(root, "claude-desktop-config.json")
	windsurfPath := filepath.Join(root, "windsurf-config.json")
	zedPath := filepath.Join(root, "zed-settings.json")
	warpPath := filepath.Join(root, "warp-mcp.json")
	kimiPath := filepath.Join(root, "kimi-settings.json")
	mistralVibePath := filepath.Join(root, "mistral-vibe-config.toml")

	reg := registry.NewRegistry()
	reg.Register(providers.NewClaudeCodeProvider())
	reg.Register(providers.NewClaudeDesktopProviderWithPath(claudeDesktopPath))
	reg.Register(providers.NewGoogleGeminiProviderWithPath(googleGeminiPath))
	reg.Register(providers.NewOpenCodeProvider())
	reg.Register(providers.NewOpenaiCodexProviderWithPath(codexPath))
	reg.Register(providers.NewCursorProvider())
	reg.Register(providers.NewWindsurfProviderWithPath(windsurfPath))
	reg.Register(providers.NewVSCodeCopilotProvider())
	reg.Register(providers.NewZedProviderWithPath(zedPath))
	reg.Register(providers.NewQwenProvider())
	reg.Register(providers.NewWarpProviderWithPath(warpPath))
	reg.Register(providers.NewKimiProviderWithPath(kimiPath))
	reg.Register(providers.NewMistralVibeProviderWithPath(mistralVibePath))
	reg.Register(providers.NewIntelliJProvider())
	return reg
}

// ── shared setup ──────────────────────────────────────────────────────────────

func (s *scenarioCtx) aCleanWorkspace() error {
	s.storePath = filepath.Join(s.root, ".iris.json")
	store, err := config.NewStore(s.storePath)
	if err != nil {
		return fmt.Errorf("NewStore: %w", err)
	}
	s.store = store
	s.cfg = &types.IrisConfig{
		Version: 1,
		Servers: make(map[string]types.MCPServer),
	}
	s.reg = buildReg(s.root)
	return nil
}

// ── add steps ─────────────────────────────────────────────────────────────────

func (s *scenarioCtx) anMCPServerWithCommandAndArgs(name, command, rawArgs string) error {
	args := splitArgs(rawArgs)
	if err := cli.RunAdd(s.cfg, s.store, name, types.MCPServer{
		Transport: types.TransportStdio,
		Command:   command,
		Args:      args,
	}); err != nil {
		return fmt.Errorf("add server %s: %w", name, err)
	}
	return nil
}

func (s *scenarioCtx) anMCPServerWithCommandAndNoArgs(name, command string) error {
	if err := cli.RunAdd(s.cfg, s.store, name, types.MCPServer{
		Transport: types.TransportStdio,
		Command:   command,
	}); err != nil {
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
	}); err != nil {
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
	})
	return nil
}

func (s *scenarioCtx) iAddAStdioServerWithCommandAndNoArgs(name, command string) error {
	s.lastErr = cli.RunAdd(s.cfg, s.store, name, types.MCPServer{
		Transport: types.TransportStdio,
		Command:   command,
	})
	return nil
}

func (s *scenarioCtx) iAddAnSSEServerWithURL(name, url string) error {
	s.lastErr = cli.RunAdd(s.cfg, s.store, name, types.MCPServer{
		Transport: types.TransportSSE,
		URL:       url,
	})
	return nil
}

func (s *scenarioCtx) iAddAStdioServerWithCommandAndEnv(name, command, rawEnv string) error {
	env := parseEnvPairs(rawEnv)
	s.lastErr = cli.RunAdd(s.cfg, s.store, name, types.MCPServer{
		Transport: types.TransportStdio,
		Command:   command,
		Env:       env,
	})
	return nil
}

func (s *scenarioCtx) iTryToAddAStdioServerWithNoCommand(name string) error {
	s.lastErr = cli.RunAdd(s.cfg, s.store, name, types.MCPServer{
		Transport: types.TransportStdio,
		// no Command
	})
	return nil
}

// anSSEServerWithURL is a Given-form step (no error capture) for SSE servers.
func (s *scenarioCtx) anSSEServerWithURL(name, url string) error {
	if err := cli.RunAdd(s.cfg, s.store, name, types.MCPServer{
		Transport: types.TransportSSE,
		URL:       url,
	}); err != nil {
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
	}); err != nil {
		return fmt.Errorf("add server %s: %w", name, err)
	}
	return nil
}

// ── remove steps ──────────────────────────────────────────────────────────────

func (s *scenarioCtx) iRemoveTheServer(name string) error {
	if err := cli.RunRemove(s.cfg, s.store, name); err != nil {
		return fmt.Errorf("remove server %s: %w", name, err)
	}
	return nil
}

func (s *scenarioCtx) iTryToRemoveTheServer(name string) error {
	s.lastErr = cli.RunRemove(s.cfg, s.store, name)
	return nil
}

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

// ── sync steps ────────────────────────────────────────────────────────────────

func (s *scenarioCtx) iSyncToAllProviders() error {
	s.reg = buildReg(s.root)
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

// ── round-trip steps ──────────────────────────────────────────────────────────

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

// ── init steps ────────────────────────────────────────────────────────────────

func (s *scenarioCtx) iRunInit() error {
	s.output.Reset()
	if err := cli.RunInitNonInteractive(s.store, s.output); err != nil {
		return fmt.Errorf("init: %w", err)
	}
	return nil
}

func (s *scenarioCtx) theIrisConfigAlreadyExistsWithOneServer() error {
	if err := cli.RunAdd(s.cfg, s.store, "existing", types.MCPServer{
		Transport: types.TransportStdio,
		Command:   "existing-cmd",
	}); err != nil {
		return fmt.Errorf("add existing: %w", err)
	}
	return nil
}

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

// ── assertion steps ───────────────────────────────────────────────────────────

// --- error assertions ---

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

// --- in-memory config assertions ---

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

// --- reloaded config assertions ---

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

// --- file-level provider assertions ---

func (s *scenarioCtx) theProviderConfigFileExists(filename string) error {
	path := filepath.Join(s.root, filename)
	if _, err := os.Stat(path); err != nil {
		return fmt.Errorf("expected file %s to exist: %w", path, err)
	}
	return nil
}

func (s *scenarioCtx) theJSONProviderFileContainsServersUnderKey(filename, rawServers, key string) error {
	path := filepath.Join(s.root, filename)
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}
	var doc map[string]map[string]json.RawMessage
	if err := json.Unmarshal(data, &doc); err != nil {
		return fmt.Errorf("parse %s: %w", path, err)
	}
	section, ok := doc[key]
	if !ok {
		return fmt.Errorf("%s: missing key %q", filename, key)
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

// --- copilot-specific assertions ---

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

// --- env var mutations ---

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

// theJSONProviderServerHasEnvVar asserts that a named server in a JSON provider
// file (under sectionKey → serverName → env) contains the given env key.
// Regex capture order: filename, serverName, sectionKey, envKey.
func (s *scenarioCtx) theJSONProviderServerHasEnvVar(filename, serverName, sectionKey, envKey string) error {
	path := filepath.Join(s.root, filename)
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}
	var doc map[string]map[string]struct {
		Env map[string]string `json:"env"`
	}
	if err := json.Unmarshal(data, &doc); err != nil {
		return fmt.Errorf("parse %s: %w", path, err)
	}
	section, ok := doc[sectionKey]
	if !ok {
		return fmt.Errorf("%s: missing key %q", filename, sectionKey)
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

// --- sync result assertions ---

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

// --- output text assertions ---

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

// --- JSON list assertions ---

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

// --- JSON sync assertions ---

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

// --- JSON status assertions ---

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

// --- init assertions ---

func (s *scenarioCtx) theIrisConfigFileIsValidJSONWithVersion1() error {
	data, err := os.ReadFile(s.storePath)
	if err != nil {
		return fmt.Errorf("read iris config: %w", err)
	}
	var cfg map[string]any
	if err := json.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("parse iris config JSON: %w", err)
	}
	v, ok := cfg["version"]
	if !ok {
		return fmt.Errorf("iris config missing \"version\" field")
	}
	// JSON numbers decode as float64.
	if v.(float64) != 1 {
		return fmt.Errorf("expected version 1, got %v", v)
	}
	return nil
}

// theJSONProviderServerHasField asserts a specific field exists in a server entry
// inside a JSON provider file, under doc[key][serverName].
func (s *scenarioCtx) theJSONProviderServerHasField(filename, serverName, key, field string) error {
	path := filepath.Join(s.root, filename)
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}
	var doc map[string]map[string]json.RawMessage
	if err := json.Unmarshal(data, &doc); err != nil {
		return fmt.Errorf("parse %s: %w", path, err)
	}
	section, ok := doc[key]
	if !ok {
		return fmt.Errorf("%s: missing key %q", filename, key)
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

// ── helpers ───────────────────────────────────────────────────────────────────

func splitArgs(raw string) []string {
	if raw == "" {
		return nil
	}
	return strings.Split(raw, ",")
}

func parseEnvPairs(raw string) map[string]string {
	env := make(map[string]string)
	for _, pair := range strings.Split(raw, ",") {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) == 2 {
			env[parts[0]] = parts[1]
		}
	}
	return env
}



// ── suite wiring ──────────────────────────────────────────────────────────────

func initializeScenario(t *testing.T) func(ctx *godog.ScenarioContext) {
	t.Helper()
	return func(sc *godog.ScenarioContext) {
		root := t.TempDir()
		s := newScenarioCtx(root)

		// setup
		sc.Step(`^a clean workspace$`, s.aCleanWorkspace)
		sc.Step(`^the iris config already exists with one server$`, s.theIrisConfigAlreadyExistsWithOneServer)

		// add — Given form (no error capture)
		sc.Step(`^an MCP server "([^"]+)" with command "([^"]+)" and args "([^"]+)"$`, s.anMCPServerWithCommandAndArgs)
		sc.Step(`^an MCP server "([^"]+)" with command "([^"]+)" and no args$`, s.anMCPServerWithCommandAndNoArgs)
		sc.Step(`^an MCP server "([^"]+)" with command "([^"]+)" and env "([^"]+)"$`, s.anMCPServerWithCommandAndEnv)
		sc.Step(`^an SSE server "([^"]+)" with URL "([^"]+)"$`, s.anSSEServerWithURL)
		sc.Step(`^the server "([^"]+)" has env var "([^"]+)" set to "([^"]+)"$`, s.theServerHasEnvVarSetTo)

		// add — When form (error capture)
		sc.Step(`^I add a stdio server "([^"]+)" with command "([^"]+)" and args "([^"]+)"$`, s.iAddAStdioServerWithCommandAndArgs)
		sc.Step(`^I add a stdio server "([^"]+)" with command "([^"]+)" and no args$`, s.iAddAStdioServerWithCommandAndNoArgs)
		sc.Step(`^I add an SSE server "([^"]+)" with url "([^"]+)"$`, s.iAddAnSSEServerWithURL)
		sc.Step(`^I add a stdio server "([^"]+)" with command "([^"]+)" and env "([^"]+)"$`, s.iAddAStdioServerWithCommandAndEnv)
		sc.Step(`^I try to add a stdio server "([^"]+)" with no command$`, s.iTryToAddAStdioServerWithNoCommand)
		sc.Step(`^I add an MCP server "([^"]+)" with command "([^"]+)" and args "([^"]+)"$`, s.iAddAnMCPServerWithCommandAndArgs)

		// remove
		sc.Step(`^I remove the server "([^"]+)"$`, s.iRemoveTheServer)
		sc.Step(`^I try to remove the server "([^"]+)"$`, s.iTryToRemoveTheServer)

		// list
		sc.Step(`^I run list$`, s.iRunList)
		sc.Step(`^I run list with JSON output$`, s.iRunListWithJSONOutput)

		// sync
		sc.Step(`^I sync to all providers$`, s.iSyncToAllProviders)
		sc.Step(`^I sync to all providers again$`, s.iSyncToAllProvidersAgain)
		sc.Step(`^I sync to all providers with JSON output$`, s.iSyncToAllProvidersWithJSONOutput)

		// status
		sc.Step(`^I run status$`, s.iRunStatus)
		sc.Step(`^I run status with JSON output$`, s.iRunStatusWithJSONOutput)
		sc.Step(`^I corrupt the provider config file "([^"]+)"$`, s.iCorruptTheProviderConfigFile)

		// init
		sc.Step(`^I run init$`, s.iRunInit)

		// reload
		sc.Step(`^I reload the config from disk$`, s.iReloadTheConfigFromDisk)

		// assertions — errors
		sc.Step(`^the last error wraps "([^"]+)"$`, s.theLastErrorWraps)
		sc.Step(`^the last error is ErrServerNotFound$`, s.theLastErrorIsErrServerNotFound)

		// assertions — in-memory config
		sc.Step(`^the config contains server "([^"]+)" with command "([^"]+)"$`, s.theConfigContainsServerWithCommand)
		sc.Step(`^the config contains server "([^"]+)" with transport "([^"]+)"$`, s.theConfigContainsServerWithTransport)
		sc.Step(`^the config contains server "([^"]+)" with env var "([^"]+)" equal to "([^"]+)"$`, s.theConfigContainsServerWithEnvVar)
		sc.Step(`^the config has exactly (\d+) server$`, s.theConfigHasExactlyNServers)
		sc.Step(`^the iris config file exists on disk$`, s.theIrisConfigFileExistsOnDisk)

		// assertions — reloaded config
		sc.Step(`^the config contains (\d+) servers$`, s.theConfigContainsNServers)
		sc.Step(`^the config does not contain server "([^"]+)"$`, s.theConfigDoesNotContainServer)
		sc.Step(`^the config contains server "([^"]+)" with command "([^"]+)"$`, s.theReloadedConfigContainsServerWithCommand)

		// assertions — provider files
		sc.Step(`^the provider config file "([^"]+)" exists$`, s.theProviderConfigFileExists)
		sc.Step(`^the JSON provider file "([^"]+)" contains servers "([^"]+)" under key "([^"]+)"$`, s.theJSONProviderFileContainsServersUnderKey)
		sc.Step(`^the JSON provider file "([^"]+)" server "([^"]+)" under key "([^"]+)" has field "([^"]+)"$`, s.theJSONProviderServerHasField)
		sc.Step(`^the opencode provider file "([^"]+)" contains servers "([^"]+)"$`, s.theOpencodeProviderFileContainsServers)
		sc.Step(`^the opencode server "([^"]+)" in file "([^"]+)" has correct field format$`, s.theOpencodeServerHasCorrectFieldFormat)
		sc.Step(`^the TOML provider file "([^"]+)" contains servers "([^"]+)"$`, s.theTOMLProviderFileContainsServers)
		sc.Step(`^the zed provider file "([^"]+)" contains servers "([^"]+)"$`, s.theZedProviderFileContainsServers)
		sc.Step(`^the TOML mistral provider file "([^"]+)" contains servers "([^"]+)"$`, s.theTOMLMistralProviderFileContainsServers)
		sc.Step(`^the copilot server "([^"]+)" in file "([^"]+)" does not have field "([^"]+)"$`, s.theCopilotServerDoesNotHaveField)
		sc.Step(`^the JSON provider file "([^"]+)" server "([^"]+)" under key "([^"]+)" has env var "([^"]+)"$`, s.theJSONProviderServerHasEnvVar)

		// assertions — sync results
		sc.Step(`^all providers report status "([^"]+)"$`, s.allProvidersReportStatus)

		// assertions — text output
		sc.Step(`^the output contains "([^"]+)"$`, s.theOutputContains)
		sc.Step(`^the output lines appear in order "([^"]+)"$`, s.theOutputLinesAppearInOrder)

		// assertions — JSON list
		sc.Step(`^the JSON output has a "servers" array$`, s.theJSONOutputHasAServersArray)
		sc.Step(`^the JSON servers array contains an entry with name "([^"]+)" and command "([^"]+)"$`, s.theJSONServersArrayContainsEntryWithNameAndCommand)

		// assertions — JSON sync
		sc.Step(`^the JSON sync output has a "results" array$`, s.theJSONSyncOutputHasAResultsArray)
		sc.Step(`^the JSON sync results contain an entry for provider "([^"]+)" with status "([^"]+)"$`, s.theJSONSyncResultsContainEntryForProviderWithStatus)

		// assertions — JSON status
		sc.Step(`^the JSON status output has a "providers" array$`, s.theJSONStatusOutputHasAProvidersArray)
		sc.Step(`^the JSON status providers contain an entry for provider "([^"]+)" with status "([^"]+)"$`, s.theJSONStatusProvidersContainEntryForProviderWithStatus)

		// assertions — text status
		sc.Step(`^the status output contains provider "([^"]+)" with status "([^"]+)"$`, s.theStatusOutputContainsProviderWithStatus)

		// assertions — init
		sc.Step(`^the iris config file is valid JSON with version 1$`, s.theIrisConfigFileIsValidJSONWithVersion1)
	}
}

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		Name:                 "iris integration",
		TestSuiteInitializer: nil,
		ScenarioInitializer:  initializeScenario(t),
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t,
		},
	}
	if suite.Run() != 0 {
		t.Fatal("BDD integration tests failed")
	}
}

// Keep the context import used by godog internally.
var _ context.Context
