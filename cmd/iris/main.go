package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/pantheon-org/iris/internal/version"
)

func main() {
	root := &cobra.Command{
		Use:     "iris",
		Short:   "Manage MCP server configs across AI providers",
		Version: version.Version,
	}

	root.AddCommand(
		&cobra.Command{
			Use:   "init",
			Short: "Scaffold .iris.json in the current project",
			RunE:  func(cmd *cobra.Command, args []string) error { return nil },
		},
		&cobra.Command{
			Use:   "add [name]",
			Short: "Add or update a server entry",
			Args:  cobra.ExactArgs(1),
			RunE:  func(cmd *cobra.Command, args []string) error { return nil },
		},
		&cobra.Command{
			Use:   "remove [name]",
			Short: "Remove a server entry",
			Args:  cobra.ExactArgs(1),
			RunE:  func(cmd *cobra.Command, args []string) error { return nil },
		},
		&cobra.Command{
			Use:   "list",
			Short: "Pretty-print all servers",
			RunE:  func(cmd *cobra.Command, args []string) error { return nil },
		},
		&cobra.Command{
			Use:   "sync",
			Short: "Re-generate all active provider config files",
			RunE:  func(cmd *cobra.Command, args []string) error { return nil },
		},
		&cobra.Command{
			Use:   "status",
			Short: "Show per-provider sync state",
			RunE:  func(cmd *cobra.Command, args []string) error { return nil },
		},
	)

	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
