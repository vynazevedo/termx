package combobox

import (
	"fmt"
	"strings"

	"github.com/vynazevedo/termx/renderer"
	"github.com/vynazevedo/termx/theme"
)

// ComboBox represents a searchable select component with custom input
type ComboBox struct {
	label        string
	options      []string
	result       *string
	value        string
	cursor       int
	filtered     []string
	showDropdown bool
	allowCustom  bool
	placeholder  string
	validator    func(string) error
	theme        *theme.Theme
	renderer     *renderer.Renderer
	maxDisplay   int
	caseSensitive bool
}

// New creates a new ComboBox instance
func New(label string, options []string, result *string) *ComboBox {
	return &ComboBox{
		label:        label,
		options:      options,
		result:       result,
		value:        "",
		cursor:       0,
		filtered:     make([]string, len(options)),
		showDropdown: false,
		allowCustom:  true,
		theme:        theme.Default(),
		renderer:     renderer.New(),
		maxDisplay:   8,
		caseSensitive: false,
	}
}

// WithPlaceholder sets a placeholder text
func (cb *ComboBox) WithPlaceholder(placeholder string) *ComboBox {
	cb.placeholder = placeholder
	return cb
}

// WithValidator sets a validation function
func (cb *ComboBox) WithValidator(validator func(string) error) *ComboBox {
	cb.validator = validator
	return cb
}

// WithTheme sets the theme
func (cb *ComboBox) WithTheme(t *theme.Theme) *ComboBox {
	cb.theme = t
	return cb
}

// WithMaxDisplay sets maximum number of options to display
func (cb *ComboBox) WithMaxDisplay(max int) *ComboBox {
	cb.maxDisplay = max
	return cb
}

// WithoutCustomInput disallows custom input (strict selection only)
func (cb *ComboBox) WithoutCustomInput() *ComboBox {
	cb.allowCustom = false
	return cb
}

// WithCaseSensitiveSearch enables case-sensitive search
func (cb *ComboBox) WithCaseSensitiveSearch() *ComboBox {
	cb.caseSensitive = true
	return cb
}

// filterOptions filters options based on current input value
func (cb *ComboBox) filterOptions() {
	cb.filtered = cb.filtered[:0]
	
	if cb.value == "" {
		cb.filtered = append(cb.filtered, cb.options...)
		return
	}
	
	searchValue := cb.value
	if !cb.caseSensitive {
		searchValue = strings.ToLower(searchValue)
	}
	
	// Exact matches first
	for _, option := range cb.options {
		compareOption := option
		if !cb.caseSensitive {
			compareOption = strings.ToLower(compareOption)
		}
		
		if compareOption == searchValue {
			cb.filtered = append(cb.filtered, option)
		}
	}
	
	// Prefix matches
	for _, option := range cb.options {
		compareOption := option
		if !cb.caseSensitive {
			compareOption = strings.ToLower(compareOption)
		}
		
		if strings.HasPrefix(compareOption, searchValue) && compareOption != searchValue {
			cb.filtered = append(cb.filtered, option)
		}
	}
	
	// Contains matches
	for _, option := range cb.options {
		compareOption := option
		if !cb.caseSensitive {
			compareOption = strings.ToLower(compareOption)
		}
		
		if strings.Contains(compareOption, searchValue) && 
		   !strings.HasPrefix(compareOption, searchValue) {
			cb.filtered = append(cb.filtered, option)
		}
	}
	
	// Remove duplicates
	seen := make(map[string]bool)
	uniqueFiltered := cb.filtered[:0]
	for _, option := range cb.filtered {
		if !seen[option] {
			seen[option] = true
			uniqueFiltered = append(uniqueFiltered, option)
		}
	}
	cb.filtered = uniqueFiltered
}

// render displays the combobox interface
func (cb *ComboBox) render() {
	cb.renderer.ClearScreen()
	
	// Label
	fmt.Printf("%s%s%s\n", cb.theme.Primary, cb.label, cb.theme.Reset)
	
	// Help text
	helpText := "Digite para buscar, ↑↓ para navegar, Enter para selecionar, Esc para cancelar"
	if cb.allowCustom {
		helpText = "Digite valor customizado ou busque, ↑↓ para navegar, Enter para confirmar"
	}
	fmt.Printf("%s%s%s\n", cb.theme.Muted, helpText, cb.theme.Reset)
	
	// Input field
	displayValue := cb.value
	if displayValue == "" && cb.placeholder != "" {
		displayValue = cb.theme.Muted + cb.placeholder + cb.theme.Reset
	} else {
		displayValue = cb.theme.Input + displayValue + cb.theme.Reset
	}
	
	fmt.Printf("\n%s> %s%s\n", cb.theme.Primary, displayValue, cb.theme.Reset)
	
	// Dropdown
	if cb.showDropdown && len(cb.filtered) > 0 {
		fmt.Printf("\n%sOpções disponíveis:%s\n", cb.theme.Secondary, cb.theme.Reset)
		
		displayCount := cb.maxDisplay
		if len(cb.filtered) < displayCount {
			displayCount = len(cb.filtered)
		}
		
		for i := 0; i < displayCount; i++ {
			option := cb.filtered[i]
			cursor := "  "
			
			if i == cb.cursor {
				cursor = cb.theme.Primary + "❯ " + cb.theme.Reset
				option = cb.theme.Highlight + option + cb.theme.Reset
			}
			
			fmt.Printf("%s%s\n", cursor, option)
		}
		
		if len(cb.filtered) > displayCount {
			remaining := len(cb.filtered) - displayCount
			fmt.Printf("%s... e mais %d opções%s\n", cb.theme.Muted, remaining, cb.theme.Reset)
		}
	} else if cb.showDropdown && len(cb.filtered) == 0 && cb.value != "" {
		if cb.allowCustom {
			fmt.Printf("\n%sNenhuma opção encontrada. Valor customizado será usado.%s\n", 
				cb.theme.Warning, cb.theme.Reset)
		} else {
			fmt.Printf("\n%sNenhuma opção encontrada.%s\n", cb.theme.Error, cb.theme.Reset)
		}
	}
	
	// Custom input indicator
	if cb.allowCustom && cb.value != "" {
		isExistingOption := false
		for _, option := range cb.options {
			if (cb.caseSensitive && option == cb.value) || 
			   (!cb.caseSensitive && strings.ToLower(option) == strings.ToLower(cb.value)) {
				isExistingOption = true
				break
			}
		}
		
		if !isExistingOption {
			fmt.Printf("\n%s💡 Valor customizado: \"%s\"%s\n", 
				cb.theme.Info, cb.value, cb.theme.Reset)
		}
	}
}

