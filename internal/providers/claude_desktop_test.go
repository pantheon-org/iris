package providers_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/pantheon-org/iris/internal/providers"
	"github.com/pantheon-org/iris/internal/types"
)

func TestClaudeDesktopProvider_Config(t *testing.T) {
	tmp := t.TempDir()
	p := providers.NewClaudeDesktopProviderWithPath(filepath.Join(tmp, "claude_desktop_config.json"))
	cfg := p.Config()
	if cfg.Name != "claude-desktop" {
		t.Fatalf("expected name %q, got %q", "claude-desktop", cfg.Name)
	}
	if cfg.SupportsProjectConfig {
		t.Fatal("expected SupportsProjectConfig=false")
	}
}

func TestClaudeDesktopProvider_GenerateParse_roundtrip(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "claude_desktop_config.json")
	p := providers.NewClaudeDesktopProviderWithPath(path)

	servers := map[string]types.MCPServer{
		"memory": {Transport: types.TransportStdio, Command: "npx", Args: []string{"-y", "@modelcontextprotocol/server-memory"}},
	}
	out, err := p.Generate(servers, "")
	if err != nil {
		t.Fatal(err)
	}

	parsed, err := p.Parse(out)
	if err != nil {
		t.Fatal(err)
	}
	if parsed["memory"].Command != "npx" {
		t.Fatalf("expected command %q, got %q", "npx", parsed["memory"].Command)
	}
}

func TestClaudeDesktopProvider_Parse_withTestdata(t *testing.T) {
	content, err := os.ReadFile("testdata/claude_desktop_input.json")
	if err != nil {
		t.Fatal(err)
	}
	tmp := t.TempDir()
	p := providers.NewClaudeDesktopProviderWithPath(filepath.Join(tmp, "claude_desktop_config.json"))
	parsed, err := p.Parse(string(content))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if len(parsed) != 2 {
		t.Fatalf("expected 2 servers, got %d", len(parsed))
	}
	if parsed["filesystem"].Command != "npx" {
		t.Errorf("filesystem.command = %q, want %q", parsed["filesystem"].Command, "npx")
	}
	if parsed["brave-search"].Env["BRAVE_API_KEY"] != "test-key" {
		t.Errorf("brave-search env = %v, want BRAVE_API_KEY=test-key", parsed["brave-search"].Env)
	}
}

func TestClaudeDesktopProvider_Generate_withTestdata_preservesGlobalShortcut(t *testing.T) {
	content, err := os.ReadFile("testdata/claude_desktop_input.json")
	if err != nil {
		t.Fatal(err)
	}
	tmp := t.TempDir()
	p := providers.NewClaudeDesktopProviderWithPath(filepath.Join(tmp, "claude_desktop_config.json"))
	servers := map[string]types.MCPServer{
		"filesystem": {Command: "npx", Args: []string{"-y", "@modelcontextprotocol/server-filesystem", "/tmp"}},
	}
	out, err := p.Generate(servers, string(content))
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	var doc map[string]json.RawMessage
	if err := json.Unmarshal([]byte(out), &doc); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if _, ok := doc["globalShortcut"]; !ok {
		t.Error("expected 'globalShortcut' key to be preserved")
	}
}

func TestClaudeDesktopProvider_Exists(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "claude_desktop_config.json")
	p := providers.NewClaudeDesktopProviderWithPath(path)

	ok, err := p.Exists(tmp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Fatal("should not exist before file is created")
	}
	if err := os.WriteFile(path, []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	ok, err = p.Exists(tmp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("should exist after file is created")
	}
}
