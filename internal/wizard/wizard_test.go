package wizard_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pantheon-org/iris/internal/config"
	"github.com/pantheon-org/iris/internal/providers"
	"github.com/pantheon-org/iris/internal/wizard"
)

func newStore(t *testing.T) *config.Store {
	t.Helper()
	store, err := config.NewStore(filepath.Join(t.TempDir(), ".iris.json"))
	require.NoError(t, err)
	return store
}

func newRegistry() *providers.Registry {
	return providers.NewRegistry()
}

func TestRunInit_happyPath_twoServers(t *testing.T) {
	store := newStore(t)
	r := wizard.NewScriptedRunner([]string{
		"yes",      // Add a server?
		"server-a", // Server name
		"npx",      // Command
		"-y foo",   // Args
		"stdio",    // Transport
		"yes",      // Add a server?
		"server-b", // Server name
		"uvx",      // Command
		"",         // Args (none)
		"stdio",    // Transport
		"no",       // Add a server?
	})

	err := wizard.RunInit(r, store, newRegistry())
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

	err := wizard.RunInit(r, store, newRegistry())
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
		"npx",      // Command
		"",         // Args (none)
		"stdio",    // Transport
		"no",       // Add a server?
	})

	err := wizard.RunInit(r, store, newRegistry())
	require.NoError(t, err)

	cfg, err := store.Load()
	require.NoError(t, err)
	assert.Len(t, cfg.Servers, 1)
	assert.Contains(t, cfg.Servers, "server-a")
}

func TestRunInit_duplicateName_overwritten(t *testing.T) {
	store := newStore(t)
	r := wizard.NewScriptedRunner([]string{
		"yes",    // Add a server?
		"my-srv", // Server name
		"npx",    // Command
		"",       // Args
		"stdio",  // Transport
		"yes",    // Add a server?
		"my-srv", // Same name
		"uvx",    // Different command
		"",       // Args
		"stdio",  // Transport
		"no",     // Add a server?
	})

	err := wizard.RunInit(r, store, newRegistry())
	require.NoError(t, err)

	cfg, err := store.Load()
	require.NoError(t, err)
	assert.Len(t, cfg.Servers, 1)
	assert.Equal(t, "uvx", cfg.Servers["my-srv"].Command)
}
