package wizard_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pantheon-org/iris/internal/config"
	"github.com/pantheon-org/iris/internal/providers"
	"github.com/pantheon-org/iris/internal/registry"
	"github.com/pantheon-org/iris/internal/wizard"
)

func newStore(t *testing.T) *config.Store {
	t.Helper()
	store, err := config.NewStore(filepath.Join(t.TempDir(), ".iris.json"))
	require.NoError(t, err)
	return store
}

func newRegistry() *registry.Registry {
	return registry.NewRegistry()
}

func TestRunInit_happyPath_twoServers(t *testing.T) {
	store := newStore(t)
	r := wizard.NewScriptedRunner([]string{
		// No provider configs exist, so PromptMultiSelect is not called.
		"yes",      // Add a server?
		"server-a", // Server name
		"stdio",    // Transport
		"npx",      // Command
		"-y foo",   // Args
		"yes",      // Add a server?
		"server-b", // Server name
		"stdio",    // Transport
		"uvx",      // Command
		"",         // Args (none)
		"no",       // Add a server?
	})

	err := wizard.RunInit(r, "", store, newRegistry())
	require.NoError(t, err)

	cfg, err := store.Load()
	require.NoError(t, err)
	assert.Len(t, cfg.Servers, 2)
	assert.Contains(t, cfg.Servers, "server-a")
	assert.Contains(t, cfg.Servers, "server-b")
	assert.Equal(t, "npx", cfg.Servers["server-a"].Command)
	assert.Equal(t, []string{"-y", "foo"}, cfg.Servers["server-a"].Args)
	assert.Equal(t, "uvx", cfg.Servers["server-b"].Command)
}

func TestRunInit_noServers_emptyConfig(t *testing.T) {
	store := newStore(t)
	r := wizard.NewScriptedRunner([]string{
		"no", // Add a server? — immediate no
	})

	err := wizard.RunInit(r, "", store, newRegistry())
	require.NoError(t, err)

	cfg, err := store.Load()
	require.NoError(t, err)
	assert.Empty(t, cfg.Servers)
}

func TestRunInit_cancelMidFlow_partialSave(t *testing.T) {
	store := newStore(t)
	r := wizard.NewScriptedRunner([]string{
		"yes",      // Add a server?
		"server-a", // Server name
		"stdio",    // Transport
		"npx",      // Command
		"",         // Args (none)
		"no",       // Add a server?
	})

	err := wizard.RunInit(r, "", store, newRegistry())
	require.NoError(t, err)

	cfg, err := store.Load()
	require.NoError(t, err)
	assert.Len(t, cfg.Servers, 1)
	assert.Contains(t, cfg.Servers, "server-a")
}

func TestRunInit_importDetectedProvider_importsServers(t *testing.T) {
	dir := t.TempDir()
	store, err := config.NewStore(filepath.Join(dir, ".iris.json"))
	require.NoError(t, err)

	mcpJSON := `{"mcpServers":{"imported-srv":{"command":"npx","args":["-y","thing"],"type":"stdio"}}}`
	require.NoError(t, os.WriteFile(filepath.Join(dir, ".mcp.json"), []byte(mcpJSON), 0o600))

	reg := registry.NewRegistry()
	// Pin global path to a non-existent file to avoid reading real ~/.claude.json.
	reg.Register(providers.NewClaudeCodeProviderWithGlobalPath(filepath.Join(dir, "no-global.json")))

	r := wizard.NewScriptedRunner([]string{
		"0",  // PromptMultiSelect: select index 0 (imported-srv [claude-code] [project])
		"no", // Add a server?
	})

	err = wizard.RunInit(r, dir, store, reg)
	require.NoError(t, err)

	cfg, err := store.Load()
	require.NoError(t, err)
	assert.Contains(t, cfg.Servers, "imported-srv")
	assert.Equal(t, "npx", cfg.Servers["imported-srv"].Command)
}

