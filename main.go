// MakeaTUI entry point
package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/makeatui/makeatui/internal/app"
)

const version = "0.1.0"

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "version", "-v", "--version":
			fmt.Printf("MakeaTUI v%s\n", version)
			os.Exit(0)
		case "help", "-h", "--help":
			printHelp()
			os.Exit(0)
		}
	}

	// Print glamorous banner
	printBanner()

	// Start the TUI application
	p := tea.NewProgram(
		app.New(),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func printBanner() {
	purple := lipgloss.NewStyle().Foreground(lipgloss.Color("#9D4EDD"))
	magenta := lipgloss.NewStyle().Foreground(lipgloss.Color("#E040FB"))

	banner := `
  ███╗   ███╗ █████╗ ██╗  ██╗███████╗ █████╗ ████████╗██╗   ██╗██╗
  ████╗ ████║██╔══██╗██║ ██╔╝██╔════╝██╔══██╗╚══██╔══╝██║   ██║██║
  ██╔████╔██║███████║█████╔╝ █████╗  ███████║   ██║   ██║   ██║██║
  ██║╚██╔╝██║██╔══██║██╔═██╗ ██╔══╝  ██╔══██║   ██║   ██║   ██║██║
  ██║ ╚═╝ ██║██║  ██║██║  ██╗███████╗██║  ██║   ██║   ╚██████╔╝██║
  ╚═╝     ╚═╝╚═╝  ╚═╝╚═╝  ╚═╝╚══════╝╚═╝  ╚═╝   ╚═╝    ╚═════╝ ╚═╝`

	fmt.Println(purple.Render(banner))
	fmt.Println(magenta.Render("  ✨ TUI Design Framework for AI Agents ✨"))
	fmt.Println()
}

func printHelp() {
	printBanner()
	help := `
USAGE:
    makeatui [command]

COMMANDS:
    (none)       Start the interactive TUI designer
    version      Show version information  
    help         Show this help message

KEYBOARD SHORTCUTS:
    Tab          Switch focus (sidebar → canvas → properties)
    ↑/k, ↓/j     Navigate up/down
    ←/h, →/l     Navigate left/right
    Enter/Space  Add selected component to canvas
    d/Delete     Delete selected component
    m            Toggle move mode
    e            Export design to Go code
    ?            Toggle help overlay
    q/Ctrl+C     Quit

DESCRIPTION:
    MakeaTUI is a visual TUI designer for AI agents.
    Design glamorous terminal interfaces and export working Go code.

    Workflow for AI agents:
    1. Generate an image of the desired interface
    2. Use MakeaTUI commands to implement the design
    3. Export working Bubble Tea code
`
	fmt.Println(help)
}

