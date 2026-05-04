package types

import (
	"fmt"
	"net/url"

	"github.com/pantheon-org/iris/internal/ierrors"
)

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
