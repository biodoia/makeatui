// Package input provides input components with mouse support
package input

import (
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/makeatui/makeatui/pkg/widgets/mouse"
)

// TextInput wraps bubbles textinput with mouse support and extra features
type TextInput struct {
	textinput.Model
	ID           string
	Label        string
	Placeholder  string
	HelpText     string
	ErrorText    string
	Required     bool
	Disabled     bool
	Masked       bool // password mode
	MaxLength    int
	MinLength    int
	Validator    func(string) error
	zone         *mouse.Zone
	labelStyle   lipgloss.Style
	inputStyle   lipgloss.Style
	errorStyle   lipgloss.Style
	helpStyle    lipgloss.Style
	focusedStyle lipgloss.Style
}

// NewTextInput creates a new text input
func NewTextInput(id string) *TextInput {
	ti := textinput.New()
	ti.Cursor.SetMode(cursor.CursorBlink)

	return &TextInput{
		Model: ti,
		ID:    id,
		labelStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E040FB")).
			Bold(true),
		inputStyle: lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7B2CBF")).
			Padding(0, 1),
		errorStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF6B6B")),
		helpStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")).
			Italic(true),
		focusedStyle: lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#E040FB")).
			Padding(0, 1),
	}
}

// SetLabel sets the label
func (t *TextInput) SetLabel(label string) *TextInput {
	t.Label = label
	return t
}

// SetPlaceholder sets placeholder text
func (t *TextInput) SetPlaceholder(text string) *TextInput {
	t.Placeholder = text
	t.Model.Placeholder = text
	return t
}

// SetPassword enables password masking
func (t *TextInput) SetPassword(masked bool) *TextInput {
	t.Masked = masked
	if masked {
		t.Model.EchoMode = textinput.EchoPassword
		t.Model.EchoCharacter = '●'
	} else {
		t.Model.EchoMode = textinput.EchoNormal
	}
	return t
}

// SetWidth sets the input width
func (t *TextInput) SetWidth(width int) *TextInput {
	t.Model.Width = width
	return t
}

// SetCharLimit sets maximum characters
func (t *TextInput) SetCharLimit(limit int) *TextInput {
	t.MaxLength = limit
	t.Model.CharLimit = limit
	return t
}

// SetRequired marks as required
func (t *TextInput) SetRequired(required bool) *TextInput {
	t.Required = required
	return t
}

// SetValidator sets validation function
func (t *TextInput) SetValidator(fn func(string) error) *TextInput {
	t.Validator = fn
	return t
}

// SetDisabled disables the input
func (t *TextInput) SetDisabled(disabled bool) *TextInput {
	t.Disabled = disabled
	return t
}

// SetHelpText sets help text
func (t *TextInput) SetHelpText(text string) *TextInput {
	t.HelpText = text
	return t
}

// Validate validates the current value
func (t *TextInput) Validate() error {
	if t.Validator != nil {
		return t.Validator(t.Model.Value())
	}
	return nil
}

// IsValid returns whether input is valid
func (t *TextInput) IsValid() bool {
	return t.Validate() == nil
}

// GetZone returns the mouse zone for this input
func (t *TextInput) GetZone(x, y, width, height int) *mouse.Zone {
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
func (t *TextInput) Update(msg tea.Msg) (*TextInput, tea.Cmd) {
	if t.Disabled {
		return t, nil
	}

	var cmd tea.Cmd
	t.Model, cmd = t.Model.Update(msg)

	// Validate on change
	if err := t.Validate(); err != nil {
		t.ErrorText = err.Error()
	} else {
		t.ErrorText = ""
	}

	return t, cmd
}

// View renders the input
func (t *TextInput) View() string {
	var parts []string

	// Label
	if t.Label != "" {
		label := t.Label
		if t.Required {
			label += " *"
		}
		parts = append(parts, t.labelStyle.Render(label))
	}

	// Input field
	style := t.inputStyle
	if t.Model.Focused() {
		style = t.focusedStyle
	}
	if t.Disabled {
		style = style.Foreground(lipgloss.Color("#666666"))
	}

	inputView := style.Render(t.Model.View())
	parts = append(parts, inputView)

	// Error message
	if t.ErrorText != "" {
		parts = append(parts, t.errorStyle.Render("⚠ "+t.ErrorText))
	}

	// Help text
	if t.HelpText != "" && t.ErrorText == "" {
		parts = append(parts, t.helpStyle.Render(t.HelpText))
	}

	return strings.Join(parts, "\n")
}

// FocusInputMsg requests focus on an input
type FocusInputMsg struct {
	ID string
}

