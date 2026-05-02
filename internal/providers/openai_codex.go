package providers

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"

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
		Name:                  "codex",
		DisplayName:           "OpenAI Codex",
		SupportsProjectConfig: true,
		GlobalConfigPath:      p.configPath,
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

	// Build ordered mcp_servers block names (sorted for determinism).
	names := make([]string, 0, len(servers))
	for name := range servers {
		names = append(names, name)
	}
	sort.Strings(names)

	var buf bytes.Buffer

	// Write top-level keys in original order.
	for _, entry := range orderedTopLevel {
		single := map[string]interface{}{entry.key: entry.val}
		if err := toml.NewEncoder(&buf).Encode(single); err != nil {
			return "", fmt.Errorf("encode top-level key %q: %w", entry.key, err)
		}
	}

	// Append mcp_servers tables in Codex's documented keyed-table format.
	for _, name := range names {
		srv := servers[name]
		fmt.Fprintf(&buf, "\n[mcp_servers.%s]\n", tomlTableKey(name))
		if srv.URL != "" {
			fmt.Fprintf(&buf, "url = %q\n", srv.URL)
		} else {
			fmt.Fprintf(&buf, "command = %q\n", srv.Command)
			if len(srv.Args) > 0 {
				buf.WriteString("args = [")
				for i, arg := range srv.Args {
					if i > 0 {
						buf.WriteString(", ")
					}
					fmt.Fprintf(&buf, "%q", arg)
				}
				buf.WriteString("]\n")
			}
			if srv.Cwd != "" {
				fmt.Fprintf(&buf, "cwd = %q\n", srv.Cwd)
			}
		}
		if len(srv.Headers) > 0 {
			buf.WriteString("http_headers = { ")
			headerKeys := sortedKeys(srv.Headers)
			for i, key := range headerKeys {
				if i > 0 {
					buf.WriteString(", ")
				}
				fmt.Fprintf(&buf, "%q = %q", key, srv.Headers[key])
			}
			buf.WriteString(" }\n")
		}
		if srv.Enabled != nil {
			fmt.Fprintf(&buf, "enabled = %t\n", *srv.Enabled)
		}
		if len(srv.Env) > 0 {
			fmt.Fprintf(&buf, "\n[mcp_servers.%s.env]\n", tomlTableKey(name))
			for _, key := range sortedKeys(srv.Env) {
				fmt.Fprintf(&buf, "%s = %q\n", tomlMapKey(key), srv.Env[key])
			}
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

func tomlTableKey(key string) string {
	if isBareTOMLKey(key) {
		return key
	}
	return fmt.Sprintf("%q", key)
}

func tomlMapKey(key string) string {
	if isBareTOMLKey(key) {
		return key
	}
	return fmt.Sprintf("%q", key)
}

func isBareTOMLKey(key string) bool {
	if key == "" {
		return false
	}
	for _, r := range key {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' || r == '-' {
			continue
		}
		return false
	}
	return !strings.Contains(key, ".")
}

func sortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
