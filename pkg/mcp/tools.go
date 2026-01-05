// Package mcp - MCP Tool definitions for AI integration
package mcp

// MCPToolSchema defines an MCP tool schema
type MCPToolSchema struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"inputSchema"`
}

// GetToolSchemas returns all MakeaTUI MCP tool schemas
func GetToolSchemas() []MCPToolSchema {
	return []MCPToolSchema{
		{
			Name:        "makeatui_create_session",
			Description: "Create a new TUI design session",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"name": map[string]string{
						"type":        "string",
						"description": "Name for the design project",
					},
				},
				"required": []string{"name"},
			},
		},
		{
			Name:        "makeatui_add_box",
			Description: "Add a box component to the TUI design",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"name":   map[string]string{"type": "string", "description": "Component name"},
					"text":   map[string]string{"type": "string", "description": "Box title/content"},
					"x":      map[string]string{"type": "integer", "description": "X position"},
					"y":      map[string]string{"type": "integer", "description": "Y position"},
					"width":  map[string]string{"type": "integer", "description": "Width"},
					"height": map[string]string{"type": "integer", "description": "Height"},
				},
				"required": []string{"name", "x", "y", "width", "height"},
			},
		},
		{
			Name:        "makeatui_add_text",
			Description: "Add a text component to the TUI design",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"name": map[string]string{"type": "string", "description": "Component name"},
					"text": map[string]string{"type": "string", "description": "Text content"},
					"x":    map[string]string{"type": "integer", "description": "X position"},
					"y":    map[string]string{"type": "integer", "description": "Y position"},
				},
				"required": []string{"name", "text", "x", "y"},
			},
		},
		{
			Name:        "makeatui_add_button",
			Description: "Add a button component to the TUI design",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"name":  map[string]string{"type": "string", "description": "Component name"},
					"label": map[string]string{"type": "string", "description": "Button label"},
					"x":     map[string]string{"type": "integer", "description": "X position"},
					"y":     map[string]string{"type": "integer", "description": "Y position"},
				},
				"required": []string{"name", "label", "x", "y"},
			},
		},
		{
			Name:        "makeatui_move_component",
			Description: "Move a component to a new position",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"component_id": map[string]string{"type": "string", "description": "Component ID"},
					"x":            map[string]string{"type": "integer", "description": "New X position"},
					"y":            map[string]string{"type": "integer", "description": "New Y position"},
				},
				"required": []string{"component_id", "x", "y"},
			},
		},
		{
			Name:        "makeatui_remove_component",
			Description: "Remove a component from the design",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"component_id": map[string]string{"type": "string", "description": "Component ID to remove"},
				},
				"required": []string{"component_id"},
			},
		},
		{
			Name:        "makeatui_generate",
			Description: "Generate a TUI layout from a natural language description",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"description": map[string]string{
						"type":        "string",
						"description": "Natural language description of the desired TUI (e.g., 'a dashboard with metrics and charts')",
					},
				},
				"required": []string{"description"},
			},
		},
		{
			Name:        "makeatui_export",
			Description: "Export the TUI design as Go code or JSON",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"format": map[string]string{
						"type":        "string",
						"description": "Export format: 'go' for Go code, 'json' for JSON",
						"enum":        "go,json",
					},
				},
			},
		},
		{
			Name:        "makeatui_apply_template",
			Description: "Apply a pre-built template to the design",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"template": map[string]string{
						"type":        "string",
						"description": "Template name (dashboard-basic, form-simple, chat-interface, etc.)",
					},
				},
				"required": []string{"template"},
			},
		},
		{
			Name:        "makeatui_get_canvas",
			Description: "Get the current canvas state as JSON",
			InputSchema: map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
		},
	}
}

// ResourceSchema defines an MCP resource
type ResourceSchema struct {
	URI         string `json:"uri"`
	Name        string `json:"name"`
	Description string `json:"description"`
	MimeType    string `json:"mimeType"`
}

// GetResourceSchemas returns MakeaTUI MCP resource schemas
func GetResourceSchemas(sessionID string) []ResourceSchema {
	return []ResourceSchema{
		{
			URI:         "makeatui://session/" + sessionID + "/canvas",
			Name:        "Current Canvas",
			Description: "The current TUI design canvas with all components",
			MimeType:    "application/json",
		},
		{
			URI:         "makeatui://session/" + sessionID + "/components",
			Name:        "Component List",
			Description: "List of all components in the current design",
			MimeType:    "application/json",
		},
		{
			URI:         "makeatui://session/" + sessionID + "/code",
			Name:        "Generated Code",
			Description: "Generated Go code for the current design",
			MimeType:    "text/x-go",
		},
	}
}

