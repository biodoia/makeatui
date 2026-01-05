// Package schema defines the data structures for TUI components
package schema

import (
	"crypto/rand"
	"encoding/hex"
	"sync/atomic"
	"time"
)

// ComponentType represents the type of UI component
type ComponentType string

const (
	TypeBox      ComponentType = "box"
	TypeText     ComponentType = "text"
	TypeButton   ComponentType = "button"
	TypeInput    ComponentType = "input"
	TypeList     ComponentType = "list"
	TypeTable    ComponentType = "table"
	TypeProgress ComponentType = "progress"
	TypeSpinner  ComponentType = "spinner"
	TypeViewport ComponentType = "viewport"
	TypeTabs     ComponentType = "tabs"
)

// Position represents a position on the canvas
type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// Size represents dimensions
type Size struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

// Border represents border configuration
type Border struct {
	Style  string `json:"style"` // none, normal, rounded, thick, double
	Color  string `json:"color"`
	Left   bool   `json:"left"`
	Right  bool   `json:"right"`
	Top    bool   `json:"top"`
	Bottom bool   `json:"bottom"`
}

// Padding represents padding values
type Padding struct {
	Top    int `json:"top"`
	Right  int `json:"right"`
	Bottom int `json:"bottom"`
	Left   int `json:"left"`
}

// Margin represents margin values
type Margin struct {
	Top    int `json:"top"`
	Right  int `json:"right"`
	Bottom int `json:"bottom"`
	Left   int `json:"left"`
}

// Style represents styling for a component
type Style struct {
	Foreground string  `json:"foreground,omitempty"`
	Background string  `json:"background,omitempty"`
	Bold       bool    `json:"bold,omitempty"`
	Italic     bool    `json:"italic,omitempty"`
	Underline  bool    `json:"underline,omitempty"`
	Border     *Border `json:"border,omitempty"`
	Padding    Padding `json:"padding"`
	Margin     Margin  `json:"margin"`
	Align      string  `json:"align,omitempty"` // left, center, right
}

// Component represents a TUI component on the canvas
type Component struct {
	ID       string        `json:"id"`
	Type     ComponentType `json:"type"`
	Name     string        `json:"name"`
	Position Position      `json:"position"`
	Size     Size          `json:"size"`
	Style    Style         `json:"style"`

	// Content for different component types
	Text        string   `json:"text,omitempty"`
	Placeholder string   `json:"placeholder,omitempty"`
	Items       []string `json:"items,omitempty"`
	Value       any      `json:"value,omitempty"`

	// Children for container components
	Children []Component `json:"children,omitempty"`

	// State
	Focused  bool `json:"focused,omitempty"`
	Selected bool `json:"selected,omitempty"`
	Disabled bool `json:"disabled,omitempty"`
}

// Canvas represents the entire design canvas
type Canvas struct {
	Name       string      `json:"name"`
	Width      int         `json:"width"`
	Height     int         `json:"height"`
	Components []Component `json:"components"`
	Theme      string      `json:"theme"`
}

// NewComponent creates a new component with default values
func NewComponent(ctype ComponentType, name string) Component {
	return Component{
		ID:   generateID(),
		Type: ctype,
		Name: name,
		Position: Position{X: 0, Y: 0},
		Size:     Size{Width: 20, Height: 3},
		Style: Style{
			Border: &Border{
				Style:  "rounded",
				Left:   true,
				Right:  true,
				Top:    true,
				Bottom: true,
			},
		},
	}
}

var idCounter uint64

// generateID generates a unique ID for components
func generateID() string {
	// Combine timestamp, counter, and random bytes for uniqueness
	counter := atomic.AddUint64(&idCounter, 1)
	timestamp := time.Now().UnixNano()

	randomBytes := make([]byte, 4)
	_, _ = rand.Read(randomBytes)

	return "comp_" + hex.EncodeToString([]byte{
		byte(timestamp >> 24),
		byte(timestamp >> 16),
		byte(counter >> 8),
		byte(counter),
	}) + hex.EncodeToString(randomBytes)
}
