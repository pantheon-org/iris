package wizard

import (
	"fmt"
	"strings"

	"github.com/pantheon-org/iris/internal/cli"
	"github.com/pantheon-org/iris/internal/config"
	"github.com/pantheon-org/iris/internal/providers"
	"github.com/pantheon-org/iris/internal/types"
)

type pendingServer struct {
	name   string
	server types.MCPServer
}

func RunInit(r Runner, store *config.Store, registry *providers.Registry) error {
	var pending []pendingServer

	for {
		more, err := r.PromptConfirm("Add a server? (yes to continue, no to finish)")
		if err != nil {
			return fmt.Errorf("prompt confirm: %w", err)
		}
		if !more {
			break
		}

		name, err := r.PromptText("Server name", "")
		if err != nil {
			return fmt.Errorf("prompt server name: %w", err)
		}

		command, err := r.PromptText("Command", "")
		if err != nil {
			return fmt.Errorf("prompt command: %w", err)
		}

		argsRaw, err := r.PromptText("Args (space-separated, leave blank for none)", "")
		if err != nil {
			return fmt.Errorf("prompt args: %w", err)
		}
		var args []string
		for _, a := range strings.Fields(argsRaw) {
			if a != "" {
				args = append(args, a)
			}
		}

		transport, err := r.PromptSelect("Transport", []string{"stdio", "sse"})
		if err != nil {
			return fmt.Errorf("prompt transport: %w", err)
		}

		pending = append(pending, pendingServer{
			name: name,
			server: types.MCPServer{
				Transport: types.Transport(transport),
				Command:   command,
				Args:      args,
			},
		})
	}

	cfg, err := store.Load()
	if err != nil || cfg == nil {
		cfg = &types.IrisConfig{Version: 1}
	}

	for _, p := range pending {
		if err := cli.RunAdd(cfg, store, p.name, p.server); err != nil {
			return fmt.Errorf("add server %q: %w", p.name, err)
		}
		if loaded, loadErr := store.Load(); loadErr == nil {
			cfg = loaded
		}
	}

	if len(pending) == 0 {
		if cfg.Servers == nil {
			cfg.Servers = make(map[string]types.MCPServer)
		}
		if err := store.Save(cfg); err != nil {
			return fmt.Errorf("save config: %w", err)
		}
	}

	_ = registry
	return nil
}
