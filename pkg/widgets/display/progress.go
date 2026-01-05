// Package display - Progress bars and indicators
package display

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ProgressBar provides a horizontal progress bar
type ProgressBar struct {
	ID        string
	Width     int
	Progress  float64 // 0.0 to 1.0
	ShowLabel bool
	ShowPercent bool
	Animated  bool
	Label     string
	style     ProgressBarStyle
}

// ProgressBarStyle holds progress bar styling
type ProgressBarStyle struct {
	Container lipgloss.Style
	Filled    lipgloss.Style
	Empty     lipgloss.Style
	Label     lipgloss.Style
	Percent   lipgloss.Style
}

// DefaultProgressBarStyle returns default styling
func DefaultProgressBarStyle() ProgressBarStyle {
	return ProgressBarStyle{
		Container: lipgloss.NewStyle(),
		Filled: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E040FB")),
		Empty: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#3C096C")),
		Label: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")),
		Percent: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9D4EDD")),
	}
}

// NewProgressBar creates a progress bar
func NewProgressBar(id string, width int) *ProgressBar {
	return &ProgressBar{
		ID:          id,
		Width:       width,
		ShowPercent: true,
		style:       DefaultProgressBarStyle(),
	}
}

// SetProgress sets progress (0.0 to 1.0)
func (p *ProgressBar) SetProgress(progress float64) *ProgressBar {
	if progress < 0 {
		progress = 0
	}
	if progress > 1 {
		progress = 1
	}
	p.Progress = progress
	return p
}

// SetLabel sets the label
func (p *ProgressBar) SetLabel(label string) *ProgressBar {
	p.Label = label
	p.ShowLabel = true
	return p
}

// View renders the progress bar
func (p *ProgressBar) View() string {
	barWidth := p.Width - 2 // For brackets
	if p.ShowPercent {
		barWidth -= 5 // For " 100%"
	}

	filled := int(float64(barWidth) * p.Progress)
	empty := barWidth - filled

	bar := "[" +
		p.style.Filled.Render(strings.Repeat("█", filled)) +
		p.style.Empty.Render(strings.Repeat("░", empty)) +
		"]"

	if p.ShowPercent {
		percent := p.style.Percent.Render(fmt.Sprintf(" %3d%%", int(p.Progress*100)))
		bar += percent
	}

	if p.ShowLabel && p.Label != "" {
		return p.style.Label.Render(p.Label) + "\n" + bar
	}

	return bar
}

// IndeterminateBar provides an indeterminate progress bar
type IndeterminateBar struct {
	ID       string
	Width    int
	Position int
	Length   int
	Speed    time.Duration
	style    ProgressBarStyle
}

// NewIndeterminateBar creates an indeterminate progress bar
func NewIndeterminateBar(id string, width int) *IndeterminateBar {
	return &IndeterminateBar{
		ID:     id,
		Width:  width,
		Length: 5,
		Speed:  100 * time.Millisecond,
		style:  DefaultProgressBarStyle(),
	}
}

// Update handles animation
func (b *IndeterminateBar) Update(msg tea.Msg) (*IndeterminateBar, tea.Cmd) {
	switch msg.(type) {
	case IndeterminateTickMsg:
		b.Position = (b.Position + 1) % (b.Width + b.Length)
	}
	return b, nil
}

// TickCmd returns animation command
func (b *IndeterminateBar) TickCmd() tea.Cmd {
	return tea.Tick(b.Speed, func(t time.Time) tea.Msg {
		return IndeterminateTickMsg{}
	})
}

// View renders the indeterminate bar
func (b *IndeterminateBar) View() string {
	barWidth := b.Width - 2

	chars := make([]string, barWidth)
	for i := range chars {
		chars[i] = "░"
	}

	// Place moving segment
	for i := 0; i < b.Length; i++ {
		pos := b.Position - i
		if pos >= 0 && pos < barWidth {
			chars[pos] = "█"
		}
	}

	bar := "[" + b.style.Filled.Render(strings.Join(chars, "")) + "]"
	return bar
}

// IndeterminateTickMsg is sent for animation
type IndeterminateTickMsg struct{}

// DownloadProgress shows download-style progress
type DownloadProgress struct {
	ID        string
	Total     int64
	Current   int64
	Speed     float64 // bytes per second
	StartTime time.Time
	Width     int
	style     ProgressBarStyle
}

// NewDownloadProgress creates a download progress indicator
func NewDownloadProgress(id string, total int64, width int) *DownloadProgress {
	return &DownloadProgress{
		ID:        id,
		Total:     total,
		Width:     width,
		StartTime: time.Now(),
		style:     DefaultProgressBarStyle(),
	}
}

// SetCurrent sets current progress
func (d *DownloadProgress) SetCurrent(current int64) *DownloadProgress {
	d.Current = current
	elapsed := time.Since(d.StartTime).Seconds()
	if elapsed > 0 {
		d.Speed = float64(current) / elapsed
	}
	return d
}

// formatBytes formats bytes to human readable
func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// View renders download progress
func (d *DownloadProgress) View() string {
	progress := float64(d.Current) / float64(d.Total)
	barWidth := d.Width - 20

	filled := int(float64(barWidth) * progress)
	empty := barWidth - filled

	bar := "[" +
		d.style.Filled.Render(strings.Repeat("█", filled)) +
		d.style.Empty.Render(strings.Repeat("░", empty)) +
		"]"

	// Stats line
	currentStr := formatBytes(d.Current)
	totalStr := formatBytes(d.Total)
	speedStr := formatBytes(int64(d.Speed)) + "/s"

	stats := fmt.Sprintf("%s / %s  %s", currentStr, totalStr, speedStr)

	return bar + " " + d.style.Percent.Render(fmt.Sprintf("%3d%%", int(progress*100))) + "\n" + stats
}

