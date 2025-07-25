package main

import (
	"fmt"
	"log"
	"strings"
	"time"
	"github.com/vynazevedo/termx"
)

func main() {
	termx.ClearScreen()
	
	termx.ASCII(`
╔╦╗┌─┐┬─┐┌┬┐╔═╗ ╦
 ║ ├┤ ├┬┘│││╔═╝ ╚╦╝
 ╩ └─┘┴└─┴ ┴╚═╝  ╩ `).WithColor("\033[36m").Render()
	
	fmt.Println("\nBiblioteca Avançada de Interface Terminal")
	fmt.Println("==========================================\n")
	
	var demo string
	err := termx.Select("Escolha uma demonstração:", []string{
		"Painel de Dados",
		"Monitor do Servidor",
		"Cliente Git",
		"Gerenciador de Banco",
		"Explorador de Arquivos",
		"Gerenciador de Tarefas",
		"Sair",
	}, &demo).Run()
	
	if err != nil {
		log.Fatal(err)
	}
	
	switch demo {
	case "Painel de Dados":
		dataDashboard()
	case "Monitor do Servidor":
		serverMonitor()
	case "Cliente Git":
		gitClient()
	case "Gerenciador de Banco":
		databaseManager()
	case "Explorador de Arquivos":
		fileExplorer()
	case "Gerenciador de Tarefas":
		taskManager()
	case "Sair":
		fmt.Println("Obrigado por testar o TermX!")
	}
}

func dataDashboard() {
	termx.ClearScreen()
	
	fmt.Println("Painel de Análises")
	fmt.Println("==================\n")
	
	fmt.Println("Usuários Ativos Diários:")
	users := []float64{1250, 1380, 1420, 1350, 1580, 1690, 1750}
	termx.Chart(users).
		WithLabels([]string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}).
		WithSize(50, 12).
		Render()
	
	fmt.Println("\nReceita por Produto:")
	revenue := []float64{45000, 38000, 52000, 29000, 41000}
	termx.Chart(revenue).
		WithLabels([]string{"Pro", "Team", "Enterprise", "Starter", "Custom"}).
		WithStyle("bar").
		WithSize(50, 10).
		Render()
	
	fmt.Println("\nRegiões com Melhor Desempenho:")
	table := termx.Table([]string{"Região", "Receita", "Crescimento", "Usuários"})
	table.AddRow("América do Norte", "$2.4M", "+15%", "45,230")
	table.AddRow("Europa", "$1.8M", "+22%", "38,150")
	table.AddRow("Ásia Pacífico", "$1.2M", "+38%", "29,870")
	table.AddRow("América Latina", "$650K", "+45%", "15,340")
	table.Render()
}

func serverMonitor() {
	termx.ClearScreen()
	
	termx.ASCII(termx.ServerRack).WithColor("\033[33m").Render()
	
	fmt.Println("\nMonitoramento do Servidor")
	fmt.Println("========================\n")
	
	servers := []struct {
		name   string
		cpu    float64
		memory float64
		disk   float64
		status string
	}{
		{"web-01", 45.2, 62.8, 71.5, "saudável"},
		{"web-02", 38.9, 55.3, 68.2, "saudável"},
		{"db-01", 78.5, 85.2, 82.1, "alerta"},
		{"cache-01", 23.4, 41.8, 35.6, "saudável"},
		{"worker-01", 91.2, 78.9, 64.3, "crítico"},
	}
	
	for _, server := range servers {
		fmt.Printf("\n%s [%s]\n", server.name, server.status)
		
		fmt.Print("CPU:    ")
		bar := termx.Progress(100).WithWidth(30).WithLabel("")
		bar.Update(int(server.cpu))
		
		fmt.Print("Memory: ")
		bar = termx.Progress(100).WithWidth(30).WithLabel("")
		bar.Update(int(server.memory))
		
		fmt.Print("Disk:   ")
		bar = termx.Progress(100).WithWidth(30).WithLabel("")
		bar.Update(int(server.disk))
	}
	
	fmt.Println("\nMétricas em Tempo Real:")
	spinner := termx.Spinner().WithStyle("dots").WithLabel("Monitorando...")
	spinner.Start()
	time.Sleep(3 * time.Second)
	spinner.Stop()
	
	fmt.Println("✓ Todos os sistemas operacionais")
}

func gitClient() {
	termx.ClearScreen()
	
	fmt.Println("Gerenciador de Repositório Git")
	fmt.Println("===============================\n")
	
	var action string
	termx.Select("O que você gostaria de fazer?", []string{
		"Ver branches",
		"Criar novo branch",
		"Confirmar alterações",
		"Push para remoto",
		"Pull do remoto",
		"Ver log",
	}, &action).Run()
	
	switch action {
	case "Ver branches":
		fmt.Println("\nBranches:")
		table := termx.Table([]string{"Branch", "Último Commit", "Autor", "Status"})
		table.AddRow("* main", "2 horas atrás", "João Silva", "atualizado")
		table.AddRow("  feature/auth", "1 dia atrás", "Maria Santos", "atrás de 3")
		table.AddRow("  bugfix/memory-leak", "3 horas atrás", "Pedro Costa", "adiantado 2")
		table.Render()
		
	case "Criar novo branch":
		var branchName string
		var fromBranch string
		
		termx.Form().
			Input("Nome do branch:", &branchName).
			Select("Criar a partir de:", []string{"main", "develop", "staging"}, &fromBranch).
			Run()
		
		spinner := termx.Spinner().WithLabel("Criando branch...")
		spinner.Start()
		time.Sleep(1 * time.Second)
		spinner.Stop()
		
		fmt.Printf("✓ Branch '%s' criado a partir de '%s'\n", branchName, fromBranch)
		
	case "Confirmar alterações":
		fmt.Println("\nArquivos alterados:")
		files := []string{
			"M  src/main.go",
			"M  src/handler.go",
			"A  src/utils.go",
			"D  old/legacy.go",
		}
		
		for _, f := range files {
			fmt.Println(f)
		}
		
		var message string
		termx.Input("\nMensagem do commit:", &message).Run()
		
		progress := termx.Progress(len(files)).WithLabel("Confirmando arquivos")
		for i := 0; i <= len(files); i++ {
			progress.Update(i)
			time.Sleep(200 * time.Millisecond)
		}
		
		fmt.Println("✓ Alterações confirmadas com sucesso")
	}
}

