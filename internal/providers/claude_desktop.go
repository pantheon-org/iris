package providers

import (
	"path/filepath"

	"github.com/pantheon-org/iris/internal/io"
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
	return filepath.Join(io.UserConfigDir(), "Claude", "claude_desktop_config.json")
}
