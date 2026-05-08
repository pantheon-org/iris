package types_test

import (
	"testing"

	"github.com/pantheon-org/iris/internal/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTomlCodec_marshall(t *testing.T) {
	c, err := types.CodecForExtension(".toml")
	require.NoError(t, err)

	original := codecFixture{Name: "alice", Age: 30}
	data, err := c.Marshal(original)
	require.NoError(t, err)

	var unmarshalled codecFixture
	err = c.Unmarshal(data, &unmarshalled)
	require.NoError(t, err)

	assert.Equal(t, original, unmarshalled)
}

func TestTomlCodec_unmarshall_invalid(t *testing.T) {
	c, err := types.CodecForExtension(".toml")
	require.NoError(t, err)

	invalidData := []byte("invalid = [unclosed list")
	var unmarshalled codecFixture
	err = c.Unmarshal(invalidData, &unmarshalled)
	require.Error(t, err)
}
