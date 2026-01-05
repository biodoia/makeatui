// Package mcp - MCP Client for connecting to MakeaTUI servers
package mcp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Client is an MCP client for MakeaTUI
type Client struct {
	baseURL    string
	httpClient *http.Client
	sessionID  string
}

// NewClient creates a new MCP client
func NewClient(serverURL string) *Client {
	return &Client{
		baseURL:    serverURL,
		httpClient: &http.Client{},
	}
}

// CreateSession creates a new design session
func (c *Client) CreateSession(name string) (string, error) {
	resp, err := c.post("/sessions", map[string]string{"name": name})
	if err != nil {
		return "", err
	}

	var result struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(resp, &result); err != nil {
		return "", err
	}

	c.sessionID = result.ID
	return result.ID, nil
}

// SetSession sets the current session
func (c *Client) SetSession(sessionID string) {
	c.sessionID = sessionID
}

// AddBox adds a box component
func (c *Client) AddBox(name, text string, x, y, width, height int) (string, error) {
	return c.addComponent("box", name, text, x, y, width, height)
}

// AddText adds a text component
func (c *Client) AddText(name, text string, x, y int) (string, error) {
	return c.addComponent("text", name, text, x, y, 0, 0)
}

// AddButton adds a button component
func (c *Client) AddButton(name, label string, x, y int) (string, error) {
	return c.addComponent("button", name, label, x, y, 0, 0)
}

func (c *Client) addComponent(ctype, name, text string, x, y, width, height int) (string, error) {
	resp, err := c.post("/tools/add_component", map[string]any{
		"session_id": c.sessionID,
		"type":       ctype,
		"name":       name,
		"text":       text,
		"x":          x,
		"y":          y,
		"width":      width,
		"height":     height,
	})
	if err != nil {
		return "", err
	}

	var result struct {
		ID string `json:"id"`
	}
	json.Unmarshal(resp, &result)
	return result.ID, nil
}

// Move moves a component
func (c *Client) Move(componentID string, x, y int) error {
	_, err := c.post("/tools/move_component", map[string]any{
		"session_id":   c.sessionID,
		"component_id": componentID,
		"x":            x,
		"y":            y,
	})
	return err
}

// Remove removes a component
func (c *Client) Remove(componentID string) error {
	_, err := c.post("/tools/remove_component", map[string]any{
		"session_id":   c.sessionID,
		"component_id": componentID,
	})
	return err
}

// SetText sets component text
func (c *Client) SetText(componentID, text string) error {
	_, err := c.post("/tools/set_text", map[string]any{
		"session_id":   c.sessionID,
		"component_id": componentID,
		"text":         text,
	})
	return err
}

// Generate generates a TUI from description
func (c *Client) Generate(description string) error {
	_, err := c.post("/tools/generate", map[string]any{
		"session_id":  c.sessionID,
		"description": description,
	})
	return err
}

// Export exports the design
func (c *Client) Export(format string) (string, error) {
	url := fmt.Sprintf("%s/tools/export?session_id=%s&format=%s", c.baseURL, c.sessionID, format)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	return string(body), err
}

// GetCanvas gets the current canvas
func (c *Client) GetCanvas() (string, error) {
	url := fmt.Sprintf("%s/resources/canvas?session_id=%s", c.baseURL, c.sessionID)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	return string(body), err
}

// ListTemplates lists available templates
func (c *Client) ListTemplates() ([]string, error) {
	url := fmt.Sprintf("%s/templates", c.baseURL)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var templates []struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&templates); err != nil {
		return nil, err
	}

	names := make([]string, len(templates))
	for i, t := range templates {
		names[i] = t.Name
	}
	return names, nil
}

// ApplyTemplate applies a template
func (c *Client) ApplyTemplate(templateName string) error {
	_, err := c.post("/templates/apply", map[string]any{
		"session_id":    c.sessionID,
		"template_name": templateName,
	})
	return err
}

// Health checks server health
func (c *Client) Health() bool {
	resp, err := c.httpClient.Get(c.baseURL + "/health")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

func (c *Client) post(path string, data any) ([]byte, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Post(c.baseURL+path, "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

