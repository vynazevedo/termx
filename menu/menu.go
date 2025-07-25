package menu

import (
	"fmt"
	"strings"

	"github.com/vynazevedo/termx/renderer"
	"github.com/vynazevedo/termx/theme"
)

// MenuItem represents a single menu item
type MenuItem struct {
	ID          string
	Label       string
	Description string
	Icon        string
	Shortcut    string
	Action      func() error
	Submenu     *Menu
	Disabled    bool
	Separator   bool
}

// Menu represents a navigable menu component
type Menu struct {
	title       string
	items       []MenuItem
	cursor      int
	result      *string
	showIcons   bool
	showDesc    bool
	showShortcuts bool
	theme       *theme.Theme
	renderer    *renderer.Renderer
	breadcrumb  []string
	parent      *Menu
	maxWidth    int
}

// New creates a new Menu instance
func New(title string, result *string) *Menu {
	return &Menu{
		title:         title,
		items:         make([]MenuItem, 0),
		cursor:        0,
		result:        result,
		showIcons:     true,
		showDesc:      true,
		showShortcuts: true,
		theme:         theme.Default(),
		renderer:      renderer.New(),
		breadcrumb:    make([]string, 0),
		maxWidth:      80,
	}
}

// WithTheme sets the menu theme
func (m *Menu) WithTheme(t *theme.Theme) *Menu {
	m.theme = t
	return m
}

// WithMaxWidth sets the maximum width for the menu
func (m *Menu) WithMaxWidth(width int) *Menu {
	m.maxWidth = width
	return m
}

// WithoutIcons disables icon display
func (m *Menu) WithoutIcons() *Menu {
	m.showIcons = false
	return m
}

// WithoutDescriptions disables description display
func (m *Menu) WithoutDescriptions() *Menu {
	m.showDesc = false
	return m
}

// WithoutShortcuts disables shortcut display
func (m *Menu) WithoutShortcuts() *Menu {
	m.showShortcuts = false
	return m
}

// AddItem adds a menu item
func (m *Menu) AddItem(id, label string) *Menu {
	m.items = append(m.items, MenuItem{
		ID:    id,
		Label: label,
	})
	return m
}

// AddItemWithIcon adds a menu item with icon
func (m *Menu) AddItemWithIcon(id, label, icon string) *Menu {
	m.items = append(m.items, MenuItem{
		ID:    id,
		Label: label,
		Icon:  icon,
	})
	return m
}

// AddItemWithDescription adds a menu item with description
func (m *Menu) AddItemWithDescription(id, label, description string) *Menu {
	m.items = append(m.items, MenuItem{
		ID:          id,
		Label:       label,
		Description: description,
	})
	return m
}

// AddItemWithAction adds a menu item with action
func (m *Menu) AddItemWithAction(id, label string, action func() error) *Menu {
	m.items = append(m.items, MenuItem{
		ID:     id,
		Label:  label,
		Action: action,
	})
	return m
}

// AddSubmenu adds a submenu item
func (m *Menu) AddSubmenu(id, label string, submenu *Menu) *Menu {
	submenu.parent = m
	m.items = append(m.items, MenuItem{
		ID:      id,
		Label:   label,
		Submenu: submenu,
	})
	return m
}

// AddSeparator adds a visual separator
func (m *Menu) AddSeparator() *Menu {
	m.items = append(m.items, MenuItem{
		Separator: true,
	})
	return m
}

// AddFullItem adds a complete menu item
func (m *Menu) AddFullItem(item MenuItem) *Menu {
	m.items = append(m.items, item)
	return m
}

