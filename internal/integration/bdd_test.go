package integration_test

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cucumber/godog"

	"github.com/pantheon-org/iris/internal/config"
	"github.com/pantheon-org/iris/internal/providers"
	"github.com/pantheon-org/iris/internal/registry"
	irisync "github.com/pantheon-org/iris/internal/sync"
	"github.com/pantheon-org/iris/internal/types"
	"github.com/pantheon-org/iris/internal/wizard"
)

// scenarioCtx holds per-scenario mutable state.
type scenarioCtx struct {
	root      string
	storePath string
	store     *config.Store
	cfg       *types.IrisConfig
	reg       *registry.Registry

	// captured output / results
	lastErr           error
	output            *bytes.Buffer
	syncResults       []irisync.SyncResult
	reloadedCfg       *types.IrisConfig
	importCandidates  []wizard.ImportCandidate
	groupedCandidates []wizard.GroupedCandidate
}

func newScenarioCtx(root string) *scenarioCtx {
	return &scenarioCtx{
		root:   root,
		output: &bytes.Buffer{},
	}
}

// buildReg constructs a full registry with all 14 providers, paths pinned under root.
func buildReg(root string) *registry.Registry {
	googleGeminiPath := filepath.Join(root, "gemini-settings.json")
	codexPath := filepath.Join(root, "codex-config.toml")
	claudeDesktopPath := filepath.Join(root, "claude-desktop-config.json")
	windsurfPath := filepath.Join(root, "windsurf-config.json")
	zedPath := filepath.Join(root, "zed-settings.json")
	warpPath := filepath.Join(root, "warp-mcp.json")
	kimiPath := filepath.Join(root, "kimi-settings.json")
	mistralVibePath := filepath.Join(root, "mistral-vibe-config.toml")

	claudeCodeGlobalPath := filepath.Join(root, "claude-global.json")
	reg := registry.NewRegistry()
	reg.Register(providers.NewClaudeCodeProviderWithGlobalPath(claudeCodeGlobalPath))
	reg.Register(providers.NewClaudeDesktopProviderWithPath(claudeDesktopPath))
	reg.Register(providers.NewGoogleGeminiProviderWithPath(googleGeminiPath))
	reg.Register(providers.NewOpenCodeProvider())
	reg.Register(providers.NewOpenaiCodexProviderWithPath(codexPath))
	reg.Register(providers.NewCursorProvider())
	reg.Register(providers.NewWindsurfProviderWithPath(windsurfPath))
	reg.Register(providers.NewVSCodeCopilotProvider())
	reg.Register(providers.NewZedProviderWithPath(zedPath))
	reg.Register(providers.NewQwenProvider())
	reg.Register(providers.NewWarpProviderWithPath(warpPath))
	reg.Register(providers.NewKimiProviderWithPath(kimiPath))
	reg.Register(providers.NewMistralVibeProviderWithPath(mistralVibePath))
	reg.Register(providers.NewIntelliJProvider())
	return reg
}

// ── shared setup ──────────────────────────────────────────────────────────────

func (s *scenarioCtx) aCleanWorkspace() error {
	s.storePath = filepath.Join(s.root, ".iris.json")
	store, err := config.NewStore(s.storePath)
	if err != nil {
		return fmt.Errorf("NewStore: %w", err)
	}
	s.store = store
	s.cfg = &types.IrisConfig{
		Version: 1,
		Servers: make(map[string]types.MCPServer),
	}
	s.reg = buildReg(s.root)
	return nil
}

// ── helpers ───────────────────────────────────────────────────────────────────

func splitArgs(raw string) []string {
	if raw == "" {
		return nil
	}
	return strings.Split(raw, ",")
}

func parseEnvPairs(raw string) map[string]string {
	env := make(map[string]string)
	for _, pair := range strings.Split(raw, ",") {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) == 2 {
			env[parts[0]] = parts[1]
		}
	}
	return env
}

// ── suite wiring ──────────────────────────────────────────────────────────────

