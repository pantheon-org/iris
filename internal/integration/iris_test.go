package integration_test

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/BurntSushi/toml"

	"github.com/pantheon-org/iris/internal/cli"
	"github.com/pantheon-org/iris/internal/config"
	"github.com/pantheon-org/iris/internal/providers"
	"github.com/pantheon-org/iris/internal/registry"
	irisync "github.com/pantheon-org/iris/internal/sync"
	"github.com/pantheon-org/iris/internal/types"
)

func buildRegistry(t *testing.T, root string) *registry.Registry {
	t.Helper()
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

func TestIris_fullPipeline_syncAllProviders(t *testing.T) {
	root := t.TempDir()

	storePath := filepath.Join(root, ".iris.json")
	store, err := config.NewStore(storePath)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}

	cfg := &types.IrisConfig{
		Version: 1,
		Servers: make(map[string]types.MCPServer),
	}

	if err := cli.RunAdd(cfg, store, "filesystem", types.MCPServer{
		Transport: types.TransportStdio,
		Command:   "npx",
		Args:      []string{"-y", "@modelcontextprotocol/server-filesystem", "/tmp"},
	}); err != nil {
		t.Fatalf("RunAdd filesystem: %v", err)
	}

	if err := cli.RunAdd(cfg, store, "fetch", types.MCPServer{
		Transport: types.TransportStdio,
		Command:   "uvx",
		Args:      []string{"mcp-server-fetch"},
	}); err != nil {
		t.Fatalf("RunAdd fetch: %v", err)
	}

	reg := buildRegistry(t, root)

	if err := cli.RunSync(root, cfg, reg, io.Discard, false); err != nil {
		t.Fatalf("RunSync (first): %v", err)
	}

	claudePath := filepath.Join(root, ".mcp.json")
	googleGeminiPath := filepath.Join(root, "gemini-settings.json")
	opencodePath := filepath.Join(root, "opencode.json")
	codexPath := filepath.Join(root, "codex-config.toml")

	for _, path := range []string{claudePath, googleGeminiPath, opencodePath, codexPath} {
		if _, err := os.Stat(path); err != nil {
			t.Errorf("expected file %s to exist: %v", path, err)
		}
	}

	assertJSONMCPServers := func(t *testing.T, path string) {
		t.Helper()
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read %s: %v", path, err)
		}
		var doc struct {
			MCPServers map[string]json.RawMessage `json:"mcpServers"`
		}
		if err := json.Unmarshal(data, &doc); err != nil {
			t.Fatalf("parse %s: %v", path, err)
		}
		for _, name := range []string{"filesystem", "fetch"} {
			if _, ok := doc.MCPServers[name]; !ok {
				t.Errorf("%s: missing server %q", path, name)
			}
		}
	}

	assertJSONMCPServers(t, claudePath)
	assertJSONMCPServers(t, googleGeminiPath)

	opencodeData, err := os.ReadFile(opencodePath)
	if err != nil {
		t.Fatalf("read opencode.json: %v", err)
	}
	var opencodeDoc struct {
		MCP map[string]json.RawMessage `json:"mcp"`
	}
	if err := json.Unmarshal(opencodeData, &opencodeDoc); err != nil {
		t.Fatalf("parse opencode.json: %v", err)
	}
	for _, name := range []string{"filesystem", "fetch"} {
		if _, ok := opencodeDoc.MCP[name]; !ok {
			t.Errorf("opencode.json: missing server %q", name)
		}
	}

	codexData, err := os.ReadFile(codexPath)
	if err != nil {
		t.Fatalf("read codex config: %v", err)
	}
	var codexDoc struct {
		MCPServers map[string]map[string]any `toml:"mcp_servers"`
	}
	if _, err := toml.Decode(string(codexData), &codexDoc); err != nil {
		t.Fatalf("parse codex config: %v", err)
	}
	for _, name := range []string{"filesystem", "fetch"} {
		if _, ok := codexDoc.MCPServers[name]; !ok {
			t.Errorf("codex config: missing server %q", name)
		}
	}

	reg2 := buildRegistry(t, root)
	if err := cli.RunSync(root, cfg, reg2, io.Discard, false); err != nil {
		t.Fatalf("RunSync (second): %v", err)
	}

	results := irisync.SyncAllProviders(root, reg2, cfg.Servers)
	for _, r := range results {
		if r.Err != nil {
			t.Errorf("provider %s: unexpected error: %v", r.ProviderName, r.Err)
		}
		if r.Status != irisync.SyncStatusUnchanged {
			t.Errorf("provider %s: expected %s, got %s", r.ProviderName, irisync.SyncStatusUnchanged, r.Status)
		}
	}
}

func TestIris_addRemove_persistsCorrectly(t *testing.T) {
	root := t.TempDir()

	storePath := filepath.Join(root, ".iris.json")
	store, err := config.NewStore(storePath)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}

	cfg := &types.IrisConfig{
		Version: 1,
		Servers: make(map[string]types.MCPServer),
	}

	servers := map[string]types.MCPServer{
		"alpha": {Transport: types.TransportStdio, Command: "cmd-alpha"},
		"beta":  {Transport: types.TransportStdio, Command: "cmd-beta"},
		"gamma": {Transport: types.TransportStdio, Command: "cmd-gamma"},
	}
	for name, srv := range servers {
		if err := cli.RunAdd(cfg, store, name, srv); err != nil {
			t.Fatalf("RunAdd %s: %v", name, err)
		}
	}

	if err := cli.RunRemove(cfg, store, "gamma"); err != nil {
		t.Fatalf("RunRemove gamma: %v", err)
	}

	store2, err := config.NewStore(storePath)
	if err != nil {
		t.Fatalf("NewStore (reload): %v", err)
	}
	loaded, err := store2.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if len(loaded.Servers) != 2 {
		t.Fatalf("expected 2 servers, got %d", len(loaded.Servers))
	}
	if _, ok := loaded.Servers["gamma"]; ok {
		t.Error("gamma should have been removed")
	}
	for _, name := range []string{"alpha", "beta"} {
		srv, ok := loaded.Servers[name]
		if !ok {
			t.Errorf("missing server %q", name)
			continue
		}
		if srv.Command != servers[name].Command {
			t.Errorf("server %q: expected command %q, got %q", name, servers[name].Command, srv.Command)
		}
	}
}
