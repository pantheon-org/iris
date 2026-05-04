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
	if cfg.Name != "moonshot-kimi" {
		t.Fatalf("expected name %q, got %q", "moonshot-kimi", cfg.Name)
	}
	if cfg.SupportsProjectConfig {
		t.Fatal("expected SupportsProjectConfig=false")
	}
}

func TestKimiProvider_GenerateParse_roundtrip(t *testing.T) {
	tmp := t.TempDir()
	p := providers.NewKimiProviderWithPath(filepath.Join(tmp, "mcp.json"))

	servers := map[string]types.MCPServer{
		"context7": {
			Transport: types.TransportSSE,
			URL:       "https://mcp.context7.com/mcp",
			Headers:   map[string]string{"CONTEXT7_API_KEY": "key"},
		},
		"chrome": {
			Transport: types.TransportStdio,
			Command:   "npx",
			Args:      []string{"chrome-devtools-mcp@latest"},
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
	if parsed["chrome"].Command != "npx" {
		t.Fatalf("expected command %q, got %q", "npx", parsed["chrome"].Command)
	}
	if parsed["context7"].Transport != types.TransportSSE {
		t.Fatalf("expected transport %q, got %q", types.TransportSSE, parsed["context7"].Transport)
	}
	if parsed["context7"].URL != "https://mcp.context7.com/mcp" {
		t.Fatalf("expected URL to roundtrip, got %q", parsed["context7"].URL)
	}
	if parsed["context7"].Headers["CONTEXT7_API_KEY"] != "key" {
		t.Fatalf("expected headers to roundtrip, got %v", parsed["context7"].Headers)
	}
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
