package integration_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pantheon-org/iris/internal/cli"
	"github.com/pantheon-org/iris/internal/config"
	"github.com/pantheon-org/iris/internal/providers"
	"github.com/pantheon-org/iris/internal/registry"
	"github.com/pantheon-org/iris/internal/types"
	"github.com/pantheon-org/iris/internal/wizard"
)

// ── non-interactive init ───────────────────────────────────────────────────────

func (s *scenarioCtx) iRunInit() error {
	s.output.Reset()
	if err := cli.RunInitNonInteractive(s.store, s.output); err != nil {
		return fmt.Errorf("init: %w", err)
	}
	return nil
}

// ── interactive init helpers ───────────────────────────────────────────────────

// irisConfigServerCount returns the number of servers in the on-disk iris config.
func (s *scenarioCtx) irisConfigServerCount() (int, error) {
	store2, err := config.NewStore(s.storePath)
	if err != nil {
		return 0, fmt.Errorf("NewStore: %w", err)
	}
	cfg, err := store2.Load()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return 0, nil
		}
		return 0, fmt.Errorf("load: %w", err)
	}
	return len(cfg.Servers), nil
}

// runInteractiveInit executes wizard.RunInit using a ScriptedRunner built from answers.
func (s *scenarioCtx) runInteractiveInit(answers []string) error {
	r := wizard.NewScriptedRunner(answers)
	if err := wizard.RunInit(r, s.root, s.store, s.reg); err != nil {
		return fmt.Errorf("RunInit: %w", err)
	}
	return nil
}

// buildIsolatedReg constructs a registry with only the providers needed for
// interactive-init BDD scenarios, all paths pinned under root so no real
// home-directory configs are read.
func buildIsolatedReg(root string) *registry.Registry {
	reg := registry.NewRegistry()
	reg.Register(providers.NewClaudeCodeProviderWithGlobalPath(filepath.Join(root, "claude-global.json")))
	reg.Register(providers.NewCursorProvider())
	reg.Register(providers.NewGoogleGeminiProviderWithPath(filepath.Join(root, "gemini-settings.json")))
	return reg
}

// ── provider setup helpers ────────────────────────────────────────────────────

func (s *scenarioCtx) noProviderConfigFilesExist() error {
	s.reg = buildIsolatedReg(s.root)
	return nil
}

func (s *scenarioCtx) aMalformedClaudeCodeProjectConfigExists() error {
	s.reg = buildIsolatedReg(s.root)
	if err := os.WriteFile(filepath.Join(s.root, ".mcp.json"), []byte(`{"mcpServers": {`), 0o600); err != nil {
		return fmt.Errorf("write malformed config: %w", err)
	}
	return nil
}

func (s *scenarioCtx) aClaudeCodeProjectConfigExistsWithServer(serverName, command, rawArgs string) error {
	s.reg = buildIsolatedReg(s.root)
	args := strings.Fields(rawArgs)
	entry := map[string]any{"command": command, "args": args, "type": "stdio"}
	doc := map[string]any{"mcpServers": map[string]any{serverName: entry}}
	data, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("marshal claude config: %w", err)
	}
	if err := os.WriteFile(filepath.Join(s.root, ".mcp.json"), data, 0o600); err != nil {
		return fmt.Errorf("write claude config: %w", err)
	}
	return nil
}

func (s *scenarioCtx) aCursorProjectConfigExistsWithServer(serverName, command, rawArgs string) error {
	// Does not reset s.reg — assumes a provider-setup step already set it via aClaudeCodeProjectConfigExistsWithServer
	// or noProviderConfigFilesExist. For standalone use it will rely on the isolated registry already set.
	args := strings.Fields(rawArgs)
	entry := map[string]any{"command": command, "args": args, "type": "stdio"}
	doc := map[string]any{"mcpServers": map[string]any{serverName: entry}}
	data, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("marshal cursor config: %w", err)
	}
	cursorDir := filepath.Join(s.root, ".cursor")
	if err := os.MkdirAll(cursorDir, 0o700); err != nil {
		return fmt.Errorf("mkdir cursor dir: %w", err)
	}
	if err := os.WriteFile(filepath.Join(cursorDir, "mcp.json"), data, 0o600); err != nil {
		return fmt.Errorf("write cursor config: %w", err)
	}
	return nil
}

