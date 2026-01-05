// Package codegen generates Go code from canvas designs
package codegen

import (
	"fmt"
	"strings"

	"github.com/makeatui/makeatui/pkg/schema"
)

// Generator generates Go code from a canvas
type Generator struct {
	Canvas schema.Canvas
}

// NewGenerator creates a new code generator
func NewGenerator(canvas schema.Canvas) *Generator {
	return &Generator{Canvas: canvas}
}

// Generate generates complete Go code for the TUI
func (g *Generator) Generate() string {
	var sb strings.Builder

	// Package and imports
	sb.WriteString(`package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

`)

	// Styles
	sb.WriteString(g.generateStyles())

	// Model
	sb.WriteString(g.generateModel())

	// Init
	sb.WriteString(g.generateInit())

	// Update
	sb.WriteString(g.generateUpdate())

	// View
	sb.WriteString(g.generateView())

	// Main
	sb.WriteString(g.generateMain())

	return sb.String()
}

func (g *Generator) generateStyles() string {
	return `// Styles - Ultraviolet theme
var (
	primaryColor   = lipgloss.Color("#9D4EDD")
	secondaryColor = lipgloss.Color("#7B2CBF")
	accentColor    = lipgloss.Color("#E040FB")
	bgColor        = lipgloss.Color("#0D0221")
	surfaceColor   = lipgloss.Color("#1A1333")
	textColor      = lipgloss.Color("#FFFFFF")
	mutedColor     = lipgloss.Color("#6B6B8D")
	borderColor    = lipgloss.Color("#3D2C5E")

	titleStyle = lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true).
		MarginBottom(1)

	boxStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Padding(1, 2)

	buttonStyle = lipgloss.NewStyle().
		Background(surfaceColor).
		Foreground(textColor).
		Padding(0, 3).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(borderColor)

	buttonActiveStyle = lipgloss.NewStyle().
		Background(primaryColor).
		Foreground(textColor).
		Padding(0, 3).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(accentColor).
		Bold(true)
)

`
}

func (g *Generator) generateModel() string {
	var sb strings.Builder
	sb.WriteString("// Model represents the application state\n")
	sb.WriteString("type model struct {\n")
	sb.WriteString("\twidth  int\n")
	sb.WriteString("\theight int\n")

	// Add state for each component
	for i, comp := range g.Canvas.Components {
		switch comp.Type {
		case schema.TypeInput:
			sb.WriteString(fmt.Sprintf("\tinput%d string\n", i))
		case schema.TypeList:
			sb.WriteString(fmt.Sprintf("\tlist%d []string\n", i))
			sb.WriteString(fmt.Sprintf("\tlist%dSelected int\n", i))
		case schema.TypeProgress:
			sb.WriteString(fmt.Sprintf("\tprogress%d float64\n", i))
		}
	}

	sb.WriteString("\tquitting bool\n")
	sb.WriteString("}\n\n")
	return sb.String()
}

func (g *Generator) generateInit() string {
	var sb strings.Builder
	sb.WriteString("func (m model) Init() tea.Cmd {\n")
	sb.WriteString("\treturn nil\n")
	sb.WriteString("}\n\n")
	return sb.String()
}

func (g *Generator) generateUpdate() string {
	return `func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		}
	}
	return m, nil
}

`
}

func (g *Generator) generateView() string {
	var sb strings.Builder
	sb.WriteString("func (m model) View() string {\n")
	sb.WriteString("\tif m.quitting {\n")
	sb.WriteString("\t\treturn \"\"\n")
	sb.WriteString("\t}\n\n")
	sb.WriteString("\tvar content string\n\n")

	// Generate view code for each component
	for i, comp := range g.Canvas.Components {
		sb.WriteString(g.generateComponentView(i, comp))
	}

	sb.WriteString("\treturn content\n")
	sb.WriteString("}\n\n")
	return sb.String()
}

func (g *Generator) generateComponentView(index int, comp schema.Component) string {
	var sb strings.Builder

	switch comp.Type {
	case schema.TypeBox:
		sb.WriteString(fmt.Sprintf("\t// Box: %s\n", comp.Name))
		sb.WriteString(fmt.Sprintf("\tbox%d := boxStyle.Width(%d).Height(%d).Render(%q)\n",
			index, comp.Size.Width, comp.Size.Height, comp.Text))
		sb.WriteString(fmt.Sprintf("\tcontent += box%d + \"\\n\"\n\n", index))

	case schema.TypeText:
		sb.WriteString(fmt.Sprintf("\t// Text: %s\n", comp.Name))
		sb.WriteString(fmt.Sprintf("\ttext%d := lipgloss.NewStyle().Foreground(textColor).Render(%q)\n",
			index, comp.Text))
		sb.WriteString(fmt.Sprintf("\tcontent += text%d + \"\\n\"\n\n", index))

	case schema.TypeButton:
		sb.WriteString(fmt.Sprintf("\t// Button: %s\n", comp.Name))
		sb.WriteString(fmt.Sprintf("\tbutton%d := buttonStyle.Render(%q)\n", index, comp.Text))
		sb.WriteString(fmt.Sprintf("\tcontent += button%d + \"\\n\"\n\n", index))

	default:
		sb.WriteString(fmt.Sprintf("\t// %s: %s (TODO: implement)\n", comp.Type, comp.Name))
		sb.WriteString(fmt.Sprintf("\tcontent += \"[%s: %s]\\n\"\n\n", comp.Type, comp.Name))
	}

	return sb.String()
}

