package providers

import (
	"os"
	"path/filepath"
	"runtime"
)

type ClaudeDesktopProvider struct {
	baseJSONProvider
}

func NewClaudeDesktopProvider() *ClaudeDesktopProvider {
	return newClaudeDesktopProviderWithPath(claudeDesktopConfigPath())
}

func newClaudeDesktopProviderWithPath(path string) *ClaudeDesktopProvider {
	p := &ClaudeDesktopProvider{}
	p.config = ProviderConfig{
		Name:                  "claude-desktop",
		DisplayName:           "Claude Desktop",
		ConfigPath:            "~/Library/Application Support/Claude/claude_desktop_config.json",
		SupportsProjectConfig: false,
		GlobalConfigPath:      path,
	}
	p.resolvedPath = func(_ string) string {
		return path
	}
	return p
}

// NewClaudeDesktopProviderWithPath creates a ClaudeDesktopProvider using a custom config path.
// Intended for use in tests.
func NewClaudeDesktopProviderWithPath(path string) *ClaudeDesktopProvider {
	return newClaudeDesktopProviderWithPath(path)
}

func claudeDesktopConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	switch runtime.GOOS {
	case "windows":
		appData := os.Getenv("APPDATA")
		if appData == "" {
			appData = filepath.Join(home, "AppData", "Roaming")
		}
		return filepath.Join(appData, "Claude", "claude_desktop_config.json")
	case "darwin":
		return filepath.Join(home, "Library", "Application Support", "Claude", "claude_desktop_config.json")
	default:
		return filepath.Join(home, ".config", "Claude", "claude_desktop_config.json")
	}
}
