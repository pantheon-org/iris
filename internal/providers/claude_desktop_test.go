package providers_test

import (
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
	if cfg.Name != "anthropic-claude-desktop" {
		t.Fatalf("expected name %q, got %q", "anthropic-claude-desktop", cfg.Name)
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
