// Package display - ListBox component (inspired by Lanterna)
package display

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/makeatui/makeatui/pkg/widgets/mouse"
)

// ListBox provides a scrollable list with selection
type ListBox struct {
	ID          string
	Width       int
	Height      int
	Items       []ListItem
	Selected    int
	MultiSelect bool
	SelectedSet map[int]bool
	ScrollY     int
	Focused     bool
	style       ListBoxStyle
	zoneManager *mouse.ZoneManager
}

// ListItem represents a list item
type ListItem struct {
	Value string
	Label string
	Icon  string
	Disabled bool
}

// ListBoxStyle holds listbox styling
type ListBoxStyle struct {
	Container lipgloss.Style
	Item      lipgloss.Style
	ItemSelected lipgloss.Style
	ItemHovered lipgloss.Style
	ItemDisabled lipgloss.Style
	ItemChecked lipgloss.Style
	Scrollbar lipgloss.Style
}

// DefaultListBoxStyle returns default styling
func DefaultListBoxStyle() ListBoxStyle {
	return ListBoxStyle{
		Container: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7B2CBF")),
		Item: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Padding(0, 1),
		ItemSelected: lipgloss.NewStyle().
			Background(lipgloss.Color("#9D4EDD")).
			Foreground(lipgloss.Color("#FFFFFF")).
			Bold(true).
			Padding(0, 1),
		ItemHovered: lipgloss.NewStyle().
			Background(lipgloss.Color("#3C096C")).
			Foreground(lipgloss.Color("#FFFFFF")).
			Padding(0, 1),
		ItemDisabled: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")).
			Padding(0, 1),
		ItemChecked: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6BCB77")).
			Padding(0, 1),
		Scrollbar: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#5A189A")),
	}
}

// NewListBox creates a listbox
func NewListBox(id string, width, height int) *ListBox {
	return &ListBox{
		ID:          id,
		Width:       width,
		Height:      height,
		Items:       []ListItem{},
		SelectedSet: make(map[int]bool),
		style:       DefaultListBoxStyle(),
		zoneManager: mouse.NewZoneManager(),
	}
}

// AddItem adds an item
func (l *ListBox) AddItem(value, label string) *ListBox {
	l.Items = append(l.Items, ListItem{Value: value, Label: label})
	return l
}

// AddItemWithIcon adds an item with icon
func (l *ListBox) AddItemWithIcon(value, label, icon string) *ListBox {
	l.Items = append(l.Items, ListItem{Value: value, Label: label, Icon: icon})
	return l
}

// SetItems sets all items
func (l *ListBox) SetItems(items []ListItem) *ListBox {
	l.Items = items
	return l
}

// GetSelected returns selected item
func (l *ListBox) GetSelected() *ListItem {
	if l.Selected >= 0 && l.Selected < len(l.Items) {
		return &l.Items[l.Selected]
	}
	return nil
}

// GetSelectedItems returns all selected items (for multi-select)
func (l *ListBox) GetSelectedItems() []ListItem {
	var selected []ListItem
	for i := range l.Items {
		if l.SelectedSet[i] {
			selected = append(selected, l.Items[i])
		}
	}
	return selected
}

// SetMultiSelect enables/disables multi-select
func (l *ListBox) SetMultiSelect(multi bool) *ListBox {
	l.MultiSelect = multi
	return l
}

// ensureVisible scrolls to keep selection visible
func (l *ListBox) ensureVisible() {
	visibleHeight := l.Height - 2
	if l.Selected < l.ScrollY {
		l.ScrollY = l.Selected
	}
	if l.Selected >= l.ScrollY+visibleHeight {
		l.ScrollY = l.Selected - visibleHeight + 1
	}
}

// Update handles messages
func (l *ListBox) Update(msg tea.Msg) (*ListBox, tea.Cmd) {
	if !l.Focused {
		return l, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			for l.Selected > 0 {
				l.Selected--
				if !l.Items[l.Selected].Disabled {
					break
				}
			}
			l.ensureVisible()
		case "down", "j":
			for l.Selected < len(l.Items)-1 {
				l.Selected++
				if !l.Items[l.Selected].Disabled {
					break
				}
			}
			l.ensureVisible()
		case "home":
			l.Selected = 0
			l.ensureVisible()
		case "end":
			l.Selected = len(l.Items) - 1
			l.ensureVisible()
		case " ":
			if l.MultiSelect {
				l.SelectedSet[l.Selected] = !l.SelectedSet[l.Selected]
			}
		case "enter":
			if !l.Items[l.Selected].Disabled {
				return l, func() tea.Msg {
					return ListBoxSelectMsg{
						ID:    l.ID,
						Index: l.Selected,
						Item:  l.Items[l.Selected],
					}
				}
			}
		}

	case tea.MouseMsg:
		if msg.Action == tea.MouseActionPress {
			clickY := msg.Y - 1 + l.ScrollY
			if clickY >= 0 && clickY < len(l.Items) {
				if !l.Items[clickY].Disabled {
					l.Selected = clickY
					if l.MultiSelect {
						l.SelectedSet[clickY] = !l.SelectedSet[clickY]
					}
				}
			}
		}

	case mouse.ScrollMsg:
		if msg.Direction == mouse.ScrollUp {
			l.ScrollY--
			if l.ScrollY < 0 {
				l.ScrollY = 0
			}
		} else {
			maxScroll := len(l.Items) - (l.Height - 2)
			if maxScroll < 0 {
				maxScroll = 0
			}
			l.ScrollY++
			if l.ScrollY > maxScroll {
				l.ScrollY = maxScroll
			}
		}
	}

	return l, nil
}

// View renders the listbox
func (l *ListBox) View() string {
	var lines []string
	visibleHeight := l.Height - 2

	for i := 0; i < visibleHeight; i++ {
		idx := l.ScrollY + i
		if idx >= len(l.Items) {
			lines = append(lines, strings.Repeat(" ", l.Width-2))
			continue
		}

		item := l.Items[idx]
		style := l.style.Item

		if item.Disabled {
			style = l.style.ItemDisabled
		} else if idx == l.Selected {
			style = l.style.ItemSelected
		} else if l.SelectedSet[idx] {
			style = l.style.ItemChecked
		}

		// Build label
		label := item.Label
		if item.Icon != "" {
			label = item.Icon + " " + label
		}

		// Multi-select checkbox
		if l.MultiSelect && !item.Disabled {
			if l.SelectedSet[idx] {
				label = "☑ " + label
			} else {
				label = "☐ " + label
			}
		}

		// Truncate if needed
		if len(label) > l.Width-4 {
			label = label[:l.Width-7] + "..."
		}

		lines = append(lines, style.Width(l.Width-2).Render(label))
	}

	content := lipgloss.JoinVertical(lipgloss.Left, lines...)
	return l.style.Container.Width(l.Width).Render(content)
}

// ListBoxSelectMsg is sent when item is selected
type ListBoxSelectMsg struct {
	ID    string
	Index int
	Item  ListItem
}

