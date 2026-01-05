// Package advanced - File Manager Update method
package advanced

import (
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/makeatui/makeatui/pkg/widgets/mouse"
)

// Update handles messages
func (fm *FileManager) Update(msg tea.Msg) (*FileManager, tea.Cmd) {
	if !fm.Focused {
		return fm, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			fm.moveCursor(-1)
		case "down", "j":
			fm.moveCursor(1)
		case "home", "g":
			fm.Cursor = 0
			fm.ensureVisible()
		case "end", "G":
			fm.Cursor = len(fm.Entries) - 1
			fm.ensureVisible()
		case "pgup":
			fm.moveCursor(-(fm.Height - 4))
		case "pgdown":
			fm.moveCursor(fm.Height - 4)
		case "enter", "l", "right":
			if fm.Cursor >= 0 && fm.Cursor < len(fm.Entries) {
				entry := fm.Entries[fm.Cursor]
				if entry.IsDir {
					_ = fm.SetDirectory(entry.Path)
				} else {
					return fm, func() tea.Msg {
						return FileOpenMsg{Path: entry.Path}
					}
				}
			}
		case "h", "left", "backspace":
			fm.GoUp()
		case " ":
			fm.toggleSelect()
		case "a":
			fm.selectAll()
		case "A":
			fm.deselectAll()
		case ".":
			fm.ShowHidden = !fm.ShowHidden
			_ = fm.Refresh()
		case "d":
			fm.ShowDetails = !fm.ShowDetails
		case "s":
			fm.cycleSortBy()
		case "r":
			fm.SortReverse = !fm.SortReverse
			_ = fm.Refresh()
		case "R":
			_ = fm.Refresh()
		case "y":
			fm.Copy()
		case "x":
			fm.Cut()
		case "p":
			return fm, fm.Paste()
		case "D":
			return fm, fm.Delete()
		case "n":
			return fm, func() tea.Msg {
				return FileCreateMsg{Dir: fm.CurrentDir}
			}
		case "~":
			home, _ := os.UserHomeDir()
			_ = fm.SetDirectory(home)
		case "/":
			_ = fm.SetDirectory("/")
		}

	case tea.MouseMsg:
		if msg.Action == tea.MouseActionPress {
			clickY := msg.Y - 2 + fm.ScrollY
			if clickY >= 0 && clickY < len(fm.Entries) {
				if fm.Cursor == clickY {
					// Double click simulation
					entry := fm.Entries[fm.Cursor]
					if entry.IsDir {
						_ = fm.SetDirectory(entry.Path)
					}
				} else {
					fm.Cursor = clickY
				}
			}
		}

	case mouse.ScrollMsg:
		if msg.Direction == mouse.ScrollUp {
			fm.moveCursor(-3)
		} else {
			fm.moveCursor(3)
		}

	case FileRefreshMsg:
		_ = fm.Refresh()
	}

	return fm, nil
}

// moveCursor moves cursor by delta
func (fm *FileManager) moveCursor(delta int) {
	fm.Cursor += delta
	if fm.Cursor < 0 {
		fm.Cursor = 0
	}
	if fm.Cursor >= len(fm.Entries) {
		fm.Cursor = len(fm.Entries) - 1
	}
	fm.ensureVisible()
}

// ensureVisible scrolls to keep cursor visible
func (fm *FileManager) ensureVisible() {
	visibleHeight := fm.Height - 4
	if fm.Cursor < fm.ScrollY {
		fm.ScrollY = fm.Cursor
	}
	if fm.Cursor >= fm.ScrollY+visibleHeight {
		fm.ScrollY = fm.Cursor - visibleHeight + 1
	}
}

// GoUp navigates to parent directory
func (fm *FileManager) GoUp() {
	parent := filepath.Dir(fm.CurrentDir)
	if parent != fm.CurrentDir {
		_ = fm.SetDirectory(parent)
	}
}

// GoBack navigates back in history
func (fm *FileManager) GoBack() {
	if fm.HistoryIdx > 0 {
		fm.HistoryIdx--
		_ = fm.SetDirectory(fm.History[fm.HistoryIdx])
	}
}

