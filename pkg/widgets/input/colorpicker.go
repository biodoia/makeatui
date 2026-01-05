// Package input - Color Picker component (inspired by PyTermGUI)
package input

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/makeatui/makeatui/pkg/widgets/mouse"
)

// ColorPicker provides color selection
type ColorPicker struct {
	ID          string
	Label       string
	Color       string // hex color
	Palette     []string
	Open        bool
	Selected    int
	Cols        int // columns in palette grid
	ShowHex     bool
	ShowPreview bool
	Disabled    bool
	style       ColorPickerStyle
}

// ColorPickerStyle holds styling
type ColorPickerStyle struct {
	Label    lipgloss.Style
	Preview  lipgloss.Style
	Palette  lipgloss.Style
	Selected lipgloss.Style
}

// DefaultColorPickerStyle returns default styling
func DefaultColorPickerStyle() ColorPickerStyle {
	return ColorPickerStyle{
		Label: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E040FB")).
			Bold(true),
		Preview: lipgloss.NewStyle().
			Width(4).
			Height(2),
		Palette: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7B2CBF")).
			Padding(1),
		Selected: lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(lipgloss.Color("#E040FB")),
	}
}

// DefaultPalette returns a default color palette
func DefaultPalette() []string {
	return []string{
		"#FF6B6B", "#FF8E72", "#FFD93D", "#6BCB77", "#4D96FF",
		"#9D4EDD", "#E040FB", "#F72585", "#3C096C", "#0D0221",
		"#FFFFFF", "#CCCCCC", "#999999", "#666666", "#333333",
		"#000000", "#FF0000", "#00FF00", "#0000FF", "#FFFF00",
	}
}

// UltravioletPalette returns MakeaTUI themed palette
func UltravioletPalette() []string {
	return []string{
		"#0D0221", "#1A0533", "#240046", "#3C096C", "#5A189A",
		"#7B2CBF", "#9D4EDD", "#C77DFF", "#E0AAFF", "#E040FB",
		"#F72585", "#B5179E", "#7209B7", "#560BAD", "#480CA8",
		"#3A0CA3", "#3F37C9", "#4361EE", "#4895EF", "#4CC9F0",
	}
}

// NewColorPicker creates a color picker
func NewColorPicker(id string) *ColorPicker {
	palette := DefaultPalette()
	return &ColorPicker{
		ID:          id,
		Color:       palette[0],
		Palette:     palette,
		Cols:        5,
		ShowHex:     true,
		ShowPreview: true,
		style:       DefaultColorPickerStyle(),
	}
}

// SetLabel sets the label
func (c *ColorPicker) SetLabel(label string) *ColorPicker {
	c.Label = label
	return c
}

// SetColor sets the selected color
func (c *ColorPicker) SetColor(color string) *ColorPicker {
	c.Color = color
	// Find in palette
	for i, col := range c.Palette {
		if col == color {
			c.Selected = i
			break
		}
	}
	return c
}

// SetPalette sets the color palette
func (c *ColorPicker) SetPalette(palette []string) *ColorPicker {
	c.Palette = palette
	return c
}

// SetCols sets grid columns
func (c *ColorPicker) SetCols(cols int) *ColorPicker {
	c.Cols = cols
	return c
}

// Toggle opens/closes the picker
func (c *ColorPicker) Toggle() *ColorPicker {
	c.Open = !c.Open
	return c
}

// SelectColor selects a color by index
func (c *ColorPicker) SelectColor(index int) *ColorPicker {
	if index >= 0 && index < len(c.Palette) {
		c.Selected = index
		c.Color = c.Palette[index]
	}
	return c
}

// GetZone returns the mouse zone
func (c *ColorPicker) GetZone(x, y, width, height int) *mouse.Zone {
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
func (c *ColorPicker) Update(msg tea.Msg) (*ColorPicker, tea.Cmd) {
	if c.Disabled {
		return c, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", " ":
			if c.Open {
				c.SelectColor(c.Selected)
				c.Open = false
			} else {
				c.Toggle()
			}
		case "esc":
			c.Open = false
		case "up", "k":
			if c.Open {
				c.Selected -= c.Cols
				if c.Selected < 0 {
					c.Selected += len(c.Palette)
				}
			}
		case "down", "j":
			if c.Open {
				c.Selected += c.Cols
				if c.Selected >= len(c.Palette) {
					c.Selected %= len(c.Palette)
				}
			}
		case "left", "h":
			if c.Open {
				c.Selected--
				if c.Selected < 0 {
					c.Selected = len(c.Palette) - 1
				}
			}
		case "right", "l":
			if c.Open {
				c.Selected++
				if c.Selected >= len(c.Palette) {
					c.Selected = 0
				}
			}
		}
	}

	return c, nil
}

