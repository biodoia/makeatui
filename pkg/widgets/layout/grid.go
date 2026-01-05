// Package layout provides layout components for TUI design
// Inspired by Ratatui, Textual, and CSS Grid/Flexbox
package layout

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Grid implements a CSS Grid-like layout system
type Grid struct {
	Columns     int
	Rows        int
	Gap         int
	ColumnGap   int
	RowGap      int
	cells       [][]string
	columnSizes []int
	rowSizes    []int
	style       lipgloss.Style
}

// NewGrid creates a new grid layout
func NewGrid(columns, rows int) *Grid {
	cells := make([][]string, rows)
	for i := range cells {
		cells[i] = make([]string, columns)
	}

	return &Grid{
		Columns:     columns,
		Rows:        rows,
		Gap:         1,
		cells:       cells,
		columnSizes: make([]int, columns),
		rowSizes:    make([]int, rows),
		style:       lipgloss.NewStyle(),
	}
}

// SetCell sets content at a specific grid position
func (g *Grid) SetCell(col, row int, content string) *Grid {
	if row >= 0 && row < g.Rows && col >= 0 && col < g.Columns {
		g.cells[row][col] = content
	}
	return g
}

// SetColumnSize sets the width of a column
func (g *Grid) SetColumnSize(col, size int) *Grid {
	if col >= 0 && col < g.Columns {
		g.columnSizes[col] = size
	}
	return g
}

// SetRowSize sets the height of a row
func (g *Grid) SetRowSize(row, size int) *Grid {
	if row >= 0 && row < g.Rows {
		g.rowSizes[row] = size
	}
	return g
}

// SetGap sets uniform gap
func (g *Grid) SetGap(gap int) *Grid {
	g.Gap = gap
	g.ColumnGap = gap
	g.RowGap = gap
	return g
}

// SetStyle sets the grid style
func (g *Grid) SetStyle(style lipgloss.Style) *Grid {
	g.style = style
	return g
}

// Render renders the grid
func (g *Grid) Render() string {
	if g.ColumnGap == 0 {
		g.ColumnGap = g.Gap
	}
	if g.RowGap == 0 {
		g.RowGap = g.Gap
	}

	// Calculate column widths if not set
	for col := 0; col < g.Columns; col++ {
		if g.columnSizes[col] == 0 {
			maxWidth := 0
			for row := 0; row < g.Rows; row++ {
				w := lipgloss.Width(g.cells[row][col])
				if w > maxWidth {
					maxWidth = w
				}
			}
			g.columnSizes[col] = maxWidth
		}
	}

	// Render rows
	var rows []string
	for row := 0; row < g.Rows; row++ {
		var cells []string
		for col := 0; col < g.Columns; col++ {
			cell := g.cells[row][col]
			cellStyle := lipgloss.NewStyle().Width(g.columnSizes[col])
			cells = append(cells, cellStyle.Render(cell))
		}
		rowStr := lipgloss.JoinHorizontal(lipgloss.Top, interleave(cells, strings.Repeat(" ", g.ColumnGap))...)
		rows = append(rows, rowStr)
	}

	// Join rows with gap
	result := strings.Join(rows, strings.Repeat("\n", g.RowGap+1))
	return g.style.Render(result)
}

// interleave inserts separator between elements
func interleave(items []string, sep string) []string {
	if len(items) == 0 {
		return items
	}
	result := make([]string, len(items)*2-1)
	for i, item := range items {
		result[i*2] = item
		if i < len(items)-1 {
			result[i*2+1] = sep
		}
	}
	return result
}

// AutoGrid creates a grid that auto-distributes items
func AutoGrid(items []string, columns, width int) string {
	if columns <= 0 {
		columns = 1
	}

	rows := (len(items) + columns - 1) / columns
	grid := NewGrid(columns, rows)

	colWidth := width / columns
	for i := range columns {
		grid.SetColumnSize(i, colWidth)
	}

	for i, item := range items {
		col := i % columns
		row := i / columns
		grid.SetCell(col, row, item)
	}

	return grid.Render()
}

