package main

import (
	"fmt"
	"log"
	"time"
	"github.com/vynazevedo/termx"
	"github.com/vynazevedo/termx/renderer"
)

type Cluster struct {
	Name      string
	Status    string
	Nodes     int
	Pods      int
	CPU       float64
	Memory    float64
	Region    string
}

type Pod struct {
	Name      string
	Namespace string
	Status    string
	Restarts  int
	CPU       string
	Memory    string
}

func main() {
	termx.ClearScreen()
	
	termx.ASCII(termx.KubernetesLogo).WithColor("\033[34m").Render()
	fmt.Println("\nGerenciador de Cluster Kubernetes v1.0")
	fmt.Println("======================================\n")
	
	var action string
	err := termx.Select("O que você gostaria de fazer?", []string{
		"Ver Visão Geral do Cluster",
		"Gerenciar Pods",
		"Escalar Deployment",
		"Ver Uso de Recursos",
		"Implantar Aplicação",
		"Sair",
	}, &action).Run()
	
	if err != nil {
		log.Fatal(err)
	}
	
	switch action {
	case "Ver Visão Geral do Cluster":
		showClusterOverview()
	case "Gerenciar Pods":
		managePods()
	case "Escalar Deployment":
		scaleDeployment()
	case "Ver Uso de Recursos":
		showResourceUsage()
	case "Implantar Aplicação":
		deployApplication()
	case "Sair":
		fmt.Println("Até logo!")
	}
}

func showClusterOverview() {
	termx.ClearScreen()
	
	clusters := []Cluster{
		{"prod-us-leste", "Saudável", 15, 287, 68.5, 82.3, "us-east-1"},
		{"prod-eu-oeste", "Saudável", 12, 203, 45.2, 67.8, "eu-west-1"},
		{"homologacao", "Alerta", 5, 89, 92.1, 78.5, "us-west-2"},
		{"dev", "Saudável", 3, 45, 23.4, 41.2, "us-east-2"},
	}
	
	fmt.Println("Visão Geral do Cluster")
	fmt.Println("======================\n")
	
	table := termx.Table([]string{"Cluster", "Status", "Nós", "Pods", "CPU %", "Memória %", "Região"})
	
	for _, c := range clusters {
		status := c.Status
		if c.Status == "Alerta" {
			status = "\033[33m" + c.Status + "\033[0m"
		} else {
			status = "\033[32m" + c.Status + "\033[0m"
		}
		
		table.AddRow(
			c.Name,
			status,
			fmt.Sprintf("%d", c.Nodes),
			fmt.Sprintf("%d", c.Pods),
			fmt.Sprintf("%.1f%%", c.CPU),
			fmt.Sprintf("%.1f%%", c.Memory),
			c.Region,
		)
	}
	
	table.Render()
	
	var selected string
	termx.Select("\nSelecione o cluster para gerenciar:", []string{
		"prod-us-leste",
		"prod-eu-oeste",
		"homologacao",
		"dev",
		"Voltar ao menu principal",
	}, &selected).Run()
	
	if selected != "Voltar ao menu principal" {
		manageCluster(selected)
	}
}

func managePods() {
	termx.ClearScreen()
	
	pods := []Pod{
		{"api-server-7d9c5b6f9-xvz2k", "default", "Running", 0, "250m", "512Mi"},
		{"frontend-5f8d7c6b5-abc123", "default", "Running", 2, "100m", "256Mi"},
		{"database-0", "default", "Running", 0, "500m", "2Gi"},
		{"cache-redis-8f9d7c6b5-def456", "default", "CrashLoopBackOff", 5, "50m", "128Mi"},
		{"worker-job-2h3j4", "jobs", "Completed", 0, "1000m", "1Gi"},
	}
	
	fmt.Println("Pod Management")
	fmt.Println("==============\n")
	
	table := termx.Table([]string{"Pod Name", "Namespace", "Status", "Restarts", "CPU", "Memory"})
	table.Interactive()
	
	for _, p := range pods {
		status := p.Status
		switch p.Status {
		case "Running":
			status = "\033[32m" + p.Status + "\033[0m"
		case "CrashLoopBackOff":
			status = "\033[31m" + p.Status + "\033[0m"
		case "Completed":
			status = "\033[34m" + p.Status + "\033[0m"
		}
		
		table.AddRow(p.Name, p.Namespace, status, fmt.Sprintf("%d", p.Restarts), p.CPU, p.Memory)
	}
	
	selected, err := table.Run()
	if err == nil && selected >= 0 {
		managePod(pods[selected])
	}
}

func managePod(pod Pod) {
	termx.ClearScreen()
	
	fmt.Printf("Managing Pod: %s\n", pod.Name)
	fmt.Println("==================\n")
	
	var action string
	termx.Select("Select action:", []string{
		"View Logs",
		"Restart Pod",
		"Delete Pod",
		"Describe Pod",
		"Port Forward",
		"Back",
	}, &action).Run()
	
	switch action {
	case "View Logs":
		fmt.Println("\nFetching logs...")
		spinner := termx.Spinner().WithLabel("Loading pod logs")
		spinner.Start()
		time.Sleep(2 * time.Second)
		spinner.Stop()
		
		fmt.Println("Recent logs:")
		fmt.Println("2024-01-15 10:23:45 INFO  Starting application...")
		fmt.Println("2024-01-15 10:23:46 INFO  Connected to database")
		fmt.Println("2024-01-15 10:23:47 INFO  Server listening on :8080")
	
	case "Restart Pod":
		var confirm bool
		termx.Confirm(fmt.Sprintf("Are you sure you want to restart %s?", pod.Name), &confirm).Run()
		
		if confirm {
			progress := termx.Progress(100).WithLabel("Restarting pod")
			for i := 0; i <= 100; i += 10 {
				progress.Update(i)
				time.Sleep(100 * time.Millisecond)
			}
			fmt.Println("\nPod restarted successfully!")
		}
	}
}

