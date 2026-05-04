package sync_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pantheon-org/iris/internal/ierrors"
	"github.com/pantheon-org/iris/internal/providers"
	"github.com/pantheon-org/iris/internal/registry"
	irisync "github.com/pantheon-org/iris/internal/sync"
	"github.com/pantheon-org/iris/internal/types"
)

var testServers = map[string]types.MCPServer{
	"test-server": {
		Transport: types.TransportStdio,
		Command:   "npx",
		Args:      []string{"-y", "some-mcp-server"},
	},
}

func TestSyncProvider_fileAbsent_createsFile(t *testing.T) {
	dir := t.TempDir()
	p := providers.NewClaudeCodeProvider()

	result := irisync.SyncProvider(dir, p, testServers)

	assert.Equal(t, irisync.SyncStatusCreated, result.Status)
	require.NoError(t, result.Err)
	_, err := os.Stat(filepath.Join(dir, ".mcp.json"))
	assert.NoError(t, err)
}

func TestSyncProvider_filePresent_contentUnchanged_returnsUnchanged(t *testing.T) {
	dir := t.TempDir()
	p := providers.NewOpenCodeProvider()

	firstResult := irisync.SyncProvider(dir, p, testServers)
	require.Equal(t, irisync.SyncStatusCreated, firstResult.Status, "expected SyncStatusCreated on first sync")

	configPath := filepath.Join(dir, "opencode.json")
	statBefore, err := os.Stat(configPath)
	require.NoError(t, err, "config file missing")

	secondResult := irisync.SyncProvider(dir, p, testServers)

	assert.Equal(t, irisync.SyncStatusUnchanged, secondResult.Status)
	require.NoError(t, secondResult.Err)

	statAfter, err := os.Stat(configPath)
	require.NoError(t, err, "config file gone")
	assert.Equal(t, statBefore.ModTime(), statAfter.ModTime(), "file was written when content was unchanged")
}

func TestSyncProvider_filePresent_contentChanged_returnsUpdated(t *testing.T) {
	dir := t.TempDir()
	p := providers.NewClaudeCodeProvider()

	firstResult := irisync.SyncProvider(dir, p, testServers)
	require.Equal(t, irisync.SyncStatusCreated, firstResult.Status, "expected SyncStatusCreated on first sync")

	updatedServers := map[string]types.MCPServer{
		"test-server": {
			Transport: types.TransportStdio,
			Command:   "uvx",
			Args:      []string{"another-mcp"},
		},
	}

	secondResult := irisync.SyncProvider(dir, p, updatedServers)

	assert.Equal(t, irisync.SyncStatusUpdated, secondResult.Status)
	require.NoError(t, secondResult.Err)
}

func TestSyncProvider_generateError_returnsError(t *testing.T) {
	dir := t.TempDir()
	p := providers.NewClaudeCodeProvider()

	badConfigPath := filepath.Join(dir, ".mcp.json")
	require.NoError(t, os.WriteFile(badConfigPath, []byte("not valid json {{{{"), 0644))

	result := irisync.SyncProvider(dir, p, testServers)

	assert.Equal(t, irisync.SyncStatusError, result.Status)
	assert.NotNil(t, result.Err)
}

func TestSyncAllProviders_multipleProviders_allResultsReturned(t *testing.T) {
	dir := t.TempDir()
	reg := registry.NewRegistry()
	reg.Register(providers.NewClaudeCodeProvider())
	reg.Register(providers.NewOpenCodeProvider())

	results := irisync.SyncAllProviders(dir, reg, testServers)

	assert.Len(t, results, 2)
}

func TestSyncAllProviders_oneErrors_errorCapturedInResult(t *testing.T) {
	dir := t.TempDir()

	badConfigPath := filepath.Join(dir, ".mcp.json")
	require.NoError(t, os.WriteFile(badConfigPath, []byte("not valid json {{{{"), 0644))

	reg := registry.NewRegistry()
	reg.Register(providers.NewClaudeCodeProvider())
	reg.Register(providers.NewOpenCodeProvider())

	results := irisync.SyncAllProviders(dir, reg, testServers)

	require.Len(t, results, 2)

	var errCount, okCount int
	for _, r := range results {
		if r.Status == irisync.SyncStatusError {
			assert.NotNil(t, r.Err, "provider %q has error status but nil Err", r.ProviderName)
			errCount++
		} else {
			okCount++
		}
	}
	assert.Equal(t, 1, errCount, "expected 1 error result")
	assert.Equal(t, 1, okCount, "expected 1 ok result")
}

