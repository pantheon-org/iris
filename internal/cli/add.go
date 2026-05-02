package cli

import (
	"fmt"

	"github.com/pantheon-org/iris/internal/config"
	"github.com/pantheon-org/iris/internal/types"
)

func RunAdd(cfg *types.IrisConfig, store *config.Store, name string, server types.MCPServer) error {
	cfg.Servers[name] = server
	if err := store.Save(cfg); err != nil {
		return fmt.Errorf("add %s: %w", name, err)
	}
	return nil
}
