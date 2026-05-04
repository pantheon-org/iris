package providers

import (
	"path/filepath"

	"github.com/pantheon-org/iris/internal/io"
)

type ClaudeCodeProvider struct {
	baseJSONProvider
}

func NewClaudeCodeProvider() *ClaudeCodeProvider {
	p := &ClaudeCodeProvider{}
	p.config = ProviderConfig{
		Name:                  NameAnthropicClaudeCode,
		DisplayName:           "Anthropic Claude Code",
		ConfigPath:            ".mcp.json",
		SupportsProjectConfig: true,
		GlobalConfigPath:      io.UserHomePath(".claude.json"),
	}
	p.resolvedPath = func(projectRoot string) string {
		if projectRoot != "" {
			return filepath.Join(projectRoot, ".mcp.json")
		}
		return io.UserHomePath(".claude.json")
	}
	return p
}
