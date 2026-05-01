package providers

import (
	"os"
	"path/filepath"
)

type QwenProvider struct {
	baseJSONProvider
}

func NewQwenProvider() *QwenProvider {
	return newQwenProviderWithPath(qwenConfigPath())
}

func newQwenProviderWithPath(path string) *QwenProvider {
	p := &QwenProvider{}
	p.config = ProviderConfig{
		Name:                  "qwen",
		DisplayName:           "Qwen Code",
		ConfigPath:            "~/.qwen/settings.json",
		SupportsProjectConfig: false,
		GlobalConfigPath:      path,
	}
	p.resolvedPath = func(_ string) string {
		return path
	}
	return p
}

// NewQwenProviderWithPath creates a QwenProvider using a custom config path.
// Intended for use in tests.
func NewQwenProviderWithPath(path string) *QwenProvider {
	return newQwenProviderWithPath(path)
}

func qwenConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(".qwen", "settings.json")
	}
	return filepath.Join(home, ".qwen", "settings.json")
}
