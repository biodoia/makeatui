// Package components - List component implementation
package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/makeatui/makeatui/internal/ui/styles"
	"github.com/makeatui/makeatui/pkg/schema"
)

// RenderList renders a list component
func RenderList(c schema.Component, theme styles.Theme, selectedIndex int) string {
	if len(c.Items) == 0 {
		return lipgloss.NewStyle().
			Foreground(theme.TextMuted).
			Italic(true).
			Render("(empty list)")
	}

	var lines []string

	normalStyle := lipgloss.NewStyle().
		Foreground(theme.TextSecondary).
		PaddingLeft(2)

	selectedStyle := lipgloss.NewStyle().
		Foreground(theme.TextPrimary).
		Background(theme.SurfaceLight).
		Bold(true).
		PaddingLeft(1)

	cursor := lipgloss.NewStyle().
		Foreground(theme.Primary).
		Bold(true).
		Render("▸")

	for i, item := range c.Items {
		if i == selectedIndex {
			lines = append(lines, cursor+" "+selectedStyle.Render(item))
		} else {
			lines = append(lines, normalStyle.Render(item))
		}
	}

	content := strings.Join(lines, "\n")

	// Apply container styling
	containerStyle := lipgloss.NewStyle().
		Width(c.Size.Width).
		Height(c.Size.Height)

	if c.Style.Border != nil {
		containerStyle = containerStyle.
			BorderStyle(GetBorderStyle(c.Style.Border.Style)).
			BorderForeground(lipgloss.Color(c.Style.Border.Color))

		if c.Style.Border.Color == "" {
			if c.Focused {
				containerStyle = containerStyle.BorderForeground(theme.Primary)
			} else {
				containerStyle = containerStyle.BorderForeground(theme.Border)
			}
		}
	}

	return containerStyle.Render(content)
}

// RenderTabs renders a tab component
func RenderTabs(c schema.Component, theme styles.Theme, activeIndex int) string {
	if len(c.Items) == 0 {
		return ""
	}

	var tabs []string

	activeStyle := lipgloss.NewStyle().
		Foreground(theme.TextPrimary).
		Background(theme.Primary).
		Padding(0, 2).
		Bold(true)

	inactiveStyle := lipgloss.NewStyle().
		Foreground(theme.TextSecondary).
		Background(theme.Surface).
		Padding(0, 2)

	for i, item := range c.Items {
		if i == activeIndex {
			tabs = append(tabs, activeStyle.Render(item))
		} else {
			tabs = append(tabs, inactiveStyle.Render(item))
		}
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, tabs...)
}

// RenderTable renders a simple table
func RenderTable(c schema.Component, theme styles.Theme) string {
	if len(c.Items) == 0 {
		return lipgloss.NewStyle().
			Foreground(theme.TextMuted).
			Italic(true).
			Render("(empty table)")
	}

	headerStyle := lipgloss.NewStyle().
		Foreground(theme.TextPrimary).
		Bold(true).
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true).
		BorderForeground(theme.Border)

	cellStyle := lipgloss.NewStyle().
		Foreground(theme.TextSecondary).
		Padding(0, 2)

	// Simple single-column table for now
	var rows []string
	for i, item := range c.Items {
		if i == 0 {
			rows = append(rows, headerStyle.Render(item))
		} else {
			rows = append(rows, cellStyle.Render(item))
		}
	}

	content := strings.Join(rows, "\n")

	containerStyle := lipgloss.NewStyle()
	if c.Style.Border != nil {
		containerStyle = containerStyle.
			BorderStyle(GetBorderStyle(c.Style.Border.Style)).
			BorderForeground(theme.Border)
	}

	return containerStyle.Render(content)
}

