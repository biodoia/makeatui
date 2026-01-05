// Example: Using MakeaTUI Agent API
// This demonstrates how an AI agent can programmatically create TUI designs
package main

import (
	"fmt"

	"github.com/makeatui/makeatui/pkg/agent"
)

func main() {
	// Create a new API instance
	api := agent.NewAPI("My Dashboard")

	// Add a title
	api.AddText("title", "📊 Dashboard", 2, 1)

	// Add a main content box
	mainBox := api.AddBox("main", "Welcome to the Dashboard!", 2, 3, 40, 10)
	fmt.Printf("Added main box with ID: %s\n", mainBox)

	// Add navigation buttons
	api.AddButton("btn_home", "🏠 Home", 2, 14)
	api.AddButton("btn_settings", "⚙️ Settings", 15, 14)
	api.AddButton("btn_help", "❓ Help", 30, 14)

	// Add a progress bar
	api.AddProgress("progress", 0.75, 2, 18, 40)

	// Add a list of items
	api.AddList("menu", []string{
		"Dashboard",
		"Analytics",
		"Reports",
		"Settings",
	}, 45, 3, 25, 12)

	// Export to JSON to see the structure
	jsonOutput, err := api.ExportJSON()
	if err != nil {
		fmt.Printf("Error exporting JSON: %v\n", err)
		return
	}

	fmt.Println("\n📋 Canvas JSON:")
	fmt.Println(jsonOutput)

	// Export to Go code
	goCode := api.Export()
	fmt.Println("\n📝 Generated Go Code (first 1000 chars):")
	if len(goCode) > 1000 {
		fmt.Println(goCode[:1000] + "...")
	} else {
		fmt.Println(goCode)
	}

	// Demonstrate undo/redo
	fmt.Println("\n🔄 Testing Undo/Redo:")
	
	// Add another button
	newBtn := api.AddButton("btn_new", "New Button", 50, 14)
	fmt.Printf("Added button: %s\n", newBtn)
	
	// Undo
	if api.Undo() {
		fmt.Println("✓ Undo successful")
	}
	
	// Redo
	if api.Redo() {
		fmt.Println("✓ Redo successful")
	}

	fmt.Println("\n✨ Done! MakeaTUI Agent API demonstration complete.")
}

