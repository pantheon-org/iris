package types_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pantheon-org/iris/internal/ierrors"
	"github.com/pantheon-org/iris/internal/types"
)

func TestMCPServer_Validate_stdio_valid(t *testing.T) {
	s := types.MCPServer{
		Transport: types.TransportStdio,
		Command:   "npx",
		Args:      []string{"-y", "@modelcontextprotocol/server-everything"},
	}
	require.NoError(t, s.Validate())
}

func TestMCPServer_Validate_stdio_missingCommand_returnsError(t *testing.T) {
	s := types.MCPServer{
		Transport: types.TransportStdio,
	}
	err := s.Validate()
	require.Error(t, err)
	assert.True(t, errors.Is(err, ierrors.ErrMalformedConfig))
}

func TestMCPServer_Validate_sse_valid(t *testing.T) {
	s := types.MCPServer{
		Transport: types.TransportSSE,
		URL:       "https://example.com/sse",
	}
	require.NoError(t, s.Validate())
}

func TestMCPServer_Validate_sse_missingURL_returnsError(t *testing.T) {
	s := types.MCPServer{
		Transport: types.TransportSSE,
	}
	err := s.Validate()
	require.Error(t, err)
	assert.True(t, errors.Is(err, ierrors.ErrMalformedConfig))
}

func TestMCPServer_Validate_sse_relativeURL_returnsError(t *testing.T) {
	s := types.MCPServer{
		Transport: types.TransportSSE,
		URL:       "/relative/path",
	}
	err := s.Validate()
	require.Error(t, err)
	assert.True(t, errors.Is(err, ierrors.ErrMalformedConfig))
}

func TestMCPServer_Validate_sse_malformedURL_returnsError(t *testing.T) {
	s := types.MCPServer{
		Transport: types.TransportSSE,
		URL:       "://not-a-url",
	}
	err := s.Validate()
	require.Error(t, err)
	assert.True(t, errors.Is(err, ierrors.ErrMalformedConfig))
}

func TestMCPServer_Validate_unknownTransport_returnsError(t *testing.T) {
	s := types.MCPServer{
		Transport: "websocket",
		Command:   "cmd",
	}
	err := s.Validate()
	require.Error(t, err)
	assert.True(t, errors.Is(err, ierrors.ErrMalformedConfig))
}

func TestMCPServer_Validate_emptyTransport_returnsError(t *testing.T) {
	s := types.MCPServer{}
	err := s.Validate()
	require.Error(t, err)
	assert.True(t, errors.Is(err, ierrors.ErrMalformedConfig))
}
