// Package mouse - Mouse event messages and commands
package mouse

import tea "github.com/charmbracelet/bubbletea"

// ClickMsg is sent when a zone is clicked
type ClickMsg struct {
	ZoneID string
	X, Y   int
	Button tea.MouseButton
}

// HoverMsg is sent when mouse enters a zone
type HoverMsg struct {
	ZoneID string
	X, Y   int
}

// LeaveMsg is sent when mouse leaves a zone
type LeaveMsg struct {
	ZoneID string
}

// DragMsg is sent during drag operations
type DragMsg struct {
	ZoneID string
	DX, DY int
	X, Y   int
}

// ScrollMsg is sent on scroll wheel
type ScrollMsg struct {
	ZoneID    string
	Direction ScrollDirection
	Delta     int
}

// ScrollDirection indicates scroll direction
type ScrollDirection int

const (
	ScrollUp ScrollDirection = iota
	ScrollDown
	ScrollLeft
	ScrollRight
)

// DoubleClickMsg is sent on double click
type DoubleClickMsg struct {
	ZoneID string
	X, Y   int
}

// RightClickMsg is sent on right click (context menu)
type RightClickMsg struct {
	ZoneID string
	X, Y   int
}

// ScrollUpCmd creates a scroll up command
func ScrollUpCmd(zoneID string) tea.Cmd {
	return func() tea.Msg {
		return ScrollMsg{ZoneID: zoneID, Direction: ScrollUp, Delta: 1}
	}
}

// ScrollDownCmd creates a scroll down command
func ScrollDownCmd(zoneID string) tea.Cmd {
	return func() tea.Msg {
		return ScrollMsg{ZoneID: zoneID, Direction: ScrollDown, Delta: 1}
	}
}

// ClickCmd creates a click command
func ClickCmd(zoneID string, x, y int) tea.Cmd {
	return func() tea.Msg {
		return ClickMsg{ZoneID: zoneID, X: x, Y: y, Button: tea.MouseButtonLeft}
	}
}

// HoverCmd creates a hover command
func HoverCmd(zoneID string, x, y int) tea.Cmd {
	return func() tea.Msg {
		return HoverMsg{ZoneID: zoneID, X: x, Y: y}
	}
}

// DragCmd creates a drag command
func DragCmd(zoneID string, dx, dy, x, y int) tea.Cmd {
	return func() tea.Msg {
		return DragMsg{ZoneID: zoneID, DX: dx, DY: dy, X: x, Y: y}
	}
}

// Cursor types for different interactions
type CursorType int

const (
	CursorDefault CursorType = iota
	CursorPointer
	CursorText
	CursorMove
	CursorResizeNS
	CursorResizeEW
	CursorResizeNWSE
	CursorResizeNESW
	CursorNotAllowed
	CursorWait
	CursorCrosshair
)

// CursorChangeMsg requests cursor change
type CursorChangeMsg struct {
	Cursor CursorType
}

// SetCursorCmd creates a cursor change command
func SetCursorCmd(cursor CursorType) tea.Cmd {
	return func() tea.Msg {
		return CursorChangeMsg{Cursor: cursor}
	}
}

// MouseState tracks overall mouse state
type MouseState struct {
	X, Y          int
	LeftPressed   bool
	RightPressed  bool
	MiddlePressed bool
	Dragging      bool
	DragStartX    int
	DragStartY    int
}

// UpdateFromMsg updates state from a mouse message
func (ms *MouseState) UpdateFromMsg(msg tea.MouseMsg) {
	ms.X = msg.X
	ms.Y = msg.Y

	switch msg.Action {
	case tea.MouseActionPress:
		switch msg.Button {
		case tea.MouseButtonLeft:
			ms.LeftPressed = true
			ms.DragStartX = msg.X
			ms.DragStartY = msg.Y
		case tea.MouseButtonRight:
			ms.RightPressed = true
		case tea.MouseButtonMiddle:
			ms.MiddlePressed = true
		}
	case tea.MouseActionRelease:
		ms.LeftPressed = false
		ms.RightPressed = false
		ms.MiddlePressed = false
		ms.Dragging = false
	case tea.MouseActionMotion:
		if ms.LeftPressed {
			ms.Dragging = true
		}
	}
}

