package io

import (
	"fmt"
	"os"
	"path/filepath"
)

func UserHomeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: could not determine home directory: %v; using current directory\n", err)
		return "."
	}
	return home
}

func UserHomePath(elem ...string) string {
	return filepath.Join(append([]string{UserHomeDir()}, elem...)...)
}

func UserConfigDir() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return UserHomePath(".config")
	}
	return configDir
}

func UserConfigPath(elem ...string) string {
	return filepath.Join(append([]string{UserConfigDir()}, elem...)...)
}
