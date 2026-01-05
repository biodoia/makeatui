// Package feedback - Loading and skeleton components
package feedback

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// SpinnerType defines spinner animation type
type SpinnerType int

const (
	SpinnerDots SpinnerType = iota
	SpinnerLine
	SpinnerPulse
	SpinnerGlobe
	SpinnerMoon
	SpinnerMonkey
	SpinnerMeter
	SpinnerHamburg
)

// SpinnerFrames for different spinner types
var SpinnerFrames = map[SpinnerType][]string{
	SpinnerDots:    {"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
	SpinnerLine:    {"-", "\\", "|", "/"},
	SpinnerPulse:   {"█", "▓", "▒", "░", "▒", "▓"},
	SpinnerGlobe:   {"🌍", "🌎", "🌏"},
	SpinnerMoon:    {"🌑", "🌒", "🌓", "🌔", "🌕", "🌖", "🌗", "🌘"},
	SpinnerMonkey:  {"🙈", "🙉", "🙊"},
	SpinnerMeter:   {"▱▱▱", "▰▱▱", "▰▰▱", "▰▰▰", "▰▰▱", "▰▱▱"},
	SpinnerHamburg: {"☱", "☲", "☴"},
}

// Spinner provides an animated loading indicator
type Spinner struct {
	ID       string
	Label    string
	Type     SpinnerType
	Frame    int
	Visible  bool
	Interval time.Duration
	style    SpinnerStyle
}

// SpinnerStyle holds spinner styling
type SpinnerStyle struct {
	Spinner lipgloss.Style
	Label   lipgloss.Style
}

// DefaultSpinnerStyle returns default styling
func DefaultSpinnerStyle() SpinnerStyle {
	return SpinnerStyle{
		Spinner: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E040FB")),
		Label: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")),
	}
}

// NewSpinner creates a new spinner
func NewSpinner(id string) *Spinner {
	return &Spinner{
		ID:       id,
		Type:     SpinnerDots,
		Visible:  true,
		Interval: 100 * time.Millisecond,
		style:    DefaultSpinnerStyle(),
	}
}

// SetLabel sets the spinner label
func (s *Spinner) SetLabel(label string) *Spinner {
	s.Label = label
	return s
}

// SetType sets the spinner type
func (s *Spinner) SetType(spinnerType SpinnerType) *Spinner {
	s.Type = spinnerType
	s.Frame = 0
	return s
}

// Tick advances the spinner animation
func (s *Spinner) Tick() *Spinner {
	frames := SpinnerFrames[s.Type]
	s.Frame = (s.Frame + 1) % len(frames)
	return s
}

// TickCmd returns a command for animation
func (s *Spinner) TickCmd() tea.Cmd {
	return tea.Tick(s.Interval, func(t time.Time) tea.Msg {
		return SpinnerTickMsg{ID: s.ID}
	})
}

// View renders the spinner
func (s *Spinner) View() string {
	if !s.Visible {
		return ""
	}

	frames := SpinnerFrames[s.Type]
	frame := frames[s.Frame%len(frames)]
	spinner := s.style.Spinner.Render(frame)

	if s.Label != "" {
		return spinner + " " + s.style.Label.Render(s.Label)
	}
	return spinner
}

// SpinnerTickMsg is sent for spinner animation
type SpinnerTickMsg struct {
	ID string
}

// Skeleton provides placeholder loading state (inspired by React/Textual)
type Skeleton struct {
	ID      string
	Width   int
	Height  int
	Type    SkeletonType
	Animate bool
	style   lipgloss.Style
}

// SkeletonType defines skeleton shape
type SkeletonType int

const (
	SkeletonText SkeletonType = iota
	SkeletonRectangle
	SkeletonCircle
	SkeletonCard
)

// NewSkeleton creates a skeleton placeholder
func NewSkeleton(id string) *Skeleton {
	return &Skeleton{
		ID:      id,
		Width:   20,
		Height:  1,
		Type:    SkeletonText,
		Animate: true,
		style: lipgloss.NewStyle().
			Background(lipgloss.Color("#3C096C")).
			Foreground(lipgloss.Color("#5A189A")),
	}
}

// SetWidth sets skeleton width
func (s *Skeleton) SetWidth(width int) *Skeleton {
	s.Width = width
	return s
}

// SetHeight sets skeleton height
func (s *Skeleton) SetHeight(height int) *Skeleton {
	s.Height = height
	return s
}

// SetType sets skeleton type
func (s *Skeleton) SetType(skeletonType SkeletonType) *Skeleton {
	s.Type = skeletonType
	return s
}

// View renders the skeleton
func (s *Skeleton) View() string {
	char := "░"
	line := strings.Repeat(char, s.Width)

	switch s.Type {
	case SkeletonText:
		return s.style.Render(line)
	case SkeletonRectangle:
		var lines []string
		for i := 0; i < s.Height; i++ {
			lines = append(lines, s.style.Render(line))
		}
		return strings.Join(lines, "\n")
	case SkeletonCircle:
		// Approximate circle with characters
		return s.style.Render("( " + strings.Repeat(char, s.Width-4) + " )")
	case SkeletonCard:
		var lines []string
		border := "┌" + strings.Repeat("─", s.Width-2) + "┐"
		lines = append(lines, s.style.Render(border))
		for i := 0; i < s.Height-2; i++ {
			lines = append(lines, s.style.Render("│"+strings.Repeat(char, s.Width-2)+"│"))
		}
		border = "└" + strings.Repeat("─", s.Width-2) + "┘"
		lines = append(lines, s.style.Render(border))
		return strings.Join(lines, "\n")
	}

	return s.style.Render(line)
}

// SkeletonGroup creates a group of skeletons
type SkeletonGroup struct {
	Skeletons []*Skeleton
	Gap       int
}

// NewSkeletonGroup creates a skeleton group
func NewSkeletonGroup() *SkeletonGroup {
	return &SkeletonGroup{
		Skeletons: []*Skeleton{},
		Gap:       1,
	}
}

// AddLine adds a text line skeleton
func (sg *SkeletonGroup) AddLine(width int) *SkeletonGroup {
	sg.Skeletons = append(sg.Skeletons, NewSkeleton("").SetWidth(width))
	return sg
}

// AddCard adds a card skeleton
func (sg *SkeletonGroup) AddCard(width, height int) *SkeletonGroup {
	sg.Skeletons = append(sg.Skeletons, NewSkeleton("").
		SetWidth(width).
		SetHeight(height).
		SetType(SkeletonCard))
	return sg
}

// View renders the skeleton group
func (sg *SkeletonGroup) View() string {
	var views []string
	for _, s := range sg.Skeletons {
		views = append(views, s.View())
	}
	return lipgloss.JoinVertical(lipgloss.Left, views...)
}

