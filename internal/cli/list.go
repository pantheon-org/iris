package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"

	"github.com/pantheon-org/iris/internal/i18n"
	"github.com/pantheon-org/iris/internal/types"
)

// ListServerEntry is the JSON representation of a single server in RunList output.
type ListServerEntry struct {
	Name      string            `json:"name"`
	Transport string            `json:"transport"`
	Command   string            `json:"command,omitempty"`
	Args      []string          `json:"args,omitempty"`
	URL       string            `json:"url,omitempty"`
	Env       map[string]string `json:"env,omitempty"`
}

// ListOutput is the JSON representation of RunList output.
type ListOutput struct {
	Servers []ListServerEntry `json:"servers"`
}

func RunList(cfg *types.IrisConfig, w io.Writer, jsonOutput bool, st *Styles) error {
	names := make([]string, 0, len(cfg.Servers))
	for name := range cfg.Servers {
		names = append(names, name)
	}
	sort.Strings(names)

	if jsonOutput {
		entries := make([]ListServerEntry, 0, len(names))
		for _, name := range names {
			srv := cfg.Servers[name]
			transport := string(srv.Transport)
			if transport == "" {
				transport = string(types.TransportStdio)
			}
			entries = append(entries, ListServerEntry{
				Name:      name,
				Transport: transport,
				Command:   srv.Command,
				Args:      srv.Args,
				URL:       srv.URL,
				Env:       srv.Env,
			})
		}
		if err := json.NewEncoder(w).Encode(ListOutput{Servers: entries}); err != nil {
			return fmt.Errorf("encode json: %w", err)
		}
		return nil
	}

	if len(cfg.Servers) == 0 {
		fmt.Fprint(w, "No servers configured.\n")
		return nil
	}

	fmt.Fprintln(w, st.Bold.Render(i18n.T("list.servers_count", len(names))))
	t := table.New().
		Border(lipgloss.HiddenBorder()).
		BorderTop(false).
		BorderBottom(false)
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

		t.Row(
			st.Accent.Render(name),
			st.Muted.Render(transport),
			st.Muted.Render(cmdLine),
		)
	}
	fmt.Fprintln(w, t.Render())
	return nil
}
