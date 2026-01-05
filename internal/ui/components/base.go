// Package components provides TUI component implementations
package components

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/makeatui/makeatui/internal/ui/styles"
	"github.com/makeatui/makeatui/pkg/schema"
)

// Renderer renders a component to a string
type Renderer interface {
	Render(c schema.Component, theme styles.Theme) string
}

// GetBorderStyle converts border style string to lipgloss border
func GetBorderStyle(style string) lipgloss.Border {
	switch style {
	case "rounded":
		return lipgloss.RoundedBorder()
	case "thick":
		return lipgloss.ThickBorder()
	case "double":
		return lipgloss.DoubleBorder()
	case "none":
		return lipgloss.HiddenBorder()
	default:
		return lipgloss.NormalBorder()
	}
}

// RenderBox renders a box component
func RenderBox(c schema.Component, theme styles.Theme) string {
	style := lipgloss.NewStyle().
		Width(c.Size.Width).
		Height(c.Size.Height).
		Padding(c.Style.Padding.Top, c.Style.Padding.Right, c.Style.Padding.Bottom, c.Style.Padding.Left).
		Margin(c.Style.Margin.Top, c.Style.Margin.Right, c.Style.Margin.Bottom, c.Style.Margin.Left)

	if c.Style.Foreground != "" {
		style = style.Foreground(lipgloss.Color(c.Style.Foreground))
	} else {
		style = style.Foreground(theme.TextPrimary)
	}

	if c.Style.Background != "" {
		style = style.Background(lipgloss.Color(c.Style.Background))
	}

	if c.Style.Border != nil {
		style = style.BorderStyle(GetBorderStyle(c.Style.Border.Style))
		if c.Style.Border.Color != "" {
			style = style.BorderForeground(lipgloss.Color(c.Style.Border.Color))
		} else if c.Selected {
			style = style.BorderForeground(theme.Primary)
		} else {
			style = style.BorderForeground(theme.Border)
		}
	}

	if c.Style.Bold {
		style = style.Bold(true)
	}
	if c.Style.Italic {
		style = style.Italic(true)
	}

	return style.Render(c.Text)
}

// RenderText renders a text component
func RenderText(c schema.Component, theme styles.Theme) string {
	style := lipgloss.NewStyle()

	if c.Style.Foreground != "" {
		style = style.Foreground(lipgloss.Color(c.Style.Foreground))
	} else {
		style = style.Foreground(theme.TextPrimary)
	}

	if c.Style.Bold {
		style = style.Bold(true)
	}
	if c.Style.Italic {
		style = style.Italic(true)
	}
	if c.Style.Underline {
		style = style.Underline(true)
	}

	switch c.Style.Align {
	case "center":
		style = style.Align(lipgloss.Center)
	case "right":
		style = style.Align(lipgloss.Right)
	default:
		style = style.Align(lipgloss.Left)
	}

	if c.Size.Width > 0 {
		style = style.Width(c.Size.Width)
	}

	return style.Render(c.Text)
}

// RenderButton renders a button component
func RenderButton(c schema.Component, theme styles.Theme) string {
	var style lipgloss.Style

	if c.Focused || c.Selected {
		style = lipgloss.NewStyle().
			Background(theme.Primary).
			Foreground(theme.TextPrimary).
			Padding(0, 3).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(theme.Accent).
			Bold(true)
	} else {
		style = lipgloss.NewStyle().
			Background(theme.Surface).
			Foreground(theme.TextPrimary).
			Padding(0, 3).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(theme.Border)
	}

	if c.Disabled {
		style = style.Foreground(theme.TextMuted)
	}

	return style.Render(c.Text)
}

// RenderProgress renders a progress bar
func RenderProgress(c schema.Component, theme styles.Theme) string {
	width := c.Size.Width
	if width == 0 {
		width = 40
	}

	percent := 0.0
	if v, ok := c.Value.(float64); ok {
		percent = v
	}

	filled := int(float64(width-2) * percent)
	empty := width - 2 - filled

	bar := "█"
	empty_bar := "░"

	progress := ""
	for i := 0; i < filled; i++ {
		progress += bar
	}
	for i := 0; i < empty; i++ {
		progress += empty_bar
	}

	style := lipgloss.NewStyle().
		Foreground(theme.Primary)

	return "[" + style.Render(progress) + "]"
}

