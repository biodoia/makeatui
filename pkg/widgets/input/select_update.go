// Package input - Select Update and View methods
package input

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/makeatui/makeatui/pkg/widgets/mouse"
)

// Update handles messages
func (s *Select) Update(msg tea.Msg) (*Select, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", " ":
			if s.Open {
				s.SelectIndex(s.hovered)
			} else {
				s.Toggle()
			}
		case "esc":
			s.Open = false
			s.SearchQuery = ""
		case "up", "k":
			if s.Open {
				s.hovered--
				if s.hovered < 0 {
					s.hovered = len(s.Options) - 1
				}
				s.ensureVisible()
			}
		case "down", "j":
			if s.Open {
				s.hovered++
				if s.hovered >= len(s.Options) {
					s.hovered = 0
				}
				s.ensureVisible()
			}
		case "backspace":
			if s.Searchable && len(s.SearchQuery) > 0 {
				s.SearchQuery = s.SearchQuery[:len(s.SearchQuery)-1]
			}
		default:
			if s.Searchable && s.Open && len(msg.String()) == 1 {
				s.SearchQuery += msg.String()
			}
		}

	case tea.MouseMsg:
		return s.handleMouse(msg)
	}

	return s, nil
}

func (s *Select) handleMouse(msg tea.MouseMsg) (*Select, tea.Cmd) {
	cmd := s.zoneManager.HandleMouse(msg)
	return s, cmd
}

func (s *Select) ensureVisible() {
	if s.hovered < s.scrollOffset {
		s.scrollOffset = s.hovered
	}
	if s.hovered >= s.scrollOffset+s.MaxVisible {
		s.scrollOffset = s.hovered - s.MaxVisible + 1
	}
}

// GetZone returns the main zone for this select
func (s *Select) GetZone(x, y, width, height int) *mouse.Zone {
	return &mouse.Zone{
		ID:     s.ID,
		X:      x,
		Y:      y,
		Width:  width,
		Height: height,
		OnClick: func() tea.Cmd {
			s.Toggle()
			return nil
		},
	}
}

// View renders the select
func (s *Select) View() string {
	var parts []string

	// Label
	if s.Label != "" {
		parts = append(parts, s.style.Label.Render(s.Label))
	}

	// Trigger button
	triggerStyle := s.style.Trigger
	if s.Open {
		triggerStyle = s.style.TriggerOpen
	}

	icon := "▼"
	if s.Open {
		icon = "▲"
	}

	displayText := s.GetLabel()
	if displayText == "" {
		displayText = s.Placeholder
	}

	trigger := triggerStyle.Render(displayText + " " + icon)
	parts = append(parts, trigger)

	// Dropdown
	if s.Open {
		var optionViews []string

		// Search box
		if s.Searchable {
			searchBox := s.style.Search.Render("🔍 " + s.SearchQuery + "█")
			optionViews = append(optionViews, searchBox)
		}

		// Options
		filtered := s.filteredOptions()
		visibleStart := s.scrollOffset
		visibleEnd := visibleStart + s.MaxVisible
		if visibleEnd > len(filtered) {
			visibleEnd = len(filtered)
		}

		for i := visibleStart; i < visibleEnd; i++ {
			if i >= len(filtered) {
				break
			}
			optIdx := filtered[i]
			opt := s.Options[optIdx]

			var optStyle lipgloss.Style
			if optIdx == s.Selected {
				optStyle = s.style.OptionSel
			} else if i == s.hovered {
				optStyle = s.style.OptionHover
			} else {
				optStyle = s.style.Option
			}

			if opt.Disabled {
				optStyle = optStyle.Foreground(lipgloss.Color("#666666"))
			}

			icon := "  "
			if optIdx == s.Selected {
				icon = "✓ "
			}
			if opt.Icon != "" {
				icon = opt.Icon + " "
			}

			optionViews = append(optionViews, optStyle.Render(icon+opt.Label))
		}

		// Scroll indicators
		if s.scrollOffset > 0 {
			optionViews = append([]string{"  ↑ more..."}, optionViews...)
		}
		if visibleEnd < len(filtered) {
			optionViews = append(optionViews, "  ↓ more...")
		}

		dropdown := s.style.Dropdown.Render(strings.Join(optionViews, "\n"))
		parts = append(parts, dropdown)
	}

	return strings.Join(parts, "\n")
}

