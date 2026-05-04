package config_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pantheon-org/iris/internal/config"
	"github.com/pantheon-org/iris/internal/ierrors"
)

type codecFixture struct {
	Name string `json:"name" yaml:"name" toml:"name"`
	Age  int    `json:"age" yaml:"age" toml:"age"`
}

func TestCodecForExtension_json_returnsCodec(t *testing.T) {
	c, err := config.CodecForExtension(".json")
	require.NoError(t, err)
	assert.Equal(t, ".json", c.Extension())
}

func TestCodecForExtension_yaml_returnsCodec(t *testing.T) {
	c, err := config.CodecForExtension(".yaml")
	require.NoError(t, err)
	assert.Equal(t, ".yaml", c.Extension())
}

func TestCodecForExtension_yml_returnsYamlCodec(t *testing.T) {
	c, err := config.CodecForExtension(".yml")
	require.NoError(t, err)
	assert.Equal(t, ".yaml", c.Extension())
}

func TestCodecForExtension_toml_returnsCodec(t *testing.T) {
	c, err := config.CodecForExtension(".toml")
	require.NoError(t, err)
	assert.Equal(t, ".toml", c.Extension())
}

func TestCodecForExtension_unknown_returnsUnsupportedFormatError(t *testing.T) {
	_, err := config.CodecForExtension(".xml")
	require.Error(t, err)
	assert.True(t, errors.Is(err, ierrors.ErrUnsupportedFormat))
}

func TestCodecForExtension_empty_returnsUnsupportedFormatError(t *testing.T) {
	_, err := config.CodecForExtension("")
	require.Error(t, err)
	assert.True(t, errors.Is(err, ierrors.ErrUnsupportedFormat))
}

func TestJsonCodec_roundTrip_equal(t *testing.T) {
	c, err := config.CodecForExtension(".json")
	require.NoError(t, err)

	original := codecFixture{Name: "alice", Age: 30}
	data, err := c.Marshal(original)
	require.NoError(t, err)

	var got codecFixture
	require.NoError(t, c.Unmarshal(data, &got))
	assert.Equal(t, original, got)
}

func TestYamlCodec_roundTrip_equal(t *testing.T) {
	c, err := config.CodecForExtension(".yaml")
	require.NoError(t, err)

	original := codecFixture{Name: "bob", Age: 25}
	data, err := c.Marshal(original)
	require.NoError(t, err)

	var got codecFixture
	require.NoError(t, c.Unmarshal(data, &got))
	assert.Equal(t, original, got)
}

func TestTomlCodec_roundTrip_equal(t *testing.T) {
	c, err := config.CodecForExtension(".toml")
	require.NoError(t, err)

	original := codecFixture{Name: "carol", Age: 40}
	data, err := c.Marshal(original)
	require.NoError(t, err)

	var got codecFixture
	require.NoError(t, c.Unmarshal(data, &got))
	assert.Equal(t, original, got)
}

func TestJsonCodec_unmarshalBadData_returnsError(t *testing.T) {
	c, err := config.CodecForExtension(".json")
	require.NoError(t, err)

	var got codecFixture
	err = c.Unmarshal([]byte(`{bad json`), &got)
	assert.Error(t, err)
}

func TestYamlCodec_unmarshalBadData_returnsError(t *testing.T) {
	c, err := config.CodecForExtension(".yaml")
	require.NoError(t, err)

	var got codecFixture
	err = c.Unmarshal([]byte(":\tbad: yaml:"), &got)
	assert.Error(t, err)
}
