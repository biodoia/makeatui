// Package mouse provides mouse interaction utilities for TUI components
package mouse

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Zone represents a clickable/interactive zone
type Zone struct {
	ID      string
	X, Y    int
	Width   int
	Height  int
	OnClick func() tea.Cmd
	OnHover func() tea.Cmd
	OnDrag  func(dx, dy int) tea.Cmd
	Data    interface{} // arbitrary data
}

// ZoneManager tracks interactive zones
type ZoneManager struct {
	zones       map[string]*Zone
	hovered     string
	pressed     string
	dragStart   *struct{ x, y int }
	mouseX      int
	mouseY      int
}

// NewZoneManager creates a new zone manager
func NewZoneManager() *ZoneManager {
	return &ZoneManager{
		zones: make(map[string]*Zone),
	}
}

// Register adds a zone
func (zm *ZoneManager) Register(zone *Zone) {
	zm.zones[zone.ID] = zone
}

// Unregister removes a zone
func (zm *ZoneManager) Unregister(id string) {
	delete(zm.zones, id)
}

// Clear removes all zones
func (zm *ZoneManager) Clear() {
	zm.zones = make(map[string]*Zone)
}

// Update updates zone bounds (useful after re-render)
func (zm *ZoneManager) Update(id string, x, y, width, height int) {
	if zone, ok := zm.zones[id]; ok {
		zone.X = x
		zone.Y = y
		zone.Width = width
		zone.Height = height
	}
}

// HitTest returns the zone at the given coordinates
func (zm *ZoneManager) HitTest(x, y int) *Zone {
	for _, zone := range zm.zones {
		if x >= zone.X && x < zone.X+zone.Width &&
			y >= zone.Y && y < zone.Y+zone.Height {
			return zone
		}
	}
	return nil
}

// HandleMouse processes a mouse event and returns appropriate commands
func (zm *ZoneManager) HandleMouse(msg tea.MouseMsg) tea.Cmd {
	zm.mouseX = msg.X
	zm.mouseY = msg.Y

	zone := zm.HitTest(msg.X, msg.Y)

	switch msg.Action {
	case tea.MouseActionPress:
		if zone != nil {
			zm.pressed = zone.ID
			zm.dragStart = &struct{ x, y int }{msg.X, msg.Y}
		}

	case tea.MouseActionRelease:
		if zone != nil && zm.pressed == zone.ID {
			zm.pressed = ""
			zm.dragStart = nil
			if zone.OnClick != nil {
				return zone.OnClick()
			}
		}
		zm.pressed = ""
		zm.dragStart = nil

	case tea.MouseActionMotion:
		// Handle hover
		newHovered := ""
		if zone != nil {
			newHovered = zone.ID
		}

		if newHovered != zm.hovered {
			zm.hovered = newHovered
			if zone != nil && zone.OnHover != nil {
				return zone.OnHover()
			}
		}

		// Handle drag
		if zm.pressed != "" && zm.dragStart != nil {
			if dragZone, ok := zm.zones[zm.pressed]; ok && dragZone.OnDrag != nil {
				dx := msg.X - zm.dragStart.x
				dy := msg.Y - zm.dragStart.y
				zm.dragStart.x = msg.X
				zm.dragStart.y = msg.Y
				return dragZone.OnDrag(dx, dy)
			}
		}
	}

	// Handle scroll wheel
	if zone != nil {
		switch msg.Button {
		case tea.MouseButtonWheelUp:
			return ScrollUpCmd(zone.ID)
		case tea.MouseButtonWheelDown:
			return ScrollDownCmd(zone.ID)
		}
	}

	return nil
}

// GetHovered returns the currently hovered zone ID
func (zm *ZoneManager) GetHovered() string {
	return zm.hovered
}

// GetPressed returns the currently pressed zone ID
func (zm *ZoneManager) GetPressed() string {
	return zm.pressed
}

// IsHovered checks if a zone is hovered
func (zm *ZoneManager) IsHovered(id string) bool {
	return zm.hovered == id
}

// IsPressed checks if a zone is pressed
func (zm *ZoneManager) IsPressed(id string) bool {
	return zm.pressed == id
}

// MousePosition returns current mouse position
func (zm *ZoneManager) MousePosition() (int, int) {
	return zm.mouseX, zm.mouseY
}

