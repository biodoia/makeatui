// Package agent provides an API for AI agents to interact with MakeaTUI
package agent

import (
	"encoding/json"
	"fmt"

	"github.com/makeatui/makeatui/pkg/schema"
)

// Command represents a command that can be executed by an AI agent
type Command struct {
	Type   string          `json:"type"`
	Params json.RawMessage `json:"params"`
}

// CommandType defines available command types
type CommandType string

const (
	CmdAddComponent    CommandType = "add_component"
	CmdRemoveComponent CommandType = "remove_component"
	CmdMoveComponent   CommandType = "move_component"
	CmdResizeComponent CommandType = "resize_component"
	CmdStyleComponent  CommandType = "style_component"
	CmdSetText         CommandType = "set_text"
	CmdExport          CommandType = "export"
	CmdSave            CommandType = "save"
	CmdLoad            CommandType = "load"
	CmdUndo            CommandType = "undo"
	CmdRedo            CommandType = "redo"
)

// AddComponentParams parameters for adding a component
type AddComponentParams struct {
	Type     schema.ComponentType `json:"type"`
	Name     string               `json:"name"`
	X        int                  `json:"x"`
	Y        int                  `json:"y"`
	Width    int                  `json:"width"`
	Height   int                  `json:"height"`
	Text     string               `json:"text,omitempty"`
	Style    *schema.Style        `json:"style,omitempty"`
}

// MoveComponentParams parameters for moving a component
type MoveComponentParams struct {
	ID string `json:"id"`
	X  int    `json:"x"`
	Y  int    `json:"y"`
}

// ResizeComponentParams parameters for resizing a component
type ResizeComponentParams struct {
	ID     string `json:"id"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

// StyleComponentParams parameters for styling a component
type StyleComponentParams struct {
	ID    string       `json:"id"`
	Style schema.Style `json:"style"`
}

// SetTextParams parameters for setting text content
type SetTextParams struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

// Session represents an AI agent's design session
type Session struct {
	Canvas     schema.Canvas
	History    []Command
	UndoStack  []schema.Canvas
	RedoStack  []schema.Canvas
}

// NewSession creates a new AI agent session
func NewSession(name string, width, height int) *Session {
	return &Session{
		Canvas: schema.Canvas{
			Name:       name,
			Width:      width,
			Height:     height,
			Components: []schema.Component{},
			Theme:      "ultraviolet",
		},
		History:   []Command{},
		UndoStack: []schema.Canvas{},
		RedoStack: []schema.Canvas{},
	}
}

// Execute executes a command on the session
func (s *Session) Execute(cmd Command) error {
	// Save state for undo
	s.saveState()

	switch CommandType(cmd.Type) {
	case CmdAddComponent:
		return s.addComponent(cmd.Params)
	case CmdRemoveComponent:
		return s.removeComponent(cmd.Params)
	case CmdMoveComponent:
		return s.moveComponent(cmd.Params)
	case CmdResizeComponent:
		return s.resizeComponent(cmd.Params)
	case CmdStyleComponent:
		return s.styleComponent(cmd.Params)
	case CmdSetText:
		return s.setText(cmd.Params)
	default:
		return fmt.Errorf("unknown command type: %s", cmd.Type)
	}
}

func (s *Session) saveState() {
	// Deep copy canvas for undo
	stateCopy := s.Canvas
	stateCopy.Components = make([]schema.Component, len(s.Canvas.Components))
	copy(stateCopy.Components, s.Canvas.Components)
	s.UndoStack = append(s.UndoStack, stateCopy)
	s.RedoStack = nil // Clear redo stack on new action
}

func (s *Session) addComponent(params json.RawMessage) error {
	var p AddComponentParams
	if err := json.Unmarshal(params, &p); err != nil {
		return err
	}

	comp := schema.NewComponent(p.Type, p.Name)
	comp.Position = schema.Position{X: p.X, Y: p.Y}
	if p.Width > 0 {
		comp.Size.Width = p.Width
	}
	if p.Height > 0 {
		comp.Size.Height = p.Height
	}
	comp.Text = p.Text
	if p.Style != nil {
		comp.Style = *p.Style
	}

	s.Canvas.Components = append(s.Canvas.Components, comp)
	return nil
}

func (s *Session) removeComponent(params json.RawMessage) error {
	var p struct{ ID string `json:"id"` }
	if err := json.Unmarshal(params, &p); err != nil {
		return err
	}

	for i, comp := range s.Canvas.Components {
		if comp.ID == p.ID {
			s.Canvas.Components = append(s.Canvas.Components[:i], s.Canvas.Components[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("component not found: %s", p.ID)
}

func (s *Session) moveComponent(params json.RawMessage) error {
	var p MoveComponentParams
	if err := json.Unmarshal(params, &p); err != nil {
		return err
	}

	for i := range s.Canvas.Components {
		if s.Canvas.Components[i].ID == p.ID {
			s.Canvas.Components[i].Position.X = p.X
			s.Canvas.Components[i].Position.Y = p.Y
			return nil
		}
	}
	return fmt.Errorf("component not found: %s", p.ID)
}

