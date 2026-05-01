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
	p := providers.NewCodexProvider()
	cfg := p.Config()

	assert.Equal(t, "codex", cfg.Name)
	assert.Equal(t, "OpenAI Codex", cfg.DisplayName)
	assert.True(t, cfg.SupportsProjectConfig)

	home, err := os.UserHomeDir()
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(home, ".codex", "config.toml"), cfg.GlobalConfigPath)
}

func TestCodexProvider_ConfigFilePath_WithProjectRoot_ReturnsProjectPath(t *testing.T) {
	p := providers.NewCodexProvider()
	got := p.ConfigFilePath("/any/project")
	want := filepath.Join("/any/project", ".codex", "config.toml")
	assert.Equal(t, want, got)
}

func TestCodexProvider_ConfigFilePath_WithEmptyRoot_ReturnsGlobalPath(t *testing.T) {
	p := providers.NewCodexProvider()

	home, err := os.UserHomeDir()
	require.NoError(t, err)
	want := filepath.Join(home, ".codex", "config.toml")
	assert.Equal(t, want, p.ConfigFilePath(""))
}

func TestCodexProvider_Exists_ReturnsFalseWhenFileAbsent(t *testing.T) {
	tmp := t.TempDir()
	p := providers.NewCodexProviderWithPath(filepath.Join(tmp, "config.toml"))
	assert.False(t, p.Exists(""))
}

func TestCodexProvider_Generate_EmptyExistingContent_ProducesCorrectTOML(t *testing.T) {
	p := providers.NewCodexProvider()
	servers := map[string]types.MCPServer{
		"my-server": {
			Transport: types.TransportStdio,
			Command:   "npx",
			Args:      []string{"-y", "@modelcontextprotocol/server-filesystem"},
		},
	}

	got, err := p.Generate(servers, "")
	require.NoError(t, err)

	assert.Contains(t, got, `name = "my-server"`)
	assert.Contains(t, got, `command = "npx"`)
	assert.Contains(t, got, `type = "stdio"`)
	assert.Contains(t, got, `[[mcp_servers]]`)
}

func TestCodexProvider_Generate_PreservesNonMcpServersKeys(t *testing.T) {
	p := providers.NewCodexProvider()

	existing := `version = 1
theme = "dark"

[[mcp_servers]]
  name = "old-server"
  command = "old"
  type = "stdio"
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
	assert.Contains(t, got, `name = "new-server"`)
	assert.NotContains(t, got, `name = "old-server"`)
}

func TestCodexProvider_Generate_WithEnv_IncludesEnvMap(t *testing.T) {
	p := providers.NewCodexProvider()
	servers := map[string]types.MCPServer{
		"env-server": {
			Transport: types.TransportStdio,
			Command:   "node",
			Env:       map[string]string{"FOO": "bar"},
		},
	}

	got, err := p.Generate(servers, "")
	require.NoError(t, err)

	assert.Contains(t, got, "FOO")
	assert.Contains(t, got, "bar")
}

func TestCodexProvider_Parse_ExtractsServersFromFixture(t *testing.T) {
	p := providers.NewCodexProvider()

	content, err := os.ReadFile("testdata/codex_input.toml")
	require.NoError(t, err)

	servers, err := p.Parse(string(content))
	require.NoError(t, err)

	require.Len(t, servers, 2)

	fs, ok := servers["filesystem"]
	require.True(t, ok, "expected key 'filesystem'")
	assert.Equal(t, "npx", fs.Command)
	assert.Equal(t, []string{"-y", "@modelcontextprotocol/server-filesystem"}, fs.Args)
	assert.Equal(t, types.TransportStdio, fs.Transport)

	ev, ok := servers["everything"]
	require.True(t, ok, "expected key 'everything'")
	assert.Equal(t, "DEBUG", func() string {
		for k := range ev.Env {
			return k
		}
		return ""
	}())
}

func TestCodexProvider_Parse_MalformedTOML_WrapsErrMalformedConfig(t *testing.T) {
	p := providers.NewCodexProvider()

	_, err := p.Parse("[[mcp_servers\nname = bad toml ][")
	require.Error(t, err)
	assert.True(t, errors.Is(err, ierrors.ErrMalformedConfig))
}

func TestCodexProvider_Generate_FixtureMatch(t *testing.T) {
	p := providers.NewCodexProvider()

	existing, err := os.ReadFile("testdata/codex_input.toml")
	require.NoError(t, err)

	servers := map[string]types.MCPServer{
		"new-server": {
			Transport: types.TransportStdio,
			Command:   "node",
			Args:      []string{"server.js"},
		},
	}

	got, err := p.Generate(servers, string(existing))
	require.NoError(t, err)

	expected, err := os.ReadFile("testdata/codex_expected.toml")
	require.NoError(t, err)

	assert.Equal(t, string(expected), got)
}
