// Package effects - Enhanced Matrix rain (inspired by ScaLaMatrixRain)
package effects

import (
	"math/rand"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// MatrixColumn represents a single column of falling characters
type MatrixColumn struct {
	chars     []rune
	positions []int
	speeds    []float64
	lengths   []int
	active    []bool
}

// EnhancedMatrix provides advanced Matrix rain effect
type EnhancedMatrix struct {
	ID          string
	Width       int
	Height      int
	columns     []MatrixColumn
	CharSets    [][]rune
	ActiveSet   int
	Speed       time.Duration
	Density     float64
	TrailLength int
	ColorMode   MatrixColorMode
	style       EnhancedMatrixStyle
}

// MatrixColorMode defines color modes
type MatrixColorMode int

const (
	MatrixColorClassic MatrixColorMode = iota // Green
	MatrixColorAmber                          // Amber/Orange
	MatrixColorCyan                           // Cyan/Blue
	MatrixColorRainbow                        // Rainbow
	MatrixColorUltraviolet                    // MakeaTUI theme
)

// EnhancedMatrixStyle holds matrix styling
type EnhancedMatrixStyle struct {
	HeadColors  []string
	BodyColors  []string
	TailColors  []string
	Background  lipgloss.Style
}

// ClassicMatrixStyle returns classic green style
func ClassicMatrixStyle() EnhancedMatrixStyle {
	return EnhancedMatrixStyle{
		HeadColors: []string{"#FFFFFF", "#AAFFAA"},
		BodyColors: []string{"#00FF00", "#00DD00", "#00BB00"},
		TailColors: []string{"#009900", "#006600", "#003300"},
		Background: lipgloss.NewStyle().Background(lipgloss.Color("#000000")),
	}
}

// UltravioletEnhancedMatrixStyle returns MakeaTUI themed style
func UltravioletEnhancedMatrixStyle() EnhancedMatrixStyle {
	return EnhancedMatrixStyle{
		HeadColors: []string{"#FFFFFF", "#FFB0FF"},
		BodyColors: []string{"#E040FB", "#9D4EDD", "#7B2CBF"},
		TailColors: []string{"#5A189A", "#3C096C", "#1A0533"},
		Background: lipgloss.NewStyle().Background(lipgloss.Color("#0D0221")),
	}
}

// Character sets
var (
	KatakanaChars = []rune("アイウエオカキクケコサシスセソタチツテトナニヌネノハヒフヘホマミムメモヤユヨラリルレロワヲン")
	LatinChars    = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	BinaryChars   = []rune("01")
	HexChars      = []rune("0123456789ABCDEF")
	SymbolChars   = []rune("!@#$%^&*()[]{}|;:,.<>?/~`")
	MixedChars    = []rune("アイウエオABCDE01234!@#$%")
)

// NewEnhancedMatrix creates enhanced matrix effect
func NewEnhancedMatrix(width, height int) *EnhancedMatrix {
	m := &EnhancedMatrix{
		ID:          "matrix",
		Width:       width,
		Height:      height,
		CharSets:    [][]rune{KatakanaChars, LatinChars, BinaryChars},
		ActiveSet:   0,
		Speed:       50 * time.Millisecond,
		Density:     0.3,
		TrailLength: 15,
		ColorMode:   MatrixColorClassic,
		style:       ClassicMatrixStyle(),
	}
	m.initColumns()
	return m
}

// initColumns initializes matrix columns
func (m *EnhancedMatrix) initColumns() {
	m.columns = make([]MatrixColumn, m.Width)
	for i := range m.columns {
		m.columns[i] = MatrixColumn{
			chars:     make([]rune, m.Height),
			positions: []int{},
			speeds:    []float64{},
			lengths:   []int{},
			active:    []bool{},
		}
		// Initialize with random characters
		for j := range m.columns[i].chars {
			m.columns[i].chars[j] = m.randomChar()
		}
	}
}

// randomChar returns a random character from active set
func (m *EnhancedMatrix) randomChar() rune {
	chars := m.CharSets[m.ActiveSet]
	return chars[rand.Intn(len(chars))]
}

// SetColorMode sets the color mode
func (m *EnhancedMatrix) SetColorMode(mode MatrixColorMode) *EnhancedMatrix {
	m.ColorMode = mode
	switch mode {
	case MatrixColorClassic:
		m.style = ClassicMatrixStyle()
	case MatrixColorUltraviolet:
		m.style = UltravioletEnhancedMatrixStyle()
	}
	return m
}

// SetCharSet sets the active character set
func (m *EnhancedMatrix) SetCharSet(index int) *EnhancedMatrix {
	if index >= 0 && index < len(m.CharSets) {
		m.ActiveSet = index
	}
	return m
}

// Update handles animation
func (m *EnhancedMatrix) Update(msg tea.Msg) (*EnhancedMatrix, tea.Cmd) {
	switch msg.(type) {
	case EnhancedMatrixTickMsg:
		// Spawn new drops
		for i := range m.columns {
			if rand.Float64() < m.Density*0.1 {
				m.columns[i].positions = append(m.columns[i].positions, 0)
				m.columns[i].speeds = append(m.columns[i].speeds, 0.5+rand.Float64())
				m.columns[i].lengths = append(m.columns[i].lengths, 5+rand.Intn(m.TrailLength))
				m.columns[i].active = append(m.columns[i].active, true)
			}
		}

		// Update drops
		for i := range m.columns {
			for j := range m.columns[i].positions {
				if m.columns[i].active[j] {
					m.columns[i].positions[j] += int(m.columns[i].speeds[j])
					if m.columns[i].positions[j] > m.Height+m.columns[i].lengths[j] {
						m.columns[i].active[j] = false
					}
				}
			}

			// Randomly change characters
			for j := range m.columns[i].chars {
				if rand.Float64() < 0.02 {
					m.columns[i].chars[j] = m.randomChar()
				}
			}
		}
	}
	return m, nil
}

// TickCmd returns animation command
func (m *EnhancedMatrix) TickCmd() tea.Cmd {
	return tea.Tick(m.Speed, func(t time.Time) tea.Msg {
		return EnhancedMatrixTickMsg{}
	})
}

// View renders the matrix
func (m *EnhancedMatrix) View() string {
	canvas := make([][]rune, m.Height)
	colors := make([][]int, m.Height) // 0=none, 1=head, 2=body, 3=tail

	for y := range canvas {
		canvas[y] = make([]rune, m.Width)
		colors[y] = make([]int, m.Width)
		for x := range canvas[y] {
			canvas[y][x] = ' '
		}
	}

	// Draw drops
	for x, col := range m.columns {
		for i, pos := range col.positions {
			if !col.active[i] {
				continue
			}
			length := col.lengths[i]

			for offset := 0; offset < length; offset++ {
				y := pos - offset
				if y >= 0 && y < m.Height {
					canvas[y][x] = col.chars[y]
					if offset == 0 {
						colors[y][x] = 1 // Head
					} else if offset < length/2 {
						colors[y][x] = 2 // Body
					} else {
						colors[y][x] = 3 // Tail
					}
				}
			}
		}
	}

	// Render
	var lines []string
	for y := range canvas {
		var line strings.Builder
		for x := range canvas[y] {
			ch := string(canvas[y][x])
			switch colors[y][x] {
			case 1:
				style := lipgloss.NewStyle().Foreground(lipgloss.Color(m.style.HeadColors[0]))
				line.WriteString(style.Render(ch))
			case 2:
				style := lipgloss.NewStyle().Foreground(lipgloss.Color(m.style.BodyColors[0]))
				line.WriteString(style.Render(ch))
			case 3:
				style := lipgloss.NewStyle().Foreground(lipgloss.Color(m.style.TailColors[0]))
				line.WriteString(style.Render(ch))
			default:
				line.WriteString(ch)
			}
		}
		lines = append(lines, line.String())
	}

	return m.style.Background.Render(strings.Join(lines, "\n"))
}

// EnhancedMatrixTickMsg is sent for animation
type EnhancedMatrixTickMsg struct{}

