package providers

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/pantheon-org/iris/internal/ierrors"
	"github.com/pantheon-org/iris/internal/io"
	"github.com/pantheon-org/iris/internal/types"
)

// Zed uses a "context_servers" key in ~/.config/zed/settings.json.
// Stdio servers: { "command": "str", "args": [...], "env": {} }
// Remote servers: { "url": "str", "headers": {} }
type ZedProvider struct {
	configPath string
}

func NewZedProvider() *ZedProvider {
	return newZedProviderWithPath(zedConfigPath())
}

func newZedProviderWithPath(path string) *ZedProvider {
	return &ZedProvider{configPath: path}
}

// NewZedProviderWithPath creates a ZedProvider using a custom config path.
// Intended for use in tests.
func NewZedProviderWithPath(path string) *ZedProvider {
	return newZedProviderWithPath(path)
}

func (p *ZedProvider) Config() ProviderConfig {
	return ProviderConfig{
		Name:                  NameZed,
		DisplayName:           "Zed",
		LocalConfigPath:       nil,
		SupportsProjectConfig: false,
		GlobalConfigPath:      homeRel(p.configPath),
	}
}

func (p *ZedProvider) ConfigFilePath(_ string) string { return p.configPath }

func (p *ZedProvider) Exists(_ string) (bool, error) {
	_, err := os.Stat(p.configPath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, fmt.Errorf("stat config: %w", err)
}

type zedContextServer struct {
	Command string            `json:"command,omitempty"`
	Args    []string          `json:"args,omitempty"`
	Env     map[string]string `json:"env,omitempty"`
	URL     string            `json:"url,omitempty"`
	Headers map[string]string `json:"headers,omitempty"`
}

func (p *ZedProvider) Generate(servers map[string]types.MCPServer, existingContent string) (string, error) {
	root := make(map[string]json.RawMessage)
	if existingContent != "" {
		if err := json.Unmarshal([]byte(existingContent), &root); err != nil {
			return "", fmt.Errorf("parse existing content: %w", ierrors.ErrMalformedConfig)
		}
	}

	out := make(map[string]zedContextServer, len(servers))
	for name, srv := range servers {
		entry := zedContextServer{
			Command: srv.Command,
			Args:    srv.Args,
			Env:     srv.Env,
			URL:     srv.URL,
			Headers: srv.Headers,
		}
		out[name] = entry
	}

	encoded, err := json.Marshal(out)
	if err != nil {
		return "", fmt.Errorf("marshal context_servers: %w", err)
	}
	root["context_servers"] = json.RawMessage(encoded)

	result, err := json.MarshalIndent(root, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshal output: %w", err)
	}
	return string(result), nil
}

func (p *ZedProvider) Parse(content string) (map[string]types.MCPServer, error) {
	var root struct {
		ContextServers map[string]zedContextServer `json:"context_servers"`
	}
	if err := json.Unmarshal([]byte(content), &root); err != nil {
		return nil, fmt.Errorf("%w: %w", ierrors.ErrMalformedConfig, err)
	}

	result := make(map[string]types.MCPServer, len(root.ContextServers))
	for name, s := range root.ContextServers {
		transport := types.TransportStdio
		if s.URL != "" {
			transport = types.TransportSSE
		}
		result[name] = types.MCPServer{
			Command:   s.Command,
			Args:      s.Args,
			Env:       s.Env,
			URL:       s.URL,
			Headers:   s.Headers,
			Transport: transport,
		}
	}
	return result, nil
}

func zedConfigPath() string { return io.UserHomePath(".config", "zed", "settings.json") }
