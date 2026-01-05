// Package advanced - File Manager View method
package advanced

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// View renders the file manager
func (fm *FileManager) View() string {
	var b strings.Builder

	// Header
	header := fm.style.Header.Width(fm.Width - 2).Render("📁 File Manager")
	b.WriteString(header)
	b.WriteString("\n")

	// Path bar
	pathBar := fm.style.PathBar.Width(fm.Width - 2).Render(fm.truncatePath(fm.CurrentDir))
	b.WriteString(pathBar)
	b.WriteString("\n")

	// Calculate visible area
	visibleHeight := fm.Height - 4 // header + path + status
	if visibleHeight < 1 {
		visibleHeight = 1
	}

	// Render entries
	for i := fm.ScrollY; i < fm.ScrollY+visibleHeight && i < len(fm.Entries); i++ {
		entry := fm.Entries[i]
		line := fm.renderEntry(entry, i == fm.Cursor)
		b.WriteString(line)
		b.WriteString("\n")
	}

	// Pad remaining space
	for i := len(fm.Entries) - fm.ScrollY; i < visibleHeight; i++ {
		b.WriteString(strings.Repeat(" ", fm.Width-2))
		b.WriteString("\n")
	}

	// Status bar
	statusBar := fm.renderStatusBar()
	b.WriteString(statusBar)

	// Apply container style
	return fm.style.Container.Width(fm.Width).Height(fm.Height).Render(b.String())
}

// renderEntry renders a single file entry
func (fm *FileManager) renderEntry(entry FileEntry, isCursor bool) string {
	width := fm.Width - 4

	// Build line content
	var line strings.Builder

	// Selection marker
	if entry.Selected {
		line.WriteString("● ")
	} else {
		line.WriteString("  ")
	}

	// Icon
	line.WriteString(entry.Icon)
	line.WriteString(" ")

	// Name
	name := entry.Name
	if entry.IsDir && entry.Name != ".." {
		name += "/"
	}

	// Details
	if fm.ShowDetails && entry.Name != ".." {
		// Calculate available space for name
		detailsWidth := 20
		nameWidth := width - 4 - detailsWidth
		if nameWidth < 10 {
			nameWidth = 10
		}

		// Truncate name if needed
		if len(name) > nameWidth {
			name = name[:nameWidth-3] + "..."
		}
		line.WriteString(fmt.Sprintf("%-*s", nameWidth, name))

		// Size
		if !entry.IsDir {
			line.WriteString(fmt.Sprintf("%10s", formatSize(entry.Size)))
		} else {
			line.WriteString(fmt.Sprintf("%10s", "<DIR>"))
		}

		// Date
		line.WriteString(" ")
		line.WriteString(entry.ModTime.Format("Jan 02 15:04"))
	} else {
		// Just name
		if len(name) > width-4 {
			name = name[:width-7] + "..."
		}
		line.WriteString(name)
	}

	// Pad to width
	content := line.String()
	if len(content) < width {
		content += strings.Repeat(" ", width-len(content))
	}

	// Apply style
	var style lipgloss.Style
	if isCursor {
		style = fm.style.ItemCursor
	} else if entry.Selected {
		style = fm.style.ItemSelected
	} else if entry.IsDir {
		style = fm.style.ItemDir
	} else if strings.HasPrefix(entry.Name, ".") {
		style = fm.style.ItemHidden
	} else if entry.Mode&0111 != 0 {
		style = fm.style.ItemExec
	} else {
		style = fm.style.Item
	}

	return style.Width(width).Render(content)
}

// renderStatusBar renders the status bar
func (fm *FileManager) renderStatusBar() string {
	// Count selected
	selectedCount := len(fm.SelectedSet)

	// Build status
	var status strings.Builder

	// Position
	status.WriteString(fmt.Sprintf("%d/%d", fm.Cursor+1, len(fm.Entries)))

	// Selected count
	if selectedCount > 0 {
		status.WriteString(fmt.Sprintf(" | %d selected", selectedCount))
	}

	// Sort info
	sortNames := []string{"Name", "Size", "Time", "Type"}
	sortDir := "↑"
	if fm.SortReverse {
		sortDir = "↓"
	}
	status.WriteString(fmt.Sprintf(" | Sort: %s%s", sortNames[fm.SortBy], sortDir))

	// Hidden files
	if fm.ShowHidden {
		status.WriteString(" | Hidden: ON")
	}

	// Clipboard
	if len(fm.Clipboard) > 0 {
		op := "Copy"
		if fm.ClipboardOp == ClipboardCut {
			op = "Cut"
		}
		status.WriteString(fmt.Sprintf(" | %s: %d", op, len(fm.Clipboard)))
	}

	return fm.style.StatusBar.Width(fm.Width - 2).Render(status.String())
}

// truncatePath truncates path to fit width
func (fm *FileManager) truncatePath(path string) string {
	maxWidth := fm.Width - 4
	if len(path) <= maxWidth {
		return path
	}

	// Show ... at beginning
	return "..." + path[len(path)-maxWidth+3:]
}

// Focus sets focus state
func (fm *FileManager) Focus() {
	fm.Focused = true
}

// Blur removes focus
func (fm *FileManager) Blur() {
	fm.Focused = false
}

// SetSize sets dimensions
func (fm *FileManager) SetSize(width, height int) {
	fm.Width = width
	fm.Height = height
}

// SetStyle sets custom style
func (fm *FileManager) SetStyle(style FileManagerStyle) {
	fm.style = style
}

