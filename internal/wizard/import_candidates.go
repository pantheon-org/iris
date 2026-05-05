package wizard

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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

// GroupedCandidate represents a unique MCP server consolidated across all providers
// that define it. The server definition is taken from the first provider encountered.
type GroupedCandidate struct {
	ServerName string
	Server     types.MCPServer
	Providers  []string // ordered list of provider names that define this server
}

// Label returns the display string shown in the multi-select prompt:
//
//	"<server-name>  [<provider1> · <provider2> · ...]"
func (g GroupedCandidate) Label() string {
	return fmt.Sprintf("%s  [%s]", g.ServerName, strings.Join(g.Providers, " · "))
}

// GroupImportCandidates collapses a flat list of ImportCandidates into one entry per
// unique server name. The server definition is taken from the first occurrence; duplicate
// provider names are deduplicated while preserving insertion order.
func GroupImportCandidates(candidates []ImportCandidate) []GroupedCandidate {
	seen := make(map[string]int) // serverName → index in result
	var result []GroupedCandidate

	for _, c := range candidates {
		if idx, ok := seen[c.ServerName]; ok {
			// Append provider only if not already listed.
			g := &result[idx]
			alreadyListed := false
			for _, p := range g.Providers {
				if p == c.ProviderName {
					alreadyListed = true
					break
				}
			}
			if !alreadyListed {
				g.Providers = append(g.Providers, c.ProviderName)
			}
		} else {
			seen[c.ServerName] = len(result)
			result = append(result, GroupedCandidate{
				ServerName: c.ServerName,
				Server:     c.Server,
				Providers:  []string{c.ProviderName},
			})
		}
	}
	return result
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
		// ConfigFilePath("") returns an absolute path when the provider has a global
		// config; project-only providers return a relative path (no root prefix).
		if absGlobal := p.ConfigFilePath(""); filepath.IsAbs(absGlobal) {
			cs, err := readCandidatesFromPath(p, absGlobal, cfg.Name, ScopeGlobal)
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
		return nil, fmt.Errorf("check exists: %w", err)
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
