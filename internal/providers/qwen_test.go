package providers_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pantheon-org/iris/internal/providers"
	"github.com/pantheon-org/iris/internal/types"
)

func TestQwenProvider_Config(t *testing.T) {
	tmp := t.TempDir()
	p := providers.NewQwenProviderWithPath(filepath.Join(tmp, "settings.json"))
	cfg := p.Config()
	if cfg.Name != "qwen" {
		t.Fatalf("expected name %q, got %q", "qwen", cfg.Name)
	}
	if !cfg.SupportsProjectConfig {
		t.Fatal("expected SupportsProjectConfig=true")
	}
}

func TestQwenProvider_ConfigFilePath_WithProjectRoot_ReturnsProjectPath(t *testing.T) {
	p := providers.NewQwenProvider()
	got := p.ConfigFilePath("/any/project")
	want := filepath.Join("/any/project", ".qwen", "settings.json")
	if got != want {
		t.Errorf("ConfigFilePath = %q, want %q", got, want)
	}
}

func TestQwenProvider_ConfigFilePath_WithEmptyRoot_ReturnsGlobalPath(t *testing.T) {
	p := providers.NewQwenProvider()
	got := p.ConfigFilePath("")
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(home, ".qwen", "settings.json")
	if got != want {
		t.Errorf("ConfigFilePath = %q, want %q", got, want)
	}
}

func TestQwenProvider_Parse_ExtractsServersFromFixture(t *testing.T) {
	content, err := os.ReadFile("testdata/qwen_input.json")
	if err != nil {
		t.Fatal(err)
	}
	tmp := t.TempDir()
	p := providers.NewQwenProviderWithPath(filepath.Join(tmp, "settings.json"))
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
	if parsed["context7"].Transport != types.TransportHTTP {
		t.Errorf("context7.transport = %q, want %q", parsed["context7"].Transport, types.TransportHTTP)
	}
}

func TestQwenProvider_Generate_FixtureMatch(t *testing.T) {
	content, err := os.ReadFile("testdata/qwen_input.json")
	if err != nil {
		t.Fatal(err)
	}
	tmp := t.TempDir()
	p := providers.NewQwenProviderWithPath(filepath.Join(tmp, "settings.json"))
	servers, err := p.Parse(string(content))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	got, err := p.Generate(servers, string(content))
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	expected, err := os.ReadFile("testdata/qwen_expected.json")
	if err != nil {
		t.Fatal(err)
	}
	if got != string(expected) {
		t.Errorf("output mismatch:\ngot:\n%s\nwant:\n%s", got, expected)
	}
}

func TestQwenProvider_Exists(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "settings.json")
	p := providers.NewQwenProviderWithPath(path)

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
