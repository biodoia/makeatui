// Package navigation - Stepper View and CommandPalette
package navigation

import (
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// View renders the stepper
func (s *Stepper) View() string {
	var parts []string

	for i, step := range s.Steps {
		// Step indicator
		stepNum := string(rune('1' + i))
		if step.Icon != "" {
			stepNum = step.Icon
		}

		stepStyle := s.style.StepPending
		labelStyle := s.style.Label

		if s.Completed[i] {
			stepStyle = s.style.StepDone
			stepNum = "✓"
		}
		if i == s.Current {
			stepStyle = s.style.StepActive
			labelStyle = s.style.LabelActive
		}

		stepView := stepStyle.Render(stepNum)
		labelView := labelStyle.Render(step.Label)

		if s.Vertical {
			// Vertical layout
			parts = append(parts, lipgloss.JoinHorizontal(lipgloss.Center,
				stepView, " ", labelView))

			if step.Description != "" {
				parts = append(parts, "    "+s.style.Desc.Render(step.Description))
			}

			if i < len(s.Steps)-1 {
				parts = append(parts, s.style.Connector.Render("    │"))
			}
		} else {
			// Horizontal layout
			parts = append(parts, stepView+" "+labelView)

			if i < len(s.Steps)-1 {
				parts = append(parts, s.style.Connector.Render(" ─── "))
			}
		}
	}

	if s.Vertical {
		return lipgloss.JoinVertical(lipgloss.Left, parts...)
	}
	return lipgloss.JoinHorizontal(lipgloss.Center, parts...)
}

// CommandPalette provides a command palette (inspired by VSCode/Textual)
type CommandPalette struct {
	ID          string
	Commands    []Command
	Query       string
	Filtered    []Command
	Selected    int
	Visible     bool
	MaxResults  int
	Placeholder string
	style       CommandPaletteStyle
}

// Command represents a palette command
type Command struct {
	ID          string
	Label       string
	Description string
	Category    string
	Shortcut    string
	Icon        string
	Score       float64
	Action      func() tea.Cmd
}

// CommandPaletteStyle holds palette styling
type CommandPaletteStyle struct {
	Container   lipgloss.Style
	Input       lipgloss.Style
	Item        lipgloss.Style
	ItemSel     lipgloss.Style
	Label       lipgloss.Style
	Description lipgloss.Style
	Category    lipgloss.Style
	Shortcut    lipgloss.Style
	Empty       lipgloss.Style
}

// DefaultCommandPaletteStyle returns default styling
func DefaultCommandPaletteStyle() CommandPaletteStyle {
	return CommandPaletteStyle{
		Container: lipgloss.NewStyle().
			Background(lipgloss.Color("#0D0221")).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#9D4EDD")).
			Padding(0, 1).
			Width(60),
		Input: lipgloss.NewStyle().
			Background(lipgloss.Color("#1A0533")).
			Foreground(lipgloss.Color("#FFFFFF")).
			Padding(0, 1).
			MarginBottom(1),
		Item: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Padding(0, 1),
		ItemSel: lipgloss.NewStyle().
			Background(lipgloss.Color("#9D4EDD")).
			Foreground(lipgloss.Color("#FFFFFF")).
			Bold(true).
			Padding(0, 1),
		Label: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")),
		Description: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")).
			Italic(true),
		Category: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7B2CBF")),
		Shortcut: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E040FB")),
		Empty: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")).
			Italic(true).
			Padding(1),
	}
}

// NewCommandPalette creates a command palette
func NewCommandPalette(id string) *CommandPalette {
	return &CommandPalette{
		ID:          id,
		Commands:    []Command{},
		Filtered:    []Command{},
		MaxResults:  10,
		Placeholder: "Type a command...",
		style:       DefaultCommandPaletteStyle(),
	}
}

// AddCommand adds a command
func (cp *CommandPalette) AddCommand(id, label, desc, category, shortcut, icon string, action func() tea.Cmd) *CommandPalette {
	cp.Commands = append(cp.Commands, Command{
		ID:          id,
		Label:       label,
		Description: desc,
		Category:    category,
		Shortcut:    shortcut,
		Icon:        icon,
		Action:      action,
	})
	return cp
}

// Show shows the palette
func (cp *CommandPalette) Show() *CommandPalette {
	cp.Visible = true
	cp.Query = ""
	cp.Selected = 0
	cp.filter()
	return cp
}

// Hide hides the palette
func (cp *CommandPalette) Hide() *CommandPalette {
	cp.Visible = false
	return cp
}

// Toggle toggles visibility
func (cp *CommandPalette) Toggle() *CommandPalette {
	if cp.Visible {
		return cp.Hide()
	}
	return cp.Show()
}

// filter filters commands based on query
func (cp *CommandPalette) filter() {
	if cp.Query == "" {
		cp.Filtered = cp.Commands
		return
	}

	query := strings.ToLower(cp.Query)
	cp.Filtered = []Command{}

	for _, cmd := range cp.Commands {
		score := fuzzyScore(strings.ToLower(cmd.Label), query)
		if score > 0 {
			cmd.Score = score
			cp.Filtered = append(cp.Filtered, cmd)
		}
	}

	// Sort by score
	sort.Slice(cp.Filtered, func(i, j int) bool {
		return cp.Filtered[i].Score > cp.Filtered[j].Score
	})

	// Limit results
	if len(cp.Filtered) > cp.MaxResults {
		cp.Filtered = cp.Filtered[:cp.MaxResults]
	}
}

// fuzzyScore calculates a simple fuzzy match score
func fuzzyScore(str, query string) float64 {
	if strings.Contains(str, query) {
		return 1.0 + float64(len(query))/float64(len(str))
	}

	// Check if all query chars exist in order
	qi := 0
	for _, c := range str {
		if qi < len(query) && byte(c) == query[qi] {
			qi++
		}
	}
	if qi == len(query) {
		return float64(len(query)) / float64(len(str))
	}

	return 0
}

// Execute executes the selected command
func (cp *CommandPalette) Execute() tea.Cmd {
	if cp.Selected >= 0 && cp.Selected < len(cp.Filtered) {
		cmd := cp.Filtered[cp.Selected]
		cp.Hide()
		if cmd.Action != nil {
			return cmd.Action()
		}
	}
	return nil
}

