package providers

import "path/filepath"

type GoogleGeminiProvider struct {
	baseJSONProvider
}

func NewGoogleGeminiProvider() *GoogleGeminiProvider {
	return newGoogleGeminiProviderWithPath(googleGeminiConfigPath())
}

func newGoogleGeminiProviderWithPath(path string) *GoogleGeminiProvider {
	p := &GoogleGeminiProvider{}
	p.config = ProviderConfig{
		Name:                  NameGoogleGemini,
		DisplayName:           "Google Gemini",
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

// NewGoogleGeminiProviderWithPath creates a GoogleGeminiProvider pinned to a fixed config path.
// Intended for use in tests.
func NewGoogleGeminiProviderWithPath(path string) *GoogleGeminiProvider {
	p := &GoogleGeminiProvider{}
	p.config = ProviderConfig{
		Name:                  NameGoogleGemini,
		DisplayName:           "Google Gemini",
		ConfigPath:            "~/.gemini/settings.json",
		SupportsProjectConfig: true,
		GlobalConfigPath:      path,
	}
	p.resolvedPath = func(_ string) string {
		return path
	}
	return p
}
