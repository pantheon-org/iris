package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/pantheon-org/iris/internal/cli"
	"github.com/pantheon-org/iris/internal/config"
	"github.com/pantheon-org/iris/internal/providers"
	"github.com/pantheon-org/iris/internal/version"
)

func loadConfig(configFlag string) (*config.Store, error) {
	store, err := config.NewStore(configFlag)
	if err != nil {
		return nil, fmt.Errorf("init store: %w", err)
	}
	return store, nil
}

func main() {
	var configFlag string

	root := &cobra.Command{
		Use:     "iris",
		Short:   "Manage MCP server configs across AI providers",
		Version: version.Version,
	}

	root.PersistentFlags().StringVar(&configFlag, "config", config.DefaultConfigFile, "path to .iris.json config file")

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
			RunE: func(cmd *cobra.Command, args []string) error {
				store, err := loadConfig(configFlag)
				if err != nil {
					return err
				}
				cfg, err := store.Load()
				if err != nil {
					return fmt.Errorf("load config: %w", err)
				}
				return cli.RunList(cfg, os.Stdout)
			},
		},
		&cobra.Command{
			Use:   "sync",
			Short: "Re-generate all active provider config files",
			RunE:  func(cmd *cobra.Command, args []string) error { return nil },
		},
		&cobra.Command{
			Use:   "status",
			Short: "Show per-provider sync state",
			RunE: func(cmd *cobra.Command, args []string) error {
				projectRoot, err := filepath.Abs(".")
				if err != nil {
					return fmt.Errorf("resolve project root: %w", err)
				}
				store, err := loadConfig(configFlag)
				if err != nil {
					return err
				}
				cfg, err := store.Load()
				if err != nil {
					return fmt.Errorf("load config: %w", err)
				}
				reg := providers.NewRegistry()
				reg.Register(providers.NewClaudeProvider())
				reg.Register(providers.NewGeminiProvider())
				reg.Register(providers.NewOpenCodeProvider())
				reg.Register(providers.NewCodexProvider())
				return cli.RunStatus(projectRoot, cfg, reg, os.Stdout)
			},
		},
	)

	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
