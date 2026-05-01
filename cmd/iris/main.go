package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/pantheon-org/iris/internal/cli"
	"github.com/pantheon-org/iris/internal/config"
	"github.com/pantheon-org/iris/internal/providers"
	"github.com/pantheon-org/iris/internal/types"
	"github.com/pantheon-org/iris/internal/version"
	"github.com/pantheon-org/iris/internal/wizard"
)

func loadConfig(configFlag string) (*config.Store, error) {
	store, err := config.NewStore(configFlag)
	if err != nil {
		return nil, fmt.Errorf("init store: %w", err)
	}
	return store, nil
}

func buildRegistry() *providers.Registry {
	reg := providers.NewRegistry()
	reg.Register(providers.NewClaudeProvider())
	reg.Register(providers.NewClaudeDesktopProvider())
	reg.Register(providers.NewGeminiProvider())
	reg.Register(providers.NewOpenCodeProvider())
	reg.Register(providers.NewCodexProvider())
	reg.Register(providers.NewCursorProvider())
	reg.Register(providers.NewWindsurfProvider())
	reg.Register(providers.NewVSCodeCopilotProvider())
	reg.Register(providers.NewZedProvider())
	reg.Register(providers.NewQwenProvider())
	reg.Register(providers.NewWarpProvider())
	reg.Register(providers.NewKimiProvider())
	reg.Register(providers.NewMistralVibeProvider())
	reg.Register(providers.NewIntelliJProvider())
	return reg
}

func main() {
	var configFlag string

	root := &cobra.Command{
		Use:     "iris",
		Short:   "Manage MCP server configs across AI providers",
		Version: version.Version,
	}

	root.PersistentFlags().StringVarP(&configFlag, "config", "C", config.DefaultConfigFile, "path to .iris.json config file")

	root.AddCommand(
		func() *cobra.Command {
			var interactive bool
			var providerNames []string
			cmd := &cobra.Command{
				Use:   "init",
				Short: "Scaffold .iris.json in the current project",
				RunE: func(cmd *cobra.Command, args []string) error {
					store, err := loadConfig(configFlag)
					if err != nil {
						return err
					}
					if interactive {
						reg := buildRegistry()
						if len(providerNames) > 0 {
							reg, err = reg.Filter(providerNames)
							if err != nil {
								return fmt.Errorf("filter providers: %w", err)
							}
						}
						projectRoot := filepath.Dir(store.Path())
						return wizard.RunInit(wizard.NewBubbleteaRunner(), projectRoot, store, reg)
					}
					return cli.RunInitNonInteractive(store, os.Stdout)
				},
			}
			cmd.Flags().BoolVarP(&interactive, "interactive", "I", false, "run interactive wizard")
			cmd.Flags().StringArrayVarP(&providerNames, "provider", "p", nil, "limit to provider(s) by name (repeatable)")
			return cmd
		}(),
		func() *cobra.Command {
			var (
				command   string
				cmdArgs   []string
				envPairs  []string
				transport string
				url       string
			)
			cmd := &cobra.Command{
				Use:   "add [name]",
				Short: "Add or update a server entry",
				Args:  cobra.ExactArgs(1),
				RunE: func(cmd *cobra.Command, args []string) error {
					store, err := loadConfig(configFlag)
					if err != nil {
						return err
					}
					cfg, err := store.Load()
					if err != nil {
						return fmt.Errorf("load config: %w", err)
					}
					envMap := make(map[string]string, len(envPairs))
					for _, pair := range envPairs {
						k, v, _ := strings.Cut(pair, "=")
						envMap[k] = v
					}
					srv := types.MCPServer{
						Transport: types.Transport(transport),
						Command:   command,
						Args:      cmdArgs,
						Env:       envMap,
						URL:       url,
					}
					return cli.RunAdd(cfg, store, args[0], srv)
				},
			}
			cmd.Flags().StringVarP(&command, "command", "c", "", "command to run (required for stdio)")
			cmd.Flags().StringArrayVarP(&cmdArgs, "args", "a", nil, "arguments for the command")
			cmd.Flags().StringArrayVarP(&envPairs, "env", "e", nil, "environment variables in key=value format")
			cmd.Flags().StringVarP(&transport, "transport", "t", string(types.TransportStdio), "transport type (stdio, sse)")
			cmd.Flags().StringVarP(&url, "url", "u", "", "URL for http/sse transport")
			return cmd
		}(),
		&cobra.Command{
			Use:   "remove [name]",
			Short: "Remove a server entry",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				store, err := loadConfig(configFlag)
				if err != nil {
					return err
				}
				cfg, err := store.Load()
				if err != nil {
					return fmt.Errorf("load config: %w", err)
				}
				return cli.RunRemove(cfg, store, args[0])
			},
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
		func() *cobra.Command {
			var providerNames []string
			cmd := &cobra.Command{
				Use:   "sync",
				Short: "Re-generate all active provider config files",
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
					reg := buildRegistry()
					if len(providerNames) > 0 {
						reg, err = reg.Filter(providerNames)
						if err != nil {
							return fmt.Errorf("filter providers: %w", err)
						}
					}
					if err := cli.RunSync(projectRoot, cfg, reg, os.Stdout); err != nil {
						fmt.Fprintln(os.Stderr, err)
						os.Exit(1)
					}
					return nil
				},
			}
			cmd.Flags().StringArrayVarP(&providerNames, "provider", "p", nil, "limit to provider(s) by name (repeatable)")
			return cmd
		}(),
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
				return cli.RunStatus(projectRoot, cfg, buildRegistry(), os.Stdout)
			},
		},
	)

	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
