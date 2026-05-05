package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"

	"github.com/pantheon-org/iris/internal/i18n"
	"github.com/pantheon-org/iris/internal/registry"
	"github.com/pantheon-org/iris/internal/types"
)

// StatusEntry is the JSON representation of a single provider in RunStatus output.
type StatusEntry struct {
	Provider string `json:"provider"`
	Status   string `json:"status"`
	Path     string `json:"path"`
}

// StatusOutput is the JSON representation of RunStatus output.
type StatusOutput struct {
	Providers []StatusEntry `json:"providers"`
}

func RunStatus(projectRoot string, cfg *types.IrisConfig, registry *registry.Registry, w io.Writer, jsonOutput bool, st *Styles) error {
	all := registry.All()
	sort.Slice(all, func(i, j int) bool {
		return all[i].Config().Name < all[j].Config().Name
	})

	if jsonOutput {
		entries := make([]StatusEntry, 0, len(all))
		for _, p := range all {
			name := p.Config().Name
			path := p.ConfigFilePath(projectRoot)

			data, err := os.ReadFile(path)
			if err != nil {
				status := i18n.T("status.error")
				if errors.Is(err, os.ErrNotExist) {
					status = i18n.T("status.missing")
				}
				entries = append(entries, StatusEntry{Provider: name, Status: status, Path: path})
				continue
			}

			existing := string(data)
			generated, err := p.Generate(cfg.Servers, existing)
			if err != nil {
				entries = append(entries, StatusEntry{Provider: name, Status: i18n.T("status.error"), Path: path})
				continue
			}

			status := i18n.T("status.synced")
			if generated != existing {
				status = i18n.T("status.desync")
			}
			entries = append(entries, StatusEntry{Provider: name, Status: status, Path: path})
		}
		if err := json.NewEncoder(w).Encode(StatusOutput{Providers: entries}); err != nil {
			return fmt.Errorf("encode json: %w", err)
		}
		return nil
	}

	fmt.Fprint(w, st.Bold.Render("Provider Status:")+"\n")
	t := table.New().
		Border(lipgloss.HiddenBorder()).
		BorderTop(false).
		BorderBottom(false)
	for _, p := range all {
		name := p.Config().Name
		path := p.ConfigFilePath(projectRoot)

		data, err := os.ReadFile(path)
		if err != nil {
			statusWord := i18n.T("status.error")
			if errors.Is(err, os.ErrNotExist) {
				statusWord = i18n.T("status.missing")
			}
			t.Row(st.Accent.Render(name), st.Err.Render(statusWord), st.Muted.Render(path))
			continue
		}

		existing := string(data)
		generated, err := p.Generate(cfg.Servers, existing)
		if err != nil {
			t.Row(st.Accent.Render(name), st.Err.Render(i18n.T("status.error")), st.Muted.Render(path))
			continue
		}

		statusWord := i18n.T("status.synced")
		statusStyled := st.Success.Render(statusWord)
		if generated != existing {
			statusWord = i18n.T("status.desync")
			statusStyled = st.Warning.Render(statusWord)
		}
		t.Row(st.Accent.Render(name), statusStyled, st.Muted.Render(path))
	}
	fmt.Fprintln(w, t.Render())
	return nil
}
