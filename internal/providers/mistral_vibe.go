package providers

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/BurntSushi/toml"

	"github.com/pantheon-org/iris/internal/ierrors"
	"github.com/pantheon-org/iris/internal/io"
	"github.com/pantheon-org/iris/internal/types"
)

// Mistral Vibe stores MCP servers in ~/.vibe/config.toml under [[mcp_servers]].
// The format mirrors Codex: a TOML array of tables with an explicit name field.

func mistralVibeConfigPath() string { return io.UserHomePath(".vibe", "config.toml") }

func NewMistralVibeProvider() *MistralVibeProvider {
	return &MistralVibeProvider{configPath: mistralVibeConfigPath()}
}

// MistralVibeProvider with pinned=true always returns configPath regardless of projectRoot.
type MistralVibeProvider struct {
	configPath string
	pinned     bool
}

// NewMistralVibeProviderWithPath creates a MistralVibeProvider pinned to a fixed config path.
// Intended for use in tests.
func NewMistralVibeProviderWithPath(path string) *MistralVibeProvider {
	return &MistralVibeProvider{configPath: path, pinned: true}
}

func (p *MistralVibeProvider) Config() ProviderConfig {
	return ProviderConfig{
		Name:                  NameMistralAIVibe,
		DisplayName:           "Mistral AI Vibe",
		ConfigPath:            "~/.vibe/config.toml",
		SupportsProjectConfig: true,
		GlobalConfigPath:      p.configPath,
	}
}

func (p *MistralVibeProvider) ConfigFilePath(projectRoot string) string {
	if !p.pinned && projectRoot != "" {
		return filepath.Join(projectRoot, ".vibe", "config.toml")
	}
	return p.configPath
}

func (p *MistralVibeProvider) Exists(projectRoot string) (bool, error) {
	_, err := os.Stat(p.ConfigFilePath(projectRoot))
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, fmt.Errorf("stat config: %w", err)
}

type vibeMCPServer struct {
	Name      string            `toml:"name"`
	Transport string            `toml:"transport"`
	Command   string            `toml:"command,omitempty"`
	Args      []string          `toml:"args,omitempty"`
	Env       map[string]string `toml:"env,omitempty"`
	URL       string            `toml:"url,omitempty"`
	Headers   map[string]string `toml:"headers,omitempty"`
}

func (p *MistralVibeProvider) Generate(servers map[string]types.MCPServer, existingContent string) (string, error) {
	var orderedTopLevel []struct {
		key string
		val interface{}
	}

	if existingContent != "" {
		var rawMap map[string]interface{}
		meta, err := toml.Decode(existingContent, &rawMap)
		if err != nil {
			return "", fmt.Errorf("%w: %w", ierrors.ErrMalformedConfig, err)
		}
		for _, key := range meta.Keys() {
			if len(key) != 1 || key[0] == "mcp_servers" {
				continue
			}
			orderedTopLevel = append(orderedTopLevel, struct {
				key string
				val interface{}
			}{key: key[0], val: rawMap[key[0]]})
		}
	}

	names := make([]string, 0, len(servers))
	for name := range servers {
		names = append(names, name)
	}
	sort.Strings(names)

	entries := make([]vibeMCPServer, 0, len(names))
	for _, name := range names {
		srv := servers[name]
		transport := string(srv.Transport)
		if transport == "" {
			transport = string(types.TransportStdio)
		}
		entries = append(entries, vibeMCPServer{
			Name:      name,
			Transport: transport,
			Command:   srv.Command,
			Args:      srv.Args,
			Env:       srv.Env,
			URL:       srv.URL,
			Headers:   srv.Headers,
		})
	}

	var buf bytes.Buffer

	// Write top-level keys in original order.
	for _, entry := range orderedTopLevel {
		single := map[string]interface{}{entry.key: entry.val}
		if err := toml.NewEncoder(&buf).Encode(single); err != nil {
			return "", fmt.Errorf("encode top-level key %q: %w", entry.key, err)
		}
	}

	// Encode [[mcp_servers]] array-of-tables using the TOML encoder for correctness.
	if len(entries) > 0 {
		if err := toml.NewEncoder(&buf).Encode(map[string]interface{}{"mcp_servers": entries}); err != nil {
			return "", fmt.Errorf("encode mcp_servers: %w", err)
		}
	}

	return buf.String(), nil
}

func (p *MistralVibeProvider) Parse(content string) (map[string]types.MCPServer, error) {
	var file struct {
		MCPServers []vibeMCPServer `toml:"mcp_servers"`
	}
	if _, err := toml.Decode(content, &file); err != nil {
		return nil, fmt.Errorf("%w: %w", ierrors.ErrMalformedConfig, err)
	}

	result := make(map[string]types.MCPServer, len(file.MCPServers))
	for _, entry := range file.MCPServers {
		result[entry.Name] = types.MCPServer{
			Transport: types.Transport(entry.Transport),
			Command:   entry.Command,
			Args:      entry.Args,
			Env:       entry.Env,
			URL:       entry.URL,
			Headers:   entry.Headers,
		}
	}
	return result, nil
}
