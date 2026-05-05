package providers

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"

	"github.com/BurntSushi/toml"

	"github.com/pantheon-org/iris/internal/ierrors"
	"github.com/pantheon-org/iris/internal/io"
	"github.com/pantheon-org/iris/internal/types"
)

type OpenaiCodexProvider struct {
	configPath string
	pinned     bool // when true, ConfigFilePath always returns configPath regardless of projectRoot
}

func codexGlobalPath() string { return io.UserHomePath(".codex", "config.toml") }

func NewOpenaiCodexProvider() *OpenaiCodexProvider {
	return &OpenaiCodexProvider{configPath: codexGlobalPath()}
}

// NewOpenaiCodexProviderWithPath creates a OpenaiCodexProvider pinned to a fixed config path.
// Intended for use in tests.
func NewOpenaiCodexProviderWithPath(path string) *OpenaiCodexProvider {
	return &OpenaiCodexProvider{configPath: path, pinned: true}
}

func (p *OpenaiCodexProvider) Config() ProviderConfig {
	return ProviderConfig{
		Name:                  NameOpenAICodex,
		DisplayName:           "OpenAI Codex",
		LocalConfigPath:       strPtr(".codex/config.toml"),
		SupportsProjectConfig: true,
		GlobalConfigPath:      homeRel(p.configPath),
		HasGlobalConfig:       true,
	}
}

func (p *OpenaiCodexProvider) ConfigFilePath(projectRoot string) string {
	if !p.pinned && projectRoot != "" {
		return filepath.Join(projectRoot, ".codex", "config.toml")
	}
	return p.configPath
}

func (p *OpenaiCodexProvider) Exists(projectRoot string) (bool, error) {
	_, err := os.Stat(p.ConfigFilePath(projectRoot))
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, fmt.Errorf("stat config: %w", err)
}

type codexMCPServer struct {
	Command     string            `toml:"command,omitempty"`
	Args        []string          `toml:"args,omitempty"`
	Env         map[string]string `toml:"env,omitempty"`
	URL         string            `toml:"url,omitempty"`
	HTTPHeaders map[string]string `toml:"http_headers,omitempty"`
	Cwd         string            `toml:"cwd,omitempty"`
	Enabled     *bool             `toml:"enabled,omitempty"`
}

func (p *OpenaiCodexProvider) Generate(servers map[string]types.MCPServer, existingContent string) (string, error) {
	// Collect top-level non-mcp_servers keys preserving original order via toml.MetaData.
	type topLevelEntry struct {
		key string
		val interface{}
	}
	var orderedTopLevel []topLevelEntry

	if existingContent != "" {
		existingServers, err := p.Parse(existingContent)
		if err != nil {
			return "", err
		}
		if reflect.DeepEqual(existingServers, servers) {
			return existingContent, nil
		}

		var rawMap map[string]interface{}
		meta, err := toml.Decode(existingContent, &rawMap)
		if err != nil {
			return "", fmt.Errorf("%w: %w", ierrors.ErrMalformedConfig, err)
		}
		for _, key := range meta.Keys() {
			// Only include top-level scalar/simple keys (skip mcp_servers and nested tables).
			if len(key) != 1 || key[0] == "mcp_servers" {
				continue
			}
			orderedTopLevel = append(orderedTopLevel, topLevelEntry{key: key[0], val: rawMap[key[0]]})
		}
	}

	// Build sorted server entries for deterministic output.
	names := make([]string, 0, len(servers))
	for name := range servers {
		names = append(names, name)
	}
	sort.Strings(names)

	mcpServers := make(map[string]codexMCPServer, len(names))
	for _, name := range names {
		srv := servers[name]
		mcpServers[name] = codexMCPServer{
			Command:     srv.Command,
			Args:        srv.Args,
			Env:         srv.Env,
			URL:         srv.URL,
			HTTPHeaders: srv.Headers,
			Cwd:         srv.Cwd,
			Enabled:     srv.Enabled,
		}
	}

	var buf bytes.Buffer

	// Write top-level keys in original order.
	for _, entry := range orderedTopLevel {
		single := map[string]interface{}{entry.key: entry.val}
		if err := toml.NewEncoder(&buf).Encode(single); err != nil {
			return "", fmt.Errorf("encode top-level key %q: %w", entry.key, err)
		}
	}

	// Encode the mcp_servers section using the TOML encoder for correctness.
	if len(mcpServers) > 0 {
		if err := toml.NewEncoder(&buf).Encode(map[string]interface{}{"mcp_servers": mcpServers}); err != nil {
			return "", fmt.Errorf("encode mcp_servers: %w", err)
		}
	}

	return buf.String(), nil
}

func (p *OpenaiCodexProvider) Parse(content string) (map[string]types.MCPServer, error) {
	var file struct {
		MCPServers map[string]codexMCPServer `toml:"mcp_servers"`
	}
	if _, err := toml.Decode(content, &file); err != nil {
		return nil, fmt.Errorf("%w: %w", ierrors.ErrMalformedConfig, err)
	}

	result := make(map[string]types.MCPServer, len(file.MCPServers))
	for name, entry := range file.MCPServers {
		transport := types.TransportStdio
		if entry.URL != "" {
			transport = types.TransportSSE
		}
		result[name] = types.MCPServer{
			Command:   entry.Command,
			Args:      entry.Args,
			Env:       entry.Env,
			URL:       entry.URL,
			Headers:   entry.HTTPHeaders,
			Cwd:       entry.Cwd,
			Enabled:   entry.Enabled,
			Transport: transport,
		}
	}
	return result, nil
}
