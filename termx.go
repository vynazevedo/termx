package termx

import (
	"github.com/vynazevedo/termx/ascii"
	"github.com/vynazevedo/termx/chart"
	"github.com/vynazevedo/termx/confirm"
	"github.com/vynazevedo/termx/form"
	"github.com/vynazevedo/termx/input"
	"github.com/vynazevedo/termx/layout"
	"github.com/vynazevedo/termx/progress"
	"github.com/vynazevedo/termx/renderer"
	"github.com/vynazevedo/termx/selector"
	"github.com/vynazevedo/termx/table"
	"github.com/vynazevedo/termx/theme"
)

var (
	// Core components
	Form    = form.New
	Input   = input.New
	Select  = selector.New
	Confirm = confirm.New
	Table   = table.New
	Chart   = chart.New
	
	// Visual components
	ASCII = ascii.New
	Box   = ascii.Box
	
	// Progress and loading
	Progress = progress.NewBar
	Spinner  = progress.NewSpinner
	
	// Layout components
	Split     = layout.NewSplit
	BoxLayout = layout.NewBox
	
	// Validators
	Required  = input.Required
	MinLength = input.MinLength
	MaxLength = input.MaxLength
	Email     = input.Email
	
	// Form helpers
	WithPlaceholder = form.WithPlaceholder
	WithValidator   = form.WithValidator
	WithMaxLength   = form.WithMaxLength
)

const (
	KubernetesLogo = ascii.KubernetesLogo
	DockerLogo     = ascii.DockerLogo
	ServerRack     = ascii.ServerRack
)

func SetTheme(t *theme.Theme) {
	theme.Set(t)
}

func GetTheme() *theme.Theme {
	return theme.Current()
}

func ClearScreen() {
	renderer.ClearScreen()
}