package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"

	"github.com/pantheon-org/iris/internal/registry"
	irisync "github.com/pantheon-org/iris/internal/sync"
	"github.com/pantheon-org/iris/internal/types"
)

// SyncResultEntry is the JSON representation of a single provider result in RunSync output.
type SyncResultEntry struct {
	Provider string `json:"provider"`
	Status   string `json:"status"`
	Path     string `json:"path"`
	Error    string `json:"error,omitempty"`
}

// SyncOutput is the JSON representation of RunSync output.
type SyncOutput struct {
	Results []SyncResultEntry `json:"results"`
}

func RunSync(projectRoot string, cfg *types.IrisConfig, registry *registry.Registry, w io.Writer, jsonOutput bool, st *Styles) error {
	results := irisync.SyncAllProviders(projectRoot, registry, cfg.Servers)

	sort.Slice(results, func(i, j int) bool {
		return results[i].ProviderName < results[j].ProviderName
	})

	var hasErr bool

	if jsonOutput {
		entries := make([]SyncResultEntry, 0, len(results))
		for _, r := range results {
			p, err := registry.Get(r.ProviderName)
			displayPath := r.ProviderName
			if err == nil {
				displayPath = p.ConfigFilePath(projectRoot)
			}

			entry := SyncResultEntry{
				Provider: r.ProviderName,
				Status:   string(r.Status),
				Path:     displayPath,
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
			return fmt.Errorf("sync completed with errors")
		}
		return nil
	}

	t := table.New().
		Border(lipgloss.HiddenBorder()).
		BorderTop(false).
		BorderBottom(false)
	for _, r := range results {
		p, err := registry.Get(r.ProviderName)
		displayPath := r.ProviderName
		if err == nil {
			displayPath = p.ConfigFilePath(projectRoot)
		}

		if r.Err != nil {
			pathCell := st.Muted.Render(displayPath) + "  " + st.Err.Render("("+r.Err.Error()+")")
			t.Row(st.Accent.Render(r.ProviderName), st.Err.Render(string(r.Status)), pathCell)
			hasErr = true
		} else {
			t.Row(st.Accent.Render(r.ProviderName), st.Success.Render(string(r.Status)), st.Muted.Render(displayPath))
		}
	}
	fmt.Fprintln(w, t.Render())

	if hasErr {
		return fmt.Errorf("sync completed with errors")
	}
	return nil
}
