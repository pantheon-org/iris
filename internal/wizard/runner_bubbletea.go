package wizard

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
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

// PromptMultiSelect prints a numbered list to stderr, then reads a comma-separated
// list of 0-based indices. Enter with no input selects nothing.
func (b *TerminalRunner) PromptMultiSelect(label string, options []string) ([]int, error) {
	fmt.Fprintf(os.Stderr, "\n%s\n", label)
	for i, o := range options {
		fmt.Fprintf(os.Stderr, "  [%d] %s\n", i, o)
	}
	fmt.Fprintf(os.Stderr, "Enter numbers separated by commas (or leave empty to skip): ")

	if !b.scanner.Scan() {
		if err := b.scanner.Err(); err != nil {
			return nil, fmt.Errorf("scan multi-select: %w", err)
		}
		return nil, nil
	}
	raw := strings.TrimSpace(b.scanner.Text())
	if raw == "" {
		return nil, nil
	}

	var result []int
	for _, part := range strings.Split(raw, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		i, err := strconv.Atoi(part)
		if err != nil {
			return nil, fmt.Errorf("invalid selection %q: %w", part, err)
		}
		if i < 0 || i >= len(options) {
			return nil, fmt.Errorf("selection %d out of range [0, %d)", i, len(options))
		}
		result = append(result, i)
	}
	return result, nil
}
