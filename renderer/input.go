package renderer

import (
	"fmt"
	"os"
)

type Key int

const (
	KeyUnknown Key = iota
	KeyArrowUp
	KeyArrowDown
	KeyArrowLeft
	KeyArrowRight
	KeyEnter
	KeyTab
	KeyBackspace
	KeyEscape
	KeySpace
	KeyCtrlC
	KeyCtrlD
	KeyHome
	KeyEnd
	KeyPageUp
	KeyPageDown
	KeyDelete
)

type InputEvent struct {
	Key  Key
	Rune rune
	Raw  []byte
}

func ReadInput() (*InputEvent, error) {
	buf := make([]byte, 256)
	n, err := os.Stdin.Read(buf)
	if err != nil {
		return nil, err
	}

	event := &InputEvent{
		Raw: buf[:n],
	}

	if n == 1 {
		switch buf[0] {
		case '\r', '\n':
			event.Key = KeyEnter
		case '\t':
			event.Key = KeyTab
		case 127, '\b':
			event.Key = KeyBackspace
		case 27:
			event.Key = KeyEscape
		case ' ':
			event.Key = KeySpace
		case 3:
			event.Key = KeyCtrlC
		case 4:
			event.Key = KeyCtrlD
		default:
			if buf[0] >= 32 && buf[0] < 127 {
				event.Rune = rune(buf[0])
			}
		}
	} else if n > 2 && buf[0] == 27 && buf[1] == '[' {
		switch buf[2] {
		case 'A':
			event.Key = KeyArrowUp
		case 'B':
			event.Key = KeyArrowDown
		case 'C':
			event.Key = KeyArrowRight
		case 'D':
			event.Key = KeyArrowLeft
		case 'H':
			event.Key = KeyHome
		case 'F':
			event.Key = KeyEnd
		case '3':
			if n > 3 && buf[3] == '~' {
				event.Key = KeyDelete
			}
		case '5':
			if n > 3 && buf[3] == '~' {
				event.Key = KeyPageUp
			}
		case '6':
			if n > 3 && buf[3] == '~' {
				event.Key = KeyPageDown
			}
		}
	}

	return event, nil
}

func (k Key) String() string {
	switch k {
	case KeyArrowUp:
		return "↑"
	case KeyArrowDown:
		return "↓"
	case KeyArrowLeft:
		return "←"
	case KeyArrowRight:
		return "→"
	case KeyEnter:
		return "Enter"
	case KeyTab:
		return "Tab"
	case KeyBackspace:
		return "Backspace"
	case KeyEscape:
		return "Esc"
	case KeySpace:
		return "Space"
	case KeyCtrlC:
		return "Ctrl+C"
	case KeyCtrlD:
		return "Ctrl+D"
	case KeyHome:
		return "Home"
	case KeyEnd:
		return "End"
	case KeyPageUp:
		return "PgUp"
	case KeyPageDown:
		return "PgDn"
	case KeyDelete:
		return "Delete"
	default:
		return fmt.Sprintf("Key(%d)", k)
	}
}