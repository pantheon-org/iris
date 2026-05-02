package cli

import (
	"errors"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/pantheon-org/iris/internal/i18n"
	"github.com/pantheon-org/iris/internal/providers"
	"github.com/pantheon-org/iris/internal/types"
)

func RunStatus(projectRoot string, cfg *types.IrisConfig, registry *providers.Registry, w io.Writer) error {
	all := registry.All()
	sort.Slice(all, func(i, j int) bool {
		return all[i].Config().Name < all[j].Config().Name
	})

	fmt.Fprint(w, "Provider Status:\n")
	for _, p := range all {
		name := p.Config().Name
		path := p.ConfigFilePath(projectRoot)
		displayPath := p.Config().ConfigPath

		data, err := os.ReadFile(path)
		if err != nil {
			status := i18n.T("status.error")
			if errors.Is(err, os.ErrNotExist) {
				status = i18n.T("status.missing")
			}
			fmt.Fprintf(w, "  %-12s  %-8s  %s\n", name, status, displayPath)
			continue
		}

		existing := string(data)
		generated, err := p.Generate(cfg.Servers, existing)
		if err != nil {
			fmt.Fprintf(w, "  %-12s  %-8s  %s\n", name, i18n.T("status.error"), displayPath)
			continue
		}

		status := "synced"
		if generated != existing {
			status = "desync"
		}
		fmt.Fprintf(w, "  %-12s  %-8s  %s\n", name, status, displayPath)
	}
	return nil
}
