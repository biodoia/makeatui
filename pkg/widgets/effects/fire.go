// Package effects - Fire and Plasma effects (inspired by asciimatics)
package effects

import (
	"math"
	"math/rand"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Fire creates an ASCII fire effect
type Fire struct {
	ID      string
	Width   int
	Height  int
	buffer  [][]float64
	Speed   time.Duration
	Palette []string
	Chars   []rune
}

// DefaultFirePalette provides fire colors
var DefaultFirePalette = []string{
	"#000000", "#1A0000", "#330000", "#4D0000", "#660000",
	"#800000", "#990000", "#B30000", "#CC0000", "#E60000",
	"#FF0000", "#FF1A00", "#FF3300", "#FF4D00", "#FF6600",
	"#FF8000", "#FF9900", "#FFB300", "#FFCC00", "#FFE600",
	"#FFFF00", "#FFFF33", "#FFFF66", "#FFFF99", "#FFFFCC",
}

// FireChars for rendering
var FireChars = []rune{' ', '.', ':', '*', 's', 'S', '#', '$', '&', '@'}

// NewFire creates a fire effect
func NewFire(width, height int) *Fire {
	f := &Fire{
		ID:      "fire",
		Width:   width,
		Height:  height,
		Speed:   50 * time.Millisecond,
		Palette: DefaultFirePalette,
		Chars:   FireChars,
		buffer:  make([][]float64, height),
	}
	for i := range f.buffer {
		f.buffer[i] = make([]float64, width)
	}
	return f
}

// Update advances the fire animation
func (f *Fire) Update(msg tea.Msg) (*Fire, tea.Cmd) {
	switch msg.(type) {
	case FireTickMsg:
		// Set bottom row to random values (heat source)
		for x := 0; x < f.Width; x++ {
			f.buffer[f.Height-1][x] = float64(rand.Intn(len(f.Palette)))
		}

		// Propagate fire upward with cooling
		for y := 0; y < f.Height-1; y++ {
			for x := 0; x < f.Width; x++ {
				// Average of neighbors below
				x1 := (x - 1 + f.Width) % f.Width
				x2 := (x + 1) % f.Width
				y1 := y + 1

				sum := f.buffer[y1][x1] + f.buffer[y1][x] + f.buffer[y1][x2]
				if y+2 < f.Height {
					sum += f.buffer[y+2][x]
					sum /= 4.0
				} else {
					sum /= 3.0
				}

				// Cooling
				cooling := rand.Float64() * 2
				f.buffer[y][x] = math.Max(0, sum-cooling)
			}
		}
	}
	return f, nil
}

// TickCmd returns animation command
func (f *Fire) TickCmd() tea.Cmd {
	return tea.Tick(f.Speed, func(t time.Time) tea.Msg {
		return FireTickMsg{}
	})
}

// View renders the fire
func (f *Fire) View() string {
	var lines []string

	for y := 0; y < f.Height; y++ {
		var line strings.Builder
		for x := 0; x < f.Width; x++ {
			val := int(f.buffer[y][x])
			if val >= len(f.Palette) {
				val = len(f.Palette) - 1
			}
			if val < 0 {
				val = 0
			}

			charIdx := val * len(f.Chars) / len(f.Palette)
			if charIdx >= len(f.Chars) {
				charIdx = len(f.Chars) - 1
			}

			style := lipgloss.NewStyle().
				Foreground(lipgloss.Color(f.Palette[val]))
			line.WriteString(style.Render(string(f.Chars[charIdx])))
		}
		lines = append(lines, line.String())
	}

	return strings.Join(lines, "\n")
}

// FireTickMsg is sent for animation
type FireTickMsg struct{}

// Plasma creates a plasma effect
type Plasma struct {
	ID      string
	Width   int
	Height  int
	time    float64
	Speed   time.Duration
	Palette []string
	Chars   []rune
}

// DefaultPlasmaPalette provides plasma colors
var DefaultPlasmaPalette = []string{
	"#0D0221", "#1A0533", "#240046", "#3C096C", "#5A189A",
	"#7B2CBF", "#9D4EDD", "#C77DFF", "#E0AAFF", "#E040FB",
	"#9D4EDD", "#7B2CBF", "#5A189A", "#3C096C", "#240046",
}

// PlasmaChars for rendering
var PlasmaChars = []rune{'░', '▒', '▓', '█'}

// NewPlasma creates a plasma effect
func NewPlasma(width, height int) *Plasma {
	return &Plasma{
		ID:      "plasma",
		Width:   width,
		Height:  height,
		Speed:   100 * time.Millisecond,
		Palette: DefaultPlasmaPalette,
		Chars:   PlasmaChars,
	}
}

// Update advances the plasma animation
func (p *Plasma) Update(msg tea.Msg) (*Plasma, tea.Cmd) {
	switch msg.(type) {
	case PlasmaTickMsg:
		p.time += 0.1
	}
	return p, nil
}

// TickCmd returns animation command
func (p *Plasma) TickCmd() tea.Cmd {
	return tea.Tick(p.Speed, func(t time.Time) tea.Msg {
		return PlasmaTickMsg{}
	})
}

// View renders the plasma
func (p *Plasma) View() string {
	var lines []string

	for y := 0; y < p.Height; y++ {
		var line strings.Builder
		for x := 0; x < p.Width; x++ {
			// Plasma formula
			v := math.Sin(float64(x)*0.1 + p.time)
			v += math.Sin(float64(y)*0.1 + p.time*0.5)
			v += math.Sin((float64(x)+float64(y))*0.1 + p.time*0.3)
			v += math.Sin(math.Sqrt(float64(x*x+y*y))*0.1 + p.time*0.2)

			// Normalize to 0-1
			v = (v + 4) / 8

			// Map to palette
			colorIdx := int(v * float64(len(p.Palette)-1))
			if colorIdx >= len(p.Palette) {
				colorIdx = len(p.Palette) - 1
			}

			charIdx := int(v * float64(len(p.Chars)-1))
			if charIdx >= len(p.Chars) {
				charIdx = len(p.Chars) - 1
			}

			style := lipgloss.NewStyle().
				Foreground(lipgloss.Color(p.Palette[colorIdx]))
			line.WriteString(style.Render(string(p.Chars[charIdx])))
		}
		lines = append(lines, line.String())
	}

	return strings.Join(lines, "\n")
}

// PlasmaTickMsg is sent for animation
type PlasmaTickMsg struct{}

