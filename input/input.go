package input

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/vynazevedo/termx/renderer"
	"github.com/vynazevedo/termx/theme"
)

type Input struct {
	Label       string
	Value       *string
	Placeholder string
	Validator   func(string) error
	Mask        bool
	MaxLength   int
	
	buffer      []rune
	cursorPos   int
	renderer    *renderer.Renderer
	error       string
}

func New(label string, value *string) *Input {
	return &Input{
		Label:     label,
		Value:     value,
		buffer:    []rune{},
		cursorPos: 0,
	}
}

func (i *Input) WithPlaceholder(placeholder string) *Input {
	i.Placeholder = placeholder
	return i
}

func (i *Input) WithValidator(validator func(string) error) *Input {
	i.Validator = validator
	return i
}

func (i *Input) WithMaxLength(maxLength int) *Input {
	i.MaxLength = maxLength
	return i
}

func (i *Input) Password() *Input {
	i.Mask = true
	return i
}

func (i *Input) Run() error {
	i.renderer = renderer.New()
	if err := i.renderer.Init(); err != nil {
		return err
	}
	defer i.renderer.Restore()

	if i.Value != nil && *i.Value != "" {
		i.buffer = []rune(*i.Value)
		i.cursorPos = len(i.buffer)
	}

	for {
		i.render()
		
		event, err := renderer.ReadInput()
		if err != nil {
			return err
		}

		switch event.Key {
		case renderer.KeyCtrlC:
			return errors.New("cancelled")
		
		case renderer.KeyEnter:
			value := string(i.buffer)
			if i.Validator != nil {
				if err := i.Validator(value); err != nil {
					i.error = err.Error()
					continue
				}
			}
			if i.Value != nil {
				*i.Value = value
			}
			return nil
		
		case renderer.KeyBackspace:
			if i.cursorPos > 0 {
				i.buffer = append(i.buffer[:i.cursorPos-1], i.buffer[i.cursorPos:]...)
				i.cursorPos--
				i.error = ""
			}
		
		case renderer.KeyDelete:
			if i.cursorPos < len(i.buffer) {
				i.buffer = append(i.buffer[:i.cursorPos], i.buffer[i.cursorPos+1:]...)
				i.error = ""
			}
		
		case renderer.KeyArrowLeft:
			if i.cursorPos > 0 {
				i.cursorPos--
			}
		
		case renderer.KeyArrowRight:
			if i.cursorPos < len(i.buffer) {
				i.cursorPos++
			}
		
		case renderer.KeyHome:
			i.cursorPos = 0
		
		case renderer.KeyEnd:
			i.cursorPos = len(i.buffer)
		
		default:
			if event.Rune != 0 {
				if i.MaxLength == 0 || len(i.buffer) < i.MaxLength {
					i.buffer = append(i.buffer[:i.cursorPos], append([]rune{event.Rune}, i.buffer[i.cursorPos:]...)...)
					i.cursorPos++
					i.error = ""
				}
			}
		}
	}
}

func (i *Input) render() {
	i.renderer.Clear()
	th := theme.Current()
	
	// Calculate centered position
	labelWidth := utf8.RuneCountInString(i.Label)
	inputWidth := 50
	totalWidth := labelWidth + inputWidth + 10
	startX := (i.renderer.Width() - totalWidth) / 2
	startY := i.renderer.Height() / 2 - 2
	
	// Label
	i.renderer.Print(startX, startY, th.Primary.Sprint(i.Label))
	
	// Input box
	boxY := startY + 1
	i.renderer.Box(startX, boxY, inputWidth, 3, "")
	
	// Value or placeholder
	valueX := startX + 2
	valueY := boxY + 1
	
	displayValue := string(i.buffer)
	if i.Mask && len(i.buffer) > 0 {
		displayValue = strings.Repeat("•", len(i.buffer))
	}
	
	if len(i.buffer) == 0 && i.Placeholder != "" {
		i.renderer.Print(valueX, valueY, th.Placeholder.Sprint(i.Placeholder))
	} else {
		i.renderer.Print(valueX, valueY, displayValue)
	}
	
	// Cursor
	cursorX := valueX + i.cursorPos
	if i.Mask && i.cursorPos > 0 {
		cursorX = valueX + i.cursorPos
	}
	i.renderer.ShowCursor()
	i.renderer.MoveCursor(cursorX, valueY)
	
	// Error message
	if i.error != "" {
		errorY := boxY + 3
		i.renderer.Print(startX, errorY, th.Error.Sprint("✗ " + i.error))
	}
	
	// Help text
	helpY := i.renderer.Height() - 2
	helpText := "Enter to confirm • Ctrl+C to cancel"
	i.renderer.PrintCentered(helpY, th.TextDim.Sprint(helpText))
}

// Validators
func Required(msg string) func(string) error {
	return func(value string) error {
		if strings.TrimSpace(value) == "" {
			return fmt.Errorf(msg)
		}
		return nil
	}
}

func MinLength(min int) func(string) error {
	return func(value string) error {
		if utf8.RuneCountInString(value) < min {
			return fmt.Errorf("must be at least %d characters", min)
		}
		return nil
	}
}

func MaxLength(max int) func(string) error {
	return func(value string) error {
		if utf8.RuneCountInString(value) > max {
			return fmt.Errorf("must be at most %d characters", max)
		}
		return nil
	}
}

func Email() func(string) error {
	return func(value string) error {
		if !strings.Contains(value, "@") || !strings.Contains(value, ".") {
			return fmt.Errorf("invalid email format")
		}
		return nil
	}
}