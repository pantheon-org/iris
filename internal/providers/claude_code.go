package providers

import "path/filepath"

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
	}
	p.resolvedPath = func(projectRoot string) string {
		return filepath.Join(projectRoot, ".mcp.json")
	}
	return p
}
