// Package effects provides visual effects (Matrix, Fire, Particles, etc.)
package effects

import (
	"math/rand"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// MatrixRain creates the classic matrix rain effect (inspired by asciimatics)
type MatrixRain struct {
	ID       string
	Width    int
	Height   int
	Density  float64
	Speed    time.Duration
	columns  []matrixColumn
	chars    []rune
	style    MatrixStyle
}

type matrixColumn struct {
	y       int
	length  int
	speed   int
	chars   []rune
	active  bool
}

// MatrixStyle holds matrix styling
type MatrixStyle struct {
	Head lipgloss.Style
	Body lipgloss.Style
	Tail lipgloss.Style
}

// DefaultMatrixStyle returns default styling
func DefaultMatrixStyle() MatrixStyle {
	return MatrixStyle{
		Head: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Bold(true),
		Body: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00")),
		Tail: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#006600")),
	}
}

// MatrixChars default character set
var MatrixChars = []rune("アイウエオカキクケコサシスセソタチツテトナニヌネノハヒフヘホマミムメモヤユヨラリルレロワヲン0123456789")

// NewMatrixRain creates a matrix rain effect
func NewMatrixRain(width, height int) *MatrixRain {
	m := &MatrixRain{
		ID:      "matrix",
		Width:   width,
		Height:  height,
		Density: 0.1,
		Speed:   50 * time.Millisecond,
		chars:   MatrixChars,
		style:   DefaultMatrixStyle(),
		columns: make([]matrixColumn, width),
	}
	m.init()
	return m
}

// init initializes columns
func (m *MatrixRain) init() {
	for i := range m.columns {
		m.columns[i] = matrixColumn{
			y:      -rand.Intn(m.Height),
			length: 5 + rand.Intn(15),
			speed:  1 + rand.Intn(3),
			chars:  make([]rune, m.Height),
			active: rand.Float64() < m.Density,
		}
		for j := range m.columns[i].chars {
			m.columns[i].chars[j] = m.chars[rand.Intn(len(m.chars))]
		}
	}
}

// Update advances the animation
func (m *MatrixRain) Update(msg tea.Msg) (*MatrixRain, tea.Cmd) {
	switch msg.(type) {
	case MatrixTickMsg:
		for i := range m.columns {
			col := &m.columns[i]
			if col.active {
				col.y += col.speed
				if col.y-col.length > m.Height {
					col.y = -rand.Intn(m.Height / 2)
					col.length = 5 + rand.Intn(15)
					col.active = rand.Float64() < m.Density
					// Randomize chars
					for j := range col.chars {
						col.chars[j] = m.chars[rand.Intn(len(m.chars))]
					}
				}
			} else {
				if rand.Float64() < 0.02 {
					col.active = true
				}
			}
		}
	}
	return m, nil
}

// TickCmd returns a command for animation
func (m *MatrixRain) TickCmd() tea.Cmd {
	return tea.Tick(m.Speed, func(t time.Time) tea.Msg {
		return MatrixTickMsg{}
	})
}

// View renders the matrix
func (m *MatrixRain) View() string {
	grid := make([][]string, m.Height)
	for i := range grid {
		grid[i] = make([]string, m.Width)
		for j := range grid[i] {
			grid[i][j] = " "
		}
	}

	for x, col := range m.columns {
		if !col.active {
			continue
		}
		for dy := 0; dy < col.length; dy++ {
			y := col.y - dy
			if y >= 0 && y < m.Height {
				char := string(col.chars[y%len(col.chars)])
				if dy == 0 {
					grid[y][x] = m.style.Head.Render(char)
				} else if dy < col.length/3 {
					grid[y][x] = m.style.Body.Render(char)
				} else {
					grid[y][x] = m.style.Tail.Render(char)
				}
			}
		}
	}

	var lines []string
	for _, row := range grid {
		lines = append(lines, strings.Join(row, ""))
	}
	return strings.Join(lines, "\n")
}

// MatrixTickMsg is sent for animation
type MatrixTickMsg struct{}

// SetStyle sets the matrix style
func (m *MatrixRain) SetStyle(style MatrixStyle) *MatrixRain {
	m.style = style
	return m
}

// UltravioletMatrixStyle returns MakeaTUI themed matrix
func UltravioletMatrixStyle() MatrixStyle {
	return MatrixStyle{
		Head: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E040FB")).
			Bold(true),
		Body: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9D4EDD")),
		Tail: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#3C096C")),
	}
}

