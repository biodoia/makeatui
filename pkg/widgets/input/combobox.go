// Package input - ComboBox component (inspired by Lanterna)
package input

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ComboBox combines text input with dropdown selection
type ComboBox struct {
	ID          string
	Width       int
	Value       string
	Options     []string
	Filtered    []string
	Selected    int
	Open        bool
	Focused     bool
	MaxVisible  int
	ScrollY     int
	Editable    bool
	style       ComboBoxStyle
}

// ComboBoxStyle holds combobox styling
type ComboBoxStyle struct {
	Container lipgloss.Style
	Input     lipgloss.Style
	InputFocused lipgloss.Style
	Dropdown  lipgloss.Style
	Option    lipgloss.Style
	OptionSelected lipgloss.Style
	Arrow     lipgloss.Style
}

// DefaultComboBoxStyle returns default styling
func DefaultComboBoxStyle() ComboBoxStyle {
	return ComboBoxStyle{
		Container: lipgloss.NewStyle(),
		Input: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7B2CBF")).
			Padding(0, 1),
		InputFocused: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#E040FB")).
			Padding(0, 1),
		Dropdown: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#5A189A")),
		Option: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Padding(0, 1),
		OptionSelected: lipgloss.NewStyle().
			Background(lipgloss.Color("#9D4EDD")).
			Foreground(lipgloss.Color("#FFFFFF")).
			Bold(true).
			Padding(0, 1),
		Arrow: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E040FB")),
	}
}

// NewComboBox creates a combobox
func NewComboBox(id string, width int) *ComboBox {
	return &ComboBox{
		ID:         id,
		Width:      width,
		Options:    []string{},
		Filtered:   []string{},
		MaxVisible: 5,
		Editable:   true,
		style:      DefaultComboBoxStyle(),
	}
}

// SetOptions sets available options
func (c *ComboBox) SetOptions(options []string) *ComboBox {
	c.Options = options
	c.filterOptions()
	return c
}

// SetValue sets current value
func (c *ComboBox) SetValue(value string) *ComboBox {
	c.Value = value
	c.filterOptions()
	return c
}

// SetEditable sets whether text can be edited
func (c *ComboBox) SetEditable(editable bool) *ComboBox {
	c.Editable = editable
	return c
}

// filterOptions filters options based on current value
func (c *ComboBox) filterOptions() {
	if c.Value == "" {
		c.Filtered = c.Options
		return
	}

	c.Filtered = []string{}
	valueLower := strings.ToLower(c.Value)
	for _, opt := range c.Options {
		if strings.Contains(strings.ToLower(opt), valueLower) {
			c.Filtered = append(c.Filtered, opt)
		}
	}
}

// Update handles messages
func (c *ComboBox) Update(msg tea.Msg) (*ComboBox, tea.Cmd) {
	if !c.Focused {
		return c, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if c.Open {
				c.Selected--
				if c.Selected < 0 {
					c.Selected = len(c.Filtered) - 1
				}
				c.ensureVisible()
			}
		case "down":
			if c.Open {
				c.Selected++
				if c.Selected >= len(c.Filtered) {
					c.Selected = 0
				}
				c.ensureVisible()
			} else {
				c.Open = true
			}
		case "enter":
			if c.Open && c.Selected >= 0 && c.Selected < len(c.Filtered) {
				c.Value = c.Filtered[c.Selected]
				c.Open = false
				return c, func() tea.Msg {
					return ComboBoxSelectMsg{ID: c.ID, Value: c.Value}
				}
			}
		case "esc":
			c.Open = false
		case "backspace":
			if c.Editable && len(c.Value) > 0 {
				c.Value = c.Value[:len(c.Value)-1]
				c.filterOptions()
				c.Open = true
			}
		default:
			if c.Editable && len(msg.String()) == 1 {
				c.Value += msg.String()
				c.filterOptions()
				c.Open = true
				c.Selected = 0
			}
		}
	}

	return c, nil
}

// ensureVisible scrolls to keep selection visible
func (c *ComboBox) ensureVisible() {
	if c.Selected < c.ScrollY {
		c.ScrollY = c.Selected
	}
	if c.Selected >= c.ScrollY+c.MaxVisible {
		c.ScrollY = c.Selected - c.MaxVisible + 1
	}
}

// View renders the combobox
func (c *ComboBox) View() string {
	// Input field
	inputStyle := c.style.Input
	if c.Focused {
		inputStyle = c.style.InputFocused
	}

	arrow := c.style.Arrow.Render("▼")
	if c.Open {
		arrow = c.style.Arrow.Render("▲")
	}

	inputWidth := c.Width - 4
	displayValue := c.Value
	if len(displayValue) > inputWidth {
		displayValue = displayValue[:inputWidth-3] + "..."
	}
	displayValue = displayValue + strings.Repeat(" ", inputWidth-len(displayValue))

	input := inputStyle.Width(c.Width - 2).Render(displayValue + " " + arrow)

	if !c.Open {
		return input
	}

	// Dropdown
	var options []string
	for i := c.ScrollY; i < c.ScrollY+c.MaxVisible && i < len(c.Filtered); i++ {
		style := c.style.Option
		if i == c.Selected {
			style = c.style.OptionSelected
		}
		opt := c.Filtered[i]
		if len(opt) > c.Width-4 {
			opt = opt[:c.Width-7] + "..."
		}
		options = append(options, style.Width(c.Width-2).Render(opt))
	}

	dropdown := c.style.Dropdown.Render(lipgloss.JoinVertical(lipgloss.Left, options...))

	return lipgloss.JoinVertical(lipgloss.Left, input, dropdown)
}

// ComboBoxSelectMsg is sent when option is selected
type ComboBoxSelectMsg struct {
	ID    string
	Value string
}

