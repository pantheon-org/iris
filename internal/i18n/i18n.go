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
	m, err := Load("en")
	if err != nil {
		panic(fmt.Sprintf("i18n: failed to load English fallback: %v", err))
	}
	fallback = m
	active = fallback
}

// Load reads and parses the locale file for code. It returns an error if the
// file cannot be read or is not valid JSON, allowing callers to distinguish
// "missing file" from "corrupted file".
func Load(code string) (map[string]string, error) {
	data, err := localesFS.ReadFile("locales/" + code + ".json")
	if err != nil {
		return nil, fmt.Errorf("i18n: read locale %q: %w", code, err)
	}
	var m map[string]string
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("i18n: parse locale %q: %w", code, err)
	}
	return m, nil
}

// Init detects the OS locale and applies lang if non-empty.
func Init(lang string) {
	if lang == "" {
		lang = DetectLocale()
	}
	SetLang(lang)
}

// SetLang switches the active locale. Unknown codes silently keep the current
// active locale unchanged. Use SetLangErr for error-aware callers.
func SetLang(code string) {
	_ = SetLangErr(code)
}

// SetLangErr switches the active locale and returns an error if the locale
// file cannot be loaded. On error the active locale is left unchanged.
func SetLangErr(code string) error {
	m, err := Load(code)
	if err != nil {
		return err
	}
	mu.Lock()
	active = m
	mu.Unlock()
	return nil
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

// normalize converts an OS locale string (e.g. "pt_BR.UTF-8") to a supported
// code. Single-pass: tracks both an exact match and a base-language fallback;
// returns the exact match immediately when found, otherwise the base match.
func normalize(raw string) string {
	if i := strings.IndexByte(raw, '.'); i != -1 {
		raw = raw[:i]
	}
	code := strings.ReplaceAll(raw, "_", "-")
	base := strings.SplitN(code, "-", 2)[0]

	var baseMatch string
	for _, s := range supported {
		if strings.EqualFold(code, s) {
			return s // exact hit — no need to continue
		}
		if baseMatch == "" && strings.EqualFold(base, s) {
			baseMatch = s
		}
	}
	return baseMatch
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
