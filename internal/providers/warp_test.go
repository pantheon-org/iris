package providers_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pantheon-org/iris/internal/providers"
	"github.com/pantheon-org/iris/internal/types"
)

func TestWarpProvider_Config(t *testing.T) {
	tmp := t.TempDir()
	p := providers.NewWarpProviderWithPath(filepath.Join(tmp, "mcp.json"))
	cfg := p.Config()
	if cfg.Name != "warp-terminal" {
		t.Fatalf("expected name %q, got %q", "warp-terminal", cfg.Name)
	}
	if cfg.SupportsProjectConfig {
		t.Fatal("expected SupportsProjectConfig=false")
	}
}

func TestWarpProvider_GenerateParse_roundtrip(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "mcp.json")
	p := providers.NewWarpProviderWithPath(path)

	servers := map[string]types.MCPServer{
		"brave": {Transport: types.TransportStdio, Command: "npx", Args: []string{"-y", "@modelcontextprotocol/server-brave-search"}, Env: map[string]string{"BRAVE_API_KEY": "key"}},
	}
	out, err := p.Generate(servers, "")
	if err != nil {
		t.Fatal(err)
	}

	parsed, err := p.Parse(out)
	if err != nil {
		t.Fatal(err)
	}
	if parsed["brave"].Env["BRAVE_API_KEY"] != "key" {
		t.Fatalf("expected env var, got %v", parsed["brave"].Env)
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
