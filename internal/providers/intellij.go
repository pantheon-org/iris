package providers

import (
	"path/filepath"

	"github.com/pantheon-org/iris/internal/types"
)

type IntelliJProvider struct {
	baseJSONProvider
}

func NewIntelliJProvider() *IntelliJProvider {
	p := &IntelliJProvider{}
	p.config = ProviderConfig{
		Name:                  types.NameIntelliJIDEA,
		DisplayName:           "IntelliJ IDEA",
		LocalConfigPath:       strPtr(".idea/mcp.json"),
		SupportsProjectConfig: true,
	}
	p.resolvedPath = func(projectRoot string) string {
		return filepath.Join(projectRoot, ".idea", "mcp.json")
	}
	return p
}
