package cli

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/pantheon-org/iris/internal/i18n"
	"github.com/pantheon-org/iris/internal/types"
)

func RunList(cfg *types.IrisConfig, w io.Writer) error {
	if len(cfg.Servers) == 0 {
		fmt.Fprint(w, "No servers configured.\n")
		return nil
	}

	names := make([]string, 0, len(cfg.Servers))
	for name := range cfg.Servers {
		names = append(names, name)
	}
	sort.Strings(names)

	fmt.Fprintln(w, i18n.T("list.servers_count", len(names)))
	for _, name := range names {
		srv := cfg.Servers[name]
		transport := string(srv.Transport)
		if transport == "" {
			transport = string(types.TransportStdio)
		}

		parts := make([]string, 0, 1+len(srv.Args))
		if srv.Command != "" {
			parts = append(parts, srv.Command)
		}
		parts = append(parts, srv.Args...)
		cmdLine := strings.Join(parts, " ")

		fmt.Fprintf(w, "  %-12s  %-6s  %s\n", name, transport, cmdLine)
	}
	return nil
}
