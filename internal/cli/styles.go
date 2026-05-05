package cli

import (
	"os"

	"github.com/charmbracelet/lipgloss"
)

type Styles struct {
	Success     lipgloss.Style
	Warning     lipgloss.Style
	Err         lipgloss.Style
	Muted       lipgloss.Style
	Accent      lipgloss.Style
	Bold        lipgloss.Style
	ScopeLocal  lipgloss.Style
	ScopeGlobal lipgloss.Style
}

func NewStyles(r *lipgloss.Renderer) *Styles {
	return &Styles{
		Success:     r.NewStyle().Foreground(lipgloss.Color("2")),
		Warning:     r.NewStyle().Foreground(lipgloss.Color("3")),
		Err:         r.NewStyle().Foreground(lipgloss.Color("1")),
		Muted:       r.NewStyle().Foreground(lipgloss.Color("7")),
		Accent:      r.NewStyle().Foreground(lipgloss.Color("6")),
		Bold:        r.NewStyle().Bold(true),
		ScopeLocal:  r.NewStyle().Foreground(lipgloss.Color("6")),
		ScopeGlobal: r.NewStyle().Foreground(lipgloss.Color("3")),
	}
}

func DefaultStyles() *Styles {
	return NewStyles(lipgloss.NewRenderer(os.Stdout))
}