func (s *scenarioCtx) aGlobalGoogleGeminiConfigExistsWithServer(serverName, command, rawArgs string) error {
	s.reg = buildIsolatedReg(s.root)
	args := strings.Fields(rawArgs)
	entry := map[string]any{"command": command, "args": args}
	doc := map[string]any{"mcpServers": map[string]any{serverName: entry}}
	data, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("marshal gemini config: %w", err)
	}
	globalPath := filepath.Join(s.root, "gemini-settings.json")
	if err := os.WriteFile(globalPath, data, 0o600); err != nil {
		return fmt.Errorf("write gemini config: %w", err)
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

// ── interactive init steps ────────────────────────────────────────────────────

func (s *scenarioCtx) iRunInteractiveInitAndSelectNoServers() error {
	return s.runInteractiveInit([]string{
		"none", // PromptMultiSelect: select nothing
		"no",   // Add a server?
	})
}

func (s *scenarioCtx) iRunInteractiveInitAndCollectImportCandidates() error {
	var err error
	s.importCandidates, err = wizard.CollectImportCandidates(s.root, s.reg)
	if err != nil {
		return fmt.Errorf("CollectImportCandidates: %w", err)
	}
	s.groupedCandidates = wizard.GroupImportCandidates(s.importCandidates)
	return nil
}

func (s *scenarioCtx) iRunInteractiveInitAndSelectServer(serverName string) error {
	candidates, err := wizard.CollectImportCandidates(s.root, s.reg)
	if err != nil {
		return fmt.Errorf("collect candidates: %w", err)
	}
	idx := -1
	for i, c := range candidates {
		if c.ServerName == serverName {
			idx = i
			break
		}
	}
	if idx < 0 {
		return fmt.Errorf("server %q not found in import candidates", serverName)
	}
	return s.runInteractiveInit([]string{
		fmt.Sprintf("%d", idx), // PromptMultiSelect: select this index
		"no",                   // Add a server?
	})
}

func (s *scenarioCtx) iRunInteractiveInitAndSelectAllDiscoveredServers() error {
	return s.runInteractiveInit([]string{
		"all", // PromptMultiSelect: select all
		"no",  // Add a server?
	})
}

func (s *scenarioCtx) iRunInteractiveInitSkipImportAndManuallyAddServer(name, command, args string) error {
	argList := strings.Fields(args)
	// No candidates exist (noProviderConfigFilesExist was the Given), so
	// PromptMultiSelect is never called — go straight to the manual loop.
	answers := []string{
		"yes",                      // Add a server?
		name,                       // Server name
		"stdio",                    // Transport
		command,                    // Command
		strings.Join(argList, " "), // Args
		"no",                       // Add a server?
	}
	return s.runInteractiveInit(answers)
}

func (s *scenarioCtx) iRunInteractiveInitImportServerAndManuallyAddServer(importName, manualName, command, args string) error {
	candidates, err := wizard.CollectImportCandidates(s.root, s.reg)
	if err != nil {
		return fmt.Errorf("collect candidates: %w", err)
	}
	idx := -1
	for i, c := range candidates {
		if c.ServerName == importName {
			idx = i
			break
		}
	}
	if idx < 0 {
		return fmt.Errorf("server %q not found in import candidates", importName)
	}
	argList := strings.Fields(args)
	answers := []string{
		fmt.Sprintf("%d", idx),     // PromptMultiSelect: select import
		"yes",                      // Add a server?
		manualName,                 // Server name
		"stdio",                    // Transport
		command,                    // Command
		strings.Join(argList, " "), // Args
		"no",                       // Add a server?
	}
	return s.runInteractiveInit(answers)
}

// ── interactive init assertions ────────────────────────────────────────────────

func (s *scenarioCtx) theIrisConfigContainsNServers(n int) error {
	count, err := s.irisConfigServerCount()
	if err != nil {
		return err
	}
	if count != n {
		return fmt.Errorf("expected %d servers, got %d", n, count)
	}
	return nil
}

func (s *scenarioCtx) theIrisConfigContainsServerWithCommand(name, command string) error {
	store2, err := config.NewStore(s.storePath)
	if err != nil {
		return fmt.Errorf("NewStore: %w", err)
	}
	cfg, err := store2.Load()
	if err != nil {
		return fmt.Errorf("load: %w", err)
	}
	srv, ok := cfg.Servers[name]
	if !ok {
		return fmt.Errorf("server %q not found in iris config", name)
	}
	if srv.Command != command {
		return fmt.Errorf("server %q: expected command %q, got %q", name, command, srv.Command)
	}
	return nil
}

func (s *scenarioCtx) theIrisConfigProvidersListContains(providerName string) error {
	store2, err := config.NewStore(s.storePath)
	if err != nil {
		return fmt.Errorf("NewStore: %w", err)
	}
	cfg, err := store2.Load()
	if err != nil {
		return fmt.Errorf("load: %w", err)
	}
	for _, p := range cfg.Providers {
		if p == providerName {
			return nil
		}
	}
	return fmt.Errorf("providers list %v does not contain %q", cfg.Providers, providerName)
}

func (s *scenarioCtx) theIrisConfigProvidersListIsSetTo(providerName string) error {
	s.cfg.Providers = []string{providerName}
	if err := s.store.Save(s.cfg); err != nil {
		return fmt.Errorf("save: %w", err)
	}
	return nil
}

func (s *scenarioCtx) theImportCandidatesIncludeEntry(serverName, providerName, scope string) error {
	for _, c := range s.importCandidates {
		if c.ServerName == serverName && c.ProviderName == providerName && string(c.Scope) == scope {
			return nil
		}
	}
	return fmt.Errorf("no candidate found for server=%q provider=%q scope=%q (got: %v)",
		serverName, providerName, scope, s.importCandidates)
}

func (s *scenarioCtx) theGroupedCandidatesContainExactlyNEntryForServer(n int, serverName string) error {
	count := 0
	for _, g := range s.groupedCandidates {
		if g.ServerName == serverName {
			count++
		}
	}
	if count != n {
		return fmt.Errorf("expected %d grouped entries for server %q, got %d", n, serverName, count)
	}
	return nil
}

func (s *scenarioCtx) theGroupedCandidateForServerListsProviders(serverName, p1, p2 string) error {
	for _, g := range s.groupedCandidates {
		if g.ServerName != serverName {
			continue
		}
		has := func(name string) bool {
			for _, p := range g.Providers {
				if p == name {
					return true
				}
			}
			return false
		}
		if !has(p1) || !has(p2) {
			return fmt.Errorf("grouped candidate %q has providers %v, expected both %q and %q",
				serverName, g.Providers, p1, p2)
		}
		return nil
	}
	return fmt.Errorf("no grouped candidate found for server %q", serverName)
}

// ── init file assertions ───────────────────────────────────────────────────────

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
