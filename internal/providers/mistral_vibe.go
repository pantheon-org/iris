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

func newMistralVibeProviderWithPath(path string) *MistralVibeProvider {
	return &MistralVibeProvider{configPath: path}
}

// NewMistralVibeProviderWithPath creates a MistralVibeProvider pinned to a fixed config path.
// Intended for use in tests.
func NewMistralVibeProviderWithPath(path string) *MistralVibeProvider {
	return &MistralVibeProvider{configPath: path, pinned: true}
}

func (p *MistralVibeProvider) Config() ProviderConfig {
	return ProviderConfig{
		Name:                  NameMistralVibe,
		DisplayName:           "Mistral Vibe",
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

func (p *MistralVibeProvider) Exists(projectRoot string) bool {
	_, err := os.Stat(p.ConfigFilePath(projectRoot))
	return err == nil
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
	for _, entry := range orderedTopLevel {
		single := map[string]interface{}{entry.key: entry.val}
		if err := toml.NewEncoder(&buf).Encode(single); err != nil {
			return "", fmt.Errorf("encode top-level key %q: %w", entry.key, err)
		}
	}

	for _, entry := range entries {
		buf.WriteString("\n[[mcp_servers]]\n")
		fmt.Fprintf(&buf, "name = %q\n", entry.Name)
		fmt.Fprintf(&buf, "transport = %q\n", entry.Transport)
		if entry.Command != "" {
			fmt.Fprintf(&buf, "command = %q\n", entry.Command)
		}
		if len(entry.Args) > 0 {
			buf.WriteString("args = [")
			for i, a := range entry.Args {
				if i > 0 {
					buf.WriteString(", ")
				}
				fmt.Fprintf(&buf, "%q", a)
			}
			buf.WriteString("]\n")
		}
		if len(entry.Env) > 0 {
			buf.WriteString("env = {")
			keys := make([]string, 0, len(entry.Env))
			for k := range entry.Env {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for i, k := range keys {
				if i > 0 {
					buf.WriteString(", ")
				}
				fmt.Fprintf(&buf, "%s = %q", k, entry.Env[k])
			}
			buf.WriteString("}\n")
		}
		if entry.URL != "" {
			fmt.Fprintf(&buf, "url = %q\n", entry.URL)
		}
		if len(entry.Headers) > 0 {
			buf.WriteString("headers = {")
			keys := make([]string, 0, len(entry.Headers))
			for k := range entry.Headers {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for i, k := range keys {
				if i > 0 {
					buf.WriteString(", ")
				}
				fmt.Fprintf(&buf, "%s = %q", k, entry.Headers[k])
			}
			buf.WriteString("}\n")
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
