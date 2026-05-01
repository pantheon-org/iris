package cli

import (
	"fmt"
	"io"
	"sort"

	"github.com/pantheon-org/iris/internal/merger"
	"github.com/pantheon-org/iris/internal/providers"
	"github.com/pantheon-org/iris/internal/types"
)

func RunSync(projectRoot string, cfg *types.IrisConfig, registry *providers.Registry, w io.Writer) error {
	results := merger.SyncAllProviders(projectRoot, registry, cfg.Servers)

	sort.Slice(results, func(i, j int) bool {
		return results[i].ProviderName < results[j].ProviderName
	})

	var hasErr bool
	for _, r := range results {
		p, err := registry.Get(r.ProviderName)
		displayPath := r.ProviderName
		if err == nil {
			displayPath = p.Config().ConfigPath
		}

		if r.Err != nil {
			fmt.Fprintf(w, "  %-12s  %-9s  %s  (%s)\n", r.ProviderName, string(r.Status), displayPath, r.Err)
			hasErr = true
		} else {
			fmt.Fprintf(w, "  %-12s  %-9s  %s\n", r.ProviderName, string(r.Status), displayPath)
		}
	}

	if hasErr {
		return fmt.Errorf("sync completed with errors")
	}
	return nil
}
