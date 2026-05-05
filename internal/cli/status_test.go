package cli_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pantheon-org/iris/internal/cli"
	irio "github.com/pantheon-org/iris/internal/io"
	"github.com/pantheon-org/iris/internal/providers"
	"github.com/pantheon-org/iris/internal/registry"
	"github.com/pantheon-org/iris/internal/types"
)

func buildTestRegistry(t *testing.T, tmpDir string) *registry.Registry {
	t.Helper()
	reg := registry.NewRegistry()
	reg.Register(providers.NewClaudeCodeProvider())
	reg.Register(providers.NewGoogleGeminiProviderWithPath(filepath.Join(tmpDir, "gemini-settings.json")))
	reg.Register(providers.NewOpenCodeProvider())
	reg.Register(providers.NewOpenaiCodexProviderWithPath(filepath.Join(tmpDir, "codex.toml")))
	return reg
}

func minimalConfig() *types.IrisConfig {
	return &types.IrisConfig{
		Version: 1,
		Servers: map[string]types.MCPServer{
			"fetch": {Transport: types.TransportStdio, Command: "uvx", Args: []string{"mcp-server-fetch"}},
		},
	}
}

func TestRunStatus_allMissing_showsMissing(t *testing.T) {
	dir := t.TempDir()
	reg := buildTestRegistry(t, dir)
	cfg := minimalConfig()
	var buf bytes.Buffer

	err := cli.RunStatus(dir, cfg, reg, &buf, false, noColourStyles())

	require.NoError(t, err)
	out := buf.String()
	assert.Contains(t, out, "Provider Status:")
	assert.Contains(t, out, "missing")
}

func TestRunStatus_filePresent_synced(t *testing.T) {
	dir := t.TempDir()
	reg := registry.NewRegistry()
	reg.Register(providers.NewClaudeCodeProvider())

	cfg := minimalConfig()

	// Generate the expected content and write it to the file.
	p, err := reg.Get("claude")
	require.NoError(t, err)
	content, err := p.Generate(cfg.Servers, "")
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(filepath.Join(dir, ".mcp.json"), []byte(content), 0o644))

	var buf bytes.Buffer
	err = cli.RunStatus(dir, cfg, reg, &buf, false, noColourStyles())

	require.NoError(t, err)
	assert.Contains(t, buf.String(), "synced")
}

func TestRunStatus_filePresent_desync(t *testing.T) {
	dir := t.TempDir()
	reg := registry.NewRegistry()
	reg.Register(providers.NewClaudeCodeProvider())

	cfg := minimalConfig()

	// Write stale content (empty JSON object) to trigger desync.
	stale, _ := json.Marshal(map[string]interface{}{})
	require.NoError(t, os.WriteFile(filepath.Join(dir, ".mcp.json"), stale, 0o644))

	var buf bytes.Buffer
	err := cli.RunStatus(dir, cfg, reg, &buf, false, noColourStyles())

	require.NoError(t, err)
	assert.Contains(t, buf.String(), "desync")
}

func TestRunStatus_readFailure_showsError(t *testing.T) {
	dir := t.TempDir()
	reg := registry.NewRegistry()
	reg.Register(providers.NewClaudeCodeProvider())

	cfg := minimalConfig()

	require.NoError(t, os.Mkdir(filepath.Join(dir, ".mcp.json"), 0o755))

	var buf bytes.Buffer
	err := cli.RunStatus(dir, cfg, reg, &buf, false, noColourStyles())

	require.NoError(t, err)
	assert.Contains(t, buf.String(), "error")
	assert.NotContains(t, buf.String(), "missing")
}

func TestRunStatus_displaysResolvedProjectPaths(t *testing.T) {
	dir := t.TempDir()
	home := irio.UserHomeDir()
	reg := registry.NewRegistry()
	reg.Register(providers.NewGoogleGeminiProvider())
	reg.Register(providers.NewOpenaiCodexProvider())

	var buf bytes.Buffer
	err := cli.RunStatus(dir, minimalConfig(), reg, &buf, false, noColourStyles())

	require.NoError(t, err)
	out := buf.String()
	assert.Contains(t, out, cli.ShortenPath(filepath.Join(dir, ".gemini", "settings.json"), home))
	assert.Contains(t, out, cli.ShortenPath(filepath.Join(dir, ".codex", "config.toml"), home))
}

