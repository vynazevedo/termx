package main

import (
	"fmt"
	"log"
	"time"

	"github.com/vynazevedo/termx"
)

type Project struct {
	Name        string
	Language    string
	Database    string
	Cloud       string
	Environment []string
	Features    []string
	Status      string
}

func main() {
	termx.ClearScreen()
	
	// Modern ASCII logo
	termx.ASCII(`
╔╦╗┌─┐┬─┐┌┬┐╦ ╦  ╔╦╗┌─┐┌─┐┬ ┬┌┐ ┌─┐┌─┐┬─┐┌┬┐
 ║ ├┤ ├┬┘││║╚╦╝   ║║├─┤└─┐├─┤├┴┐│ │├─┤├┬┘ ││
 ╩ └─┘┴└─┴ ╩ ╩   ═╩╝┴ ┴└─┘┴ ┴└─┘└─┘┴ ┴┴└──┴┘`).WithColor("\033[35m").Render()
	
	fmt.Printf("\n%s🚀 Dashboard Interativo de Criação de Projetos%s\n", "\033[36m", "\033[0m")
	fmt.Printf("%s━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━%s\n\n", "\033[36m", "\033[0m")
	
	// Create project instance
	project := &Project{}
	
	// Step 1: Main menu
	mainMenuLoop(project)
}

func mainMenuLoop(project *Project) {
	for {
		var action string
		err := termx.Select("O que você gostaria de fazer?", []string{
			"🆕 Criar Novo Projeto",
			"📋 Usar Template",
			"📥 Importar Projeto",
			"🕒 Projetos Recentes",
			"⚙️ Configurações",
			"🚪 Sair",
		}, &action).Run()
		
		if err != nil {
			log.Fatal(err)
		}
		
		switch action {
		case "🆕 Criar Novo Projeto":
			createProjectWizard(project)
		case "📋 Usar Template":
			templateSelector(project)
		case "📥 Importar Projeto":
			importProject(project)
		case "🕒 Projetos Recentes":
			showRecentProjects()
		case "⚙️ Configurações":
			showSettings()
		case "🚪 Sair":
			fmt.Println("\n👋 Obrigado por usar o TermX Dashboard!")
			return
		}
	}
}

func createProjectWizard(project *Project) {
	termx.ClearScreen()
	
	fmt.Printf("%s✨ Assistente de Criação de Projeto%s\n", "\033[35m", "\033[0m")
	fmt.Printf("%s════════════════════════════════════%s\n\n", "\033[35m", "\033[0m")
	
	// Step-by-step form with basic components
	var confirmed bool
	
	err := termx.Form().
		Input("Nome do projeto:", &project.Name).
		Select("Linguagem principal:", []string{
			"Go", "Python", "JavaScript", "TypeScript", "Rust",
			"Java", "C#", "PHP", "Ruby", "Swift",
		}, &project.Language).
		Select("Banco de dados:", []string{
			"PostgreSQL", "MySQL", "MongoDB", "Redis", "SQLite",
		}, &project.Database).
		Select("Provedor de nuvem:", []string{
			"AWS", "Google Cloud", "Azure", "DigitalOcean",
		}, &project.Cloud).
		Confirm("Continuar com a criação?", &confirmed).
		Run()
	
	if err != nil {
		fmt.Printf("❌ Erro: %v\n", err)
		return
	}
	
	if !confirmed {
		fmt.Println("❌ Criação cancelada pelo usuário.")
		return
	}
	
	// Show creation progress
	createProjectWithProgress(project)
}

func createProjectWithProgress(project *Project) {
	termx.ClearScreen()
	
	fmt.Printf("%s🔧 Criando projeto \"%s\"%s\n\n", "\033[33m", project.Name, "\033[0m")
	
	steps := []struct {
		name    string
		message string
		delay   time.Duration
	}{
		{"Inicializando estrutura", "Criando diretórios e arquivos base", 1 * time.Second},
		{"Configurando linguagem", fmt.Sprintf("Preparando ambiente %s", project.Language), 800 * time.Millisecond},
		{"Configurando banco", fmt.Sprintf("Integrando %s", project.Database), 1200 * time.Millisecond},
		{"Configurando deploy", fmt.Sprintf("Preparando deploy no %s", project.Cloud), 900 * time.Millisecond},
		{"Instalando dependências", "Baixando pacotes necessários", 1500 * time.Millisecond},
		{"Finalizando", "Aplicando configurações finais", 600 * time.Millisecond},
	}
	
	for i, step := range steps {
		// Simple progress simulation
		fmt.Printf("🔄 %s...\n", step.message)
		time.Sleep(step.delay)
		fmt.Printf("✅ %s\n", step.name)
		
		// Update progress bar
		progress := int(float64(i+1) / float64(len(steps)) * 100)
		fmt.Printf("\nProgresso geral: ")
		termx.Progress(100).WithWidth(40).Update(progress)
		fmt.Println()
	}
	
	// Show success
	project.Status = "Criado com sucesso"
	showProjectSummary(project)
}

