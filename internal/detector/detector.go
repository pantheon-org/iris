package detector

import (
	"fmt"

	"github.com/pantheon-org/iris/internal/providers"
	"github.com/pantheon-org/iris/internal/registry"
)

// Detect returns all providers from registry whose config file exists in projectRoot.
// Providers with SupportsProjectConfig=false are skipped (they are global, always synced).
// If any provider's Exists() returns an IO error (i.e. not simply "file not found"),
// Detect surfaces that error immediately rather than silently skipping the provider.
func Detect(projectRoot string, registry *registry.Registry) ([]providers.Provider, error) {
	var detected []providers.Provider
	for _, p := range registry.All() {
		if !p.Config().SupportsProjectConfig {
			continue
		}
		ok, err := p.Exists(projectRoot)
		if err != nil {
			return nil, fmt.Errorf("check provider %s: %w", p.Config().Name, err)
		}
		if ok {
			detected = append(detected, p)
		}
	}
	return detected, nil
}
