package types

import (
	"fmt"
	"net/url"

	"github.com/pantheon-org/iris/internal/ierrors"
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

// Validate checks that the MCPServer has a valid transport and all required
// fields for that transport. It returns a wrapped ErrMalformedConfig on any
// validation failure.
func (s MCPServer) Validate() error {
	switch s.Transport {
	case TransportStdio:
		if s.Command == "" {
			return fmt.Errorf("stdio server requires a command: %w", ierrors.ErrMalformedConfig)
		}
	case TransportSSE:
		if s.URL == "" {
			return fmt.Errorf("sse server requires a url: %w", ierrors.ErrMalformedConfig)
		}
		u, err := url.ParseRequestURI(s.URL)
		if err != nil || !u.IsAbs() {
			return fmt.Errorf("sse server url %q must be an absolute URI: %w", s.URL, ierrors.ErrMalformedConfig)
		}
	default:
		return fmt.Errorf("unknown transport %q: %w", s.Transport, ierrors.ErrMalformedConfig)
	}
	return nil
}
