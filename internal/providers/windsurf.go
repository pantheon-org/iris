package providers

import "github.com/pantheon-org/iris/internal/io"

type WindsurfProvider struct {
	baseJSONProvider
}

func NewWindsurfProvider() *WindsurfProvider {
	return newWindsurfProviderWithPath(windsurfConfigPath())
}

func newWindsurfProviderWithPath(path string) *WindsurfProvider {
	p := &WindsurfProvider{}
	p.config = ProviderConfig{
		Name:                  NameWindsurf,
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

func windsurfConfigPath() string { return io.UserHomePath(".codeium", "windsurf", "mcp_config.json") }
