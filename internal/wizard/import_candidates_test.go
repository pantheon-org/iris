package wizard_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pantheon-org/iris/internal/providers"
	"github.com/pantheon-org/iris/internal/registry"
	"github.com/pantheon-org/iris/internal/types"
	"github.com/pantheon-org/iris/internal/wizard"
)

// newIsolatedClaudeCodeProvider returns a ClaudeCodeProvider whose global config path
// points at a non-existent file in dir, preventing reads from the real ~/.claude.json.
func newIsolatedClaudeCodeProvider(dir string) *providers.ClaudeCodeProvider {
	return providers.NewClaudeCodeProviderWithGlobalPath(filepath.Join(dir, "no-global.json"))
}

func TestCollectImportCandidates_noProviderConfigs_returnsEmpty(t *testing.T) {
	dir := t.TempDir()
	reg := registry.NewRegistry()
	reg.Register(newIsolatedClaudeCodeProvider(dir))

	got, err := wizard.CollectImportCandidates(dir, reg)
	require.NoError(t, err)
	assert.Empty(t, got)
}

func TestCollectImportCandidates_projectConfig_returnsProjectScopedCandidate(t *testing.T) {
	dir := t.TempDir()
	mcpJSON := `{"mcpServers":{"fmt":{"command":"npx","args":["-y","foo"],"type":"stdio"}}}`
	require.NoError(t, os.WriteFile(filepath.Join(dir, ".mcp.json"), []byte(mcpJSON), 0o600))

	reg := registry.NewRegistry()
	reg.Register(newIsolatedClaudeCodeProvider(dir))

	got, err := wizard.CollectImportCandidates(dir, reg)
	require.NoError(t, err)

	require.Len(t, got, 1)
	assert.Equal(t, "fmt", got[0].ServerName)
	assert.Equal(t, "claude", got[0].ProviderName)
	assert.Equal(t, wizard.ScopeProject, got[0].Scope)
	assert.Equal(t, "npx", got[0].Server.Command)
}

func TestCollectImportCandidates_globalConfig_returnsGlobalScopedCandidate(t *testing.T) {
	dir := t.TempDir()
	globalPath := filepath.Join(dir, "gemini-settings.json")
	geminiJSON := `{"mcpServers":{"gsrv":{"command":"python","args":["-m","gemini_mcp"]}}}`
	require.NoError(t, os.WriteFile(globalPath, []byte(geminiJSON), 0o600))

	reg := registry.NewRegistry()
	reg.Register(providers.NewGoogleGeminiProviderWithPath(globalPath))

	got, err := wizard.CollectImportCandidates("", reg)
	require.NoError(t, err)

	require.Len(t, got, 1)
	assert.Equal(t, "gsrv", got[0].ServerName)
	assert.Equal(t, wizard.ScopeGlobal, got[0].Scope)
}

func TestCollectImportCandidates_multipleProviders_returnsAllCandidates(t *testing.T) {
	dir := t.TempDir()

	mcpJSON := `{"mcpServers":{"fmt":{"command":"npx","args":["-y","foo"],"type":"stdio"}}}`
	require.NoError(t, os.WriteFile(filepath.Join(dir, ".mcp.json"), []byte(mcpJSON), 0o600))

	cursorDir := filepath.Join(dir, ".cursor")
	require.NoError(t, os.MkdirAll(cursorDir, 0o700))
	cursorJSON := `{"mcpServers":{"github":{"command":"uvx","args":["mcp-server-github"],"type":"stdio"}}}`
	require.NoError(t, os.WriteFile(filepath.Join(cursorDir, "mcp.json"), []byte(cursorJSON), 0o600))

	reg := registry.NewRegistry()
	reg.Register(newIsolatedClaudeCodeProvider(dir))
	reg.Register(providers.NewCursorProvider())

	got, err := wizard.CollectImportCandidates(dir, reg)
	require.NoError(t, err)
	assert.Len(t, got, 2)

	names := make(map[string]bool)
	for _, c := range got {
		names[c.ServerName] = true
	}
	assert.True(t, names["fmt"])
	assert.True(t, names["github"])
}

