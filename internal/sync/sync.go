package sync

import (
	"os"
	"path/filepath"

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

type SyncResult struct {
	ProviderName string
	Status       SyncStatus
	Err          error
}

func SyncProvider(projectRoot string, p providers.Provider, servers map[string]types.MCPServer) SyncResult {
	name := p.Config().Name
	configPath := p.ConfigFilePath(projectRoot)

	existingContent := ""
	fileAbsent := false

	raw, err := os.ReadFile(configPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return SyncResult{ProviderName: name, Status: SyncStatusError, Err: err}
		}
		fileAbsent = true
	} else {
		existingContent = string(raw)
	}

	newContent, err := p.Generate(servers, existingContent)
	if err != nil {
		return SyncResult{ProviderName: name, Status: SyncStatusError, Err: err}
	}

	if !fileAbsent && newContent == existingContent {
		return SyncResult{ProviderName: name, Status: SyncStatusUnchanged}
	}

	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return SyncResult{ProviderName: name, Status: SyncStatusError, Err: err}
	}

	if err := os.WriteFile(configPath, []byte(newContent), 0644); err != nil {
		return SyncResult{ProviderName: name, Status: SyncStatusError, Err: err}
	}

	if fileAbsent {
		return SyncResult{ProviderName: name, Status: SyncStatusCreated}
	}
	return SyncResult{ProviderName: name, Status: SyncStatusUpdated}
}

func SyncAllProviders(projectRoot string, registry *registry.Registry, servers map[string]types.MCPServer) []SyncResult {
	all := registry.All()
	results := make([]SyncResult, 0, len(all))
	for _, p := range all {
		results = append(results, SyncProvider(projectRoot, p, servers))
	}
	return results
}
