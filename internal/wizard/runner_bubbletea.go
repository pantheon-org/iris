package wizard

import (
	"fmt"

	"github.com/charmbracelet/huh"
)

// TerminalRunner is the production terminal UI runner, backed by charmbracelet/huh.
type TerminalRunner struct{}

func NewTerminalRunner() *TerminalRunner {
	return &TerminalRunner{}
}

func (b *TerminalRunner) PromptText(label, placeholder string) (string, error) {
	var val string
	err := huh.NewInput().
		Title(label).
		Placeholder(placeholder).
		Value(&val).
		Run()
	if err != nil {
		return "", fmt.Errorf("input %q: %w", label, err)
	}
	return val, nil
}

func (b *TerminalRunner) PromptSelect(label string, options []string) (string, error) {
	opts := make([]huh.Option[string], len(options))
	for i, o := range options {
		opts[i] = huh.NewOption(o, o)
	}
	var val string
	err := huh.NewSelect[string]().
		Title(label).
		Options(opts...).
		Value(&val).
		Run()
	if err != nil {
		return "", fmt.Errorf("select %q: %w", label, err)
	}
	return val, nil
}

func (b *TerminalRunner) PromptConfirm(label string) (bool, error) {
	var val bool
	err := huh.NewConfirm().
		Title(label).
		Value(&val).
		Run()
	if err != nil {
		return false, fmt.Errorf("confirm %q: %w", label, err)
	}
	return val, nil
}

// PromptMultiSelect presents a checkbox list and returns the 0-based indices of chosen items.
func (b *TerminalRunner) PromptMultiSelect(label string, options []string) ([]int, error) {
	huhOpts := make([]huh.Option[int], len(options))
	for i, o := range options {
		huhOpts[i] = huh.NewOption(o, i)
	}
	var selected []int
	err := huh.NewMultiSelect[int]().
		Title(label).
		Options(huhOpts...).
		Value(&selected).
		Run()
	if err != nil {
		return nil, fmt.Errorf("multi-select %q: %w", label, err)
	}
	return selected, nil
}
