package providers

type WarpProvider struct {
	baseJSONProvider
}

func NewWarpProvider() *WarpProvider {
	return newWarpProviderWithPath(warpConfigPath())
}

func newWarpProviderWithPath(path string) *WarpProvider {
	p := &WarpProvider{}
	p.config = ProviderConfig{
		Name:                  "warp",
		DisplayName:           "Warp",
		ConfigPath:            "~/.warp/mcp.json",
		SupportsProjectConfig: false,
		GlobalConfigPath:      path,
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

func warpConfigPath() string { return homePath(".warp", "mcp.json") }
