package providers

import "path/filepath"

type ClaudeProvider struct {
	baseJSONProvider
}

func NewClaudeProvider() *ClaudeProvider {
	p := &ClaudeProvider{}
	p.config = ProviderConfig{
		Name:                  "claude",
		DisplayName:           "Claude",
		ConfigPath:            ".mcp.json",
		SupportsProjectConfig: true,
	}
	p.resolvedPath = func(projectRoot string) string {
		return filepath.Join(projectRoot, ".mcp.json")
	}
	return p
}