// validateInput validates the current input
func (cb *ComboBox) validateInput() error {
	if !cb.allowCustom && cb.value != "" {
		// Check if value exists in options
		found := false
		for _, option := range cb.options {
			if (cb.caseSensitive && option == cb.value) || 
			   (!cb.caseSensitive && strings.ToLower(option) == strings.ToLower(cb.value)) {
				found = true
				break
			}
		}
		
		if !found {
			return fmt.Errorf("valor deve ser selecionado da lista de opções")
		}
	}
	
	if cb.validator != nil {
		return cb.validator(cb.value)
	}
	
	return nil
}

// Run executes the combobox interaction
func (cb *ComboBox) Run() error {
	defer cb.renderer.Close()
	
	cb.filterOptions()
	
	for {
		cb.render()
		
		key, err := cb.renderer.ReadKey()
		if err != nil {
			return err
		}
		
		switch key {
		case renderer.KeyUp:
			if cb.showDropdown && cb.cursor > 0 {
				cb.cursor--
			} else if !cb.showDropdown {
				cb.showDropdown = true
				cb.cursor = 0
			}
			
		case renderer.KeyDown:
			if cb.showDropdown && cb.cursor < len(cb.filtered)-1 {
				cb.cursor++
			} else if !cb.showDropdown {
				cb.showDropdown = true
				cb.cursor = 0
			}
			
		case renderer.KeyTab:
			if len(cb.filtered) > 0 && cb.cursor < len(cb.filtered) {
				cb.value = cb.filtered[cb.cursor]
				cb.showDropdown = false
				cb.filterOptions()
			}
			
		case renderer.KeyEnter:
			if cb.showDropdown && len(cb.filtered) > 0 && cb.cursor < len(cb.filtered) {
				cb.value = cb.filtered[cb.cursor]
				cb.showDropdown = false
			}
			
			if err := cb.validateInput(); err != nil {
				fmt.Printf("\n%s%s%s\n", cb.theme.Error, err.Error(), cb.theme.Reset)
				fmt.Println("Pressione qualquer tecla para continuar...")
				cb.renderer.ReadKey()
				continue
			}
			
			*cb.result = cb.value
			return nil
			
		case renderer.KeyBackspace:
			if len(cb.value) > 0 {
				cb.value = cb.value[:len(cb.value)-1]
				cb.filterOptions()
				cb.cursor = 0
				cb.showDropdown = len(cb.filtered) > 0
			}
			
		case renderer.KeyEscape:
			cb.showDropdown = false
			if cb.value == "" {
				return fmt.Errorf("operação cancelada")
			}
			
		default:
			// Handle regular character input
			if len(key) == 1 && key[0] >= 32 && key[0] <= 126 {
				cb.value += key
				cb.filterOptions()
				cb.cursor = 0
				cb.showDropdown = true
			}
		}
	}
}

// Predefined combobox configurations
func Countries(result *string) *ComboBox {
	options := []string{
		"Brasil", "Estados Unidos", "Canadá", "México", "Argentina",
		"Reino Unido", "França", "Alemanha", "Itália", "Espanha",
		"China", "Japão", "Coreia do Sul", "Índia", "Austrália",
	}
	return New("Selecione o país:", options, result).
		WithPlaceholder("Digite ou selecione um país")
}

func ProgrammingLanguages(result *string) *ComboBox {
	options := []string{
		"Go", "Python", "JavaScript", "TypeScript", "Rust",
		"Java", "C++", "C#", "PHP", "Ruby", "Swift", "Kotlin",
	}
	return New("Linguagem de programação:", options, result).
		WithPlaceholder("Digite ou selecione uma linguagem").
		WithoutCustomInput()
}

func Databases(result *string) *ComboBox {
	options := []string{
		"PostgreSQL", "MySQL", "SQLite", "MongoDB", "Redis",
		"Elasticsearch", "CockroachDB", "InfluxDB", "Cassandra",
	}
	return New("Banco de dados:", options, result).
		WithPlaceholder("Digite ou selecione um banco de dados")
}

func CloudProviders(result *string) *ComboBox {
	options := []string{
		"AWS", "Google Cloud", "Microsoft Azure", "DigitalOcean",
		"Heroku", "Vercel", "Netlify", "Railway", "Fly.io",
	}
	return New("Provedor de nuvem:", options, result).
		WithPlaceholder("Digite ou selecione um provedor")
}