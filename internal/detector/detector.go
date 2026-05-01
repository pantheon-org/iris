package detector

import "github.com/pantheon-org/iris/internal/providers"

// Detect returns all providers from registry whose config file exists in projectRoot.
// Providers with SupportsProjectConfig=false are skipped (they are global, always synced).
func Detect(projectRoot string, registry *providers.Registry) []providers.Provider {
	var detected []providers.Provider
	for _, p := range registry.All() {
		if !p.Config().SupportsProjectConfig {
			continue
		}
		if p.Exists(projectRoot) {
			detected = append(detected, p)
		}
	}
	return detected
}
