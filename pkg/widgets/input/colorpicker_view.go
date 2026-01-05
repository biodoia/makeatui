// Package input - ColorPicker View method
package input

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// View renders the color picker
func (c *ColorPicker) View() string {
	var parts []string

	// Label
	if c.Label != "" {
		parts = append(parts, c.style.Label.Render(c.Label))
	}

	// Preview and trigger
	previewStyle := c.style.Preview.
		Background(lipgloss.Color(c.Color))
	preview := previewStyle.Render("  ")

	triggerText := fmt.Sprintf("%s ▼", c.Color)
	if !c.ShowHex {
		triggerText = "▼"
	}

	trigger := lipgloss.JoinHorizontal(
		lipgloss.Center,
		preview,
		" ",
		triggerText,
	)

	parts = append(parts, trigger)

	// Palette
	if c.Open {
		var rows []string
		var row []string

		for i, color := range c.Palette {
			colorStyle := lipgloss.NewStyle().
				Background(lipgloss.Color(color)).
				Width(3).
				Height(1)

			if i == c.Selected {
				colorStyle = colorStyle.
					Border(lipgloss.NormalBorder()).
					BorderForeground(lipgloss.Color("#FFFFFF"))
			}

			swatch := colorStyle.Render("   ")
			row = append(row, swatch)

			if (i+1)%c.Cols == 0 || i == len(c.Palette)-1 {
				rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top, row...))
				row = []string{}
			}
		}

		palette := lipgloss.JoinVertical(lipgloss.Left, rows...)
		parts = append(parts, c.style.Palette.Render(palette))
	}

	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}

// Toggle is a switch component (inspired by Textual)
type Toggle struct {
	ID       string
	Label    string
	On       bool
	Disabled bool
	OnLabel  string
	OffLabel string
	style    ToggleStyle
}

// ToggleStyle holds styling
type ToggleStyle struct {
	Label   lipgloss.Style
	Track   lipgloss.Style
	TrackOn lipgloss.Style
	Thumb   lipgloss.Style
}

// DefaultToggleStyle returns default styling
func DefaultToggleStyle() ToggleStyle {
	return ToggleStyle{
		Label: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")),
		Track: lipgloss.NewStyle().
			Background(lipgloss.Color("#3C096C")).
			Foreground(lipgloss.Color("#666666")),
		TrackOn: lipgloss.NewStyle().
			Background(lipgloss.Color("#9D4EDD")).
			Foreground(lipgloss.Color("#FFFFFF")),
		Thumb: lipgloss.NewStyle().
			Background(lipgloss.Color("#FFFFFF")).
			Foreground(lipgloss.Color("#000000")),
	}
}

// NewToggle creates a new toggle
func NewToggle(id, label string) *Toggle {
	return &Toggle{
		ID:       id,
		Label:    label,
		OnLabel:  "ON",
		OffLabel: "OFF",
		style:    DefaultToggleStyle(),
	}
}

// SetOn sets the toggle state
func (t *Toggle) SetOn(on bool) *Toggle {
	t.On = on
	return t
}

// SetLabels sets on/off labels
func (t *Toggle) SetLabels(on, off string) *Toggle {
	t.OnLabel = on
	t.OffLabel = off
	return t
}

// Toggle toggles the switch
func (t *Toggle) Toggle() *Toggle {
	if !t.Disabled {
		t.On = !t.On
	}
	return t
}

// View renders the toggle
func (t *Toggle) View() string {
	var track string

	if t.On {
		track = t.style.TrackOn.Render(" " + t.OnLabel + " ") +
			t.style.Thumb.Render("●")
	} else {
		track = t.style.Thumb.Render("●") +
			t.style.Track.Render(" " + t.OffLabel + " ")
	}

	if t.Disabled {
		track = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")).
			Render(track)
	}

	if t.Label != "" {
		return lipgloss.JoinHorizontal(
			lipgloss.Center,
			t.style.Label.Render(t.Label+" "),
			track,
		)
	}

	return track
}

