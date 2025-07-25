package confirm

import (
	"errors"
	"strings"

	"github.com/vynazevedo/termx/renderer"
	"github.com/vynazevedo/termx/theme"
)

type Confirm struct {
	Label    string
	Result   *bool
	Default  bool
	
	selected bool
	renderer *renderer.Renderer
}

func New(label string, result *bool) *Confirm {
	return &Confirm{
		Label:    label,
		Result:   result,
		Default:  false,
		selected: false,
	}
}

func (c *Confirm) WithDefault(defaultValue bool) *Confirm {
	c.Default = defaultValue
	c.selected = defaultValue
	return c
}

func (c *Confirm) Run() error {
	c.renderer = renderer.New()
	if err := c.renderer.Init(); err != nil {
		return err
	}
	defer c.renderer.Restore()

	for {
		c.render()
		
		event, err := renderer.ReadInput()
		if err != nil {
			return err
		}

		switch event.Key {
		case renderer.KeyCtrlC:
			return errors.New("cancelled")
		
		case renderer.KeyEnter:
			if c.Result != nil {
				*c.Result = c.selected
			}
			return nil
		
		case renderer.KeyArrowLeft, renderer.KeyArrowRight:
			c.selected = !c.selected
		
		case renderer.KeyTab:
			c.selected = !c.selected
		
		default:
			// Handle Y/N shortcuts
			if event.Rune != 0 {
				lower := strings.ToLower(string(event.Rune))
				switch lower {
				case "y":
					c.selected = true
					if c.Result != nil {
						*c.Result = true
					}
					return nil
				case "n":
					c.selected = false
					if c.Result != nil {
						*c.Result = false
					}
					return nil
				}
			}
		}
	}
}

func (c *Confirm) render() {
	c.renderer.Clear()
	th := theme.Current()
	
	// Calculate centered position
	boxWidth := 40
	boxHeight := 5
	
	startX := (c.renderer.Width() - boxWidth) / 2
	startY := (c.renderer.Height() - boxHeight) / 2
	
	// Draw box
	c.renderer.Box(startX, startY, boxWidth, boxHeight, "")
	
	// Label
	labelY := startY + 1
	c.renderer.PrintCentered(labelY, th.Primary.Sprint(c.Label))
	
	// Options
	optionsY := startY + 3
	optionsX := startX + (boxWidth / 2) - 10
	
	yesStyle := th.Text
	noStyle := th.Text
	
	if c.selected {
		yesStyle = th.Selected
	} else {
		noStyle = th.Selected
	}
	
	c.renderer.Print(optionsX, optionsY, yesStyle.Sprint("  Yes (Y)  "))
	c.renderer.Print(optionsX+12, optionsY, noStyle.Sprint("  No (N)  "))
	
	// Help text
	helpY := c.renderer.Height() - 2
	helpText := "←→/Tab to toggle • Y/N shortcuts • Enter to confirm • Ctrl+C to cancel"
	c.renderer.PrintCentered(helpY, th.TextDim.Sprint(helpText))
}