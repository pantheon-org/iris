package providers_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pantheon-org/iris/internal/ierrors"
	"github.com/pantheon-org/iris/internal/providers"
)

func TestValidateProjectRoot_TraversalDotDot_ReturnsErrPathTraversal(t *testing.T) {
	err := providers.ValidateProjectRoot("/some/project/../../etc")
	if err == nil {
		t.Fatal("expected error for path traversal, got nil")
	}
	if !errors.Is(err, ierrors.ErrPathTraversal) {
		t.Errorf("error = %v, want wrapping ErrPathTraversal", err)
	}
}

func TestValidateProjectRoot_RelativeDotDot_ReturnsErrPathTraversal(t *testing.T) {
	err := providers.ValidateProjectRoot("../../etc/passwd")
	if err == nil {
		t.Fatal("expected error for relative traversal path, got nil")
	}
	if !errors.Is(err, ierrors.ErrPathTraversal) {
		t.Errorf("error = %v, want wrapping ErrPathTraversal", err)
	}
}

func TestValidateProjectRoot_CleanAbsolutePath_ReturnsNil(t *testing.T) {
	err := providers.ValidateProjectRoot("/some/clean/project")
	if err != nil {
		t.Errorf("unexpected error for clean path: %v", err)
	}
}

func TestValidateProjectRoot_EmptyString_ReturnsNil(t *testing.T) {
	// Empty projectRoot is used to select global config path; should be allowed.
	err := providers.ValidateProjectRoot("")
	if err != nil {
		t.Errorf("unexpected error for empty projectRoot: %v", err)
	}
}

func TestClaudeProvider_SafeConfigFilePath_TraversalRoot_ReturnsError(t *testing.T) {
	p := providers.NewClaudeProvider()
	_, err := p.SafeConfigFilePath("/legit/project/../../etc")
	if err == nil {
		t.Fatal("expected error for path traversal, got nil")
	}
	if !errors.Is(err, ierrors.ErrPathTraversal) {
		t.Errorf("error = %v, want wrapping ErrPathTraversal", err)
	}
}

func TestClaudeProvider_SafeConfigFilePath_CleanRoot_ReturnsPath(t *testing.T) {
	p := providers.NewClaudeProvider()
	got, err := p.SafeConfigFilePath("/some/project")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == "" {
		t.Error("expected non-empty path")
	}
}

func TestOpenCodeProvider_Config_UsesExpectedConfigFileName(t *testing.T) {
	cfg := providers.NewOpenCodeProvider().Config()
	assert.NotEmpty(t, cfg.GlobalConfigPath)
	assert.Contains(t, cfg.GlobalConfigPath, "opencode")
	assert.Contains(t, cfg.GlobalConfigPath, "opencode.json")
}

func TestClaudeDesktopProvider_Config_UsesExpectedConfigFileName(t *testing.T) {
	cfg := providers.NewClaudeDesktopProvider().Config()
	assert.NotEmpty(t, cfg.GlobalConfigPath)
	assert.Contains(t, cfg.GlobalConfigPath, "Claude")
	assert.Contains(t, cfg.GlobalConfigPath, "claude_desktop_config.json")
}
