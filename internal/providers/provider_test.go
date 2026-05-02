package providers_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pantheon-org/iris/internal/providers"
)

func TestOpenCodeProvider_Config_UsesExpectedConfigFileName(t *testing.T) {
	cfg := providers.NewOpenCodeProvider().Config()
	assert.NotEmpty(t, cfg.GlobalConfigPath)
	assert.Contains(t, cfg.GlobalConfigPath, "opencode")
	assert.Contains(t, cfg.GlobalConfigPath, "opencode.json")
}

func TestClaudeDesktopProvider_Config_UsesExpectedConfigFileName(t *testing.T) {
	cfg := providers.NewClaudeDesktopProvider().Config()
	assert.NotEmpty(t, cfg.GlobalConfigPath)
	assert.Contains(t, cfg.GlobalConfigPath, "Claude")
	assert.Contains(t, cfg.GlobalConfigPath, "claude_desktop_config.json")
}
