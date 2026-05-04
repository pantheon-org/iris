package registry_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pantheon-org/iris/internal/ierrors"
	"github.com/pantheon-org/iris/internal/providers"
	"github.com/pantheon-org/iris/internal/registry"
	"github.com/pantheon-org/iris/internal/types"
)

type mockProvider struct {
	name string
}

func (m *mockProvider) Config() providers.ProviderConfig {
	return providers.ProviderConfig{Name: m.name, DisplayName: m.name}
}

func (m *mockProvider) Generate(_ map[string]types.MCPServer, _ string) (string, error) {
	return "", nil
}

func (m *mockProvider) Parse(_ string) (map[string]types.MCPServer, error) {
	return nil, nil
}

func (m *mockProvider) ConfigFilePath(_ string) string { return "" }

func (m *mockProvider) Exists(_ string) (bool, error) { return false, nil }

func TestRegistry_NewRegistry_IsEmpty(t *testing.T) {
	r := registry.NewRegistry()
	assert.Empty(t, r.All())
	assert.Empty(t, r.Names())
}

func TestRegistry_RegisterAndGet_ReturnsProvider(t *testing.T) {
	r := registry.NewRegistry()
	p := &mockProvider{name: "claude"}
	r.Register(p)

	got, err := r.Get("claude")
	require.NoError(t, err)
	assert.Equal(t, p, got)
}

func TestRegistry_Get_UnknownName_WrapsErrProviderNotFound(t *testing.T) {
	r := registry.NewRegistry()

	_, err := r.Get("unknown")
	require.Error(t, err)
	assert.True(t, errors.Is(err, ierrors.ErrProviderNotFound))
}

func TestRegistry_All_ReturnsAllProviders(t *testing.T) {
	r := registry.NewRegistry()
	r.Register(&mockProvider{name: "claude"})
	r.Register(&mockProvider{name: "cursor"})

	all := r.All()
	assert.Len(t, all, 2)
}

func TestRegistry_Names_ReturnsAllProviderNames(t *testing.T) {
	r := registry.NewRegistry()
	r.Register(&mockProvider{name: "claude"})
	r.Register(&mockProvider{name: "gemini"})

	names := r.Names()
	assert.Len(t, names, 2)
	assert.Contains(t, names, "claude")
	assert.Contains(t, names, "gemini")
}

func TestRegistry_Filter_ReturnsSubset(t *testing.T) {
	r := registry.NewRegistry()
	r.Register(&mockProvider{name: "claude"})
	r.Register(&mockProvider{name: "gemini"})
	r.Register(&mockProvider{name: "cursor"})

	filtered, err := r.Filter([]string{"claude", "cursor"})
	require.NoError(t, err)
	assert.Len(t, filtered.All(), 2)
	assert.Contains(t, filtered.Names(), "claude")
	assert.Contains(t, filtered.Names(), "cursor")
}

func TestRegistry_Filter_UnknownName_WrapsErrProviderNotFound(t *testing.T) {
	r := registry.NewRegistry()
	r.Register(&mockProvider{name: "claude"})

	_, err := r.Filter([]string{"claude", "nope"})
	require.Error(t, err)
	assert.True(t, errors.Is(err, ierrors.ErrProviderNotFound))
}

func TestRegistry_Filter_UnknownName_ErrorContainsProviderName(t *testing.T) {
	r := registry.NewRegistry()
	r.Register(&mockProvider{name: "claude"})

	_, err := r.Filter([]string{"claude", "unknown-provider"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unknown-provider",
		"error message should contain the missing provider name for easier debugging")
	assert.True(t, errors.Is(err, ierrors.ErrProviderNotFound),
		"error should still unwrap to ErrProviderNotFound via errors.Is")
}

func TestRegistry_Filter_EmptySlice_ReturnsEmptyRegistry(t *testing.T) {
	r := registry.NewRegistry()
	r.Register(&mockProvider{name: "claude"})

	filtered, err := r.Filter([]string{})
	require.NoError(t, err)
	assert.Empty(t, filtered.All())
}
