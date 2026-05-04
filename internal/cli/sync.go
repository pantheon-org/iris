package cli

import (
	"fmt"
	"io"
	"sort"

	"github.com/pantheon-org/iris/internal/registry"
	irisync "github.com/pantheon-org/iris/internal/sync"
	"github.com/pantheon-org/iris/internal/types"
)

func RunSync(projectRoot string, cfg *types.IrisConfig, registry *registry.Registry, w io.Writer) error {
	results := irisync.SyncAllProviders(projectRoot, registry, cfg.Servers)

	sort.Slice(results, func(i, j int) bool {
		return results[i].ProviderName < results[j].ProviderName
	})

	var hasErr bool
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
