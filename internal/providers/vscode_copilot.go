package providers

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pantheon-org/iris/internal/ierrors"
	"github.com/pantheon-org/iris/internal/types"
)

// VSCode Copilot uses a "servers" key (not "mcpServers") and each entry
// has a sibling "type" field alongside "command"/"args"/"env".
type VSCodeCopilotProvider struct {
	configPath func(projectRoot string) string
}

func NewVSCodeCopilotProvider() *VSCodeCopilotProvider {
	return &VSCodeCopilotProvider{
		configPath: func(projectRoot string) string {
			return filepath.Join(projectRoot, ".vscode", "mcp.json")
		},
	}
}

func (p *VSCodeCopilotProvider) Config() ProviderConfig {
	return ProviderConfig{
		Name:                  NameGitHubCopilot,
		DisplayName:           "GitHub Copilot",
		ConfigPath:            ".vscode/mcp.json",
		SupportsProjectConfig: true,
	}
}

func (p *VSCodeCopilotProvider) ConfigFilePath(projectRoot string) string {
	return p.configPath(projectRoot)
}

func (p *VSCodeCopilotProvider) Exists(projectRoot string) (bool, error) {
	_, err := os.Stat(p.ConfigFilePath(projectRoot))
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, fmt.Errorf("stat config: %w", err)
}

type vscodeServer struct {
	Type    string            `json:"type,omitempty"`
	Command string            `json:"command,omitempty"`
	Args    []string          `json:"args,omitempty"`
	Env     map[string]string `json:"env,omitempty"`
	URL     string            `json:"url,omitempty"`
}

func (p *VSCodeCopilotProvider) Generate(servers map[string]types.MCPServer, existingContent string) (string, error) {
	root := make(map[string]json.RawMessage)
	if existingContent != "" {
		if err := json.Unmarshal([]byte(existingContent), &root); err != nil {
			return "", fmt.Errorf("parse existing content: %w", ierrors.ErrMalformedConfig)
		}
	}

	out := make(map[string]vscodeServer, len(servers))
	for name, srv := range servers {
		entry := vscodeServer{
			Command: srv.Command,
			Args:    srv.Args,
			Env:     srv.Env,
			URL:     srv.URL,
		}
		switch srv.Transport {
		case types.TransportHTTP:
			// VS Code Copilot uses "http" for Streamable HTTP (tries HTTP, falls back to SSE).
			entry.Type = "http"
		case types.TransportSSE:
			entry.Type = "sse"
		default:
			if srv.URL != "" {
				entry.Type = "sse"
			} else {
				entry.Type = "stdio"
			}
		}
		out[name] = entry
	}

	encoded, err := json.Marshal(out)
	if err != nil {
		return "", fmt.Errorf("marshal servers: %w", err)
	}
	root["servers"] = json.RawMessage(encoded)

	result, err := json.MarshalIndent(root, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshal output: %w", err)
	}
	return string(result), nil
}

func (p *VSCodeCopilotProvider) Parse(content string) (map[string]types.MCPServer, error) {
	var root struct {
		Servers map[string]vscodeServer `json:"servers"`
	}
	if err := json.Unmarshal([]byte(content), &root); err != nil {
		return nil, fmt.Errorf("%w: %w", ierrors.ErrMalformedConfig, err)
	}

	result := make(map[string]types.MCPServer, len(root.Servers))
	for name, s := range root.Servers {
		var transport types.Transport
		switch s.Type {
		case "http":
			transport = types.TransportHTTP
		case "sse":
			transport = types.TransportSSE
		default:
			// Infer from URL if type is missing or unrecognised.
			if s.URL != "" {
				transport = types.TransportSSE
			} else {
				transport = types.TransportStdio
			}
		}
		result[name] = types.MCPServer{
			Command:   s.Command,
			Args:      s.Args,
			Env:       s.Env,
			URL:       s.URL,
			Transport: transport,
		}
	}
	return result, nil
}
