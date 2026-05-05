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

func TestRunRemove_existingServer_removedFromDisk(t *testing.T) {
	store, path := newTempStore(t)
	cfg, err := store.Load()
	require.NoError(t, err)

	srv := types.MCPServer{Transport: types.TransportStdio, Command: "uvx"}
	require.NoError(t, cli.RunAdd(cfg, store, "fetch", srv, nil, nil))

	cfg, err = store.Load()
	require.NoError(t, err)

	err = cli.RunRemove(cfg, store, "fetch", nil, nil)
	require.NoError(t, err)

	reloadStore, err := config.NewStore(path)
	require.NoError(t, err)
	reloaded, err := reloadStore.Load()
	require.NoError(t, err)

	_, ok := reloaded.Servers["fetch"]
	assert.False(t, ok, "server 'fetch' should be absent after remove")
}

func TestRunRemove_missingServer_returnsErrServerNotFound(t *testing.T) {
	store, _ := newTempStore(t)
	cfg, err := store.Load()
	require.NoError(t, err)

	err = cli.RunRemove(cfg, store, "nonexistent", nil, nil)
	require.Error(t, err)
	assert.True(t, errors.Is(err, ierrors.ErrServerNotFound))
}

func TestRunRemove_withWriter_printsSuccess(t *testing.T) {
	store, _ := newTempStore(t)
	cfg, err := store.Load()
	require.NoError(t, err)

	srv := types.MCPServer{Transport: types.TransportStdio, Command: "uvx"}
	require.NoError(t, cli.RunAdd(cfg, store, "fetch", srv, nil, nil))

	cfg, err = store.Load()
	require.NoError(t, err)

	var buf bytes.Buffer
	err = cli.RunRemove(cfg, store, "fetch", &buf, noColourStyles())
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "fetch")
}

func TestRunRemove_saveFails_returnsError(t *testing.T) {
	store, path := newTempStore(t)
	cfg, err := store.Load()
	require.NoError(t, err)

	srv := types.MCPServer{Transport: types.TransportStdio, Command: "uvx"}
	require.NoError(t, cli.RunAdd(cfg, store, "fetch", srv, nil, nil))

	cfg, err = store.Load()
	require.NoError(t, err)

	require.NoError(t, os.Chmod(filepath.Dir(path), 0o555))
	t.Cleanup(func() { _ = os.Chmod(filepath.Dir(path), 0o755) })

	err = cli.RunRemove(cfg, store, "fetch", nil, nil)
	require.Error(t, err)
	assert.True(t, errors.Is(err, ierrors.ErrConfigPermission))
}