func TestCollectImportCandidates_sameServerInTwoProviders_returnsTwo(t *testing.T) {
	dir := t.TempDir()

	mcpJSON := `{"mcpServers":{"shared":{"command":"npx","args":["-y","foo"],"type":"stdio"}}}`
	require.NoError(t, os.WriteFile(filepath.Join(dir, ".mcp.json"), []byte(mcpJSON), 0o600))

	cursorDir := filepath.Join(dir, ".cursor")
	require.NoError(t, os.MkdirAll(cursorDir, 0o700))
	require.NoError(t, os.WriteFile(filepath.Join(cursorDir, "mcp.json"), []byte(mcpJSON), 0o600))

	reg := registry.NewRegistry()
	reg.Register(newIsolatedClaudeCodeProvider(dir))
	reg.Register(providers.NewCursorProvider())

	got, err := wizard.CollectImportCandidates(dir, reg)
	require.NoError(t, err)
	assert.Len(t, got, 2)
}

func TestCollectImportCandidates_malformedProjectConfig_skippedSilently(t *testing.T) {
	dir := t.TempDir()
	// Write a syntactically invalid JSON file where the provider config should be.
	require.NoError(t, os.WriteFile(filepath.Join(dir, ".mcp.json"), []byte(`{"mcpServers": {`), 0o600))

	reg := registry.NewRegistry()
	reg.Register(newIsolatedClaudeCodeProvider(dir))

	got, err := wizard.CollectImportCandidates(dir, reg)
	require.NoError(t, err, "malformed provider config must not abort init")
	assert.Empty(t, got)
}

func TestCollectImportCandidates_malformedGlobalConfig_skippedSilently(t *testing.T) {
	dir := t.TempDir()
	globalPath := filepath.Join(dir, "gemini-settings.json")
	require.NoError(t, os.WriteFile(globalPath, []byte(`{"mcpServers": {`), 0o600))

	reg := registry.NewRegistry()
	reg.Register(providers.NewGoogleGeminiProviderWithPath(globalPath))

	got, err := wizard.CollectImportCandidates("", reg)
	require.NoError(t, err, "malformed global config must not abort init")
	assert.Empty(t, got)
}

func TestCollectImportCandidates_oneMalformedOneValid_returnsValid(t *testing.T) {
	dir := t.TempDir()

	// Claude Code — malformed.
	require.NoError(t, os.WriteFile(filepath.Join(dir, ".mcp.json"), []byte(`{bad json}`), 0o600))

	// Cursor — valid.
	cursorDir := filepath.Join(dir, ".cursor")
	require.NoError(t, os.MkdirAll(cursorDir, 0o700))
	cursorJSON := `{"mcpServers":{"github":{"command":"uvx","args":["mcp-server-github"],"type":"stdio"}}}`
	require.NoError(t, os.WriteFile(filepath.Join(cursorDir, "mcp.json"), []byte(cursorJSON), 0o600))

	reg := registry.NewRegistry()
	reg.Register(newIsolatedClaudeCodeProvider(dir))
	reg.Register(providers.NewCursorProvider())

	got, err := wizard.CollectImportCandidates(dir, reg)
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, "github", got[0].ServerName)
}

func TestImportCandidate_Label_formatIsNameProviderScope(t *testing.T) {
	c := wizard.ImportCandidate{
		ServerName:   "fmt",
		ProviderName: "claude",
		Scope:        wizard.ScopeProject,
	}
	assert.Equal(t, "fmt  [claude] [project]", c.Label())
}

// ── GroupImportCandidates ─────────────────────────────────────────────────────

func TestGroupImportCandidates_empty_returnsEmpty(t *testing.T) {
	got := wizard.GroupImportCandidates(nil)
	assert.Empty(t, got)
}

func TestGroupImportCandidates_uniqueNames_oneEntryEach(t *testing.T) {
	candidates := []wizard.ImportCandidate{
		{ServerName: "fmt", ProviderName: "claude", Scope: wizard.ScopeProject},
		{ServerName: "github", ProviderName: "cursor", Scope: wizard.ScopeProject},
	}
	got := wizard.GroupImportCandidates(candidates)
	require.Len(t, got, 2)
	assert.Equal(t, "fmt", got[0].ServerName)
	assert.Equal(t, []string{"claude"}, got[0].Providers)
	assert.Equal(t, "github", got[1].ServerName)
}

