package multiselect

import (
	"fmt"
	"strings"

	"github.com/vynazevedo/termx/renderer"
	"github.com/vynazevedo/termx/theme"
)

// MultiSelect represents a multi-selection component
type MultiSelect struct {
	label       string
	options     []string
	selected    map[int]bool
	result      *[]string
	cursor      int
	searchMode  bool
	searchTerm  string
	filtered    []int
	placeholder string
	validator   func([]string) error
	theme       *theme.Theme
	renderer    *renderer.Renderer
	minSelect   int
	maxSelect   int
	showHelp    bool
}

// New creates a new MultiSelect instance
func New(label string, options []string, result *[]string) *MultiSelect {
	return &MultiSelect{
		label:      label,
		options:    options,
		selected:   make(map[int]bool),
		result:     result,
		cursor:     0,
		searchMode: false,
		filtered:   make([]int, len(options)),
		theme:      theme.Default,
		renderer:   renderer.New(),
		minSelect:  0,
		maxSelect:  len(options),
		showHelp:   true,
	}
}

// WithPlaceholder sets a placeholder text
func (ms *MultiSelect) WithPlaceholder(placeholder string) *MultiSelect {
	ms.placeholder = placeholder
	return ms
}

// WithValidator sets a validation function
func (ms *MultiSelect) WithValidator(validator func([]string) error) *MultiSelect {
	ms.validator = validator
	return ms
}

// WithTheme sets the theme
func (ms *MultiSelect) WithTheme(t *theme.Theme) *MultiSelect {
	ms.theme = t
	return ms
}

// WithMinSelect sets minimum number of selections required
func (ms *MultiSelect) WithMinSelect(min int) *MultiSelect {
	ms.minSelect = min
	return ms
}

// WithMaxSelect sets maximum number of selections allowed
func (ms *MultiSelect) WithMaxSelect(max int) *MultiSelect {
	ms.maxSelect = max
	return ms
}

// WithoutHelp disables help text display
func (ms *MultiSelect) WithoutHelp() *MultiSelect {
	ms.showHelp = false
	return ms
}

// initFiltered initializes the filtered options list
func (ms *MultiSelect) initFiltered() {
	for i := range ms.options {
		ms.filtered[i] = i
	}
}

// filterOptions filters options based on search term
func (ms *MultiSelect) filterOptions() {
	if ms.searchTerm == "" {
		ms.initFiltered()
		return
	}

	ms.filtered = ms.filtered[:0]
	searchLower := strings.ToLower(ms.searchTerm)
	
	for i, option := range ms.options {
		if strings.Contains(strings.ToLower(option), searchLower) {
			ms.filtered = append(ms.filtered, i)
		}
	}
}

// render displays the multi-select interface
func (ms *MultiSelect) render() {
	ms.renderer.ClearScreen()
	
	// Header
	fmt.Printf("%s\n", ms.theme.Primary.Sprint(ms.label))
	
	if ms.showHelp {
		helpText := "Use ↑↓ para navegar, Space para selecionar, / para buscar, Enter para confirmar, Esc para cancelar"
		fmt.Printf("%s\n", ms.theme.Muted.Sprint(helpText))
	}
	
	// Search bar
	if ms.searchMode {
		fmt.Printf("\nBuscar: %s\n", ms.theme.Primary.Sprint(ms.searchTerm))
	} else if ms.searchTerm != "" {
		fmt.Printf("\n%s\n", ms.theme.Muted.Sprint(fmt.Sprintf("Buscar: %s (pressione / para editar)", ms.searchTerm)))
	}
	
	// Selection count
	selectedCount := len(ms.getSelectedValues())
	fmt.Printf("\n%s\n", ms.theme.Secondary.Sprint(fmt.Sprintf("Selecionados: %d/%d", selectedCount, len(ms.options))))
	
	if ms.placeholder != "" && selectedCount == 0 {
		fmt.Printf("%s\n", ms.theme.Muted.Sprint(ms.placeholder))
	}
	
	fmt.Println()
	
	// Options list
	visibleOptions := ms.filtered
	if len(visibleOptions) == 0 {
		fmt.Printf("%s\n", ms.theme.Error.Sprint("Nenhuma opção encontrada"))
		return
	}
	
	// Calculate display window
	maxDisplay := 10
	start := 0
	if ms.cursor >= maxDisplay {
		start = ms.cursor - maxDisplay + 1
	}
	end := start + maxDisplay
	if end > len(visibleOptions) {
		end = len(visibleOptions)
	}
	
	for i := start; i < end; i++ {
		optionIndex := visibleOptions[i]
		option := ms.options[optionIndex]
		
		// Cursor indicator
		cursor := "  "
		if i == ms.cursor {
			cursor = ms.theme.Primary.Sprint("❯ ")
		}
		
		// Selection indicator
		checkbox := "☐"
		if ms.selected[optionIndex] {
			checkbox = ms.theme.Success.Sprint("☑")
		}
		
		// Highlight current option
		if i == ms.cursor {
			option = ms.theme.Selected.Sprint(option)
		}
		
		fmt.Printf("%s%s %s\n", cursor, checkbox, option)
	}
	
	// Show more indicator
	if end < len(visibleOptions) {
		fmt.Printf("%s\n", ms.theme.Muted.Sprint(fmt.Sprintf("... e mais %d opções", len(visibleOptions)-end)))
	}
}

