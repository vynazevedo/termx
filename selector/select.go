package selector

import (
	"errors"
	"fmt"
	"strings"

	"github.com/vynazevedo/termx/renderer"
	"github.com/vynazevedo/termx/theme"
)

type Select struct {
	Label    string
	Options  []string
	Selected *string
	
	currentIndex int
	renderer     *renderer.Renderer
	filter       string
	filtered     []int
}

func New(label string, options []string, selected *string) *Select {
	s := &Select{
		Label:    label,
		Options:  options,
		Selected: selected,
	}
	
	// Set initial index if value exists
	if selected != nil && *selected != "" {
		for i, opt := range options {
			if opt == *selected {
				s.currentIndex = i
				break
			}
		}
	}
	
	s.updateFiltered()
	return s
}

func (s *Select) updateFiltered() {
	s.filtered = []int{}
	filter := strings.ToLower(s.filter)
	
	for i, opt := range s.Options {
		if filter == "" || strings.Contains(strings.ToLower(opt), filter) {
			s.filtered = append(s.filtered, i)
		}
	}
	
	// Reset current index if out of bounds
	if len(s.filtered) > 0 {
		found := false
		for idx, optIdx := range s.filtered {
			if optIdx == s.currentIndex {
				s.currentIndex = idx
				found = true
				break
			}
		}
		if !found {
			s.currentIndex = 0
		}
	}
}

func (s *Select) Run() error {
	s.renderer = renderer.New()
	if err := s.renderer.Init(); err != nil {
		return err
	}
	defer s.renderer.Restore()

	for {
		s.render()
		
		event, err := renderer.ReadInput()
		if err != nil {
			return err
		}

		switch event.Key {
		case renderer.KeyCtrlC:
			return errors.New("cancelled")
		
		case renderer.KeyEnter:
			if len(s.filtered) > 0 && s.currentIndex < len(s.filtered) {
				selectedOption := s.Options[s.filtered[s.currentIndex]]
				if s.Selected != nil {
					*s.Selected = selectedOption
				}
				return nil
			}
		
		case renderer.KeyArrowUp:
			if s.currentIndex > 0 {
				s.currentIndex--
			}
		
		case renderer.KeyArrowDown:
			if s.currentIndex < len(s.filtered)-1 {
				s.currentIndex++
			}
		
		case renderer.KeyBackspace:
			if len(s.filter) > 0 {
				s.filter = s.filter[:len(s.filter)-1]
				s.updateFiltered()
			}
		
		case renderer.KeyEscape:
			if s.filter != "" {
				s.filter = ""
				s.updateFiltered()
			}
		
		default:
			if event.Rune != 0 {
				s.filter += string(event.Rune)
				s.updateFiltered()
			}
		}
	}
}

func (s *Select) render() {
	s.renderer.Clear()
	th := theme.Current()
	
	// Calculate dimensions
	maxOptionLen := 0
	for _, opt := range s.Options {
		if len(opt) > maxOptionLen {
			maxOptionLen = len(opt)
		}
	}
	
	boxWidth := maxOptionLen + 10
	if boxWidth < 40 {
		boxWidth = 40
	}
	
	visibleItems := 7
	boxHeight := visibleItems + 4
	
	startX := (s.renderer.Width() - boxWidth) / 2
	startY := (s.renderer.Height() - boxHeight) / 2
	
	// Label
	s.renderer.PrintCentered(startY-2, th.Primary.Sprint(s.Label))
	
	// Main box
	s.renderer.Box(startX, startY, boxWidth, boxHeight, "")
	
	// Filter display
	filterY := startY + 1
	filterX := startX + 2
	
	if s.filter != "" {
		s.renderer.Print(filterX, filterY, th.Info.Sprint("Filter: ") + s.filter)
	} else {
		s.renderer.Print(filterX, filterY, th.TextDim.Sprint("Type to filter..."))
	}
	
	// Options
	optionsY := startY + 3
	
	// Calculate visible range
	startIdx := 0
	if s.currentIndex >= visibleItems {
		startIdx = s.currentIndex - visibleItems + 1
	}
	
	endIdx := startIdx + visibleItems
	if endIdx > len(s.filtered) {
		endIdx = len(s.filtered)
		if endIdx-startIdx < visibleItems && startIdx > 0 {
			startIdx = endIdx - visibleItems
			if startIdx < 0 {
				startIdx = 0
			}
		}
	}
	
	// Render visible options
	for i := startIdx; i < endIdx; i++ {
		y := optionsY + (i - startIdx)
		optionIdx := s.filtered[i]
		option := s.Options[optionIdx]
		
		if i == s.currentIndex {
			// Selected item
			s.renderer.Print(startX+1, y, strings.Repeat(" ", boxWidth-2))
			s.renderer.Print(startX+2, y, th.Selected.Sprint(fmt.Sprintf("▶ %s", option)))
		} else {
			s.renderer.Print(startX+2, y, fmt.Sprintf("  %s", option))
		}
	}
	
	// Scrollbar indicator
	if len(s.filtered) > visibleItems {
		scrollY := optionsY
		scrollHeight := visibleItems
		scrollPos := int(float64(s.currentIndex) / float64(len(s.filtered)-1) * float64(scrollHeight-1))
		
		for i := 0; i < scrollHeight; i++ {
			x := startX + boxWidth - 3
			if i == scrollPos {
				s.renderer.Print(x, scrollY+i, th.Primary.Sprint("█"))
			} else {
				s.renderer.Print(x, scrollY+i, th.Border.Sprint("│"))
			}
		}
	}
	
	// Status line
	statusY := startY + boxHeight - 2
	status := fmt.Sprintf("%d/%d items", len(s.filtered), len(s.Options))
	s.renderer.Print(startX+2, statusY, th.TextDim.Sprint(status))
	
	// Help text
	helpY := s.renderer.Height() - 2
	helpText := "↑↓ Navigate • Enter Select • Type to filter • Esc Clear filter • Ctrl+C Cancel"
	s.renderer.PrintCentered(helpY, th.TextDim.Sprint(helpText))
}