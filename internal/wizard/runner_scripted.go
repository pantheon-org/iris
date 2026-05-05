package wizard

import (
	"fmt"
	"strconv"
	"strings"
)

type ScriptedRunner struct {
	answers []string
	pos     int
}

func NewScriptedRunner(answers []string) *ScriptedRunner {
	return &ScriptedRunner{answers: answers}
}

func (s *ScriptedRunner) next(label string) (string, error) {
	if s.pos >= len(s.answers) {
		return "", fmt.Errorf("scripted runner exhausted: no answer for %q", label)
	}
	ans := s.answers[s.pos]
	s.pos++
	return ans, nil
}

func (s *ScriptedRunner) PromptText(label, _ string) (string, error) {
	return s.next(label)
}

func (s *ScriptedRunner) PromptSelect(label string, _ []string) (string, error) {
	return s.next(label)
}

func (s *ScriptedRunner) PromptConfirm(label string) (bool, error) {
	ans, err := s.next(label)
	if err != nil {
		return false, err
	}
	switch strings.ToLower(strings.TrimSpace(ans)) {
	case "true", "yes", "y", "1":
		return true, nil
	default:
		return false, nil
	}
}

// PromptMultiSelect accepts a comma-separated string of 0-based indices (e.g. "0,2")
// or the keyword "all" to select everything, or "" / "none" to select nothing.
func (s *ScriptedRunner) PromptMultiSelect(label string, options []string) ([]int, error) {
	ans, err := s.next(label)
	if err != nil {
		return nil, err
	}
	ans = strings.TrimSpace(ans)
	if ans == "" || strings.ToLower(ans) == "none" {
		return nil, nil
	}
	if strings.ToLower(ans) == "all" {
		idx := make([]int, len(options))
		for i := range options {
			idx[i] = i
		}
		return idx, nil
	}
	var result []int
	for _, part := range strings.Split(ans, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		i, err := strconv.Atoi(part)
		if err != nil {
			return nil, fmt.Errorf("invalid selection %q for %q: %w", part, label, err)
		}
		if i < 0 || i >= len(options) {
			return nil, fmt.Errorf("selection %d out of range [0, %d) for %q", i, len(options), label)
		}
		result = append(result, i)
	}
	return result, nil
}
