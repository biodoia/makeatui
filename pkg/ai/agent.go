// Package ai provides an AI agent specialized in TUI design
package ai

import (
	"fmt"
	"strings"

	"github.com/makeatui/makeatui/pkg/agent"
	"github.com/makeatui/makeatui/pkg/templates"
)

// TUIAgent is an AI agent specialized in generating TUI interfaces
type TUIAgent struct {
	api            *agent.API
	templateEngine *templates.TemplateEngine
	context        *DesignContext
}

// DesignContext holds information about the desired TUI design
type DesignContext struct {
	UseCase      string
	Description  string
	Style        string
	Components   []string
	ColorScheme  string
	Width        int
	Height       int
}

// UseCaseType represents common TUI use cases
type UseCaseType string

const (
	UseCaseDashboard UseCaseType = "dashboard"
	UseCaseForm      UseCaseType = "form"
	UseCaseBrowser   UseCaseType = "browser"
	UseCaseChat      UseCaseType = "chat"
	UseCaseMonitor   UseCaseType = "monitor"
	UseCaseEditor    UseCaseType = "editor"
	UseCaseWizard    UseCaseType = "wizard"
	UseCaseSettings  UseCaseType = "settings"
	UseCaseCustom    UseCaseType = "custom"
)

// NewTUIAgent creates a new TUI design agent
func NewTUIAgent() *TUIAgent {
	return &TUIAgent{
		api:            agent.NewAPI("AI-Generated TUI"),
		templateEngine: templates.NewTemplateEngine(),
		context:        &DesignContext{Width: 80, Height: 24},
	}
}

// SetContext sets the design context
func (a *TUIAgent) SetContext(ctx *DesignContext) {
	a.context = ctx
}

// AnalyzeDescription analyzes a natural language description
func (a *TUIAgent) AnalyzeDescription(description string) UseCaseType {
	desc := strings.ToLower(description)

	switch {
	case containsAny(desc, "dashboard", "metrics", "stats", "overview"):
		return UseCaseDashboard
	case containsAny(desc, "form", "input", "registration", "signup"):
		return UseCaseForm
	case containsAny(desc, "file", "browse", "explorer", "list"):
		return UseCaseBrowser
	case containsAny(desc, "chat", "message", "conversation"):
		return UseCaseChat
	case containsAny(desc, "monitor", "log", "process", "system"):
		return UseCaseMonitor
	case containsAny(desc, "edit", "code", "text", "write"):
		return UseCaseEditor
	case containsAny(desc, "wizard", "step", "setup", "install"):
		return UseCaseWizard
	case containsAny(desc, "settings", "config", "preferences", "options"):
		return UseCaseSettings
	default:
		return UseCaseCustom
	}
}

// GenerateFromDescription generates a TUI from natural language
func (a *TUIAgent) GenerateFromDescription(description string) (*agent.API, error) {
	useCase := a.AnalyzeDescription(description)
	return a.GenerateFromUseCase(useCase, description)
}

// GenerateFromUseCase generates a TUI based on use case
func (a *TUIAgent) GenerateFromUseCase(useCase UseCaseType, title string) (*agent.API, error) {
	a.api.Clear()

	switch useCase {
	case UseCaseDashboard:
		a.generateDashboard(title)
	case UseCaseForm:
		a.generateForm(title)
	case UseCaseBrowser:
		a.generateBrowser(title)
	case UseCaseChat:
		a.generateChat(title)
	case UseCaseMonitor:
		a.generateMonitor(title)
	case UseCaseWizard:
		a.generateWizard(title)
	case UseCaseSettings:
		a.generateSettings(title)
	default:
		a.generateCustom(title)
	}

	return a.api, nil
}

func (a *TUIAgent) generateDashboard(title string) {
	w, h := a.context.Width, a.context.Height

	// Header
	a.api.AddBox("header", title, 0, 0, w, 3)

	// Sidebar
	sidebarW := w / 4
	a.api.AddBox("sidebar", "Navigation", 0, 4, sidebarW, h-5)
	a.api.AddList("nav-menu", []string{"Home", "Analytics", "Reports", "Settings"}, 1, 6, sidebarW-2, 8)

	// Main content area
	mainX := sidebarW + 1
	mainW := w - sidebarW - 1
	a.api.AddBox("main", "Dashboard", mainX, 4, mainW, h-5)

	// Metrics row
	cardW := (mainW - 4) / 3
	a.api.AddBox("card1", "Users", mainX+1, 6, cardW, 4)
	a.api.AddBox("card2", "Revenue", mainX+cardW+2, 6, cardW, 4)
	a.api.AddBox("card3", "Growth", mainX+cardW*2+3, 6, cardW, 4)

	// Chart area
	a.api.AddBox("chart", "Activity", mainX+1, 11, mainW-2, h-13)
}