// render displays the menu interface
func (m *Menu) render() {
	m.renderer.ClearScreen()
	
	// Breadcrumb
	if len(m.breadcrumb) > 0 {
		breadcrumbStr := strings.Join(m.breadcrumb, " > ")
		fmt.Printf("%s%s%s\n", m.theme.Muted, breadcrumbStr, m.theme.Reset)
	}
	
	// Title
	fmt.Printf("%s%s%s\n", m.theme.Primary, m.title, m.theme.Reset)
	
	// Border
	titleLen := len(m.title)
	if titleLen < m.maxWidth {
		border := strings.Repeat("═", titleLen)
		fmt.Printf("%s%s%s\n", m.theme.Secondary, border, m.theme.Reset)
	}
	
	fmt.Println()
	
	// Menu items
	for i, item := range m.items {
		if item.Separator {
			fmt.Printf("%s%s%s\n", m.theme.Muted, 
				strings.Repeat("─", m.maxWidth/2), m.theme.Reset)
			continue
		}
		
		// Skip cursor for disabled items
		if item.Disabled && i == m.cursor {
			m.moveCursorToNext()
		}
		
		// Cursor indicator
		cursor := "  "
		if i == m.cursor && !item.Disabled {
			cursor = m.theme.Primary + "❯ " + m.theme.Reset
		}
		
		// Icon
		icon := ""
		if m.showIcons && item.Icon != "" {
			icon = item.Icon + " "
		} else if m.showIcons && item.Submenu != nil {
			icon = "📁 "
		} else if m.showIcons {
			icon = "• "
		}
		
		// Label
		label := item.Label
		if item.Disabled {
			label = m.theme.Muted + label + m.theme.Reset
		} else if i == m.cursor {
			label = m.theme.Highlight + label + m.theme.Reset
		}
		
		// Shortcut
		shortcut := ""
		if m.showShortcuts && item.Shortcut != "" {
			shortcut = fmt.Sprintf(" %s[%s]%s", m.theme.Muted, item.Shortcut, m.theme.Reset)
		}
		
		// Submenu indicator
		submenuIndicator := ""
		if item.Submenu != nil {
			submenuIndicator = fmt.Sprintf(" %s▶%s", m.theme.Secondary, m.theme.Reset)
		}
		
		fmt.Printf("%s%s%s%s%s%s\n", cursor, icon, label, shortcut, submenuIndicator, m.theme.Reset)
		
		// Description
		if m.showDesc && item.Description != "" && !item.Disabled {
			desc := item.Description
			if len(desc) > m.maxWidth-6 {
				desc = desc[:m.maxWidth-9] + "..."
			}
			fmt.Printf("    %s%s%s\n", m.theme.Muted, desc, m.theme.Reset)
		}
	}
	
	// Help text
	fmt.Printf("\n%sUse ↑↓ para navegar, Enter para selecionar, Esc para voltar/sair%s\n", 
		m.theme.Muted, m.theme.Reset)
	
	if m.parent != nil {
		fmt.Printf("%s← Voltar para menu anterior%s\n", m.theme.Muted, m.theme.Reset)
	}
}

// moveCursorToNext moves cursor to next non-disabled, non-separator item
func (m *Menu) moveCursorToNext() {
	start := m.cursor
	for {
		m.cursor++
		if m.cursor >= len(m.items) {
			m.cursor = 0
		}
		
		if m.cursor == start {
			break // Prevent infinite loop
		}
		
		item := m.items[m.cursor]
		if !item.Disabled && !item.Separator {
			break
		}
	}
}

// moveCursorToPrev moves cursor to previous non-disabled, non-separator item
func (m *Menu) moveCursorToPrev() {
	start := m.cursor
	for {
		m.cursor--
		if m.cursor < 0 {
			m.cursor = len(m.items) - 1
		}
		
		if m.cursor == start {
			break // Prevent infinite loop
		}
		
		item := m.items[m.cursor]
		if !item.Disabled && !item.Separator {
			break
		}
	}
}

