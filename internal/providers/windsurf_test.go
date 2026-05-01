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

func TestWindsurfProvider_GenerateParse_roundtrip(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "mcp_config.json")
	p := providers.NewWindsurfProviderWithPath(path)

	servers := map[string]types.MCPServer{
		"fetch": {Transport: types.TransportStdio, Command: "uvx", Args: []string{"mcp-server-fetch"}},
	}
	out, err := p.Generate(servers, "")
	if err != nil {
		t.Fatal(err)
	}

	parsed, err := p.Parse(out)
	if err != nil {
		t.Fatal(err)
	}
	if parsed["fetch"].Command != "uvx" {
		t.Fatalf("expected command %q, got %q", "uvx", parsed["fetch"].Command)
	}
	if len(parsed["fetch"].Args) != 1 || parsed["fetch"].Args[0] != "mcp-server-fetch" {
		t.Fatalf("unexpected args: %v", parsed["fetch"].Args)
	}
}

func TestWindsurfProvider_Exists(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "mcp_config.json")
	p := providers.NewWindsurfProviderWithPath(path)

	if p.Exists(tmp) {
		t.Fatal("should not exist before file is created")
	}
	if err := os.WriteFile(path, []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	if !p.Exists(tmp) {
		t.Fatal("should exist after file is created")
	}
}
