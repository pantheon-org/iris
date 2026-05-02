package i18n_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pantheon-org/iris/internal/i18n"
)

func TestT_knownKey_returnsEnglishByDefault(t *testing.T) {
	i18n.SetLang("en")
	assert.Equal(t, "Transport", i18n.T("wizard.transport"))
}

func TestT_unknownKey_returnsKeyItself(t *testing.T) {
	i18n.SetLang("en")
	assert.Equal(t, "no.such.key", i18n.T("no.such.key"))
}

func TestT_withArgs_interpolates(t *testing.T) {
	i18n.SetLang("en")
	assert.Equal(t, "Initialized /tmp/foo", i18n.T("init.initialized", "/tmp/foo"))
}

func TestSetLang_knownCode_switchesLocale(t *testing.T) {
	i18n.SetLang("fr")
	assert.Equal(t, "Transport", i18n.T("wizard.transport"))
	assert.Equal(t, "Commande", i18n.T("wizard.command"))
	i18n.SetLang("en")
}

func TestSetLang_unknownCode_fallsBackToEnglish(t *testing.T) {
	i18n.SetLang("en") // ensure active is English before the test
	i18n.SetLang("xx") // unknown: should leave active unchanged
	assert.Equal(t, "Transport", i18n.T("wizard.transport"))
}

func TestDetectLocale_langEnv_returnsBaseCode(t *testing.T) {
	t.Setenv("LC_ALL", "")
	t.Setenv("LANGUAGE", "")
	t.Setenv("LANG", "fr_FR.UTF-8")
	assert.Equal(t, "fr", i18n.DetectLocale())
}

func TestDetectLocale_ptBR_returnsFullCode(t *testing.T) {
	t.Setenv("LC_ALL", "")
	t.Setenv("LANGUAGE", "")
	t.Setenv("LANG", "pt_BR.UTF-8")
	assert.Equal(t, "pt-BR", i18n.DetectLocale())
}

func TestDetectLocale_zhCN_returnsFullCode(t *testing.T) {
	t.Setenv("LC_ALL", "")
	t.Setenv("LANGUAGE", "")
	t.Setenv("LANG", "zh_CN.UTF-8")
	assert.Equal(t, "zh-CN", i18n.DetectLocale())
}

func TestDetectLocale_lcAllTakesPrecedence(t *testing.T) {
	t.Setenv("LC_ALL", "de_DE.UTF-8")
	t.Setenv("LANG", "en_US.UTF-8")
	t.Setenv("LANGUAGE", "")
	assert.Equal(t, "de", i18n.DetectLocale())
}

func TestDetectLocale_noEnvVars_returnsEn(t *testing.T) {
	t.Setenv("LC_ALL", "")
	t.Setenv("LANG", "")
	t.Setenv("LANGUAGE", "")
	assert.Equal(t, "en", i18n.DetectLocale())
}

func TestDetectLocale_unsupportedLocale_returnsEn(t *testing.T) {
	t.Setenv("LC_ALL", "")
	t.Setenv("LANGUAGE", "")
	t.Setenv("LANG", "xx_XX.UTF-8")
	assert.Equal(t, "en", i18n.DetectLocale())
}

// Finding 2: Load() must propagate errors instead of silently returning empty map.

func TestLoad_missingFile_returnsError(t *testing.T) {
	m, err := i18n.Load("xx-MISSING")
	assert.Error(t, err)
	assert.Nil(t, m)
}

func TestLoad_knownCode_returnsMapAndNoError(t *testing.T) {
	m, err := i18n.Load("en")
	assert.NoError(t, err)
	assert.NotEmpty(t, m)
}

func TestSetLangErr_missingCode_returnsError(t *testing.T) {
	err := i18n.SetLangErr("xx-MISSING")
	assert.Error(t, err)
}

func TestSetLangErr_knownCode_returnsNoError(t *testing.T) {
	err := i18n.SetLangErr("fr")
	assert.NoError(t, err)
	i18n.SetLang("en") // restore
}

// Finding 10: normalize() single-pass optimisation — observable behaviour must be unchanged.

func TestNormalize_exactMatch_preferred(t *testing.T) {
	// "pt-BR" should match exactly, not fall back to base "pt" (which isn't supported).
	t.Setenv("LC_ALL", "")
	t.Setenv("LANGUAGE", "")
	t.Setenv("LANG", "pt_BR.UTF-8")
	assert.Equal(t, "pt-BR", i18n.DetectLocale())
}

func TestNormalize_baseLanguageFallback_matchesWhenNoExact(t *testing.T) {
	// "fr-CA" has no exact match; should fall back to base "fr".
	t.Setenv("LC_ALL", "")
	t.Setenv("LANGUAGE", "")
	t.Setenv("LANG", "fr_CA.UTF-8")
	assert.Equal(t, "fr", i18n.DetectLocale())
}

func TestNormalize_exactMatchBeatsBase(t *testing.T) {
	// "zh-CN" exact match should win over "zh" (unsupported base).
	t.Setenv("LC_ALL", "")
	t.Setenv("LANGUAGE", "")
	t.Setenv("LANG", "zh_CN.UTF-8")
	assert.Equal(t, "zh-CN", i18n.DetectLocale())
}
