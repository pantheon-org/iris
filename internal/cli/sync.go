package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"

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

func RunSync(projectRoot string, cfg *types.IrisConfig, registry *registry.Registry, w io.Writer, jsonOutput bool) error {
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

	maxWidth := 12 // minimum
	for _, r := range results {
		if len(r.ProviderName) > maxWidth {
			maxWidth = len(r.ProviderName)
		}
	}
	fmtStr := fmt.Sprintf("  %%-%ds  %%-9s  %%s\n", maxWidth)
	fmtStrErr := fmt.Sprintf("  %%-%ds  %%-9s  %%s  (%%s)\n", maxWidth)
	for _, r := range results {
		p, err := registry.Get(r.ProviderName)
		displayPath := r.ProviderName
		if err == nil {
			displayPath = p.ConfigFilePath(projectRoot)
		}

		if r.Err != nil {
			fmt.Fprintf(w, fmtStrErr, r.ProviderName, string(r.Status), displayPath, r.Err)
			hasErr = true
		} else {
			fmt.Fprintf(w, fmtStr, r.ProviderName, string(r.Status), displayPath)
		}
	}

	if hasErr {
		return fmt.Errorf("sync completed with errors")
	}
	return nil
}
