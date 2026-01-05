// Package layout - Split pane layouts
package layout

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// SplitDirection defines split orientation
type SplitDirection int

const (
	SplitHorizontal SplitDirection = iota
	SplitVertical
)

// Split implements a split pane layout (inspired by Ratatui/AppCUI)
type Split struct {
	Direction   SplitDirection
	Ratio       float64 // 0.0 to 1.0, position of split
	MinSize     int     // minimum size for either pane
	ShowDivider bool
	DividerChar string
	FirstPane   string
	SecondPane  string
	width       int
	height      int
	style       lipgloss.Style
}

// NewSplit creates a new split layout
func NewSplit(direction SplitDirection, ratio float64) *Split {
	return &Split{
		Direction:   direction,
		Ratio:       ratio,
		MinSize:     1,
		ShowDivider: true,
		DividerChar: "│",
		style:       lipgloss.NewStyle(),
	}
}

// SetPanes sets both pane contents
func (s *Split) SetPanes(first, second string) *Split {
	s.FirstPane = first
	s.SecondPane = second
	return s
}

// SetSize sets the split size
func (s *Split) SetSize(width, height int) *Split {
	s.width = width
	s.height = height
	return s
}

// SetDivider configures the divider
func (s *Split) SetDivider(show bool, char string) *Split {
	s.ShowDivider = show
	s.DividerChar = char
	return s
}

// Render renders the split layout
func (s *Split) Render() string {
	dividerSize := 0
	if s.ShowDivider {
		dividerSize = 1
	}

	switch s.Direction {
	case SplitHorizontal:
		return s.renderHorizontal(dividerSize)
	case SplitVertical:
		return s.renderVertical(dividerSize)
	}
	return ""
}

func (s *Split) renderHorizontal(dividerSize int) string {
	firstWidth := int(float64(s.width) * s.Ratio)
	if firstWidth < s.MinSize {
		firstWidth = s.MinSize
	}
	secondWidth := s.width - firstWidth - dividerSize
	if secondWidth < s.MinSize {
		secondWidth = s.MinSize
		firstWidth = s.width - secondWidth - dividerSize
	}

	firstStyle := lipgloss.NewStyle().Width(firstWidth).Height(s.height)
	secondStyle := lipgloss.NewStyle().Width(secondWidth).Height(s.height)

	first := firstStyle.Render(s.FirstPane)
	second := secondStyle.Render(s.SecondPane)

	if s.ShowDivider {
		divider := strings.Repeat(s.DividerChar+"\n", s.height)
		divider = strings.TrimSuffix(divider, "\n")
		return lipgloss.JoinHorizontal(lipgloss.Top, first, divider, second)
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, first, second)
}

func (s *Split) renderVertical(dividerSize int) string {
	firstHeight := int(float64(s.height) * s.Ratio)
	if firstHeight < s.MinSize {
		firstHeight = s.MinSize
	}
	secondHeight := s.height - firstHeight - dividerSize
	if secondHeight < s.MinSize {
		secondHeight = s.MinSize
		firstHeight = s.height - secondHeight - dividerSize
	}

	firstStyle := lipgloss.NewStyle().Width(s.width).Height(firstHeight)
	secondStyle := lipgloss.NewStyle().Width(s.width).Height(secondHeight)

	first := firstStyle.Render(s.FirstPane)
	second := secondStyle.Render(s.SecondPane)

	if s.ShowDivider {
		divider := strings.Repeat("─", s.width)
		return lipgloss.JoinVertical(lipgloss.Left, first, divider, second)
	}
	return lipgloss.JoinVertical(lipgloss.Left, first, second)
}

// TripleSplit creates a three-pane layout
type TripleSplit struct {
	Direction SplitDirection
	Ratio1    float64 // first split
	Ratio2    float64 // second split (of remaining space)
	Panes     [3]string
	width     int
	height    int
}

// NewTripleSplit creates a new triple split
func NewTripleSplit(direction SplitDirection) *TripleSplit {
	return &TripleSplit{
		Direction: direction,
		Ratio1:    0.33,
		Ratio2:    0.5,
	}
}

// SetPanes sets all three panes
func (t *TripleSplit) SetPanes(first, second, third string) *TripleSplit {
	t.Panes = [3]string{first, second, third}
	return t
}

// SetSize sets dimensions
func (t *TripleSplit) SetSize(width, height int) *TripleSplit {
	t.width = width
	t.height = height
	return t
}

// Render renders the triple split
func (t *TripleSplit) Render() string {
	split1 := NewSplit(t.Direction, t.Ratio1)
	split2 := NewSplit(t.Direction, t.Ratio2)

	if t.Direction == SplitHorizontal {
		remainingWidth := t.width - int(float64(t.width)*t.Ratio1) - 1
		split2.SetSize(remainingWidth, t.height)
		split2.SetPanes(t.Panes[1], t.Panes[2])
		secondContent := split2.Render()

		split1.SetSize(t.width, t.height)
		split1.SetPanes(t.Panes[0], secondContent)
	} else {
		remainingHeight := t.height - int(float64(t.height)*t.Ratio1) - 1
		split2.SetSize(t.width, remainingHeight)
		split2.SetPanes(t.Panes[1], t.Panes[2])
		secondContent := split2.Render()

		split1.SetSize(t.width, t.height)
		split1.SetPanes(t.Panes[0], secondContent)
	}

	return split1.Render()
}

