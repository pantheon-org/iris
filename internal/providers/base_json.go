package providers

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pantheon-org/iris/internal/ierrors"
	"github.com/pantheon-org/iris/internal/types"
)

type mcpServerJSON struct {
	Command string            `json:"command,omitempty"`
	Args    []string          `json:"args,omitempty"`
	Env     map[string]string `json:"env,omitempty"`
	Type    string            `json:"type,omitempty"`
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

func (b *baseJSONProvider) Exists(projectRoot string) bool {
	_, err := os.Stat(b.ConfigFilePath(projectRoot))
	return err == nil
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
		mcpServers[name] = mcpServerJSON{
			Command: srv.Command,
			Args:    srv.Args,
			Env:     srv.Env,
			Type:    string(types.TransportStdio),
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
		result[name] = types.MCPServer{
			Command:   s.Command,
			Args:      s.Args,
			Env:       s.Env,
			Transport: types.Transport(s.Type),
		}
	}
	return result, nil
}

func geminiConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(".gemini", "settings.json")
	}
	return filepath.Join(home, ".gemini", "settings.json")
}
