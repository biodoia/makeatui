// Package feedback - Tooltip and Popover components
package feedback

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// TooltipPosition defines tooltip position
type TooltipPosition int

const (
	TooltipTop TooltipPosition = iota
	TooltipBottom
	TooltipLeft
	TooltipRight
)

// Tooltip provides a simple tooltip
type Tooltip struct {
	ID       string
	Content  string
	Position TooltipPosition
	Visible  bool
	X, Y     int
	style    TooltipStyle
}

// TooltipStyle holds tooltip styling
type TooltipStyle struct {
	Container lipgloss.Style
	Arrow     lipgloss.Style
}

// DefaultTooltipStyle returns default styling
func DefaultTooltipStyle() TooltipStyle {
	return TooltipStyle{
		Container: lipgloss.NewStyle().
			Background(lipgloss.Color("#3C096C")).
			Foreground(lipgloss.Color("#FFFFFF")).
			Padding(0, 1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7B2CBF")),
		Arrow: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7B2CBF")),
	}
}

// Arrow characters for different positions
var TooltipArrows = map[TooltipPosition]string{
	TooltipTop:    "▼",
	TooltipBottom: "▲",
	TooltipLeft:   "▶",
	TooltipRight:  "◀",
}

// NewTooltip creates a new tooltip
func NewTooltip(id, content string) *Tooltip {
	return &Tooltip{
		ID:       id,
		Content:  content,
		Position: TooltipTop,
		style:    DefaultTooltipStyle(),
	}
}

// SetPosition sets tooltip position
func (t *Tooltip) SetPosition(pos TooltipPosition) *Tooltip {
	t.Position = pos
	return t
}

// SetContent sets tooltip content
func (t *Tooltip) SetContent(content string) *Tooltip {
	t.Content = content
	return t
}

// Show shows the tooltip at coordinates
func (t *Tooltip) Show(x, y int) *Tooltip {
	t.Visible = true
	t.X = x
	t.Y = y
	return t
}

// Hide hides the tooltip
func (t *Tooltip) Hide() *Tooltip {
	t.Visible = false
	return t
}

// View renders the tooltip
func (t *Tooltip) View() string {
	if !t.Visible {
		return ""
	}

	content := t.style.Container.Render(t.Content)
	arrow := t.style.Arrow.Render(TooltipArrows[t.Position])

	switch t.Position {
	case TooltipTop:
		return lipgloss.JoinVertical(lipgloss.Center, content, arrow)
	case TooltipBottom:
		return lipgloss.JoinVertical(lipgloss.Center, arrow, content)
	case TooltipLeft:
		return lipgloss.JoinHorizontal(lipgloss.Center, content, arrow)
	case TooltipRight:
		return lipgloss.JoinHorizontal(lipgloss.Center, arrow, content)
	}

	return content
}

// Popover provides a rich popover/dropdown
type Popover struct {
	ID        string
	Title     string
	Content   string
	Visible   bool
	Width     int
	X, Y      int
	Closable  bool
	style     PopoverStyle
}

// PopoverStyle holds popover styling
type PopoverStyle struct {
	Container lipgloss.Style
	Title     lipgloss.Style
	Content   lipgloss.Style
	Close     lipgloss.Style
}

// DefaultPopoverStyle returns default styling
func DefaultPopoverStyle() PopoverStyle {
	return PopoverStyle{
		Container: lipgloss.NewStyle().
			Background(lipgloss.Color("#1A0533")).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#9D4EDD")).
			Padding(1),
		Title: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E040FB")).
			Bold(true).
			MarginBottom(1),
		Content: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")),
		Close: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")),
	}
}

// NewPopover creates a new popover
func NewPopover(id string) *Popover {
	return &Popover{
		ID:       id,
		Width:    40,
		Closable: true,
		style:    DefaultPopoverStyle(),
	}
}

// SetTitle sets popover title
func (p *Popover) SetTitle(title string) *Popover {
	p.Title = title
	return p
}

// SetContent sets popover content
func (p *Popover) SetContent(content string) *Popover {
	p.Content = content
	return p
}

// SetWidth sets popover width
func (p *Popover) SetWidth(width int) *Popover {
	p.Width = width
	return p
}

// Show shows the popover
func (p *Popover) Show(x, y int) *Popover {
	p.Visible = true
	p.X = x
	p.Y = y
	return p
}

// Hide hides the popover
func (p *Popover) Hide() *Popover {
	p.Visible = false
	return p
}

// Toggle toggles visibility
func (p *Popover) Toggle() *Popover {
	p.Visible = !p.Visible
	return p
}

// View renders the popover
func (p *Popover) View() string {
	if !p.Visible {
		return ""
	}

	var parts []string

	// Title row
	if p.Title != "" || p.Closable {
		titleRow := p.style.Title.Width(p.Width - 4).Render(p.Title)
		if p.Closable {
			closeBtn := p.style.Close.Render("[×]")
			titleRow = lipgloss.JoinHorizontal(lipgloss.Top, titleRow, closeBtn)
		}
		parts = append(parts, titleRow)
	}

	// Content
	content := p.style.Content.Width(p.Width - 4).Render(p.Content)
	parts = append(parts, content)

	inner := strings.Join(parts, "\n")
	return p.style.Container.Width(p.Width).Render(inner)
}

// Badge provides a small label/badge
type Badge struct {
	Label string
	Color string
	style lipgloss.Style
}

// NewBadge creates a badge
func NewBadge(label, color string) *Badge {
	return &Badge{
		Label: label,
		Color: color,
		style: lipgloss.NewStyle().
			Padding(0, 1).
			Bold(true),
	}
}

// View renders the badge
func (b *Badge) View() string {
	return b.style.
		Background(lipgloss.Color(b.Color)).
		Foreground(lipgloss.Color("#FFFFFF")).
		Render(b.Label)
}

// PredefinedBadges
func BadgeSuccess(label string) *Badge { return NewBadge(label, "#6BCB77") }
func BadgeWarning(label string) *Badge { return NewBadge(label, "#FFD93D") }
func BadgeError(label string) *Badge   { return NewBadge(label, "#FF6B6B") }
func BadgeInfo(label string) *Badge    { return NewBadge(label, "#4CC9F0") }
func BadgePrimary(label string) *Badge { return NewBadge(label, "#9D4EDD") }

