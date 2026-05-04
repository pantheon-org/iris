package providers_test

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pantheon-org/iris/internal/ierrors"
	"github.com/pantheon-org/iris/internal/providers"
	"github.com/pantheon-org/iris/internal/types"
)

func TestGoogleGeminiProvider_Config_ReturnsCorrectProviderConfig(t *testing.T) {
	p := providers.NewGoogleGeminiProvider()
	cfg := p.Config()

	if cfg.Name != "gemini" {
		t.Errorf("Name = %q, want %q", cfg.Name, "gemini")
	}
	if cfg.DisplayName != "Google Gemini" {
		t.Errorf("DisplayName = %q, want %q", cfg.DisplayName, "Google Gemini")
	}
	if !cfg.SupportsProjectConfig {
		t.Error("SupportsProjectConfig = false, want true")
	}
	if !strings.Contains(cfg.GlobalConfigPath, filepath.Join(".gemini", "settings.json")) {
		t.Errorf("GlobalConfigPath = %q, want path containing .gemini/settings.json", cfg.GlobalConfigPath)
	}
}

func TestGoogleGeminiProvider_ConfigFilePath_WithProjectRoot_ReturnsProjectPath(t *testing.T) {
	p := providers.NewGoogleGeminiProvider()
	got := p.ConfigFilePath("/any/project")
	want := filepath.Join("/any/project", ".gemini", "settings.json")
	if got != want {
		t.Errorf("ConfigFilePath = %q, want %q", got, want)
	}
}

func TestGoogleGeminiProvider_ConfigFilePath_WithEmptyRoot_ReturnsGlobalPath(t *testing.T) {
	p := providers.NewGoogleGeminiProvider()
	got := p.ConfigFilePath("")

	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(home, ".gemini", "settings.json")
	if got != want {
		t.Errorf("ConfigFilePath = %q, want %q", got, want)
	}
}

func TestGoogleGeminiProvider_Exists_ReturnsFalseWhenAbsent(t *testing.T) {
	tmp := t.TempDir()
	p := providers.NewGoogleGeminiProviderWithPath(filepath.Join(tmp, "settings.json"))
	ok, err := p.Exists(tmp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Error("Exists = true, want false for missing file")
	}
}

func TestGoogleGeminiProvider_Generate_WithEmptyContent_ProducesCorrectJSON(t *testing.T) {
	p := providers.NewGoogleGeminiProvider()
	servers := map[string]types.MCPServer{
		"memory": {
			Command:   "npx",
			Args:      []string{"-y", "@modelcontextprotocol/server-memory"},
			Transport: types.TransportStdio,
		},
	}

	out, err := p.Generate(servers, "")
	if err != nil {
		t.Fatalf("Generate error: %v", err)
	}

	var doc map[string]json.RawMessage
	if err := json.Unmarshal([]byte(out), &doc); err != nil {
		t.Fatalf("output not valid JSON: %v", err)
	}
	if _, ok := doc["mcpServers"]; !ok {
		t.Error("output missing mcpServers key")
	}
}

func TestGoogleGeminiProvider_Generate_WithExistingContent_PreservesNonMCPKeys(t *testing.T) {
	p := providers.NewGoogleGeminiProvider()

	input, err := os.ReadFile("testdata/google_gemini_input.json")
	if err != nil {
		t.Fatal(err)
	}

	servers := map[string]types.MCPServer{
		"memory": {
			Command:   "npx",
			Args:      []string{"-y", "@modelcontextprotocol/server-memory"},
			Transport: types.TransportStdio,
		},
	}

	out, err := p.Generate(servers, string(input))
	if err != nil {
		t.Fatalf("Generate error: %v", err)
	}

	var doc map[string]json.RawMessage
	if err := json.Unmarshal([]byte(out), &doc); err != nil {
		t.Fatalf("output not valid JSON: %v", err)
	}

	if _, ok := doc["theme"]; !ok {
		t.Error("Generate dropped theme key from existing content")
	}
	if _, ok := doc["mcpServers"]; !ok {
		t.Error("Generate missing mcpServers key")
	}
}

func TestGoogleGeminiProvider_Parse_ExtractsServersFromFixture(t *testing.T) {
	p := providers.NewGoogleGeminiProvider()

	content, err := os.ReadFile("testdata/google_gemini_input.json")
	if err != nil {
		t.Fatal(err)
	}

	servers, err := p.Parse(string(content))
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if len(servers) != 1 {
		t.Errorf("len(servers) = %d, want 1", len(servers))
	}
	mem, ok := servers["memory"]
	if !ok {
		t.Fatal("missing memory server")
	}
	if mem.Command != "npx" {
		t.Errorf("memory.Command = %q, want %q", mem.Command, "npx")
	}
}

func TestGoogleGeminiProvider_Parse_MalformedJSON_ReturnsErrMalformedConfig(t *testing.T) {
	p := providers.NewGoogleGeminiProvider()
	_, err := p.Parse("{not valid json")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, ierrors.ErrMalformedConfig) {
		t.Errorf("error = %v, want wrapping ErrMalformedConfig", err)
	}
}

func TestGoogleGeminiProvider_Generate_SSEServer_writesUrl(t *testing.T) {
	p := providers.NewGoogleGeminiProvider()
	servers := map[string]types.MCPServer{
		"remote": {Transport: types.TransportSSE, URL: "https://example.com/sse"},
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
	if entry.URL != "https://example.com/sse" {
		t.Fatalf("expected url=%q, got %q", "https://example.com/sse", entry.URL)
	}
	if entry.HTTPUrl != "" {
		t.Fatalf("expected httpUrl to be empty, got %q", entry.HTTPUrl)
	}
}

func TestGoogleGeminiProvider_Generate_HTTPServer_writesHttpUrl(t *testing.T) {
	p := providers.NewGoogleGeminiProvider()
	servers := map[string]types.MCPServer{
		"remote": {Transport: types.TransportHTTP, URL: "https://example.com/mcp"},
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
	if entry.HTTPUrl != "https://example.com/mcp" {
		t.Fatalf("expected httpUrl=%q, got %q", "https://example.com/mcp", entry.HTTPUrl)
	}
	if entry.URL != "" {
		t.Fatalf("expected url to be empty, got %q", entry.URL)
	}
}

func TestGoogleGeminiProvider_Parse_httpUrl_setsTransportHTTP(t *testing.T) {
	p := providers.NewGoogleGeminiProvider()
	input := `{"mcpServers":{"remote":{"httpUrl":"https://example.com/mcp"}}}`
	parsed, err := p.Parse(input)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	srv := parsed["remote"]
	if srv.Transport != types.TransportHTTP {
		t.Fatalf("expected TransportHTTP, got %q", srv.Transport)
	}
	if srv.URL != "https://example.com/mcp" {
		t.Fatalf("expected URL=%q, got %q", "https://example.com/mcp", srv.URL)
	}
}

func TestGoogleGeminiProvider_Generate_noTypeField(t *testing.T) {
	p := providers.NewGoogleGeminiProvider()
	servers := map[string]types.MCPServer{
		"local": {Transport: types.TransportStdio, Command: "npx", Args: []string{"-y", "server"}},
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
