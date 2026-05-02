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
	"github.com/pantheon-org/iris/internal/providers"
	"github.com/pantheon-org/iris/internal/types"
)

func buildTestRegistry(t *testing.T, tmpDir string) *providers.Registry {
	t.Helper()
	reg := providers.NewRegistry()
	reg.Register(providers.NewClaudeProvider())
	reg.Register(providers.NewGeminiProviderWithPath(filepath.Join(tmpDir, "gemini-settings.json")))
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

	err := cli.RunStatus(dir, cfg, reg, &buf)

	require.NoError(t, err)
	out := buf.String()
	assert.Contains(t, out, "Provider Status:")
	assert.Contains(t, out, "missing")
}

func TestRunStatus_filePresent_synced(t *testing.T) {
	dir := t.TempDir()
	reg := providers.NewRegistry()
	reg.Register(providers.NewClaudeProvider())

	cfg := minimalConfig()

	// Generate the expected content and write it to the file.
	p, err := reg.Get("claude")
	require.NoError(t, err)
	content, err := p.Generate(cfg.Servers, "")
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(filepath.Join(dir, ".mcp.json"), []byte(content), 0o644))

	var buf bytes.Buffer
	err = cli.RunStatus(dir, cfg, reg, &buf)

	require.NoError(t, err)
	assert.Contains(t, buf.String(), "synced")
}

func TestRunStatus_filePresent_desync(t *testing.T) {
	dir := t.TempDir()
	reg := providers.NewRegistry()
	reg.Register(providers.NewClaudeProvider())

	cfg := minimalConfig()

	// Write stale content (empty JSON object) to trigger desync.
	stale, _ := json.Marshal(map[string]interface{}{})
	require.NoError(t, os.WriteFile(filepath.Join(dir, ".mcp.json"), stale, 0o644))

	var buf bytes.Buffer
	err := cli.RunStatus(dir, cfg, reg, &buf)

	require.NoError(t, err)
	assert.Contains(t, buf.String(), "desync")
}

func TestRunStatus_readFailure_showsError(t *testing.T) {
	dir := t.TempDir()
	reg := providers.NewRegistry()
	reg.Register(providers.NewClaudeProvider())

	cfg := minimalConfig()

	require.NoError(t, os.Mkdir(filepath.Join(dir, ".mcp.json"), 0o755))

	var buf bytes.Buffer
	err := cli.RunStatus(dir, cfg, reg, &buf)

	require.NoError(t, err)
	assert.Contains(t, buf.String(), "error")
	assert.NotContains(t, buf.String(), "missing")
}
