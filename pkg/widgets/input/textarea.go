// Package input - TextArea component (multiline input)
package input

import (
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/makeatui/makeatui/pkg/widgets/mouse"
)

// TextArea provides multiline text input with mouse support
type TextArea struct {
	textarea.Model
	ID           string
	Label        string
	Placeholder  string
	HelpText     string
	ErrorText    string
	Required     bool
	Disabled     bool
	MaxLength    int
	ShowLineNums bool
	labelStyle   lipgloss.Style
	errorStyle   lipgloss.Style
	helpStyle    lipgloss.Style
}

// NewTextArea creates a new text area
func NewTextArea(id string) *TextArea {
	ta := textarea.New()
	ta.ShowLineNumbers = true
	ta.CharLimit = 0 // no limit

	return &TextArea{
		Model:        ta,
		ID:           id,
		ShowLineNums: true,
		labelStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E040FB")).
			Bold(true),
		errorStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF6B6B")),
		helpStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")).
			Italic(true),
	}
}

// SetLabel sets the label
func (t *TextArea) SetLabel(label string) *TextArea {
	t.Label = label
	return t
}

// SetPlaceholder sets placeholder text
func (t *TextArea) SetPlaceholder(text string) *TextArea {
	t.Placeholder = text
	t.Model.Placeholder = text
	return t
}

// SetSize sets dimensions
func (t *TextArea) SetSize(width, height int) *TextArea {
	t.Model.SetWidth(width)
	t.Model.SetHeight(height)
	return t
}

// SetMaxLength sets character limit
func (t *TextArea) SetMaxLength(limit int) *TextArea {
	t.MaxLength = limit
	t.Model.CharLimit = limit
	return t
}

// SetShowLineNumbers toggles line numbers
func (t *TextArea) SetShowLineNumbers(show bool) *TextArea {
	t.ShowLineNums = show
	t.Model.ShowLineNumbers = show
	return t
}

// SetRequired marks as required
func (t *TextArea) SetRequired(required bool) *TextArea {
	t.Required = required
	return t
}

// SetDisabled disables the textarea
func (t *TextArea) SetDisabled(disabled bool) *TextArea {
	t.Disabled = disabled
	return t
}

// SetValue sets the content
func (t *TextArea) SetValue(value string) *TextArea {
	t.Model.SetValue(value)
	return t
}

// LineCount returns number of lines
func (t *TextArea) LineCount() int {
	return len(strings.Split(t.Model.Value(), "\n"))
}

// GetZone returns the mouse zone
func (t *TextArea) GetZone(x, y, width, height int) *mouse.Zone {
	return &mouse.Zone{
		ID:     t.ID,
		X:      x,
		Y:      y,
		Width:  width,
		Height: height,
		OnClick: func() tea.Cmd {
			return func() tea.Msg {
				return FocusInputMsg{ID: t.ID}
			}
		},
	}
}

// Update handles messages
func (t *TextArea) Update(msg tea.Msg) (*TextArea, tea.Cmd) {
	if t.Disabled {
		return t, nil
	}

	var cmd tea.Cmd
	t.Model, cmd = t.Model.Update(msg)
	return t, cmd
}

// View renders the textarea
func (t *TextArea) View() string {
	var parts []string

	// Label
	if t.Label != "" {
		label := t.Label
		if t.Required {
			label += " *"
		}
		parts = append(parts, t.labelStyle.Render(label))
	}

	// TextArea
	parts = append(parts, t.Model.View())

	// Character count
	if t.MaxLength > 0 {
		count := len(t.Model.Value())
		countStr := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")).
			Render(strings.Repeat(" ", 2) + string(rune('0'+count%10)) + "/" + string(rune('0'+t.MaxLength%10)))
		parts = append(parts, countStr)
	}

	// Error
	if t.ErrorText != "" {
		parts = append(parts, t.errorStyle.Render("⚠ "+t.ErrorText))
	}

	// Help
	if t.HelpText != "" && t.ErrorText == "" {
		parts = append(parts, t.helpStyle.Render(t.HelpText))
	}

	return strings.Join(parts, "\n")
}