func (a *TUIAgent) generateForm(title string) {
	w, h := a.context.Width, a.context.Height

	formW := w * 3 / 4
	formH := h - 4
	formX := (w - formW) / 2
	formY := 2

	a.api.AddBox("form", title, formX, formY, formW, formH)

	// Form fields
	labelX := formX + 3
	inputX := formX + 15

	a.api.AddText("label1", "Name:", labelX, formY+3)
	a.api.AddBox("input1", "", inputX, formY+2, formW-20, 3)

	a.api.AddText("label2", "Email:", labelX, formY+6)
	a.api.AddBox("input2", "", inputX, formY+5, formW-20, 3)

	a.api.AddText("label3", "Password:", labelX, formY+9)
	a.api.AddBox("input3", "", inputX, formY+8, formW-20, 3)

	// Buttons
	btnY := formY + formH - 4
	a.api.AddButton("submit", "Submit", formX+formW/2-15, btnY)
	a.api.AddButton("cancel", "Cancel", formX+formW/2+5, btnY)
}

func (a *TUIAgent) generateBrowser(title string) {
	w, h := a.context.Width, a.context.Height

	listW := w * 2 / 5
	a.api.AddBox("list", title, 0, 0, listW, h-1)
	a.api.AddList("items", []string{"Item 1", "Item 2", "Item 3"}, 1, 2, listW-2, h-5)

	a.api.AddBox("preview", "Preview", listW+1, 0, w-listW-1, h-1)
	a.api.AddText("preview-text", "Select an item to preview", listW+3, 3)
}

func (a *TUIAgent) generateChat(title string) {
	w, h := a.context.Width, a.context.Height

	a.api.AddBox("messages", title, 0, 0, w, h-4)
	a.api.AddBox("input", "", 0, h-3, w, 3)
	a.api.AddText("prompt", "> Type your message...", 2, h-2)
}

func (a *TUIAgent) generateMonitor(title string) {
	w, h := a.context.Width, a.context.Height

	half := w / 2
	topH := h * 2 / 3

	a.api.AddBox("processes", "Processes", 0, 0, half, topH)
	a.api.AddBox("logs", "Logs", half+1, 0, w-half-1, topH)
	a.api.AddBox("resources", "System Resources", 0, topH+1, w, h-topH-2)

	a.api.AddProgress("cpu", 0.45, 2, topH+3, 25)
	a.api.AddProgress("mem", 0.62, 30, topH+3, 25)
	a.api.AddProgress("disk", 0.33, 58, topH+3, 20)
}

func (a *TUIAgent) generateWizard(title string) {
	w, h := a.context.Width, a.context.Height

	a.api.AddBox("container", title, 5, 1, w-10, h-3)
	a.api.AddText("step", "Step 1 of 3", 8, 3)
	a.api.AddProgress("progress", 0.33, 8, 5, w-20)
	a.api.AddBox("content", "", 8, 7, w-18, h-12)
	a.api.AddButton("prev", "Previous", 8, h-4)
	a.api.AddButton("next", "Next", w-20, h-4)
}

func (a *TUIAgent) generateSettings(title string) {
	w, h := a.context.Width, a.context.Height

	sideW := w / 3
	a.api.AddBox("categories", "Categories", 0, 0, sideW, h-1)
	a.api.AddList("cat-list", []string{"General", "Appearance", "Shortcuts", "Advanced"}, 1, 2, sideW-2, h-5)

	a.api.AddBox("options", title, sideW+1, 0, w-sideW-1, h-1)
}

func (a *TUIAgent) generateCustom(title string) {
	w, h := a.context.Width, a.context.Height
	a.api.AddBox("main", title, 0, 0, w, h-1)
	a.api.AddText("content", "Custom TUI - Add your components", 2, 2)
}

// Export generates the Go code
func (a *TUIAgent) Export() string {
	return a.api.Export()
}

// GetAPI returns the underlying API
func (a *TUIAgent) GetAPI() *agent.API {
	return a.api
}

func containsAny(s string, substrs ...string) bool {
	for _, sub := range substrs {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}

// Suggestion represents a design suggestion
type Suggestion struct {
	Type        string
	Description string
	Action      func(*agent.API)
}

// GetSuggestions returns design suggestions based on context
func (a *TUIAgent) GetSuggestions() []Suggestion {
	suggestions := []Suggestion{
		{
			Type:        "add_header",
			Description: "Add a header bar at the top",
			Action: func(api *agent.API) {
				api.AddBox("header", "Header", 0, 0, a.context.Width, 3)
			},
		},
		{
			Type:        "add_sidebar",
			Description: "Add a navigation sidebar",
			Action: func(api *agent.API) {
				api.AddBox("sidebar", "Menu", 0, 3, 20, a.context.Height-4)
			},
		},
		{
			Type:        "add_footer",
			Description: "Add a status bar at the bottom",
			Action: func(api *agent.API) {
				api.AddBox("footer", "Status", 0, a.context.Height-2, a.context.Width, 1)
			},
		},
	}

	return suggestions
}

// ApplySuggestion applies a suggestion
func (a *TUIAgent) ApplySuggestion(suggestion Suggestion) {
	suggestion.Action(a.api)
}

// DescribeLayout returns a description of the current layout
func (a *TUIAgent) DescribeLayout() string {
	components := a.api.ListComponents()
	if len(components) == 0 {
		return "Empty canvas"
	}

	var desc strings.Builder
	desc.WriteString(fmt.Sprintf("Layout with %d components:\n", len(components)))

	for _, c := range components {
		desc.WriteString(fmt.Sprintf("- %s (%s) at (%d,%d) size %dx%d\n",
			c.Name, c.Type, c.Position.X, c.Position.Y, c.Size.Width, c.Size.Height))
	}

	return desc.String()
}

