package types

type Transport string

const (
	TransportStdio Transport = "stdio"
	TransportSSE   Transport = "sse"
)

type MCPServer struct {
	Transport Transport         `json:"transport" yaml:"transport" toml:"transport"`
	Command   string            `json:"command,omitempty" yaml:"command,omitempty" toml:"command,omitempty"`
	Args      []string          `json:"args,omitempty" yaml:"args,omitempty" toml:"args,omitempty"`
	Env       map[string]string `json:"env,omitempty" yaml:"env,omitempty" toml:"env,omitempty"`
	URL       string            `json:"url,omitempty" yaml:"url,omitempty" toml:"url,omitempty"`
	Headers   map[string]string `json:"headers,omitempty" yaml:"headers,omitempty" toml:"headers,omitempty"`
	Cwd       string            `json:"cwd,omitempty" yaml:"cwd,omitempty" toml:"cwd,omitempty"`
	Enabled   *bool             `json:"enabled,omitempty" yaml:"enabled,omitempty" toml:"enabled,omitempty"`
}

type IrisConfig struct {
	Version   int                  `json:"version" yaml:"version" toml:"version"`
	Lang      string               `json:"lang,omitempty" yaml:"lang,omitempty" toml:"lang,omitempty"`
	Providers []string             `json:"providers" yaml:"providers" toml:"providers"`
	Servers   map[string]MCPServer `json:"servers" yaml:"servers" toml:"servers"`
}
