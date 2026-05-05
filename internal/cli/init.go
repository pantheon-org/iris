package cli

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/pantheon-org/iris/internal/config"
	"github.com/pantheon-org/iris/internal/i18n"
	"github.com/pantheon-org/iris/internal/types"
)

func RunInitNonInteractive(store *config.Store, w io.Writer, st *Styles) error {
	_, err := store.Load()
	if err == nil {
		fmt.Fprintln(w, st.Warning.Render(i18n.T("init.already_exists", store.Path())))
		return nil
	}
	if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("check config: %w", err)
	}

	cfg := &types.IrisConfig{
		Version: 1,
		Servers: map[string]types.MCPServer{},
	}
	if err := store.Save(cfg); err != nil {
		return fmt.Errorf("save config: %w", err)
	}
	fmt.Fprintln(w, st.Success.Render(i18n.T("init.initialized", store.Path())))
	return nil
}
