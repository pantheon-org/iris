package providers_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/pantheon-org/iris/internal/providers"
	"github.com/pantheon-org/iris/internal/types"
)

func TestVSCodeCopilotProvider_Config(t *testing.T) {
	p := providers.NewVSCodeCopilotProvider()
	cfg := p.Config()
	if cfg.Name != "vscode-copilot" {
		t.Fatalf("expected name %q, got %q", "vscode-copilot", cfg.Name)
	}
	if !cfg.SupportsProjectConfig {
		t.Fatal("expected SupportsProjectConfig=true")
	}
}

func TestVSCodeCopilotProvider_ConfigFilePath(t *testing.T) {
	p := providers.NewVSCodeCopilotProvider()
	got := p.ConfigFilePath("/project")
	want := filepath.Join("/project", ".vscode", "mcp.json")
	if got != want {
		t.Fatalf("want %q, got %q", want, got)
	}
}

func TestVSCodeCopilotProvider_Generate_usesServersKey(t *testing.T) {
	p := providers.NewVSCodeCopilotProvider()
	servers := map[string]types.MCPServer{
		"fs": {Transport: types.TransportStdio, Command: "node", Args: []string{"server.js"}},
	}
	out, err := p.Generate(servers, "")
	if err != nil {
		t.Fatal(err)
	}

	var doc map[string]json.RawMessage
	if err := json.Unmarshal([]byte(out), &doc); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out)
	}
	if _, ok := doc["servers"]; !ok {
		t.Fatalf("expected top-level 'servers' key, got keys: %v", keysOf(doc))
	}
	if _, ok := doc["mcpServers"]; ok {
		t.Fatal("should NOT have 'mcpServers' key")
	}
}

func TestVSCodeCopilotProvider_Generate_hasTypeField(t *testing.T) {
	p := providers.NewVSCodeCopilotProvider()
	servers := map[string]types.MCPServer{
		"fs": {Transport: types.TransportStdio, Command: "node"},
	}
	out, err := p.Generate(servers, "")
	if err != nil {
		t.Fatal(err)
	}

	var doc struct {
		Servers map[string]struct {
			Type string `json:"type"`
		} `json:"servers"`
	}
	if err := json.Unmarshal([]byte(out), &doc); err != nil {
		t.Fatal(err)
	}
	if doc.Servers["fs"].Type != "stdio" {
		t.Fatalf("expected type=stdio, got %q", doc.Servers["fs"].Type)
	}
}

func TestVSCodeCopilotProvider_Parse(t *testing.T) {
	p := providers.NewVSCodeCopilotProvider()
	input := `{"servers":{"fs":{"type":"stdio","command":"node","args":["server.js"]}}}`
	parsed, err := p.Parse(input)
	if err != nil {
		t.Fatal(err)
	}
	if parsed["fs"].Command != "node" {
		t.Fatalf("expected command %q, got %q", "node", parsed["fs"].Command)
	}
	if parsed["fs"].Transport != types.TransportStdio {
		t.Fatalf("expected transport stdio, got %q", parsed["fs"].Transport)
	}
}

func TestVSCodeCopilotProvider_GenerateParse_preservesExistingKeys(t *testing.T) {
	p := providers.NewVSCodeCopilotProvider()
	existing := `{"inputs":[],"servers":{}}`
	servers := map[string]types.MCPServer{
		"fs": {Transport: types.TransportStdio, Command: "node"},
	}
	out, err := p.Generate(servers, existing)
	if err != nil {
		t.Fatal(err)
	}

	var doc map[string]json.RawMessage
	if err := json.Unmarshal([]byte(out), &doc); err != nil {
		t.Fatal(err)
	}
	if _, ok := doc["inputs"]; !ok {
		t.Fatal("expected 'inputs' key to be preserved")
	}
}

func TestVSCodeCopilotProvider_Exists(t *testing.T) {
	tmp := t.TempDir()
	p := providers.NewVSCodeCopilotProvider()

	if p.Exists(tmp) {
		t.Fatal("should not exist before file is created")
	}

	dir := filepath.Join(tmp, ".vscode")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "mcp.json"), []byte(`{"servers":{}}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if !p.Exists(tmp) {
		t.Fatal("should exist after file is created")
	}
}

func keysOf(m map[string]json.RawMessage) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
