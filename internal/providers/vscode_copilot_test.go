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

func TestVSCodeCopilotProvider_Config(t *testing.T) {
	p := providers.NewVSCodeCopilotProvider()
	cfg := p.Config()
	if cfg.Name != "copilot" {
		t.Fatalf("expected name %q, got %q", "copilot", cfg.Name)
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

func TestVSCodeCopilotProvider_Parse_ExtractsServersFromFixture(t *testing.T) {
	content, err := os.ReadFile("testdata/vscode_copilot_input.json")
	if err != nil {
		t.Fatal(err)
	}
	p := providers.NewVSCodeCopilotProvider()
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
	if parsed["filesystem"].Transport != types.TransportStdio {
		t.Errorf("filesystem.transport = %q, want stdio", parsed["filesystem"].Transport)
	}
	if parsed["context7"].URL != "https://mcp.context7.com/mcp" {
		t.Errorf("context7.url = %q, want URL", parsed["context7"].URL)
	}
}

func TestVSCodeCopilotProvider_Generate_FixtureMatch(t *testing.T) {
	content, err := os.ReadFile("testdata/vscode_copilot_input.json")
	if err != nil {
		t.Fatal(err)
	}
	p := providers.NewVSCodeCopilotProvider()
	servers, err := p.Parse(string(content))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	got, err := p.Generate(servers, string(content))
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	expected, err := os.ReadFile("testdata/vscode_copilot_expected.json")
	if err != nil {
		t.Fatal(err)
	}
	requireJSONEqual(t, string(expected), got)
}

func TestVSCodeCopilotProvider_Parse_malformedInput_returnsError(t *testing.T) {
	p := providers.NewVSCodeCopilotProvider()
	_, err := p.Parse("not json at all")
	if err == nil {
		t.Fatal("Parse: expected error for malformed input, got nil")
	}
	if !errors.Is(err, ierrors.ErrMalformedConfig) {
		t.Errorf("Parse: error does not wrap ErrMalformedConfig; got: %v", err)
	}
}

func TestVSCodeCopilotProvider_Exists(t *testing.T) {
	tmp := t.TempDir()
	p := providers.NewVSCodeCopilotProvider()

	ok, err := p.Exists(tmp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Fatal("should not exist before file is created")
	}

	dir := filepath.Join(tmp, ".vscode")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "mcp.json"), []byte(`{"servers":{}}`), 0o644); err != nil {
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

func keysOf(m map[string]json.RawMessage) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
