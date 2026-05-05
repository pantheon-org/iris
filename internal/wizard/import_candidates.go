package wizard

import (
	"errors"
	"fmt"
	"os"

	"github.com/pantheon-org/iris/internal/ierrors"
	"github.com/pantheon-org/iris/internal/providers"
	"github.com/pantheon-org/iris/internal/registry"
	"github.com/pantheon-org/iris/internal/types"
)

// Scope indicates whether a provider config was found at the project or global level.
type Scope string

const (
	ScopeProject Scope = "project"
	ScopeGlobal  Scope = "global"
)

// ImportCandidate is a single MCP server discovered in an existing provider config.
type ImportCandidate struct {
	ServerName   string
	Server       types.MCPServer
	ProviderName string
	Scope        Scope
}

// Label returns the display string shown in the multi-select prompt:
//
//	"<server-name>  [<provider>] [<scope>]"
func (c ImportCandidate) Label() string {
	return fmt.Sprintf("%s  [%s] [%s]", c.ServerName, c.ProviderName, c.Scope)
}

// CollectImportCandidates scans every provider in the registry for both project-level
// and global config files, parses them, and returns one ImportCandidate per discovered
// server. Providers that have no config file on disk are silently skipped.
func CollectImportCandidates(projectRoot string, reg *registry.Registry) ([]ImportCandidate, error) {
	var candidates []ImportCandidate

	for _, p := range reg.All() {
		cfg := p.Config()

		// Project-scoped config.
		if cfg.SupportsProjectConfig && projectRoot != "" {
			cs, err := readCandidates(p, projectRoot, ScopeProject)
			if err != nil {
				return nil, fmt.Errorf("read project config for %s: %w", cfg.Name, err)
			}
			candidates = append(candidates, cs...)
		}

		// Global config (when provider exposes one).
		if cfg.GlobalConfigPath != "" {
			cs, err := readCandidatesFromPath(p, cfg.GlobalConfigPath, cfg.Name, ScopeGlobal)
			if err != nil {
				return nil, fmt.Errorf("read global config for %s: %w", cfg.Name, err)
			}
			candidates = append(candidates, cs...)
		}
	}

	return candidates, nil
}

// readCandidates loads a provider's project-level config file and returns candidates.
func readCandidates(p providers.Provider, projectRoot string, scope Scope) ([]ImportCandidate, error) {
	ok, err := p.Exists(projectRoot)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}
	return readCandidatesFromPath(p, p.ConfigFilePath(projectRoot), p.Config().Name, scope)
}

// readCandidatesFromPath reads a config file at path, parses it, and returns candidates.
func readCandidatesFromPath(p providers.Provider, path, providerName string, scope Scope) ([]ImportCandidate, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	servers, err := p.Parse(string(data))
	if err != nil {
		// A malformed config is not fatal — skip this provider's file and continue.
		if errors.Is(err, ierrors.ErrMalformedConfig) {
			return nil, nil
		}
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	var out []ImportCandidate
	for name, srv := range servers {
		out = append(out, ImportCandidate{
			ServerName:   name,
			Server:       srv,
			ProviderName: providerName,
			Scope:        scope,
		})
	}
	return out, nil
}
