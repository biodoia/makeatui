// Package canvas implements the design canvas area
package canvas

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/makeatui/makeatui/internal/ui/components"
	"github.com/makeatui/makeatui/internal/ui/styles"
	"github.com/makeatui/makeatui/pkg/schema"
)

// Canvas represents the main design area
type Canvas struct {
	Width      int
	Height     int
	Components []schema.Component
	Selected   int
	Theme      styles.Theme
	CursorX    int
	CursorY    int
	Mode       Mode
}

// Mode represents canvas interaction mode
type Mode int

const (
	ModeSelect Mode = iota
	ModeDraw
	ModeMove
	ModeResize
)

// New creates a new canvas
func New(width, height int, theme styles.Theme) *Canvas {
	return &Canvas{
		Width:      width,
		Height:     height,
		Components: []schema.Component{},
		Selected:   -1,
		Theme:      theme,
		CursorX:    width / 2,
		CursorY:    height / 2,
		Mode:       ModeSelect,
	}
}

// AddComponent adds a component to the canvas
func (c *Canvas) AddComponent(comp schema.Component) {
	comp.Position.X = c.CursorX
	comp.Position.Y = c.CursorY
	c.Components = append(c.Components, comp)
	c.Selected = len(c.Components) - 1
}

// RemoveSelected removes the selected component
func (c *Canvas) RemoveSelected() {
	if c.Selected >= 0 && c.Selected < len(c.Components) {
		c.Components = append(c.Components[:c.Selected], c.Components[c.Selected+1:]...)
		if c.Selected >= len(c.Components) {
			c.Selected = len(c.Components) - 1
		}
	}
}

// MoveSelected moves the selected component
func (c *Canvas) MoveSelected(dx, dy int) {
	if c.Selected >= 0 && c.Selected < len(c.Components) {
		c.Components[c.Selected].Position.X += dx
		c.Components[c.Selected].Position.Y += dy
	}
}

// GetSelected returns the selected component
func (c *Canvas) GetSelected() *schema.Component {
	if c.Selected >= 0 && c.Selected < len(c.Components) {
		return &c.Components[c.Selected]
	}
	return nil
}

// Render renders the canvas with all components
func (c *Canvas) Render() string {
	// Create the grid
	grid := make([][]rune, c.Height)
	for y := range grid {
		grid[y] = make([]rune, c.Width)
		for x := range grid[y] {
			grid[y][x] = ' '
		}
	}

	// Draw grid dots for visual reference
	for y := 0; y < c.Height; y += 4 {
		for x := 0; x < c.Width; x += 8 {
			if x < c.Width && y < c.Height {
				grid[y][x] = '·'
			}
		}
	}

	// Convert grid to string
	var lines []string
	for _, row := range grid {
		lines = append(lines, string(row))
	}
	background := strings.Join(lines, "\n")

	// Render components on top
	renderedComps := []string{}
	for i, comp := range c.Components {
		comp.Selected = (i == c.Selected)
		rendered := c.renderComponent(comp)
		renderedComps = append(renderedComps, rendered)
	}

	// Draw cursor
	cursorStyle := lipgloss.NewStyle().
		Foreground(c.Theme.Accent).
		Bold(true)
	cursor := cursorStyle.Render("╋")

	// Combine background and components
	result := background
	for i, comp := range c.Components {
		result = c.placeAt(result, renderedComps[i], comp.Position.X, comp.Position.Y)
	}

	// Place cursor
	result = c.placeAt(result, cursor, c.CursorX, c.CursorY)

	// Add mode indicator
	modeStyle := lipgloss.NewStyle().
		Foreground(c.Theme.TextMuted).
		Italic(true)
	modeText := fmt.Sprintf(" Mode: %s | Components: %d ", c.modeString(), len(c.Components))

	return result + "\n" + modeStyle.Render(modeText)
}

func (c *Canvas) renderComponent(comp schema.Component) string {
	switch comp.Type {
	case schema.TypeBox:
		return components.RenderBox(comp, c.Theme)
	case schema.TypeText:
		return components.RenderText(comp, c.Theme)
	case schema.TypeButton:
		return components.RenderButton(comp, c.Theme)
	case schema.TypeProgress:
		return components.RenderProgress(comp, c.Theme)
	case schema.TypeList:
		return components.RenderList(comp, c.Theme, 0)
	default:
		return components.RenderBox(comp, c.Theme)
	}
}

func (c *Canvas) placeAt(base, overlay string, x, y int) string {
	// Simple placement - in production use proper overlay logic
	return base // Simplified for now
}

func (c *Canvas) modeString() string {
	switch c.Mode {
	case ModeSelect:
		return "SELECT"
	case ModeDraw:
		return "DRAW"
	case ModeMove:
		return "MOVE"
	case ModeResize:
		return "RESIZE"
	default:
		return "UNKNOWN"
	}
}

