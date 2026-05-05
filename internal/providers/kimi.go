package providers

import (
	"github.com/pantheon-org/iris/internal/io"
	"github.com/pantheon-org/iris/internal/types"
)

type KimiProvider struct {
	baseJSONProvider
}

func NewKimiProvider() *KimiProvider {
	return newKimiProviderWithPath(kimiConfigPath())
}

func newKimiProviderWithPath(path string) *KimiProvider {
	p := &KimiProvider{}
	p.config = ProviderConfig{
		Name:                  types.NameMoonshotKimi,
		DisplayName:           "Moonshot Kimi",
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

// NewKimiProviderWithPath creates a KimiProvider using a custom config path.
// Intended for use in tests.
func NewKimiProviderWithPath(path string) *KimiProvider {
	return newKimiProviderWithPath(path)
}

func kimiConfigPath() string { return io.UserHomePath(".kimi", "mcp.json") }
