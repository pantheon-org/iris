package wizard

import (
	"fmt"
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
