package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sort"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"

	"github.com/pantheon-org/iris/internal/i18n"
	irio "github.com/pantheon-org/iris/internal/io"
	"github.com/pantheon-org/iris/internal/registry"
	irisync "github.com/pantheon-org/iris/internal/sync"
	"github.com/pantheon-org/iris/internal/types"
)

// SyncResultEntry is the JSON representation of a single provider result in RunSync output.
type SyncResultEntry struct {
	Provider string `json:"provider"`
	Scope    string `json:"scope,omitempty"`
	Status   string `json:"status"`
	Path     string `json:"path"`
	Error    string `json:"error,omitempty"`
}

// SyncOutput is the JSON representation of RunSync output.
type SyncOutput struct {
	Results []SyncResultEntry `json:"results"`
}

func RunSync(projectRoot string, cfg *types.IrisConfig, reg *registry.Registry, w io.Writer, scope irisync.SyncScope, jsonOutput bool, st *Styles) error {
	results := irisync.SyncAllProviders(projectRoot, scope, reg, cfg.Servers)

	sort.Slice(results, func(i, j int) bool {
		if results[i].ProviderName != results[j].ProviderName {
			return results[i].ProviderName < results[j].ProviderName
		}
		return results[i].Scope < results[j].Scope
	})

	// Show scope column when results carry scope information.
	showScope := false
	for _, r := range results {
		if r.Scope != "" {
			showScope = true
			break
		}
	}

	var hasErr bool

	home := irio.UserHomeDir()

	if jsonOutput {
		entries := make([]SyncResultEntry, 0, len(results))
		for _, r := range results {
			entry := SyncResultEntry{
				Provider: r.ProviderName,
				Scope:    r.Scope,
				Status:   string(r.Status),
				Path:     r.Path,
			}
			if r.Err != nil {
				entry.Error = r.Err.Error()
				hasErr = true
			}
			entries = append(entries, entry)
		}
		if err := json.NewEncoder(w).Encode(SyncOutput{Results: entries}); err != nil {
			return fmt.Errorf("encode json: %w", err)
		}
		if hasErr {
			return errors.New(i18n.T("sync.completed_with_errors"))
		}
		return nil
	}

	t := table.New().
		Border(lipgloss.HiddenBorder()).
		BorderTop(false).
		BorderBottom(false)
	for _, r := range results {
		displayPath := ShortenPath(r.Path, home)
		if r.Err != nil {
			pathCell := st.Muted.Render(displayPath) + "  " + st.Err.Render("("+r.Err.Error()+")")
			if showScope {
				t.Row(st.Accent.Render(r.ProviderName), scopeStyle(r.Scope, st), st.Err.Render(string(r.Status)), pathCell)
			} else {
				t.Row(st.Accent.Render(r.ProviderName), st.Err.Render(string(r.Status)), pathCell)
			}
			hasErr = true
		} else {
			if showScope {
				t.Row(st.Accent.Render(r.ProviderName), scopeStyle(r.Scope, st), st.Success.Render(string(r.Status)), st.Muted.Render(displayPath))
			} else {
				t.Row(st.Accent.Render(r.ProviderName), st.Success.Render(string(r.Status)), st.Muted.Render(displayPath))
			}
		}
	}
	fmt.Fprintln(w, t.Render())

	if hasErr {
		return errors.New(i18n.T("sync.completed_with_errors"))
	}
	return nil
}
