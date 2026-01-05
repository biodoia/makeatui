// Package display - Sparkline View and additional chart components
package display

import (
	"fmt"
	"math"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// View renders the sparkline
func (s *Sparkline) View() string {
	if len(s.Data) == 0 {
		return strings.Repeat(" ", s.Width)
	}

	// Calculate min/max
	min, max := s.Min, s.Max
	if min == 0 && max == 0 {
		min, max = s.Data[0], s.Data[0]
		for _, v := range s.Data {
			if v < min {
				min = v
			}
			if v > max {
				max = v
			}
		}
	}

	// Normalize and render
	var result strings.Builder
	dataLen := len(s.Data)
	step := float64(dataLen) / float64(s.Width)
	if step < 1 {
		step = 1
	}

	for i := 0; i < s.Width && i < dataLen; i++ {
		idx := int(float64(i) * step)
		if idx >= dataLen {
			idx = dataLen - 1
		}

		value := s.Data[idx]
		normalized := 0.0
		if max != min {
			normalized = (value - min) / (max - min)
		}

		blockIdx := int(normalized * float64(len(barBlocks)-1))
		if blockIdx >= len(barBlocks) {
			blockIdx = len(barBlocks) - 1
		}
		if blockIdx < 0 {
			blockIdx = 0
		}

		result.WriteString(s.style.Line.Render(barBlocks[blockIdx]))
	}

	sparkline := result.String()

	if s.ShowMinMax {
		minMax := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")).
			Render(fmt.Sprintf(" [%.1f-%.1f]", min, max))
		return sparkline + minMax
	}

	return sparkline
}

// Gauge renders a progress gauge (inspired by Ratatui)
type Gauge struct {
	ID         string
	Label      string
	Value      float64 // 0-100
	Width      int
	ShowPercent bool
	ShowLabel   bool
	style      GaugeStyle
}

// GaugeStyle holds styling
type GaugeStyle struct {
	Label    lipgloss.Style
	Fill     lipgloss.Style
	Empty    lipgloss.Style
	Percent  lipgloss.Style
	Border   lipgloss.Style
	FillChar string
	EmptyChar string
}

// DefaultGaugeStyle returns default styling
func DefaultGaugeStyle() GaugeStyle {
	return GaugeStyle{
		Label: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E040FB")).
			Bold(true),
		Fill: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9D4EDD")),
		Empty: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#3C096C")),
		Percent: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Bold(true),
		Border: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7B2CBF")),
		FillChar:  "█",
		EmptyChar: "░",
	}
}

// NewGauge creates a new gauge
func NewGauge(id string) *Gauge {
	return &Gauge{
		ID:          id,
		Width:       30,
		ShowPercent: true,
		style:       DefaultGaugeStyle(),
	}
}

// SetValue sets the gauge value (0-100)
func (g *Gauge) SetValue(value float64) *Gauge {
	g.Value = math.Max(0, math.Min(100, value))
	return g
}

// SetLabel sets the label
func (g *Gauge) SetLabel(label string) *Gauge {
	g.Label = label
	g.ShowLabel = true
	return g
}

// SetWidth sets the gauge width
func (g *Gauge) SetWidth(width int) *Gauge {
	g.Width = width
	return g
}

// View renders the gauge
func (g *Gauge) View() string {
	var parts []string

	// Label
	if g.ShowLabel && g.Label != "" {
		parts = append(parts, g.style.Label.Render(g.Label))
	}

	// Calculate fill width
	fillWidth := int((g.Value / 100.0) * float64(g.Width))
	emptyWidth := g.Width - fillWidth

	// Render gauge bar
	fill := g.style.Fill.Render(strings.Repeat(g.style.FillChar, fillWidth))
	empty := g.style.Empty.Render(strings.Repeat(g.style.EmptyChar, emptyWidth))

	gauge := g.style.Border.Render("[") + fill + empty + g.style.Border.Render("]")

	// Percentage
	if g.ShowPercent {
		percent := g.style.Percent.Render(fmt.Sprintf(" %.0f%%", g.Value))
		gauge += percent
	}

	parts = append(parts, gauge)

	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}

// ProgressRing renders a circular progress indicator (text-based)
type ProgressRing struct {
	ID          string
	Value       float64 // 0-100
	Size        int     // diameter in characters
	ShowPercent bool
	style       lipgloss.Style
}

// NewProgressRing creates a new progress ring
func NewProgressRing(id string) *ProgressRing {
	return &ProgressRing{
		ID:          id,
		Size:        5,
		ShowPercent: true,
		style: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E040FB")),
	}
}

// SetValue sets the value
func (p *ProgressRing) SetValue(value float64) *ProgressRing {
	p.Value = math.Max(0, math.Min(100, value))
	return p
}

// View renders the progress ring
func (p *ProgressRing) View() string {
	// Simple circular representation using Unicode
	chars := []string{"○", "◔", "◑", "◕", "●"}
	idx := int((p.Value / 100.0) * float64(len(chars)-1))
	if idx >= len(chars) {
		idx = len(chars) - 1
	}

	ring := p.style.Render(chars[idx])

	if p.ShowPercent {
		return ring + " " + fmt.Sprintf("%.0f%%", p.Value)
	}
	return ring
}

