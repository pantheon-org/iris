package sync_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

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

	if result.Status != irisync.SyncStatusCreated {
		t.Fatalf("expected SyncStatusCreated, got %q", result.Status)
	}
	if result.Err != nil {
		t.Fatalf("unexpected error: %v", result.Err)
	}
	if _, err := os.Stat(filepath.Join(dir, ".mcp.json")); err != nil {
		t.Fatalf("expected config file to exist: %v", err)
	}
}

func TestSyncProvider_filePresent_contentUnchanged_returnsUnchanged(t *testing.T) {
	dir := t.TempDir()
	p := providers.NewOpenCodeProvider()

	firstResult := irisync.SyncProvider(dir, p, testServers)
	if firstResult.Status != irisync.SyncStatusCreated {
		t.Fatalf("expected SyncStatusCreated on first sync, got %q", firstResult.Status)
	}

	configPath := filepath.Join(dir, "opencode.json")
	statBefore, err := os.Stat(configPath)
	if err != nil {
		t.Fatalf("config file missing: %v", err)
	}

	secondResult := irisync.SyncProvider(dir, p, testServers)

	if secondResult.Status != irisync.SyncStatusUnchanged {
		t.Fatalf("expected SyncStatusUnchanged, got %q", secondResult.Status)
	}
	if secondResult.Err != nil {
		t.Fatalf("unexpected error: %v", secondResult.Err)
	}

	statAfter, err := os.Stat(configPath)
	if err != nil {
		t.Fatalf("config file gone: %v", err)
	}
	if statAfter.ModTime() != statBefore.ModTime() {
		t.Fatal("file was written when content was unchanged")
	}
}

func TestSyncProvider_filePresent_contentChanged_returnsUpdated(t *testing.T) {
	dir := t.TempDir()
	p := providers.NewClaudeCodeProvider()

	firstResult := irisync.SyncProvider(dir, p, testServers)
	if firstResult.Status != irisync.SyncStatusCreated {
		t.Fatalf("expected SyncStatusCreated on first sync, got %q", firstResult.Status)
	}

	updatedServers := map[string]types.MCPServer{
		"test-server": {
			Transport: types.TransportStdio,
			Command:   "uvx",
			Args:      []string{"another-mcp"},
		},
	}

	secondResult := irisync.SyncProvider(dir, p, updatedServers)

	if secondResult.Status != irisync.SyncStatusUpdated {
		t.Fatalf("expected SyncStatusUpdated, got %q", secondResult.Status)
	}
	if secondResult.Err != nil {
		t.Fatalf("unexpected error: %v", secondResult.Err)
	}
}

func TestSyncProvider_generateError_returnsError(t *testing.T) {
	dir := t.TempDir()
	p := providers.NewClaudeCodeProvider()

	badConfigPath := filepath.Join(dir, ".mcp.json")
	if err := os.WriteFile(badConfigPath, []byte("not valid json {{{{"), 0644); err != nil {
		t.Fatalf("failed to write bad config: %v", err)
	}

	result := irisync.SyncProvider(dir, p, testServers)

	if result.Status != irisync.SyncStatusError {
		t.Fatalf("expected SyncStatusError, got %q", result.Status)
	}
	if result.Err == nil {
		t.Fatal("expected non-nil error")
	}
}

func TestSyncAllProviders_multipleProviders_allResultsReturned(t *testing.T) {
	dir := t.TempDir()
	reg := registry.NewRegistry()
	reg.Register(providers.NewClaudeCodeProvider())
	reg.Register(providers.NewOpenCodeProvider())

	results := irisync.SyncAllProviders(dir, reg, testServers)

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestSyncAllProviders_oneErrors_errorCapturedInResult(t *testing.T) {
	dir := t.TempDir()

	badConfigPath := filepath.Join(dir, ".mcp.json")
	if err := os.WriteFile(badConfigPath, []byte("not valid json {{{{"), 0644); err != nil {
		t.Fatalf("failed to write bad config: %v", err)
	}

	reg := registry.NewRegistry()
	reg.Register(providers.NewClaudeCodeProvider())
	reg.Register(providers.NewOpenCodeProvider())

	results := irisync.SyncAllProviders(dir, reg, testServers)

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}

	var errCount, okCount int
	for _, r := range results {
		if r.Status == irisync.SyncStatusError {
			if r.Err == nil {
				t.Errorf("provider %q has error status but nil Err", r.ProviderName)
			}
			errCount++
		} else {
			okCount++
		}
	}
	if errCount != 1 {
		t.Fatalf("expected 1 error result, got %d", errCount)
	}
	if okCount != 1 {
		t.Fatalf("expected 1 ok result, got %d", okCount)
	}
}

