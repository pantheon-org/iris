package types

import "github.com/pantheon-org/iris/internal/ierrors"

type Transport string

const (
	TransportStdio Transport = "stdio"
	TransportSSE   Transport = "sse"
	// TransportHTTP represents the MCP Streamable HTTP transport (as distinct from SSE).
	// Used by providers such as Gemini CLI and Qwen Code which expose an "httpUrl" field.
	TransportHTTP Transport = "http"
)

func (t Transport) Validate() error {
	switch t {
	case TransportStdio, TransportSSE, TransportHTTP:
		return nil
	default:
		return ierrors.ErrInvalidTransport
	}
}
