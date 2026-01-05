// Package app - Export functionality
package app

import (
	"os"

	"github.com/makeatui/makeatui/internal/codegen"
	"github.com/makeatui/makeatui/pkg/schema"
)

// Export exports the current canvas to Go code
func (m *Model) Export(filename string) error {
	canvas := schema.Canvas{
		Name:       m.projectName,
		Width:      m.canvas.Width,
		Height:     m.canvas.Height,
		Components: m.canvas.Components,
		Theme:      m.theme.Name,
	}

	gen := codegen.NewGenerator(canvas)
	code := gen.Generate()

	return os.WriteFile(filename, []byte(code), 0644)
}

// ExportString returns the generated code as a string
func (m *Model) ExportString() string {
	canvas := schema.Canvas{
		Name:       m.projectName,
		Width:      m.canvas.Width,
		Height:     m.canvas.Height,
		Components: m.canvas.Components,
		Theme:      m.theme.Name,
	}

	gen := codegen.NewGenerator(canvas)
	return gen.Generate()
}

// GetCanvasSchema returns the canvas schema for JSON export
func (m *Model) GetCanvasSchema() schema.Canvas {
	return schema.Canvas{
		Name:       m.projectName,
		Width:      m.canvas.Width,
		Height:     m.canvas.Height,
		Components: m.canvas.Components,
		Theme:      m.theme.Name,
	}
}

