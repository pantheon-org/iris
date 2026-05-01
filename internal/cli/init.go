package cli

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/pantheon-org/iris/internal/config"
	"github.com/pantheon-org/iris/internal/types"
)

func RunInitNonInteractive(store *config.Store, w io.Writer) error {
	_, err := store.Load()
	if err == nil {
		fmt.Fprintf(w, "Config already exists at %s\n", store.Path())
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
	fmt.Fprintf(w, "Initialized %s\n", store.Path())
	return nil
}
