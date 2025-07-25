package chart

import (
	"fmt"
	"math"
	"strings"
	"github.com/vynazevedo/termx/renderer"
	"github.com/vynazevedo/termx/theme"
)

type Chart struct {
	data   []float64
	labels []string
	width  int
	height int
	style  string
}

func New(data []float64) *Chart {
	return &Chart{
		data:   data,
		width:  60,
		height: 15,
		style:  "bar",
	}
}

func (c *Chart) WithLabels(labels []string) *Chart {
	c.labels = labels
	return c
}

func (c *Chart) WithSize(width, height int) *Chart {
	c.width = width
	c.height = height
	return c
}

func (c *Chart) WithStyle(style string) *Chart {
	c.style = style
	return c
}

func (c *Chart) Render() {
	r := renderer.New()
	defer r.Close()
	
	switch c.style {
	case "line":
		c.renderLine(r)
	case "scatter":
		c.renderScatter(r)
	default:
		c.renderBar(r)
	}
}

func (c *Chart) renderBar(r *renderer.Renderer) {
	if len(c.data) == 0 {
		return
	}
	
	max := c.data[0]
	for _, v := range c.data {
		if v > max {
			max = v
		}
	}
	
	barWidth := c.width / len(c.data)
	if barWidth < 3 {
		barWidth = 3
	}
	
	th := theme.Current()
	
	for row := c.height; row > 0; row-- {
		threshold := (float64(row) / float64(c.height)) * max
		
		r.WriteStyled(fmt.Sprintf("%6.1f │", threshold), th.Muted.Foreground)
		
		for _, value := range c.data {
			if value >= threshold {
				r.WriteStyled(strings.Repeat("█", barWidth-1)+" ", th.Primary.Foreground)
			} else {
				r.Write(strings.Repeat(" ", barWidth))
			}
		}
		r.NewLine()
	}
	
	r.WriteStyled("       └"+strings.Repeat("─", len(c.data)*barWidth), th.Muted.Foreground)
	r.NewLine()
	
	if c.labels != nil {
		r.Write("        ")
		for i, label := range c.labels {
			if i < len(c.data) {
				r.Write(fmt.Sprintf("%-*s", barWidth, label))
			}
		}
		r.NewLine()
	}
}

func (c *Chart) renderLine(r *renderer.Renderer) {
	if len(c.data) == 0 {
		return
	}
	
	max := c.data[0]
	min := c.data[0]
	for _, v := range c.data {
		if v > max {
			max = v
		}
		if v < min {
			min = v
		}
	}
	
	grid := make([][]string, c.height)
	for i := range grid {
		grid[i] = make([]string, c.width)
		for j := range grid[i] {
			grid[i][j] = " "
		}
	}
	
	chars := map[string]string{
		"up":    "╱",
		"down":  "╲",
		"flat":  "─",
		"point": "●",
	}
	
	scale := float64(c.height-1) / (max - min)
	xStep := float64(c.width-1) / float64(len(c.data)-1)
	
	for i := 0; i < len(c.data)-1; i++ {
		x1 := int(float64(i) * xStep)
		y1 := c.height - 1 - int((c.data[i]-min)*scale)
		x2 := int(float64(i+1) * xStep)
		y2 := c.height - 1 - int((c.data[i+1]-min)*scale)
		
		if y1 >= 0 && y1 < c.height && x1 >= 0 && x1 < c.width {
			grid[y1][x1] = chars["point"]
		}
		
		dx := x2 - x1
		dy := y2 - y1
		steps := int(math.Max(math.Abs(float64(dx)), math.Abs(float64(dy))))
		
		for step := 1; step < steps; step++ {
			x := x1 + (dx*step)/steps
			y := y1 + (dy*step)/steps
			
			if y >= 0 && y < c.height && x >= 0 && x < c.width {
				if dy > 0 {
					grid[y][x] = chars["down"]
				} else if dy < 0 {
					grid[y][x] = chars["up"]
				} else {
					grid[y][x] = chars["flat"]
				}
			}
		}
	}
	
	if len(c.data) > 0 {
		lastX := int(float64(len(c.data)-1) * xStep)
		lastY := c.height - 1 - int((c.data[len(c.data)-1]-min)*scale)
		if lastY >= 0 && lastY < c.height && lastX >= 0 && lastX < c.width {
			grid[lastY][lastX] = chars["point"]
		}
	}
	
	th := theme.Current()
	
	for i, row := range grid {
		value := max - (float64(i)/float64(c.height-1))*(max-min)
		r.WriteStyled(fmt.Sprintf("%6.1f │", value), th.Muted.Foreground)
		
		for _, cell := range row {
			if cell == " " {
				r.Write(cell)
			} else {
				r.WriteStyled(cell, th.Primary.Foreground)
			}
		}
		r.NewLine()
	}
	
	r.WriteStyled("       └"+strings.Repeat("─", c.width), th.Muted.Foreground)
	r.NewLine()
}

func (c *Chart) renderScatter(r *renderer.Renderer) {
	if len(c.data) == 0 {
		return
	}
	
	max := c.data[0]
	for _, v := range c.data {
		if v > max {
			max = v
		}
	}
	
	grid := make([][]bool, c.height)
	for i := range grid {
		grid[i] = make([]bool, c.width)
	}
	
	xStep := float64(c.width) / float64(len(c.data))
	scale := float64(c.height-1) / max
	
	for i, value := range c.data {
		x := int(float64(i) * xStep)
		y := c.height - 1 - int(value*scale)
		
		if x >= 0 && x < c.width && y >= 0 && y < c.height {
			grid[y][x] = true
		}
	}
	
	th := theme.Current()
	
	for row := 0; row < c.height; row++ {
		value := max - (float64(row)/float64(c.height-1))*max
		r.WriteStyled(fmt.Sprintf("%6.1f │", value), th.Muted.Foreground)
		
		for col := 0; col < c.width; col++ {
			if grid[row][col] {
				r.WriteStyled("●", th.Primary.Foreground)
			} else {
				r.Write(" ")
			}
		}
		r.NewLine()
	}
	
	r.WriteStyled("       └"+strings.Repeat("─", c.width), th.Muted.Foreground)
	r.NewLine()
}