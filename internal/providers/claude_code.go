package providers

import (
	"path/filepath"

	"github.com/pantheon-org/iris/internal/io"
)

type ClaudeCodeProvider struct {
	baseJSONProvider
}

func NewClaudeCodeProvider() *ClaudeCodeProvider {
	return newClaudeCodeProviderWithGlobalPath(io.UserHomePath(".claude.json"))
}

// NewClaudeCodeProviderWithGlobalPath creates a ClaudeCodeProvider with a fixed global
// config path. Intended for use in tests to avoid reading the real ~/.claude.json.
func NewClaudeCodeProviderWithGlobalPath(globalPath string) *ClaudeCodeProvider {
	return newClaudeCodeProviderWithGlobalPath(globalPath)
}

func newClaudeCodeProviderWithGlobalPath(globalPath string) *ClaudeCodeProvider {
	p := &ClaudeCodeProvider{}
	p.config = ProviderConfig{
		Name:                  NameAnthropicClaudeCode,
		DisplayName:           "Anthropic Claude Code",
		LocalConfigPath:       strPtr(".mcp.json"),
		SupportsProjectConfig: true,
		GlobalConfigPath:      homeRel(globalPath),
	}
	p.resolvedPath = func(projectRoot string) string {
		if projectRoot != "" {
			return filepath.Join(projectRoot, ".mcp.json")
		}
		return globalPath
	}
	return p
}
