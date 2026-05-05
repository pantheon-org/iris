package wizard

import (
	"fmt"
	"sort"

	"github.com/pantheon-org/iris/internal/i18n"
	irisync "github.com/pantheon-org/iris/internal/sync"
)

// SyncSelection holds the choices made during the interactive sync wizard.
type SyncSelection struct {
	ProviderNames []string // empty = all providers
	Scope         irisync.SyncScope
}

// RunSyncInteractive prompts the user to select which providers and scope to sync,
// then returns their choices for the caller to act on.
func RunSyncInteractive(r Runner, allProviderNames []string) (SyncSelection, error) {
	sorted := make([]string, len(allProviderNames))
	copy(sorted, allProviderNames)
	sort.Strings(sorted)

	// Provider selection.
	selectedIdxs, err := r.PromptMultiSelect(i18n.T("wizard.sync_providers"), sorted)
	if err != nil {
		return SyncSelection{}, fmt.Errorf("prompt providers: %w", err)
	}
	var chosen []string
	for _, idx := range selectedIdxs {
		chosen = append(chosen, sorted[idx])
	}

	// Scope selection.
	scopeOptions := []string{
		i18n.T("wizard.sync_scope.both"),
		i18n.T("wizard.sync_scope.global"),
		i18n.T("wizard.sync_scope.local"),
	}
	scopeChoice, err := r.PromptSelect(i18n.T("wizard.sync_scope"), scopeOptions)
	if err != nil {
		return SyncSelection{}, fmt.Errorf("prompt scope: %w", err)
	}

	var scope irisync.SyncScope
	switch scopeChoice {
	case i18n.T("wizard.sync_scope.global"):
		scope = irisync.ScopeGlobal
	case i18n.T("wizard.sync_scope.local"):
		scope = irisync.ScopeLocal
	default:
		scope = irisync.ScopeAll
	}

	return SyncSelection{ProviderNames: chosen, Scope: scope}, nil
}
