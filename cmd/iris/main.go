package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/pantheon-org/iris/internal/cli"
	"github.com/pantheon-org/iris/internal/config"
	"github.com/pantheon-org/iris/internal/i18n"
	"github.com/pantheon-org/iris/internal/ierrors"
	"github.com/pantheon-org/iris/internal/providers"
	"github.com/pantheon-org/iris/internal/registry"
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

func buildRegistry() *registry.Registry {
	reg := registry.NewRegistry()
	reg.Register(providers.NewClaudeCodeProvider())
	reg.Register(providers.NewClaudeDesktopProvider())
	reg.Register(providers.NewGoogleGeminiProvider())
	reg.Register(providers.NewOpenCodeProvider())
	reg.Register(providers.NewOpenaiCodexProvider())
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

// parseLangArg scans os.Args for --lang value or --lang=value before cobra runs,
// so that command Short descriptions are already translated when cobra builds help.
func parseLangArg(args []string) string {
	for i, a := range args {
		if a == "--lang" && i+1 < len(args) {
			return args[i+1]
		}
		if strings.HasPrefix(a, "--lang=") {
			return strings.TrimPrefix(a, "--lang=")
		}
	}
	return ""
}

func main() {
	i18n.Init(parseLangArg(os.Args[1:]))

	var configFlag string

	root := &cobra.Command{
		Use:     "iris",
		Short:   i18n.T("cmd.root"),
		Version: version.Version,
		// Apply lang from .iris.json when --lang was not given on the CLI.
		// PersistentPreRunE runs before every subcommand, so all output is translated.
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if parseLangArg(os.Args[1:]) != "" {
				return nil // --lang already applied at startup
			}
			store, err := loadConfig(configFlag)
			if err != nil {
				return err
			}
			cfg, err := store.Load()
			if err != nil {
				if errors.Is(err, os.ErrNotExist) {
					return nil // config doesn't exist yet (e.g. during init)
				}
				return fmt.Errorf("load config: %w", err)
			}
			if cfg.Lang != "" {
				i18n.SetLang(cfg.Lang)
			}
			return nil
		},
	}

	root.PersistentFlags().StringVarP(&configFlag, "config", "C", config.DefaultConfigFile, i18n.T("flag.config"))
	root.PersistentFlags().String("lang", "", i18n.T("flag.lang"))

	root.AddCommand(
		func() *cobra.Command {
			var interactive bool
			var providerNames []string
			cmd := &cobra.Command{
				Use:   "init",
				Short: i18n.T("cmd.init"),
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
						return wizard.RunInit(wizard.NewTerminalRunner(), projectRoot, store, reg)
					}
					return cli.RunInitNonInteractive(store, os.Stdout, cli.DefaultStyles())
				},
			}
			cmd.Flags().BoolVarP(&interactive, "interactive", "I", false, i18n.T("flag.interactive"))
			cmd.Flags().StringArrayVarP(&providerNames, "provider", "p", nil, i18n.T("flag.provider"))
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
				Short: i18n.T("cmd.add"),
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
					return cli.RunAdd(cfg, store, args[0], srv, os.Stdout, cli.DefaultStyles())
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
			Short: i18n.T("cmd.remove"),
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
				return cli.RunRemove(cfg, store, args[0], os.Stdout, cli.DefaultStyles())
			},
		},
		func() *cobra.Command {
			var jsonOutput bool
			cmd := &cobra.Command{
				Use:   "list",
				Short: i18n.T("cmd.list"),
				RunE: func(cmd *cobra.Command, args []string) error {
					store, err := loadConfig(configFlag)
					if err != nil {
						return err
					}
					cfg, err := store.Load()
					if err != nil {
						return fmt.Errorf("load config: %w", err)
					}
					return cli.RunList(cfg, os.Stdout, jsonOutput, cli.DefaultStyles())
				},
			}
			cmd.Flags().BoolVar(&jsonOutput, "json", false, i18n.T("flag.json"))
			return cmd
		}(),
		func() *cobra.Command {
			var (
				providerNames []string
				jsonOutput    bool
			)
			cmd := &cobra.Command{
				Use:   "sync",
				Short: i18n.T("cmd.sync"),
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
					// --provider flag takes precedence; fall back to providers saved in .iris.json.
					targetProviders := providerNames
					if len(targetProviders) == 0 {
						targetProviders = cfg.Providers
					}
					if len(targetProviders) > 0 {
						reg, err = reg.Filter(targetProviders)
						if err != nil {
							return fmt.Errorf("filter providers: %w", err)
						}
					}
					return cli.RunSync(projectRoot, cfg, reg, os.Stdout, jsonOutput, cli.DefaultStyles())
				},
			}
			cmd.Flags().StringArrayVarP(&providerNames, "provider", "p", nil, i18n.T("flag.provider"))
			cmd.Flags().BoolVar(&jsonOutput, "json", false, i18n.T("flag.json"))
			return cmd
		}(),
		func() *cobra.Command {
			var jsonOutput bool
			cmd := &cobra.Command{
				Use:   "status",
				Short: i18n.T("cmd.status"),
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
					return cli.RunStatus(projectRoot, cfg, buildRegistry(), os.Stdout, jsonOutput, cli.DefaultStyles())
				},
			}
			cmd.Flags().BoolVar(&jsonOutput, "json", false, i18n.T("flag.json"))
			return cmd
		}(),
	)

	if err := root.Execute(); err != nil {
		code := cli.ExitGeneral
		if errors.Is(err, ierrors.ErrServerNotFound) || errors.Is(err, ierrors.ErrProviderNotFound) {
			code = cli.ExitNotFound
		} else if errors.Is(err, ierrors.ErrConfigPermission) {
			code = cli.ExitPermission
		}
		os.Exit(code)
	}
}
