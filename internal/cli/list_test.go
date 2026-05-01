package cli_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pantheon-org/iris/internal/cli"
	"github.com/pantheon-org/iris/internal/types"
)

func TestRunList_noServers_printsMessage(t *testing.T) {
	cfg := &types.IrisConfig{Servers: map[string]types.MCPServer{}}
	var buf bytes.Buffer

	err := cli.RunList(cfg, &buf)

	require.NoError(t, err)
	assert.Equal(t, "No servers configured.\n", buf.String())
}

func TestRunList_singleServer_correctFormat(t *testing.T) {
	cfg := &types.IrisConfig{
		Servers: map[string]types.MCPServer{
			"fetch": {Transport: types.TransportStdio, Command: "uvx", Args: []string{"mcp-server-fetch"}},
		},
	}
	var buf bytes.Buffer

	err := cli.RunList(cfg, &buf)

	require.NoError(t, err)
	out := buf.String()
	assert.Contains(t, out, "Servers (1):")
	assert.Contains(t, out, "fetch")
	assert.Contains(t, out, "stdio")
	assert.Contains(t, out, "uvx mcp-server-fetch")
}

func TestRunList_multipleServers_sortedAlphabetically(t *testing.T) {
	cfg := &types.IrisConfig{
		Servers: map[string]types.MCPServer{
			"my-server":  {Transport: types.TransportStdio, Command: "npx", Args: []string{"-y", "@mcp/server"}},
			"fetch":      {Transport: types.TransportStdio, Command: "uvx", Args: []string{"mcp-server-fetch"}},
			"filesystem": {Transport: types.TransportStdio, Command: "npx", Args: []string{"-y", "@mcp/filesystem"}},
		},
	}
	var buf bytes.Buffer

	err := cli.RunList(cfg, &buf)

	require.NoError(t, err)
	out := buf.String()
	assert.Contains(t, out, "Servers (3):")

	fetchPos := bytes.Index(buf.Bytes(), []byte("fetch"))
	fsPos := bytes.Index(buf.Bytes(), []byte("filesystem"))
	myPos := bytes.Index(buf.Bytes(), []byte("my-server"))
	assert.Less(t, fetchPos, fsPos, "fetch should come before filesystem")
	assert.Less(t, fsPos, myPos, "filesystem should come before my-server")
}
