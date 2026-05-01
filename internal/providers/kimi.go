package providers

import (
	"os"
	"path/filepath"
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
		Name:                  "kimi",
		DisplayName:           "Kimi Code",
		ConfigPath:            "~/.kimi/mcp.json",
		SupportsProjectConfig: false,
		GlobalConfigPath:      path,
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

func kimiConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(".kimi", "mcp.json")
	}
	return filepath.Join(home, ".kimi", "mcp.json")
}
