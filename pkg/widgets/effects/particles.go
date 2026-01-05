// Package effects - Particle system (inspired by asciimatics)
package effects

import (
	"math"
	"math/rand"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Particle represents a single particle
type Particle struct {
	X, Y     float64
	VX, VY   float64
	Life     float64
	MaxLife  float64
	Char     rune
	Color    string
}

// ParticleSystem manages particles
type ParticleSystem struct {
	ID         string
	Width      int
	Height     int
	Particles  []Particle
	MaxParticles int
	EmitRate   int
	Gravity    float64
	EmitterX   float64
	EmitterY   float64
	Spread     float64
	Speed      time.Duration
	chars      []rune
	colors     []string
}

// NewParticleSystem creates a particle system
func NewParticleSystem(width, height int) *ParticleSystem {
	return &ParticleSystem{
		ID:           "particles",
		Width:        width,
		Height:       height,
		Particles:    []Particle{},
		MaxParticles: 100,
		EmitRate:     5,
		Gravity:      0.1,
		EmitterX:     float64(width) / 2,
		EmitterY:     float64(height),
		Spread:       0.5,
		Speed:        50 * time.Millisecond,
		chars:        []rune{'*', '•', '◦', '○', '✦', '✧', '✶', '✷'},
		colors:       []string{"#E040FB", "#9D4EDD", "#7B2CBF", "#5A189A", "#FFFFFF"},
	}
}

// SetEmitter sets emitter position
func (ps *ParticleSystem) SetEmitter(x, y float64) *ParticleSystem {
	ps.EmitterX = x
	ps.EmitterY = y
	return ps
}

// SetGravity sets gravity
func (ps *ParticleSystem) SetGravity(gravity float64) *ParticleSystem {
	ps.Gravity = gravity
	return ps
}

// SetColors sets particle colors
func (ps *ParticleSystem) SetColors(colors []string) *ParticleSystem {
	ps.colors = colors
	return ps
}

// emit creates new particles
func (ps *ParticleSystem) emit() {
	for i := 0; i < ps.EmitRate && len(ps.Particles) < ps.MaxParticles; i++ {
		angle := -math.Pi/2 + (rand.Float64()-0.5)*ps.Spread*math.Pi
		speed := 1 + rand.Float64()*2

		ps.Particles = append(ps.Particles, Particle{
			X:       ps.EmitterX,
			Y:       ps.EmitterY,
			VX:      math.Cos(angle) * speed,
			VY:      math.Sin(angle) * speed,
			Life:    1.0,
			MaxLife: 1.0,
			Char:    ps.chars[rand.Intn(len(ps.chars))],
			Color:   ps.colors[rand.Intn(len(ps.colors))],
		})
	}
}

// update moves particles
func (ps *ParticleSystem) update() {
	var alive []Particle

	for _, p := range ps.Particles {
		p.VY += ps.Gravity
		p.X += p.VX
		p.Y += p.VY
		p.Life -= 0.02

		if p.Life > 0 && p.Y >= 0 && p.Y < float64(ps.Height) &&
			p.X >= 0 && p.X < float64(ps.Width) {
			alive = append(alive, p)
		}
	}

	ps.Particles = alive
}

// Update handles messages
func (ps *ParticleSystem) Update(msg tea.Msg) (*ParticleSystem, tea.Cmd) {
	switch msg.(type) {
	case ParticleTickMsg:
		ps.emit()
		ps.update()
	}
	return ps, nil
}

// TickCmd returns animation command
func (ps *ParticleSystem) TickCmd() tea.Cmd {
	return tea.Tick(ps.Speed, func(t time.Time) tea.Msg {
		return ParticleTickMsg{}
	})
}

// View renders particles
func (ps *ParticleSystem) View() string {
	grid := make([][]string, ps.Height)
	for i := range grid {
		grid[i] = make([]string, ps.Width)
		for j := range grid[i] {
			grid[i][j] = " "
		}
	}

	for _, p := range ps.Particles {
		x, y := int(p.X), int(p.Y)
		if x >= 0 && x < ps.Width && y >= 0 && y < ps.Height {
			// Fade color based on life
			style := lipgloss.NewStyle().Foreground(lipgloss.Color(p.Color))
			if p.Life < 0.5 {
				style = style.Faint(true)
			}
			grid[y][x] = style.Render(string(p.Char))
		}
	}

	var lines []string
	for _, row := range grid {
		lines = append(lines, strings.Join(row, ""))
	}
	return strings.Join(lines, "\n")
}

// ParticleTickMsg is sent for animation
type ParticleTickMsg struct{}

// Firework creates a firework effect
type Firework struct {
	*ParticleSystem
	exploded bool
	targetY  float64
}

// NewFirework creates a firework
func NewFirework(width, height int) *Firework {
	ps := NewParticleSystem(width, height)
	ps.EmitterY = float64(height)
	ps.EmitterX = float64(width/2) + float64(rand.Intn(width/2)-width/4)
	ps.Gravity = 0.05
	ps.EmitRate = 0

	// Launch particle
	ps.Particles = append(ps.Particles, Particle{
		X:       ps.EmitterX,
		Y:       float64(height),
		VX:      0,
		VY:      -2 - rand.Float64(),
		Life:    1.0,
		MaxLife: 1.0,
		Char:    '▲',
		Color:   "#FFFFFF",
	})

	return &Firework{
		ParticleSystem: ps,
		targetY:        float64(height/4) + float64(rand.Intn(height/4)),
	}
}

// Update handles firework animation
func (fw *Firework) Update(msg tea.Msg) (*Firework, tea.Cmd) {
	switch msg.(type) {
	case ParticleTickMsg:
		if !fw.exploded && len(fw.Particles) > 0 {
			if fw.Particles[0].Y <= fw.targetY {
				fw.explode()
			}
		}
		fw.ParticleSystem.update()
	}
	return fw, nil
}

// explode creates explosion particles
func (fw *Firework) explode() {
	if len(fw.Particles) == 0 {
		return
	}

	x, y := fw.Particles[0].X, fw.Particles[0].Y
	fw.Particles = []Particle{}
	fw.exploded = true

	// Create explosion
	for i := 0; i < 30; i++ {
		angle := float64(i) * (2 * math.Pi / 30)
		speed := 0.5 + rand.Float64()*1.5

		fw.Particles = append(fw.Particles, Particle{
			X:       x,
			Y:       y,
			VX:      math.Cos(angle) * speed,
			VY:      math.Sin(angle) * speed,
			Life:    1.0,
			MaxLife: 1.0,
			Char:    fw.chars[rand.Intn(len(fw.chars))],
			Color:   fw.colors[rand.Intn(len(fw.colors))],
		})
	}
}

