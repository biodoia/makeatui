// Package input - Spinner/Number input (inspired by Lanterna)
package input

import (
	"fmt"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Spinner provides a numeric spinner input
type Spinner struct {
	ID       string
	Value    int
	Min      int
	Max      int
	Step     int
	Width    int
	Focused  bool
	Editable bool
	editing  bool
	editBuf  string
	style    SpinnerStyle
}

// SpinnerStyle holds spinner styling
type SpinnerStyle struct {
	Container lipgloss.Style
	Value     lipgloss.Style
	ValueFocused lipgloss.Style
	Button    lipgloss.Style
	ButtonHover lipgloss.Style
}

// DefaultSpinnerStyle returns default styling
func DefaultSpinnerStyle() SpinnerStyle {
	return SpinnerStyle{
		Container: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7B2CBF")),
		Value: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Align(lipgloss.Center),
		ValueFocused: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E040FB")).
			Bold(true).
			Align(lipgloss.Center),
		Button: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9D4EDD")),
		ButtonHover: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E040FB")).
			Bold(true),
	}
}

// NewSpinner creates a spinner
func NewSpinner(id string, min, max, initial int) *Spinner {
	return &Spinner{
		ID:       id,
		Value:    initial,
		Min:      min,
		Max:      max,
		Step:     1,
		Width:    10,
		Editable: true,
		style:    DefaultSpinnerStyle(),
	}
}

// SetStep sets increment step
func (s *Spinner) SetStep(step int) *Spinner {
	s.Step = step
	return s
}

// SetWidth sets display width
func (s *Spinner) SetWidth(width int) *Spinner {
	s.Width = width
	return s
}

// Increment increases value
func (s *Spinner) Increment() *Spinner {
	s.Value += s.Step
	if s.Value > s.Max {
		s.Value = s.Max
	}
	return s
}

// Decrement decreases value
func (s *Spinner) Decrement() *Spinner {
	s.Value -= s.Step
	if s.Value < s.Min {
		s.Value = s.Min
	}
	return s
}

// Update handles messages
func (s *Spinner) Update(msg tea.Msg) (*Spinner, tea.Cmd) {
	if !s.Focused {
		return s, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if s.editing {
			switch msg.String() {
			case "enter":
				if val, err := strconv.Atoi(s.editBuf); err == nil {
					s.Value = val
					if s.Value < s.Min {
						s.Value = s.Min
					}
					if s.Value > s.Max {
						s.Value = s.Max
					}
				}
				s.editing = false
				return s, func() tea.Msg {
					return SpinnerChangeMsg{ID: s.ID, Value: s.Value}
				}
			case "esc":
				s.editing = false
			case "backspace":
				if len(s.editBuf) > 0 {
					s.editBuf = s.editBuf[:len(s.editBuf)-1]
				}
			default:
				if len(msg.String()) == 1 {
					ch := msg.String()[0]
					if ch >= '0' && ch <= '9' || (ch == '-' && len(s.editBuf) == 0) {
						s.editBuf += msg.String()
					}
				}
			}
		} else {
			switch msg.String() {
			case "up", "k", "+":
				s.Increment()
				return s, func() tea.Msg {
					return SpinnerChangeMsg{ID: s.ID, Value: s.Value}
				}
			case "down", "j", "-":
				s.Decrement()
				return s, func() tea.Msg {
					return SpinnerChangeMsg{ID: s.ID, Value: s.Value}
				}
			case "enter":
				if s.Editable {
					s.editing = true
					s.editBuf = fmt.Sprintf("%d", s.Value)
				}
			}
		}
	}

	return s, nil
}

// View renders the spinner
func (s *Spinner) View() string {
	upBtn := s.style.Button.Render("▲")
	downBtn := s.style.Button.Render("▼")

	valueStyle := s.style.Value
	if s.Focused {
		valueStyle = s.style.ValueFocused
	}

	var valueStr string
	if s.editing {
		valueStr = s.editBuf + "▌"
	} else {
		valueStr = fmt.Sprintf("%d", s.Value)
	}

	value := valueStyle.Width(s.Width - 4).Render(valueStr)

	content := upBtn + " " + value + " " + downBtn
	return s.style.Container.Render(content)
}

// SpinnerChangeMsg is sent when value changes
type SpinnerChangeMsg struct {
	ID    string
	Value int
}