// Refresh reloads current directory
func (fm *FileManager) Refresh() error {
	return fm.SetDirectory(fm.CurrentDir)
}

// SetDirectory changes to a directory
func (fm *FileManager) SetDirectory(path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	entries, err := fm.readDirectory(absPath)
	if err != nil {
		return err
	}

	// Add to history
	if fm.CurrentDir != "" && fm.CurrentDir != absPath {
		fm.History = append(fm.History[:fm.HistoryIdx+1], fm.CurrentDir)
		fm.HistoryIdx = len(fm.History) - 1
	}

	fm.CurrentDir = absPath
	fm.Entries = entries
	fm.Cursor = 0
	fm.ScrollY = 0

	return nil
}

// toggleSelect toggles selection of current item
func (fm *FileManager) toggleSelect() {
	if fm.Cursor >= 0 && fm.Cursor < len(fm.Entries) {
		entry := &fm.Entries[fm.Cursor]
		if entry.Name != ".." {
			entry.Selected = !entry.Selected
			fm.SelectedSet[entry.Path] = entry.Selected
			if !entry.Selected {
				delete(fm.SelectedSet, entry.Path)
			}
		}
	}
}

// selectAll selects all items
func (fm *FileManager) selectAll() {
	for i := range fm.Entries {
		if fm.Entries[i].Name != ".." {
			fm.Entries[i].Selected = true
			fm.SelectedSet[fm.Entries[i].Path] = true
		}
	}
}

// deselectAll deselects all items
func (fm *FileManager) deselectAll() {
	for i := range fm.Entries {
		fm.Entries[i].Selected = false
	}
	fm.SelectedSet = make(map[string]bool)
}

// cycleSortBy cycles through sort options
func (fm *FileManager) cycleSortBy() {
	fm.SortBy = (fm.SortBy + 1) % 4
	_ = fm.Refresh()
}

// Copy copies selected files to clipboard
func (fm *FileManager) Copy() {
	fm.Clipboard = fm.getSelectedPaths()
	fm.ClipboardOp = ClipboardCopy
}

// Cut cuts selected files to clipboard
func (fm *FileManager) Cut() {
	fm.Clipboard = fm.getSelectedPaths()
	fm.ClipboardOp = ClipboardCut
}

// getSelectedPaths returns paths of selected items
func (fm *FileManager) getSelectedPaths() []string {
	var paths []string
	for path := range fm.SelectedSet {
		paths = append(paths, path)
	}
	if len(paths) == 0 && fm.Cursor >= 0 && fm.Cursor < len(fm.Entries) {
		entry := fm.Entries[fm.Cursor]
		if entry.Name != ".." {
			paths = append(paths, entry.Path)
		}
	}
	return paths
}

// Paste pastes files from clipboard
func (fm *FileManager) Paste() tea.Cmd {
	if len(fm.Clipboard) == 0 {
		return nil
	}
	return func() tea.Msg {
		return FilePasteMsg{
			Sources: fm.Clipboard,
			Dest:    fm.CurrentDir,
			Op:      fm.ClipboardOp,
		}
	}
}

// Delete deletes selected files
func (fm *FileManager) Delete() tea.Cmd {
	paths := fm.getSelectedPaths()
	if len(paths) == 0 {
		return nil
	}
	return func() tea.Msg {
		return FileDeleteMsg{Paths: paths}
	}
}

// GetSelected returns currently selected entry
func (fm *FileManager) GetSelected() *FileEntry {
	if fm.Cursor >= 0 && fm.Cursor < len(fm.Entries) {
		return &fm.Entries[fm.Cursor]
	}
	return nil
}

// formatSize formats file size
func formatSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}

// Message types

// FileOpenMsg is sent when a file is opened
type FileOpenMsg struct {
	Path string
}

// FileCreateMsg is sent to create a new file
type FileCreateMsg struct {
	Dir string
}

// FilePasteMsg is sent to paste files
type FilePasteMsg struct {
	Sources []string
	Dest    string
	Op      ClipboardOp
}

// FileDeleteMsg is sent to delete files
type FileDeleteMsg struct {
	Paths []string
}

// FileRefreshMsg is sent to refresh file list
type FileRefreshMsg struct{}

