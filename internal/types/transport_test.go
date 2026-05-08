package types_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pantheon-org/iris/internal/ierrors"
	"github.com/pantheon-org/iris/internal/types"
)

func TestTransportValidate_Stdio_ReturnsNil(t *testing.T) {
	require.NoError(t, types.TransportStdio.Validate())
}

func TestTransportValidate_SSE_ReturnsNil(t *testing.T) {
	require.NoError(t, types.TransportSSE.Validate())
}

func TestTransportValidate_HTTP_ReturnsNil(t *testing.T) {
	require.NoError(t, types.TransportHTTP.Validate())
}

func TestTransportValidate_Unknown_ReturnsErrInvalidTransport(t *testing.T) {
	err := types.Transport("grpc").Validate()
	require.Error(t, err)
	assert.True(t, errors.Is(err, ierrors.ErrInvalidTransport))
}

func TestTransportValidate_Empty_ReturnsErrInvalidTransport(t *testing.T) {
	err := types.Transport("").Validate()
	require.Error(t, err)
	assert.True(t, errors.Is(err, ierrors.ErrInvalidTransport))
}