func TestRunStatus_showsLocalScope_forProjectPath(t *testing.T) {
	dir := t.TempDir()
	reg := registry.NewRegistry()
	// Gemini resolves to project path when projectRoot is set → scope = local.
	reg.Register(providers.NewGoogleGeminiProvider())

	var buf bytes.Buffer
	err := cli.RunStatus(dir, minimalConfig(), reg, &buf, false, noColourStyles())

	require.NoError(t, err)
	assert.Contains(t, buf.String(), "local")
}

func TestRunStatus_showsGlobalScope_whenNoProjectRoot(t *testing.T) {
	reg := registry.NewRegistry()
	// Empty projectRoot → Gemini resolves to its global config path → scope = global.
	reg.Register(providers.NewGoogleGeminiProvider())

	var buf bytes.Buffer
	err := cli.RunStatus("", minimalConfig(), reg, &buf, false, noColourStyles())

	require.NoError(t, err)
	assert.Contains(t, buf.String(), "global")
}

func TestRunStatus_jsonOutput_includesScope(t *testing.T) {
	dir := t.TempDir()
	reg := registry.NewRegistry()
	reg.Register(providers.NewGoogleGeminiProvider())
	var buf bytes.Buffer

	err := cli.RunStatus(dir, minimalConfig(), reg, &buf, true, noColourStyles())

	require.NoError(t, err)
	var out cli.StatusOutput
	require.NoError(t, json.NewDecoder(&buf).Decode(&out))
	// Gemini has both local and global configs — expect one row per scope.
	require.Len(t, out.Providers, 2)
	assert.Equal(t, "local", out.Providers[0].Scope)
	assert.Equal(t, "global", out.Providers[1].Scope)
}

func TestShortenPath_replacesHomePrefix(t *testing.T) {
	tests := []struct {
		path, home, want string
	}{
		{"/home/user/.config/foo", "/home/user", "~/.config/foo"},
		{"/home/user", "/home/user", "~"},
		{"/tmp/foo", "/home/user", "/tmp/foo"},
		{"/home/user/.config/foo", "", "/home/user/.config/foo"},
		{"/home/user/.config/foo", ".", "/home/user/.config/foo"},
	}
	for _, tc := range tests {
		got := cli.ShortenPath(tc.path, tc.home)
		assert.Equal(t, tc.want, got, "path=%q home=%q", tc.path, tc.home)
	}
}

func TestRunStatus_jsonOutput_validJSON(t *testing.T) {
	dir := t.TempDir()
	reg := registry.NewRegistry()
	reg.Register(providers.NewOpenaiCodexProviderWithPath(filepath.Join(dir, "codex.toml")))
	cfg := minimalConfig()
	var buf bytes.Buffer

	err := cli.RunStatus(dir, cfg, reg, &buf, true, noColourStyles())

	require.NoError(t, err)

	var out cli.StatusOutput
	require.NoError(t, json.NewDecoder(&buf).Decode(&out))
	require.Len(t, out.Providers, 1)
	assert.Equal(t, "codex", out.Providers[0].Provider)
	assert.Equal(t, "missing", out.Providers[0].Status)
	assert.NotEmpty(t, out.Providers[0].Path)
}

func TestRunStatus_jsonOutput_synced(t *testing.T) {
	dir := t.TempDir()
	reg := registry.NewRegistry()
	reg.Register(providers.NewClaudeCodeProvider())
	cfg := minimalConfig()

	p, err := reg.Get("claude")
	require.NoError(t, err)
	content, err := p.Generate(cfg.Servers, "")
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(filepath.Join(dir, ".mcp.json"), []byte(content), 0o644))

	var buf bytes.Buffer
	err = cli.RunStatus(dir, cfg, reg, &buf, true, noColourStyles())

	require.NoError(t, err)

	var out cli.StatusOutput
	require.NoError(t, json.NewDecoder(&buf).Decode(&out))
	// Claude has both local and global configs — local (index 0) should be synced.
	require.Len(t, out.Providers, 2)
	assert.Equal(t, "local", out.Providers[0].Scope)
	assert.Equal(t, "synced", out.Providers[0].Status)
	assert.Equal(t, "global", out.Providers[1].Scope)
}
