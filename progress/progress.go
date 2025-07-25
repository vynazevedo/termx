package progress

import (
	"fmt"
	"strings"
	"time"
	"github.com/vynazevedo/termx/renderer"
	"github.com/vynazevedo/termx/theme"
)

type Bar struct {
	total    int
	current  int
	width    int
	label    string
	showPercent bool
	char     string
	emptyChar string
}

func NewBar(total int) *Bar {
	return &Bar{
		total:    total,
		current:  0,
		width:    40,
		showPercent: true,
		char:     "█",
		emptyChar: "░",
	}
}

func (b *Bar) WithLabel(label string) *Bar {
	b.label = label
	return b
}

func (b *Bar) WithWidth(width int) *Bar {
	b.width = width
	return b
}

func (b *Bar) WithChar(char, empty string) *Bar {
	b.char = char
	b.emptyChar = empty
	return b
}

func (b *Bar) Update(current int) {
	b.current = current
	b.Render()
}

func (b *Bar) Increment() {
	b.current++
	b.Render()
}

func (b *Bar) Render() {
	r := renderer.New()
	defer r.Close()
	
	percent := float64(b.current) / float64(b.total)
	filled := int(percent * float64(b.width))
	
	r.MoveCursorUp(1)
	r.ClearLine()
	
	if b.label != "" {
		r.Write(b.label + " ")
	}
	
	r.Write("[")
	r.WriteStyled(strings.Repeat(b.char, filled), theme.Current().Success.Foreground)
	r.Write(strings.Repeat(b.emptyChar, b.width-filled))
	r.Write("]")
	
	if b.showPercent {
		r.Write(fmt.Sprintf(" %3.0f%%", percent*100))
	}
	
	r.NewLine()
}

type Spinner struct {
	frames  []string
	label   string
	delay   time.Duration
	running bool
	style   string
}

var spinnerStyles = map[string][]string{
	"dots":    {"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
	"line":    {"-", "\\", "|", "/"},
	"arrow":   {"←", "↖", "↑", "↗", "→", "↘", "↓", "↙"},
	"circle":  {"◐", "◓", "◑", "◒"},
	"box":     {"▖", "▘", "▝", "▗"},
	"bounce":  {"[    ]", "[=   ]", "[==  ]", "[=== ]", "[ ===]", "[  ==]", "[   =]"},
}

func NewSpinner() *Spinner {
	return &Spinner{
		frames: spinnerStyles["dots"],
		delay:  100 * time.Millisecond,
		style:  "dots",
	}
}

func (s *Spinner) WithStyle(style string) *Spinner {
	if frames, ok := spinnerStyles[style]; ok {
		s.frames = frames
		s.style = style
	}
	return s
}

func (s *Spinner) WithLabel(label string) *Spinner {
	s.label = label
	return s
}

func (s *Spinner) Start() {
	s.running = true
	r := renderer.New()
	
	go func() {
		i := 0
		for s.running {
			r.MoveCursorUp(1)
			r.ClearLine()
			
			r.WriteStyled(s.frames[i], theme.Current().Primary.Foreground)
			if s.label != "" {
				r.Write(" " + s.label)
			}
			r.NewLine()
			
			i = (i + 1) % len(s.frames)
			time.Sleep(s.delay)
		}
		r.Close()
	}()
}

func (s *Spinner) Stop() {
	s.running = false
	time.Sleep(s.delay)
	
	r := renderer.New()
	defer r.Close()
	r.MoveCursorUp(1)
	r.ClearLine()
}