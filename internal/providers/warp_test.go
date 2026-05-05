package providers_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pantheon-org/iris/internal/providers"
)

func TestWarpProvider_Config(t *testing.T) {
	tmp := t.TempDir()
	p := providers.NewWarpProviderWithPath(filepath.Join(tmp, "mcp.json"))
	cfg := p.Config()
	if cfg.Name != "warp" {
		t.Fatalf("expected name %q, got %q", "warp", cfg.Name)
	}
	if cfg.SupportsProjectConfig {
		t.Fatal("expected SupportsProjectConfig=false")
	}
}

func TestWarpProvider_Parse_ExtractsServersFromFixture(t *testing.T) {
	content, err := os.ReadFile("testdata/warp_input.json")
	if err != nil {
		t.Fatal(err)
	}
	tmp := t.TempDir()
	p := providers.NewWarpProviderWithPath(filepath.Join(tmp, "mcp.json"))
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
	if parsed["fetch"].Command != "uvx" {
		t.Errorf("fetch.command = %q, want %q", parsed["fetch"].Command, "uvx")
	}
}

func TestWarpProvider_Generate_FixtureMatch(t *testing.T) {
	content, err := os.ReadFile("testdata/warp_input.json")
	if err != nil {
		t.Fatal(err)
	}
	tmp := t.TempDir()
	p := providers.NewWarpProviderWithPath(filepath.Join(tmp, "mcp.json"))
	servers, err := p.Parse(string(content))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	got, err := p.Generate(servers, string(content))
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	expected, err := os.ReadFile("testdata/warp_expected.json")
	if err != nil {
		t.Fatal(err)
	}
	if got != string(expected) {
		t.Errorf("output mismatch:\ngot:\n%s\nwant:\n%s", got, expected)
	}
}

func TestWarpProvider_Exists(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "mcp.json")
	p := providers.NewWarpProviderWithPath(path)

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
