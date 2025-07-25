package renderer

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"golang.org/x/term"
)

type Renderer struct {
	width  int
	height int
	oldState *term.State
}

func New() *Renderer {
	r := &Renderer{}
	r.updateDimensions()
	return r
}

func (r *Renderer) Init() error {
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		return fmt.Errorf("not running in a terminal")
	}

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}
	r.oldState = oldState

	r.HideCursor()
	r.Clear()
	return nil
}

func (r *Renderer) Restore() error {
	r.ShowCursor()
	if r.oldState != nil {
		return term.Restore(int(os.Stdin.Fd()), r.oldState)
	}
	return nil
}

func (r *Renderer) updateDimensions() {
	width, height, _ := term.GetSize(int(os.Stdout.Fd()))
	r.width = width
	r.height = height
}

func (r *Renderer) Width() int {
	return r.width
}

func (r *Renderer) Height() int {
	return r.height
}

func (r *Renderer) Clear() {
	fmt.Print("\033[2J\033[H")
}

func (r *Renderer) ClearLine() {
	fmt.Print("\033[2K")
}

func (r *Renderer) MoveCursor(x, y int) {
	fmt.Printf("\033[%d;%dH", y+1, x+1)
}

func (r *Renderer) HideCursor() {
	fmt.Print("\033[?25l")
}

func (r *Renderer) ShowCursor() {
	fmt.Print("\033[?25h")
}

func (r *Renderer) Print(x, y int, text string) {
	r.MoveCursor(x, y)
	fmt.Print(text)
}

func (r *Renderer) PrintCentered(y int, text string) {
	x := (r.width - len(stripANSI(text))) / 2
	if x < 0 {
		x = 0
	}
	r.Print(x, y, text)
}

func (r *Renderer) Box(x, y, width, height int, title string) {
	r.Print(x, y, "┌"+strings.Repeat("─", width-2)+"┐")
	
	if title != "" {
		titleLen := len(stripANSI(title))
		titleX := x + (width-titleLen-2)/2
		r.Print(titleX, y, "┤"+title+"├")
	}
	
	for i := 1; i < height-1; i++ {
		r.Print(x, y+i, "│"+strings.Repeat(" ", width-2)+"│")
	}
	
	r.Print(x, y+height-1, "└"+strings.Repeat("─", width-2)+"┘")
}

func stripANSI(text string) string {
	var result strings.Builder
	i := 0
	for i < len(text) {
		if text[i] == '\033' {
			for i < len(text) && text[i] != 'm' {
				i++
			}
			if i < len(text) {
				i++
			}
		} else {
			result.WriteByte(text[i])
			i++
		}
	}
	return result.String()
}

func (r *Renderer) Write(text string) {
	fmt.Print(text)
}

func (r *Renderer) WriteStyled(text, color string) {
	fmt.Print(color + text + "\033[0m")
}

func (r *Renderer) NewLine() {
	fmt.Print("\n")
}

func (r *Renderer) MoveCursorUp(lines int) {
	fmt.Printf("\033[%dA", lines)
}

func (r *Renderer) Close() {
	r.Restore()
}

func (r *Renderer) ReadKey() (Key, error) {
	event, err := ReadInput()
	if err != nil {
		return KeyUnknown, err
	}
	return event.Key, nil
}

func (r *Renderer) ClearScreen() {
	r.Clear()
}

func ClearScreen() {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	default:
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}