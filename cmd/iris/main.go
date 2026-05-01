package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/pantheon-org/iris/internal/version"
)

func main() {
	root := &cobra.Command{
		Use:     "iris",
		Short:   "Manage MCP server configs across AI providers",
		Version: version.Version,
	}
	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
