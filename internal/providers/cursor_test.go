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

func TestCursorProvider_Parse_withTestdata(t *testing.T) {
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
	if parsed["brave-search"].Env["BRAVE_API_KEY"] != "test-key" {
		t.Errorf("brave-search env = %v, want BRAVE_API_KEY=test-key", parsed["brave-search"].Env)
	}
}

func TestCursorProvider_Generate_stripsTypeField(t *testing.T) {
	content, err := os.ReadFile("testdata/cursor_input.json")
	if err != nil {
		t.Fatal(err)
	}
	p := providers.NewCursorProvider()
	servers := map[string]types.MCPServer{
		"filesystem": {Command: "npx", Args: []string{"-y", "@modelcontextprotocol/server-filesystem", "/tmp"}},
	}
	out, err := p.Generate(servers, string(content))
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	var doc struct {
		MCPServers map[string]struct {
			Type string `json:"type"`
		} `json:"mcpServers"`
	}
	if err := json.Unmarshal([]byte(out), &doc); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if doc.MCPServers["filesystem"].Type != "" {
		t.Errorf("expected type field to be stripped on generate, got %q", doc.MCPServers["filesystem"].Type)
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
