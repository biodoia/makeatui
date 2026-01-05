// Package advanced - Window Manager (inspired by Lanterna)
package advanced

import (
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// WindowManager manages multiple windows
type WindowManager struct {
	ID         string
	Width      int
	Height     int
	Windows    []*Window
	FocusedIdx int
	Background string
	style      WindowManagerStyle
}

// WindowManagerStyle holds manager styling
type WindowManagerStyle struct {
	Background lipgloss.Style
	Overlay    lipgloss.Style
}

// DefaultWindowManagerStyle returns default styling
func DefaultWindowManagerStyle() WindowManagerStyle {
	return WindowManagerStyle{
		Background: lipgloss.NewStyle().
			Background(lipgloss.Color("#0D0221")),
		Overlay: lipgloss.NewStyle().
			Background(lipgloss.Color("#00000080")),
	}
}

// NewWindowManager creates a window manager
func NewWindowManager(id string, width, height int) *WindowManager {
	return &WindowManager{
		ID:         id,
		Width:      width,
		Height:     height,
		Windows:    []*Window{},
		FocusedIdx: -1,
		style:      DefaultWindowManagerStyle(),
	}
}

// AddWindow adds a window
func (wm *WindowManager) AddWindow(w *Window) *WindowManager {
	w.zIndex = len(wm.Windows)
	wm.Windows = append(wm.Windows, w)
	wm.FocusWindow(len(wm.Windows) - 1)
	return wm
}

// RemoveWindow removes a window by ID
func (wm *WindowManager) RemoveWindow(id string) *WindowManager {
	for i, w := range wm.Windows {
		if w.ID == id {
			wm.Windows = append(wm.Windows[:i], wm.Windows[i+1:]...)
			break
		}
	}
	if wm.FocusedIdx >= len(wm.Windows) {
		wm.FocusedIdx = len(wm.Windows) - 1
	}
	return wm
}

// FocusWindow focuses a window by index
func (wm *WindowManager) FocusWindow(idx int) *WindowManager {
	if idx < 0 || idx >= len(wm.Windows) {
		return wm
	}

	// Unfocus all
	for _, w := range wm.Windows {
		w.Focused = false
	}

	// Focus selected and bring to front
	wm.Windows[idx].Focused = true
	wm.Windows[idx].zIndex = len(wm.Windows)
	wm.FocusedIdx = idx

	// Sort by zIndex
	sort.Slice(wm.Windows, func(i, j int) bool {
		return wm.Windows[i].zIndex < wm.Windows[j].zIndex
	})

	// Find new focused index
	for i, w := range wm.Windows {
		if w.Focused {
			wm.FocusedIdx = i
			break
		}
	}

	return wm
}

// FocusNext focuses the next window
func (wm *WindowManager) FocusNext() *WindowManager {
	if len(wm.Windows) == 0 {
		return wm
	}
	next := (wm.FocusedIdx + 1) % len(wm.Windows)
	return wm.FocusWindow(next)
}

// FocusPrev focuses the previous window
func (wm *WindowManager) FocusPrev() *WindowManager {
	if len(wm.Windows) == 0 {
		return wm
	}
	prev := wm.FocusedIdx - 1
	if prev < 0 {
		prev = len(wm.Windows) - 1
	}
	return wm.FocusWindow(prev)
}

// GetFocused returns the focused window
func (wm *WindowManager) GetFocused() *Window {
	if wm.FocusedIdx >= 0 && wm.FocusedIdx < len(wm.Windows) {
		return wm.Windows[wm.FocusedIdx]
	}
	return nil
}

// Update handles messages
func (wm *WindowManager) Update(msg tea.Msg) (*WindowManager, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			wm.FocusNext()
		case "shift+tab":
			wm.FocusPrev()
		}

	case tea.MouseMsg:
		// Check if click is on any window (reverse order for z-index)
		for i := len(wm.Windows) - 1; i >= 0; i-- {
			w := wm.Windows[i]
			if w.Visible &&
				msg.X >= w.X && msg.X < w.X+w.Width &&
				msg.Y >= w.Y && msg.Y < w.Y+w.Height {
				if msg.Action == tea.MouseActionPress {
					wm.FocusWindow(i)
				}
				break
			}
		}

	case WindowClosedMsg:
		wm.RemoveWindow(msg.ID)
	}

	// Update focused window
	if focused := wm.GetFocused(); focused != nil {
		var cmd tea.Cmd
		_, cmd = focused.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	return wm, tea.Batch(cmds...)
}

// View renders all windows
func (wm *WindowManager) View() string {
	// Create background canvas
	canvas := make([][]rune, wm.Height)
	for y := range canvas {
		canvas[y] = make([]rune, wm.Width)
		for x := range canvas[y] {
			canvas[y][x] = ' '
		}
	}

	// Render each window onto canvas (in z-order)
	for _, w := range wm.Windows {
		if !w.Visible {
			continue
		}

		windowView := w.View()
		windowLines := strings.Split(windowView, "\n")

		for dy, line := range windowLines {
			y := w.Y + dy
			if y < 0 || y >= wm.Height {
				continue
			}

			for dx, ch := range line {
				x := w.X + dx
				if x < 0 || x >= wm.Width {
					continue
				}
				canvas[y][x] = ch
			}
		}
	}

	// Convert canvas to string
	var lines []string
	for _, row := range canvas {
		lines = append(lines, string(row))
	}

	return wm.style.Background.Render(strings.Join(lines, "\n"))
}

