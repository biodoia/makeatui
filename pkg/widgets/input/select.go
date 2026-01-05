// Package input - Select/Dropdown component (inspired by Textual/PyTermGUI)
package input

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/makeatui/makeatui/pkg/widgets/mouse"
)

// SelectOption represents a select option
type SelectOption struct {
	Value    string
	Label    string
	Disabled bool
	Icon     string
}

// Select provides a dropdown selection component
type Select struct {
	ID           string
	Label        string
	Options      []SelectOption
	Selected     int
	Open         bool
	Searchable   bool
	SearchQuery  string
	Placeholder  string
	MaxVisible   int // max visible options when open
	scrollOffset int
	hovered      int
	style        SelectStyle
	zoneManager  *mouse.ZoneManager
}

// SelectStyle holds styling options
type SelectStyle struct {
	Label       lipgloss.Style
	Trigger     lipgloss.Style
	TriggerOpen lipgloss.Style
	Option      lipgloss.Style
	OptionHover lipgloss.Style
	OptionSel   lipgloss.Style
	Dropdown    lipgloss.Style
	Search      lipgloss.Style
}

// DefaultSelectStyle returns default styling
func DefaultSelectStyle() SelectStyle {
	return SelectStyle{
		Label: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E040FB")).
			Bold(true),
		Trigger: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7B2CBF")).
			Padding(0, 1).
			Width(30),
		TriggerOpen: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#E040FB")).
			Padding(0, 1).
			Width(30),
		Option: lipgloss.NewStyle().
			Padding(0, 1),
		OptionHover: lipgloss.NewStyle().
			Background(lipgloss.Color("#3C096C")).
			Foreground(lipgloss.Color("#FFFFFF")).
			Padding(0, 1),
		OptionSel: lipgloss.NewStyle().
			Background(lipgloss.Color("#9D4EDD")).
			Foreground(lipgloss.Color("#FFFFFF")).
			Padding(0, 1).
			Bold(true),
		Dropdown: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7B2CBF")),
		Search: lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#666666")).
			Padding(0, 1),
	}
}

// NewSelect creates a new select component
func NewSelect(id string) *Select {
	return &Select{
		ID:          id,
		Options:     []SelectOption{},
		Selected:    -1,
		MaxVisible:  5,
		hovered:     0,
		style:       DefaultSelectStyle(),
		zoneManager: mouse.NewZoneManager(),
	}
}

// SetLabel sets the label
func (s *Select) SetLabel(label string) *Select {
	s.Label = label
	return s
}

// AddOption adds an option
func (s *Select) AddOption(value, label string) *Select {
	s.Options = append(s.Options, SelectOption{Value: value, Label: label})
	return s
}

// AddOptions adds multiple options
func (s *Select) AddOptions(options ...SelectOption) *Select {
	s.Options = append(s.Options, options...)
	return s
}

// SetOptions replaces all options
func (s *Select) SetOptions(options []SelectOption) *Select {
	s.Options = options
	return s
}

// SetPlaceholder sets placeholder text
func (s *Select) SetPlaceholder(text string) *Select {
	s.Placeholder = text
	return s
}

// SetSearchable enables search
func (s *Select) SetSearchable(searchable bool) *Select {
	s.Searchable = searchable
	return s
}

// SetMaxVisible sets max visible options
func (s *Select) SetMaxVisible(max int) *Select {
	s.MaxVisible = max
	return s
}

// Toggle opens/closes the dropdown
func (s *Select) Toggle() *Select {
	s.Open = !s.Open
	if s.Open {
		s.hovered = s.Selected
		if s.hovered < 0 {
			s.hovered = 0
		}
	}
	return s
}

// SelectIndex selects an option by index
func (s *Select) SelectIndex(index int) *Select {
	if index >= 0 && index < len(s.Options) && !s.Options[index].Disabled {
		s.Selected = index
		s.Open = false
	}
	return s
}

// GetValue returns the selected value
func (s *Select) GetValue() string {
	if s.Selected >= 0 && s.Selected < len(s.Options) {
		return s.Options[s.Selected].Value
	}
	return ""
}

// GetLabel returns the selected label
func (s *Select) GetLabel() string {
	if s.Selected >= 0 && s.Selected < len(s.Options) {
		return s.Options[s.Selected].Label
	}
	return s.Placeholder
}

// filteredOptions returns options matching search
func (s *Select) filteredOptions() []int {
	if s.SearchQuery == "" {
		indices := make([]int, len(s.Options))
		for i := range indices {
			indices[i] = i
		}
		return indices
	}

	var indices []int
	query := strings.ToLower(s.SearchQuery)
	for i, opt := range s.Options {
		if strings.Contains(strings.ToLower(opt.Label), query) {
			indices = append(indices, i)
		}
	}
	return indices
}

