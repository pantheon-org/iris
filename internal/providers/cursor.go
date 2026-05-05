package providers

import (
	"path/filepath"

	"github.com/pantheon-org/iris/internal/types"
)

type CursorProvider struct {
	baseJSONProvider
}

func NewCursorProvider() *CursorProvider {
	p := &CursorProvider{}
	p.config = ProviderConfig{
		Name:                  types.NameAnysphereCursor,
		DisplayName:           "Anysphere Cursor",
		LocalConfigPath:       strPtr(".cursor/mcp.json"),
		SupportsProjectConfig: true,
	}
	p.resolvedPath = func(projectRoot string) string {
		return filepath.Join(projectRoot, ".cursor", "mcp.json")
	}
	return p
}
