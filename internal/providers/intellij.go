package providers

import "path/filepath"

type IntelliJProvider struct {
	baseJSONProvider
}

func NewIntelliJProvider() *IntelliJProvider {
	p := &IntelliJProvider{}
	p.config = ProviderConfig{
		Name:                  NameIntelliJIDEA,
		DisplayName:           "IntelliJ IDEA",
		LocalConfigPath:       strPtr(".idea/mcp.json"),
		SupportsProjectConfig: true,
	}
	p.resolvedPath = func(projectRoot string) string {
		return filepath.Join(projectRoot, ".idea", "mcp.json")
	}
	return p
}
