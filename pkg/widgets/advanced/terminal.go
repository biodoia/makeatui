// Package advanced provides advanced widgets (Terminal, FileExplorer, etc.)
package advanced

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/makeatui/makeatui/pkg/widgets/mouse"
)

// Terminal provides an embedded terminal widget
type Terminal struct {
	ID         string
	Width      int
	Height     int
	Lines      []string
	Input      string
	Prompt     string
	History    []string
	historyIdx int
	scrollY    int
	style      TerminalStyle
	zoneManager *mouse.ZoneManager
}

// TerminalStyle holds terminal styling
type TerminalStyle struct {
	Container lipgloss.Style
	Prompt    lipgloss.Style
	Input     lipgloss.Style
	Output    lipgloss.Style
	Error     lipgloss.Style
	Cursor    lipgloss.Style
}

// DefaultTerminalStyle returns default styling
func DefaultTerminalStyle() TerminalStyle {
	return TerminalStyle{
		Container: lipgloss.NewStyle().
			Background(lipgloss.Color("#0D0221")).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7B2CBF")).
			Padding(0, 1),
		Prompt: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E040FB")).
			Bold(true),
		Input: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")),
		Output: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#CCCCCC")),
		Error: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF6B6B")),
		Cursor: lipgloss.NewStyle().
			Background(lipgloss.Color("#E040FB")).
			Foreground(lipgloss.Color("#000000")),
	}
}

// NewTerminal creates a terminal widget
func NewTerminal(id string, width, height int) *Terminal {
	return &Terminal{
		ID:          id,
		Width:       width,
		Height:      height,
		Lines:       []string{},
		Prompt:      "$ ",
		History:     []string{},
		historyIdx:  -1,
		style:       DefaultTerminalStyle(),
		zoneManager: mouse.NewZoneManager(),
	}
}

// SetPrompt sets the prompt
func (t *Terminal) SetPrompt(prompt string) *Terminal {
	t.Prompt = prompt
	return t
}

// Write adds output to terminal
func (t *Terminal) Write(text string) *Terminal {
	lines := strings.Split(text, "\n")
	t.Lines = append(t.Lines, lines...)
	t.scrollToBottom()
	return t
}

// WriteError adds error output
func (t *Terminal) WriteError(text string) *Terminal {
	t.Lines = append(t.Lines, t.style.Error.Render(text))
	t.scrollToBottom()
	return t
}

// Clear clears the terminal
func (t *Terminal) Clear() *Terminal {
	t.Lines = []string{}
	t.scrollY = 0
	return t
}

// scrollToBottom scrolls to latest output
func (t *Terminal) scrollToBottom() {
	if len(t.Lines) > t.Height-2 {
		t.scrollY = len(t.Lines) - (t.Height - 2)
	}
}

// Submit submits the current input
func (t *Terminal) Submit() (string, tea.Cmd) {
	cmd := t.Input
	if cmd != "" {
		t.History = append(t.History, cmd)
		t.historyIdx = -1
		t.Lines = append(t.Lines, t.style.Prompt.Render(t.Prompt)+t.Input)
	}
	t.Input = ""
	t.scrollToBottom()
	return cmd, nil
}

// GetZone returns the mouse zone
func (t *Terminal) GetZone(x, y, width, height int) *mouse.Zone {
	return &mouse.Zone{
		ID:     t.ID,
		X:      x,
		Y:      y,
		Width:  width,
		Height: height,
	}
}

// Update handles messages
func (t *Terminal) Update(msg tea.Msg) (*Terminal, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			_, cmd := t.Submit()
			return t, cmd
		case "backspace":
			if len(t.Input) > 0 {
				t.Input = t.Input[:len(t.Input)-1]
			}
		case "up":
			if len(t.History) > 0 {
				if t.historyIdx == -1 {
					t.historyIdx = len(t.History) - 1
				} else if t.historyIdx > 0 {
					t.historyIdx--
				}
				t.Input = t.History[t.historyIdx]
			}
		case "down":
			if t.historyIdx >= 0 {
				t.historyIdx++
				if t.historyIdx >= len(t.History) {
					t.historyIdx = -1
					t.Input = ""
				} else {
					t.Input = t.History[t.historyIdx]
				}
			}
		case "ctrl+l":
			t.Clear()
		case "ctrl+c":
			t.Input = ""
		default:
			if len(msg.String()) == 1 {
				t.Input += msg.String()
			}
		}

	case mouse.ScrollMsg:
		if msg.Direction == mouse.ScrollUp {
			t.scrollY--
			if t.scrollY < 0 {
				t.scrollY = 0
			}
		} else {
			t.scrollY++
			maxScroll := len(t.Lines) - (t.Height - 2)
			if maxScroll < 0 {
				maxScroll = 0
			}
			if t.scrollY > maxScroll {
				t.scrollY = maxScroll
			}
		}
	}

	return t, nil
}

// View renders the terminal
func (t *Terminal) View() string {
	var lines []string

	// Output lines
	visibleHeight := t.Height - 2 // -2 for input and border
	endIdx := t.scrollY + visibleHeight
	if endIdx > len(t.Lines) {
		endIdx = len(t.Lines)
	}

	for i := t.scrollY; i < endIdx; i++ {
		line := t.Lines[i]
		if len(line) > t.Width-4 {
			line = line[:t.Width-4]
		}
		lines = append(lines, t.style.Output.Render(line))
	}

	// Pad if needed
	for len(lines) < visibleHeight {
		lines = append(lines, "")
	}

	// Input line
	cursor := t.style.Cursor.Render(" ")
	inputLine := t.style.Prompt.Render(t.Prompt) +
		t.style.Input.Render(t.Input) + cursor
	lines = append(lines, inputLine)

	content := lipgloss.JoinVertical(lipgloss.Left, lines...)
	return t.style.Container.Width(t.Width).Height(t.Height).Render(content)
}

