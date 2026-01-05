// Package feedback - Modal dialog component
package feedback

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ModalButton represents a modal button
type ModalButton struct {
	Label   string
	Value   string
	Primary bool
	Danger  bool
}

// Modal provides a modal dialog
type Modal struct {
	ID          string
	Title       string
	Content     string
	Buttons     []ModalButton
	Visible     bool
	Width       int
	FocusedBtn  int
	Closable    bool
	Overlay     bool
	style       ModalStyle
}

// ModalStyle holds modal styling
type ModalStyle struct {
	Overlay   lipgloss.Style
	Container lipgloss.Style
	Title     lipgloss.Style
	Content   lipgloss.Style
	Button    lipgloss.Style
	ButtonPri lipgloss.Style
	ButtonDng lipgloss.Style
	ButtonFoc lipgloss.Style
	Close     lipgloss.Style
}

// DefaultModalStyle returns default styling
func DefaultModalStyle() ModalStyle {
	return ModalStyle{
		Overlay: lipgloss.NewStyle().
			Background(lipgloss.Color("#000000")),
		Container: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#9D4EDD")).
			Background(lipgloss.Color("#0D0221")).
			Padding(1, 2),
		Title: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E040FB")).
			Bold(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			BorderForeground(lipgloss.Color("#3C096C")).
			MarginBottom(1),
		Content: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			MarginBottom(1),
		Button: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#3C096C")).
			Padding(0, 2).
			MarginRight(1),
		ButtonPri: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#9D4EDD")).
			Padding(0, 2).
			MarginRight(1).
			Bold(true),
		ButtonDng: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#FF6B6B")).
			Padding(0, 2).
			MarginRight(1),
		ButtonFoc: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#000000")).
			Background(lipgloss.Color("#E040FB")).
			Padding(0, 2).
			MarginRight(1).
			Bold(true),
		Close: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")),
	}
}

// NewModal creates a new modal
func NewModal(id, title string) *Modal {
	return &Modal{
		ID:       id,
		Title:    title,
		Buttons:  []ModalButton{},
		Width:    50,
		Closable: true,
		Overlay:  true,
		style:    DefaultModalStyle(),
	}
}

// SetContent sets modal content
func (m *Modal) SetContent(content string) *Modal {
	m.Content = content
	return m
}

// SetWidth sets modal width
func (m *Modal) SetWidth(width int) *Modal {
	m.Width = width
	return m
}

// AddButton adds a button
func (m *Modal) AddButton(label, value string, primary bool) *Modal {
	m.Buttons = append(m.Buttons, ModalButton{
		Label:   label,
		Value:   value,
		Primary: primary,
	})
	return m
}

// AddDangerButton adds a danger button
func (m *Modal) AddDangerButton(label, value string) *Modal {
	m.Buttons = append(m.Buttons, ModalButton{
		Label:  label,
		Value:  value,
		Danger: true,
	})
	return m
}

// Show shows the modal
func (m *Modal) Show() *Modal {
	m.Visible = true
	m.FocusedBtn = 0
	return m
}

// Hide hides the modal
func (m *Modal) Hide() *Modal {
	m.Visible = false
	return m
}

// Update handles messages
func (m *Modal) Update(msg tea.Msg) (*Modal, tea.Cmd) {
	if !m.Visible {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.Closable {
				m.Hide()
			}
		case "tab", "right", "l":
			m.FocusedBtn++
			if m.FocusedBtn >= len(m.Buttons) {
				m.FocusedBtn = 0
			}
		case "shift+tab", "left", "h":
			m.FocusedBtn--
			if m.FocusedBtn < 0 {
				m.FocusedBtn = len(m.Buttons) - 1
			}
		case "enter", " ":
			if m.FocusedBtn >= 0 && m.FocusedBtn < len(m.Buttons) {
				btn := m.Buttons[m.FocusedBtn]
				m.Hide()
				return m, func() tea.Msg {
					return ModalResultMsg{
						ModalID: m.ID,
						Button:  btn.Value,
					}
				}
			}
		}
	}

	return m, nil
}

// View renders the modal
func (m *Modal) View() string {
	if !m.Visible {
		return ""
	}

	// Title
	title := m.style.Title.Width(m.Width).Render(m.Title)

	// Close button
	closeBtn := ""
	if m.Closable {
		closeBtn = m.style.Close.Render(" [×]")
	}

	// Content
	content := m.style.Content.Width(m.Width).Render(m.Content)

	// Buttons
	var buttons []string
	for i, btn := range m.Buttons {
		style := m.style.Button
		if btn.Primary {
			style = m.style.ButtonPri
		}
		if btn.Danger {
			style = m.style.ButtonDng
		}
		if i == m.FocusedBtn {
			style = m.style.ButtonFoc
		}
		buttons = append(buttons, style.Render(btn.Label))
	}
	buttonRow := lipgloss.JoinHorizontal(lipgloss.Center, buttons...)

	// Combine
	modal := lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.JoinHorizontal(lipgloss.Top, title, closeBtn),
		content,
		buttonRow,
	)

	return m.style.Container.Width(m.Width + 4).Render(modal)
}

// ModalResultMsg is sent when a modal button is clicked
type ModalResultMsg struct {
	ModalID string
	Button  string
}

// Confirm creates a confirmation modal
func Confirm(title, message string) *Modal {
	return NewModal("confirm", title).
		SetContent(message).
		AddButton("Cancel", "cancel", false).
		AddButton("OK", "ok", true)
}

// Alert creates an alert modal
func Alert(title, message string) *Modal {
	return NewModal("alert", title).
		SetContent(message).
		AddButton("OK", "ok", true)
}

// Prompt creates a prompt modal (simplified - for full version use with TextInput)
func Prompt(title, message string) *Modal {
	return NewModal("prompt", title).
		SetContent(message).
		AddButton("Cancel", "cancel", false).
		AddButton("Submit", "submit", true)
}

