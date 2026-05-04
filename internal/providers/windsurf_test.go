package providers_test

import (
	"encoding/json"
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

func TestWindsurfProvider_Generate_remoteServer_writesServerUrl(t *testing.T) {
	p := providers.NewWindsurfProviderWithPath(t.TempDir() + "/mcp_config.json")
	servers := map[string]types.MCPServer{
		"remote": {Transport: types.TransportSSE, URL: "https://example.com/mcp"},
	}
	out, err := p.Generate(servers, "")
	if err != nil {
		t.Fatal(err)
	}

	var doc struct {
		MCPServers map[string]struct {
			ServerURL string `json:"serverUrl"`
			URL       string `json:"url"`
		} `json:"mcpServers"`
	}
	if err := json.Unmarshal([]byte(out), &doc); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	entry := doc.MCPServers["remote"]
	if entry.ServerURL != "https://example.com/mcp" {
		t.Fatalf("expected serverUrl=%q, got %q", "https://example.com/mcp", entry.ServerURL)
	}
	if entry.URL != "" {
		t.Fatalf("expected url to be empty, got %q", entry.URL)
	}
}

func TestWindsurfProvider_Parse_acceptsServerUrlOrUrl(t *testing.T) {
	p := providers.NewWindsurfProviderWithPath(t.TempDir() + "/mcp_config.json")

	// Test serverUrl (primary field)
	withServerURL := `{"mcpServers":{"remote":{"serverUrl":"https://example.com/mcp"}}}`
	parsed, err := p.Parse(withServerURL)
	if err != nil {
		t.Fatalf("Parse(serverUrl): %v", err)
	}
	if parsed["remote"].URL != "https://example.com/mcp" {
		t.Fatalf("serverUrl: expected URL=%q, got %q", "https://example.com/mcp", parsed["remote"].URL)
	}
	if parsed["remote"].Transport != types.TransportSSE {
		t.Fatalf("serverUrl: expected TransportSSE, got %q", parsed["remote"].Transport)
	}

	// Test fallback url field
	withURL := `{"mcpServers":{"remote":{"url":"https://fallback.com/mcp"}}}`
	parsed, err = p.Parse(withURL)
	if err != nil {
		t.Fatalf("Parse(url): %v", err)
	}
	if parsed["remote"].URL != "https://fallback.com/mcp" {
		t.Fatalf("url: expected URL=%q, got %q", "https://fallback.com/mcp", parsed["remote"].URL)
	}
}

func TestWindsurfProvider_Exists(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "mcp_config.json")
	p := providers.NewWindsurfProviderWithPath(path)

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
