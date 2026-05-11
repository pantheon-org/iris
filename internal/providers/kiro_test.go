package providers_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pantheon-org/iris/internal/providers"
)

func TestKiroProvider_Config(t *testing.T) {
	tmp := t.TempDir()
	p := providers.NewKiroProviderWithPath(filepath.Join(tmp, "mcp.json"))
	cfg := p.Config()
	if cfg.Name != "kiro" {
		t.Fatalf("expected name %q, got %q", "kiro", cfg.Name)
	}
	if !cfg.SupportsProjectConfig {
		t.Fatal("expected SupportsProjectConfig=true")
	}
	if cfg.LocalConfigPath == nil || *cfg.LocalConfigPath != ".kiro/settings/mcp.json" {
		t.Fatalf("expected local config .kiro/settings/mcp.json, got %v", cfg.LocalConfigPath)
	}
}

func TestKiroProvider_Parse_ExtractsServersFromFixture(t *testing.T) {
	content, err := os.ReadFile("fixtures/kiro_input.json")
	if err != nil {
		t.Fatal(err)
	}
	tmp := t.TempDir()
	p := providers.NewKiroProviderWithPath(filepath.Join(tmp, "mcp.json"))
	parsed, err := p.Parse(string(content))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if len(parsed) != 2 {
		t.Fatalf("expected 2 servers, got %d", len(parsed))
	}
	if parsed["web-search"].Command != "npx" {
		t.Errorf("web-search.command = %q, want %q", parsed["web-search"].Command, "npx")
	}
	if parsed["context7"].URL != "https://mcp.context7.com/mcp" {
		t.Errorf("context7.url = %q, want URL", parsed["context7"].URL)
	}
}

func TestKiroProvider_Generate_FixtureMatch(t *testing.T) {
	content, err := os.ReadFile("fixtures/kiro_input.json")
	if err != nil {
		t.Fatal(err)
	}
	tmp := t.TempDir()
	p := providers.NewKiroProviderWithPath(filepath.Join(tmp, "mcp.json"))
	servers, err := p.Parse(string(content))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	got, err := p.Generate(servers, string(content))
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	expected, err := os.ReadFile("fixtures/kiro_expected.json")
	if err != nil {
		t.Fatal(err)
	}
	requireJSONEqual(t, string(expected), got)
}

func TestKiroProvider_Exists(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, ".kiro", "settings", "mcp.json")
	p := providers.NewKiroProviderWithPath(path)

	ok, err := p.Exists(tmp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Fatal("should not exist before file is created")
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(`{"mcpServers":{}}`), 0o644); err != nil {
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
