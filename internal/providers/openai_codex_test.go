package providers_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pantheon-org/iris/internal/ierrors"
	"github.com/pantheon-org/iris/internal/providers"
	"github.com/pantheon-org/iris/internal/types"
)

func TestCodexProvider_Config_ReturnsCorrectProviderConfig(t *testing.T) {
	p := providers.NewOpenaiCodexProvider()
	cfg := p.Config()

	assert.Equal(t, "codex", cfg.Name)
	assert.Equal(t, "OpenAI Codex", cfg.DisplayName)
	assert.True(t, cfg.SupportsProjectConfig)

	home, err := os.UserHomeDir()
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(home, ".codex", "config.toml"), cfg.GlobalConfigPath)
}

func TestCodexProvider_ConfigFilePath_WithProjectRoot_ReturnsProjectPath(t *testing.T) {
	p := providers.NewOpenaiCodexProvider()
	got := p.ConfigFilePath("/any/project")
	want := filepath.Join("/any/project", ".codex", "config.toml")
	assert.Equal(t, want, got)
}

func TestCodexProvider_ConfigFilePath_WithEmptyRoot_ReturnsGlobalPath(t *testing.T) {
	p := providers.NewOpenaiCodexProvider()

	home, err := os.UserHomeDir()
	require.NoError(t, err)
	want := filepath.Join(home, ".codex", "config.toml")
	assert.Equal(t, want, p.ConfigFilePath(""))
}

func TestCodexProvider_Exists_ReturnsFalseWhenFileAbsent(t *testing.T) {
	tmp := t.TempDir()
	p := providers.NewOpenaiCodexProviderWithPath(filepath.Join(tmp, "config.toml"))
	ok, err := p.Exists("")
	require.NoError(t, err)
	assert.False(t, ok)
}

func TestCodexProvider_Generate_EmptyExistingContent_ProducesCorrectTOML(t *testing.T) {
	p := providers.NewOpenaiCodexProvider()
	servers := map[string]types.MCPServer{
		"my-server": {
			Transport: types.TransportStdio,
			Command:   "npx",
			Args:      []string{"-y", "@modelcontextprotocol/server-filesystem"},
		},
	}

	got, err := p.Generate(servers, "")
	require.NoError(t, err)

	assert.Contains(t, got, `[mcp_servers.my-server]`)
	assert.Contains(t, got, `command = "npx"`)
	assert.Contains(t, got, `args = ["-y", "@modelcontextprotocol/server-filesystem"]`)
}

func TestCodexProvider_Generate_PreservesNonMcpServersKeys(t *testing.T) {
	p := providers.NewOpenaiCodexProvider()

	existing := `version = 1
theme = "dark"

[mcp_servers.old-server]
command = "old"
`

	servers := map[string]types.MCPServer{
		"new-server": {
			Transport: types.TransportStdio,
			Command:   "node",
			Args:      []string{"server.js"},
		},
	}

	got, err := p.Generate(servers, existing)
	require.NoError(t, err)

	assert.Contains(t, got, "version = 1")
	assert.Contains(t, got, `theme = "dark"`)
	assert.Contains(t, got, `[mcp_servers.new-server]`)
	assert.NotContains(t, got, `[mcp_servers.old-server]`)
}

func TestCodexProvider_Generate_WithEnv_IncludesEnvMap(t *testing.T) {
	p := providers.NewOpenaiCodexProvider()
	servers := map[string]types.MCPServer{
		"env-server": {
			Transport: types.TransportStdio,
			Command:   "node",
			Env:       map[string]string{"FOO": "bar"},
		},
	}

	got, err := p.Generate(servers, "")
	require.NoError(t, err)

	assert.Contains(t, got, `[mcp_servers.env-server]`)
	assert.Contains(t, got, `[mcp_servers.env-server.env]`)
	assert.Contains(t, got, "FOO")
	assert.Contains(t, got, "bar")
}

func TestCodexProvider_Parse_ExtractsServersFromFixture(t *testing.T) {
	p := providers.NewOpenaiCodexProvider()

	content, err := os.ReadFile("testdata/openai_codex_input.toml")
	require.NoError(t, err)

	servers, err := p.Parse(string(content))
	require.NoError(t, err)

	require.Len(t, servers, 3)

	context7, ok := servers["context7"]
	require.True(t, ok, "expected key 'context7'")
	assert.Equal(t, "npx", context7.Command)
	assert.Equal(t, []string{"-y", "@upstash/context7-mcp"}, context7.Args)
	assert.Equal(t, "MY_ENV_VALUE", context7.Env["MY_ENV_VAR"])
	assert.Equal(t, types.TransportStdio, context7.Transport)

	figma, ok := servers["figma"]
	require.True(t, ok, "expected key 'figma'")
	assert.Equal(t, "https://mcp.figma.com/mcp", figma.URL)
	assert.Equal(t, "us-east-1", figma.Headers["X-Figma-Region"])
	assert.Equal(t, types.TransportSSE, figma.Transport)

	chrome, ok := servers["chrome_devtools"]
	require.True(t, ok, "expected key 'chrome_devtools'")
	require.NotNil(t, chrome.Enabled)
	assert.True(t, *chrome.Enabled)
	assert.Equal(t, "http://localhost:3000/mcp", chrome.URL)
}

func TestCodexProvider_Parse_MalformedTOML_WrapsErrMalformedConfig(t *testing.T) {
	p := providers.NewOpenaiCodexProvider()

	_, err := p.Parse("[[mcp_servers\nname = bad toml ][")
	require.Error(t, err)
	assert.True(t, errors.Is(err, ierrors.ErrMalformedConfig))
}

func TestCodexProvider_Generate_FixtureMatch(t *testing.T) {
	p := providers.NewOpenaiCodexProvider()

	existing, err := os.ReadFile("testdata/openai_codex_input.toml")
	require.NoError(t, err)

	servers, err := p.Parse(string(existing))
	require.NoError(t, err)

	got, err := p.Generate(servers, string(existing))
	require.NoError(t, err)

	expected, err := os.ReadFile("testdata/openai_codex_expected.toml")
	require.NoError(t, err)

	assert.Equal(t, string(expected), got)
}
