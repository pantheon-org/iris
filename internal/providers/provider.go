package providers

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pantheon-org/iris/internal/ierrors"
	"github.com/pantheon-org/iris/internal/types"
)

// homeDir returns the current user's home directory using os.UserHomeDir,
// which is OS-agnostic (Linux, macOS, Windows). Falls back to "." on error
// so callers always get a usable path rather than a silent empty-string root.
func homeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "."
	}
	return home
}

// homePath joins the user's home directory with the given path elements.
func homePath(elem ...string) string {
	return filepath.Join(append([]string{homeDir()}, elem...)...)
}

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
