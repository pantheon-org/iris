package providers_test

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/pantheon-org/iris/internal/ierrors"
	"github.com/pantheon-org/iris/internal/providers"
	"github.com/pantheon-org/iris/internal/types"
)

func TestIntelliJProvider_Config_ReturnsCorrectProviderConfig(t *testing.T) {
	p := providers.NewIntelliJProvider()
	cfg := p.Config()

	if cfg.Name != "intellij" {
		t.Errorf("Name = %q, want %q", cfg.Name, "intellij")
	}
	if cfg.DisplayName != "IntelliJ IDEA" {
		t.Errorf("DisplayName = %q, want %q", cfg.DisplayName, "IntelliJ IDEA")
	}
	if !cfg.SupportsProjectConfig {
		t.Error("SupportsProjectConfig = false, want true")
	}
	if cfg.LocalConfigPath == nil || *cfg.LocalConfigPath != ".idea/mcp.json" {
		t.Errorf("LocalConfigPath = %v, want %q", cfg.LocalConfigPath, ".idea/mcp.json")
	}
}

func TestIntelliJProvider_ConfigFilePath_ReturnsProjectRelativePath(t *testing.T) {
	p := providers.NewIntelliJProvider()
	got := p.ConfigFilePath("/some/project")
	want := "/some/project/.idea/mcp.json"
	if got != want {
		t.Errorf("ConfigFilePath = %q, want %q", got, want)
	}
}

func TestIntelliJProvider_Exists_ReturnsFalseWhenAbsent(t *testing.T) {
	p := providers.NewIntelliJProvider()
	tmp := t.TempDir()
	ok, err := p.Exists(tmp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Error("Exists = true, want false for missing file")
	}
}

func TestIntelliJProvider_Exists_ReturnsTrueWhenPresent(t *testing.T) {
	p := providers.NewIntelliJProvider()
	tmp := t.TempDir()
	if err := os.MkdirAll(filepath.Join(tmp, ".idea"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmp, ".idea", "mcp.json"), []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	ok, err := p.Exists(tmp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Error("Exists = false, want true when file present")
	}
}

func TestIntelliJProvider_Generate_WithEmptyContent_ProducesCorrectJSON(t *testing.T) {
	p := providers.NewIntelliJProvider()
	servers := map[string]types.MCPServer{
		"my-server": {
			Command:   "npx",
			Args:      []string{"-y", "@modelcontextprotocol/server-filesystem"},
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

func TestIntelliJProvider_Generate_WithExistingContent_PreservesNonMCPKeys(t *testing.T) {
	p := providers.NewIntelliJProvider()

	input, err := os.ReadFile("testdata/intellij_input.json")
	if err != nil {
		t.Fatal(err)
	}

	servers := map[string]types.MCPServer{
		"filesystem": {
			Command:   "npx",
			Args:      []string{"-y", "@modelcontextprotocol/server-filesystem"},
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

	if _, ok := doc["extraKey"]; !ok {
		t.Error("Generate dropped extraKey from existing content")
	}
	if _, ok := doc["mcpServers"]; !ok {
		t.Error("Generate missing mcpServers key")
	}
}

func TestIntelliJProvider_Parse_ExtractsServersFromFixture(t *testing.T) {
	p := providers.NewIntelliJProvider()

	content, err := os.ReadFile("testdata/intellij_input.json")
	if err != nil {
		t.Fatal(err)
	}

	servers, err := p.Parse(string(content))
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if len(servers) != 2 {
		t.Errorf("len(servers) = %d, want 2", len(servers))
	}
	fs, ok := servers["filesystem"]
	if !ok {
		t.Fatal("missing filesystem server")
	}
	if fs.Command != "npx" {
		t.Errorf("filesystem.Command = %q, want %q", fs.Command, "npx")
	}
}

func TestIntelliJProvider_Generate_FixtureMatch(t *testing.T) {
	p := providers.NewIntelliJProvider()
	content, err := os.ReadFile("testdata/intellij_input.json")
	if err != nil {
		t.Fatal(err)
	}
	servers, err := p.Parse(string(content))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	got, err := p.Generate(servers, string(content))
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	expected, err := os.ReadFile("testdata/intellij_expected.json")
	if err != nil {
		t.Fatal(err)
	}
	if got != string(expected) {
		t.Errorf("output mismatch:\ngot:\n%s\nwant:\n%s", got, expected)
	}
}

func TestIntelliJProvider_Parse_MalformedJSON_ReturnsErrMalformedConfig(t *testing.T) {
	p := providers.NewIntelliJProvider()
	_, err := p.Parse("{not valid json")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, ierrors.ErrMalformedConfig) {
		t.Errorf("error = %v, want wrapping ErrMalformedConfig", err)
	}
}
