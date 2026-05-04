package providers

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pantheon-org/iris/internal/ierrors"
	"github.com/pantheon-org/iris/internal/io"
	"github.com/pantheon-org/iris/internal/types"
)

type OpenCodeProvider struct{}

func NewOpenCodeProvider() *OpenCodeProvider {
	return &OpenCodeProvider{}
}

func (p *OpenCodeProvider) Config() ProviderConfig {
	return ProviderConfig{
		Name:                  NameAnomalycoOpenCode,
		DisplayName:           "Anomalyco OpenCode",
		ConfigPath:            "opencode.json",
		SupportsProjectConfig: true,
		GlobalConfigPath:      io.UserConfigPath("opencode", "opencode.json"),
	}
}

func (p *OpenCodeProvider) ConfigFilePath(projectRoot string) string {
	return filepath.Join(projectRoot, "opencode.json")
}

func (p *OpenCodeProvider) Exists(projectRoot string) (bool, error) {
	_, err := os.Stat(p.ConfigFilePath(projectRoot))
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, fmt.Errorf("stat config: %w", err)
}

type opencodeMCPEntry struct {
	Command     []string          `json:"command,omitempty"`
	Type        string            `json:"type,omitempty"`
	Enabled     *bool             `json:"enabled,omitempty"`
	Environment map[string]string `json:"environment,omitempty"`
	URL         string            `json:"url,omitempty"`
	Headers     map[string]string `json:"headers,omitempty"`
	Cwd         string            `json:"cwd,omitempty"`
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

		var enabled *bool
		if srv.Enabled != nil {
			v := *srv.Enabled
			enabled = &v
		} else {
			v := true
			enabled = &v
		}

		entryType := "local"
		if srv.Transport == types.TransportSSE || (srv.Transport == "" && srv.URL != "") {
			entryType = "remote"
		}

		mcp[name] = opencodeMCPEntry{
			Command:     cmd,
			Type:        entryType,
			Enabled:     enabled,
			Environment: srv.Env,
			URL:         srv.URL,
			Headers:     srv.Headers,
			Cwd:         srv.Cwd,
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

		transport := types.TransportStdio
		if entry.Type == "remote" || entry.URL != "" {
			transport = types.TransportSSE
		}

		enabled := entry.Enabled
		servers[name] = types.MCPServer{
			Command:   cmd,
			Args:      args,
			Env:       entry.Environment,
			URL:       entry.URL,
			Headers:   entry.Headers,
			Cwd:       entry.Cwd,
			Enabled:   enabled,
			Transport: transport,
		}
	}
	return servers, nil
}
