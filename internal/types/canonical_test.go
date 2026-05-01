package types_test

import (
	"encoding/json"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/pantheon-org/iris/internal/types"
)

func fixture() types.IrisConfig {
	enabled := true
	return types.IrisConfig{
		Version:   1,
		Providers: []string{"claude", "cursor"},
		Servers: map[string]types.MCPServer{
			"my-server": {
				Transport: types.TransportStdio,
				Command:   "npx",
				Args:      []string{"-y", "@modelcontextprotocol/server-everything"},
				Env:       map[string]string{"DEBUG": "1"},
				Cwd:       "/tmp",
				Enabled:   &enabled,
			},
			"remote-server": {
				Transport: types.TransportSSE,
				URL:       "https://example.com/sse",
				Headers:   map[string]string{"Authorization": "Bearer token"},
			},
		},
	}
}

func TestIrisConfig_JSON_roundTrip(t *testing.T) {
	original := fixture()

	data, err := json.Marshal(original)
	require.NoError(t, err)

	var got types.IrisConfig
	require.NoError(t, json.Unmarshal(data, &got))

	assert.Equal(t, original, got)
}

func TestIrisConfig_YAML_roundTrip(t *testing.T) {
	original := fixture()

	data, err := yaml.Marshal(original)
	require.NoError(t, err)

	var got types.IrisConfig
	require.NoError(t, yaml.Unmarshal(data, &got))

	assert.Equal(t, original, got)
}

func TestIrisConfig_TOML_roundTrip(t *testing.T) {
	original := fixture()

	var buf []byte
	bufWriter := &tomlBuffer{}
	err := toml.NewEncoder(bufWriter).Encode(original)
	require.NoError(t, err)
	buf = bufWriter.Bytes()

	var got types.IrisConfig
	require.NoError(t, toml.Unmarshal(buf, &got))

	assert.Equal(t, original, got)
}

type tomlBuffer struct {
	data []byte
}

func (b *tomlBuffer) Write(p []byte) (int, error) {
	b.data = append(b.data, p...)
	return len(p), nil
}

func (b *tomlBuffer) Bytes() []byte {
	return b.data
}
