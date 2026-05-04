package providers_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/pantheon-org/iris/internal/providers"
	"github.com/pantheon-org/iris/internal/types"
)

func TestZedProvider_Config(t *testing.T) {
	tmp := t.TempDir()
	p := providers.NewZedProviderWithPath(filepath.Join(tmp, "settings.json"))
	cfg := p.Config()
	if cfg.Name != "zed" {
		t.Fatalf("expected name %q, got %q", "zed", cfg.Name)
	}
	if cfg.SupportsProjectConfig {
		t.Fatal("expected SupportsProjectConfig=false")
	}
}

func TestZedProvider_Generate_usesContextServersKey(t *testing.T) {
	tmp := t.TempDir()
	p := providers.NewZedProviderWithPath(filepath.Join(tmp, "settings.json"))
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
	if _, ok := doc["context_servers"]; !ok {
		t.Fatalf("expected 'context_servers' key, got keys: %v", keysOf(doc))
	}
}

func TestZedProvider_Parse(t *testing.T) {
	tmp := t.TempDir()
	p := providers.NewZedProviderWithPath(filepath.Join(tmp, "settings.json"))
	input := `{"context_servers":{"fs":{"command":"node","args":["server.js"],"env":{}}}}`
	parsed, err := p.Parse(input)
	if err != nil {
		t.Fatal(err)
	}
	if parsed["fs"].Command != "node" {
		t.Fatalf("expected command %q, got %q", "node", parsed["fs"].Command)
	}
}

func TestZedProvider_GenerateParse_preservesExistingKeys(t *testing.T) {
	tmp := t.TempDir()
	p := providers.NewZedProviderWithPath(filepath.Join(tmp, "settings.json"))
	existing := `{"theme":"One Dark","context_servers":{}}`
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
	if _, ok := doc["theme"]; !ok {
		t.Fatal("expected 'theme' key to be preserved")
	}
}

func TestZedProvider_Exists(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "settings.json")
	p := providers.NewZedProviderWithPath(path)

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
