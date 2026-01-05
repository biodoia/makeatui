// Package effects - CRT/Retro terminal effects (inspired by cool-retro-term)
package effects

import (
	"math"
	"math/rand"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// CRTEffect applies CRT monitor effects to content
type CRTEffect struct {
	ID         string
	Width      int
	Height     int
	Content    string
	Scanlines  bool
	Flicker    bool
	Glow       bool
	Curvature  bool
	Noise      float64
	Bloom      float64
	Brightness float64
	flickerVal float64
	time       float64
	Speed      time.Duration
	style      CRTStyle
}

// CRTStyle holds CRT effect styling
type CRTStyle struct {
	Base       lipgloss.Style
	Scanline   lipgloss.Style
	GlowColor  string
	NoiseChars []rune
}

// DefaultCRTStyle returns default CRT styling
func DefaultCRTStyle() CRTStyle {
	return CRTStyle{
		Base: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00")),
		Scanline: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#004400")),
		GlowColor:  "#00FF00",
		NoiseChars: []rune{' ', '.', ':', '░'},
	}
}

// AmberCRTStyle returns amber phosphor style
func AmberCRTStyle() CRTStyle {
	return CRTStyle{
		Base: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFB000")),
		Scanline: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#664400")),
		GlowColor:  "#FFB000",
		NoiseChars: []rune{' ', '.', ':', '░'},
	}
}

// UltravioletCRTStyle returns MakeaTUI themed CRT
func UltravioletCRTStyle() CRTStyle {
	return CRTStyle{
		Base: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E040FB")),
		Scanline: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#5A189A")),
		GlowColor:  "#E040FB",
		NoiseChars: []rune{' ', '.', ':', '░'},
	}
}

// NewCRTEffect creates a CRT effect wrapper
func NewCRTEffect(width, height int) *CRTEffect {
	return &CRTEffect{
		ID:         "crt",
		Width:      width,
		Height:     height,
		Scanlines:  true,
		Flicker:    true,
		Glow:       true,
		Curvature:  false,
		Noise:      0.02,
		Bloom:      0.3,
		Brightness: 1.0,
		flickerVal: 1.0,
		Speed:      50 * time.Millisecond,
		style:      DefaultCRTStyle(),
	}
}

// SetContent sets the content to display
func (c *CRTEffect) SetContent(content string) *CRTEffect {
	c.Content = content
	return c
}

// SetScanlines enables/disables scanlines
func (c *CRTEffect) SetScanlines(enable bool) *CRTEffect {
	c.Scanlines = enable
	return c
}

// SetFlicker enables/disables flicker
func (c *CRTEffect) SetFlicker(enable bool) *CRTEffect {
	c.Flicker = enable
	return c
}

// SetNoise sets noise level (0.0-1.0)
func (c *CRTEffect) SetNoise(level float64) *CRTEffect {
	c.Noise = level
	return c
}

// SetStyle sets the CRT style
func (c *CRTEffect) SetStyle(style CRTStyle) *CRTEffect {
	c.style = style
	return c
}

// Update handles animation
func (c *CRTEffect) Update(msg tea.Msg) (*CRTEffect, tea.Cmd) {
	switch msg.(type) {
	case CRTTickMsg:
		c.time += 0.1
		if c.Flicker {
			c.flickerVal = 0.95 + rand.Float64()*0.1
		}
	}
	return c, nil
}

// TickCmd returns animation command
func (c *CRTEffect) TickCmd() tea.Cmd {
	return tea.Tick(c.Speed, func(t time.Time) tea.Msg {
		return CRTTickMsg{}
	})
}

// View renders the CRT effect
func (c *CRTEffect) View() string {
	lines := strings.Split(c.Content, "\n")

	// Pad lines to height
	for len(lines) < c.Height {
		lines = append(lines, "")
	}

	var output []string
	for y, line := range lines {
		if y >= c.Height {
			break
		}

		// Pad line to width
		for len(line) < c.Width {
			line += " "
		}
		if len(line) > c.Width {
			line = line[:c.Width]
		}

		// Apply effects
		processedLine := c.processLine(line, y)
		output = append(output, processedLine)
	}

	return strings.Join(output, "\n")
}

// processLine applies CRT effects to a single line
func (c *CRTEffect) processLine(line string, y int) string {
	var result strings.Builder

	for x, ch := range line {
		// Apply brightness and flicker
		brightness := c.Brightness * c.flickerVal

		// Add noise
		if c.Noise > 0 && rand.Float64() < c.Noise {
			noiseChar := c.style.NoiseChars[rand.Intn(len(c.style.NoiseChars))]
			ch = noiseChar
		}

		// Scanline effect (dim every other line)
		style := c.style.Base
		if c.Scanlines && y%2 == 1 {
			style = c.style.Scanline
		}

		// Apply curvature darkening at edges
		if c.Curvature {
			distX := float64(x-c.Width/2) / float64(c.Width/2)
			distY := float64(y-c.Height/2) / float64(c.Height/2)
			dist := math.Sqrt(distX*distX + distY*distY)
			if dist > 0.8 {
				brightness *= 1.0 - (dist-0.8)*2
			}
		}

		// Apply brightness
		if brightness < 0.5 {
			style = style.Faint(true)
		}

		result.WriteString(style.Render(string(ch)))
	}

	return result.String()
}

// CRTTickMsg is sent for animation
type CRTTickMsg struct{}

