package theme

import "fmt"

type Theme struct {
	Primary      Color
	Secondary    Color
	Success      Color
	Error        Color
	Warning      Color
	Info         Color
	Text         Color
	TextDim      Color
	Background   Color
	Border       Color
	Cursor       Color
	Selected     Color
	Placeholder  Color
	Muted        Color
}

type Color struct {
	Foreground string
	Background string
}

var Default = &Theme{
	Primary:     Color{Foreground: "\033[36m"},       // Cyan
	Secondary:   Color{Foreground: "\033[35m"},       // Magenta
	Success:     Color{Foreground: "\033[32m"},       // Green
	Error:       Color{Foreground: "\033[31m"},       // Red
	Warning:     Color{Foreground: "\033[33m"},       // Yellow
	Info:        Color{Foreground: "\033[34m"},       // Blue
	Text:        Color{Foreground: "\033[37m"},       // White
	TextDim:     Color{Foreground: "\033[90m"},       // Bright Black (Gray)
	Background:  Color{Background: "\033[40m"},       // Black
	Border:      Color{Foreground: "\033[90m"},       // Gray
	Cursor:      Color{Foreground: "\033[36m"},       // Cyan
	Selected:    Color{Foreground: "\033[30m", Background: "\033[46m"}, // Black on Cyan
	Placeholder: Color{Foreground: "\033[90m"},       // Gray
	Muted:       Color{Foreground: "\033[90m"},       // Gray
}

var current = Default

func Current() *Theme {
	return current
}

func Set(theme *Theme) {
	current = theme
}

func (c Color) Sprint(text string) string {
	return c.Foreground + c.Background + text + "\033[0m"
}

func (c Color) Sprintf(format string, args ...interface{}) string {
	return c.Sprint(fmt.Sprintf(format, args...))
}

func Bold(text string) string {
	return "\033[1m" + text + "\033[22m"
}

func Underline(text string) string {
	return "\033[4m" + text + "\033[24m"
}

func Italic(text string) string {
	return "\033[3m" + text + "\033[23m"
}

func Reset() string {
	return "\033[0m"
}