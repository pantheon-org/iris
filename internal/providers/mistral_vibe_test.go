package providers_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/pantheon-org/iris/internal/ierrors"
	"github.com/pantheon-org/iris/internal/providers"
)

func TestMistralVibeProvider_Config(t *testing.T) {
	tmp := t.TempDir()
	p := providers.NewMistralVibeProviderWithPath(filepath.Join(tmp, "config.toml"))
	cfg := p.Config()
	if cfg.Name != "mistral-vibe" {
		t.Fatalf("expected name %q, got %q", "mistral-vibe", cfg.Name)
	}
	if !cfg.SupportsProjectConfig {
		t.Fatal("expected SupportsProjectConfig=true")
	}
}

func TestMistralVibeProvider_ConfigFilePath_WithProjectRoot_ReturnsProjectPath(t *testing.T) {
	p := providers.NewMistralVibeProvider()
	got := p.ConfigFilePath("/any/project")
	want := filepath.Join("/any/project", ".vibe", "config.toml")
	if got != want {
		t.Errorf("ConfigFilePath = %q, want %q", got, want)
	}
}

func TestMistralVibeProvider_ConfigFilePath_WithEmptyRoot_ReturnsGlobalPath(t *testing.T) {
	p := providers.NewMistralVibeProvider()
	got := p.ConfigFilePath("")
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(home, ".vibe", "config.toml")
	if got != want {
		t.Errorf("ConfigFilePath = %q, want %q", got, want)
	}
}

func TestMistralVibeProvider_Parse_ExtractsServersFromFixture(t *testing.T) {
	content, err := os.ReadFile("testdata/mistral_vibe_input.toml")
	if err != nil {
		t.Fatal(err)
	}
	tmp := t.TempDir()
	p := providers.NewMistralVibeProviderWithPath(filepath.Join(tmp, "config.toml"))
	parsed, err := p.Parse(string(content))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if len(parsed) != 2 {
		t.Fatalf("expected 2 servers, got %d", len(parsed))
	}
	if parsed["fetch"].Command != "uvx" {
		t.Errorf("fetch.command = %q, want %q", parsed["fetch"].Command, "uvx")
	}
	if parsed["context7"].URL != "https://mcp.context7.com/mcp" {
		t.Errorf("context7.url = %q, want URL", parsed["context7"].URL)
	}
}

func TestMistralVibeProvider_Generate_FixtureMatch(t *testing.T) {
	content, err := os.ReadFile("testdata/mistral_vibe_input.toml")
	if err != nil {
		t.Fatal(err)
	}
	tmp := t.TempDir()
	p := providers.NewMistralVibeProviderWithPath(filepath.Join(tmp, "config.toml"))
	servers, err := p.Parse(string(content))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	got, err := p.Generate(servers, string(content))
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	expected, err := os.ReadFile("testdata/mistral_vibe_expected.toml")
	if err != nil {
		t.Fatal(err)
	}
	if got != string(expected) {
		t.Errorf("output mismatch:\ngot:\n%s\nwant:\n%s", got, expected)
	}
}

func TestMistralVibeProvider_Parse_malformedInput_returnsError(t *testing.T) {
	tmp := t.TempDir()
	p := providers.NewMistralVibeProviderWithPath(filepath.Join(tmp, "config.toml"))
	_, err := p.Parse("not toml at all \x00\xff")
	if err == nil {
		t.Fatal("Parse: expected error for malformed input, got nil")
	}
	if !errors.Is(err, ierrors.ErrMalformedConfig) {
		t.Errorf("Parse: error does not wrap ErrMalformedConfig; got: %v", err)
	}
}

func TestMistralVibeProvider_Exists(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "config.toml")
	p := providers.NewMistralVibeProviderWithPath(path)

	ok, err := p.Exists(tmp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Fatal("should not exist before file is created")
	}
	if err := os.WriteFile(path, []byte(""), 0o644); err != nil {
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
