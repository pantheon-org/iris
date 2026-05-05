package providers_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pantheon-org/iris/internal/providers"
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

func TestClaudeDesktopProvider_Parse_ExtractsServersFromFixture(t *testing.T) {
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

func TestClaudeDesktopProvider_Generate_FixtureMatch(t *testing.T) {
	content, err := os.ReadFile("testdata/claude_desktop_input.json")
	if err != nil {
		t.Fatal(err)
	}
	tmp := t.TempDir()
	p := providers.NewClaudeDesktopProviderWithPath(filepath.Join(tmp, "claude_desktop_config.json"))
	servers, err := p.Parse(string(content))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	got, err := p.Generate(servers, string(content))
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	expected, err := os.ReadFile("testdata/claude_desktop_expected.json")
	if err != nil {
		t.Fatal(err)
	}
	requireJSONEqual(t, string(expected), got)
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
