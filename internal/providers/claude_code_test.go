package providers_test

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/pantheon-org/iris/internal/ierrors"
	"github.com/pantheon-org/iris/internal/providers"
	"github.com/pantheon-org/iris/internal/types"
)

func TestClaudeCodeProvider_Config_ReturnsCorrectProviderConfig(t *testing.T) {
	p := providers.NewClaudeCodeProvider()
	cfg := p.Config()

	if cfg.Name != "anthropic-claude-code" {
		t.Errorf("Name = %q, want %q", cfg.Name, "anthropic-claude-code")
	}
	if cfg.DisplayName != "Anthropic Claude Code" {
		t.Errorf("DisplayName = %q, want %q", cfg.DisplayName, "Anthropic Claude Code")
	}
	if !cfg.SupportsProjectConfig {
		t.Error("SupportsProjectConfig = false, want true")
	}
	if cfg.ConfigPath != ".mcp.json" {
		t.Errorf("ConfigPath = %q, want %q", cfg.ConfigPath, ".mcp.json")
	}
}

func TestClaudeCodeProvider_ConfigFilePath_ReturnsProjectRelativePath(t *testing.T) {
	p := providers.NewClaudeCodeProvider()
	got := p.ConfigFilePath("/some/project")
	want := "/some/project/.mcp.json"
	if got != want {
		t.Errorf("ConfigFilePath = %q, want %q", got, want)
	}
}

func TestClaudeCodeProvider_Exists_ReturnsFalseWhenAbsent(t *testing.T) {
	p := providers.NewClaudeCodeProvider()
	tmp := t.TempDir()
	ok, err := p.Exists(tmp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Error("Exists = true, want false for missing file")
	}
}

func TestClaudeCodeProvider_Exists_ReturnsTrueWhenPresent(t *testing.T) {
	p := providers.NewClaudeCodeProvider()
	tmp := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmp, ".mcp.json"), []byte("{}"), 0o644); err != nil {
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

func TestClaudeCodeProvider_Generate_WithEmptyContent_ProducesCorrectJSON(t *testing.T) {
	p := providers.NewClaudeCodeProvider()
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

func TestClaudeCodeProvider_Generate_WithExistingContent_PreservesNonMCPKeys(t *testing.T) {
	p := providers.NewClaudeCodeProvider()

	input, err := os.ReadFile("testdata/claude_code_input.json")
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

func TestClaudeCodeProvider_Parse_ExtractsServersFromFixture(t *testing.T) {
	p := providers.NewClaudeCodeProvider()

	content, err := os.ReadFile("testdata/claude_code_input.json")
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

func TestClaudeCodeProvider_Parse_MalformedJSON_ReturnsErrMalformedConfig(t *testing.T) {
	p := providers.NewClaudeCodeProvider()
	_, err := p.Parse("{not valid json")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, ierrors.ErrMalformedConfig) {
		t.Errorf("error = %v, want wrapping ErrMalformedConfig", err)
	}
}

func TestClaudeCodeProvider_GenerateParse_PreservesRemoteServerFields(t *testing.T) {
	p := providers.NewClaudeCodeProvider()
	enabled := false
	servers := map[string]types.MCPServer{
		"remote": {
			Transport: types.TransportSSE,
			URL:       "https://example.com/mcp",
			Headers:   map[string]string{"Authorization": "Bearer token"},
			Cwd:       "/tmp/work",
			Enabled:   &enabled,
		},
	}

	out, err := p.Generate(servers, "")
	if err != nil {
		t.Fatalf("Generate error: %v", err)
	}

	parsed, err := p.Parse(out)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if !reflect.DeepEqual(parsed, servers) {
		t.Fatalf("roundtrip mismatch:\n got: %#v\nwant: %#v", parsed, servers)
	}
}
