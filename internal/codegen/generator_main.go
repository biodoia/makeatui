// Package codegen - Main function generation
package codegen

func (g *Generator) generateMain() string {
	return `func main() {
	p := tea.NewProgram(model{}, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
`
}

// GenerateJSON generates JSON representation of the canvas
func (g *Generator) GenerateJSON() string {
	// Use encoding/json in real implementation
	return "{}" // Simplified for now
}

