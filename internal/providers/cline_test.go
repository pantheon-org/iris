package providers_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pantheon-org/iris/internal/providers"
)

func TestClineProvider_Config(t *testing.T) {
	tmp := t.TempDir()
	p := providers.NewClineProviderWithPath(filepath.Join(tmp, "cline_mcp_settings.json"))
	cfg := p.Config()
	if cfg.Name != "cline" {
		t.Fatalf("expected name %q, got %q", "cline", cfg.Name)
	}
	if cfg.SupportsProjectConfig {
		t.Fatal("expected SupportsProjectConfig=false")
	}
}

func TestClineProvider_Parse_ExtractsServersFromFixture(t *testing.T) {
	content, err := os.ReadFile("fixtures/cline_input.json")
	if err != nil {
		t.Fatal(err)
	}
	tmp := t.TempDir()
	p := providers.NewClineProviderWithPath(filepath.Join(tmp, "cline_mcp_settings.json"))
	parsed, err := p.Parse(string(content))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if len(parsed) != 3 {
		t.Fatalf("expected 3 servers, got %d", len(parsed))
	}
	if parsed["filesystem"].Command != "npx" {
		t.Errorf("filesystem.command = %q, want %q", parsed["filesystem"].Command, "npx")
	}
	if parsed["fetch"].Command != "uvx" {
		t.Errorf("fetch.command = %q, want %q", parsed["fetch"].Command, "uvx")
	}
	if parsed["brave-search"].Command != "uvx" {
		t.Errorf("brave-search.command = %q, want %q", parsed["brave-search"].Command, "uvx")
	}
}

func TestClineProvider_Generate_FixtureMatch(t *testing.T) {
	content, err := os.ReadFile("fixtures/cline_input.json")
	if err != nil {
		t.Fatal(err)
	}
	tmp := t.TempDir()
	p := providers.NewClineProviderWithPath(filepath.Join(tmp, "cline_mcp_settings.json"))
	servers, err := p.Parse(string(content))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	got, err := p.Generate(servers, string(content))
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	expected, err := os.ReadFile("fixtures/cline_expected.json")
	if err != nil {
		t.Fatal(err)
	}
	requireJSONEqual(t, string(expected), got)
}

func TestClineProvider_Exists(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "cline_mcp_settings.json")
	p := providers.NewClineProviderWithPath(path)

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
