package providers_test

import (
	"encoding/json"
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

func TestQwenProvider_GenerateParse_roundtrip(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "settings.json")
	p := providers.NewQwenProviderWithPath(path)

	servers := map[string]types.MCPServer{
		"fs": {Transport: types.TransportStdio, Command: "npx", Args: []string{"-y", "@modelcontextprotocol/server-filesystem", "/tmp"}},
	}
	out, err := p.Generate(servers, "")
	if err != nil {
		t.Fatal(err)
	}

	parsed, err := p.Parse(out)
	if err != nil {
		t.Fatal(err)
	}
	if parsed["fs"].Command != "npx" {
		t.Fatalf("expected command %q, got %q", "npx", parsed["fs"].Command)
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

func TestQwenProvider_Generate_HTTPServer_writesHttpUrl(t *testing.T) {
	p := providers.NewQwenProviderWithPath(t.TempDir() + "/settings.json")
	servers := map[string]types.MCPServer{
		"remote": {Transport: types.TransportHTTP, URL: "https://api.example.com/mcp"},
	}
	out, err := p.Generate(servers, "")
	if err != nil {
		t.Fatal(err)
	}
	var doc struct {
		MCPServers map[string]struct {
			URL     string `json:"url"`
			HTTPUrl string `json:"httpUrl"`
		} `json:"mcpServers"`
	}
	if err := json.Unmarshal([]byte(out), &doc); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	entry := doc.MCPServers["remote"]
	if entry.HTTPUrl != "https://api.example.com/mcp" {
		t.Fatalf("expected httpUrl=%q, got %q", "https://api.example.com/mcp", entry.HTTPUrl)
	}
	if entry.URL != "" {
		t.Fatalf("expected url to be empty, got %q", entry.URL)
	}
}

func TestQwenProvider_Parse_httpUrl_setsTransportHTTP(t *testing.T) {
	p := providers.NewQwenProviderWithPath(t.TempDir() + "/settings.json")
	input := `{"mcpServers":{"remote":{"httpUrl":"https://api.example.com/mcp"}}}`
	parsed, err := p.Parse(input)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	srv := parsed["remote"]
	if srv.Transport != types.TransportHTTP {
		t.Fatalf("expected TransportHTTP, got %q", srv.Transport)
	}
	if srv.URL != "https://api.example.com/mcp" {
		t.Fatalf("expected URL=%q, got %q", "https://api.example.com/mcp", srv.URL)
	}
}

func TestQwenProvider_Parse_withTestdata(t *testing.T) {
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

func TestQwenProvider_Generate_noTypeField(t *testing.T) {
	p := providers.NewQwenProviderWithPath(t.TempDir() + "/settings.json")
	servers := map[string]types.MCPServer{
		"local": {Transport: types.TransportStdio, Command: "npx"},
	}
	out, err := p.Generate(servers, "")
	if err != nil {
		t.Fatal(err)
	}
	var doc struct {
		MCPServers map[string]struct {
			Type string `json:"type"`
		} `json:"mcpServers"`
	}
	if err := json.Unmarshal([]byte(out), &doc); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if doc.MCPServers["local"].Type != "" {
		t.Fatalf("expected no type field, got %q", doc.MCPServers["local"].Type)
	}
}
