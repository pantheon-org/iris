package providers_test

import (
	"encoding/json"
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

func TestCursorProvider_GenerateParse_roundtrip(t *testing.T) {
	p := providers.NewCursorProvider()
	servers := map[string]types.MCPServer{
		"fs": {Transport: types.TransportStdio, Command: "node", Args: []string{"server.js"}},
	}
	out, err := p.Generate(servers, "")
	if err != nil {
		t.Fatal(err)
	}

	var doc struct {
		MCPServers map[string]json.RawMessage `json:"mcpServers"`
	}
	if err := json.Unmarshal([]byte(out), &doc); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out)
	}
	if _, ok := doc.MCPServers["fs"]; !ok {
		t.Fatalf("expected mcpServers.fs in output")
	}

	parsed, err := p.Parse(out)
	if err != nil {
		t.Fatal(err)
	}
	if parsed["fs"].Command != "node" {
		t.Fatalf("expected command %q, got %q", "node", parsed["fs"].Command)
	}
}

func TestCursorProvider_Exists(t *testing.T) {
	tmp := t.TempDir()
	p := providers.NewCursorProvider()

	if p.Exists(tmp) {
		t.Fatal("should not exist before file is created")
	}

	dir := filepath.Join(tmp, ".cursor")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "mcp.json"), []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	if !p.Exists(tmp) {
		t.Fatal("should exist after file is created")
	}
}
