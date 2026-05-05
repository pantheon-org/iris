package providers_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pantheon-org/iris/internal/ierrors"
	"github.com/pantheon-org/iris/internal/providers"
)

// requireJSONEqual asserts that two JSON strings are semantically equal,
// ignoring field ordering differences produced by different marshal paths.
func requireJSONEqual(t *testing.T, want, got string) {
	t.Helper()
	var wantVal, gotVal any
	require.NoError(t, json.Unmarshal([]byte(want), &wantVal), "want: invalid JSON")
	require.NoError(t, json.Unmarshal([]byte(got), &gotVal), "got: invalid JSON")
	assert.Equal(t, wantVal, gotVal)
}

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
	p := providers.NewClaudeCodeProvider()
	_, err := p.SafeConfigFilePath("/legit/project/../../etc")
	if err == nil {
		t.Fatal("expected error for path traversal, got nil")
	}
	if !errors.Is(err, ierrors.ErrPathTraversal) {
		t.Errorf("error = %v, want wrapping ErrPathTraversal", err)
	}
}

func TestClaudeProvider_SafeConfigFilePath_CleanRoot_ReturnsPath(t *testing.T) {
	p := providers.NewClaudeCodeProvider()
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
	require.NotNil(t, cfg.GlobalConfigPath)
	assert.Contains(t, *cfg.GlobalConfigPath, "opencode")
	assert.Contains(t, *cfg.GlobalConfigPath, "opencode.json")
}

func TestClaudeDesktopProvider_Config_UsesExpectedConfigFileName(t *testing.T) {
	cfg := providers.NewClaudeDesktopProvider().Config()
	require.NotNil(t, cfg.GlobalConfigPath)
	assert.Contains(t, *cfg.GlobalConfigPath, "Claude")
	assert.Contains(t, *cfg.GlobalConfigPath, "claude_desktop_config.json")
}
