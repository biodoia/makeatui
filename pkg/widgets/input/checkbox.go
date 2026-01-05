// Package input - Checkbox and Toggle components
package input

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/makeatui/makeatui/pkg/widgets/mouse"
)

// CheckboxIcon defines checkbox icons
type CheckboxIcon struct {
	Checked   string
	Unchecked string
	Mixed     string // for indeterminate state
}

// DefaultCheckboxIcon returns default icons
func DefaultCheckboxIcon() CheckboxIcon {
	return CheckboxIcon{
		Checked:   "☑",
		Unchecked: "☐",
		Mixed:     "▣",
	}
}

// RoundCheckboxIcon returns round icons
func RoundCheckboxIcon() CheckboxIcon {
	return CheckboxIcon{
		Checked:   "◉",
		Unchecked: "○",
		Mixed:     "◐",
	}
}

// Checkbox provides a checkbox component
type Checkbox struct {
	ID          string
	Label       string
	Checked     bool
	Disabled    bool
	Indeterminate bool // mixed state
	Icon        CheckboxIcon
	style       lipgloss.Style
	checkedStyle lipgloss.Style
	labelStyle   lipgloss.Style
}

// NewCheckbox creates a new checkbox
func NewCheckbox(id, label string) *Checkbox {
	return &Checkbox{
		ID:    id,
		Label: label,
		Icon:  DefaultCheckboxIcon(),
		style: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7B2CBF")),
		checkedStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E040FB")),
		labelStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")),
	}
}

// SetChecked sets the checked state
func (c *Checkbox) SetChecked(checked bool) *Checkbox {
	c.Checked = checked
	c.Indeterminate = false
	return c
}

// SetIndeterminate sets mixed state
func (c *Checkbox) SetIndeterminate() *Checkbox {
	c.Indeterminate = true
	return c
}

// SetDisabled disables the checkbox
func (c *Checkbox) SetDisabled(disabled bool) *Checkbox {
	c.Disabled = disabled
	return c
}

// SetIcon sets custom icons
func (c *Checkbox) SetIcon(icon CheckboxIcon) *Checkbox {
	c.Icon = icon
	return c
}

// Toggle toggles the checkbox
func (c *Checkbox) Toggle() *Checkbox {
	if !c.Disabled {
		c.Checked = !c.Checked
		c.Indeterminate = false
	}
	return c
}

// GetZone returns the mouse zone
func (c *Checkbox) GetZone(x, y, width, height int) *mouse.Zone {
	return &mouse.Zone{
		ID:     c.ID,
		X:      x,
		Y:      y,
		Width:  width,
		Height: height,
		OnClick: func() tea.Cmd {
			c.Toggle()
			return nil
		},
	}
}

// Update handles messages
func (c *Checkbox) Update(msg tea.Msg) (*Checkbox, tea.Cmd) {
	if c.Disabled {
		return c, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", " ":
			c.Toggle()
		}
	}

	return c, nil
}

// View renders the checkbox
func (c *Checkbox) View() string {
	var icon string
	var iconStyle lipgloss.Style

	if c.Indeterminate {
		icon = c.Icon.Mixed
		iconStyle = c.style
	} else if c.Checked {
		icon = c.Icon.Checked
		iconStyle = c.checkedStyle
	} else {
		icon = c.Icon.Unchecked
		iconStyle = c.style
	}

	if c.Disabled {
		iconStyle = iconStyle.Foreground(lipgloss.Color("#666666"))
	}

	labelStyle := c.labelStyle
	if c.Disabled {
		labelStyle = labelStyle.Foreground(lipgloss.Color("#666666"))
	}

	return iconStyle.Render(icon) + " " + labelStyle.Render(c.Label)
}

// CheckboxGroup manages multiple checkboxes
type CheckboxGroup struct {
	ID        string
	Label     string
	Checkboxes []*Checkbox
	MinSelect int // minimum selections
	MaxSelect int // maximum selections (0 = unlimited)
	style     lipgloss.Style
}

// NewCheckboxGroup creates a checkbox group
func NewCheckboxGroup(id string) *CheckboxGroup {
	return &CheckboxGroup{
		ID:         id,
		Checkboxes: []*Checkbox{},
		style: lipgloss.NewStyle().
			PaddingLeft(1),
	}
}

// SetLabel sets the group label
func (g *CheckboxGroup) SetLabel(label string) *CheckboxGroup {
	g.Label = label
	return g
}

// Add adds a checkbox
func (g *CheckboxGroup) Add(id, label string) *CheckboxGroup {
	g.Checkboxes = append(g.Checkboxes, NewCheckbox(id, label))
	return g
}

// GetSelected returns selected values
func (g *CheckboxGroup) GetSelected() []string {
	var selected []string
	for _, cb := range g.Checkboxes {
		if cb.Checked {
			selected = append(selected, cb.ID)
		}
	}
	return selected
}

// View renders the group
func (g *CheckboxGroup) View() string {
	var parts []string

	if g.Label != "" {
		labelStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#E040FB"))
		parts = append(parts, labelStyle.Render(g.Label))
	}

	for _, cb := range g.Checkboxes {
		parts = append(parts, g.style.Render(cb.View()))
	}

	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}

