// Package display - Table Update and View methods
package display

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/makeatui/makeatui/pkg/widgets/mouse"
)

// GetZone returns the mouse zone
func (t *Table) GetZone(x, y, width, height int) *mouse.Zone {
	return &mouse.Zone{
		ID:     t.ID,
		X:      x,
		Y:      y,
		Width:  width,
		Height: height,
	}
}

// Update handles messages
func (t *Table) Update(msg tea.Msg) (*Table, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			t.SelectedRow--
			if t.SelectedRow < 0 {
				t.SelectedRow = len(t.Rows) - 1
			}
			t.ensureVisible()
		case "down", "j":
			t.SelectedRow++
			if t.SelectedRow >= len(t.Rows) {
				t.SelectedRow = 0
			}
			t.ensureVisible()
		case "home", "g":
			t.SelectedRow = 0
			t.ensureVisible()
		case "end", "G":
			t.SelectedRow = len(t.Rows) - 1
			t.ensureVisible()
		case "pgup":
			t.SelectedRow -= t.Height
			if t.SelectedRow < 0 {
				t.SelectedRow = 0
			}
			t.ensureVisible()
		case "pgdown":
			t.SelectedRow += t.Height
			if t.SelectedRow >= len(t.Rows) {
				t.SelectedRow = len(t.Rows) - 1
			}
			t.ensureVisible()
		case " ":
			if t.MultiSelect && t.SelectedRow >= 0 && t.SelectedRow < len(t.Rows) {
				t.Rows[t.SelectedRow].Selected = !t.Rows[t.SelectedRow].Selected
			}
		case "a":
			if t.MultiSelect {
				t.SelectAll()
			}
		case "A":
			if t.MultiSelect {
				t.ClearSelection()
			}
		}

	case mouse.ScrollMsg:
		if msg.Direction == mouse.ScrollUp {
			t.scrollOffset--
			if t.scrollOffset < 0 {
				t.scrollOffset = 0
			}
		} else if msg.Direction == mouse.ScrollDown {
			maxOffset := len(t.Rows) - t.Height
			if maxOffset < 0 {
				maxOffset = 0
			}
			t.scrollOffset++
			if t.scrollOffset > maxOffset {
				t.scrollOffset = maxOffset
			}
		}
	}

	return t, nil
}

// View renders the table
func (t *Table) View() string {
	var parts []string

	// Header
	if t.ShowHeader {
		var headerCells []string
		for _, col := range t.Columns {
			title := col.Title
			if t.Sortable && t.sortColumn >= 0 {
				// Add sort indicator
				for i, c := range t.Columns {
					if c.Key == col.Key && i == t.sortColumn {
						if t.sortAsc {
							title += " ▲"
						} else {
							title += " ▼"
						}
					}
				}
			}
			cell := t.style.HeaderCell.Width(col.Width).Render(title)
			headerCells = append(headerCells, cell)
		}
		header := t.style.Header.Render(strings.Join(headerCells, ""))
		parts = append(parts, header)
	}

	// Rows
	visibleEnd := t.scrollOffset + t.Height
	if visibleEnd > len(t.Rows) {
		visibleEnd = len(t.Rows)
	}

	for i := t.scrollOffset; i < visibleEnd; i++ {
		row := t.Rows[i]
		var cells []string

		// Determine row style
		rowStyle := t.style.Row
		if t.Striped && i%2 == 1 {
			rowStyle = t.style.RowAlt
		}
		if i == t.SelectedRow && t.Focused {
			rowStyle = t.style.RowSelected
		}
		if row.Selected {
			rowStyle = t.style.RowSelected
		}

		for _, col := range t.Columns {
			content := row.Data[col.Key]
			cell := t.renderCell(content, col, t.style.Cell)
			cells = append(cells, cell)
		}

		rowView := rowStyle.Render(strings.Join(cells, ""))
		parts = append(parts, rowView)
	}

	// Scroll indicator
	if len(t.Rows) > t.Height {
		scrollInfo := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")).
			Italic(true).
			Render("  " + string(rune('0'+t.scrollOffset+1)) + "-" +
				string(rune('0'+visibleEnd)) + " of " +
				string(rune('0'+len(t.Rows))))
		parts = append(parts, scrollInfo)
	}

	result := lipgloss.JoinVertical(lipgloss.Left, parts...)

	if t.ShowBorder {
		result = t.style.Border.Render(result)
	}

	return result
}

