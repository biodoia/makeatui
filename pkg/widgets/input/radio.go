// Package input - Radio button component
package input

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/makeatui/makeatui/pkg/widgets/mouse"
)

// RadioIcon defines radio button icons
type RadioIcon struct {
	Selected   string
	Unselected string
}

// DefaultRadioIcon returns default icons
func DefaultRadioIcon() RadioIcon {
	return RadioIcon{
		Selected:   "◉",
		Unselected: "○",
	}
}

// RadioOption represents a radio option
type RadioOption struct {
	Value    string
	Label    string
	Disabled bool
}

// RadioGroup provides a radio button group
type RadioGroup struct {
	ID       string
	Label    string
	Options  []RadioOption
	Selected int
	Focused  int
	Disabled bool
	Horizontal bool
	Icon     RadioIcon
	style    RadioStyle
	zoneManager *mouse.ZoneManager
}

// RadioStyle holds styling
type RadioStyle struct {
	Label       lipgloss.Style
	Option      lipgloss.Style
	OptionFocus lipgloss.Style
	OptionSel   lipgloss.Style
	Icon        lipgloss.Style
	IconSel     lipgloss.Style
}

// DefaultRadioStyle returns default styling
func DefaultRadioStyle() RadioStyle {
	return RadioStyle{
		Label: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E040FB")).
			Bold(true),
		Option: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")),
		OptionFocus: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#3C096C")),
		OptionSel: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E040FB")),
		Icon: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7B2CBF")),
		IconSel: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E040FB")),
	}
}

// NewRadioGroup creates a new radio group
func NewRadioGroup(id string) *RadioGroup {
	return &RadioGroup{
		ID:          id,
		Options:     []RadioOption{},
		Selected:    -1,
		Focused:     0,
		Icon:        DefaultRadioIcon(),
		style:       DefaultRadioStyle(),
		zoneManager: mouse.NewZoneManager(),
	}
}

// SetLabel sets the group label
func (r *RadioGroup) SetLabel(label string) *RadioGroup {
	r.Label = label
	return r
}

// AddOption adds an option
func (r *RadioGroup) AddOption(value, label string) *RadioGroup {
	r.Options = append(r.Options, RadioOption{Value: value, Label: label})
	return r
}

// SetHorizontal sets horizontal layout
func (r *RadioGroup) SetHorizontal(horizontal bool) *RadioGroup {
	r.Horizontal = horizontal
	return r
}

// SetDisabled disables the group
func (r *RadioGroup) SetDisabled(disabled bool) *RadioGroup {
	r.Disabled = disabled
	return r
}

// SetIcon sets custom icons
func (r *RadioGroup) SetIcon(icon RadioIcon) *RadioGroup {
	r.Icon = icon
	return r
}

// Select selects an option by index
func (r *RadioGroup) Select(index int) *RadioGroup {
	if index >= 0 && index < len(r.Options) && !r.Options[index].Disabled {
		r.Selected = index
	}
	return r
}

// SelectByValue selects an option by value
func (r *RadioGroup) SelectByValue(value string) *RadioGroup {
	for i, opt := range r.Options {
		if opt.Value == value && !opt.Disabled {
			r.Selected = i
			break
		}
	}
	return r
}

// GetValue returns the selected value
func (r *RadioGroup) GetValue() string {
	if r.Selected >= 0 && r.Selected < len(r.Options) {
		return r.Options[r.Selected].Value
	}
	return ""
}

// GetZone returns the mouse zone
func (r *RadioGroup) GetZone(x, y, width, height int) *mouse.Zone {
	return &mouse.Zone{
		ID:     r.ID,
		X:      x,
		Y:      y,
		Width:  width,
		Height: height,
	}
}

// Update handles messages
func (r *RadioGroup) Update(msg tea.Msg) (*RadioGroup, tea.Cmd) {
	if r.Disabled {
		return r, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k", "left", "h":
			r.Focused--
			if r.Focused < 0 {
				r.Focused = len(r.Options) - 1
			}
		case "down", "j", "right", "l":
			r.Focused++
			if r.Focused >= len(r.Options) {
				r.Focused = 0
			}
		case "enter", " ":
			r.Select(r.Focused)
		}
	}

	return r, nil
}

// View renders the radio group
func (r *RadioGroup) View() string {
	var parts []string

	// Label
	if r.Label != "" {
		parts = append(parts, r.style.Label.Render(r.Label))
	}

	// Options
	var optionViews []string
	for i, opt := range r.Options {
		icon := r.Icon.Unselected
		iconStyle := r.style.Icon
		labelStyle := r.style.Option

		if i == r.Selected {
			icon = r.Icon.Selected
			iconStyle = r.style.IconSel
			labelStyle = r.style.OptionSel
		}

		if i == r.Focused {
			labelStyle = r.style.OptionFocus
		}

		if opt.Disabled {
			iconStyle = iconStyle.Foreground(lipgloss.Color("#666666"))
			labelStyle = labelStyle.Foreground(lipgloss.Color("#666666"))
		}

		optionView := iconStyle.Render(icon) + " " + labelStyle.Render(opt.Label)
		optionViews = append(optionViews, optionView)
	}

	if r.Horizontal {
		parts = append(parts, lipgloss.JoinHorizontal(lipgloss.Top, interleaveStrings(optionViews, "  ")...))
	} else {
		parts = append(parts, lipgloss.JoinVertical(lipgloss.Left, optionViews...))
	}

	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}

func interleaveStrings(items []string, sep string) []string {
	if len(items) <= 1 {
		return items
	}
	result := make([]string, len(items)*2-1)
	for i, item := range items {
		result[i*2] = item
		if i < len(items)-1 {
			result[i*2+1] = sep
		}
	}
	return result
}

