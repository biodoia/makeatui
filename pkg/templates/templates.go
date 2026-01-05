// Package templates provides pre-built TUI layout templates
package templates

import (
	"github.com/makeatui/makeatui/pkg/agent"
	"github.com/makeatui/makeatui/pkg/schema"
)

// Template represents a reusable TUI layout template
type Template struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Category    TemplateCategory  `json:"category"`
	Preview     string            `json:"preview"`
	Canvas      schema.Canvas     `json:"canvas"`
}

// TemplateCategory categorizes templates
type TemplateCategory string

const (
	CategoryDashboard  TemplateCategory = "dashboard"
	CategoryForm       TemplateCategory = "form"
	CategoryList       TemplateCategory = "list"
	CategoryWizard     TemplateCategory = "wizard"
	CategoryMonitor    TemplateCategory = "monitor"
	CategoryChat       TemplateCategory = "chat"
	CategoryEditor     TemplateCategory = "editor"
	CategorySettings   TemplateCategory = "settings"
)

// TemplateEngine manages and applies templates
type TemplateEngine struct {
	templates map[string]*Template
}

// NewTemplateEngine creates a new template engine with built-in templates
func NewTemplateEngine() *TemplateEngine {
	e := &TemplateEngine{
		templates: make(map[string]*Template),
	}
	e.registerBuiltins()
	return e
}

// Register adds a template
func (e *TemplateEngine) Register(t *Template) {
	e.templates[t.Name] = t
}

// Get returns a template by name
func (e *TemplateEngine) Get(name string) *Template {
	return e.templates[name]
}

// List returns all templates
func (e *TemplateEngine) List() []*Template {
	result := make([]*Template, 0, len(e.templates))
	for _, t := range e.templates {
		result = append(result, t)
	}
	return result
}

// ListByCategory returns templates in a category
func (e *TemplateEngine) ListByCategory(category TemplateCategory) []*Template {
	var result []*Template
	for _, t := range e.templates {
		if t.Category == category {
			result = append(result, t)
		}
	}
	return result
}

// Apply creates an API with the template already applied
func (e *TemplateEngine) Apply(templateName string) (*agent.API, error) {
	t := e.Get(templateName)
	if t == nil {
		return nil, ErrTemplateNotFound
	}

	api := agent.NewAPI(t.Name)
	api.ImportCanvas(t.Canvas)
	return api, nil
}

// registerBuiltins adds all built-in templates
func (e *TemplateEngine) registerBuiltins() {
	e.Register(DashboardBasic())
	e.Register(DashboardMetrics())
	e.Register(FormSimple())
	e.Register(FormMultiStep())
	e.Register(ListBrowser())
	e.Register(ChatInterface())
	e.Register(MonitorSystem())
	e.Register(SettingsPanel())
}

// ErrTemplateNotFound indicates the template doesn't exist
var ErrTemplateNotFound = templateError("template not found")

type templateError string
func (e templateError) Error() string { return string(e) }

// BuildTemplate creates a template using the fluent API
type TemplateBuilder struct {
	template *Template
	api      *agent.API
}

// NewTemplateBuilder starts building a template
func NewTemplateBuilder(name, description string, category TemplateCategory) *TemplateBuilder {
	return &TemplateBuilder{
		template: &Template{
			Name:        name,
			Description: description,
			Category:    category,
		},
		api: agent.NewAPI(name),
	}
}

// AddBox adds a box to the template
func (b *TemplateBuilder) AddBox(name, title string, x, y, w, h int) *TemplateBuilder {
	b.api.AddBox(name, title, x, y, w, h)
	return b
}

// AddText adds text to the template
func (b *TemplateBuilder) AddText(name, content string, x, y int) *TemplateBuilder {
	b.api.AddText(name, content, x, y)
	return b
}

// AddButton adds a button to the template
func (b *TemplateBuilder) AddButton(name, label string, x, y int) *TemplateBuilder {
	b.api.AddButton(name, label, x, y)
	return b
}

// AddList adds a list to the template
func (b *TemplateBuilder) AddList(name string, items []string, x, y, w, h int) *TemplateBuilder {
	b.api.AddList(name, items, x, y, w, h)
	return b
}

// AddProgress adds a progress bar to the template
func (b *TemplateBuilder) AddProgress(name string, value float64, x, y, w int) *TemplateBuilder {
	b.api.AddProgress(name, value, x, y, w)
	return b
}

// SetPreview sets the ASCII preview
func (b *TemplateBuilder) SetPreview(preview string) *TemplateBuilder {
	b.template.Preview = preview
	return b
}

// Build finalizes the template
func (b *TemplateBuilder) Build() *Template {
	b.template.Canvas = *b.api.GetCanvas()
	return b.template
}

