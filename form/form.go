package form

import (
	"github.com/vynazevedo/termx/confirm"
	"github.com/vynazevedo/termx/input"
	"github.com/vynazevedo/termx/selector"
)

type Form struct {
	steps []Step
}

type Step interface {
	Run() error
}

func New() *Form {
	return &Form{
		steps: []Step{},
	}
}

func (f *Form) Input(label string, value *string) *Form {
	f.steps = append(f.steps, input.New(label, value))
	return f
}

func (f *Form) InputWithOptions(label string, value *string, opts ...func(*input.Input)) *Form {
	i := input.New(label, value)
	for _, opt := range opts {
		opt(i)
	}
	f.steps = append(f.steps, i)
	return f
}

func (f *Form) Password(label string, value *string) *Form {
	f.steps = append(f.steps, input.New(label, value).Password())
	return f
}

func (f *Form) Select(label string, options []string, selected *string) *Form {
	f.steps = append(f.steps, selector.New(label, options, selected))
	return f
}

func (f *Form) Confirm(label string, result *bool) *Form {
	f.steps = append(f.steps, confirm.New(label, result))
	return f
}

func (f *Form) Run() error {
	for _, step := range f.steps {
		if err := step.Run(); err != nil {
			return err
		}
	}
	return nil
}

// Helper functions for input options
func WithPlaceholder(placeholder string) func(*input.Input) {
	return func(i *input.Input) {
		i.WithPlaceholder(placeholder)
	}
}

func WithValidator(validator func(string) error) func(*input.Input) {
	return func(i *input.Input) {
		i.WithValidator(validator)
	}
}

func WithMaxLength(maxLength int) func(*input.Input) {
	return func(i *input.Input) {
		i.WithMaxLength(maxLength)
	}
}