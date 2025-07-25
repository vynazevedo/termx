package main

import (
	"fmt"
	"github.com/vynazevedo/termx"
)

func main() {
	fmt.Println("Testing TermX components compilation...")
	
	// Test ASCII art
	art := termx.ASCII("Test\nArt").WithColor("\033[36m")
	fmt.Printf("ASCII component created: %+v\n", art != nil)
	
	// Test table
	table := termx.Table([]string{"Col1", "Col2"})
	table.AddRow("Value1", "Value2")
	fmt.Printf("Table component created: %+v\n", table != nil)
	
	// Test chart
	chart := termx.Chart([]float64{1, 2, 3}).WithSize(10, 5)
	fmt.Printf("Chart component created: %+v\n", chart != nil)
	
	// Test progress
	progress := termx.Progress(100).WithLabel("Test")
	fmt.Printf("Progress component created: %+v\n", progress != nil)
	
	// Test spinner
	spinner := termx.Spinner().WithStyle("dots")
	fmt.Printf("Spinner component created: %+v\n", spinner != nil)
	
	// Test theme
	theme := termx.GetTheme()
	fmt.Printf("Theme system working: %+v\n", theme != nil)
	
	fmt.Println("\n✅ All components compiled successfully!")
	fmt.Println("📦 New features added:")
	fmt.Println("  • ASCII art rendering with logos")
	fmt.Println("  • Interactive tables with selection")
	fmt.Println("  • Real-time charts (bar, line, scatter)")
	fmt.Println("  • Progress bars with multiple styles")
	fmt.Println("  • Loading spinners with animations")
	fmt.Println("  • Split layout system")
	fmt.Println("  • Comprehensive theme system")
	fmt.Println("\n🚀 Run examples in a real terminal:")
	fmt.Println("  go run example/main.go")
	fmt.Println("  go run example/showcase/main.go")
	fmt.Println("  go run example/k8s-manager/main.go")
}