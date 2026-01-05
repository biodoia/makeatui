// Package display - Status bar component
package display

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// StatusBar provides a status bar at bottom of screen
type StatusBar struct {
	ID       string
	Width    int
	Left     []StatusItem
	Center   []StatusItem
	Right    []StatusItem
	style    StatusBarStyle
}

// StatusItem represents a status bar item
type StatusItem struct {
	Text  string
	Icon  string
	Style lipgloss.Style
}

// StatusBarStyle holds status bar styling
type StatusBarStyle struct {
	Container lipgloss.Style
	Item      lipgloss.Style
	Separator lipgloss.Style
}

// DefaultStatusBarStyle returns default styling
func DefaultStatusBarStyle() StatusBarStyle {
	return StatusBarStyle{
		Container: lipgloss.NewStyle().
			Background(lipgloss.Color("#1A0533")).
			Foreground(lipgloss.Color("#FFFFFF")),
		Item: lipgloss.NewStyle().
			Padding(0, 1),
		Separator: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#5A189A")),
	}
}

// NewStatusBar creates a status bar
func NewStatusBar(id string, width int) *StatusBar {
	return &StatusBar{
		ID:     id,
		Width:  width,
		Left:   []StatusItem{},
		Center: []StatusItem{},
		Right:  []StatusItem{},
		style:  DefaultStatusBarStyle(),
	}
}

// AddLeft adds item to left section
func (s *StatusBar) AddLeft(text string) *StatusBar {
	s.Left = append(s.Left, StatusItem{Text: text})
	return s
}

// AddLeftWithIcon adds item with icon to left section
func (s *StatusBar) AddLeftWithIcon(icon, text string) *StatusBar {
	s.Left = append(s.Left, StatusItem{Text: text, Icon: icon})
	return s
}

// AddCenter adds item to center section
func (s *StatusBar) AddCenter(text string) *StatusBar {
	s.Center = append(s.Center, StatusItem{Text: text})
	return s
}

// AddRight adds item to right section
func (s *StatusBar) AddRight(text string) *StatusBar {
	s.Right = append(s.Right, StatusItem{Text: text})
	return s
}

// AddRightWithIcon adds item with icon to right section
func (s *StatusBar) AddRightWithIcon(icon, text string) *StatusBar {
	s.Right = append(s.Right, StatusItem{Text: text, Icon: icon})
	return s
}

// SetStyle sets status bar style
func (s *StatusBar) SetStyle(style StatusBarStyle) *StatusBar {
	s.style = style
	return s
}

// renderItems renders a list of items
func (s *StatusBar) renderItems(items []StatusItem) string {
	var parts []string
	for _, item := range items {
		text := item.Text
		if item.Icon != "" {
			text = item.Icon + " " + text
		}
		style := s.style.Item
		if item.Style.Value() != "" {
			style = item.Style
		}
		parts = append(parts, style.Render(text))
	}
	sep := s.style.Separator.Render(" │ ")
	return strings.Join(parts, sep)
}

// View renders the status bar
func (s *StatusBar) View() string {
	left := s.renderItems(s.Left)
	center := s.renderItems(s.Center)
	right := s.renderItems(s.Right)

	leftLen := lipgloss.Width(left)
	centerLen := lipgloss.Width(center)
	rightLen := lipgloss.Width(right)

	// Calculate padding
	totalContent := leftLen + centerLen + rightLen
	remaining := s.Width - totalContent

	if remaining < 0 {
		// Content too wide, just join
		return s.style.Container.Width(s.Width).Render(left + center + right)
	}

	// Distribute space
	leftPad := remaining / 2
	rightPad := remaining - leftPad

	content := left +
		strings.Repeat(" ", leftPad) +
		center +
		strings.Repeat(" ", rightPad) +
		right

	return s.style.Container.Width(s.Width).Render(content)
}

// KeyHint provides keyboard shortcut hints
type KeyHint struct {
	Key   string
	Label string
}

// KeyHints displays keyboard shortcuts
type KeyHints struct {
	ID     string
	Width  int
	Hints  []KeyHint
	style  KeyHintsStyle
}

// KeyHintsStyle holds key hints styling
type KeyHintsStyle struct {
	Container lipgloss.Style
	Key       lipgloss.Style
	Label     lipgloss.Style
	Separator lipgloss.Style
}

// DefaultKeyHintsStyle returns default styling
func DefaultKeyHintsStyle() KeyHintsStyle {
	return KeyHintsStyle{
		Container: lipgloss.NewStyle().
			Background(lipgloss.Color("#0D0221")),
		Key: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E040FB")).
			Bold(true),
		Label: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888")),
		Separator: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#3C096C")),
	}
}

// NewKeyHints creates key hints display
func NewKeyHints(id string, width int) *KeyHints {
	return &KeyHints{
		ID:    id,
		Width: width,
		Hints: []KeyHint{},
		style: DefaultKeyHintsStyle(),
	}
}

// AddHint adds a key hint
func (k *KeyHints) AddHint(key, label string) *KeyHints {
	k.Hints = append(k.Hints, KeyHint{Key: key, Label: label})
	return k
}

// View renders key hints
func (k *KeyHints) View() string {
	var parts []string
	for _, hint := range k.Hints {
		keyStr := k.style.Key.Render(hint.Key)
		labelStr := k.style.Label.Render(hint.Label)
		parts = append(parts, keyStr+" "+labelStr)
	}

	sep := k.style.Separator.Render("  ")
	content := strings.Join(parts, sep)

	return k.style.Container.Width(k.Width).Render(content)
}

