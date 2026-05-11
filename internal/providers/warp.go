package providers

import (
	"github.com/pantheon-org/iris/internal/io"
	"github.com/pantheon-org/iris/internal/types"
)

type WarpProvider struct {
	baseJSONProvider
}

func warpConfigPath() string { return io.UserHomePath(".warp", "mcp.json") }

func NewWarpProvider() *WarpProvider {
	path := warpConfigPath()
	return newWarpProviderWithPath(&path)
}

func newWarpProviderWithPath(path *string) *WarpProvider {
	if path == nil {
		defaultPath := warpConfigPath()
		path = &defaultPath
	}

	p := &WarpProvider{}
	p.config = ProviderConfig{
		Name:                  types.NameWarpTerminal,
		DisplayName:           "Warp Terminal",
		LocalConfigPath:       nil,
		SupportsProjectConfig: false,
		GlobalConfigPath:      homeRel(*path),
		HasGlobalConfig:       true,
	}
	p.resolvedPath = func(_ string) string {
		return *path
	}
	return p
}

// NewWarpProviderWithPath creates a WarpProvider using a custom config path.
// Intended for use in tests.
func NewWarpProviderWithPath(path string) *WarpProvider {
	return newWarpProviderWithPath(&path)
}
