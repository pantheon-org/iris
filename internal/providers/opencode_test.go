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

func TestOpenCodeProvider_Config_ReturnsCorrectConfig(t *testing.T) {
	p := providers.NewOpenCodeProvider()
	cfg := p.Config()

	if cfg.Name != "opencode" {
		t.Errorf("Name: got %q, want %q", cfg.Name, "opencode")
	}
	if cfg.ConfigPath != "opencode.json" {
		t.Errorf("ConfigPath: got %q, want %q", cfg.ConfigPath, "opencode.json")
	}
	if !cfg.SupportsProjectConfig {
		t.Error("SupportsProjectConfig: want true")
	}
	if cfg.GlobalConfigPath == "" {
		t.Error("GlobalConfigPath: want non-empty")
	}
}

func TestOpenCodeProvider_ConfigFilePath_ReturnsProjectPath(t *testing.T) {
	p := providers.NewOpenCodeProvider()
	got := p.ConfigFilePath("/some/project")
	want := filepath.Join("/some/project", "opencode.json")
	if got != want {
		t.Errorf("ConfigFilePath: got %q, want %q", got, want)
	}
}

func TestOpenCodeProvider_Exists_ReturnsFalseWhenAbsent(t *testing.T) {
	p := providers.NewOpenCodeProvider()
	dir := t.TempDir()
	if p.Exists(dir) {
		t.Error("Exists: want false for missing file")
	}
}

func TestOpenCodeProvider_Exists_ReturnsTrueWhenPresent(t *testing.T) {
	p := providers.NewOpenCodeProvider()
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "opencode.json"), []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	if !p.Exists(dir) {
		t.Error("Exists: want true when file present")
	}
}

func TestOpenCodeProvider_Generate_EmptyExistingContent_ProducesCorrectJSON(t *testing.T) {
	p := providers.NewOpenCodeProvider()
	enabled := true
	servers := map[string]types.MCPServer{
		"my-server": {
			Command: "npx",
			Args:    []string{"-y", "@modelcontextprotocol/server-filesystem"},
			Enabled: &enabled,
			Env:     map[string]string{"FOO": "bar"},
		},
	}

	out, err := p.Generate(servers, "")
	if err != nil {
		t.Fatalf("Generate: unexpected error: %v", err)
	}

	var result map[string]json.RawMessage
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("Generate output is not valid JSON: %v", err)
	}
	if _, ok := result["mcp"]; !ok {
		t.Error("Generate: missing top-level 'mcp' key")
	}

	var mcp map[string]struct {
		Command     []string          `json:"command"`
		Type        string            `json:"type"`
		Enabled     bool              `json:"enabled"`
		Environment map[string]string `json:"environment"`
	}
	if err := json.Unmarshal(result["mcp"], &mcp); err != nil {
		t.Fatalf("Generate: cannot parse mcp block: %v", err)
	}

	srv, ok := mcp["my-server"]
	if !ok {
		t.Fatal("Generate: missing 'my-server' in mcp")
	}
	if len(srv.Command) < 1 || srv.Command[0] != "npx" {
		t.Errorf("Generate: command[0] = %q, want %q", srv.Command[0], "npx")
	}
	if len(srv.Command) != 3 {
		t.Errorf("Generate: command length = %d, want 3 (cmd+2 args)", len(srv.Command))
	}
	if srv.Type != "local" {
		t.Errorf("Generate: type = %q, want %q", srv.Type, "local")
	}
	if !srv.Enabled {
		t.Error("Generate: enabled = false, want true")
	}
	if srv.Environment["FOO"] != "bar" {
		t.Errorf("Generate: environment[FOO] = %q, want %q", srv.Environment["FOO"], "bar")
	}
}

func TestOpenCodeProvider_Generate_NilEnabled_DefaultsToTrue(t *testing.T) {
	p := providers.NewOpenCodeProvider()
	servers := map[string]types.MCPServer{
		"no-flag": {Command: "echo"},
	}

	out, err := p.Generate(servers, "")
	if err != nil {
		t.Fatalf("Generate: unexpected error: %v", err)
	}

	var root struct {
		MCP map[string]struct {
			Enabled bool `json:"enabled"`
		} `json:"mcp"`
	}
	if err := json.Unmarshal([]byte(out), &root); err != nil {
		t.Fatalf("cannot parse output: %v", err)
	}
	if !root.MCP["no-flag"].Enabled {
		t.Error("Generate: nil Enabled should default to true")
	}
}

