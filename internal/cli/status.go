package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"

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

func RunStatus(projectRoot string, cfg *types.IrisConfig, registry *registry.Registry, w io.Writer, jsonOutput bool) error {
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

	fmt.Fprint(w, "Provider Status:\n")
	maxWidth := 12 // minimum
	for _, p := range all {
		if len(p.Config().Name) > maxWidth {
			maxWidth = len(p.Config().Name)
		}
	}
	fmtStr := fmt.Sprintf("  %%-%ds  %%-8s  %%s\n", maxWidth)
	for _, p := range all {
		name := p.Config().Name
		path := p.ConfigFilePath(projectRoot)
		displayPath := path

		data, err := os.ReadFile(path)
		if err != nil {
			status := i18n.T("status.error")
			if errors.Is(err, os.ErrNotExist) {
				status = i18n.T("status.missing")
			}
			fmt.Fprintf(w, fmtStr, name, status, displayPath)
			continue
		}

		existing := string(data)
		generated, err := p.Generate(cfg.Servers, existing)
		if err != nil {
			fmt.Fprintf(w, fmtStr, name, i18n.T("status.error"), displayPath)
			continue
		}

		status := i18n.T("status.synced")
		if generated != existing {
			status = i18n.T("status.desync")
		}
		fmt.Fprintf(w, fmtStr, name, status, displayPath)
	}
	return nil
}
