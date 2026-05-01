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

func TestGeminiProvider_Config_ReturnsCorrectProviderConfig(t *testing.T) {
	p := providers.NewGeminiProvider()
	cfg := p.Config()

	if cfg.Name != "gemini" {
		t.Errorf("Name = %q, want %q", cfg.Name, "gemini")
	}
	if cfg.DisplayName != "Gemini" {
		t.Errorf("DisplayName = %q, want %q", cfg.DisplayName, "Gemini")
	}
	if cfg.SupportsProjectConfig {
		t.Error("SupportsProjectConfig = true, want false")
	}
	if !strings.Contains(cfg.GlobalConfigPath, filepath.Join(".config", "gemini", "settings.json")) {
		t.Errorf("GlobalConfigPath = %q, want path containing .config/gemini/settings.json", cfg.GlobalConfigPath)
	}
}

func TestGeminiProvider_ConfigFilePath_ReturnsGlobalPath(t *testing.T) {
	p := providers.NewGeminiProvider()
	got := p.ConfigFilePath("/any/project")

	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(home, ".config", "gemini", "settings.json")
	if got != want {
		t.Errorf("ConfigFilePath = %q, want %q", got, want)
	}
}

func TestGeminiProvider_Exists_ReturnsFalseWhenAbsent(t *testing.T) {
	tmp := t.TempDir()
	p := providers.NewGeminiProviderWithPath(filepath.Join(tmp, "settings.json"))
	if p.Exists(tmp) {
		t.Error("Exists = true, want false for missing file")
	}
}

func TestGeminiProvider_Generate_WithEmptyContent_ProducesCorrectJSON(t *testing.T) {
	p := providers.NewGeminiProvider()
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

func TestGeminiProvider_Generate_WithExistingContent_PreservesNonMCPKeys(t *testing.T) {
	p := providers.NewGeminiProvider()

	input, err := os.ReadFile("testdata/gemini_input.json")
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

func TestGeminiProvider_Parse_ExtractsServersFromFixture(t *testing.T) {
	p := providers.NewGeminiProvider()

	content, err := os.ReadFile("testdata/gemini_input.json")
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

func TestGeminiProvider_Parse_MalformedJSON_ReturnsErrMalformedConfig(t *testing.T) {
	p := providers.NewGeminiProvider()
	_, err := p.Parse("{not valid json")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, ierrors.ErrMalformedConfig) {
		t.Errorf("error = %v, want wrapping ErrMalformedConfig", err)
	}
}
