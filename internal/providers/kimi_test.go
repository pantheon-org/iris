package providers_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pantheon-org/iris/internal/providers"
	"github.com/pantheon-org/iris/internal/types"
)

func TestKimiProvider_Config(t *testing.T) {
	tmp := t.TempDir()
	p := providers.NewKimiProviderWithPath(filepath.Join(tmp, "mcp.json"))
	cfg := p.Config()
	if cfg.Name != "kimi" {
		t.Fatalf("expected name %q, got %q", "kimi", cfg.Name)
	}
	if cfg.SupportsProjectConfig {
		t.Fatal("expected SupportsProjectConfig=false")
	}
}

func TestKimiProvider_Parse_ExtractsServersFromFixture(t *testing.T) {
	content, err := os.ReadFile("testdata/kimi_input.json")
	if err != nil {
		t.Fatal(err)
	}
	tmp := t.TempDir()
	p := providers.NewKimiProviderWithPath(filepath.Join(tmp, "mcp.json"))
	parsed, err := p.Parse(string(content))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if len(parsed) != 2 {
		t.Fatalf("expected 2 servers, got %d", len(parsed))
	}
	if parsed["chrome-devtools"].Command != "npx" {
		t.Errorf("chrome-devtools.command = %q, want %q", parsed["chrome-devtools"].Command, "npx")
	}
	if parsed["context7"].URL != "https://mcp.context7.com/mcp" {
		t.Errorf("context7.url = %q, want URL", parsed["context7"].URL)
	}
	if parsed["context7"].Transport != types.TransportSSE {
		t.Errorf("context7.transport = %q, want sse", parsed["context7"].Transport)
	}
}

func TestKimiProvider_Generate_FixtureMatch(t *testing.T) {
	content, err := os.ReadFile("testdata/kimi_input.json")
	if err != nil {
		t.Fatal(err)
	}
	tmp := t.TempDir()
	p := providers.NewKimiProviderWithPath(filepath.Join(tmp, "mcp.json"))
	servers, err := p.Parse(string(content))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	got, err := p.Generate(servers, string(content))
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	expected, err := os.ReadFile("testdata/kimi_expected.json")
	if err != nil {
		t.Fatal(err)
	}
	requireJSONEqual(t, string(expected), got)
}

func TestKimiProvider_Exists(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "mcp.json")
	p := providers.NewKimiProviderWithPath(path)

	ok, err := p.Exists(tmp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Fatal("should not exist before file is created")
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
