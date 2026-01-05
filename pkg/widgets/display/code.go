// Package display - Code viewer and JSON viewer components
package display

import (
	"encoding/json"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

// CodeView renders syntax-highlighted code
type CodeView struct {
	ID          string
	Code        string
	Language    string
	ShowLineNum bool
	HighlightLines []int
	StartLine   int
	Theme       string
	style       CodeStyle
}

// CodeStyle holds styling
type CodeStyle struct {
	LineNum     lipgloss.Style
	LineNumHL   lipgloss.Style
	Code        lipgloss.Style
	CodeHL      lipgloss.Style
	Border      lipgloss.Style
}

// DefaultCodeStyle returns default styling
func DefaultCodeStyle() CodeStyle {
	return CodeStyle{
		LineNum: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")).
			Width(4).
			Align(lipgloss.Right),
		LineNumHL: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E040FB")).
			Background(lipgloss.Color("#3C096C")).
			Width(4).
			Align(lipgloss.Right),
		Code: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")),
		CodeHL: lipgloss.NewStyle().
			Background(lipgloss.Color("#3C096C")).
			Foreground(lipgloss.Color("#FFFFFF")),
		Border: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7B2CBF")),
	}
}

// NewCodeView creates a code viewer
func NewCodeView(id string) *CodeView {
	return &CodeView{
		ID:          id,
		ShowLineNum: true,
		StartLine:   1,
		Theme:       "dracula",
		style:       DefaultCodeStyle(),
	}
}

// SetCode sets the code content
func (c *CodeView) SetCode(code, language string) *CodeView {
	c.Code = code
	c.Language = language
	return c
}

// SetHighlightLines sets lines to highlight
func (c *CodeView) SetHighlightLines(lines []int) *CodeView {
	c.HighlightLines = lines
	return c
}

// isHighlighted checks if a line should be highlighted
func (c *CodeView) isHighlighted(lineNum int) bool {
	for _, hl := range c.HighlightLines {
		if hl == lineNum {
			return true
		}
	}
	return false
}

// View renders the code
func (c *CodeView) View() string {
	lines := strings.Split(c.Code, "\n")
	var result []string

	for i, line := range lines {
		lineNum := c.StartLine + i
		var lineView string

		if c.ShowLineNum {
			numStyle := c.style.LineNum
			if c.isHighlighted(lineNum) {
				numStyle = c.style.LineNumHL
			}
			lineView = numStyle.Render(string(rune('0'+lineNum%10))) + " │ "
		}

		codeStyle := c.style.Code
		if c.isHighlighted(lineNum) {
			codeStyle = c.style.CodeHL
		}
		lineView += codeStyle.Render(line)

		result = append(result, lineView)
	}

	return c.style.Border.Render(strings.Join(result, "\n"))
}

// JSONView renders formatted JSON
type JSONView struct {
	ID       string
	Data     interface{}
	Indent   int
	Expanded bool
	MaxDepth int
	style    JSONStyle
}

// JSONStyle holds styling
type JSONStyle struct {
	Key      lipgloss.Style
	String   lipgloss.Style
	Number   lipgloss.Style
	Bool     lipgloss.Style
	Null     lipgloss.Style
	Bracket  lipgloss.Style
}

// DefaultJSONStyle returns default styling
func DefaultJSONStyle() JSONStyle {
	return JSONStyle{
		Key: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E040FB")),
		String: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#98C379")),
		Number: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#D19A66")),
		Bool: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#56B6C2")),
		Null: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")).
			Italic(true),
		Bracket: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ABB2BF")),
	}
}

// NewJSONView creates a JSON viewer
func NewJSONView(id string) *JSONView {
	return &JSONView{
		ID:       id,
		Indent:   2,
		Expanded: true,
		MaxDepth: 10,
		style:    DefaultJSONStyle(),
	}
}

// SetData sets the JSON data
func (j *JSONView) SetData(data interface{}) *JSONView {
	j.Data = data
	return j
}

// SetJSON parses and sets JSON string
func (j *JSONView) SetJSON(jsonStr string) error {
	var data interface{}
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		return err
	}
	j.Data = data
	return nil
}

// View renders the JSON with syntax highlighting
func (j *JSONView) View() string {
	if j.Data == nil {
		return j.style.Null.Render("null")
	}

	// Pretty print JSON
	formatted, err := json.MarshalIndent(j.Data, "", strings.Repeat(" ", j.Indent))
	if err != nil {
		return err.Error()
	}

	// Use glamour for rendering if available
	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(80),
	)
	if err == nil {
		md := "```json\n" + string(formatted) + "\n```"
		out, _ := renderer.Render(md)
		return strings.TrimSpace(out)
	}

	return string(formatted)
}

