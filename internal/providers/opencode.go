package providers

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pantheon-org/iris/internal/ierrors"
	"github.com/pantheon-org/iris/internal/types"
)

type OpenCodeProvider struct{}

func NewOpenCodeProvider() *OpenCodeProvider {
	return &OpenCodeProvider{}
}

func (p *OpenCodeProvider) Config() ProviderConfig {
	home, _ := os.UserHomeDir()
	return ProviderConfig{
		Name:                  "opencode",
		DisplayName:           "OpenCode",
		ConfigPath:            "opencode.json",
		SupportsProjectConfig: true,
		GlobalConfigPath:      filepath.Join(home, ".config", "opencode", "opencode.json"),
	}
}

func (p *OpenCodeProvider) ConfigFilePath(projectRoot string) string {
	return filepath.Join(projectRoot, "opencode.json")
}

func (p *OpenCodeProvider) Exists(projectRoot string) bool {
	_, err := os.Stat(p.ConfigFilePath(projectRoot))
	return err == nil
}

type opencodeMCPEntry struct {
	Command     []string          `json:"command"`
	Type        string            `json:"type"`
	Enabled     bool              `json:"enabled"`
	Environment map[string]string `json:"environment"`
}

func (p *OpenCodeProvider) Generate(servers map[string]types.MCPServer, existingContent string) (string, error) {
	root := make(map[string]json.RawMessage)

	if existingContent != "" {
		if err := json.Unmarshal([]byte(existingContent), &root); err != nil {
			return "", fmt.Errorf("parse existing content: %w: %w", ierrors.ErrMalformedConfig, err)
		}
	}

	mcp := make(map[string]opencodeMCPEntry, len(servers))
	for name, srv := range servers {
		cmd := make([]string, 0, 1+len(srv.Args))
		if srv.Command != "" {
			cmd = append(cmd, srv.Command)
		}
		cmd = append(cmd, srv.Args...)

		enabled := true
		if srv.Enabled != nil {
			enabled = *srv.Enabled
		}

		env := srv.Env
		if env == nil {
			env = map[string]string{}
		}

		mcp[name] = opencodeMCPEntry{
			Command:     cmd,
			Type:        "local",
			Enabled:     enabled,
			Environment: env,
		}
	}

	mcpBytes, err := json.Marshal(mcp)
	if err != nil {
		return "", fmt.Errorf("marshal mcp block: %w", err)
	}
	root["mcp"] = json.RawMessage(mcpBytes)

	out, err := json.MarshalIndent(root, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshal output: %w", err)
	}
	return string(out), nil
}

func (p *OpenCodeProvider) Parse(content string) (map[string]types.MCPServer, error) {
	var root struct {
		MCP map[string]opencodeMCPEntry `json:"mcp"`
	}
	if err := json.Unmarshal([]byte(content), &root); err != nil {
		return nil, fmt.Errorf("%w: %w", ierrors.ErrMalformedConfig, err)
	}

	servers := make(map[string]types.MCPServer, len(root.MCP))
	for name, entry := range root.MCP {
		var cmd string
		var args []string
		if len(entry.Command) > 0 {
			cmd = entry.Command[0]
			args = entry.Command[1:]
		}

		enabled := entry.Enabled
		servers[name] = types.MCPServer{
			Command: cmd,
			Args:    args,
			Env:     entry.Environment,
			Enabled: &enabled,
		}
	}
	return servers, nil
}
