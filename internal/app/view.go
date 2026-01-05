// Package app - View rendering for the application
package app

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// renderView builds the complete UI
func (m Model) renderView() string {
	// Calculate dimensions
	sidebarWidth := 28
	canvasWidth := m.width - sidebarWidth - 4
	if canvasWidth < 40 {
		canvasWidth = 40
	}

	contentHeight := m.height - 5 // toolbar + statusbar

	// Build each section
	toolbar := m.renderToolbar()
	sidebar := m.renderSidebar(sidebarWidth, contentHeight)
	canvasView := m.renderCanvas(canvasWidth, contentHeight)
	statusBar := m.renderStatusBar()
	helpOverlay := ""
	if m.showHelp {
		helpOverlay = m.renderHelp()
	}

	// Layout: Toolbar on top, then sidebar + canvas side by side, statusbar at bottom
	mainContent := lipgloss.JoinHorizontal(lipgloss.Top, sidebar, canvasView)
	fullView := lipgloss.JoinVertical(lipgloss.Left, toolbar, mainContent, statusBar)

	if m.showHelp {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, helpOverlay)
	}

	return fullView
}

// renderToolbar renders the top toolbar
func (m Model) renderToolbar() string {
	logo := lipgloss.NewStyle().
		Foreground(m.theme.Primary).
		Bold(true).
		Render("✨ MakeaTUI")

	projectName := lipgloss.NewStyle().
		Foreground(m.theme.TextSecondary).
		Render(" │ " + m.projectName)

	spacer := lipgloss.NewStyle().
		Width(m.width - lipgloss.Width(logo) - lipgloss.Width(projectName) - 20).
		Render("")

	actions := lipgloss.NewStyle().
		Foreground(m.theme.TextMuted).
		Render("[?] Help  [e] Export  [q] Quit")

	toolbar := lipgloss.NewStyle().
		Background(m.theme.SurfaceLight).
		Width(m.width).
		Padding(0, 2).
		Render(logo + projectName + spacer + actions)

	return toolbar
}

// renderSidebar renders the component palette
func (m Model) renderSidebar(width, height int) string {
	titleStyle := lipgloss.NewStyle().
		Foreground(m.theme.Primary).
		Bold(true).
		MarginBottom(1)

	title := titleStyle.Render("📦 Components")

	var items []string
	for i, comp := range m.components {
		itemStyle := lipgloss.NewStyle().
			Foreground(m.theme.TextSecondary).
			PaddingLeft(1)

		if i == m.selected && m.focus == FocusSidebar {
			itemStyle = lipgloss.NewStyle().
				Foreground(m.theme.TextPrimary).
				Background(m.theme.SurfaceLight).
				Bold(true).
				PaddingLeft(1)
		}

		cursor := " "
		if i == m.selected && m.focus == FocusSidebar {
			cursor = lipgloss.NewStyle().
				Foreground(m.theme.Accent).
				Render("▸")
		}

		items = append(items, cursor+itemStyle.Render(fmt.Sprintf(" %s %s", comp.Icon, comp.Name)))
	}

	list := strings.Join(items, "\n")
	content := title + "\n\n" + list

	borderColor := m.theme.Border
	if m.focus == FocusSidebar {
		borderColor = m.theme.Primary
	}

	sidebarStyle := lipgloss.NewStyle().
		Width(width).
		Height(height).
		Background(m.theme.Surface).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Padding(1, 2)

	return sidebarStyle.Render(content)
}

// renderCanvas renders the main canvas area
func (m Model) renderCanvas(width, height int) string {
	borderColor := m.theme.Border
	if m.focus == FocusCanvas {
		borderColor = m.theme.Primary
	}

	canvasContent := m.canvas.Render()

	canvasStyle := lipgloss.NewStyle().
		Width(width).
		Height(height).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Padding(0, 1)

	return canvasStyle.Render(canvasContent)
}

// renderStatusBar renders the bottom status bar
func (m Model) renderStatusBar() string {
	focusNames := []string{"SIDEBAR", "CANVAS", "PROPERTIES"}

	leftInfo := lipgloss.NewStyle().
		Foreground(m.theme.TextPrimary).
		Background(m.theme.Primary).
		Padding(0, 2).
		Render(fmt.Sprintf(" %s ", focusNames[m.focus]))

	modeInfo := lipgloss.NewStyle().
		Foreground(m.theme.TextSecondary).
		Background(m.theme.SurfaceLight).
		Padding(0, 2).
		Render(fmt.Sprintf(" Components: %d ", len(m.canvas.Components)))

	cursorInfo := lipgloss.NewStyle().
		Foreground(m.theme.TextMuted).
		Render(fmt.Sprintf(" Cursor: (%d, %d) ", m.canvas.CursorX, m.canvas.CursorY))

	spacerWidth := m.width - lipgloss.Width(leftInfo) - lipgloss.Width(modeInfo) - lipgloss.Width(cursorInfo)
	spacer := lipgloss.NewStyle().
		Background(m.theme.SurfaceLight).
		Width(spacerWidth).
		Render("")

	return lipgloss.JoinHorizontal(lipgloss.Top, leftInfo, modeInfo, spacer, cursorInfo)
}

// renderHelp renders the help overlay
func (m Model) renderHelp() string {
	helpStyle := lipgloss.NewStyle().
		Background(m.theme.Surface).
		BorderStyle(lipgloss.DoubleBorder()).
		BorderForeground(m.theme.Primary).
		Padding(2, 4).
		Width(50)

	title := lipgloss.NewStyle().
		Foreground(m.theme.Primary).
		Bold(true).
		Render("⌨️  Keyboard Shortcuts\n\n")

	shortcuts := `
↑/k, ↓/j     Navigate up/down
←/h, →/l     Navigate left/right
Tab          Switch focus area
Enter/Space  Add selected component
d/Delete     Delete selected
m            Toggle move mode
e            Export to Go code
?            Toggle this help
q/Ctrl+C     Quit
`
	return helpStyle.Render(title + shortcuts)
}

