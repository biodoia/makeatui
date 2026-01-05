// Package advanced - Window management (inspired by Lanterna)
package advanced

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/makeatui/makeatui/pkg/widgets/mouse"
)

// Window represents a draggable, resizable window
type Window struct {
	ID         string
	Title      string
	X, Y       int
	Width      int
	Height     int
	MinWidth   int
	MinHeight  int
	Content    string
	Focused    bool
	Closable   bool
	Resizable  bool
	Draggable  bool
	Visible    bool
	dragging   bool
	resizing   bool
	dragStartX int
	dragStartY int
	zIndex     int
	style      WindowStyle
}

// WindowStyle holds window styling
type WindowStyle struct {
	Border         lipgloss.Style
	BorderFocused  lipgloss.Style
	Title          lipgloss.Style
	TitleFocused   lipgloss.Style
	Content        lipgloss.Style
	CloseButton    lipgloss.Style
	ResizeHandle   lipgloss.Style
}

// DefaultWindowStyle returns default styling
func DefaultWindowStyle() WindowStyle {
	return WindowStyle{
		Border: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#666666")),
		BorderFocused: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#9D4EDD")),
		Title: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888")).
			Bold(true),
		TitleFocused: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E040FB")).
			Bold(true),
		Content: lipgloss.NewStyle().
			Padding(0, 1),
		CloseButton: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF6B6B")).
			Bold(true),
		ResizeHandle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")),
	}
}

// NewWindow creates a new window
func NewWindow(id, title string, x, y, width, height int) *Window {
	return &Window{
		ID:        id,
		Title:     title,
		X:         x,
		Y:         y,
		Width:     width,
		Height:    height,
		MinWidth:  10,
		MinHeight: 5,
		Visible:   true,
		Closable:  true,
		Resizable: true,
		Draggable: true,
		style:     DefaultWindowStyle(),
	}
}

// SetContent sets window content
func (w *Window) SetContent(content string) *Window {
	w.Content = content
	return w
}

// SetFocused sets focus state
func (w *Window) SetFocused(focused bool) *Window {
	w.Focused = focused
	return w
}

// SetPosition sets window position
func (w *Window) SetPosition(x, y int) *Window {
	w.X = x
	w.Y = y
	return w
}

// SetSize sets window size
func (w *Window) SetSize(width, height int) *Window {
	if width >= w.MinWidth {
		w.Width = width
	}
	if height >= w.MinHeight {
		w.Height = height
	}
	return w
}

// Close closes the window
func (w *Window) Close() *Window {
	w.Visible = false
	return w
}

// GetZone returns the mouse zone for the window
func (w *Window) GetZone(offsetX, offsetY int) *mouse.Zone {
	return &mouse.Zone{
		ID:     w.ID,
		X:      w.X + offsetX,
		Y:      w.Y + offsetY,
		Width:  w.Width,
		Height: w.Height,
	}
}

// Update handles messages
func (w *Window) Update(msg tea.Msg) (*Window, tea.Cmd) {
	if !w.Visible {
		return w, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if w.Focused {
			switch msg.String() {
			case "q", "esc":
				if w.Closable {
					w.Close()
					return w, func() tea.Msg {
						return WindowClosedMsg{ID: w.ID}
					}
				}
			}
		}

	case tea.MouseMsg:
		if w.Draggable && msg.Action == tea.MouseActionPress {
			// Check if click is on title bar
			if msg.Y == w.Y && msg.X >= w.X && msg.X < w.X+w.Width {
				w.dragging = true
				w.dragStartX = msg.X - w.X
				w.dragStartY = msg.Y - w.Y
			}
		}

		if w.dragging && msg.Action == tea.MouseActionMotion {
			w.X = msg.X - w.dragStartX
			w.Y = msg.Y - w.dragStartY
		}

		if msg.Action == tea.MouseActionRelease {
			w.dragging = false
			w.resizing = false
		}
	}

	return w, nil
}

// View renders the window
func (w *Window) View() string {
	if !w.Visible {
		return ""
	}

	// Choose styles based on focus
	borderStyle := w.style.Border
	titleStyle := w.style.Title
	if w.Focused {
		borderStyle = w.style.BorderFocused
		titleStyle = w.style.TitleFocused
	}

	// Build title bar
	closeBtn := ""
	if w.Closable {
		closeBtn = w.style.CloseButton.Render(" ✕")
	}

	titleWidth := w.Width - 4 - len(closeBtn)
	title := w.Title
	if len(title) > titleWidth {
		title = title[:titleWidth-3] + "..."
	}

	titleBar := titleStyle.Render(title) + strings.Repeat(" ", titleWidth-len(title)) + closeBtn

	// Build content
	contentLines := strings.Split(w.Content, "\n")
	contentHeight := w.Height - 3 // -3 for borders and title

	// Pad or truncate content
	var paddedContent []string
	for i := 0; i < contentHeight; i++ {
		if i < len(contentLines) {
			line := contentLines[i]
			if len(line) > w.Width-4 {
				line = line[:w.Width-4]
			}
			paddedContent = append(paddedContent, line)
		} else {
			paddedContent = append(paddedContent, "")
		}
	}

	content := w.style.Content.Render(strings.Join(paddedContent, "\n"))

	// Combine title and content
	inner := lipgloss.JoinVertical(lipgloss.Left, titleBar, content)

	// Add resize handle if resizable
	if w.Resizable {
		// Handle at bottom-right corner
	}

	return borderStyle.Width(w.Width).Render(inner)
}

// WindowClosedMsg is sent when window is closed
type WindowClosedMsg struct {
	ID string
}

// WindowFocusedMsg is sent when window gains focus
type WindowFocusedMsg struct {
	ID string
}

