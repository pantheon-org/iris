package io

import (
	"os"
	"path/filepath"
)

func UserHomeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
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
