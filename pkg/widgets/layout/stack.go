// Package layout - Stack and layer components
package layout

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Stack implements a z-index layered stack (like Textual's layers)
type Stack struct {
	layers []Layer
	width  int
	height int
}

// Layer represents a layer in the stack
type Layer struct {
	Content string
	X, Y    int
	Z       int // z-index for ordering
	Visible bool
}

// NewStack creates a new stack
func NewStack(width, height int) *Stack {
	return &Stack{
		layers: []Layer{},
		width:  width,
		height: height,
	}
}

// AddLayer adds a layer to the stack
func (s *Stack) AddLayer(content string, x, y, z int) *Stack {
	s.layers = append(s.layers, Layer{
		Content: content,
		X:       x,
		Y:       y,
		Z:       z,
		Visible: true,
	})
	return s
}

// Render composites all layers
func (s *Stack) Render() string {
	// Create base canvas
	canvas := make([][]rune, s.height)
	for i := range canvas {
		canvas[i] = []rune(strings.Repeat(" ", s.width))
	}

	// Sort layers by Z (simple bubble sort)
	sorted := make([]Layer, len(s.layers))
	copy(sorted, s.layers)
	for i := 0; i < len(sorted)-1; i++ {
		for j := 0; j < len(sorted)-i-1; j++ {
			if sorted[j].Z > sorted[j+1].Z {
				sorted[j], sorted[j+1] = sorted[j+1], sorted[j]
			}
		}
	}

	// Render layers from bottom to top
	for _, layer := range sorted {
		if !layer.Visible {
			continue
		}
		lines := strings.Split(layer.Content, "\n")
		for ly, line := range lines {
			y := layer.Y + ly
			if y < 0 || y >= s.height {
				continue
			}
			for lx, ch := range line {
				x := layer.X + lx
				if x < 0 || x >= s.width {
					continue
				}
				canvas[y][x] = ch
			}
		}
	}

	// Convert back to string
	var result strings.Builder
	for i, row := range canvas {
		result.WriteString(string(row))
		if i < len(canvas)-1 {
			result.WriteRune('\n')
		}
	}
	return result.String()
}

// VStack is a vertical stack (simpler API)
type VStack struct {
	items   []string
	gap     int
	align   lipgloss.Position
	style   lipgloss.Style
}

// NewVStack creates a vertical stack
func NewVStack() *VStack {
	return &VStack{
		items: []string{},
		gap:   0,
		align: lipgloss.Left,
	}
}

// Add adds items to the stack
func (v *VStack) Add(items ...string) *VStack {
	v.items = append(v.items, items...)
	return v
}

// SetGap sets the gap between items
func (v *VStack) SetGap(gap int) *VStack {
	v.gap = gap
	return v
}

// SetAlign sets alignment
func (v *VStack) SetAlign(align lipgloss.Position) *VStack {
	v.align = align
	return v
}

// Render renders the stack
func (v *VStack) Render() string {
	if len(v.items) == 0 {
		return ""
	}

	if v.gap > 0 {
		expanded := make([]string, 0, len(v.items)*2-1)
		for i, item := range v.items {
			expanded = append(expanded, item)
			if i < len(v.items)-1 {
				expanded = append(expanded, strings.Repeat("\n", v.gap-1))
			}
		}
		return v.style.Render(lipgloss.JoinVertical(v.align, expanded...))
	}
	return v.style.Render(lipgloss.JoinVertical(v.align, v.items...))
}

// HStack is a horizontal stack
type HStack struct {
	items []string
	gap   int
	align lipgloss.Position
	style lipgloss.Style
}

// NewHStack creates a horizontal stack
func NewHStack() *HStack {
	return &HStack{
		items: []string{},
		gap:   1,
		align: lipgloss.Top,
	}
}

// Add adds items to the stack
func (h *HStack) Add(items ...string) *HStack {
	h.items = append(h.items, items...)
	return h
}

// SetGap sets the gap between items
func (h *HStack) SetGap(gap int) *HStack {
	h.gap = gap
	return h
}

// Render renders the stack
func (h *HStack) Render() string {
	if len(h.items) == 0 {
		return ""
	}

	if h.gap > 0 {
		spacer := strings.Repeat(" ", h.gap)
		expanded := make([]string, 0, len(h.items)*2-1)
		for i, item := range h.items {
			expanded = append(expanded, item)
			if i < len(h.items)-1 {
				expanded = append(expanded, spacer)
			}
		}
		return h.style.Render(lipgloss.JoinHorizontal(h.align, expanded...))
	}
	return h.style.Render(lipgloss.JoinHorizontal(h.align, h.items...))
}