func TestSyncAllProviders_oneErrors_doesNotReturnError(t *testing.T) {
	dir := t.TempDir()

	badConfigPath := filepath.Join(dir, ".mcp.json")
	if err := os.WriteFile(badConfigPath, []byte("not valid json {{{{"), 0644); err != nil {
		t.Fatalf("failed to write bad config: %v", err)
	}

	reg := registry.NewRegistry()
	reg.Register(providers.NewClaudeCodeProvider())

	results := irisync.SyncAllProviders(dir, reg, testServers)

	found := false
	for _, r := range results {
		if r.Status == irisync.SyncStatusError && errors.Is(r.Err, r.Err) {
			found = true
		}
	}
	if !found {
		t.Fatal("expected error captured in result, not propagated")
	}
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

	if len(results) != tc.wantLen {
		t.Fatalf("expected %d results, got %d", tc.wantLen, len(results))
	}

	var errCount, okCount int
	for _, r := range results {
		if r.Status == irisync.SyncStatusError {
			if r.Err == nil {
				t.Errorf("provider %q: SyncStatusError but Err is nil", r.ProviderName)
			}
			errCount++
		} else {
			if r.ProviderName == "" {
				t.Errorf("non-error result has empty ProviderName")
			}
			okCount++
		}
	}

	if errCount != tc.wantErrCount {
		t.Errorf("expected %d error result(s), got %d", tc.wantErrCount, errCount)
	}
	if okCount != tc.wantOkCount {
		t.Errorf("expected %d ok result(s), got %d", tc.wantOkCount, okCount)
	}
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
			if err := os.WriteFile(filepath.Join(dir, ".mcp.json"), []byte("{{{invalid}}}"), 0644); err != nil {
				t.Fatalf("setup: %v", err)
			}
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
			if err := os.WriteFile(filepath.Join(dir, ".mcp.json"), []byte("{{{invalid}}}"), 0644); err != nil {
				t.Fatalf("setup: %v", err)
			}
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
			if err := os.WriteFile(filepath.Join(dir, "opencode.json"), []byte("{{{invalid}}}"), 0644); err != nil {
				t.Fatalf("setup: %v", err)
			}
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
			if err := os.WriteFile(filepath.Join(dir, ".mcp.json"), []byte("{{{invalid}}}"), 0644); err != nil {
				t.Fatalf("setup .mcp.json: %v", err)
			}
			if err := os.WriteFile(filepath.Join(dir, "opencode.json"), []byte("{{{invalid}}}"), 0644); err != nil {
				t.Fatalf("setup opencode.json: %v", err)
			}
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

func TestSyncAllProviders_errorsAreContainedInResults_noPanic(t *testing.T) {
	// Verify SyncAllProviders never panics and always returns len(results) == len(providers),
	// even when every provider encounters an error.
	dir := t.TempDir()

	for _, name := range []string{".mcp.json", "opencode.json", ".cursor/mcp.json"} {
		fullPath := filepath.Join(dir, name)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			t.Fatalf("mkdir %s: %v", filepath.Dir(fullPath), err)
		}
		if err := os.WriteFile(fullPath, []byte("{{{invalid}}}"), 0644); err != nil {
			t.Fatalf("write %s: %v", name, err)
		}
	}

	reg := registry.NewRegistry()
	reg.Register(providers.NewClaudeCodeProvider())
	reg.Register(providers.NewOpenCodeProvider())
	reg.Register(providers.NewCursorProvider())

	results := irisync.SyncAllProviders(dir, reg, testServers)

	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	for _, r := range results {
		if r.Status != irisync.SyncStatusError {
			t.Errorf("provider %q: expected SyncStatusError, got %q", r.ProviderName, r.Status)
		}
		if r.Err == nil {
			t.Errorf("provider %q: expected non-nil Err", r.ProviderName)
		}
	}
}
