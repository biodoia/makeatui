// Package display - Chart components (Bar, Line, Sparkline) inspired by Ratatui
package display

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// BarChart renders horizontal bar charts
type BarChart struct {
	ID        string
	Title     string
	Data      []BarData
	MaxValue  float64
	Width     int
	ShowValue bool
	ShowLabel bool
	Sorted    bool
	style     BarChartStyle
}

// BarData represents a bar in the chart
type BarData struct {
	Label string
	Value float64
	Color string
}

// BarChartStyle holds styling
type BarChartStyle struct {
	Title     lipgloss.Style
	Bar       lipgloss.Style
	Label     lipgloss.Style
	Value     lipgloss.Style
	BarChar   string
	EmptyChar string
}

// DefaultBarChartStyle returns default styling
func DefaultBarChartStyle() BarChartStyle {
	return BarChartStyle{
		Title: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E040FB")).
			Bold(true),
		Bar: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9D4EDD")),
		Label: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")),
		Value: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7B2CBF")),
		BarChar:   "█",
		EmptyChar: "░",
	}
}

// NewBarChart creates a new bar chart
func NewBarChart(id string) *BarChart {
	return &BarChart{
		ID:        id,
		Data:      []BarData{},
		Width:     40,
		ShowValue: true,
		ShowLabel: true,
		style:     DefaultBarChartStyle(),
	}
}

// AddBar adds a bar to the chart
func (b *BarChart) AddBar(label string, value float64, color string) *BarChart {
	b.Data = append(b.Data, BarData{Label: label, Value: value, Color: color})
	return b
}

// SetWidth sets the chart width
func (b *BarChart) SetWidth(width int) *BarChart {
	b.Width = width
	return b
}

// SetTitle sets the title
func (b *BarChart) SetTitle(title string) *BarChart {
	b.Title = title
	return b
}

// View renders the bar chart
func (b *BarChart) View() string {
	if len(b.Data) == 0 {
		return ""
	}

	var lines []string

	// Title
	if b.Title != "" {
		lines = append(lines, b.style.Title.Render(b.Title))
		lines = append(lines, "")
	}

	// Find max value
	maxVal := b.MaxValue
	if maxVal == 0 {
		for _, d := range b.Data {
			if d.Value > maxVal {
				maxVal = d.Value
			}
		}
	}

	// Find max label width
	maxLabelWidth := 0
	for _, d := range b.Data {
		if len(d.Label) > maxLabelWidth {
			maxLabelWidth = len(d.Label)
		}
	}

	// Render bars
	barWidth := b.Width - maxLabelWidth - 10
	if barWidth < 10 {
		barWidth = 10
	}

	for _, d := range b.Data {
		// Label
		label := b.style.Label.Width(maxLabelWidth).Render(d.Label)

		// Bar
		ratio := d.Value / maxVal
		filledWidth := int(ratio * float64(barWidth))
		emptyWidth := barWidth - filledWidth

		barStyle := b.style.Bar
		if d.Color != "" {
			barStyle = barStyle.Foreground(lipgloss.Color(d.Color))
		}

		bar := barStyle.Render(strings.Repeat(b.style.BarChar, filledWidth))
		empty := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#333333")).
			Render(strings.Repeat(b.style.EmptyChar, emptyWidth))

		// Value
		value := ""
		if b.ShowValue {
			value = b.style.Value.Render(" " + formatFloat(d.Value))
		}

		lines = append(lines, label+" "+bar+empty+value)
	}

	return strings.Join(lines, "\n")
}

// Sparkline renders a compact line chart (inspired by Ratatui)
type Sparkline struct {
	ID       string
	Data     []float64
	Width    int
	Height   int
	Min      float64
	Max      float64
	ShowMinMax bool
	style    SparklineStyle
}

// SparklineStyle holds styling
type SparklineStyle struct {
	Line lipgloss.Style
	Fill lipgloss.Style
}

// DefaultSparklineStyle returns default styling
func DefaultSparklineStyle() SparklineStyle {
	return SparklineStyle{
		Line: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E040FB")),
		Fill: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#3C096C")),
	}
}

// Braille and bar characters for sparkline
var barBlocks = []string{" ", "▁", "▂", "▃", "▄", "▅", "▆", "▇", "█"}

// NewSparkline creates a new sparkline
func NewSparkline(id string) *Sparkline {
	return &Sparkline{
		ID:     id,
		Data:   []float64{},
		Width:  40,
		Height: 1,
		style:  DefaultSparklineStyle(),
	}
}

// SetData sets the data points
func (s *Sparkline) SetData(data []float64) *Sparkline {
	s.Data = data
	return s
}

// AddPoint adds a data point
func (s *Sparkline) AddPoint(value float64) *Sparkline {
	s.Data = append(s.Data, value)
	return s
}

// SetWidth sets the width
func (s *Sparkline) SetWidth(width int) *Sparkline {
	s.Width = width
	return s
}

func formatFloat(f float64) string {
	if f == float64(int(f)) {
		return strings.TrimRight(strings.TrimRight(
			strings.Replace(string(rune(int(f)+'0')), ".", "", 1), "0"), ".")
	}
	return string(rune(int(f)+'0')) + "." + string(rune(int((f-float64(int(f)))*10)+'0'))
}