func TestOpenCodeProvider_Generate_PreservesNonMCPKeys(t *testing.T) {
	p := providers.NewOpenCodeProvider()
	existing := `{"theme":"dark","fontSize":14,"mcp":{}}`
	servers := map[string]types.MCPServer{
		"s": {Command: "cmd"},
	}

	out, err := p.Generate(servers, existing)
	if err != nil {
		t.Fatalf("Generate: unexpected error: %v", err)
	}

	var result map[string]json.RawMessage
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if _, ok := result["theme"]; !ok {
		t.Error("Generate: 'theme' key lost when preserving existing content")
	}
	if _, ok := result["fontSize"]; !ok {
		t.Error("Generate: 'fontSize' key lost when preserving existing content")
	}
}

func TestOpenCodeProvider_Generate_MergesCommandAndArgs(t *testing.T) {
	p := providers.NewOpenCodeProvider()
	servers := map[string]types.MCPServer{
		"srv": {Command: "node", Args: []string{"server.js", "--port", "3000"}},
	}

	out, err := p.Generate(servers, "")
	if err != nil {
		t.Fatalf("Generate: unexpected error: %v", err)
	}

	var root struct {
		MCP map[string]struct {
			Command []string `json:"command"`
		} `json:"mcp"`
	}
	if err := json.Unmarshal([]byte(out), &root); err != nil {
		t.Fatalf("cannot parse output: %v", err)
	}
	want := []string{"node", "server.js", "--port", "3000"}
	got := root.MCP["srv"].Command
	if len(got) != len(want) {
		t.Fatalf("command length: got %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("command[%d]: got %q, want %q", i, got[i], want[i])
		}
	}
}

func TestOpenCodeProvider_Parse_CorrectlyExtractsServers(t *testing.T) {
	input, err := os.ReadFile("testdata/opencode_input.json")
	if err != nil {
		t.Fatalf("reading fixture: %v", err)
	}

	p := providers.NewOpenCodeProvider()
	servers, err := p.Parse(string(input))
	if err != nil {
		t.Fatalf("Parse: unexpected error: %v", err)
	}

	if len(servers) != 2 {
		t.Fatalf("Parse: got %d servers, want 2", len(servers))
	}

	fs, ok := servers["filesystem"]
	if !ok {
		t.Fatal("Parse: missing 'filesystem' server")
	}
	if fs.Command != "npx" {
		t.Errorf("Parse: Command = %q, want %q", fs.Command, "npx")
	}
	wantArgs := []string{"-y", "@modelcontextprotocol/server-filesystem", "/tmp"}
	if len(fs.Args) != len(wantArgs) {
		t.Fatalf("Parse: Args length = %d, want %d", len(fs.Args), len(wantArgs))
	}
	for i, a := range wantArgs {
		if fs.Args[i] != a {
			t.Errorf("Parse: Args[%d] = %q, want %q", i, fs.Args[i], a)
		}
	}

	fetch, ok := servers["fetch"]
	if !ok {
		t.Fatal("Parse: missing 'fetch' server")
	}
	if fetch.Env["PROXY"] != "http://proxy:8080" {
		t.Errorf("Parse: Env[PROXY] = %q, want %q", fetch.Env["PROXY"], "http://proxy:8080")
	}
}

func TestOpenCodeProvider_Parse_MalformedJSON_ReturnsWrappedError(t *testing.T) {
	p := providers.NewOpenCodeProvider()
	_, err := p.Parse(`{not valid json`)
	if err == nil {
		t.Fatal("Parse: expected error for malformed JSON, got nil")
	}
	if !errors.Is(err, ierrors.ErrMalformedConfig) {
		t.Errorf("Parse: error does not wrap ErrMalformedConfig; got: %v", err)
	}
}
