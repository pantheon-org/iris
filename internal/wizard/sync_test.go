package wizard_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	irisync "github.com/pantheon-org/iris/internal/sync"
	"github.com/pantheon-org/iris/internal/wizard"
)

func TestRunSyncInteractive_SelectsProvidersAndScope(t *testing.T) {
	r := wizard.NewScriptedRunner([]string{
		"0,2",                          // PromptMultiSelect: pick index 0 ("claude") and 2 ("opencode") from sorted list
		"Local only (project configs)", // PromptSelect: local scope
	})

	sel, err := wizard.RunSyncInteractive(r, []string{"gemini", "claude", "opencode"})
	require.NoError(t, err)

	assert.Equal(t, []string{"claude", "opencode"}, sel.ProviderNames)
	assert.Equal(t, irisync.ScopeLocal, sel.Scope)
}

func TestRunSyncInteractive_AllProviders_EmptySelection(t *testing.T) {
	r := wizard.NewScriptedRunner([]string{
		"none",                  // PromptMultiSelect: select nothing
		"Both global and local", // PromptSelect: all scope
	})

	sel, err := wizard.RunSyncInteractive(r, []string{"claude", "gemini"})
	require.NoError(t, err)

	assert.Empty(t, sel.ProviderNames)
	assert.Equal(t, irisync.ScopeAll, sel.Scope)
}

func TestRunSyncInteractive_GlobalScope(t *testing.T) {
	r := wizard.NewScriptedRunner([]string{
		"all",                                  // PromptMultiSelect: all providers
		"Global only (home-directory configs)", // PromptSelect: global scope
	})

	sel, err := wizard.RunSyncInteractive(r, []string{"claude", "gemini"})
	require.NoError(t, err)

	assert.Equal(t, []string{"claude", "gemini"}, sel.ProviderNames)
	assert.Equal(t, irisync.ScopeGlobal, sel.Scope)
}

func TestRunSyncInteractive_SortsProvidersAlphabetically(t *testing.T) {
	r := wizard.NewScriptedRunner([]string{
		"0", // first sorted entry = "claude"
		"Both global and local",
	})

	sel, err := wizard.RunSyncInteractive(r, []string{"zed", "claude", "gemini"})
	require.NoError(t, err)

	assert.Equal(t, []string{"claude"}, sel.ProviderNames)
}
