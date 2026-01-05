// Package input - Slider/Range component (inspired by PyTermGUI)
package input

import (
	"fmt"
	"math"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/makeatui/makeatui/pkg/widgets/mouse"
)

// Slider provides a range slider component
type Slider struct {
	ID          string
	Label       string
	Min         float64
	Max         float64
	Step        float64
	Value       float64
	Width       int
	ShowValue   bool
	ShowMinMax  bool
	Disabled    bool
	Vertical    bool
	trackStyle  lipgloss.Style
	fillStyle   lipgloss.Style
	thumbStyle  lipgloss.Style
	labelStyle  lipgloss.Style
	thumbChar   string
	trackChar   string
	fillChar    string
}

// NewSlider creates a new slider
func NewSlider(id string, min, max float64) *Slider {
	return &Slider{
		ID:        id,
		Min:       min,
		Max:       max,
		Step:      1,
		Value:     min,
		Width:     20,
		ShowValue: true,
		ShowMinMax: false,
		thumbChar: "●",
		trackChar: "─",
		fillChar:  "━",
		trackStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#3C096C")),
		fillStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9D4EDD")),
		thumbStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E040FB")).
			Bold(true),
		labelStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E040FB")).
			Bold(true),
	}
}

// SetLabel sets the label
func (s *Slider) SetLabel(label string) *Slider {
	s.Label = label
	return s
}

// SetStep sets the step increment
func (s *Slider) SetStep(step float64) *Slider {
	s.Step = step
	return s
}

// SetValue sets the current value
func (s *Slider) SetValue(value float64) *Slider {
	s.Value = clamp(value, s.Min, s.Max)
	return s
}

// SetWidth sets the slider width
func (s *Slider) SetWidth(width int) *Slider {
	s.Width = width
	return s
}

// SetShowValue toggles value display
func (s *Slider) SetShowValue(show bool) *Slider {
	s.ShowValue = show
	return s
}

// SetVertical sets vertical orientation
func (s *Slider) SetVertical(vertical bool) *Slider {
	s.Vertical = vertical
	return s
}

// SetDisabled disables the slider
func (s *Slider) SetDisabled(disabled bool) *Slider {
	s.Disabled = disabled
	return s
}

// SetThumbChar sets the thumb character
func (s *Slider) SetThumbChar(char string) *Slider {
	s.thumbChar = char
	return s
}

// Increment increases value by step
func (s *Slider) Increment() *Slider {
	s.Value = clamp(s.Value+s.Step, s.Min, s.Max)
	return s
}

// Decrement decreases value by step
func (s *Slider) Decrement() *Slider {
	s.Value = clamp(s.Value-s.Step, s.Min, s.Max)
	return s
}

// Percentage returns current value as percentage
func (s *Slider) Percentage() float64 {
	return (s.Value - s.Min) / (s.Max - s.Min)
}

// GetZone returns the mouse zone
func (s *Slider) GetZone(x, y, width, height int) *mouse.Zone {
	return &mouse.Zone{
		ID:     s.ID,
		X:      x,
		Y:      y,
		Width:  width,
		Height: height,
		OnDrag: func(dx, dy int) tea.Cmd {
			if s.Disabled {
				return nil
			}
			// Convert drag to value change
			if s.Vertical {
				change := float64(-dy) / float64(height) * (s.Max - s.Min)
				s.SetValue(s.Value + change)
			} else {
				change := float64(dx) / float64(width) * (s.Max - s.Min)
				s.SetValue(s.Value + change)
			}
			return nil
		},
	}
}

// Update handles messages
func (s *Slider) Update(msg tea.Msg) (*Slider, tea.Cmd) {
	if s.Disabled {
		return s, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "right", "l", "up", "k":
			s.Increment()
		case "left", "h", "down", "j":
			s.Decrement()
		case "home":
			s.SetValue(s.Min)
		case "end":
			s.SetValue(s.Max)
		}
	}

	return s, nil
}

// View renders the slider
func (s *Slider) View() string {
	var parts []string

	// Label
	if s.Label != "" {
		parts = append(parts, s.labelStyle.Render(s.Label))
	}

	// Slider track
	fillWidth := int(s.Percentage() * float64(s.Width))
	trackWidth := s.Width - fillWidth - 1 // -1 for thumb

	if fillWidth < 0 {
		fillWidth = 0
	}
	if trackWidth < 0 {
		trackWidth = 0
	}

	fill := s.fillStyle.Render(strings.Repeat(s.fillChar, fillWidth))
	thumb := s.thumbStyle.Render(s.thumbChar)
	track := s.trackStyle.Render(strings.Repeat(s.trackChar, trackWidth))

	slider := fill + thumb + track

	// Min/Max labels
	if s.ShowMinMax {
		minLabel := fmt.Sprintf("%.0f", s.Min)
		maxLabel := fmt.Sprintf("%.0f", s.Max)
		slider = minLabel + " " + slider + " " + maxLabel
	}

	// Value display
	if s.ShowValue {
		valueStr := fmt.Sprintf(" %.1f", s.Value)
		slider += lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Render(valueStr)
	}

	parts = append(parts, slider)

	return strings.Join(parts, "\n")
}

func clamp(value, min, max float64) float64 {
	return math.Min(math.Max(value, min), max)
}