// Run executes the menu interaction
func (m *Menu) Run() error {
	defer m.renderer.Close()
	
	// Skip to first valid item
	if len(m.items) > 0 {
		for m.cursor < len(m.items) {
			item := m.items[m.cursor]
			if !item.Disabled && !item.Separator {
				break
			}
			m.cursor++
		}
	}
	
	for {
		m.render()
		
		key, err := m.renderer.ReadKey()
		if err != nil {
			return err
		}
		
		switch key {
		case renderer.KeyUp:
			m.moveCursorToPrev()
			
		case renderer.KeyDown:
			m.moveCursorToNext()
			
		case renderer.KeyEnter:
			if m.cursor >= len(m.items) {
				continue
			}
			
			item := m.items[m.cursor]
			if item.Disabled || item.Separator {
				continue
			}
			
			// Handle submenu
			if item.Submenu != nil {
				item.Submenu.breadcrumb = append(m.breadcrumb, m.title)
				err := item.Submenu.Run()
				if err != nil && err.Error() != "voltar" {
					return err
				}
				continue
			}
			
			// Handle action
			if item.Action != nil {
				err := item.Action()
				if err != nil {
					fmt.Printf("\n%s%s%s\n", m.theme.Error, err.Error(), m.theme.Reset)
					fmt.Println("Pressione qualquer tecla para continuar...")
					m.renderer.ReadKey()
					continue
				}
			}
			
			// Return selection
			if m.result != nil {
				*m.result = item.ID
			}
			return nil
			
		case renderer.KeyEscape:
			if m.parent != nil {
				return fmt.Errorf("voltar")
			}
			return fmt.Errorf("operação cancelada")
		}
		
		// Handle shortcut keys
		if len(key) == 1 {
			for i, item := range m.items {
				if !item.Disabled && !item.Separator && 
				   strings.ToLower(item.Shortcut) == strings.ToLower(key) {
					m.cursor = i
					
					// Execute the selection
					if item.Submenu != nil {
						item.Submenu.breadcrumb = append(m.breadcrumb, m.title)
						err := item.Submenu.Run()
						if err != nil && err.Error() != "voltar" {
							return err
						}
						continue
					}
					
					if item.Action != nil {
						err := item.Action()
						if err != nil {
							fmt.Printf("\n%s%s%s\n", m.theme.Error, err.Error(), m.theme.Reset)
							fmt.Println("Pressione qualquer tecla para continuar...")
							m.renderer.ReadKey()
							continue
						}
					}
					
					if m.result != nil {
						*m.result = item.ID
					}
					return nil
				}
			}
		}
	}
}

// Predefined menu configurations
func MainMenu(result *string) *Menu {
	return New("Menu Principal", result).
		AddItemWithIcon("new", "Novo Projeto", "📝").
		AddItemWithIcon("open", "Abrir Projeto", "📂").
		AddItemWithIcon("recent", "Projetos Recentes", "🕒").
		AddSeparator().
		AddItemWithIcon("settings", "Configurações", "⚙️").
		AddItemWithIcon("about", "Sobre", "ℹ️").
		AddSeparator().
		AddItemWithIcon("exit", "Sair", "🚪")
}

func FileMenu(result *string) *Menu {
	return New("Arquivo", result).
		AddFullItem(MenuItem{
			ID:          "new",
			Label:       "Novo",
			Icon:        "📄",
			Shortcut:    "n",
			Description: "Criar um novo arquivo",
		}).
		AddFullItem(MenuItem{
			ID:          "open",
			Label:       "Abrir",
			Icon:        "📂",
			Shortcut:    "o",
			Description: "Abrir arquivo existente",
		}).
		AddFullItem(MenuItem{
			ID:          "save",
			Label:       "Salvar",
			Icon:        "💾",
			Shortcut:    "s",
			Description: "Salvar arquivo atual",
		}).
		AddSeparator().
		AddFullItem(MenuItem{
			ID:          "exit",
			Label:       "Sair",
			Icon:        "🚪",
			Shortcut:    "q",
			Description: "Sair da aplicação",
		})
}

func ToolsMenu(result *string) *Menu {
	return New("Ferramentas", result).
		AddItemWithDescription("git", "Git", "Controle de versão").
		AddItemWithDescription("docker", "Docker", "Containerização").
		AddItemWithDescription("k8s", "Kubernetes", "Orquestração de containers").
		AddSeparator().
		AddItemWithDescription("lint", "Linter", "Análise de código").
		AddItemWithDescription("test", "Testes", "Executar testes").
		AddItemWithDescription("build", "Build", "Compilar projeto")
}