package providers

import (
	"os"
	"path/filepath"
)

type WindsurfProvider struct {
	baseJSONProvider
}

func NewWindsurfProvider() *WindsurfProvider {
	return newWindsurfProviderWithPath(windsurfConfigPath())
}

func newWindsurfProviderWithPath(path string) *WindsurfProvider {
	p := &WindsurfProvider{}
	p.config = ProviderConfig{
		Name:                  "windsurf",
		DisplayName:           "Windsurf",
		ConfigPath:            "~/.codeium/windsurf/mcp_config.json",
		SupportsProjectConfig: false,
		GlobalConfigPath:      path,
	}
	p.resolvedPath = func(_ string) string {
		return path
	}
	return p
}

// NewWindsurfProviderWithPath creates a WindsurfProvider using a custom config path.
// Intended for use in tests.
func NewWindsurfProviderWithPath(path string) *WindsurfProvider {
	return newWindsurfProviderWithPath(path)
}

func windsurfConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(".codeium", "windsurf", "mcp_config.json")
	}
	return filepath.Join(home, ".codeium", "windsurf", "mcp_config.json")
}