func TestRunInit_importDetectedProvider_declineImport(t *testing.T) {
	dir := t.TempDir()
	store, err := config.NewStore(filepath.Join(dir, ".iris.json"))
	require.NoError(t, err)

	mcpJSON := `{"mcpServers":{"imported-srv":{"command":"npx","args":[],"type":"stdio"}}}`
	require.NoError(t, os.WriteFile(filepath.Join(dir, ".mcp.json"), []byte(mcpJSON), 0o600))

	reg := registry.NewRegistry()
	reg.Register(providers.NewClaudeCodeProviderWithGlobalPath(filepath.Join(dir, "no-global.json")))

	r := wizard.NewScriptedRunner([]string{
		"none", // PromptMultiSelect: select nothing
		"no",   // Add a server?
	})

	err = wizard.RunInit(r, dir, store, reg)
	require.NoError(t, err)

	cfg, err := store.Load()
	require.NoError(t, err)
	assert.Empty(t, cfg.Servers)
}

func TestRunInit_duplicateName_overwritten(t *testing.T) {
	store := newStore(t)
	r := wizard.NewScriptedRunner([]string{
		"yes",    // Add a server?
		"my-srv", // Server name
		"stdio",  // Transport
		"npx",    // Command
		"",       // Args
		"yes",    // Add a server?
		"my-srv", // Same name
		"stdio",  // Transport
		"uvx",    // Different command
		"",       // Args
		"no",     // Add a server?
	})

	err := wizard.RunInit(r, "", store, newRegistry())
	require.NoError(t, err)

	cfg, err := store.Load()
	require.NoError(t, err)
	assert.Len(t, cfg.Servers, 1)
	assert.Equal(t, "uvx", cfg.Servers["my-srv"].Command)
}

func TestRunInit_malformedProviderConfig_skippedAndProceedsNormally(t *testing.T) {
	dir := t.TempDir()
	store, err := config.NewStore(filepath.Join(dir, ".iris.json"))
	require.NoError(t, err)

	// Write a malformed Claude Code config.
	require.NoError(t, os.WriteFile(filepath.Join(dir, ".mcp.json"), []byte(`{bad json}`), 0o600))

	reg := registry.NewRegistry()
	reg.Register(providers.NewClaudeCodeProviderWithGlobalPath(filepath.Join(dir, "no-global.json")))

	// No PromptMultiSelect call expected (malformed config yields zero candidates).
	r := wizard.NewScriptedRunner([]string{
		"no", // Add a server?
	})

	err = wizard.RunInit(r, dir, store, reg)
	require.NoError(t, err, "malformed provider config must not abort init")

	cfg, err := store.Load()
	require.NoError(t, err)
	assert.Empty(t, cfg.Servers)
}

func TestRunInit_unreadableExistingConfig_returnsError(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, ".iris.json")

	// Write a config file that exists but is unreadable (no read permission).
	require.NoError(t, os.WriteFile(cfgPath, []byte(`{"version":1,"servers":{}}`), 0o000))
	t.Cleanup(func() { _ = os.Chmod(cfgPath, 0o600) })

	store, err := config.NewStore(cfgPath)
	require.NoError(t, err)

	r := wizard.NewScriptedRunner([]string{
		"no", // Add a server?
	})

	err = wizard.RunInit(r, dir, store, newRegistry())
	require.Error(t, err)
}

func TestRunInit_sseServer_promptsForURLAndPersistsRemoteConfig(t *testing.T) {
	store := newStore(t)
	r := wizard.NewScriptedRunner([]string{
		"yes",                 // Add a server?
		"remote-srv",          // Server name
		"sse",                 // Transport
		"https://example/mcp", // URL
		"no",                  // Add a server?
	})

	err := wizard.RunInit(r, "", store, newRegistry())
	require.NoError(t, err)

	cfg, err := store.Load()
	require.NoError(t, err)

	srv, ok := cfg.Servers["remote-srv"]
	require.True(t, ok)
	assert.Equal(t, "https://example/mcp", srv.URL)
	assert.Equal(t, "sse", string(srv.Transport))
	assert.Empty(t, srv.Command)
	assert.Empty(t, srv.Args)
}
