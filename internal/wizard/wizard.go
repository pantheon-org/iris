package wizard

import (
	"fmt"
	"os"
	"strings"

	"github.com/pantheon-org/iris/internal/cli"
	"github.com/pantheon-org/iris/internal/config"
	"github.com/pantheon-org/iris/internal/detector"
	"github.com/pantheon-org/iris/internal/i18n"
	"github.com/pantheon-org/iris/internal/providers"
	"github.com/pantheon-org/iris/internal/types"
)

type pendingServer struct {
	name   string
	server types.MCPServer
}

func RunInit(r Runner, projectRoot string, store *config.Store, registry *providers.Registry) error {
	var pending []pendingServer

	// Detect installed harnesses and offer to import their existing MCP servers.
	for _, p := range detector.Detect(projectRoot, registry) {
		importIt, err := r.PromptConfirm(i18n.T("wizard.import_prompt", p.Config().DisplayName))
		if err != nil {
			return fmt.Errorf("prompt import %s: %w", p.Config().Name, err)
		}
		if !importIt {
			continue
		}
		filePath := p.ConfigFilePath(projectRoot)
		content, err := readFile(filePath)
		if err != nil {
			return fmt.Errorf("read %s config: %w", p.Config().Name, err)
		}
		servers, err := p.Parse(content)
		if err != nil {
			return fmt.Errorf("parse %s config: %w", p.Config().Name, err)
		}
		for name, srv := range servers {
			pending = append(pending, pendingServer{name: name, server: srv})
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

	return nil
}

func readFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read file %s: %w", path, err)
	}
	return string(data), nil
}
