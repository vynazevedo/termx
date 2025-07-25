package layout

import (
	"strings"
	"github.com/vynazevedo/termx/renderer"
	"golang.org/x/term"
)

type Layout interface {
	Render(r *renderer.Renderer)
	SetContent(content string)
}

type Split struct {
	direction string
	ratio     float64
	left      Layout
	right     Layout
	width     int
	height    int
}

func NewSplit(direction string) *Split {
	w, h, _ := term.GetSize(0)
	return &Split{
		direction: direction,
		ratio:     0.5,
		width:     w,
		height:    h,
	}
}

func (s *Split) WithRatio(ratio float64) *Split {
	s.ratio = ratio
	return s
}

func (s *Split) SetLeft(l Layout) *Split {
	s.left = l
	return s
}

func (s *Split) SetRight(l Layout) *Split {
	s.right = l
	return s
}

func (s *Split) Render(r *renderer.Renderer) {
	if s.direction == "horizontal" {
		splitCol := int(float64(s.width) * s.ratio)
		
		if s.left != nil {
			s.renderPane(r, s.left, 0, 0, splitCol, s.height)
		}
		
		for i := 0; i < s.height; i++ {
			r.MoveCursor(splitCol, i)
			r.Write("│")
		}
		
		if s.right != nil {
			s.renderPane(r, s.right, splitCol+1, 0, s.width-splitCol-1, s.height)
		}
	} else {
		splitRow := int(float64(s.height) * s.ratio)
		
		if s.left != nil {
			s.renderPane(r, s.left, 0, 0, s.width, splitRow)
		}
		
		r.MoveCursor(0, splitRow)
		r.Write(strings.Repeat("─", s.width))
		
		if s.right != nil {
			s.renderPane(r, s.right, 0, splitRow+1, s.width, s.height-splitRow-1)
		}
	}
}

func (s *Split) renderPane(r *renderer.Renderer, layout Layout, x, y, w, h int) {
	r.MoveCursor(x, y)
	layout.Render(r)
}

type Box struct {
	content string
	title   string
	border  bool
}

func NewBox(title string) *Box {
	return &Box{
		title:  title,
		border: true,
	}
}

func (b *Box) SetContent(content string) {
	b.content = content
}

func (b *Box) Render(r *renderer.Renderer) {
	lines := strings.Split(b.content, "\n")
	maxWidth := 0
	for _, line := range lines {
		if len(line) > maxWidth {
			maxWidth = len(line)
		}
	}
	
	if b.border {
		r.Write("┌─ " + b.title + " " + strings.Repeat("─", maxWidth-len(b.title)-2) + "┐")
		r.NewLine()
		
		for _, line := range lines {
			r.Write("│ " + line + strings.Repeat(" ", maxWidth-len(line)) + " │")
			r.NewLine()
		}
		
		r.Write("└" + strings.Repeat("─", maxWidth+2) + "┘")
		r.NewLine()
	} else {
		for _, line := range lines {
			r.Write(line)
			r.NewLine()
		}
	}
}