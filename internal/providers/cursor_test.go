package providers_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pantheon-org/iris/internal/providers"
	"github.com/pantheon-org/iris/internal/types"
)

func TestCursorProvider_Config(t *testing.T) {
	p := providers.NewCursorProvider()
	cfg := p.Config()
	if cfg.Name != "cursor" {
		t.Fatalf("expected name %q, got %q", "cursor", cfg.Name)
	}
	if !cfg.SupportsProjectConfig {
		t.Fatal("expected SupportsProjectConfig=true")
	}
}

func TestCursorProvider_ConfigFilePath(t *testing.T) {
	p := providers.NewCursorProvider()
	got := p.ConfigFilePath("/project")
	want := filepath.Join("/project", ".cursor", "mcp.json")
	if got != want {
		t.Fatalf("want %q, got %q", want, got)
	}
}

func TestCursorProvider_Parse_ExtractsServersFromFixture(t *testing.T) {
	content, err := os.ReadFile("testdata/cursor_input.json")
	if err != nil {
		t.Fatal(err)
	}
	p := providers.NewCursorProvider()
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
	if parsed["filesystem"].Transport != types.TransportStdio {
		t.Errorf("filesystem.transport = %q, want stdio", parsed["filesystem"].Transport)
	}
	if parsed["brave-search"].Env["BRAVE_API_KEY"] != "test-key" {
		t.Errorf("brave-search env = %v, want BRAVE_API_KEY=test-key", parsed["brave-search"].Env)
	}
}

func TestCursorProvider_Generate_FixtureMatch(t *testing.T) {
	content, err := os.ReadFile("testdata/cursor_input.json")
	if err != nil {
		t.Fatal(err)
	}
	p := providers.NewCursorProvider()
	servers, err := p.Parse(string(content))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	got, err := p.Generate(servers, string(content))
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	expected, err := os.ReadFile("testdata/cursor_expected.json")
	if err != nil {
		t.Fatal(err)
	}
	if got != string(expected) {
		t.Errorf("output mismatch:\ngot:\n%s\nwant:\n%s", got, expected)
	}
}

func TestCursorProvider_Exists(t *testing.T) {
	tmp := t.TempDir()
	p := providers.NewCursorProvider()

	ok, err := p.Exists(tmp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Fatal("should not exist before file is created")
	}

	dir := filepath.Join(tmp, ".cursor")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "mcp.json"), []byte("{}"), 0o644); err != nil {
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
