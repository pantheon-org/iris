package providers

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/pantheon-org/iris/internal/ierrors"
	"github.com/pantheon-org/iris/internal/io"
	"github.com/pantheon-org/iris/internal/types"
)

func googleGeminiConfigPath() string { return io.UserHomePath(".gemini", "settings.json") }

// GoogleGeminiProvider handles Gemini CLI's settings.json.
// Remote servers support both "url" (SSE) and "httpUrl" (HTTP streaming) per official docs.
type GoogleGeminiProvider struct {
	config       ProviderConfig
	resolvedPath func(projectRoot string) string
}

func NewGoogleGeminiProvider() *GoogleGeminiProvider {
	return newGoogleGeminiProviderWithPath(googleGeminiConfigPath())
}

func newGoogleGeminiProviderWithPath(path string) *GoogleGeminiProvider {
	p := &GoogleGeminiProvider{}
	p.config = ProviderConfig{
		Name:                  NameGoogleGemini,
		DisplayName:           "Google Gemini",
		ConfigPath:            "~/.gemini/settings.json",
		SupportsProjectConfig: true,
		GlobalConfigPath:      path,
	}
	p.resolvedPath = func(projectRoot string) string {
		if projectRoot != "" {
			return filepath.Join(projectRoot, ".gemini", "settings.json")
		}
		return path
	}
	return p
}

// NewGoogleGeminiProviderWithPath creates a GoogleGeminiProvider pinned to a fixed config path.
// Intended for use in tests.
func NewGoogleGeminiProviderWithPath(path string) *GoogleGeminiProvider {
	p := &GoogleGeminiProvider{}
	p.config = ProviderConfig{
		Name:                  NameGoogleGemini,
		DisplayName:           "Google Gemini",
		ConfigPath:            "~/.gemini/settings.json",
		SupportsProjectConfig: true,
		GlobalConfigPath:      path,
	}
	p.resolvedPath = func(_ string) string {
		return path
	}
	return p
}

func (p *GoogleGeminiProvider) Config() ProviderConfig { return p.config }

func (p *GoogleGeminiProvider) ConfigFilePath(projectRoot string) string {
	return p.resolvedPath(projectRoot)
}

func (p *GoogleGeminiProvider) Exists(projectRoot string) (bool, error) {
	if err := ValidateProjectRoot(projectRoot); err != nil {
		return false, err
	}
	return existsOnDisk(p.resolvedPath(projectRoot))
}

// geminiServerJSON is the wire format for a Gemini CLI MCP server entry.
// "url" is the SSE endpoint; "httpUrl" is the HTTP streaming endpoint.
// Transport is not specified explicitly — it is inferred from which URL field is set.
type geminiServerJSON struct {
	Command string            `json:"command,omitempty"`
	Args    []string          `json:"args,omitempty"`
	Env     map[string]string `json:"env,omitempty"`
	URL     string            `json:"url,omitempty"`
	HTTPUrl string            `json:"httpUrl,omitempty"`
	Headers map[string]string `json:"headers,omitempty"`
	Cwd     string            `json:"cwd,omitempty"`
}

func (p *GoogleGeminiProvider) Generate(servers map[string]types.MCPServer, existingContent string) (string, error) {
	raw := make(map[string]json.RawMessage)
	if existingContent != "" {
		if err := json.Unmarshal([]byte(existingContent), &raw); err != nil {
			return "", fmt.Errorf("parse existing content: %w", ierrors.ErrMalformedConfig)
		}
	}

	mcpServers := make(map[string]geminiServerJSON, len(servers))
	for name, srv := range servers {
		entry := geminiServerJSON{
			Command: srv.Command,
			Args:    srv.Args,
			Env:     srv.Env,
			Headers: srv.Headers,
			Cwd:     srv.Cwd,
		}
		switch srv.Transport {
		case types.TransportHTTP:
			entry.HTTPUrl = srv.URL
		case types.TransportSSE:
			entry.URL = srv.URL
		default:
			// Infer from URL presence: no transport set but URL given → SSE
			if srv.URL != "" {
				entry.URL = srv.URL
			}
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

func (p *GoogleGeminiProvider) Parse(content string) (map[string]types.MCPServer, error) {
	var doc struct {
		MCPServers map[string]geminiServerJSON `json:"mcpServers"`
	}
	if err := json.Unmarshal([]byte(content), &doc); err != nil {
		return nil, fmt.Errorf("%w: %w", ierrors.ErrMalformedConfig, err)
	}

	result := make(map[string]types.MCPServer, len(doc.MCPServers))
	for name, s := range doc.MCPServers {
		var transport types.Transport
		var remoteURL string
		switch {
		case s.HTTPUrl != "":
			transport = types.TransportHTTP
			remoteURL = s.HTTPUrl
		case s.URL != "":
			transport = types.TransportSSE
			remoteURL = s.URL
		default:
			transport = types.TransportStdio
		}
		result[name] = types.MCPServer{
			Command:   s.Command,
			Args:      s.Args,
			Env:       s.Env,
			URL:       remoteURL,
			Headers:   s.Headers,
			Cwd:       s.Cwd,
			Transport: transport,
		}
	}
	return result, nil
}
