// Package navigation - CommandPalette Update and View methods
package navigation

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Update handles messages for CommandPalette
func (cp *CommandPalette) Update(msg tea.Msg) (*CommandPalette, tea.Cmd) {
	if !cp.Visible {
		return cp, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			cp.Hide()
		case "enter":
			return cp, cp.Execute()
		case "up", "ctrl+p":
			cp.Selected--
			if cp.Selected < 0 {
				cp.Selected = len(cp.Filtered) - 1
			}
		case "down", "ctrl+n":
			cp.Selected++
			if cp.Selected >= len(cp.Filtered) {
				cp.Selected = 0
			}
		case "backspace":
			if len(cp.Query) > 0 {
				cp.Query = cp.Query[:len(cp.Query)-1]
				cp.filter()
				cp.Selected = 0
			}
		default:
			if len(msg.String()) == 1 {
				cp.Query += msg.String()
				cp.filter()
				cp.Selected = 0
			}
		}
	}

	return cp, nil
}

// View renders the command palette
func (cp *CommandPalette) View() string {
	if !cp.Visible {
		return ""
	}

	var parts []string

	// Input field
	inputContent := cp.Query
	if inputContent == "" {
		inputContent = cp.Placeholder
	}
	input := cp.style.Input.Render("> " + inputContent + "█")
	parts = append(parts, input)

	// Results
	if len(cp.Filtered) == 0 {
		parts = append(parts, cp.style.Empty.Render("No commands found"))
	} else {
		for i, cmd := range cp.Filtered {
			style := cp.style.Item
			if i == cp.Selected {
				style = cp.style.ItemSel
			}

			// Build item
			var itemParts []string

			// Icon
			if cmd.Icon != "" {
				itemParts = append(itemParts, cmd.Icon+" ")
			}

			// Category
			if cmd.Category != "" {
				itemParts = append(itemParts, cp.style.Category.Render(cmd.Category+": "))
			}

			// Label
			itemParts = append(itemParts, cp.style.Label.Render(cmd.Label))

			// Description
			if cmd.Description != "" {
				itemParts = append(itemParts, " "+cp.style.Description.Render(cmd.Description))
			}

			// Shortcut
			if cmd.Shortcut != "" {
				itemParts = append(itemParts, " "+cp.style.Shortcut.Render("["+cmd.Shortcut+"]"))
			}

			itemView := style.Render(strings.Join(itemParts, ""))
			parts = append(parts, itemView)
		}
	}

	return cp.style.Container.Render(lipgloss.JoinVertical(lipgloss.Left, parts...))
}

// Tabs provides a tab navigation component
type Tabs struct {
	ID       string
	Items    []TabItem
	Active   int
	style    TabsStyle
}

// TabItem represents a tab
type TabItem struct {
	ID      string
	Label   string
	Icon    string
	Badge   string
	Closable bool
}

// TabsStyle holds tabs styling
type TabsStyle struct {
	Container lipgloss.Style
	Tab       lipgloss.Style
	TabActive lipgloss.Style
	TabHover  lipgloss.Style
	Icon      lipgloss.Style
	Badge     lipgloss.Style
	Close     lipgloss.Style
}

// DefaultTabsStyle returns default styling
func DefaultTabsStyle() TabsStyle {
	return TabsStyle{
		Container: lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			BorderForeground(lipgloss.Color("#3C096C")),
		Tab: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7B2CBF")).
			Padding(0, 2),
		TabActive: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E040FB")).
			Background(lipgloss.Color("#3C096C")).
			Bold(true).
			Padding(0, 2),
		TabHover: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#1A0533")).
			Padding(0, 2),
		Icon: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9D4EDD")).
			MarginRight(1),
		Badge: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#FF6B6B")).
			Padding(0, 1),
		Close: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")).
			MarginLeft(1),
	}
}

// NewTabs creates tabs
func NewTabs(id string) *Tabs {
	return &Tabs{
		ID:    id,
		Items: []TabItem{},
		style: DefaultTabsStyle(),
	}
}

// AddTab adds a tab
func (t *Tabs) AddTab(id, label, icon string) *Tabs {
	t.Items = append(t.Items, TabItem{
		ID:    id,
		Label: label,
		Icon:  icon,
	})
	return t
}

// SetActive sets active tab
func (t *Tabs) SetActive(index int) *Tabs {
	if index >= 0 && index < len(t.Items) {
		t.Active = index
	}
	return t
}

// Next moves to next tab
func (t *Tabs) Next() *Tabs {
	t.Active++
	if t.Active >= len(t.Items) {
		t.Active = 0
	}
	return t
}

// Prev moves to previous tab
func (t *Tabs) Prev() *Tabs {
	t.Active--
	if t.Active < 0 {
		t.Active = len(t.Items) - 1
	}
	return t
}

// Update handles messages
func (t *Tabs) Update(msg tea.Msg) (*Tabs, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "right", "l":
			t.Next()
		case "shift+tab", "left", "h":
			t.Prev()
		}
	}
	return t, nil
}

// View renders the tabs
func (t *Tabs) View() string {
	var tabs []string

	for i, item := range t.Items {
		style := t.style.Tab
		if i == t.Active {
			style = t.style.TabActive
		}

		var parts []string
		if item.Icon != "" {
			parts = append(parts, t.style.Icon.Render(item.Icon))
		}
		parts = append(parts, item.Label)
		if item.Badge != "" {
			parts = append(parts, " "+t.style.Badge.Render(item.Badge))
		}
		if item.Closable {
			parts = append(parts, t.style.Close.Render("×"))
		}

		tabs = append(tabs, style.Render(strings.Join(parts, "")))
	}

	return t.style.Container.Render(lipgloss.JoinHorizontal(lipgloss.Bottom, tabs...))
}

