package wizard

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/pantheon-org/iris/internal/config"
	"github.com/pantheon-org/iris/internal/i18n"
	"github.com/pantheon-org/iris/internal/registry"
	"github.com/pantheon-org/iris/internal/types"
)

type pendingServer struct {
	name   string
	server types.MCPServer
}

func RunInit(r Runner, projectRoot string, store *config.Store, registry *registry.Registry) error {
	var pending []pendingServer

	// Collect all MCP servers found across every provider config (project + global).
	candidates, err := CollectImportCandidates(projectRoot, registry)
	if err != nil {
		return fmt.Errorf("collect import candidates: %w", err)
	}
	grouped := GroupImportCandidates(candidates)
	if len(grouped) > 0 {
		labels := make([]string, len(grouped))
		for i, g := range grouped {
			labels[i] = g.Label()
		}
		selected, err := r.PromptMultiSelect(i18n.T("wizard.select_import"), labels)
		if err != nil {
			return fmt.Errorf("prompt select import: %w", err)
		}
		for _, idx := range selected {
			g := grouped[idx]
			pending = append(pending, pendingServer{name: g.ServerName, server: g.Server})
		}
	}

	for {
		more, err := r.PromptConfirm(i18n.T("wizard.add_server"))
		if err != nil {
			return fmt.Errorf("prompt confirm: %w", err)
		}
		if !more {
			break
		}

		name, err := r.PromptText(i18n.T("wizard.server_name"), "")
		if err != nil {
			return fmt.Errorf("prompt server name: %w", err)
		}

		transport, err := r.PromptSelect(i18n.T("wizard.transport"), []string{"stdio", "sse"})
		if err != nil {
			return fmt.Errorf("prompt transport: %w", err)
		}

		server := types.MCPServer{Transport: types.Transport(transport)}

		switch server.Transport {
		case types.TransportSSE:
			url, err := r.PromptText(i18n.T("wizard.url"), "")
			if err != nil {
				return fmt.Errorf("prompt url: %w", err)
			}
			server.URL = url
		default:
			command, err := r.PromptText(i18n.T("wizard.command"), "")
			if err != nil {
				return fmt.Errorf("prompt command: %w", err)
			}

			argsRaw, err := r.PromptText(i18n.T("wizard.args"), "")
			if err != nil {
				return fmt.Errorf("prompt args: %w", err)
			}
			var args []string
			for _, a := range strings.Fields(argsRaw) {
				if a != "" {
					args = append(args, a)
				}
			}

			server.Command = command
			server.Args = args
		}

		pending = append(pending, pendingServer{
			name:   name,
			server: server,
		})
	}

	cfg, err := store.Load()
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("load config: %w", err)
		}
		cfg = types.NewIrisConfig()
	}

	for _, p := range pending {
		if err := p.server.Validate(); err != nil {
			return fmt.Errorf("add server %q: %w", p.name, err)
		}
		cfg.Servers[p.name] = p.server
	}
	if err := store.Save(cfg); err != nil {
		return fmt.Errorf("save config: %w", err)
	}

	return nil
}