func TestGroupImportCandidates_sameNameTwoProviders_collapsedIntoOne(t *testing.T) {
	candidates := []wizard.ImportCandidate{
		{ServerName: "shared", ProviderName: "claude", Scope: wizard.ScopeProject},
		{ServerName: "shared", ProviderName: "cursor", Scope: wizard.ScopeProject},
	}
	got := wizard.GroupImportCandidates(candidates)
	require.Len(t, got, 1)
	assert.Equal(t, "shared", got[0].ServerName)
	assert.Equal(t, []string{"claude", "cursor"}, got[0].Providers)
}

func TestGroupImportCandidates_sameNameSameProvider_deduplicatesProvider(t *testing.T) {
	candidates := []wizard.ImportCandidate{
		{ServerName: "shared", ProviderName: "claude", Scope: wizard.ScopeProject},
		{ServerName: "shared", ProviderName: "claude", Scope: wizard.ScopeGlobal},
	}
	got := wizard.GroupImportCandidates(candidates)
	require.Len(t, got, 1)
	assert.Equal(t, []string{"claude"}, got[0].Providers)
}

func TestGroupImportCandidates_preservesFirstDefinition(t *testing.T) {
	srv1 := types.MCPServer{Command: "npx"}
	srv2 := types.MCPServer{Command: "uvx"}
	candidates := []wizard.ImportCandidate{
		{ServerName: "shared", Server: srv1, ProviderName: "claude"},
		{ServerName: "shared", Server: srv2, ProviderName: "cursor"},
	}
	got := wizard.GroupImportCandidates(candidates)
	require.Len(t, got, 1)
	assert.Equal(t, "npx", got[0].Server.Command)
}

func TestGroupedCandidate_Label_singleProvider(t *testing.T) {
	g := wizard.GroupedCandidate{ServerName: "fmt", Providers: []string{"claude"}}
	assert.Equal(t, "fmt  [claude]", g.Label())
}

func TestGroupedCandidate_Label_multipleProviders(t *testing.T) {
	g := wizard.GroupedCandidate{ServerName: "fmt", Providers: []string{"claude", "cursor", "codex"}}
	assert.Equal(t, "fmt  [claude · cursor · codex]", g.Label())
}

// ── PromptMultiSelect scripted runner ─────────────────────────────────────────

func TestScriptedRunner_PromptMultiSelect_emptyAnswer_returnsNil(t *testing.T) {
	r := wizard.NewScriptedRunner([]string{""})
	got, err := r.PromptMultiSelect("pick", []string{"a", "b", "c"})
	require.NoError(t, err)
	assert.Nil(t, got)
}

func TestScriptedRunner_PromptMultiSelect_noneKeyword_returnsNil(t *testing.T) {
	r := wizard.NewScriptedRunner([]string{"none"})
	got, err := r.PromptMultiSelect("pick", []string{"a", "b"})
	require.NoError(t, err)
	assert.Nil(t, got)
}

func TestScriptedRunner_PromptMultiSelect_allKeyword_returnsAllIndices(t *testing.T) {
	r := wizard.NewScriptedRunner([]string{"all"})
	got, err := r.PromptMultiSelect("pick", []string{"a", "b", "c"})
	require.NoError(t, err)
	assert.Equal(t, []int{0, 1, 2}, got)
}

func TestScriptedRunner_PromptMultiSelect_commaSeparatedIndices_returnsSelected(t *testing.T) {
	r := wizard.NewScriptedRunner([]string{"0,2"})
	got, err := r.PromptMultiSelect("pick", []string{"a", "b", "c"})
	require.NoError(t, err)
	assert.Equal(t, []int{0, 2}, got)
}

func TestScriptedRunner_PromptMultiSelect_outOfRange_returnsError(t *testing.T) {
	r := wizard.NewScriptedRunner([]string{"5"})
	_, err := r.PromptMultiSelect("pick", []string{"a", "b"})
	require.Error(t, err)
}
