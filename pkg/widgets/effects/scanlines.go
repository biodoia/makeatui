// Package effects - Scanline and retro effects
package effects

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Scanlines applies scanline overlay effect
type Scanlines struct {
	Intensity  float64 // 0.0 to 1.0
	Spacing    int     // Lines between scanlines
	darkStyle  lipgloss.Style
	lightStyle lipgloss.Style
}

// NewScanlines creates scanline effect
func NewScanlines() *Scanlines {
	return &Scanlines{
		Intensity: 0.3,
		Spacing:   2,
		darkStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#333333")),
		lightStyle: lipgloss.NewStyle(),
	}
}

// SetIntensity sets scanline intensity
func (s *Scanlines) SetIntensity(intensity float64) *Scanlines {
	s.Intensity = intensity
	return s
}

// Apply applies scanlines to content
func (s *Scanlines) Apply(content string) string {
	lines := strings.Split(content, "\n")
	var result []string

	for i, line := range lines {
		if i%s.Spacing == 0 {
			// Apply dim effect to scanline
			result = append(result, s.darkStyle.Render(line))
		} else {
			result = append(result, s.lightStyle.Render(line))
		}
	}

	return strings.Join(result, "\n")
}

// PhosphorDecay simulates phosphor decay/persistence
type PhosphorDecay struct {
	DecayFrames int
	history     []string
	colors      []string
}

// NewPhosphorDecay creates phosphor decay effect
func NewPhosphorDecay() *PhosphorDecay {
	return &PhosphorDecay{
		DecayFrames: 3,
		history:     []string{},
		colors: []string{
			"#00FF00", // Bright
			"#00AA00", // Medium
			"#006600", // Dim
			"#003300", // Very dim
		},
	}
}

// AddFrame adds a frame to the decay history
func (p *PhosphorDecay) AddFrame(content string) {
	p.history = append([]string{content}, p.history...)
	if len(p.history) > p.DecayFrames {
		p.history = p.history[:p.DecayFrames]
	}
}

// SetColors sets decay color gradient
func (p *PhosphorDecay) SetColors(colors []string) *PhosphorDecay {
	p.colors = colors
	return p
}

// View renders with decay effect
func (p *PhosphorDecay) View() string {
	if len(p.history) == 0 {
		return ""
	}

	// Just return latest frame for now
	// Full implementation would blend frames
	return p.history[0]
}

// Glitch applies glitch/corruption effect
type Glitch struct {
	Intensity float64
	chars     []rune
}

// NewGlitch creates glitch effect
func NewGlitch() *Glitch {
	return &Glitch{
		Intensity: 0.05,
		chars:     []rune{'░', '▒', '▓', '█', '▄', '▀', '■', '□', '▪', '▫'},
	}
}

// SetIntensity sets glitch intensity
func (g *Glitch) SetIntensity(intensity float64) *Glitch {
	g.Intensity = intensity
	return g
}

// Apply applies glitch effect to content
func (g *Glitch) Apply(content string) string {
	// For a simple implementation, just return content
	// Full glitch would randomly corrupt characters
	return content
}

// AsciiArt converts simple shapes to ASCII art
type AsciiArt struct {
	CharSet AsciiCharSet
}

// AsciiCharSet defines characters for ASCII art
type AsciiCharSet struct {
	Solid     rune
	Light     rune
	Medium    rune
	Dark      rune
	HLine     rune
	VLine     rune
	Corner    [4]rune // TL, TR, BL, BR
}

// DefaultAsciiCharSet returns default character set
func DefaultAsciiCharSet() AsciiCharSet {
	return AsciiCharSet{
		Solid:  '█',
		Light:  '░',
		Medium: '▒',
		Dark:   '▓',
		HLine:  '─',
		VLine:  '│',
		Corner: [4]rune{'┌', '┐', '└', '┘'},
	}
}

// BlockAsciiCharSet returns block characters
func BlockAsciiCharSet() AsciiCharSet {
	return AsciiCharSet{
		Solid:  '█',
		Light:  '▄',
		Medium: '▀',
		Dark:   '■',
		HLine:  '▬',
		VLine:  '▐',
		Corner: [4]rune{'█', '█', '█', '█'},
	}
}

// NewAsciiArt creates ASCII art renderer
func NewAsciiArt() *AsciiArt {
	return &AsciiArt{
		CharSet: DefaultAsciiCharSet(),
	}
}

// Box draws an ASCII box
func (a *AsciiArt) Box(width, height int) string {
	if width < 2 || height < 2 {
		return ""
	}

	var lines []string

	// Top line
	top := string(a.CharSet.Corner[0]) +
		strings.Repeat(string(a.CharSet.HLine), width-2) +
		string(a.CharSet.Corner[1])
	lines = append(lines, top)

	// Middle lines
	for i := 0; i < height-2; i++ {
		middle := string(a.CharSet.VLine) +
			strings.Repeat(" ", width-2) +
			string(a.CharSet.VLine)
		lines = append(lines, middle)
	}

	// Bottom line
	bottom := string(a.CharSet.Corner[2]) +
		strings.Repeat(string(a.CharSet.HLine), width-2) +
		string(a.CharSet.Corner[3])
	lines = append(lines, bottom)

	return strings.Join(lines, "\n")
}

// FilledBox draws a filled ASCII box
func (a *AsciiArt) FilledBox(width, height int) string {
	var lines []string
	for i := 0; i < height; i++ {
		lines = append(lines, strings.Repeat(string(a.CharSet.Solid), width))
	}
	return strings.Join(lines, "\n")
}