func databaseManager() {
	termx.ClearScreen()
	
	fmt.Println("Gerenciador de Banco de Dados")
	fmt.Println("=============================\n")
	
	var (
		host     string
		port     string
		database string
		username string
		password string
	)
	
	termx.Form().
		Input("Host:", &host).
		Input("Porta:", &port).
		Select("Banco de Dados:", []string{"usuarios", "produtos", "pedidos", "analitica"}, &database).
		Input("Usuário:", &username).
		Input("Senha:", &password).
		Run()
	
	spinner := termx.Spinner().WithLabel("Conectando ao banco de dados...")
	spinner.Start()
	time.Sleep(2 * time.Second)
	spinner.Stop()
	
	fmt.Println("✓ Conectado com sucesso\n")
	
	fmt.Println("Tabelas no banco de dados:")
	table := termx.Table([]string{"Tabela", "Linhas", "Tamanho", "Última Modificação"}).Interactive()
	table.AddRow("usuarios", "15,234", "2.4 MB", "2 minutos atrás")
	table.AddRow("produtos", "8,921", "15.8 MB", "1 hora atrás")
	table.AddRow("pedidos", "45,123", "38.2 MB", "5 minutos atrás")
	table.AddRow("sessoes", "128,543", "125.6 MB", "agora mesmo")
	
	selected, _ := table.Run()
	if selected >= 0 {
		fmt.Println("\nConsulta executada com sucesso!")
	}
}

func fileExplorer() {
	termx.ClearScreen()
	
	fmt.Println("File Explorer")
	fmt.Println("=============\n")
	
	currentPath := "/home/user/projects"
	
	for {
		fmt.Printf("Current: %s\n\n", currentPath)
		
		files := []string{
			"[DIR] src/",
			"[DIR] tests/",
			"[DIR] docs/",
			"[FILE] README.md",
			"[FILE] main.go",
			"[FILE] go.mod",
			"[FILE] .gitignore",
			"[UP] .. (pai)",
		}
		
		var selected string
		err := termx.Select("Selecione arquivo/pasta:", files, &selected).Run()
		if err != nil {
			break
		}
		
		if selected == "[UP] .. (pai)" {
			currentPath = "/home/usuario"
			termx.ClearScreen()
			fmt.Println("Explorador de Arquivos")
			fmt.Println("=====================\n")
		} else if strings.HasPrefix(selected, "[DIR]") {
			folder := selected[6:len(selected)-1]
			currentPath = currentPath + "/" + folder
			termx.ClearScreen()
			fmt.Println("Explorador de Arquivos")
			fmt.Println("=====================\n")
		} else {
			fmt.Printf("\nArquivo: %s\n", selected[7:])
			fmt.Println("Tamanho: 1.2 KB")
			fmt.Println("Modificado: 2 horas atrás")
			fmt.Println("Permissões: -rw-r--r--")
			
			var action string
			termx.Select("\nAção:", []string{"Abrir", "Editar", "Deletar", "Voltar"}, &action).Run()
			
			if action == "Voltar" {
				termx.ClearScreen()
				fmt.Println("Explorador de Arquivos")
				fmt.Println("=====================\n")
			}
		}
	}
}

func taskManager() {
	termx.ClearScreen()
	
	fmt.Println("Gerenciador de Tarefas")
	fmt.Println("=====================\n")
	
	fmt.Println("Visão Geral do Sistema:")
	fmt.Printf("Uso de CPU: ")
	termx.Progress(100).WithWidth(30).Update(67)
	fmt.Printf("Memória:    ")
	termx.Progress(100).WithWidth(30).Update(82)
	fmt.Printf("E/S Disco: ")
	termx.Progress(100).WithWidth(30).Update(45)
	
	fmt.Println("\nProcessos em Execução:")
	table := termx.Table([]string{"PID", "Nome", "CPU %", "Memória", "Status"}).Interactive()
	table.AddRow("1234", "chrome", "15.2%", "2.1 GB", "Executando")
	table.AddRow("5678", "code", "8.5%", "1.5 GB", "Executando")
	table.AddRow("9012", "docker", "22.1%", "3.2 GB", "Executando")
	table.AddRow("3456", "postgres", "5.8%", "512 MB", "Executando")
	table.AddRow("7890", "node", "12.3%", "768 MB", "Executando")
	
	selected, _ := table.Run()
	if selected >= 0 {
		var action string
		termx.Select("Ação do processo:", []string{
			"Ver detalhes",
			"Finalizar processo",
			"Alterar prioridade",
			"Cancelar",
		}, &action).Run()
		
		if action == "Finalizar processo" {
			var confirm bool
			termx.Confirm("Tem certeza?", &confirm).Run()
			if confirm {
				fmt.Println("\n✓ Processo finalizado")
			}
		}
	}
}