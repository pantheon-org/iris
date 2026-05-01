package providers

type GeminiProvider struct {
	baseJSONProvider
}

func NewGeminiProvider() *GeminiProvider {
	return newGeminiProviderWithPath(geminiConfigPath())
}

func newGeminiProviderWithPath(path string) *GeminiProvider {
	p := &GeminiProvider{}
	p.config = ProviderConfig{
		Name:                  "gemini",
		DisplayName:           "Gemini",
		ConfigPath:            "~/.config/gemini/settings.json",
		SupportsProjectConfig: false,
		GlobalConfigPath:      path,
	}
	p.resolvedPath = func(_ string) string {
		return path
	}
	return p
}

// NewGeminiProviderWithPath creates a GeminiProvider using a custom config path.
// Intended for use in tests.
func NewGeminiProviderWithPath(path string) *GeminiProvider {
	return newGeminiProviderWithPath(path)
}
