// Package markdown provides glamourous Markdown rendering using Glamour
package markdown

import (
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/makeatui/makeatui/internal/ui/styles"
)

// Renderer wraps glamour for styled Markdown rendering
type Renderer struct {
	glamour *glamour.TermRenderer
	width   int
	style   string
}

// StylePreset defines available style presets
type StylePreset string

const (
	StyleDark     StylePreset = "dark"
	StyleLight    StylePreset = "light"
	StyleDracula  StylePreset = "dracula"
	StylePink     StylePreset = "pink"
	StyleNoTTY    StylePreset = "notty"
	StyleUltraviolet StylePreset = "ultraviolet" // Custom
)

// NewRenderer creates a new Markdown renderer
func NewRenderer(width int, preset StylePreset) (*Renderer, error) {
	var styleName string

	switch preset {
	case StyleUltraviolet:
		// Use dark as base, we'll customize
		styleName = "dark"
	case StyleDark, StyleLight, StyleDracula, StylePink, StyleNoTTY:
		styleName = string(preset)
	default:
		styleName = "dark"
	}

	r, err := glamour.NewTermRenderer(
		glamour.WithStylePath(styleName),
		glamour.WithWordWrap(width),
	)
	if err != nil {
		return nil, err
	}

	return &Renderer{
		glamour: r,
		width:   width,
		style:   styleName,
	}, nil
}

// Render renders Markdown to styled terminal output
func (r *Renderer) Render(markdown string) (string, error) {
	return r.glamour.Render(markdown)
}

// RenderWithTheme renders Markdown with the MakeaTUI theme applied
func (r *Renderer) RenderWithTheme(markdown string, theme styles.Theme) (string, error) {
	rendered, err := r.glamour.Render(markdown)
	if err != nil {
		return "", err
	}

	// Apply additional styling
	container := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(theme.Border).
		Padding(1, 2).
		Width(r.width)

	return container.Render(rendered), nil
}

// RenderInline renders a short Markdown snippet inline
func (r *Renderer) RenderInline(markdown string) string {
	rendered, err := r.glamour.Render(markdown)
	if err != nil {
		return markdown
	}
	return rendered
}

// QuickRender is a convenience function for one-off rendering
func QuickRender(markdown string, width int) (string, error) {
	r, err := NewRenderer(width, StyleDark)
	if err != nil {
		return "", err
	}
	return r.Render(markdown)
}

// RenderHelp renders help text as Markdown
func RenderHelp(title, content string, theme styles.Theme) string {
	md := "# " + title + "\n\n" + content

	r, err := NewRenderer(60, StyleDark)
	if err != nil {
		return content
	}

	rendered, err := r.RenderWithTheme(md, theme)
	if err != nil {
		return content
	}

	return rendered
}

// RenderCodeBlock renders code with syntax highlighting
func RenderCodeBlock(code, language string, width int) (string, error) {
	md := "```" + language + "\n" + code + "\n```"
	return QuickRender(md, width)
}

// RenderTable renders a simple Markdown table
func RenderTable(headers []string, rows [][]string, width int) (string, error) {
	var md string

	// Header
	md += "| " 
	for _, h := range headers {
		md += h + " | "
	}
	md += "\n|"
	for range headers {
		md += " --- |"
	}
	md += "\n"

	// Rows
	for _, row := range rows {
		md += "| "
		for _, cell := range row {
			md += cell + " | "
		}
		md += "\n"
	}

	return QuickRender(md, width)
}

// ComponentDoc generates documentation for a component
func ComponentDoc(name, description string, properties map[string]string) string {
	md := "## " + name + "\n\n"
	md += description + "\n\n"

	if len(properties) > 0 {
		md += "### Properties\n\n"
		for prop, desc := range properties {
			md += "- **" + prop + "**: " + desc + "\n"
		}
	}

	return md
}

