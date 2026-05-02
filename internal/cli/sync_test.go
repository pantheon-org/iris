package cli_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pantheon-org/iris/internal/cli"
	"github.com/pantheon-org/iris/internal/providers"
	"github.com/pantheon-org/iris/internal/types"
)

func buildSyncRegistry(t *testing.T, tmpDir string) *providers.Registry {
	t.Helper()
	reg := providers.NewRegistry()
	reg.Register(providers.NewClaudeProvider())
	reg.Register(providers.NewGeminiProviderWithPath(filepath.Join(tmpDir, "gemini-settings.json")))
	reg.Register(providers.NewOpenCodeProvider())
	reg.Register(providers.NewOpenaiCodexProviderWithPath(filepath.Join(tmpDir, "codex.toml")))
	return reg
}

func syncConfig() *types.IrisConfig {
	return &types.IrisConfig{
		Version: 1,
		Servers: map[string]types.MCPServer{
			"fetch": {Transport: types.TransportStdio, Command: "uvx", Args: []string{"mcp-server-fetch"}},
		},
	}
}

func TestRunSync_allCreated_printsCreatedStatus(t *testing.T) {
	dir := t.TempDir()
	reg := buildSyncRegistry(t, dir)
	cfg := syncConfig()
	var buf bytes.Buffer

	err := cli.RunSync(dir, cfg, reg, &buf)

	require.NoError(t, err)
	out := buf.String()
	assert.Contains(t, out, "created")
}

func TestRunSync_unchanged_printsUnchangedStatus(t *testing.T) {
	dir := t.TempDir()
	reg := buildSyncRegistry(t, dir)
	cfg := syncConfig()

	require.NoError(t, cli.RunSync(dir, cfg, reg, &bytes.Buffer{}))

	var buf bytes.Buffer
	err := cli.RunSync(dir, cfg, reg, &buf)

	require.NoError(t, err)
	assert.Contains(t, buf.String(), "unchanged")
}

func TestRunSync_updated_printsUpdatedStatus(t *testing.T) {
	dir := t.TempDir()
	reg := providers.NewRegistry()
	reg.Register(providers.NewClaudeProvider())

	cfg := syncConfig()
	require.NoError(t, cli.RunSync(dir, cfg, reg, &bytes.Buffer{}))

	updatedCfg := &types.IrisConfig{
		Version: 1,
		Servers: map[string]types.MCPServer{
			"fetch": {Transport: types.TransportStdio, Command: "uvx", Args: []string{"mcp-server-fetch"}},
			"extra": {Transport: types.TransportStdio, Command: "node", Args: []string{"extra.js"}},
		},
	}

	var buf bytes.Buffer
	err := cli.RunSync(dir, updatedCfg, reg, &buf)

	require.NoError(t, err)
	assert.Contains(t, buf.String(), "updated")
}

func TestRunSync_providerError_returnsError(t *testing.T) {
	if os.Getuid() == 0 {
		t.Skip("chmod 0555 has no effect as root")
	}

	dir := t.TempDir()

	// Point codex to a file inside a read-only dir that does NOT yet exist,
	// so SyncProvider must attempt to create it (triggering a write error).
	lockedDir := filepath.Join(dir, "locked")
	require.NoError(t, os.MkdirAll(lockedDir, 0o555))
	t.Cleanup(func() { _ = os.Chmod(lockedDir, 0o755) })

	lockedFile := filepath.Join(lockedDir, "codex.toml")

	reg := providers.NewRegistry()
	reg.Register(providers.NewClaudeProvider())
	reg.Register(providers.NewOpenaiCodexProviderWithPath(lockedFile))

	cfg := &types.IrisConfig{
		Version: 1,
		Servers: map[string]types.MCPServer{
			"fetch": {Transport: types.TransportStdio, Command: "uvx", Args: []string{"mcp-server-fetch"}},
		},
	}

	var buf bytes.Buffer
	err := cli.RunSync(dir, cfg, reg, &buf)

	require.Error(t, err)
	assert.Contains(t, buf.String(), "error")
}

func TestRunSync_outputSortedAlphabetically(t *testing.T) {
	dir := t.TempDir()
	reg := buildSyncRegistry(t, dir)
	cfg := syncConfig()
	var buf bytes.Buffer

	require.NoError(t, cli.RunSync(dir, cfg, reg, &buf))

	out := buf.String()
	claudeIdx := indexOf(out, "claude")
	codexIdx := indexOf(out, "codex")
	geminiIdx := indexOf(out, "gemini")
	opencodeIdx := indexOf(out, "opencode")

	assert.True(t, claudeIdx < codexIdx, "claude should appear before codex")
	assert.True(t, codexIdx < geminiIdx, "codex should appear before gemini")
	assert.True(t, geminiIdx < opencodeIdx, "gemini should appear before opencode")
}

func TestRunSync_displaysResolvedProjectPaths(t *testing.T) {
	dir := t.TempDir()
	reg := providers.NewRegistry()
	reg.Register(providers.NewGeminiProvider())
	reg.Register(providers.NewOpenaiCodexProvider())

	var buf bytes.Buffer
	err := cli.RunSync(dir, syncConfig(), reg, &buf)

	require.NoError(t, err)
	out := buf.String()
	assert.Contains(t, out, filepath.Join(dir, ".gemini", "settings.json"))
	assert.Contains(t, out, filepath.Join(dir, ".codex", "config.toml"))
	assert.NotContains(t, out, "~/.gemini/settings.json")
}

func TestRunSync_displaysPinnedProviderPathOnError(t *testing.T) {
	if os.Getuid() == 0 {
		t.Skip("chmod 0555 has no effect as root")
	}

	dir := t.TempDir()
	lockedDir := filepath.Join(dir, "locked")
	require.NoError(t, os.MkdirAll(lockedDir, 0o555))
	t.Cleanup(func() { _ = os.Chmod(lockedDir, 0o755) })

	lockedFile := filepath.Join(lockedDir, "codex.toml")

	reg := providers.NewRegistry()
	reg.Register(providers.NewOpenaiCodexProviderWithPath(lockedFile))

	var buf bytes.Buffer
	err := cli.RunSync(dir, syncConfig(), reg, &buf)

	require.Error(t, err)
	assert.Contains(t, buf.String(), lockedFile)
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
