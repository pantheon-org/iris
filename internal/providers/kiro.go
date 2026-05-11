package providers

import (
	"path/filepath"

	"github.com/pantheon-org/iris/internal/io"
	"github.com/pantheon-org/iris/internal/types"
)

type KiroProvider struct {
	baseJSONProvider
}

func NewKiroProvider() *KiroProvider {
	return newKiroProviderWithPath(kiroConfigPath())
}

func newKiroProviderWithPath(path string) *KiroProvider {
	p := &KiroProvider{}
	p.config = ProviderConfig{
		Name:                  types.NameKiro,
		DisplayName:           "Kiro",
		LocalConfigPath:       strPtr(".kiro/settings/mcp.json"),
		SupportsProjectConfig: true,
		GlobalConfigPath:      homeRel(path),
		HasGlobalConfig:       true,
	}
	p.resolvedPath = func(projectRoot string) string {
		return filepath.Join(projectRoot, ".kiro", "settings", "mcp.json")
	}
	return p
}

func NewKiroProviderWithPath(path string) *KiroProvider {
	return newKiroProviderWithPath(path)
}

func kiroConfigPath() string { return io.UserHomePath(".kiro", "settings", "mcp.json") }
