package wizard

type Runner interface {
	PromptText(label, placeholder string) (string, error)
	PromptSelect(label string, options []string) (string, error)
	PromptConfirm(label string) (bool, error)
}
