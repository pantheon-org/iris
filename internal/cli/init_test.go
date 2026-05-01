package cli_test

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pantheon-org/iris/internal/cli"
	"github.com/pantheon-org/iris/internal/config"
)

func newEmptyStore(t *testing.T) *config.Store {
	t.Helper()
	path := filepath.Join(t.TempDir(), ".iris.json")
	store, err := config.NewStore(path)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	return store
}

func TestRunInitNonInteractive_noExistingFile_createsConfig(t *testing.T) {
	store := newEmptyStore(t)
	var buf bytes.Buffer

	if err := cli.RunInitNonInteractive(store, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "Initialized") {
		t.Errorf("expected 'Initialized' in output, got %q", buf.String())
	}

	cfg, err := store.Load()
	if err != nil {
		t.Fatalf("config not created: %v", err)
	}
	if cfg.Version != 1 {
		t.Errorf("expected Version=1, got %d", cfg.Version)
	}
	if cfg.Servers == nil {
		t.Error("expected non-nil Servers map")
	}
}

func TestRunInitNonInteractive_existingFile_skips(t *testing.T) {
	store := newEmptyStore(t)
	var buf bytes.Buffer

	if err := cli.RunInitNonInteractive(store, &buf); err != nil {
		t.Fatalf("first call failed: %v", err)
	}
	buf.Reset()

	if err := cli.RunInitNonInteractive(store, &buf); err != nil {
		t.Fatalf("second call failed: %v", err)
	}

	if !strings.Contains(buf.String(), "already exists") {
		t.Errorf("expected 'already exists' in output, got %q", buf.String())
	}
}
