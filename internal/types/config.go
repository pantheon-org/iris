package types

type IrisConfig struct {
	Version   int                  `json:"version" yaml:"version" toml:"version"`
	Lang      string               `json:"lang,omitempty" yaml:"lang,omitempty" toml:"lang,omitempty"`
	Providers []string             `json:"providers,omitempty" yaml:"providers,omitempty" toml:"providers,omitempty"`
	Servers   map[string]MCPServer `json:"servers" yaml:"servers" toml:"servers"`
}

// NewIrisConfig returns an IrisConfig with version 1 and a non-nil Servers map.
func NewIrisConfig() *IrisConfig {
	return &IrisConfig{
		Version: 1,
		Servers: make(map[string]MCPServer),
	}
}
