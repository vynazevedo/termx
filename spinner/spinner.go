package spinner

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/vynazevedo/termx/renderer"
	"github.com/vynazevedo/termx/theme"
)

// SpinnerStyle defines different spinner animations
type SpinnerStyle string

const (
	Dots     SpinnerStyle = "dots"
	Line     SpinnerStyle = "line"
	Circle   SpinnerStyle = "circle"
	Arrow    SpinnerStyle = "arrow"
	Clock    SpinnerStyle = "clock"
	Bounce   SpinnerStyle = "bounce"
	Pulse    SpinnerStyle = "pulse"
	Growing  SpinnerStyle = "growing"
)

// Spinner represents a loading spinner component
type Spinner struct {
	style    SpinnerStyle
	label    string
	color    string
	speed    time.Duration
	active   bool
	done     chan bool
	renderer *renderer.Renderer
	mu       sync.RWMutex
	theme    *theme.Theme
}

// spinnerFrames contains animation frames for each style
var spinnerFrames = map[SpinnerStyle][]string{
	Dots:    {"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
	Line:    {"|", "/", "-", "\\"},
	Circle:  {"◐", "◓", "◑", "◒"},
	Arrow:   {"←", "↖", "↑", "↗", "→", "↘", "↓", "↙"},
	Clock:   {"🕐", "🕑", "🕒", "🕓", "🕔", "🕕", "🕖", "🕗", "🕘", "🕙", "🕚", "🕛"},
	Bounce:  {"⠁", "⠂", "⠄", "⠂"},
	Pulse:   {"●", "○", "●", "○"},
	Growing: {"▁", "▃", "▄", "▅", "▆", "▇", "█", "▇", "▆", "▅", "▄", "▃"},
}

// New creates a new spinner instance
func New() *Spinner {
	return &Spinner{
		style:    Dots,
		label:    "Carregando...",
		color:    "\033[36m", // cyan
		speed:    100 * time.Millisecond,
		done:     make(chan bool),
		renderer: renderer.New(),
		theme:    theme.Default,
	}
}

// WithStyle sets the spinner animation style
func (s *Spinner) WithStyle(style SpinnerStyle) *Spinner {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.style = style
	return s
}

// WithLabel sets the spinner label text
func (s *Spinner) WithLabel(label string) *Spinner {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.label = label
	return s
}

// WithColor sets the spinner color
func (s *Spinner) WithColor(color string) *Spinner {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.color = color
	return s
}

// WithSpeed sets the animation speed
func (s *Spinner) WithSpeed(speed time.Duration) *Spinner {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.speed = speed
	return s
}

// WithTheme sets the spinner theme
func (s *Spinner) WithTheme(t *theme.Theme) *Spinner {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.theme = t
	s.color = t.Primary.Foreground
	return s
}

// Start begins the spinner animation
func (s *Spinner) Start() {
	s.mu.Lock()
	if s.active {
		s.mu.Unlock()
		return
	}
	s.active = true
	s.mu.Unlock()

	go s.animate()
}

// Stop stops the spinner animation
func (s *Spinner) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if !s.active {
		return
	}
	
	s.active = false
	s.done <- true
	
	// Clear the spinner line
	fmt.Print("\r" + strings.Repeat(" ", len(s.label)+10) + "\r")
}

// StopWithMessage stops the spinner and shows a completion message
func (s *Spinner) StopWithMessage(message string) {
	s.Stop()
	fmt.Printf("✓ %s%s\n", s.theme.Success.Sprint(message), theme.Reset())
}

// StopWithError stops the spinner and shows an error message
func (s *Spinner) StopWithError(message string) {
	s.Stop()
	fmt.Printf("✗ %s%s\n", s.theme.Error.Sprint(message), theme.Reset())
}

// IsActive returns whether the spinner is currently running
func (s *Spinner) IsActive() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.active
}

// animate runs the spinner animation loop
func (s *Spinner) animate() {
	frames := spinnerFrames[s.style]
	frameIndex := 0
	ticker := time.NewTicker(s.speed)
	defer ticker.Stop()

	for {
		select {
		case <-s.done:
			return
		case <-ticker.C:
			s.mu.RLock()
			if !s.active {
				s.mu.RUnlock()
				return
			}
			
			frame := frames[frameIndex%len(frames)]
			output := fmt.Sprintf("\r%s%s%s %s", s.color, frame, theme.Reset(), s.label)
			fmt.Print(output)
			
			frameIndex++
			s.mu.RUnlock()
		}
	}
}

// Render renders the spinner at current state (for static display)
func (s *Spinner) Render() {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	frames := spinnerFrames[s.style]
	frame := frames[0] // Use first frame for static display
	output := fmt.Sprintf("%s%s%s %s", s.color, frame, theme.Reset(), s.label)
	fmt.Print(output)
}

// Close cleans up the spinner resources
func (s *Spinner) Close() {
	s.Stop()
	if s.renderer != nil {
		s.renderer.Close()
	}
}

// Predefined spinner configurations
func SpinnerDots() *Spinner {
	return New().WithStyle(Dots).WithLabel("Processando...")
}

func Loading() *Spinner {
	return New().WithStyle(Dots).WithLabel("Carregando...")
}

func Processing() *Spinner {
	return New().WithStyle(Line).WithLabel("Processando...")
}

func Downloading() *Spinner {
	return New().WithStyle(Growing).WithLabel("Baixando...")
}

func Installing() *Spinner {
	return New().WithStyle(Circle).WithLabel("Instalando...")
}

func Connecting() *Spinner {
	return New().WithStyle(Pulse).WithLabel("Conectando...")
}