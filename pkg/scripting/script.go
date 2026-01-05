// Package scripting - Script execution engine
package scripting

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/makeatui/makeatui/pkg/agent"
	"github.com/makeatui/makeatui/pkg/schema"
)

// Script represents a MakeaTUI automation script
type Script struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Version     string    `json:"version"`
	Steps       []Step    `json:"steps"`
}

// Step represents a single script step
type Step struct {
	Action     string            `json:"action"`
	Component  string            `json:"component,omitempty"`
	Properties map[string]any    `json:"properties,omitempty"`
	Condition  string            `json:"condition,omitempty"`
	Loop       *LoopConfig       `json:"loop,omitempty"`
}

// LoopConfig defines loop behavior
type LoopConfig struct {
	Times int      `json:"times,omitempty"`
	Items []string `json:"items,omitempty"`
}

// ScriptEngine executes MakeaTUI scripts
type ScriptEngine struct {
	api      *agent.API
	gum      *Gum
	vars     map[string]any
	dryRun   bool
	verbose  bool
}

// NewScriptEngine creates a new script engine
func NewScriptEngine(projectName string) *ScriptEngine {
	return &ScriptEngine{
		api:    agent.NewAPI(projectName),
		gum:    NewGum(),
		vars:   make(map[string]any),
	}
}

// SetDryRun enables dry-run mode (no actual changes)
func (e *ScriptEngine) SetDryRun(dryRun bool) {
	e.dryRun = dryRun
}

// SetVerbose enables verbose output
func (e *ScriptEngine) SetVerbose(verbose bool) {
	e.verbose = verbose
}

// SetVar sets a variable
func (e *ScriptEngine) SetVar(name string, value any) {
	e.vars[name] = value
}

// LoadScript loads a script from a file
func (e *ScriptEngine) LoadScript(path string) (*Script, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var script Script
	if err := json.Unmarshal(data, &script); err != nil {
		return nil, err
	}

	return &script, nil
}

// Execute runs a script
func (e *ScriptEngine) Execute(script *Script) error {
	if e.verbose {
		fmt.Printf("🚀 Executing script: %s\n", script.Name)
	}

	for i, step := range script.Steps {
		if e.verbose {
			fmt.Printf("  Step %d: %s\n", i+1, step.Action)
		}

		if err := e.executeStep(step); err != nil {
			return fmt.Errorf("step %d failed: %w", i+1, err)
		}
	}

	if e.verbose {
		fmt.Println("✅ Script completed successfully")
	}

	return nil
}

func (e *ScriptEngine) executeStep(step Step) error {
	if e.dryRun {
		fmt.Printf("  [DRY-RUN] Would execute: %s\n", step.Action)
		return nil
	}

	switch step.Action {
	case "add_box":
		return e.addComponent(schema.TypeBox, step)
	case "add_text":
		return e.addComponent(schema.TypeText, step)
	case "add_button":
		return e.addComponent(schema.TypeButton, step)
	case "add_list":
		return e.addComponent(schema.TypeList, step)
	case "add_progress":
		return e.addComponent(schema.TypeProgress, step)
	case "set_theme":
		if theme, ok := step.Properties["theme"].(string); ok {
			e.api.SetTheme(theme)
		}
	case "clear":
		e.api.Clear()
	case "prompt":
		return e.handlePrompt(step)
	case "export":
		return e.handleExport(step)
	default:
		return fmt.Errorf("unknown action: %s", step.Action)
	}

	return nil
}

func (e *ScriptEngine) addComponent(ctype schema.ComponentType, step Step) error {
	name := getString(step.Properties, "name", "component")
	text := getString(step.Properties, "text", "")
	x := getInt(step.Properties, "x", 0)
	y := getInt(step.Properties, "y", 0)
	width := getInt(step.Properties, "width", 20)
	height := getInt(step.Properties, "height", 5)

	switch ctype {
	case schema.TypeBox:
		e.api.AddBox(name, text, x, y, width, height)
	case schema.TypeText:
		e.api.AddText(name, text, x, y)
	case schema.TypeButton:
		e.api.AddButton(name, text, x, y)
	}

	return nil
}

func (e *ScriptEngine) handlePrompt(step Step) error {
	if !IsGumInstalled() {
		return fmt.Errorf("gum is not installed")
	}

	promptType := getString(step.Properties, "type", "input")
	prompt := getString(step.Properties, "prompt", "Enter value:")
	varName := getString(step.Properties, "var", "result")

	var result string
	var err error

	switch promptType {
	case "input":
		result, err = e.gum.Input(prompt, "")
	case "confirm":
		confirmed, err := e.gum.Confirm(prompt)
		if err == nil {
			result = fmt.Sprintf("%v", confirmed)
		}
	case "choose":
		options := getStringSlice(step.Properties, "options")
		result, err = e.gum.Choose(prompt, options)
	}

	if err != nil {
		return err
	}

	e.vars[varName] = result
	return nil
}

func (e *ScriptEngine) handleExport(step Step) error {
	format := getString(step.Properties, "format", "go")
	output := getString(step.Properties, "output", "output.go")

	var content string
	if format == "json" {
		content, _ = e.api.ExportJSON()
		if !strings.HasSuffix(output, ".json") {
			output += ".json"
		}
	} else {
		content = e.api.Export()
	}

	return os.WriteFile(output, []byte(content), 0644)
}

// GetAPI returns the underlying agent API
func (e *ScriptEngine) GetAPI() *agent.API {
	return e.api
}

// Helper functions
func getString(m map[string]any, key, def string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return def
}

func getInt(m map[string]any, key string, def int) int {
	if v, ok := m[key].(float64); ok {
		return int(v)
	}
	return def
}

func getStringSlice(m map[string]any, key string) []string {
	if v, ok := m[key].([]any); ok {
		result := make([]string, len(v))
		for i, item := range v {
			result[i] = fmt.Sprintf("%v", item)
		}
		return result
	}
	return nil
}

