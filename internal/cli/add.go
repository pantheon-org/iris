package cli

import (
	"fmt"
	"io"

	"github.com/pantheon-org/iris/internal/config"
	"github.com/pantheon-org/iris/internal/i18n"
	"github.com/pantheon-org/iris/internal/types"
)

func RunAdd(cfg *types.IrisConfig, store *config.Store, name string, server types.MCPServer, w io.Writer, st *Styles) error {
	if err := server.Validate(); err != nil {
		return fmt.Errorf("add %s: %w", name, err)
	}
	cfg.Servers[name] = server
	if err := store.Save(cfg); err != nil {
		return fmt.Errorf("add %s: %w", name, err)
	}
	if w != nil && st != nil {
		fmt.Fprintln(w, st.Success.Render(i18n.T("add.added", name)))
	}
	return nil
}
