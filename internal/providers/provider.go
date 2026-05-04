package providers

import (
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
	Exists(projectRoot string) (bool, error)
}
