package providers

import "path/filepath"

type CursorProvider struct {
	baseJSONProvider
}

func NewCursorProvider() *CursorProvider {
	p := &CursorProvider{}
	p.config = ProviderConfig{
		Name:                  NameAnyspereCursor,
		DisplayName:           "Anysphere Cursor",
		ConfigPath:            ".cursor/mcp.json",
		SupportsProjectConfig: true,
	}
	p.resolvedPath = func(projectRoot string) string {
		return filepath.Join(projectRoot, ".cursor", "mcp.json")
	}
	return p
}
