// Package navigation provides navigation components
package navigation

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/makeatui/makeatui/pkg/widgets/mouse"
)

// MenuItem represents a menu item
type MenuItem struct {
	ID        string
	Label     string
	Icon      string
	Shortcut  string
	Disabled  bool
	Separator bool
	SubMenu   *Menu
	Action    func() tea.Cmd
}

// Menu provides a menu component (horizontal or vertical)
type Menu struct {
	ID         string
	Items      []*MenuItem
	Selected   int
	Horizontal bool
	Open       bool
	OpenSub    int
	style      MenuStyle
	zoneManager *mouse.ZoneManager
}

// MenuStyle holds menu styling
type MenuStyle struct {
	Container lipgloss.Style
	Item      lipgloss.Style
	ItemSel   lipgloss.Style
	ItemDis   lipgloss.Style
	Icon      lipgloss.Style
	Shortcut  lipgloss.Style
	Separator lipgloss.Style
	SubMenu   lipgloss.Style
}

// DefaultMenuStyle returns default styling
func DefaultMenuStyle() MenuStyle {
	return MenuStyle{
		Container: lipgloss.NewStyle().
			Background(lipgloss.Color("#1A0533")).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7B2CBF")),
		Item: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Padding(0, 2),
		ItemSel: lipgloss.NewStyle().
			Background(lipgloss.Color("#9D4EDD")).
			Foreground(lipgloss.Color("#FFFFFF")).
			Bold(true).
			Padding(0, 2),
		ItemDis: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")).
			Padding(0, 2),
		Icon: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E040FB")).
			MarginRight(1),
		Shortcut: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7B2CBF")).
			Italic(true),
		Separator: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#3C096C")),
		SubMenu: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")),
	}
}

// NewMenu creates a new menu
func NewMenu(id string) *Menu {
	return &Menu{
		ID:          id,
		Items:       []*MenuItem{},
		Selected:    0,
		OpenSub:     -1,
		style:       DefaultMenuStyle(),
		zoneManager: mouse.NewZoneManager(),
	}
}

// AddItem adds a menu item
func (m *Menu) AddItem(id, label, icon, shortcut string) *Menu {
	m.Items = append(m.Items, &MenuItem{
		ID:       id,
		Label:    label,
		Icon:     icon,
		Shortcut: shortcut,
	})
	return m
}

// AddSeparator adds a separator
func (m *Menu) AddSeparator() *Menu {
	m.Items = append(m.Items, &MenuItem{Separator: true})
	return m
}

// AddSubMenu adds a submenu
func (m *Menu) AddSubMenu(id, label, icon string, subMenu *Menu) *Menu {
	m.Items = append(m.Items, &MenuItem{
		ID:      id,
		Label:   label,
		Icon:    icon,
		SubMenu: subMenu,
	})
	return m
}

// SetHorizontal sets horizontal layout
func (m *Menu) SetHorizontal(horizontal bool) *Menu {
	m.Horizontal = horizontal
	return m
}

// GetSelectedItem returns the selected item
func (m *Menu) GetSelectedItem() *MenuItem {
	if m.Selected >= 0 && m.Selected < len(m.Items) {
		return m.Items[m.Selected]
	}
	return nil
}

// findNextSelectable finds next non-separator, non-disabled item
func (m *Menu) findNextSelectable(start, direction int) int {
	for i := 0; i < len(m.Items); i++ {
		idx := (start + direction*(i+1) + len(m.Items)) % len(m.Items)
		item := m.Items[idx]
		if !item.Separator && !item.Disabled {
			return idx
		}
	}
	return start
}

// Update handles messages
func (m *Menu) Update(msg tea.Msg) (*Menu, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if !m.Horizontal {
				m.Selected = m.findNextSelectable(m.Selected, -1)
			}
		case "down", "j":
			if !m.Horizontal {
				m.Selected = m.findNextSelectable(m.Selected, 1)
			}
		case "left", "h":
			if m.Horizontal {
				m.Selected = m.findNextSelectable(m.Selected, -1)
			} else if m.OpenSub >= 0 {
				m.OpenSub = -1
			}
		case "right", "l":
			if m.Horizontal {
				m.Selected = m.findNextSelectable(m.Selected, 1)
			} else if item := m.GetSelectedItem(); item != nil && item.SubMenu != nil {
				m.OpenSub = m.Selected
			}
		case "enter", " ":
			if item := m.GetSelectedItem(); item != nil {
				if item.SubMenu != nil {
					m.OpenSub = m.Selected
				} else if item.Action != nil {
					return m, item.Action()
				}
			}
		case "esc":
			if m.OpenSub >= 0 {
				m.OpenSub = -1
			} else {
				m.Open = false
			}
		}
	}

	return m, nil
}

// View renders the menu
func (m *Menu) View() string {
	var items []string

	for i, item := range m.Items {
		if item.Separator {
			if m.Horizontal {
				items = append(items, m.style.Separator.Render(" │ "))
			} else {
				items = append(items, m.style.Separator.Render(strings.Repeat("─", 20)))
			}
			continue
		}

		style := m.style.Item
		if item.Disabled {
			style = m.style.ItemDis
		} else if i == m.Selected {
			style = m.style.ItemSel
		}

		// Build item content
		var parts []string
		if item.Icon != "" {
			parts = append(parts, m.style.Icon.Render(item.Icon))
		}
		parts = append(parts, item.Label)
		if item.Shortcut != "" {
			parts = append(parts, m.style.Shortcut.Render("  "+item.Shortcut))
		}
		if item.SubMenu != nil {
			parts = append(parts, m.style.SubMenu.Render(" ▶"))
		}

		itemView := style.Render(strings.Join(parts, ""))
		items = append(items, itemView)
	}

	var menu string
	if m.Horizontal {
		menu = lipgloss.JoinHorizontal(lipgloss.Center, items...)
	} else {
		menu = lipgloss.JoinVertical(lipgloss.Left, items...)
		menu = m.style.Container.Render(menu)
	}

	return menu
}

