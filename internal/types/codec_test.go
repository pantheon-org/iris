package types_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pantheon-org/iris/internal/ierrors"
	"github.com/pantheon-org/iris/internal/types"
)

type codecFixture struct {
	Name string `json:"name" yaml:"name" toml:"name"`
	Age  int    `json:"age" yaml:"age" toml:"age"`
}

func TestCodecForExtension_json_returnsCodec(t *testing.T) {
	c, err := types.CodecForExtension(".json")
	require.NoError(t, err)
	assert.Equal(t, ".json", c.Extension())
}

func TestCodecForExtension_yaml_returnsCodec(t *testing.T) {
	c, err := types.CodecForExtension(".yaml")
	require.NoError(t, err)
	assert.Equal(t, ".yaml", c.Extension())
}

func TestCodecForExtension_yml_returnsYamlCodec(t *testing.T) {
	c, err := types.CodecForExtension(".yml")
	require.NoError(t, err)
	assert.Equal(t, ".yaml", c.Extension())
}

func TestCodecForExtension_toml_returnsCodec(t *testing.T) {
	c, err := types.CodecForExtension(".toml")
	require.NoError(t, err)
	assert.Equal(t, ".toml", c.Extension())
}

func TestCodecForExtension_unknown_returnsUnsupportedFormatError(t *testing.T) {
	_, err := types.CodecForExtension(".xml")
	require.Error(t, err)
	assert.True(t, errors.Is(err, ierrors.ErrUnsupportedFormat))
}

func TestCodecForExtension_empty_returnsUnsupportedFormatError(t *testing.T) {
	_, err := types.CodecForExtension("")
	require.Error(t, err)
	assert.True(t, errors.Is(err, ierrors.ErrUnsupportedFormat))
}
