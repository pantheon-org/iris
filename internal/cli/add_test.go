package cli_test

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pantheon-org/iris/internal/cli"
	"github.com/pantheon-org/iris/internal/config"
	"github.com/pantheon-org/iris/internal/ierrors"
	"github.com/pantheon-org/iris/internal/types"
)

func newTempStore(t *testing.T) (*config.Store, string) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".iris.json")
	initial := &types.IrisConfig{
		Version: 1,
		Servers: map[string]types.MCPServer{},
	}
	store, err := config.NewStore(path)
	require.NoError(t, err)
	require.NoError(t, store.Save(initial))
	return store, path
}

func TestRunAdd_newServer_savedToDisk(t *testing.T) {
	store, path := newTempStore(t)
	cfg, err := store.Load()
	require.NoError(t, err)

	srv := types.MCPServer{Transport: types.TransportStdio, Command: "uvx", Args: []string{"mcp-server-fetch"}}
	err = cli.RunAdd(cfg, store, "fetch", srv, nil, nil)
	require.NoError(t, err)

	reloadStore, err := config.NewStore(path)
	require.NoError(t, err)
	reloaded, err := reloadStore.Load()
	require.NoError(t, err)

	got, ok := reloaded.Servers["fetch"]
	require.True(t, ok, "server 'fetch' should be present after add")
	assert.Equal(t, "uvx", got.Command)
	assert.Equal(t, []string{"mcp-server-fetch"}, got.Args)
}

func TestRunAdd_existingServer_overwritten(t *testing.T) {
	store, path := newTempStore(t)
	cfg, err := store.Load()
	require.NoError(t, err)

	first := types.MCPServer{Transport: types.TransportStdio, Command: "old-cmd"}
	require.NoError(t, cli.RunAdd(cfg, store, "myserver", first, nil, nil))

	cfg, err = store.Load()
	require.NoError(t, err)

	second := types.MCPServer{Transport: types.TransportStdio, Command: "new-cmd", Args: []string{"--flag"}}
	require.NoError(t, cli.RunAdd(cfg, store, "myserver", second, nil, nil))

	reloadStore, err := config.NewStore(path)
	require.NoError(t, err)
	reloaded, err := reloadStore.Load()
	require.NoError(t, err)

	got := reloaded.Servers["myserver"]
	assert.Equal(t, "new-cmd", got.Command)
	assert.Equal(t, []string{"--flag"}, got.Args)
}

func TestRunAdd_withWriter_printsSuccess(t *testing.T) {
	store, _ := newTempStore(t)
	cfg, err := store.Load()
	require.NoError(t, err)

	var buf bytes.Buffer
	srv := types.MCPServer{Transport: types.TransportStdio, Command: "uvx"}
	err = cli.RunAdd(cfg, store, "fetch", srv, &buf, noColourStyles())
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "fetch")
}

func TestRunAdd_saveFails_returnsError(t *testing.T) {
	store, path := newTempStore(t)
	cfg, err := store.Load()
	require.NoError(t, err)

	require.NoError(t, os.Chmod(filepath.Dir(path), 0o555))
	t.Cleanup(func() { _ = os.Chmod(filepath.Dir(path), 0o755) })

	srv := types.MCPServer{Transport: types.TransportStdio, Command: "uvx"}
	err = cli.RunAdd(cfg, store, "fetch", srv, nil, nil)
	require.Error(t, err)
	assert.True(t, errors.Is(err, ierrors.ErrConfigPermission))
}
