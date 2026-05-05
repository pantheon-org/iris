package providers_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/pantheon-org/iris/internal/ierrors"
	"github.com/pantheon-org/iris/internal/providers"
)

func TestZedProvider_Config(t *testing.T) {
	tmp := t.TempDir()
	p := providers.NewZedProviderWithPath(filepath.Join(tmp, "settings.json"))
	cfg := p.Config()
	if cfg.Name != "zed" {
		t.Fatalf("expected name %q, got %q", "zed", cfg.Name)
	}
	if cfg.SupportsProjectConfig {
		t.Fatal("expected SupportsProjectConfig=false")
	}
}

func TestZedProvider_Parse_ExtractsServersFromFixture(t *testing.T) {
	content, err := os.ReadFile("testdata/zed_input.json")
	if err != nil {
		t.Fatal(err)
	}
	tmp := t.TempDir()
	p := providers.NewZedProviderWithPath(filepath.Join(tmp, "settings.json"))
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
}

func TestZedProvider_Generate_FixtureMatch(t *testing.T) {
	content, err := os.ReadFile("testdata/zed_input.json")
	if err != nil {
		t.Fatal(err)
	}
	tmp := t.TempDir()
	p := providers.NewZedProviderWithPath(filepath.Join(tmp, "settings.json"))
	servers, err := p.Parse(string(content))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	got, err := p.Generate(servers, string(content))
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	expected, err := os.ReadFile("testdata/zed_expected.json")
	if err != nil {
		t.Fatal(err)
	}
	if got != string(expected) {
		t.Errorf("output mismatch:\ngot:\n%s\nwant:\n%s", got, expected)
	}
}

func TestZedProvider_Parse_malformedInput_returnsError(t *testing.T) {
	tmp := t.TempDir()
	p := providers.NewZedProviderWithPath(filepath.Join(tmp, "settings.json"))
	_, err := p.Parse("not json at all")
	if err == nil {
		t.Fatal("Parse: expected error for malformed input, got nil")
	}
	if !errors.Is(err, ierrors.ErrMalformedConfig) {
		t.Errorf("Parse: error does not wrap ErrMalformedConfig; got: %v", err)
	}
}

func TestZedProvider_Exists(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "settings.json")
	p := providers.NewZedProviderWithPath(path)

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
