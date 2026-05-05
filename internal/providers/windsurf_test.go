package providers_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pantheon-org/iris/internal/providers"
	"github.com/pantheon-org/iris/internal/types"
)

func TestWindsurfProvider_Config(t *testing.T) {
	tmp := t.TempDir()
	p := providers.NewWindsurfProviderWithPath(filepath.Join(tmp, "mcp_config.json"))
	cfg := p.Config()
	if cfg.Name != "windsurf" {
		t.Fatalf("expected name %q, got %q", "windsurf", cfg.Name)
	}
	if cfg.SupportsProjectConfig {
		t.Fatal("expected SupportsProjectConfig=false")
	}
}

func TestWindsurfProvider_Parse_ExtractsServersFromFixture(t *testing.T) {
	content, err := os.ReadFile("testdata/windsurf_input.json")
	if err != nil {
		t.Fatal(err)
	}
	tmp := t.TempDir()
	p := providers.NewWindsurfProviderWithPath(filepath.Join(tmp, "mcp_config.json"))
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
	if parsed["context7"].URL != "https://mcp.context7.com/mcp" {
		t.Errorf("context7.url = %q, want URL", parsed["context7"].URL)
	}
	if parsed["context7"].Transport != types.TransportSSE {
		t.Errorf("context7.transport = %q, want sse", parsed["context7"].Transport)
	}
}

func TestWindsurfProvider_Generate_FixtureMatch(t *testing.T) {
	content, err := os.ReadFile("testdata/windsurf_input.json")
	if err != nil {
		t.Fatal(err)
	}
	tmp := t.TempDir()
	p := providers.NewWindsurfProviderWithPath(filepath.Join(tmp, "mcp_config.json"))
	servers, err := p.Parse(string(content))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	got, err := p.Generate(servers, string(content))
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	expected, err := os.ReadFile("testdata/windsurf_expected.json")
	if err != nil {
		t.Fatal(err)
	}
	requireJSONEqual(t, string(expected), got)
}

func TestWindsurfProvider_Exists(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "mcp_config.json")
	p := providers.NewWindsurfProviderWithPath(path)

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
