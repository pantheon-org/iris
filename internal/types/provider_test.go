package types_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pantheon-org/iris/internal/types"
)

func TestProviderNames_ReturnsStringSlice(t *testing.T) {
	input := []types.ProviderName{
		types.NameAnthropicClaudeCode,
		types.NameGoogleGemini,
		types.NameAnomalycoOpenCode,
	}
	got := types.ProviderNames(input)
	assert.Equal(t, []string{"claude", "gemini", "opencode"}, got)
}

func TestProviderNames_EmptySlice_ReturnsEmpty(t *testing.T) {
	got := types.ProviderNames([]types.ProviderName{})
	assert.Empty(t, got)
	assert.NotNil(t, got)
}