// getSelectedValues returns the currently selected values
func (ms *MultiSelect) getSelectedValues() []string {
	var selected []string
	for i := range ms.options {
		if ms.selected[i] {
			selected = append(selected, ms.options[i])
		}
	}
	return selected
}

// validateSelection validates the current selection
func (ms *MultiSelect) validateSelection() error {
	selected := ms.getSelectedValues()
	
	// Check minimum selections
	if len(selected) < ms.minSelect {
		return fmt.Errorf("pelo menos %d opções devem ser selecionadas", ms.minSelect)
	}
	
	// Check maximum selections
	if len(selected) > ms.maxSelect {
		return fmt.Errorf("no máximo %d opções podem ser selecionadas", ms.maxSelect)
	}
	
	// Custom validation
	if ms.validator != nil {
		return ms.validator(selected)
	}
	
	return nil
}

// Run executes the multi-select interaction
func (ms *MultiSelect) Run() error {
	defer ms.renderer.Close()
	
	ms.initFiltered()
	ms.filterOptions()
	
	for {
		ms.render()
		
		key, err := ms.renderer.ReadKey()
		if err != nil {
			return err
		}
		
		if ms.searchMode {
			switch key {
			case renderer.KeyEscape:
				ms.searchMode = false
			case renderer.KeyEnter:
				ms.searchMode = false
				ms.filterOptions()
				ms.cursor = 0
			case renderer.KeyBackspace:
				if len(ms.searchTerm) > 0 {
					ms.searchTerm = ms.searchTerm[:len(ms.searchTerm)-1]
					ms.filterOptions()
					ms.cursor = 0
				}
			default:
				if len(key) == 1 && key[0] >= 32 && key[0] <= 126 {
					ms.searchTerm += key
					ms.filterOptions()
					ms.cursor = 0
				}
			}
			continue
		}
		
		switch key {
		case renderer.KeyUp:
			if ms.cursor > 0 {
				ms.cursor--
			}
		case renderer.KeyDown:
			if ms.cursor < len(ms.filtered)-1 {
				ms.cursor++
			}
		case renderer.KeySpace:
			if len(ms.filtered) > 0 {
				optionIndex := ms.filtered[ms.cursor]
				if ms.selected[optionIndex] {
					delete(ms.selected, optionIndex)
				} else {
					// Check if we can select more
					if len(ms.getSelectedValues()) < ms.maxSelect {
						ms.selected[optionIndex] = true
					}
				}
			}
		case "/":
			ms.searchMode = true
		case "c":
			if len(ms.filtered) > 0 {
				ms.searchTerm = ""
				ms.filterOptions()
				ms.cursor = 0
			}
		case "a":
			// Select all visible options
			for _, optionIndex := range ms.filtered {
				if len(ms.getSelectedValues()) < ms.maxSelect {
					ms.selected[optionIndex] = true
				} else {
					break
				}
			}
		case "n":
			// Deselect all
			ms.selected = make(map[int]bool)
		case renderer.KeyEnter:
			if err := ms.validateSelection(); err != nil {
				fmt.Printf("\n%s%s%s\n", ms.theme.Error, err.Error(), ms.theme.Reset)
				fmt.Println("Pressione qualquer tecla para continuar...")
				ms.renderer.ReadKey()
				continue
			}
			
			selected := ms.getSelectedValues()
			*ms.result = selected
			return nil
		case renderer.KeyEscape:
			return fmt.Errorf("operação cancelada")
		}
	}
}

// Predefined multi-select configurations
func Technologies(result *[]string) *MultiSelect {
	options := []string{
		"Go", "Python", "JavaScript", "TypeScript", "Rust",
		"Java", "C++", "C#", "PHP", "Ruby",
		"Docker", "Kubernetes", "AWS", "GCP", "Azure",
		"PostgreSQL", "MySQL", "MongoDB", "Redis", "Elasticsearch",
	}
	return New("Selecione as tecnologias:", options, result).
		WithPlaceholder("Nenhuma tecnologia selecionada").
		WithMinSelect(1).
		WithMaxSelect(5)
}

func Environments(result *[]string) *MultiSelect {
	options := []string{"development", "staging", "production", "testing"}
	return New("Selecione os ambientes:", options, result).
		WithMinSelect(1)
}

func Features(result *[]string) *MultiSelect {
	options := []string{
		"Autenticação", "Autorização", "Cache", "Logging",
		"Monitoramento", "Métricas", "Backup", "Recuperação",
		"API REST", "GraphQL", "WebSocket", "gRPC",
	}
	return New("Selecione os recursos:", options, result).
		WithPlaceholder("Nenhum recurso selecionado")
}