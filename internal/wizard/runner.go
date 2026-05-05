package wizard

type Runner interface {
	PromptText(label, placeholder string) (string, error)
	PromptSelect(label string, options []string) (string, error)
	PromptConfirm(label string) (bool, error)
	// PromptMultiSelect presents a numbered list and returns the 0-based indices
	// of the chosen items. An empty selection is valid (no items imported).
	PromptMultiSelect(label string, options []string) ([]int, error)
}
