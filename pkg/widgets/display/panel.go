// Package display - Panel component (inspired by Lanterna)
package display

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Panel is a container with border and title
type Panel struct {
	ID       string
	Title    string
	Content  string
	Width    int
	Height   int
	style    PanelStyle
}

// PanelStyle holds panel styling
type PanelStyle struct {
	Container lipgloss.Style
	Title     lipgloss.Style
	Content   lipgloss.Style
}

// DefaultPanelStyle returns default styling
func DefaultPanelStyle() PanelStyle {
	return PanelStyle{
		Container: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7B2CBF")),
		Title: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E040FB")).
			Bold(true).
			Padding(0, 1),
		Content: lipgloss.NewStyle().
			Padding(0, 1),
	}
}

// NewPanel creates a panel
func NewPanel(id, title string, width, height int) *Panel {
	return &Panel{
		ID:     id,
		Title:  title,
		Width:  width,
		Height: height,
		style:  DefaultPanelStyle(),
	}
}

// SetContent sets panel content
func (p *Panel) SetContent(content string) *Panel {
	p.Content = content
	return p
}

// SetTitle sets panel title
func (p *Panel) SetTitle(title string) *Panel {
	p.Title = title
	return p
}

// SetStyle sets panel style
func (p *Panel) SetStyle(style PanelStyle) *Panel {
	p.style = style
	return p
}

// View renders the panel
func (p *Panel) View() string {
	// Title bar
	title := p.style.Title.Render(p.Title)

	// Content
	contentLines := strings.Split(p.Content, "\n")
	contentHeight := p.Height - 3 // -3 for borders and title

	// Pad or truncate content
	var paddedContent []string
	for i := 0; i < contentHeight; i++ {
		if i < len(contentLines) {
			line := contentLines[i]
			if len(line) > p.Width-4 {
				line = line[:p.Width-4]
			}
			paddedContent = append(paddedContent, line)
		} else {
			paddedContent = append(paddedContent, "")
		}
	}

	content := p.style.Content.Render(strings.Join(paddedContent, "\n"))

	inner := lipgloss.JoinVertical(lipgloss.Left, title, content)
	return p.style.Container.Width(p.Width).Height(p.Height).Render(inner)
}

// ActionList is a list of actions/buttons (inspired by Lanterna)
type ActionList struct {
	ID      string
	Actions []Action
	Layout  ActionLayout
	style   ActionListStyle
}

// Action represents an action button
type Action struct {
	Key     string
	Label   string
	Enabled bool
}

// ActionLayout defines layout direction
type ActionLayout int

const (
	ActionLayoutHorizontal ActionLayout = iota
	ActionLayoutVertical
)

// ActionListStyle holds action list styling
type ActionListStyle struct {
	Container lipgloss.Style
	Action    lipgloss.Style
	ActionDisabled lipgloss.Style
	Key       lipgloss.Style
}

// DefaultActionListStyle returns default styling
func DefaultActionListStyle() ActionListStyle {
	return ActionListStyle{
		Container: lipgloss.NewStyle().
			Padding(0, 1),
		Action: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Padding(0, 1),
		ActionDisabled: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")).
			Padding(0, 1),
		Key: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E040FB")).
			Bold(true),
	}
}

// NewActionList creates an action list
func NewActionList(id string) *ActionList {
	return &ActionList{
		ID:      id,
		Actions: []Action{},
		Layout:  ActionLayoutHorizontal,
		style:   DefaultActionListStyle(),
	}
}

// AddAction adds an action
func (a *ActionList) AddAction(key, label string) *ActionList {
	a.Actions = append(a.Actions, Action{Key: key, Label: label, Enabled: true})
	return a
}

// SetLayout sets layout direction
func (a *ActionList) SetLayout(layout ActionLayout) *ActionList {
	a.Layout = layout
	return a
}

// View renders the action list
func (a *ActionList) View() string {
	var items []string

	for _, action := range a.Actions {
		style := a.style.Action
		if !action.Enabled {
			style = a.style.ActionDisabled
		}

		keyStr := a.style.Key.Render("[" + action.Key + "]")
		labelStr := style.Render(action.Label)
		items = append(items, keyStr+" "+labelStr)
	}

	if a.Layout == ActionLayoutHorizontal {
		return a.style.Container.Render(strings.Join(items, "  "))
	}
	return a.style.Container.Render(lipgloss.JoinVertical(lipgloss.Left, items...))
}

