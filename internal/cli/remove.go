package cli

import (
	"fmt"

	"github.com/pantheon-org/iris/internal/config"
	"github.com/pantheon-org/iris/internal/ierrors"
	"github.com/pantheon-org/iris/internal/types"
)

func RunRemove(cfg *types.IrisConfig, store *config.Store, name string) error {
	if _, ok := cfg.Servers[name]; !ok {
		return fmt.Errorf("remove %s: %w", name, ierrors.ErrServerNotFound)
	}
	delete(cfg.Servers, name)
	if err := store.Save(cfg); err != nil {
		return fmt.Errorf("remove %s: %w", name, err)
	}
	return nil
}
