// Package effects - Typing and text effects
package effects

import (
	"math/rand"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TypingEffect provides typewriter-style text animation
type TypingEffect struct {
	ID         string
	Text       string
	Displayed  string
	Cursor     string
	ShowCursor bool
	Speed      time.Duration
	Jitter     time.Duration
	Complete   bool
	charIndex  int
	style      lipgloss.Style
}

// NewTypingEffect creates a typing effect
func NewTypingEffect(text string) *TypingEffect {
	return &TypingEffect{
		ID:         "typing",
		Text:       text,
		Cursor:     "█",
		ShowCursor: true,
		Speed:      50 * time.Millisecond,
		Jitter:     30 * time.Millisecond,
		style: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E040FB")),
	}
}

// SetSpeed sets typing speed
func (t *TypingEffect) SetSpeed(speed time.Duration) *TypingEffect {
	t.Speed = speed
	return t
}

// SetCursor sets the cursor character
func (t *TypingEffect) SetCursor(cursor string) *TypingEffect {
	t.Cursor = cursor
	return t
}

// SetStyle sets text style
func (t *TypingEffect) SetStyle(style lipgloss.Style) *TypingEffect {
	t.style = style
	return t
}

// Reset resets the animation
func (t *TypingEffect) Reset() *TypingEffect {
	t.charIndex = 0
	t.Displayed = ""
	t.Complete = false
	return t
}

// Update advances the typing animation
func (t *TypingEffect) Update(msg tea.Msg) (*TypingEffect, tea.Cmd) {
	switch msg.(type) {
	case TypingTickMsg:
		if t.charIndex < len(t.Text) {
			t.charIndex++
			t.Displayed = t.Text[:t.charIndex]
		} else {
			t.Complete = true
		}
	}
	return t, nil
}

// TickCmd returns animation command
func (t *TypingEffect) TickCmd() tea.Cmd {
	jitter := time.Duration(rand.Int63n(int64(t.Jitter)))
	return tea.Tick(t.Speed+jitter, func(tm time.Time) tea.Msg {
		return TypingTickMsg{ID: t.ID}
	})
}

// View renders the typing effect
func (t *TypingEffect) View() string {
	cursor := ""
	if t.ShowCursor && !t.Complete {
		cursor = t.Cursor
	}
	return t.style.Render(t.Displayed + cursor)
}

// TypingTickMsg is sent for animation
type TypingTickMsg struct {
	ID string
}

// Gradient creates a color gradient text effect
type Gradient struct {
	Text   string
	Colors []string
}

// NewGradient creates a gradient text
func NewGradient(text string, colors []string) *Gradient {
	return &Gradient{
		Text:   text,
		Colors: colors,
	}
}

// UltravioletGradient returns MakeaTUI themed gradient
func UltravioletGradient() []string {
	return []string{
		"#3C096C", "#5A189A", "#7B2CBF", "#9D4EDD", "#C77DFF", "#E040FB",
	}
}

// NeonGradient returns neon gradient
func NeonGradient() []string {
	return []string{
		"#FF0080", "#FF00FF", "#8000FF", "#0080FF", "#00FFFF", "#00FF80",
	}
}

// View renders the gradient text
func (g *Gradient) View() string {
	if len(g.Colors) == 0 || len(g.Text) == 0 {
		return g.Text
	}

	var result strings.Builder
	for i, char := range g.Text {
		colorIdx := i * len(g.Colors) / len(g.Text)
		if colorIdx >= len(g.Colors) {
			colorIdx = len(g.Colors) - 1
		}
		style := lipgloss.NewStyle().Foreground(lipgloss.Color(g.Colors[colorIdx]))
		result.WriteString(style.Render(string(char)))
	}
	return result.String()
}

// Rainbow creates animated rainbow text
type Rainbow struct {
	Text    string
	offset  int
	Speed   time.Duration
	Colors  []string
}

// RainbowColors default rainbow
var RainbowColors = []string{
	"#FF0000", "#FF7F00", "#FFFF00", "#00FF00", "#0000FF", "#4B0082", "#9400D3",
}

// NewRainbow creates a rainbow effect
func NewRainbow(text string) *Rainbow {
	return &Rainbow{
		Text:   text,
		Speed:  100 * time.Millisecond,
		Colors: RainbowColors,
	}
}

// Update advances rainbow animation
func (r *Rainbow) Update(msg tea.Msg) (*Rainbow, tea.Cmd) {
	switch msg.(type) {
	case RainbowTickMsg:
		r.offset = (r.offset + 1) % len(r.Colors)
	}
	return r, nil
}

// TickCmd returns animation command
func (r *Rainbow) TickCmd() tea.Cmd {
	return tea.Tick(r.Speed, func(t time.Time) tea.Msg {
		return RainbowTickMsg{}
	})
}

// View renders rainbow text
func (r *Rainbow) View() string {
	if len(r.Colors) == 0 || len(r.Text) == 0 {
		return r.Text
	}

	var result strings.Builder
	for i, char := range r.Text {
		colorIdx := (i + r.offset) % len(r.Colors)
		style := lipgloss.NewStyle().Foreground(lipgloss.Color(r.Colors[colorIdx]))
		result.WriteString(style.Render(string(char)))
	}
	return result.String()
}

// RainbowTickMsg is sent for animation
type RainbowTickMsg struct{}

// Blink creates blinking text
type Blink struct {
	Text    string
	Visible bool
	Speed   time.Duration
	style   lipgloss.Style
}

// NewBlink creates blinking text
func NewBlink(text string) *Blink {
	return &Blink{
		Text:    text,
		Visible: true,
		Speed:   500 * time.Millisecond,
		style:   lipgloss.NewStyle().Foreground(lipgloss.Color("#E040FB")),
	}
}

// Update toggles visibility
func (b *Blink) Update(msg tea.Msg) (*Blink, tea.Cmd) {
	switch msg.(type) {
	case BlinkTickMsg:
		b.Visible = !b.Visible
	}
	return b, nil
}

// TickCmd returns animation command
func (b *Blink) TickCmd() tea.Cmd {
	return tea.Tick(b.Speed, func(t time.Time) tea.Msg {
		return BlinkTickMsg{}
	})
}

// View renders blinking text
func (b *Blink) View() string {
	if b.Visible {
		return b.style.Render(b.Text)
	}
	return strings.Repeat(" ", len(b.Text))
}

// BlinkTickMsg is sent for animation
type BlinkTickMsg struct{}

