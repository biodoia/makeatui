// Package display provides display/output components
package display

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/makeatui/makeatui/pkg/widgets/mouse"
)

// TableColumn defines a table column
type TableColumn struct {
	Key      string
	Title    string
	Width    int
	MinWidth int
	MaxWidth int
	Align    lipgloss.Position
	Sortable bool
}

// TableRow represents a table row
type TableRow struct {
	Data     map[string]string
	Selected bool
	Disabled bool
}

// Table provides a data table component (inspired by Textual/Ratatui)
type Table struct {
	ID           string
	Columns      []TableColumn
	Rows         []TableRow
	SelectedRow  int
	SelectedCol  int
	Focused      bool
	Sortable     bool
	Selectable   bool
	MultiSelect  bool
	ShowHeader   bool
	ShowBorder   bool
	Striped      bool
	Height       int // visible rows
	scrollOffset int
	sortColumn   int
	sortAsc      bool
	style        TableStyle
	zoneManager  *mouse.ZoneManager
}

// TableStyle holds table styling
type TableStyle struct {
	Header       lipgloss.Style
	HeaderCell   lipgloss.Style
	Row          lipgloss.Style
	RowAlt       lipgloss.Style
	RowSelected  lipgloss.Style
	RowHover     lipgloss.Style
	Cell         lipgloss.Style
	CellSelected lipgloss.Style
	Border       lipgloss.Style
}

// DefaultTableStyle returns default styling
func DefaultTableStyle() TableStyle {
	return TableStyle{
		Header: lipgloss.NewStyle().
			Background(lipgloss.Color("#3C096C")).
			Foreground(lipgloss.Color("#E040FB")).
			Bold(true),
		HeaderCell: lipgloss.NewStyle().
			Padding(0, 1),
		Row: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")),
		RowAlt: lipgloss.NewStyle().
			Background(lipgloss.Color("#1A0533")).
			Foreground(lipgloss.Color("#FFFFFF")),
		RowSelected: lipgloss.NewStyle().
			Background(lipgloss.Color("#9D4EDD")).
			Foreground(lipgloss.Color("#FFFFFF")).
			Bold(true),
		RowHover: lipgloss.NewStyle().
			Background(lipgloss.Color("#5A189A")).
			Foreground(lipgloss.Color("#FFFFFF")),
		Cell: lipgloss.NewStyle().
			Padding(0, 1),
		CellSelected: lipgloss.NewStyle().
			Padding(0, 1).
			Background(lipgloss.Color("#7B2CBF")),
		Border: lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7B2CBF")),
	}
}

// NewTable creates a new table
func NewTable(id string) *Table {
	return &Table{
		ID:          id,
		Columns:     []TableColumn{},
		Rows:        []TableRow{},
		SelectedRow: 0,
		ShowHeader:  true,
		ShowBorder:  true,
		Striped:     true,
		Selectable:  true,
		Height:      10,
		sortColumn:  -1,
		style:       DefaultTableStyle(),
		zoneManager: mouse.NewZoneManager(),
	}
}

// AddColumn adds a column
func (t *Table) AddColumn(key, title string, width int) *Table {
	t.Columns = append(t.Columns, TableColumn{
		Key:   key,
		Title: title,
		Width: width,
		Align: lipgloss.Left,
	})
	return t
}

// AddRow adds a row
func (t *Table) AddRow(data map[string]string) *Table {
	t.Rows = append(t.Rows, TableRow{Data: data})
	return t
}

// SetRows sets all rows
func (t *Table) SetRows(rows []TableRow) *Table {
	t.Rows = rows
	return t
}

// SetHeight sets visible rows
func (t *Table) SetHeight(height int) *Table {
	t.Height = height
	return t
}

// SetStriped enables striped rows
func (t *Table) SetStriped(striped bool) *Table {
	t.Striped = striped
	return t
}

// GetSelectedRow returns the selected row data
func (t *Table) GetSelectedRow() *TableRow {
	if t.SelectedRow >= 0 && t.SelectedRow < len(t.Rows) {
		return &t.Rows[t.SelectedRow]
	}
	return nil
}

// GetSelectedRows returns all selected rows (multi-select)
func (t *Table) GetSelectedRows() []TableRow {
	var selected []TableRow
	for _, row := range t.Rows {
		if row.Selected {
			selected = append(selected, row)
		}
	}
	return selected
}

// ClearSelection clears all selections
func (t *Table) ClearSelection() *Table {
	for i := range t.Rows {
		t.Rows[i].Selected = false
	}
	return t
}

// SelectAll selects all rows
func (t *Table) SelectAll() *Table {
	for i := range t.Rows {
		t.Rows[i].Selected = true
	}
	return t
}

// ensureVisible scrolls to keep selected row visible
func (t *Table) ensureVisible() {
	if t.SelectedRow < t.scrollOffset {
		t.scrollOffset = t.SelectedRow
	}
	if t.SelectedRow >= t.scrollOffset+t.Height {
		t.scrollOffset = t.SelectedRow - t.Height + 1
	}
}

// renderCell renders a single cell
func (t *Table) renderCell(content string, col TableColumn, style lipgloss.Style) string {
	cellStyle := style.Width(col.Width).Align(col.Align)
	// Truncate if needed
	if len(content) > col.Width {
		content = content[:col.Width-1] + "…"
	}
	return cellStyle.Render(content)
}

// TotalWidth returns total table width
func (t *Table) TotalWidth() int {
	width := 0
	for _, col := range t.Columns {
		width += col.Width + 2 // +2 for padding
	}
	return width
}

// RowCount returns number of rows
func (t *Table) RowCount() int {
	return len(t.Rows)
}