func showProjectSummary(project *Project) {
	termx.ClearScreen()
	
	fmt.Printf("%s🎉 Projeto criado com sucesso!%s\n", "\033[32m", "\033[0m")
	fmt.Printf("%s═════════════════════════════%s\n\n", "\033[32m", "\033[0m")
	
	// Project summary table
	table := termx.Table([]string{"Configuração", "Valor"})
	table.AddRow("📁 Nome", project.Name)
	table.AddRow("💻 Linguagem", project.Language)
	table.AddRow("🗄️  Banco de dados", project.Database)
	table.AddRow("☁️  Nuvem", project.Cloud)
	table.AddRow("✅ Status", project.Status)
	table.Render()
	
	// Next steps
	fmt.Printf("\n%s🚀 Próximos passos:%s\n", "\033[36m", "\033[0m")
	fmt.Println("• Clone o repositório criado")
	fmt.Println("• Configure as variáveis de ambiente")
	fmt.Println("• Execute o primeiro deploy")
	fmt.Println("• Configure o monitoramento")
	
	fmt.Printf("\n%sPressione qualquer tecla para continuar...%s", "\033[90m", "\033[0m")
	fmt.Scanln()
}

func templateSelector(project *Project) {
	termx.ClearScreen()
	
	fmt.Printf("%s📋 Seletor de Templates%s\n", "\033[34m", "\033[0m")
	fmt.Printf("%s══════════════════════%s\n\n", "\033[34m", "\033[0m")
	
	templates := []string{
		"🚀 API REST com Go + PostgreSQL",
		"⚛️  React SPA + Node.js + MongoDB",
		"🐍 Django + PostgreSQL + Redis",
		"☁️  Microserviços com Kubernetes",
		"📱 Mobile API com FastAPI",
		"🔥 Real-time com WebSocket + Redis",
		"🤖 Bot do Discord com Python",
		"📊 Dashboard Analytics com Next.js",
	}
	
	var selectedTemplate string
	err := termx.Select("Escolha um template:", templates, &selectedTemplate).Run()
	
	if err != nil {
		fmt.Printf("❌ Erro: %v\n", err)
		return
	}
	
	// Simulate template application
	fmt.Printf("🔄 Aplicando template selecionado...\n")
	time.Sleep(2 * time.Second)
	fmt.Printf("✅ Template aplicado com sucesso!\n")
	
	project.Name = "projeto-template"
	project.Status = "Criado a partir de template"
	
	fmt.Printf("\n✨ Template \"%s\" aplicado!\n", selectedTemplate)
	fmt.Printf("%sPressione qualquer tecla para continuar...%s", "\033[90m", "\033[0m")
	fmt.Scanln()
}

func importProject(project *Project) {
	termx.ClearScreen()
	
	fmt.Printf("%s📥 Importar Projeto Existente%s\n", "\033[33m", "\033[0m")
	fmt.Printf("%s═══════════════════════════%s\n\n", "\033[33m", "\033[0m")
	
	var repoUrl string
	var importType string
	
	err := termx.Form().
		Select("Tipo de importação:", []string{
			"Git Repository (HTTPS)",
			"Git Repository (SSH)",
			"Arquivo ZIP",
			"Diretório local",
		}, &importType).
		Input("URL/Caminho do repositório:", &repoUrl).
		Run()
	
	if err != nil {
		fmt.Printf("❌ Erro: %v\n", err)
		return
	}
	
	// Simulate import
	fmt.Printf("🔄 Importando projeto...\n")
	time.Sleep(3 * time.Second)
	fmt.Printf("✅ Projeto importado com sucesso!\n")
	
	fmt.Printf("\n✅ Projeto importado de: %s\n", repoUrl)
	fmt.Printf("%sPressione qualquer tecla para continuar...%s", "\033[90m", "\033[0m")
	fmt.Scanln()
}

