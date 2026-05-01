package i18n

import (
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
)

//go:embed locales
var localesFS embed.FS

var supported = []string{
	"de", "en", "es", "fr", "it", "ja", "ko",
	"nl", "pl", "pt-BR", "ru", "tr", "zh-CN", "zh-TW",
}

var (
	mu       sync.RWMutex
	active   map[string]string
	fallback map[string]string
)

func init() {
	fallback = load("en")
	active = fallback
}

func load(code string) map[string]string {
	data, err := localesFS.ReadFile("locales/" + code + ".json")
	if err != nil {
		return map[string]string{}
	}
	var m map[string]string
	if err := json.Unmarshal(data, &m); err != nil {
		return map[string]string{}
	}
	return m
}

// Init detects the OS locale and applies lang if non-empty.
func Init(lang string) {
	if lang == "" {
		lang = DetectLocale()
	}
	SetLang(lang)
}

// SetLang switches the active locale. Unknown codes silently fall back to English.
func SetLang(code string) {
	m := load(code)
	if len(m) == 0 {
		m = fallback
	}
	mu.Lock()
	active = m
	mu.Unlock()
}

// DetectLocale reads LC_ALL, LANG, or LANGUAGE from the environment and maps
// the value to a supported locale code. Returns "en" if nothing matches.
func DetectLocale() string {
	for _, env := range []string{"LC_ALL", "LANG", "LANGUAGE"} {
		if v := os.Getenv(env); v != "" {
			if code := normalize(v); code != "" {
				return code
			}
		}
	}
	return "en"
}

// normalize converts an OS locale string (e.g. "pt_BR.UTF-8") to a supported code.
func normalize(raw string) string {
	if i := strings.IndexByte(raw, '.'); i != -1 {
		raw = raw[:i]
	}
	code := strings.ReplaceAll(raw, "_", "-")

	for _, s := range supported {
		if strings.EqualFold(code, s) {
			return s
		}
	}

	// Fall back to the base language tag (e.g. "en-US" → "en").
	base := strings.SplitN(code, "-", 2)[0]
	for _, s := range supported {
		if strings.EqualFold(base, s) {
			return s
		}
	}
	return ""
}

// T returns the translated string for key. Falls back to English, then the key itself.
// If args are provided the result is passed through fmt.Sprintf.
func T(key string, args ...any) string {
	mu.RLock()
	s, ok := active[key]
	mu.RUnlock()
	if !ok {
		s = fallback[key]
	}
	if s == "" {
		s = key
	}
	if len(args) > 0 {
		return fmt.Sprintf(s, args...)
	}
	return s
}
