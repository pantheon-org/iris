package config_test

import (
	"errors"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pantheon-org/iris/internal/config"
	"github.com/pantheon-org/iris/internal/ierrors"
	"github.com/pantheon-org/iris/internal/types"
)

func irisFixture() *types.IrisConfig {
	enabled := true
	return &types.IrisConfig{
		Version:   1,
		Providers: []string{"claude"},
		Servers: map[string]types.MCPServer{
			"test-server": {
				Transport: types.TransportStdio,
				Command:   "node",
				Args:      []string{"server.js"},
				Env:       map[string]string{"PORT": "3000"},
				Enabled:   &enabled,
			},
		},
	}
}

func TestNewStore_defaultPath_resolvesToDotIrisJson(t *testing.T) {
	s, err := config.NewStore("")
	require.NoError(t, err)
	assert.Equal(t, config.DefaultConfigFile, s.Path())
}

func TestNewStore_jsonExtension_succeeds(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	s, err := config.NewStore(path)
	require.NoError(t, err)
	assert.Equal(t, path, s.Path())
}

func TestNewStore_yamlExtension_succeeds(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	s, err := config.NewStore(path)
	require.NoError(t, err)
	assert.Equal(t, path, s.Path())
}

func TestNewStore_ymlExtension_succeeds(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yml")
	s, err := config.NewStore(path)
	require.NoError(t, err)
	assert.NotNil(t, s)
}

func TestNewStore_tomlExtension_succeeds(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	s, err := config.NewStore(path)
	require.NoError(t, err)
	assert.NotNil(t, s)
}

func TestNewStore_unknownExtension_returnsMalformedConfigError(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.xml")
	_, err := config.NewStore(path)
	require.Error(t, err)
	assert.True(t, errors.Is(err, ierrors.ErrMalformedConfig))
}

func TestStore_Load_nonExistentFile_returnsError(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "does-not-exist.json")
	s, err := config.NewStore(path)
	require.NoError(t, err)

	_, err = s.Load()
	assert.Error(t, err)
}

func TestStore_Save_createsFileIfAbsent_json(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "new-config.json")
	s, err := config.NewStore(path)
	require.NoError(t, err)

	err = s.Save(irisFixture())
	require.NoError(t, err)

	_, err = os.Stat(path)
	assert.NoError(t, err)
}

func TestStore_SaveLoad_roundTrip_json(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	s, err := config.NewStore(path)
	require.NoError(t, err)

	original := irisFixture()
	require.NoError(t, s.Save(original))

	got, err := s.Load()
	require.NoError(t, err)
	assert.Equal(t, original, got)
}

func TestStore_SaveLoad_roundTrip_yaml(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	s, err := config.NewStore(path)
	require.NoError(t, err)

	original := irisFixture()
	require.NoError(t, s.Save(original))

	got, err := s.Load()
	require.NoError(t, err)
	assert.Equal(t, original, got)
}

func TestStore_SaveLoad_roundTrip_toml(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	s, err := config.NewStore(path)
	require.NoError(t, err)

	original := irisFixture()
	require.NoError(t, s.Save(original))

	got, err := s.Load()
	require.NoError(t, err)
	assert.Equal(t, original, got)
}

func TestStore_Load_noServersKey_serversMapIsNonNil(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "minimal.json")
	require.NoError(t, os.WriteFile(path, []byte(`{"version":1,"providers":[]}`), 0o644))

	s, err := config.NewStore(path)
	require.NoError(t, err)

	cfg, err := s.Load()
	require.NoError(t, err)
	assert.NotNil(t, cfg.Servers, "Servers must be non-nil even when key is absent in config file")
}

func TestStore_Load_malformedContent_returnsMalformedConfigError(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	require.NoError(t, os.WriteFile(path, []byte(`{not valid json`), 0o644))

	s, err := config.NewStore(path)
	require.NoError(t, err)

	_, err = s.Load()
	require.Error(t, err)
	assert.True(t, errors.Is(err, ierrors.ErrMalformedConfig))
}

func TestStore_Load_permissionDenied_returnsConfigPermissionError(t *testing.T) {
	if os.Getuid() == 0 {
		t.Skip("running as root; permission test not meaningful")
	}

	dir := t.TempDir()
	path := filepath.Join(dir, "noperm.json")
	require.NoError(t, os.WriteFile(path, []byte(`{}`), 0o000))

	s, err := config.NewStore(path)
	require.NoError(t, err)

	_, err = s.Load()
	require.Error(t, err)
	assert.True(t, errors.Is(err, ierrors.ErrConfigPermission))
}

// TestStore_Save_concurrent_noDataRace verifies that concurrent calls to Save do not
// race on the underlying file. Run with: go test -race ./internal/config/...
func TestStore_Save_concurrent_noDataRace(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "concurrent.json")

	s, err := config.NewStore(path)
	require.NoError(t, err)

	const goroutines = 20
	var wg sync.WaitGroup
	wg.Add(goroutines)

	errs := make([]error, goroutines)
	for i := range goroutines {
		go func(idx int) {
			defer wg.Done()
			errs[idx] = s.Save(irisFixture())
		}(i)
	}
	wg.Wait()

	for i, e := range errs {
		assert.NoError(t, e, "goroutine %d returned error", i)
	}

	// Final file must be valid and round-trip correctly.
	got, err := s.Load()
	require.NoError(t, err)
	assert.Equal(t, irisFixture(), got)
}
