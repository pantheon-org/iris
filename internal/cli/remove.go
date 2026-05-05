package cli

import (
	"fmt"
	"io"

	"github.com/pantheon-org/iris/internal/config"
	"github.com/pantheon-org/iris/internal/i18n"
	"github.com/pantheon-org/iris/internal/ierrors"
	"github.com/pantheon-org/iris/internal/types"
)

func RunRemove(cfg *types.IrisConfig, store *config.Store, name string, w io.Writer, st *Styles) error {
	if _, ok := cfg.Servers[name]; !ok {
		return fmt.Errorf("remove %s: %w", name, ierrors.ErrServerNotFound)
	}
	delete(cfg.Servers, name)
	if err := store.Save(cfg); err != nil {
		return fmt.Errorf("remove %s: %w", name, err)
	}
	if w != nil && st != nil {
		fmt.Fprintln(w, st.Success.Render(i18n.T("remove.removed", name)))
	}
	return nil
}
