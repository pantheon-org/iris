package providers

import (
	"github.com/pantheon-org/iris/internal/io"
	"github.com/pantheon-org/iris/internal/types"
)

type ClineProvider struct {
	baseJSONProvider
}

func clineConfigPath() string {
	return io.UserHomePath(
		"Library", "Application Support", "Code", "User", "globalStorage",
		"saoudrizwan.claude-dev", "settings", "cline_mcp_settings.json",
	)
}

func NewClineProvider() *ClineProvider {
	path := clineConfigPath()
	return newClineProviderWithPath(&path)
}

func newClineProviderWithPath(path *string) *ClineProvider {
	if path == nil {
		defaultPath := clineConfigPath()
		path = &defaultPath
	}

	p := &ClineProvider{}
	p.config = ProviderConfig{
		Name:                  types.NameCline,
		DisplayName:           "Cline",
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

func NewClineProviderWithPath(path string) *ClineProvider {
	return newClineProviderWithPath(&path)
}
