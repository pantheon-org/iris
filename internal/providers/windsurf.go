package providers

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/pantheon-org/iris/internal/ierrors"
	"github.com/pantheon-org/iris/internal/io"
	"github.com/pantheon-org/iris/internal/types"
)

// WindsurfProvider handles Windsurf's mcp_config.json.
// Remote servers use "serverUrl" (primary) rather than "url" per official docs.
// On parse we accept both "serverUrl" and "url" for compatibility.
type WindsurfProvider struct {
	configPath string
}

func NewWindsurfProvider() *WindsurfProvider {
	return newWindsurfProviderWithPath(windsurfConfigPath())
}

func newWindsurfProviderWithPath(path string) *WindsurfProvider {
	return &WindsurfProvider{configPath: path}
}

// NewWindsurfProviderWithPath creates a WindsurfProvider using a custom config path.
// Intended for use in tests.
func NewWindsurfProviderWithPath(path string) *WindsurfProvider {
	return newWindsurfProviderWithPath(path)
}

func windsurfConfigPath() string { return io.UserHomePath(".codeium", "windsurf", "mcp_config.json") }

func (p *WindsurfProvider) Config() ProviderConfig {
	return ProviderConfig{
		Name:                  NameOpenAIWindsurf,
		DisplayName:           "OpenAI Windsurf",
		LocalConfigPath:       nil,
		SupportsProjectConfig: false,
		GlobalConfigPath:      homeRel(p.configPath),
	}
}

func (p *WindsurfProvider) ConfigFilePath(_ string) string { return p.configPath }

func (p *WindsurfProvider) Exists(_ string) (bool, error) {
	_, err := os.Stat(p.configPath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, fmt.Errorf("stat config: %w", err)
}

// windsurfServerJSON is the wire format for a single Windsurf MCP server entry.
// Windsurf uses "serverUrl" as the primary remote URL field (official docs), but also
// accepts "url". On Generate we always write "serverUrl" for remote servers.
type windsurfServerJSON struct {
	Command   string            `json:"command,omitempty"`
	Args      []string          `json:"args,omitempty"`
	Env       map[string]string `json:"env,omitempty"`
	ServerURL string            `json:"serverUrl,omitempty"`
	URL       string            `json:"url,omitempty"` // accepted on parse only; not written
	Headers   map[string]string `json:"headers,omitempty"`
}

func (p *WindsurfProvider) Generate(servers map[string]types.MCPServer, existingContent string) (string, error) {
	raw := make(map[string]json.RawMessage)
	if existingContent != "" {
		if err := json.Unmarshal([]byte(existingContent), &raw); err != nil {
			return "", fmt.Errorf("parse existing content: %w", ierrors.ErrMalformedConfig)
		}
	}

	mcpServers := make(map[string]windsurfServerJSON, len(servers))
	for name, srv := range servers {
		entry := windsurfServerJSON{
			Command: srv.Command,
			Args:    srv.Args,
			Env:     srv.Env,
			Headers: srv.Headers,
		}
		if srv.Transport == types.TransportSSE || (srv.Transport == "" && srv.URL != "") {
			entry.ServerURL = srv.URL
		}
		mcpServers[name] = entry
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

func (p *WindsurfProvider) Parse(content string) (map[string]types.MCPServer, error) {
	var doc struct {
		MCPServers map[string]windsurfServerJSON `json:"mcpServers"`
	}
	if err := json.Unmarshal([]byte(content), &doc); err != nil {
		return nil, fmt.Errorf("parse config: %w", ierrors.ErrMalformedConfig)
	}

	result := make(map[string]types.MCPServer, len(doc.MCPServers))
	for name, s := range doc.MCPServers {
		// Accept serverUrl (primary) or url (fallback) as the remote URL.
		remoteURL := s.ServerURL
		if remoteURL == "" {
			remoteURL = s.URL
		}
		transport := types.TransportStdio
		if remoteURL != "" {
			transport = types.TransportSSE
		}
		result[name] = types.MCPServer{
			Command:   s.Command,
			Args:      s.Args,
			Env:       s.Env,
			URL:       remoteURL,
			Headers:   s.Headers,
			Transport: transport,
		}
	}
	return result, nil
}
