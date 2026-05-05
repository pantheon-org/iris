package cli_test

import (
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"

	"github.com/pantheon-org/iris/internal/cli"
)

func noColourStyles() *cli.Styles {
	r := lipgloss.NewRenderer(os.Stdout)
	r.SetColorProfile(termenv.Ascii)
	return cli.NewStyles(r)
}