func TestSyncAllProviders_oneErrors_doesNotReturnError(t *testing.T) {
	dir := t.TempDir()

	badConfigPath := filepath.Join(dir, ".mcp.json")
	require.NoError(t, os.WriteFile(badConfigPath, []byte("not valid json {{{{"), 0644))

	reg := registry.NewRegistry()
	reg.Register(providers.NewClaudeCodeProvider())

	results := irisync.SyncAllProviders(dir, reg, testServers)

	found := false
	for _, r := range results {
		if r.Status == irisync.SyncStatusError && r.Err != nil {
			found = true
		}
	}
	assert.True(t, found, "expected error captured in result, not propagated")
}

// Table-driven tests for SyncAllProviders failure modes and edge cases.

type syncAllProvidersTC struct {
	name         string
	setup        func(dir string) // optional pre-test filesystem mutations
	providers    func() []providers.Provider
	wantLen      int
	wantErrCount int
	wantOkCount  int
}

func runSyncAllProvidersTC(t *testing.T, tc syncAllProvidersTC) {
	t.Helper()
	dir := t.TempDir()

	if tc.setup != nil {
		tc.setup(dir)
	}

	reg := registry.NewRegistry()
	for _, p := range tc.providers() {
		reg.Register(p)
	}

	results := irisync.SyncAllProviders(dir, reg, testServers)

	require.Len(t, results, tc.wantLen)

	var errCount, okCount int
	for _, r := range results {
		if r.Status == irisync.SyncStatusError {
			assert.NotNil(t, r.Err, "provider %q: SyncStatusError but Err is nil", r.ProviderName)
			errCount++
		} else {
			assert.NotEmpty(t, r.ProviderName, "non-error result has empty ProviderName")
			okCount++
		}
	}

	assert.Equal(t, tc.wantErrCount, errCount, "error result count")
	assert.Equal(t, tc.wantOkCount, okCount, "ok result count")
}

func TestSyncAllProviders_emptyRegistry_returnsEmptySlice(t *testing.T) {
	runSyncAllProvidersTC(t, syncAllProvidersTC{
		name:         "empty registry returns empty slice",
		providers:    func() []providers.Provider { return nil },
		wantLen:      0,
		wantErrCount: 0,
		wantOkCount:  0,
	})
}

func TestSyncAllProviders_singleProvider_allSucceed(t *testing.T) {
	runSyncAllProvidersTC(t, syncAllProvidersTC{
		name:         "single provider succeeds",
		providers:    func() []providers.Provider { return []providers.Provider{providers.NewClaudeCodeProvider()} },
		wantLen:      1,
		wantErrCount: 0,
		wantOkCount:  1,
	})
}

func TestSyncAllProviders_singleProvider_fails(t *testing.T) {
	runSyncAllProvidersTC(t, syncAllProvidersTC{
		name: "single provider with bad config fails",
		setup: func(dir string) {
			// Claude uses .mcp.json — write corrupt JSON to trigger a Generate error.
			require.NoError(t, os.WriteFile(filepath.Join(dir, ".mcp.json"), []byte("{{{invalid}}}"), 0644))
		},
		providers:    func() []providers.Provider { return []providers.Provider{providers.NewClaudeCodeProvider()} },
		wantLen:      1,
		wantErrCount: 1,
		wantOkCount:  0,
	})
}

func TestSyncAllProviders_multipleProviders_allSucceed(t *testing.T) {
	runSyncAllProvidersTC(t, syncAllProvidersTC{
		name: "three providers all succeed",
		providers: func() []providers.Provider {
			return []providers.Provider{
				providers.NewClaudeCodeProvider(),
				providers.NewOpenCodeProvider(),
				providers.NewCursorProvider(),
			}
		},
		wantLen:      3,
		wantErrCount: 0,
		wantOkCount:  3,
	})
}

