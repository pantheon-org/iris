package providers

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/BurntSushi/toml"

	"github.com/pantheon-org/iris/internal/ierrors"
	"github.com/pantheon-org/iris/internal/types"
)

type CodexProvider struct {
	configPath string
}

func codexGlobalPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".codex", "config.toml")
}

func NewCodexProvider() *CodexProvider {
	return &CodexProvider{configPath: codexGlobalPath()}
}

// NewCodexProviderWithPath creates a CodexProvider using a custom config path.
// Intended for use in tests.
func NewCodexProviderWithPath(path string) *CodexProvider {
	return &CodexProvider{configPath: path}
}

func (p *CodexProvider) Config() ProviderConfig {
	return ProviderConfig{
		Name:                  "codex",
		DisplayName:           "OpenAI Codex",
		SupportsProjectConfig: false,
		GlobalConfigPath:      p.configPath,
	}
}

func (p *CodexProvider) ConfigFilePath(_ string) string {
	return p.configPath
}

func (p *CodexProvider) Exists(projectRoot string) bool {
	_, err := os.Stat(p.ConfigFilePath(projectRoot))
	return err == nil
}

type codexMCPServer struct {
	Name    string            `toml:"name"`
	Command string            `toml:"command"`
	Args    []string          `toml:"args,omitempty"`
	Env     map[string]string `toml:"env,omitempty"`
	Type    string            `toml:"type"`
}

// codexConfig is used only for parsing existing TOML to extract non-mcp_servers keys.
// We use map[string]interface{} for the "other" keys approach via a raw decode.
type codexFileRaw struct {
	MCPServers []codexMCPServer       `toml:"mcp_servers"`
	Other      map[string]interface{} `toml:"-"`
}

func (p *CodexProvider) Generate(servers map[string]types.MCPServer, existingContent string) (string, error) {
	// Collect top-level non-mcp_servers keys preserving original order via toml.MetaData.
	type topLevelEntry struct {
		key string
		val interface{}
	}
	var orderedTopLevel []topLevelEntry

	if existingContent != "" {
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

	// Build ordered mcp_servers slice (sorted by name for determinism).
	names := make([]string, 0, len(servers))
	for name := range servers {
		names = append(names, name)
	}
	sort.Strings(names)

	entries := make([]codexMCPServer, 0, len(names))
	for _, name := range names {
		srv := servers[name]
		t := string(srv.Transport)
		if t == "" {
			t = string(types.TransportStdio)
		}
		entries = append(entries, codexMCPServer{
			Name:    name,
			Command: srv.Command,
			Args:    srv.Args,
			Env:     srv.Env,
			Type:    t,
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

	// Append mcp_servers array table entries using the struct encoder for consistent indentation.
	for _, entry := range entries {
		buf.WriteString("\n[[mcp_servers]]\n")
		enc := toml.NewEncoder(&buf)
		enc.Indent = "  "
		if err := enc.Encode(entry); err != nil {
			return "", fmt.Errorf("encode mcp_servers entry: %w", err)
		}
	}

	return buf.String(), nil
}

func (p *CodexProvider) Parse(content string) (map[string]types.MCPServer, error) {
	var file struct {
		MCPServers []codexMCPServer `toml:"mcp_servers"`
	}
	if _, err := toml.Decode(content, &file); err != nil {
		return nil, fmt.Errorf("%w: %w", ierrors.ErrMalformedConfig, err)
	}

	result := make(map[string]types.MCPServer, len(file.MCPServers))
	for _, entry := range file.MCPServers {
		result[entry.Name] = types.MCPServer{
			Command:   entry.Command,
			Args:      entry.Args,
			Env:       entry.Env,
			Transport: types.Transport(entry.Type),
		}
	}
	return result, nil
}
