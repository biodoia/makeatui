// Package display - Label and text display components
package display

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Label provides a simple text label
type Label struct {
	ID      string
	Text    string
	Width   int
	Align   lipgloss.Position
	style   lipgloss.Style
}

// NewLabel creates a label
func NewLabel(id, text string) *Label {
	return &Label{
		ID:    id,
		Text:  text,
		Align: lipgloss.Left,
		style: lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")),
	}
}

// SetWidth sets label width
func (l *Label) SetWidth(width int) *Label {
	l.Width = width
	return l
}

// SetAlign sets text alignment
func (l *Label) SetAlign(align lipgloss.Position) *Label {
	l.Align = align
	return l
}

// SetStyle sets label style
func (l *Label) SetStyle(style lipgloss.Style) *Label {
	l.style = style
	return l
}

// SetText sets label text
func (l *Label) SetText(text string) *Label {
	l.Text = text
	return l
}

// View renders the label
func (l *Label) View() string {
	style := l.style.Align(l.Align)
	if l.Width > 0 {
		style = style.Width(l.Width)
	}
	return style.Render(l.Text)
}

// Heading provides a styled heading
type Heading struct {
	ID    string
	Text  string
	Level int // 1-6
	Width int
	style HeadingStyle
}

// HeadingStyle holds heading styling
type HeadingStyle struct {
	H1 lipgloss.Style
	H2 lipgloss.Style
	H3 lipgloss.Style
	H4 lipgloss.Style
	H5 lipgloss.Style
	H6 lipgloss.Style
}

// DefaultHeadingStyle returns default styling
func DefaultHeadingStyle() HeadingStyle {
	return HeadingStyle{
		H1: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E040FB")).
			Bold(true).
			Underline(true),
		H2: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9D4EDD")).
			Bold(true),
		H3: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7B2CBF")).
			Bold(true),
		H4: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#5A189A")),
		H5: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#3C096C")),
		H6: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#1A0533")),
	}
}

// NewHeading creates a heading
func NewHeading(id, text string, level int) *Heading {
	if level < 1 {
		level = 1
	}
	if level > 6 {
		level = 6
	}
	return &Heading{
		ID:    id,
		Text:  text,
		Level: level,
		style: DefaultHeadingStyle(),
	}
}

// SetWidth sets heading width
func (h *Heading) SetWidth(width int) *Heading {
	h.Width = width
	return h
}

// View renders the heading
func (h *Heading) View() string {
	var style lipgloss.Style
	switch h.Level {
	case 1:
		style = h.style.H1
	case 2:
		style = h.style.H2
	case 3:
		style = h.style.H3
	case 4:
		style = h.style.H4
	case 5:
		style = h.style.H5
	default:
		style = h.style.H6
	}

	if h.Width > 0 {
		style = style.Width(h.Width)
	}

	return style.Render(h.Text)
}

// Divider provides a horizontal divider
type Divider struct {
	ID     string
	Width  int
	Char   rune
	Label  string
	style  lipgloss.Style
}

// NewDivider creates a divider
func NewDivider(id string, width int) *Divider {
	return &Divider{
		ID:    id,
		Width: width,
		Char:  '─',
		style: lipgloss.NewStyle().Foreground(lipgloss.Color("#5A189A")),
	}
}

// SetChar sets divider character
func (d *Divider) SetChar(ch rune) *Divider {
	d.Char = ch
	return d
}

// SetLabel sets center label
func (d *Divider) SetLabel(label string) *Divider {
	d.Label = label
	return d
}

// View renders the divider
func (d *Divider) View() string {
	if d.Label == "" {
		return d.style.Render(strings.Repeat(string(d.Char), d.Width))
	}

	labelLen := len(d.Label) + 2 // +2 for spaces
	sideLen := (d.Width - labelLen) / 2

	left := strings.Repeat(string(d.Char), sideLen)
	right := strings.Repeat(string(d.Char), d.Width-sideLen-labelLen)

	return d.style.Render(left + " " + d.Label + " " + right)
}

