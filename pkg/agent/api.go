// Package agent provides a high-level API for AI agents to create TUI designs
package agent

import (
	"encoding/json"
	
	"github.com/makeatui/makeatui/pkg/schema"
)

// API provides a simplified interface for AI agents
type API struct {
	session *Session
}

// NewAPI creates a new API instance
func NewAPI(projectName string) *API {
	return &API{
		session: NewSession(projectName, 80, 24),
	}
}

// AddBox adds a box component to the canvas
func (a *API) AddBox(name, text string, x, y, width, height int) string {
	params := AddComponentParams{
		Type:   schema.TypeBox,
		Name:   name,
		X:      x,
		Y:      y,
		Width:  width,
		Height: height,
		Text:   text,
	}
	data, _ := json.Marshal(params)
	cmd := Command{Type: string(CmdAddComponent), Params: data}
	a.session.Execute(cmd)
	
	// Return the ID of the last added component
	comps := a.session.ListComponents()
	if len(comps) > 0 {
		return comps[len(comps)-1].ID
	}
	return ""
}

// AddText adds a text component
func (a *API) AddText(name, text string, x, y int) string {
	params := AddComponentParams{
		Type: schema.TypeText,
		Name: name,
		X:    x,
		Y:    y,
		Text: text,
	}
	data, _ := json.Marshal(params)
	cmd := Command{Type: string(CmdAddComponent), Params: data}
	a.session.Execute(cmd)
	
	comps := a.session.ListComponents()
	if len(comps) > 0 {
		return comps[len(comps)-1].ID
	}
	return ""
}

// AddButton adds a button component
func (a *API) AddButton(name, label string, x, y int) string {
	params := AddComponentParams{
		Type:   schema.TypeButton,
		Name:   name,
		X:      x,
		Y:      y,
		Width:  len(label) + 6,
		Height: 3,
		Text:   label,
	}
	data, _ := json.Marshal(params)
	cmd := Command{Type: string(CmdAddComponent), Params: data}
	a.session.Execute(cmd)
	
	comps := a.session.ListComponents()
	if len(comps) > 0 {
		return comps[len(comps)-1].ID
	}
	return ""
}

// AddList adds a list component
func (a *API) AddList(name string, items []string, x, y, width, height int) string {
	comp := schema.NewComponent(schema.TypeList, name)
	comp.Position = schema.Position{X: x, Y: y}
	comp.Size = schema.Size{Width: width, Height: height}
	comp.Items = items
	a.session.Canvas.Components = append(a.session.Canvas.Components, comp)
	return comp.ID
}

// AddProgress adds a progress bar
func (a *API) AddProgress(name string, value float64, x, y, width int) string {
	comp := schema.NewComponent(schema.TypeProgress, name)
	comp.Position = schema.Position{X: x, Y: y}
	comp.Size = schema.Size{Width: width, Height: 1}
	comp.Value = value
	a.session.Canvas.Components = append(a.session.Canvas.Components, comp)
	return comp.ID
}

// Move moves a component to a new position
func (a *API) Move(id string, x, y int) error {
	params := MoveComponentParams{ID: id, X: x, Y: y}
	data, _ := json.Marshal(params)
	cmd := Command{Type: string(CmdMoveComponent), Params: data}
	return a.session.Execute(cmd)
}

// Resize resizes a component
func (a *API) Resize(id string, width, height int) error {
	params := ResizeComponentParams{ID: id, Width: width, Height: height}
	data, _ := json.Marshal(params)
	cmd := Command{Type: string(CmdResizeComponent), Params: data}
	return a.session.Execute(cmd)
}

// SetText sets the text content of a component
func (a *API) SetText(id, text string) error {
	params := SetTextParams{ID: id, Text: text}
	data, _ := json.Marshal(params)
	cmd := Command{Type: string(CmdSetText), Params: data}
	return a.session.Execute(cmd)
}

// Delete removes a component
func (a *API) Delete(id string) error {
	params := struct{ ID string `json:"id"` }{ID: id}
	data, _ := json.Marshal(params)
	cmd := Command{Type: string(CmdRemoveComponent), Params: data}
	return a.session.Execute(cmd)
}

// Undo reverts the last action
func (a *API) Undo() bool {
	return a.session.Undo()
}

// Redo reapplies the last undone action
func (a *API) Redo() bool {
	return a.session.Redo()
}

// Export generates Go code
func (a *API) Export() string {
	return a.session.Export()
}

// ExportJSON exports the canvas as JSON
func (a *API) ExportJSON() (string, error) {
	return a.session.ExportJSON()
}

// Clear removes all components
func (a *API) Clear() {
	a.session.Clear()
}

// SetTheme sets the theme
func (a *API) SetTheme(theme string) {
	a.session.SetTheme(theme)
}

// GetCanvas returns the current canvas
func (a *API) GetCanvas() *schema.Canvas {
	return &a.session.Canvas
}

// ImportCanvas imports a canvas directly
func (a *API) ImportCanvas(canvas schema.Canvas) {
	a.session.Canvas = canvas
}

// ListComponents returns all components
func (a *API) ListComponents() []schema.Component {
	return a.session.ListComponents()
}
