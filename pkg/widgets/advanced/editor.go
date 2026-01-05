// Package advanced - Text Editor widget (inspired by r3bl_tui and Lanterna)
package advanced

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/makeatui/makeatui/pkg/widgets/mouse"
)

// Editor provides a full-featured text editor
type Editor struct {
	ID            string
	Width         int
	Height        int
	Lines         []string
	CursorX       int
	CursorY       int
	ScrollX       int
	ScrollY       int
	Selection     *Selection
	ShowLineNums  bool
	SyntaxHighlight bool
	ReadOnly      bool
	Modified      bool
	TabSize       int
	style         EditorStyle
	zoneManager   *mouse.ZoneManager
}

// Selection represents text selection
type Selection struct {
	StartX int
	StartY int
	EndX   int
	EndY   int
}

// EditorStyle holds editor styling
type EditorStyle struct {
	Container  lipgloss.Style
	LineNumber lipgloss.Style
	Cursor     lipgloss.Style
	Selection  lipgloss.Style
	Text       lipgloss.Style
	CurrentLine lipgloss.Style
	Gutter     lipgloss.Style
}

// DefaultEditorStyle returns default styling
func DefaultEditorStyle() EditorStyle {
	return EditorStyle{
		Container: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7B2CBF")),
		LineNumber: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")).
			Width(4).
			Align(lipgloss.Right),
		Cursor: lipgloss.NewStyle().
			Background(lipgloss.Color("#E040FB")).
			Foreground(lipgloss.Color("#000000")),
		Selection: lipgloss.NewStyle().
			Background(lipgloss.Color("#5A189A")),
		Text: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")),
		CurrentLine: lipgloss.NewStyle().
			Background(lipgloss.Color("#1A0533")),
		Gutter: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#3C096C")),
	}
}

// NewEditor creates a text editor
func NewEditor(id string, width, height int) *Editor {
	return &Editor{
		ID:            id,
		Width:         width,
		Height:        height,
		Lines:         []string{""},
		ShowLineNums:  true,
		TabSize:       4,
		style:         DefaultEditorStyle(),
		zoneManager:   mouse.NewZoneManager(),
	}
}

// SetContent sets editor content
func (e *Editor) SetContent(content string) *Editor {
	e.Lines = strings.Split(content, "\n")
	if len(e.Lines) == 0 {
		e.Lines = []string{""}
	}
	e.CursorX = 0
	e.CursorY = 0
	e.Modified = false
	return e
}

// GetContent returns editor content
func (e *Editor) GetContent() string {
	return strings.Join(e.Lines, "\n")
}

// SetReadOnly sets read-only mode
func (e *Editor) SetReadOnly(readOnly bool) *Editor {
	e.ReadOnly = readOnly
	return e
}

// currentLine returns the current line
func (e *Editor) currentLine() string {
	if e.CursorY >= 0 && e.CursorY < len(e.Lines) {
		return e.Lines[e.CursorY]
	}
	return ""
}

// insertChar inserts a character at cursor
func (e *Editor) insertChar(ch rune) {
	if e.ReadOnly {
		return
	}
	line := e.currentLine()
	if e.CursorX > len(line) {
		e.CursorX = len(line)
	}
	e.Lines[e.CursorY] = line[:e.CursorX] + string(ch) + line[e.CursorX:]
	e.CursorX++
	e.Modified = true
}

// deleteChar deletes character before cursor
func (e *Editor) deleteChar() {
	if e.ReadOnly {
		return
	}
	if e.CursorX > 0 {
		line := e.currentLine()
		e.Lines[e.CursorY] = line[:e.CursorX-1] + line[e.CursorX:]
		e.CursorX--
		e.Modified = true
	} else if e.CursorY > 0 {
		// Join with previous line
		prevLine := e.Lines[e.CursorY-1]
		e.CursorX = len(prevLine)
		e.Lines[e.CursorY-1] = prevLine + e.currentLine()
		e.Lines = append(e.Lines[:e.CursorY], e.Lines[e.CursorY+1:]...)
		e.CursorY--
		e.Modified = true
	}
}

// newLine inserts a new line
func (e *Editor) newLine() {
	if e.ReadOnly {
		return
	}
	line := e.currentLine()
	before := line[:e.CursorX]
	after := line[e.CursorX:]

	e.Lines[e.CursorY] = before
	newLines := make([]string, 0, len(e.Lines)+1)
	newLines = append(newLines, e.Lines[:e.CursorY+1]...)
	newLines = append(newLines, after)
	newLines = append(newLines, e.Lines[e.CursorY+1:]...)
	e.Lines = newLines

	e.CursorY++
	e.CursorX = 0
	e.Modified = true
}

// ensureVisible scrolls to keep cursor visible
func (e *Editor) ensureVisible() {
	visibleHeight := e.Height - 2

	if e.CursorY < e.ScrollY {
		e.ScrollY = e.CursorY
	}
	if e.CursorY >= e.ScrollY+visibleHeight {
		e.ScrollY = e.CursorY - visibleHeight + 1
	}
}

// Update handles messages
func (e *Editor) Update(msg tea.Msg) (*Editor, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if e.CursorY > 0 {
				e.CursorY--
				if e.CursorX > len(e.currentLine()) {
					e.CursorX = len(e.currentLine())
				}
			}
		case "down":
			if e.CursorY < len(e.Lines)-1 {
				e.CursorY++
				if e.CursorX > len(e.currentLine()) {
					e.CursorX = len(e.currentLine())
				}
			}
		case "left":
			if e.CursorX > 0 {
				e.CursorX--
			} else if e.CursorY > 0 {
				e.CursorY--
				e.CursorX = len(e.currentLine())
			}
		case "right":
			if e.CursorX < len(e.currentLine()) {
				e.CursorX++
			} else if e.CursorY < len(e.Lines)-1 {
				e.CursorY++
				e.CursorX = 0
			}
		case "home":
			e.CursorX = 0
		case "end":
			e.CursorX = len(e.currentLine())
		case "backspace":
			e.deleteChar()
		case "enter":
			e.newLine()
		case "tab":
			for i := 0; i < e.TabSize; i++ {
				e.insertChar(' ')
			}
		default:
			if len(msg.String()) == 1 {
				e.insertChar(rune(msg.String()[0]))
			}
		}

		e.ensureVisible()

	case tea.MouseMsg:
		// Handle click to position cursor
		if msg.Action == tea.MouseActionPress {
			// Calculate line number width
			gutterWidth := 0
			if e.ShowLineNums {
				gutterWidth = 5
			}
			clickX := msg.X - gutterWidth
			clickY := msg.Y + e.ScrollY

			if clickY >= 0 && clickY < len(e.Lines) {
				e.CursorY = clickY
				e.CursorX = clickX
				if e.CursorX > len(e.currentLine()) {
					e.CursorX = len(e.currentLine())
				}
				if e.CursorX < 0 {
					e.CursorX = 0
				}
			}
		}
	}

	return e, nil
}

