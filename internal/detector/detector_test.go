package detector_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/pantheon-org/iris/internal/detector"
	"github.com/pantheon-org/iris/internal/providers"
	"github.com/pantheon-org/iris/internal/registry"
	"github.com/pantheon-org/iris/internal/types"
)

func newTestRegistry(t *testing.T) *registry.Registry {
	t.Helper()
	tmp := t.TempDir()

	reg := registry.NewRegistry()
	reg.Register(providers.NewClaudeCodeProvider())
	reg.Register(providers.NewOpenCodeProvider())
	reg.Register(providers.NewGeminiProviderWithPath(filepath.Join(tmp, "gemini-settings.json")))
	reg.Register(providers.NewOpenaiCodexProviderWithPath(filepath.Join(tmp, "codex-config.toml")))
	return reg
}

func TestDetect_EmptyDir_NoProvidersDetected(t *testing.T) {
	root := t.TempDir()
	reg := newTestRegistry(t)

	got, err := detector.Detect(root, reg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

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

	got, err := detector.Detect(root, reg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

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

	got, err := detector.Detect(root, reg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

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

	got, err := detector.Detect(root, reg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

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

func TestDetect_GeminiProjectConfig_Detected(t *testing.T) {
	root := t.TempDir()

	reg := registry.NewRegistry()
	reg.Register(providers.NewGeminiProvider())

	geminiDir := filepath.Join(root, ".gemini")
	if err := os.MkdirAll(geminiDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(geminiDir, "settings.json"), []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}

	got, err := detector.Detect(root, reg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(got) != 1 {
		t.Fatalf("expected 1 provider, got %d", len(got))
	}
	if got[0].Config().Name != "gemini" {
		t.Errorf("expected gemini, got %q", got[0].Config().Name)
	}
}

func TestDetect_GeminiProjectConfig_AbsentNotDetected(t *testing.T) {
	root := t.TempDir()

	reg := registry.NewRegistry()
	reg.Register(providers.NewGeminiProvider())

	got, err := detector.Detect(root, reg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(got) != 0 {
		t.Errorf("expected 0 providers, got %d", len(got))
	}
}

func TestDetect_CodexProjectConfig_Detected(t *testing.T) {
	root := t.TempDir()

	reg := registry.NewRegistry()
	reg.Register(providers.NewOpenaiCodexProvider())

	codexDir := filepath.Join(root, ".codex")
	if err := os.MkdirAll(codexDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(codexDir, "config.toml"), []byte(""), 0o644); err != nil {
		t.Fatal(err)
	}

	got, err := detector.Detect(root, reg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(got) != 1 {
		t.Fatalf("expected 1 provider, got %d", len(got))
	}
	if got[0].Config().Name != "codex" {
		t.Errorf("expected codex, got %q", got[0].Config().Name)
	}
}

func TestDetect_CodexProjectConfig_AbsentNotDetected(t *testing.T) {
	root := t.TempDir()

	reg := registry.NewRegistry()
	reg.Register(providers.NewOpenaiCodexProvider())

	got, err := detector.Detect(root, reg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(got) != 0 {
		t.Errorf("expected 0 providers, got %d", len(got))
	}
}

func TestDetect_QwenProjectConfig_Detected(t *testing.T) {
	root := t.TempDir()

	reg := registry.NewRegistry()
	reg.Register(providers.NewQwenProvider())

	qwenDir := filepath.Join(root, ".qwen")
	if err := os.MkdirAll(qwenDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(qwenDir, "settings.json"), []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}

	got, err := detector.Detect(root, reg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(got) != 1 {
		t.Fatalf("expected 1 provider, got %d", len(got))
	}
	if got[0].Config().Name != "qwen" {
		t.Errorf("expected qwen, got %q", got[0].Config().Name)
	}
}

func TestDetect_QwenProjectConfig_AbsentNotDetected(t *testing.T) {
	root := t.TempDir()

	reg := registry.NewRegistry()
	reg.Register(providers.NewQwenProvider())

	got, err := detector.Detect(root, reg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(got) != 0 {
		t.Errorf("expected 0 providers, got %d", len(got))
	}
}

func TestDetect_MistralVibeProjectConfig_Detected(t *testing.T) {
	root := t.TempDir()

	reg := registry.NewRegistry()
	reg.Register(providers.NewMistralVibeProvider())

	vibeDir := filepath.Join(root, ".vibe")
	if err := os.MkdirAll(vibeDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(vibeDir, "config.toml"), []byte(""), 0o644); err != nil {
		t.Fatal(err)
	}

	got, err := detector.Detect(root, reg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(got) != 1 {
		t.Fatalf("expected 1 provider, got %d", len(got))
	}
	if got[0].Config().Name != "mistral-vibe" {
		t.Errorf("expected mistral-vibe, got %q", got[0].Config().Name)
	}
}

func TestDetect_MistralVibeProjectConfig_AbsentNotDetected(t *testing.T) {
	root := t.TempDir()

	reg := registry.NewRegistry()
	reg.Register(providers.NewMistralVibeProvider())

	got, err := detector.Detect(root, reg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(got) != 0 {
		t.Errorf("expected 0 providers, got %d", len(got))
	}
}

// ioErrorProvider is a test stub whose Exists always returns an IO error.
type ioErrorProvider struct {
	existsErr error
}

func (p *ioErrorProvider) Config() providers.ProviderConfig {
	return providers.ProviderConfig{
		Name:                  "io-error-stub",
		DisplayName:           "IO Error Stub",
		SupportsProjectConfig: true,
	}
}

func (p *ioErrorProvider) Exists(_ string) (bool, error) { return false, p.existsErr }

func (p *ioErrorProvider) ConfigFilePath(_ string) string { return "" }

func (p *ioErrorProvider) Generate(_ map[string]types.MCPServer, _ string) (string, error) {
	return "", nil
}

func (p *ioErrorProvider) Parse(_ string) (map[string]types.MCPServer, error) { return nil, nil }

func TestDetect_ExistsIOError_SurfacedToCallerNotSwallowed(t *testing.T) {
	root := t.TempDir()

	sentinel := errors.New("permission denied")
	stub := &ioErrorProvider{existsErr: sentinel}

	reg := registry.NewRegistry()
	reg.Register(stub)

	_, err := detector.Detect(root, reg)
	if err == nil {
		t.Fatal("expected error from Detect, got nil")
	}
	if !errors.Is(err, sentinel) {
		t.Errorf("expected sentinel error in chain, got: %v", err)
	}
}
