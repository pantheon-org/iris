package providers

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/pantheon-org/iris/internal/ierrors"
	"github.com/pantheon-org/iris/internal/io"
	"github.com/pantheon-org/iris/internal/types"
)

type mcpServerJSON struct {
	Command string            `json:"command,omitempty"`
	Args    []string          `json:"args,omitempty"`
	Env     map[string]string `json:"env,omitempty"`
	Type    string            `json:"type,omitempty"`
	URL     string            `json:"url,omitempty"`
	Headers map[string]string `json:"headers,omitempty"`
	Cwd     string            `json:"cwd,omitempty"`
	Enabled *bool             `json:"enabled,omitempty"`
}

type baseJSONProvider struct {
	config       ProviderConfig
	resolvedPath func(projectRoot string) string
}

func (b *baseJSONProvider) Config() ProviderConfig {
	return b.config
}

func (b *baseJSONProvider) ConfigFilePath(projectRoot string) string {
	return b.resolvedPath(projectRoot)
}

// SafeConfigFilePath validates projectRoot for path traversal before returning
// the config file path. Returns ErrPathTraversal if projectRoot contains ".."
// components.
func (b *baseJSONProvider) SafeConfigFilePath(projectRoot string) (string, error) {
	if err := ValidateProjectRoot(projectRoot); err != nil {
		return "", err
	}
	return b.resolvedPath(projectRoot), nil
}

func (b *baseJSONProvider) Exists(projectRoot string) (bool, error) {
	_, err := os.Stat(b.ConfigFilePath(projectRoot))
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, fmt.Errorf("stat config: %w", err)
}

func (b *baseJSONProvider) Generate(servers map[string]types.MCPServer, existingContent string) (string, error) {
	raw := make(map[string]json.RawMessage)

	if existingContent != "" {
		if err := json.Unmarshal([]byte(existingContent), &raw); err != nil {
			return "", fmt.Errorf("parse existing content: %w", ierrors.ErrMalformedConfig)
		}
	}

	mcpServers := make(map[string]mcpServerJSON, len(servers))
	for name, srv := range servers {
		serverType := string(srv.Transport)
		if serverType == "" {
			serverType = string(types.TransportStdio)
		}
		if srv.URL != "" {
			serverType = string(types.TransportSSE)
		}
		mcpServers[name] = mcpServerJSON{
			Command: srv.Command,
			Args:    srv.Args,
			Env:     srv.Env,
			Type:    serverType,
			URL:     srv.URL,
			Headers: srv.Headers,
			Cwd:     srv.Cwd,
			Enabled: srv.Enabled,
		}
	}

	encoded, err := json.Marshal(mcpServers)
	if err != nil {
		return "", fmt.Errorf("marshal mcpServers: %w", err)
	}
	raw["mcpServers"] = json.RawMessage(encoded)

	out, err := json.MarshalIndent(raw, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshal output: %w", err)
	}
	return string(out), nil
}

func (b *baseJSONProvider) Parse(content string) (map[string]types.MCPServer, error) {
	var doc struct {
		MCPServers map[string]mcpServerJSON `json:"mcpServers"`
	}
	if err := json.Unmarshal([]byte(content), &doc); err != nil {
		return nil, fmt.Errorf("parse config: %w", ierrors.ErrMalformedConfig)
	}

	result := make(map[string]types.MCPServer, len(doc.MCPServers))
	for name, s := range doc.MCPServers {
		transport := types.TransportStdio
		if s.Type != "" {
			transport = types.Transport(s.Type)
		}
		if s.URL != "" {
			transport = types.TransportSSE
		}
		result[name] = types.MCPServer{
			Command:   s.Command,
			Args:      s.Args,
			Env:       s.Env,
			URL:       s.URL,
			Headers:   s.Headers,
			Cwd:       s.Cwd,
			Enabled:   s.Enabled,
			Transport: transport,
		}
	}
	return result, nil
}

func googleGeminiConfigPath() string { return io.UserHomePath(".gemini", "settings.json") }
