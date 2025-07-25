package table

import (
	"fmt"
	"strings"
	"github.com/vynazevedo/termx/renderer"
	"github.com/vynazevedo/termx/theme"
)

type Table struct {
	headers []string
	rows    [][]string
	widths  []int
	border  bool
	compact bool
	selectedRow int
	interactive bool
}

func New(headers []string) *Table {
	return &Table{
		headers: headers,
		rows:    [][]string{},
		widths:  make([]int, len(headers)),
		border:  true,
		selectedRow: -1,
	}
}

func (t *Table) AddRow(row ...string) *Table {
	if len(row) == len(t.headers) {
		t.rows = append(t.rows, row)
		for i, cell := range row {
			if len(cell) > t.widths[i] {
				t.widths[i] = len(cell)
			}
		}
	}
	return t
}

func (t *Table) WithBorder(border bool) *Table {
	t.border = border
	return t
}

func (t *Table) Compact() *Table {
	t.compact = true
	return t
}

func (t *Table) Interactive() *Table {
	t.interactive = true
	t.selectedRow = 0
	return t
}

func (t *Table) Run() (int, error) {
	if !t.interactive {
		t.Render()
		return -1, nil
	}
	
	r := renderer.New()
	defer r.Close()
	
	for {
		t.render(r)
		
		key, err := r.ReadKey()
		if err != nil {
			return -1, err
		}
		
		switch key {
		case renderer.KeyArrowUp:
			if t.selectedRow > 0 {
				t.selectedRow--
			}
		case renderer.KeyArrowDown:
			if t.selectedRow < len(t.rows)-1 {
				t.selectedRow++
			}
		case renderer.KeyEnter:
			return t.selectedRow, nil
		case renderer.KeyEscape:
			return -1, nil
		}
		
		r.ClearScreen()
	}
}

func (t *Table) Render() {
	r := renderer.New()
	defer r.Close()
	t.render(r)
}

func (t *Table) render(r *renderer.Renderer) {
	th := theme.Current()
	
	for i, h := range t.headers {
		if len(h) > t.widths[i] {
			t.widths[i] = len(h)
		}
	}
	
	if t.border {
		t.renderBorder(r, "top")
	}
	
	for i, h := range t.headers {
		if i == 0 && t.border {
			r.Write("│ ")
		}
		r.WriteStyled(fmt.Sprintf("%-*s", t.widths[i], h), th.Primary.Foreground)
		if i < len(t.headers)-1 {
			r.Write(" │ ")
		} else if t.border {
			r.Write(" │")
		}
	}
	r.NewLine()
	
	if t.border {
		t.renderBorder(r, "middle")
	} else if !t.compact {
		r.WriteStyled(strings.Repeat("─", t.totalWidth()), th.Muted.Foreground)
		r.NewLine()
	}
	
	for idx, row := range t.rows {
		isSelected := t.interactive && idx == t.selectedRow
		
		for i, cell := range row {
			if i == 0 && t.border {
				r.Write("│ ")
			}
			
			cellText := fmt.Sprintf("%-*s", t.widths[i], cell)
			if isSelected {
				r.WriteStyled(cellText, th.Success.Foreground)
			} else {
				r.Write(cellText)
			}
			
			if i < len(row)-1 {
				r.Write(" │ ")
			} else if t.border {
				r.Write(" │")
			}
		}
		r.NewLine()
	}
	
	if t.border {
		t.renderBorder(r, "bottom")
	}
}

func (t *Table) renderBorder(r *renderer.Renderer, position string) {
	chars := map[string][]string{
		"top":    {"┌", "┬", "┐", "─"},
		"middle": {"├", "┼", "┤", "─"},
		"bottom": {"└", "┴", "┘", "─"},
	}
	
	c := chars[position]
	r.Write(c[0])
	for i, w := range t.widths {
		r.Write(strings.Repeat(c[3], w+2))
		if i < len(t.widths)-1 {
			r.Write(c[1])
		}
	}
	r.Write(c[2])
	r.NewLine()
}

func (t *Table) totalWidth() int {
	total := 0
	for _, w := range t.widths {
		total += w + 3
	}
	return total - 1
}