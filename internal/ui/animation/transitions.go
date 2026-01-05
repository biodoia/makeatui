// Package animation - Transition effects for TUI components
package animation

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// TransitionType defines the type of transition
type TransitionType int

const (
	TransitionFade TransitionType = iota
	TransitionSlideLeft
	TransitionSlideRight
	TransitionSlideUp
	TransitionSlideDown
	TransitionZoom
	TransitionPop
)

// Transition manages animated transitions between views
type Transition struct {
	Type     TransitionType
	spring   *Spring
	progress float64
	active   bool
	from     string
	to       string
}

// NewTransition creates a new transition
func NewTransition(transitionType TransitionType) *Transition {
	return &Transition{
		Type:   transitionType,
		spring: NewSpring(SmoothSpring),
	}
}

// Start begins a transition from one view to another
func (t *Transition) Start(from, to string) {
	t.from = from
	t.to = to
	t.progress = 0
	t.spring.position = 0
	t.spring.SetTarget(1.0)
	t.active = true
}

// Update advances the transition
func (t *Transition) Update() tea.Cmd {
	if !t.active {
		return nil
	}

	t.progress = t.spring.Update()

	if t.spring.IsSettled() {
		t.active = false
		t.progress = 1.0
	}

	return nil
}

// View returns the current transition frame
func (t *Transition) View(width, height int) string {
	if !t.active {
		return t.to
	}

	switch t.Type {
	case TransitionFade:
		return t.fadeBetween()
	case TransitionSlideLeft:
		return t.slideHorizontal(width, true)
	case TransitionSlideRight:
		return t.slideHorizontal(width, false)
	case TransitionSlideUp:
		return t.slideVertical(height, true)
	case TransitionSlideDown:
		return t.slideVertical(height, false)
	case TransitionPop:
		return t.popTransition()
	default:
		return t.to
	}
}

// IsActive returns true if transition is in progress
func (t *Transition) IsActive() bool {
	return t.active
}

// Progress returns the current progress (0-1)
func (t *Transition) Progress() float64 {
	return t.progress
}

func (t *Transition) fadeBetween() string {
	// For TUI, we can't do true fading, so we crossfade by showing
	// more of the "to" content as progress increases
	if t.progress < 0.5 {
		return t.from
	}
	return t.to
}

func (t *Transition) slideHorizontal(width int, left bool) string {
	offset := int(float64(width) * (1 - t.progress))
	if !left {
		offset = -offset
	}

	_ = strings.Split(t.from, "\n") // fromLines reserved for future use
	toLines := strings.Split(t.to, "\n")

	var result []string
	maxLines := len(toLines)

	for i := 0; i < maxLines; i++ {
		toLine := ""
		if i < len(toLines) {
			toLine = toLines[i]
		}

		// Simple slide effect
		if left {
			result = append(result, padOrTruncate(toLine, width, offset))
		} else {
			result = append(result, padOrTruncate(toLine, width, offset))
		}
	}

	return strings.Join(result, "\n")
}

func (t *Transition) slideVertical(height int, up bool) string {
	offset := int(float64(height) * (1 - t.progress))
	if !up {
		offset = -offset
	}

	toLines := strings.Split(t.to, "\n")

	if up {
		// Slide content up from bottom
		padding := make([]string, max(0, height-len(toLines)-offset))
		return strings.Join(append(padding, toLines...), "\n")
	}

	// Slide content down from top
	padding := make([]string, max(0, offset))
	return strings.Join(append(padding, toLines...), "\n")
}

func (t *Transition) popTransition() string {
	// Pop effect - scale from center (simulated)
	if t.progress < 0.3 {
		return "" // Start empty
	}
	return t.to
}

func padOrTruncate(s string, width, offset int) string {
	if offset > 0 {
		padding := strings.Repeat(" ", min(offset, width))
		s = padding + s
	}
	if len(s) > width {
		s = s[:width]
	}
	return s
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

