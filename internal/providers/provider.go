package providers

import (
	"path/filepath"
	"strings"

	irio "github.com/pantheon-org/iris/internal/io"
	"github.com/pantheon-org/iris/internal/types"
)

type ProviderConfig struct {
	Name                  string
	DisplayName           string
	LocalConfigPath       *string // project-relative path (e.g. ".mcp.json"); nil = project config not supported
	SupportsProjectConfig bool
	GlobalConfigPath      *string // home-relative path (e.g. ".config/tool/config.json"); nil = no global config
}

func strPtr(s string) *string { return &s }

// homeRel returns a pointer to the home-relative form of absPath (e.g. ".config/tool/config.json").
// Returns nil if absPath is not under the user's home directory.
func homeRel(absPath string) *string {
	home := irio.UserHomeDir()
	rel, err := filepath.Rel(home, absPath)
	if err != nil || strings.HasPrefix(rel, "..") {
		return nil
	}
	return strPtr(rel)
}

type Provider interface {
	Config() ProviderConfig
	Generate(servers map[string]types.MCPServer, existingContent string) (string, error)
	Parse(content string) (map[string]types.MCPServer, error)
	ConfigFilePath(projectRoot string) string
	Exists(projectRoot string) (bool, error)
}