func TestSyncAllProviders_oneProviderFails_remainingProvidersRun(t *testing.T) {
	// Claude (.mcp.json) fails; OpenCode (opencode.json) should still succeed.
	runSyncAllProvidersTC(t, syncAllProvidersTC{
		name: "one provider fails, remaining providers still run",
		setup: func(dir string) {
			require.NoError(t, os.WriteFile(filepath.Join(dir, ".mcp.json"), []byte("{{{invalid}}}"), 0644))
		},
		providers: func() []providers.Provider {
			return []providers.Provider{
				providers.NewClaudeCodeProvider(),
				providers.NewOpenCodeProvider(),
			}
		},
		wantLen:      2,
		wantErrCount: 1,
		wantOkCount:  1,
	})
}

func TestSyncAllProviders_lastProviderFails_allResultsPresent(t *testing.T) {
	// OpenCode (opencode.json) is pre-populated with bad JSON; Claude (.mcp.json) succeeds.
	runSyncAllProvidersTC(t, syncAllProvidersTC{
		name: "last provider fails, all results still present",
		setup: func(dir string) {
			require.NoError(t, os.WriteFile(filepath.Join(dir, "opencode.json"), []byte("{{{invalid}}}"), 0644))
		},
		providers: func() []providers.Provider {
			return []providers.Provider{
				providers.NewClaudeCodeProvider(),
				providers.NewOpenCodeProvider(),
			}
		},
		wantLen:      2,
		wantErrCount: 1,
		wantOkCount:  1,
	})
}

func TestSyncAllProviders_allProvidersFail_allResultsHaveErrors(t *testing.T) {
	runSyncAllProvidersTC(t, syncAllProvidersTC{
		name: "all providers fail, all results have error status",
		setup: func(dir string) {
			require.NoError(t, os.WriteFile(filepath.Join(dir, ".mcp.json"), []byte("{{{invalid}}}"), 0644), "setup .mcp.json")
			require.NoError(t, os.WriteFile(filepath.Join(dir, "opencode.json"), []byte("{{{invalid}}}"), 0644), "setup opencode.json")
		},
		providers: func() []providers.Provider {
			return []providers.Provider{
				providers.NewClaudeCodeProvider(),
				providers.NewOpenCodeProvider(),
			}
		},
		wantLen:      2,
		wantErrCount: 2,
		wantOkCount:  0,
	})
}

func TestSyncProvider_symlinkTarget_returnsError(t *testing.T) {
	dir := t.TempDir()
	p := providers.NewClaudeCodeProvider()

	// Create a real file that the symlink will point to.
	realFile := filepath.Join(dir, "real-config.json")
	require.NoError(t, os.WriteFile(realFile, []byte("{}"), 0644), "setup: write real file")

	// Create a symlink at the provider's config path pointing to the real file.
	configPath := p.ConfigFilePath(dir)
	require.NoError(t, os.Symlink(realFile, configPath), "setup: create symlink")

	result := irisync.SyncProvider(dir, p, testServers)

	assert.Equal(t, irisync.SyncStatusError, result.Status)
	assert.True(t, errors.Is(result.Err, ierrors.ErrSymlinkNotAllowed), "expected error wrapping ErrSymlinkNotAllowed, got: %v", result.Err)
}

func TestSyncAllProviders_errorsAreContainedInResults_noPanic(t *testing.T) {
	// Verify SyncAllProviders never panics and always returns len(results) == len(providers),
	// even when every provider encounters an error.
	dir := t.TempDir()

	for _, name := range []string{".mcp.json", "opencode.json", ".cursor/mcp.json"} {
		fullPath := filepath.Join(dir, name)
		require.NoError(t, os.MkdirAll(filepath.Dir(fullPath), 0755), "mkdir %s", filepath.Dir(fullPath))
		require.NoError(t, os.WriteFile(fullPath, []byte("{{{invalid}}}"), 0644), "write %s", name)
	}

	reg := registry.NewRegistry()
	reg.Register(providers.NewClaudeCodeProvider())
	reg.Register(providers.NewOpenCodeProvider())
	reg.Register(providers.NewCursorProvider())

	results := irisync.SyncAllProviders(dir, reg, testServers)

	require.Len(t, results, 3)
	for _, r := range results {
		assert.Equal(t, irisync.SyncStatusError, r.Status, "provider %q: expected SyncStatusError", r.ProviderName)
		assert.NotNil(t, r.Err, "provider %q: expected non-nil Err", r.ProviderName)
	}
}
