// Package advanced - Editor View method
package advanced

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// View renders the editor
func (e *Editor) View() string {
	var lines []string

	visibleHeight := e.Height - 2
	gutterWidth := 0
	if e.ShowLineNums {
		gutterWidth = 5
	}
	textWidth := e.Width - gutterWidth - 2

	for i := 0; i < visibleHeight; i++ {
		lineNum := e.ScrollY + i

		if lineNum >= len(e.Lines) {
			// Empty line
			gutter := ""
			if e.ShowLineNums {
				gutter = e.style.Gutter.Width(gutterWidth).Render("~")
			}
			lines = append(lines, gutter+strings.Repeat(" ", textWidth))
			continue
		}

		// Line number
		gutter := ""
		if e.ShowLineNums {
			numStyle := e.style.LineNumber
			if lineNum == e.CursorY {
				numStyle = numStyle.Foreground(lipgloss.Color("#E040FB"))
			}
			gutter = numStyle.Width(gutterWidth-1).Render(fmt.Sprintf("%d", lineNum+1)) + " "
		}

		// Line content
		line := e.Lines[lineNum]

		// Apply horizontal scroll
		if e.ScrollX > 0 && len(line) > e.ScrollX {
			line = line[e.ScrollX:]
		} else if e.ScrollX > 0 {
			line = ""
		}

		// Truncate to fit
		if len(line) > textWidth {
			line = line[:textWidth]
		}

		// Apply cursor
		var renderedLine string
		if lineNum == e.CursorY {
			// Current line with cursor
			cursorPos := e.CursorX - e.ScrollX
			runes := []rune(line)

			for x := 0; x < textWidth; x++ {
				ch := " "
				if x < len(runes) {
					ch = string(runes[x])
				}

				if x == cursorPos {
					renderedLine += e.style.Cursor.Render(ch)
				} else {
					renderedLine += e.style.Text.Render(ch)
				}
			}
		} else {
			// Normal line
			renderedLine = e.style.Text.Render(line)
			// Pad to width
			if len(line) < textWidth {
				renderedLine += strings.Repeat(" ", textWidth-len(line))
			}
		}

		lines = append(lines, gutter+renderedLine)
	}

	// Status bar
	statusStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#666666")).
		Background(lipgloss.Color("#1A0533"))

	modifiedIndicator := ""
	if e.Modified {
		modifiedIndicator = "[+] "
	}

	readOnlyIndicator := ""
	if e.ReadOnly {
		readOnlyIndicator = "[RO] "
	}

	status := statusStyle.Width(e.Width - 2).Render(
		fmt.Sprintf(" %s%sLn %d, Col %d",
			modifiedIndicator,
			readOnlyIndicator,
			e.CursorY+1,
			e.CursorX+1,
		),
	)
	lines = append(lines, status)

	content := lipgloss.JoinVertical(lipgloss.Left, lines...)
	return e.style.Container.Width(e.Width).Render(content)
}

// MessageBox provides a simple message box dialog
type MessageBox struct {
	ID      string
	Title   string
	Message string
	Buttons []string
	Selected int
	Width   int
	style   MessageBoxStyle
}

// MessageBoxStyle holds message box styling
type MessageBoxStyle struct {
	Container lipgloss.Style
	Title     lipgloss.Style
	Message   lipgloss.Style
	Button    lipgloss.Style
	ButtonSelected lipgloss.Style
}

// DefaultMessageBoxStyle returns default styling
func DefaultMessageBoxStyle() MessageBoxStyle {
	return MessageBoxStyle{
		Container: lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(lipgloss.Color("#9D4EDD")).
			Padding(1, 2),
		Title: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E040FB")).
			Bold(true),
		Message: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")),
		Button: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888")).
			Padding(0, 2),
		ButtonSelected: lipgloss.NewStyle().
			Background(lipgloss.Color("#9D4EDD")).
			Foreground(lipgloss.Color("#FFFFFF")).
			Bold(true).
			Padding(0, 2),
	}
}

// NewMessageBox creates a message box
func NewMessageBox(id, title, message string, buttons []string) *MessageBox {
	return &MessageBox{
		ID:       id,
		Title:    title,
		Message:  message,
		Buttons:  buttons,
		Width:    40,
		style:    DefaultMessageBoxStyle(),
	}
}

// View renders the message box
func (m *MessageBox) View() string {
	title := m.style.Title.Render(m.Title)
	message := m.style.Message.Render(m.Message)

	var buttonViews []string
	for i, btn := range m.Buttons {
		style := m.style.Button
		if i == m.Selected {
			style = m.style.ButtonSelected
		}
		buttonViews = append(buttonViews, style.Render(btn))
	}

	buttons := lipgloss.JoinHorizontal(lipgloss.Center, buttonViews...)

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		message,
		"",
		buttons,
	)

	return m.style.Container.Width(m.Width).Render(content)
}

