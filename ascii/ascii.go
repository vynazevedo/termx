package ascii

import (
	"strings"
	"github.com/vynazevedo/termx/renderer"
	"github.com/vynazevedo/termx/theme"
)

type Art struct {
	content []string
	color   string
}

func New(art string) *Art {
	return &Art{
		content: strings.Split(art, "\n"),
		color:   theme.Current().Primary.Foreground,
	}
}

func (a *Art) WithColor(color string) *Art {
	a.color = color
	return a
}

func (a *Art) Render() {
	r := renderer.New()
	defer r.Close()
	
	for _, line := range a.content {
		r.WriteStyled(line, a.color)
		r.NewLine()
	}
}

func Box(title string, width, height int) *Art {
	art := make([]string, height)
	art[0] = "┌" + strings.Repeat("─", width-2) + "┐"
	for i := 1; i < height-1; i++ {
		art[i] = "│" + strings.Repeat(" ", width-2) + "│"
	}
	art[height-1] = "└" + strings.Repeat("─", width-2) + "┘"
	
	if title != "" && len(title) < width-4 {
		art[0] = art[0][:2] + " " + title + " " + art[0][len(title)+4:]
	}
	
	return &Art{content: art}
}

const (
	KubernetesLogo = `    ⎈
   ╱ ╲
  ╱   ╲
 │  K8s │
  ╲   ╱
   ╲ ╱`

	DockerLogo = `   🐳
  ╱─╲
 │DOC│
 │KER│
  ╲─╱`

	ServerRack = ` ┌─────┐
 │ ▪▪▪ │
 ├─────┤
 │ ▪▪▪ │
 ├─────┤
 │ ▪▪▪ │
 └─────┘`
)