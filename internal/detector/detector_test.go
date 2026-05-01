package detector_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pantheon-org/iris/internal/detector"
	"github.com/pantheon-org/iris/internal/providers"
)

func newTestRegistry(t *testing.T) *providers.Registry {
	t.Helper()
	tmp := t.TempDir()

	reg := providers.NewRegistry()
	reg.Register(providers.NewClaudeProvider())
	reg.Register(providers.NewOpenCodeProvider())
	reg.Register(providers.NewGeminiProviderWithPath(filepath.Join(tmp, "gemini-settings.json")))
	reg.Register(providers.NewCodexProviderWithPath(filepath.Join(tmp, "codex-config.toml")))
	return reg
}

func TestDetect_EmptyDir_NoProvidersDetected(t *testing.T) {
	root := t.TempDir()
	reg := newTestRegistry(t)

	got := detector.Detect(root, reg)

	if len(got) != 0 {
		t.Errorf("expected 0 providers, got %d", len(got))
	}
}

func TestDetect_MCPJsonPresent_ClaudeDetected(t *testing.T) {
	root := t.TempDir()
	reg := newTestRegistry(t)

	if err := os.WriteFile(filepath.Join(root, ".mcp.json"), []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}

	got := detector.Detect(root, reg)

	if len(got) != 1 {
		t.Fatalf("expected 1 provider, got %d", len(got))
	}
	if got[0].Config().Name != "claude" {
		t.Errorf("expected claude, got %q", got[0].Config().Name)
	}
}

func TestDetect_OpenCodeJsonPresent_OpenCodeDetected(t *testing.T) {
	root := t.TempDir()
	reg := newTestRegistry(t)

	if err := os.WriteFile(filepath.Join(root, "opencode.json"), []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}

	got := detector.Detect(root, reg)

	if len(got) != 1 {
		t.Fatalf("expected 1 provider, got %d", len(got))
	}
	if got[0].Config().Name != "opencode" {
		t.Errorf("expected opencode, got %q", got[0].Config().Name)
	}
}

func TestDetect_BothProjectFilesPresent_BothDetected(t *testing.T) {
	root := t.TempDir()
	reg := newTestRegistry(t)

	if err := os.WriteFile(filepath.Join(root, ".mcp.json"), []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "opencode.json"), []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}

	got := detector.Detect(root, reg)

	if len(got) != 2 {
		t.Fatalf("expected 2 providers, got %d", len(got))
	}
	names := map[string]bool{}
	for _, p := range got {
		names[p.Config().Name] = true
	}
	if !names["claude"] {
		t.Error("expected claude in results")
	}
	if !names["opencode"] {
		t.Error("expected opencode in results")
	}
}

func TestDetect_GeminiGlobalFile_NeverDetected(t *testing.T) {
	root := t.TempDir()
	geminiTmp := t.TempDir()

	reg := providers.NewRegistry()
	geminiPath := filepath.Join(geminiTmp, "gemini-settings.json")
	if err := os.WriteFile(geminiPath, []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	reg.Register(providers.NewGeminiProviderWithPath(geminiPath))

	got := detector.Detect(root, reg)

	if len(got) != 0 {
		t.Errorf("expected 0 providers (gemini is global-only), got %d", len(got))
	}
}

func TestDetect_CodexGlobalFile_NeverDetected(t *testing.T) {
	root := t.TempDir()
	codexTmp := t.TempDir()

	reg := providers.NewRegistry()
	codexPath := filepath.Join(codexTmp, "codex-config.toml")
	if err := os.WriteFile(codexPath, []byte(""), 0o644); err != nil {
		t.Fatal(err)
	}
	reg.Register(providers.NewCodexProviderWithPath(codexPath))

	got := detector.Detect(root, reg)

	if len(got) != 0 {
		t.Errorf("expected 0 providers (codex is global-only), got %d", len(got))
	}
}
