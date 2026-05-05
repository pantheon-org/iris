package providers

import "github.com/pantheon-org/iris/internal/io"

type WarpProvider struct {
	baseJSONProvider
}

func NewWarpProvider() *WarpProvider {
	return newWarpProviderWithPath(warpConfigPath())
}

func newWarpProviderWithPath(path string) *WarpProvider {
	p := &WarpProvider{}
	p.config = ProviderConfig{
		Name:                  NameWarpTerminal,
		DisplayName:           "Warp Terminal",
		LocalConfigPath:       nil,
		SupportsProjectConfig: false,
		GlobalConfigPath:      homeRel(path),
		HasGlobalConfig:       true,
	}
	p.resolvedPath = func(_ string) string {
		return path
	}
	return p
}

// NewWarpProviderWithPath creates a WarpProvider using a custom config path.
// Intended for use in tests.
func NewWarpProviderWithPath(path string) *WarpProvider {
	return newWarpProviderWithPath(path)
}

func warpConfigPath() string { return io.UserHomePath(".warp", "mcp.json") }
