package providers_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pantheon-org/iris/internal/providers"
	"github.com/pantheon-org/iris/internal/types"
)

func TestMistralVibeProvider_Config(t *testing.T) {
	tmp := t.TempDir()
	p := providers.NewMistralVibeProviderWithPath(filepath.Join(tmp, "config.toml"))
	cfg := p.Config()
	if cfg.Name != "mistral-vibe" {
		t.Fatalf("expected name %q, got %q", "mistral-vibe", cfg.Name)
	}
	if cfg.SupportsProjectConfig {
		t.Fatal("expected SupportsProjectConfig=false")
	}
}

func TestMistralVibeProvider_GenerateParse_roundtrip(t *testing.T) {
	tmp := t.TempDir()
	p := providers.NewMistralVibeProviderWithPath(filepath.Join(tmp, "config.toml"))

	servers := map[string]types.MCPServer{
		"fetch": {
			Transport: types.TransportStdio,
			Command:   "uvx",
			Args:      []string{"mcp-server-fetch"},
			Env:       map[string]string{"DEBUG": "1"},
		},
	}
	out, err := p.Generate(servers, "")
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(out, "[[mcp_servers]]") {
		t.Fatalf("expected [[mcp_servers]] block, got:\n%s", out)
	}

	parsed, err := p.Parse(out)
	if err != nil {
		t.Fatal(err)
	}
	if parsed["fetch"].Command != "uvx" {
		t.Fatalf("expected command %q, got %q", "uvx", parsed["fetch"].Command)
	}
	if parsed["fetch"].Env["DEBUG"] != "1" {
		t.Fatalf("expected env var, got %v", parsed["fetch"].Env)
	}
}

func TestMistralVibeProvider_Generate_preservesExistingKeys(t *testing.T) {
	tmp := t.TempDir()
	p := providers.NewMistralVibeProviderWithPath(filepath.Join(tmp, "config.toml"))

	existing := `model = "codestral-latest"
api_key_env = "MISTRAL_API_KEY"
`
	servers := map[string]types.MCPServer{
		"fetch": {Transport: types.TransportStdio, Command: "uvx", Args: []string{"mcp-server-fetch"}},
	}
	out, err := p.Generate(servers, existing)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "model") {
		t.Fatalf("expected existing 'model' key to be preserved, got:\n%s", out)
	}
	if !strings.Contains(out, "[[mcp_servers]]") {
		t.Fatalf("expected [[mcp_servers]] block, got:\n%s", out)
	}
}

func TestMistralVibeProvider_Generate_httpTransport(t *testing.T) {
	tmp := t.TempDir()
	p := providers.NewMistralVibeProviderWithPath(filepath.Join(tmp, "config.toml"))

	servers := map[string]types.MCPServer{
		"context7": {
			Transport: types.TransportSSE,
			URL:       "https://mcp.context7.com/mcp",
		},
	}
	out, err := p.Generate(servers, "")
	if err != nil {
		t.Fatal(err)
	}

	parsed, err := p.Parse(out)
	if err != nil {
		t.Fatal(err)
	}
	if parsed["context7"].URL != "https://mcp.context7.com/mcp" {
		t.Fatalf("expected URL, got %q", parsed["context7"].URL)
	}
}

func TestMistralVibeProvider_Exists(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "config.toml")
	p := providers.NewMistralVibeProviderWithPath(path)

	if p.Exists(tmp) {
		t.Fatal("should not exist before file is created")
	}
	if err := os.WriteFile(path, []byte(""), 0o644); err != nil {
		t.Fatal(err)
	}
	if !p.Exists(tmp) {
		t.Fatal("should exist after file is created")
	}
}
