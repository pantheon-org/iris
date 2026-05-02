package providers

import (
	"fmt"

	"github.com/pantheon-org/iris/internal/ierrors"
	"github.com/pantheon-org/iris/internal/types"
)

type ProviderConfig struct {
	Name                  string
	DisplayName           string
	ConfigPath            string
	SupportsProjectConfig bool
	GlobalConfigPath      string
}

type Provider interface {
	Config() ProviderConfig
	Generate(servers map[string]types.MCPServer, existingContent string) (string, error)
	Parse(content string) (map[string]types.MCPServer, error)
	ConfigFilePath(projectRoot string) string
	Exists(projectRoot string) bool
}

type Registry struct {
	providers map[string]Provider
}

func NewRegistry() *Registry {
	return &Registry{providers: make(map[string]Provider)}
}

func (r *Registry) Register(p Provider) {
	r.providers[p.Config().Name] = p
}

func (r *Registry) Get(name string) (Provider, error) {
	p, ok := r.providers[name]
	if !ok {
		return nil, fmt.Errorf("provider %q: %w", name, ierrors.ErrProviderNotFound)
	}
	return p, nil
}

func (r *Registry) All() []Provider {
	result := make([]Provider, 0, len(r.providers))
	for _, p := range r.providers {
		result = append(result, p)
	}
	return result
}

func (r *Registry) Names() []string {
	names := make([]string, 0, len(r.providers))
	for name := range r.providers {
		names = append(names, name)
	}
	return names
}

// Filter returns a new Registry containing only the named providers.
// Returns an error wrapping ErrProviderNotFound if any name is absent.
func (r *Registry) Filter(names []string) (*Registry, error) {
	filtered := NewRegistry()
	for _, name := range names {
		p, err := r.Get(name)
		if err != nil {
			return nil, err
		}
		filtered.Register(p)
	}
	return filtered, nil
}
