// Package agent - Additional command implementations
package agent

import (
	"encoding/json"
	"fmt"

	"github.com/makeatui/makeatui/internal/codegen"
	"github.com/makeatui/makeatui/pkg/schema"
)

func (s *Session) resizeComponent(params json.RawMessage) error {
	var p ResizeComponentParams
	if err := json.Unmarshal(params, &p); err != nil {
		return err
	}

	for i := range s.Canvas.Components {
		if s.Canvas.Components[i].ID == p.ID {
			s.Canvas.Components[i].Size.Width = p.Width
			s.Canvas.Components[i].Size.Height = p.Height
			return nil
		}
	}
	return fmt.Errorf("component not found: %s", p.ID)
}

func (s *Session) styleComponent(params json.RawMessage) error {
	var p StyleComponentParams
	if err := json.Unmarshal(params, &p); err != nil {
		return err
	}

	for i := range s.Canvas.Components {
		if s.Canvas.Components[i].ID == p.ID {
			s.Canvas.Components[i].Style = p.Style
			return nil
		}
	}
	return fmt.Errorf("component not found: %s", p.ID)
}

func (s *Session) setText(params json.RawMessage) error {
	var p SetTextParams
	if err := json.Unmarshal(params, &p); err != nil {
		return err
	}

	for i := range s.Canvas.Components {
		if s.Canvas.Components[i].ID == p.ID {
			s.Canvas.Components[i].Text = p.Text
			return nil
		}
	}
	return fmt.Errorf("component not found: %s", p.ID)
}

// Undo reverts the last action
func (s *Session) Undo() bool {
	if len(s.UndoStack) == 0 {
		return false
	}

	// Save current state for redo
	s.RedoStack = append(s.RedoStack, s.Canvas)

	// Restore previous state
	s.Canvas = s.UndoStack[len(s.UndoStack)-1]
	s.UndoStack = s.UndoStack[:len(s.UndoStack)-1]
	return true
}

// Redo reapplies the last undone action
func (s *Session) Redo() bool {
	if len(s.RedoStack) == 0 {
		return false
	}

	// Save current state for undo
	s.UndoStack = append(s.UndoStack, s.Canvas)

	// Restore redo state
	s.Canvas = s.RedoStack[len(s.RedoStack)-1]
	s.RedoStack = s.RedoStack[:len(s.RedoStack)-1]
	return true
}

// Export generates Go code from the current canvas
func (s *Session) Export() string {
	gen := codegen.NewGenerator(s.Canvas)
	return gen.Generate()
}

// ExportJSON exports the canvas as JSON
func (s *Session) ExportJSON() (string, error) {
	data, err := json.MarshalIndent(s.Canvas, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// LoadFromJSON loads a canvas from JSON
func (s *Session) LoadFromJSON(data string) error {
	var canvas schema.Canvas
	if err := json.Unmarshal([]byte(data), &canvas); err != nil {
		return err
	}
	s.Canvas = canvas
	return nil
}

// GetComponent returns a component by ID
func (s *Session) GetComponent(id string) *schema.Component {
	for i := range s.Canvas.Components {
		if s.Canvas.Components[i].ID == id {
			return &s.Canvas.Components[i]
		}
	}
	return nil
}

// ListComponents returns all components
func (s *Session) ListComponents() []schema.Component {
	return s.Canvas.Components
}

// Clear removes all components from the canvas
func (s *Session) Clear() {
	s.saveState()
	s.Canvas.Components = []schema.Component{}
}

// SetTheme changes the canvas theme
func (s *Session) SetTheme(theme string) {
	s.Canvas.Theme = theme
}

// Resize changes the canvas dimensions
func (s *Session) Resize(width, height int) {
	s.Canvas.Width = width
	s.Canvas.Height = height
}