func scaleDeployment() {
	termx.ClearScreen()
	
	fmt.Println("Scale Deployment")
	fmt.Println("================\n")
	
	var deployment string
	var replicas string
	
	termx.Form().
		Select("Select deployment:", []string{
			"api-server",
			"frontend",
			"worker",
			"cache-redis",
		}, &deployment).
		Input("Number of replicas:", &replicas).
		Run()
	
	fmt.Printf("\nScaling %s to %s replicas...\n", deployment, replicas)
	
	progress := termx.Progress(100).WithLabel("Scaling deployment")
	for i := 0; i <= 100; i += 10 {
		progress.Update(i)
		time.Sleep(200 * time.Millisecond)
	}
	
	fmt.Println("\nDeployment scaled successfully!")
}

func showResourceUsage() {
	termx.ClearScreen()
	
	fmt.Println("Resource Usage Overview")
	fmt.Println("======================\n")
	
	fmt.Println("CPU Usage (last hour):")
	cpuData := []float64{45, 52, 48, 65, 72, 68, 71, 69, 73, 78, 82, 79}
	termx.Chart(cpuData).
		WithStyle("line").
		WithSize(50, 10).
		Render()
	
	fmt.Println("\nMemory Usage by Namespace:")
	memData := []float64{2.5, 1.8, 3.2, 0.9, 1.5}
	termx.Chart(memData).
		WithLabels([]string{"default", "kube-sys", "monitoring", "ingress", "logging"}).
		WithSize(50, 12).
		Render()
}

func deployApplication() {
	termx.ClearScreen()
	
	fmt.Println("Deploy New Application")
	fmt.Println("=====================\n")
	
	var (
		appName   string
		namespace string
		image     string
		replicas  string
		cpu       string
		memory    string
		expose    bool
		port      string
	)
	
	termx.Form().
		Input("Application name:", &appName).
		Select("Namespace:", []string{"default", "production", "staging", "development"}, &namespace).
		Input("Docker image:", &image).
		Input("Number of replicas:", &replicas).
		Input("CPU request (e.g., 100m):", &cpu).
		Input("Memory request (e.g., 256Mi):", &memory).
		Confirm("Expose as service?", &expose).
		Run()
	
	if expose {
		termx.Input("Service port:", &port).Run()
	}
	
	fmt.Println("\nDeployment Summary:")
	fmt.Println("==================")
	
	summary := termx.Table([]string{"Property", "Value"})
	summary.AddRow("Application", appName)
	summary.AddRow("Namespace", namespace)
	summary.AddRow("Image", image)
	summary.AddRow("Replicas", replicas)
	summary.AddRow("CPU", cpu)
	summary.AddRow("Memory", memory)
	if expose {
		summary.AddRow("Service Port", port)
	}
	summary.Render()
	
	var confirm bool
	termx.Confirm("\nDeploy this application?", &confirm).Run()
	
	if confirm {
		fmt.Println("\nDeploying application...")
		
		steps := []string{
			"Creating namespace",
			"Pulling Docker image",
			"Creating deployment",
			"Starting pods",
			"Creating service",
			"Configuring ingress",
		}
		
		for i, step := range steps {
			if !expose && i >= 4 {
				break
			}
			
			spinner := termx.Spinner().
				WithStyle("dots").
				WithLabel(step)
			spinner.Start()
			time.Sleep(1 * time.Second)
			spinner.Stop()
			fmt.Printf("✓ %s\n", step)
		}
		
		fmt.Println("\nApplication deployed successfully!")
		fmt.Printf("Access your application at: http://%s.cluster.local:%s\n", appName, port)
	}
}

func manageCluster(clusterName string) {
	termx.ClearScreen()
	
	fmt.Printf("Managing Cluster: %s\n", clusterName)
	fmt.Println("====================\n")
	
	layout := termx.Split("horizontal").WithRatio(0.3)
	
	leftBox := termx.BoxLayout("Quick Stats")
	leftBox.SetContent(fmt.Sprintf(
		"Nodes: 15\nPods: 287\nServices: 42\nIngresses: 18\n\nAlerts: 3",
	))
	
	rightBox := termx.BoxLayout("Recent Events")
	rightBox.SetContent(
		"10:45 - Pod api-server scaled up\n" +
		"10:42 - Node k8s-worker-03 joined\n" +
		"10:38 - ConfigMap updated\n" +
		"10:35 - Deployment frontend updated\n" +
		"10:30 - Backup completed",
	)
	
	layout.SetLeft(leftBox).SetRight(rightBox)
	
	r := renderer.New()
	defer r.Close()
	layout.Render(r)
}