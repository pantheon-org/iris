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
		SupportsProjectConfig: true,
		GlobalConfigPath:      path,
	}
	p.resolvedPath = func(projectRoot string) string {
		if projectRoot != "" {
			return filepath.Join(projectRoot, ".qwen", "settings.json")
		}
		return path
	}
	return p
}

// NewQwenProviderWithPath creates a QwenProvider pinned to a fixed config path.
// Intended for use in tests.
func NewQwenProviderWithPath(path string) *QwenProvider {
	p := &QwenProvider{}
	p.config = ProviderConfig{
		Name:                  "qwen",
		DisplayName:           "Qwen Code",
		ConfigPath:            "~/.qwen/settings.json",
		SupportsProjectConfig: true,
		GlobalConfigPath:      path,
	}
	p.resolvedPath = func(_ string) string {
		return path
	}
	return p
}

func qwenConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(".qwen", "settings.json")
	}
	return filepath.Join(home, ".qwen", "settings.json")
}
