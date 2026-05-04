package wizard

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// TerminalRunner is the production terminal UI runner.
type TerminalRunner struct {
	scanner *bufio.Scanner
}

func NewTerminalRunner() *TerminalRunner {
	return &TerminalRunner{scanner: bufio.NewScanner(os.Stdin)}
}

func (b *TerminalRunner) prompt(label, hint string) (string, error) {
	if hint != "" {
		fmt.Fprintf(os.Stderr, "%s [%s]: ", label, hint)
	} else {
		fmt.Fprintf(os.Stderr, "%s: ", label)
	}
	if !b.scanner.Scan() {
		if err := b.scanner.Err(); err != nil {
			return "", fmt.Errorf("scan %q: %w", label, err)
		}
		return "", fmt.Errorf("unexpected EOF reading %q", label)
	}
	return strings.TrimSpace(b.scanner.Text()), nil
}

func (b *TerminalRunner) PromptText(label, placeholder string) (string, error) {
	return b.prompt(label, placeholder)
}

func (b *TerminalRunner) PromptSelect(label string, options []string) (string, error) {
	hint := strings.Join(options, "/")
	return b.prompt(label, hint)
}

func (b *TerminalRunner) PromptConfirm(label string) (bool, error) {
	ans, err := b.prompt(label, "yes/no")
	if err != nil {
		return false, err
	}
	switch strings.ToLower(ans) {
	case "true", "yes", "y", "1":
		return true, nil
	default:
		return false, nil
	}
}
