package providers

import "path/filepath"

type GeminiProvider struct {
	baseJSONProvider
}

func NewGeminiProvider() *GeminiProvider {
	return newGeminiProviderWithPath(geminiConfigPath())
}

func newGeminiProviderWithPath(path string) *GeminiProvider {
	p := &GeminiProvider{}
	p.config = ProviderConfig{
		Name:                  NameGemini,
		DisplayName:           "Gemini",
		ConfigPath:            "~/.gemini/settings.json",
		SupportsProjectConfig: true,
		GlobalConfigPath:      path,
	}
	p.resolvedPath = func(projectRoot string) string {
		if projectRoot != "" {
			return filepath.Join(projectRoot, ".gemini", "settings.json")
		}
		return path
	}
	return p
}

// NewGeminiProviderWithPath creates a GeminiProvider pinned to a fixed config path.
// Intended for use in tests.
func NewGeminiProviderWithPath(path string) *GeminiProvider {
	p := &GeminiProvider{}
	p.config = ProviderConfig{
		Name:                  NameGemini,
		DisplayName:           "Gemini",
		ConfigPath:            "~/.gemini/settings.json",
		SupportsProjectConfig: true,
		GlobalConfigPath:      path,
	}
	p.resolvedPath = func(_ string) string {
		return path
	}
	return p
}
