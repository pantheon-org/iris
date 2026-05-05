package sync

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pantheon-org/iris/internal/ierrors"
	"github.com/pantheon-org/iris/internal/providers"
	"github.com/pantheon-org/iris/internal/registry"
	"github.com/pantheon-org/iris/internal/types"
)

type SyncStatus string

const (
	SyncStatusCreated   SyncStatus = "created"
	SyncStatusUpdated   SyncStatus = "updated"
	SyncStatusUnchanged SyncStatus = "unchanged"
	SyncStatusError     SyncStatus = "error"
)

// SyncScope controls which config paths are synced per provider.
type SyncScope int

const (
	// ScopeAll syncs both global and local config paths for every provider.
	ScopeAll SyncScope = iota
	// ScopeGlobal syncs only the global (home-directory) config path.
	ScopeGlobal
	// ScopeLocal syncs only the project-local config path.
	ScopeLocal
)

type SyncResult struct {
	ProviderName string
	Scope        string // "global" | "local" | ""
	Path         string // resolved config file path
	Status       SyncStatus
	Err          error
}

func SyncProvider(projectRoot string, p providers.Provider, servers map[string]types.MCPServer) SyncResult {
	return syncProviderScoped(p.ConfigFilePath(projectRoot), "", p, servers)
}

// syncProviderScoped syncs p at a resolved configPath, tagging the result with scope.
func syncProviderScoped(configPath, scopeTag string, p providers.Provider, servers map[string]types.MCPServer) SyncResult {
	name := p.Config().Name
	base := SyncResult{ProviderName: name, Scope: scopeTag, Path: configPath}

	existingContent := ""
	fileAbsent := false

	raw, err := os.ReadFile(configPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return SyncResult{ProviderName: name, Scope: scopeTag, Path: configPath, Status: SyncStatusError, Err: err}
		}
		fileAbsent = true
	} else {
		existingContent = string(raw)
	}

	newContent, err := p.Generate(servers, existingContent)
	if err != nil {
		return SyncResult{ProviderName: name, Scope: scopeTag, Path: configPath, Status: SyncStatusError, Err: err}
	}

	if !fileAbsent && newContent == existingContent {
		base.Status = SyncStatusUnchanged
		return base
	}

	if info, err := os.Lstat(configPath); err == nil && info.Mode()&os.ModeSymlink != 0 {
		return SyncResult{
			ProviderName: name,
			Scope:        scopeTag,
			Path:         configPath,
			Status:       SyncStatusError,
			Err:          fmt.Errorf("config path %s is a symlink: %w", configPath, ierrors.ErrSymlinkNotAllowed),
		}
	}

	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return SyncResult{ProviderName: name, Scope: scopeTag, Path: configPath, Status: SyncStatusError, Err: err}
	}

	if err := os.WriteFile(configPath, []byte(newContent), 0644); err != nil {
		return SyncResult{ProviderName: name, Scope: scopeTag, Path: configPath, Status: SyncStatusError, Err: err}
	}

	if fileAbsent {
		base.Status = SyncStatusCreated
		return base
	}
	base.Status = SyncStatusUpdated
	return base
}

func SyncAllProviders(projectRoot string, scope SyncScope, reg *registry.Registry, servers map[string]types.MCPServer) []SyncResult {
	all := reg.All()
	results := make([]SyncResult, 0, len(all))
	for _, p := range all {
		cfg := p.Config()
		switch scope {
		case ScopeGlobal:
			if !cfg.HasGlobalConfig {
				continue
			}
			results = append(results, syncProviderScoped(p.ConfigFilePath(""), "global", p, servers))
		case ScopeLocal:
			if cfg.LocalConfigPath == nil {
				// Global-only providers: fall back to their primary path.
				results = append(results, syncProviderScoped(p.ConfigFilePath(projectRoot), "", p, servers))
				continue
			}
			// Use ConfigFilePath so providers with custom path logic (e.g. pinned paths) are respected.
			results = append(results, syncProviderScoped(p.ConfigFilePath(projectRoot), "local", p, servers))
		default: // ScopeAll
			hasGlobal := cfg.HasGlobalConfig
			hasLocal := cfg.LocalConfigPath != nil && projectRoot != ""

			if hasGlobal && hasLocal {
				results = append(results, syncProviderScoped(p.ConfigFilePath(""), "global", p, servers))
				// Use ConfigFilePath for local so custom path logic is respected.
				results = append(results, syncProviderScoped(p.ConfigFilePath(projectRoot), "local", p, servers))
			} else if hasGlobal {
				results = append(results, syncProviderScoped(p.ConfigFilePath(""), "global", p, servers))
			} else {
				// Provider has only local or unscoped path.
				results = append(results, syncProviderScoped(p.ConfigFilePath(projectRoot), "", p, servers))
			}
		}
	}
	return results
}
