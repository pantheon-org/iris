package merger_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/pantheon-org/iris/internal/merger"
	"github.com/pantheon-org/iris/internal/providers"
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
	p := providers.NewClaudeProvider()

	result := merger.SyncProvider(dir, p, testServers)

	if result.Status != merger.SyncStatusCreated {
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

	firstResult := merger.SyncProvider(dir, p, testServers)
	if firstResult.Status != merger.SyncStatusCreated {
		t.Fatalf("expected SyncStatusCreated on first sync, got %q", firstResult.Status)
	}

	configPath := filepath.Join(dir, "opencode.json")
	statBefore, err := os.Stat(configPath)
	if err != nil {
		t.Fatalf("config file missing: %v", err)
	}

	secondResult := merger.SyncProvider(dir, p, testServers)

	if secondResult.Status != merger.SyncStatusUnchanged {
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
	p := providers.NewClaudeProvider()

	firstResult := merger.SyncProvider(dir, p, testServers)
	if firstResult.Status != merger.SyncStatusCreated {
		t.Fatalf("expected SyncStatusCreated on first sync, got %q", firstResult.Status)
	}

	updatedServers := map[string]types.MCPServer{
		"test-server": {
			Transport: types.TransportStdio,
			Command:   "uvx",
			Args:      []string{"another-mcp"},
		},
	}

	secondResult := merger.SyncProvider(dir, p, updatedServers)

	if secondResult.Status != merger.SyncStatusUpdated {
		t.Fatalf("expected SyncStatusUpdated, got %q", secondResult.Status)
	}
	if secondResult.Err != nil {
		t.Fatalf("unexpected error: %v", secondResult.Err)
	}
}

func TestSyncProvider_generateError_returnsError(t *testing.T) {
	dir := t.TempDir()
	p := providers.NewClaudeProvider()

	badConfigPath := filepath.Join(dir, ".mcp.json")
	if err := os.WriteFile(badConfigPath, []byte("not valid json {{{{"), 0644); err != nil {
		t.Fatalf("failed to write bad config: %v", err)
	}

	result := merger.SyncProvider(dir, p, testServers)

	if result.Status != merger.SyncStatusError {
		t.Fatalf("expected SyncStatusError, got %q", result.Status)
	}
	if result.Err == nil {
		t.Fatal("expected non-nil error")
	}
}

func TestSyncAllProviders_multipleProviders_allResultsReturned(t *testing.T) {
	dir := t.TempDir()
	registry := providers.NewRegistry()
	registry.Register(providers.NewClaudeProvider())
	registry.Register(providers.NewOpenCodeProvider())

	results := merger.SyncAllProviders(dir, registry, testServers)

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

	registry := providers.NewRegistry()
	registry.Register(providers.NewClaudeProvider())
	registry.Register(providers.NewOpenCodeProvider())

	results := merger.SyncAllProviders(dir, registry, testServers)

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}

	var errCount, okCount int
	for _, r := range results {
		if r.Status == merger.SyncStatusError {
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

	registry := providers.NewRegistry()
	registry.Register(providers.NewClaudeProvider())

	results := merger.SyncAllProviders(dir, registry, testServers)

	found := false
	for _, r := range results {
		if r.Status == merger.SyncStatusError && errors.Is(r.Err, r.Err) {
			found = true
		}
	}
	if !found {
		t.Fatal("expected error captured in result, not propagated")
	}
}
