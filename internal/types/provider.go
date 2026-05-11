package types

type ProviderName string

func ProviderNames(n []ProviderName) []string {
	result := make([]string, len(n))
	for i, p := range n {
		result[i] = string(p)
	}
	return result
}

const (
	NameAnthropicClaudeCode    = "claude"
	NameAnthropicClaudeDesktop = "claude-desktop"
	NameGoogleGemini           = "gemini"
	NameAnomalycoOpenCode      = "opencode"
	NameOpenAICodex            = "codex"
	NameAnysphereCursor        = "cursor"
	NameOpenAIWindsurf         = "windsurf"
	NameGitHubCopilot          = "copilot"
	NameZed                    = "zed"
	NameAlibabaQwenCode        = "qwen"
	NameWarpTerminal           = "warp"
	NameMoonshotKimi           = "kimi"
	NameMistralAIVibe          = "mistral-vibe"
	NameIntelliJIDEA           = "intellij"
	NameCline                  = "cline"
	NameKiro                   = "kiro"
)

type Provider struct {
	Name             ProviderName `json:"name" yaml:"name" toml:"name"`
	LocalConfigPath  string       `json:"localConfigPath,omitempty" yaml:"localConfigPath,omitempty" toml:"localConfigPath,omitempty"`
	GlobalConfigPath string       `json:"globalConfigPath,omitempty" yaml:"globalConfigPath,omitempty" toml:"globalConfigPath,omitempty"`
}