func showRecentProjects() {
	termx.ClearScreen()
	
	fmt.Printf("%s🕒 Projetos Recentes%s\n", "\033[36m", "\033[0m")
	fmt.Printf("%s══════════════════%s\n\n", "\033[36m", "\033[0m")
	
	// Recent projects table
	table := termx.Table([]string{"Projeto", "Linguagem", "Última modificação", "Status"}).Interactive()
	table.AddRow("🚀 api-gateway", "Go", "2 horas atrás", "✅ Ativo")
	table.AddRow("⚛️  dashboard-react", "TypeScript", "1 dia atrás", "🔄 Em desenvolvimento")
	table.AddRow("🐍 ml-pipeline", "Python", "3 dias atrás", "⏸️  Pausado")
	table.AddRow("📱 mobile-app", "Swift", "1 semana atrás", "🚀 Deploy")
	table.AddRow("🤖 chat-bot", "Node.js", "2 semanas atrás", "✅ Produção")
	
	selected, err := table.Run()
	if err != nil {
		return
	}
	
	if selected >= 0 {
		projects := []string{"api-gateway", "dashboard-react", "ml-pipeline", "mobile-app", "chat-bot"}
		fmt.Printf("\n📂 Abrindo projeto: %s\n", projects[selected])
		fmt.Printf("%sPressione qualquer tecla para continuar...%s", "\033[90m", "\033[0m")
		fmt.Scanln()
	}
}

func showSettings() {
	termx.ClearScreen()
	
	fmt.Printf("%s⚙️  Configurações Globais%s\n", "\033[35m", "\033[0m")
	fmt.Printf("%s═══════════════════════%s\n\n", "\033[35m", "\033[0m")
	
	var setting string
	err := termx.Select("Configurações:", []string{
		"🎨 Tema",
		"📝 Editor",
		"🔧 Git",
		"☁️  Cloud",
		"💾 Backup",
		"🔔 Notificações",
		"ℹ️  Sobre",
		"↩️  Voltar",
	}, &setting).Run()
	
	if err != nil {
		return
	}
	
	switch setting {
	case "🎨 Tema":
		changeTheme()
	case "📝 Editor":
		configureEditor()
	case "ℹ️  Sobre":
		showAbout()
	default:
		fmt.Printf("⚙️ Configurando: %s\n", setting)
		fmt.Printf("%sPressione qualquer tecla para continuar...%s", "\033[90m", "\033[0m")
		fmt.Scanln()
	}
}

func changeTheme() {
	var theme string
	themes := []string{
		"🌙 Dark (padrão)",
		"☀️  Light",
		"🌊 Ocean Blue",
		"🍃 Forest Green",
		"🔥 Fire Red",
		"💜 Purple Haze",
	}
	
	termx.Select("Escolha um tema:", themes, &theme).Run()
	fmt.Printf("🎨 Tema alterado para: %s\n", theme)
	fmt.Printf("%sPressione qualquer tecla para continuar...%s", "\033[90m", "\033[0m")
	fmt.Scanln()
}

func configureEditor() {
	var editor string
	editors := []string{"VS Code", "Vim", "Neovim", "Emacs", "Sublime Text", "IntelliJ IDEA"}
	
	termx.Select("Editor preferido:", editors, &editor).Run()
	
	fmt.Printf("📝 Editor configurado: %s\n", editor)
	fmt.Printf("%sPressione qualquer tecla para continuar...%s", "\033[90m", "\033[0m")
	fmt.Scanln()
}

func showAbout() {
	termx.ClearScreen()
	
	fmt.Printf("%sℹ️  Sobre o TermX Dashboard%s\n", "\033[36m", "\033[0m")
	fmt.Printf("%s═══════════════════════════%s\n\n", "\033[36m", "\033[0m")
	
	info := termx.Table([]string{"Informação", "Valor"})
	info.AddRow("📦 Versão", "2.0.0")
	info.AddRow("👨‍💻 Desenvolvedor", "TermX Team")
	info.AddRow("📅 Lançamento", "2025")
	info.AddRow("🛠️  Linguagem", "Go 1.21+")
	info.AddRow("📜 Licença", "MIT")
	info.AddRow("🌐 Website", "github.com/vynazevedo/termx")
	info.Render()
	
	fmt.Printf("\n%s💡 TermX Dashboard é um exemplo das capacidades da biblioteca TermX%s\n", "\033[33m", "\033[0m")
	fmt.Printf("%sObrigado por usar nosso software!%s\n\n", "\033[32m", "\033[0m")
	
	fmt.Printf("%sPressione qualquer tecla para continuar...%s", "\033[90m", "\033[0m")
	fmt.Scanln()
}