func initializeScenario(t *testing.T) func(ctx *godog.ScenarioContext) {
	t.Helper()
	return func(sc *godog.ScenarioContext) {
		root := t.TempDir()
		s := newScenarioCtx(root)

		// setup
		sc.Step(`^a clean workspace$`, s.aCleanWorkspace)
		sc.Step(`^the iris config already exists with one server$`, s.theIrisConfigAlreadyExistsWithOneServer)

		// add — Given form (no error capture)
		sc.Step(`^an MCP server "([^"]+)" with command "([^"]+)" and args "([^"]+)"$`, s.anMCPServerWithCommandAndArgs)
		sc.Step(`^an MCP server "([^"]+)" with command "([^"]+)" and no args$`, s.anMCPServerWithCommandAndNoArgs)
		sc.Step(`^an MCP server "([^"]+)" with command "([^"]+)" and env "([^"]+)"$`, s.anMCPServerWithCommandAndEnv)
		sc.Step(`^an SSE server "([^"]+)" with URL "([^"]+)"$`, s.anSSEServerWithURL)
		sc.Step(`^the server "([^"]+)" has env var "([^"]+)" set to "([^"]+)"$`, s.theServerHasEnvVarSetTo)

		// add — When form (error capture)
		sc.Step(`^I add a stdio server "([^"]+)" with command "([^"]+)" and args "([^"]+)"$`, s.iAddAStdioServerWithCommandAndArgs)
		sc.Step(`^I add a stdio server "([^"]+)" with command "([^"]+)" and no args$`, s.iAddAStdioServerWithCommandAndNoArgs)
		sc.Step(`^I add an SSE server "([^"]+)" with url "([^"]+)"$`, s.iAddAnSSEServerWithURL)
		sc.Step(`^I add a stdio server "([^"]+)" with command "([^"]+)" and env "([^"]+)"$`, s.iAddAStdioServerWithCommandAndEnv)
		sc.Step(`^I try to add a stdio server "([^"]+)" with no command$`, s.iTryToAddAStdioServerWithNoCommand)
		sc.Step(`^I add an MCP server "([^"]+)" with command "([^"]+)" and args "([^"]+)"$`, s.iAddAnMCPServerWithCommandAndArgs)

		// remove
		sc.Step(`^I remove the server "([^"]+)"$`, s.iRemoveTheServer)
		sc.Step(`^I try to remove the server "([^"]+)"$`, s.iTryToRemoveTheServer)

		// list
		sc.Step(`^I run list$`, s.iRunList)
		sc.Step(`^I run list with JSON output$`, s.iRunListWithJSONOutput)

		// sync
		sc.Step(`^I sync to all providers$`, s.iSyncToAllProviders)
		sc.Step(`^I sync to all providers again$`, s.iSyncToAllProvidersAgain)
		sc.Step(`^I sync to all providers with JSON output$`, s.iSyncToAllProvidersWithJSONOutput)

		// status
		sc.Step(`^I run status$`, s.iRunStatus)
		sc.Step(`^I run status with JSON output$`, s.iRunStatusWithJSONOutput)
		sc.Step(`^I corrupt the provider config file "([^"]+)"$`, s.iCorruptTheProviderConfigFile)

		// init (non-interactive)
		sc.Step(`^I run init$`, s.iRunInit)

		// init (interactive)
		sc.Step(`^no provider config files exist$`, s.noProviderConfigFilesExist)
		sc.Step(`^a malformed Claude Code project config exists$`, s.aMalformedClaudeCodeProjectConfigExists)
		sc.Step(`^a Claude Code project config exists with server "([^"]+)" command "([^"]+)" args "([^"]+)"$`, s.aClaudeCodeProjectConfigExistsWithServer)
		sc.Step(`^a Cursor project config exists with server "([^"]+)" command "([^"]+)" args "([^"]+)"$`, s.aCursorProjectConfigExistsWithServer)
		sc.Step(`^a global Google Gemini config exists with server "([^"]+)" command "([^"]+)" args "([^"]+)"$`, s.aGlobalGoogleGeminiConfigExistsWithServer)
		sc.Step(`^I run interactive init and select no servers$`, s.iRunInteractiveInitAndSelectNoServers)
		sc.Step(`^I run interactive init and collect the import candidates$`, s.iRunInteractiveInitAndCollectImportCandidates)
		sc.Step(`^I run interactive init and select server "([^"]+)"$`, s.iRunInteractiveInitAndSelectServer)
		sc.Step(`^I run interactive init and select all discovered servers$`, s.iRunInteractiveInitAndSelectAllDiscoveredServers)
		sc.Step(`^I run interactive init, skip import, and manually add server "([^"]+)" command "([^"]+)" args "([^"]+)"$`, s.iRunInteractiveInitSkipImportAndManuallyAddServer)
		sc.Step(`^I run interactive init, import server "([^"]+)", and manually add server "([^"]+)" command "([^"]+)" args "([^"]+)"$`, s.iRunInteractiveInitImportServerAndManuallyAddServer)

		// reload
		sc.Step(`^I reload the config from disk$`, s.iReloadTheConfigFromDisk)

		// assertions — errors
		sc.Step(`^the last error wraps "([^"]+)"$`, s.theLastErrorWraps)
		sc.Step(`^the last error is ErrServerNotFound$`, s.theLastErrorIsErrServerNotFound)

		// assertions — in-memory config
		sc.Step(`^the config contains server "([^"]+)" with command "([^"]+)"$`, s.theConfigContainsServerWithCommand)
		sc.Step(`^the config contains server "([^"]+)" with transport "([^"]+)"$`, s.theConfigContainsServerWithTransport)
		sc.Step(`^the config contains server "([^"]+)" with env var "([^"]+)" equal to "([^"]+)"$`, s.theConfigContainsServerWithEnvVar)
		sc.Step(`^the config has exactly (\d+) server$`, s.theConfigHasExactlyNServers)
		sc.Step(`^the iris config file exists on disk$`, s.theIrisConfigFileExistsOnDisk)

		// assertions — reloaded config
		sc.Step(`^the config contains (\d+) servers$`, s.theConfigContainsNServers)
		sc.Step(`^the config does not contain server "([^"]+)"$`, s.theConfigDoesNotContainServer)
		sc.Step(`^the config contains server "([^"]+)" with command "([^"]+)"$`, s.theReloadedConfigContainsServerWithCommand)

		// assertions — provider files
		sc.Step(`^the provider config file "([^"]+)" exists$`, s.theProviderConfigFileExists)
		sc.Step(`^the provider config file "([^"]+)" does not exist$`, s.theProviderConfigFileDoesNotExist)
		sc.Step(`^the JSON provider file "([^"]+)" contains servers "([^"]+)" under key "([^"]+)"$`, s.theJSONProviderFileContainsServersUnderKey)
		sc.Step(`^the JSON provider file "([^"]+)" server "([^"]+)" under key "([^"]+)" has field "([^"]+)"$`, s.theJSONProviderServerHasField)
		sc.Step(`^the JSON provider file "([^"]+)" still has key "([^"]+)"$`, s.theJSONProviderFileStillHasKey)
		sc.Step(`^a provider file "([^"]+)" exists with extra key "([^"]+)" set to "([^"]+)"$`, s.aProviderFileExistsWithExtraKey)
		sc.Step(`^the opencode provider file "([^"]+)" contains servers "([^"]+)"$`, s.theOpencodeProviderFileContainsServers)
		sc.Step(`^the opencode server "([^"]+)" in file "([^"]+)" has correct field format$`, s.theOpencodeServerHasCorrectFieldFormat)
		sc.Step(`^the TOML provider file "([^"]+)" contains servers "([^"]+)"$`, s.theTOMLProviderFileContainsServers)
		sc.Step(`^the zed provider file "([^"]+)" contains servers "([^"]+)"$`, s.theZedProviderFileContainsServers)
		sc.Step(`^the TOML mistral provider file "([^"]+)" contains servers "([^"]+)"$`, s.theTOMLMistralProviderFileContainsServers)
		sc.Step(`^the copilot server "([^"]+)" in file "([^"]+)" does not have field "([^"]+)"$`, s.theCopilotServerDoesNotHaveField)
		sc.Step(`^the JSON provider file "([^"]+)" server "([^"]+)" under key "([^"]+)" has env var "([^"]+)"$`, s.theJSONProviderServerHasEnvVar)

		// assertions — sync results
		sc.Step(`^all providers report status "([^"]+)"$`, s.allProvidersReportStatus)

		// assertions — text output
		sc.Step(`^the output contains "([^"]+)"$`, s.theOutputContains)
		sc.Step(`^the output lines appear in order "([^"]+)"$`, s.theOutputLinesAppearInOrder)

		// assertions — JSON list
		sc.Step(`^the JSON output has a "servers" array$`, s.theJSONOutputHasAServersArray)
		sc.Step(`^the JSON servers array contains an entry with name "([^"]+)" and command "([^"]+)"$`, s.theJSONServersArrayContainsEntryWithNameAndCommand)

		// assertions — JSON sync
		sc.Step(`^the JSON sync output has a "results" array$`, s.theJSONSyncOutputHasAResultsArray)
		sc.Step(`^the JSON sync results contain an entry for provider "([^"]+)" with status "([^"]+)"$`, s.theJSONSyncResultsContainEntryForProviderWithStatus)

		// assertions — JSON status
		sc.Step(`^the JSON status output has a "providers" array$`, s.theJSONStatusOutputHasAProvidersArray)
		sc.Step(`^the JSON status providers contain an entry for provider "([^"]+)" with status "([^"]+)"$`, s.theJSONStatusProvidersContainEntryForProviderWithStatus)

		// assertions — text status
		sc.Step(`^the status output contains provider "([^"]+)" with status "([^"]+)"$`, s.theStatusOutputContainsProviderWithStatus)

		// assertions — init
		sc.Step(`^the iris config file is valid JSON with version 1$`, s.theIrisConfigFileIsValidJSONWithVersion1)
		sc.Step(`^the iris config contains (\d+) servers$`, s.theIrisConfigContainsNServers)
		sc.Step(`^the iris config contains server "([^"]+)" with command "([^"]+)"$`, s.theIrisConfigContainsServerWithCommand)
		sc.Step(`^the import candidates include an entry for server "([^"]+)" from provider "([^"]+)" with scope "([^"]+)"$`, s.theImportCandidatesIncludeEntry)
		sc.Step(`^the grouped candidates contain exactly (\d+) entry for server "([^"]+)"$`, s.theGroupedCandidatesContainExactlyNEntryForServer)
		sc.Step(`^the grouped candidate for server "([^"]+)" lists providers "([^"]+)" and "([^"]+)"$`, s.theGroupedCandidateForServerListsProviders)
		sc.Step(`^the iris config providers list contains "([^"]+)"$`, s.theIrisConfigProvidersListContains)
		sc.Step(`^the iris config providers list is set to "([^"]+)"$`, s.theIrisConfigProvidersListIsSetTo)
	}
}

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		Name:                 "iris integration",
		TestSuiteInitializer: nil,
		ScenarioInitializer:  initializeScenario(t),
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t,
		},
	}
	if suite.Run() != 0 {
		t.Fatal("BDD integration tests failed")
	}
}

// Keep the context import used by godog internally.
var _ context.Context
