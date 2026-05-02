package io_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pantheon-org/iris/internal/io"
)

func TestUserHomePath_UsesUserHomeDir(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("UserHomeDir: %v", err)
	}

	got := io.UserHomePath("test", "config.json")
	want := filepath.Join(home, "test", "config.json")
	if got != want {
		t.Fatalf("UserHomePath() = %q, want %q", got, want)
	}
}

func TestUserConfigPath_UsesUserConfigDir(t *testing.T) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		t.Fatalf("UserConfigDir: %v", err)
	}

	got := io.UserConfigPath("test", "config.json")
	want := filepath.Join(configDir, "test", "config.json")
	if got != want {
		t.Fatalf("UserConfigPath() = %q, want %q", got, want)
	}
}
