// Package advanced - LogViewer View method
package advanced

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// View renders the log viewer
func (lv *LogViewer) View() string {
	filtered := lv.getFilteredEntries()
	var lines []string

	// Calculate visible range
	visibleHeight := lv.Height - 2
	endIdx := lv.scrollY + visibleHeight
	if endIdx > len(filtered) {
		endIdx = len(filtered)
	}

	for i := lv.scrollY; i < endIdx; i++ {
		entry := filtered[i]
		var parts []string

		// Timestamp
		if lv.ShowTime {
			timeStr := entry.Timestamp.Format("15:04:05")
			parts = append(parts, lv.style.Timestamp.Render(timeStr))
		}

		// Level
		if lv.ShowLevel {
			levelStyle := lv.getLevelStyle(entry.Level)
			levelStr := LogLevelLabels[entry.Level]
			parts = append(parts, levelStyle.Width(5).Render(levelStr))
		}

		// Source
		if lv.ShowSource && entry.Source != "" {
			parts = append(parts, lv.style.Source.Render("["+entry.Source+"]"))
		}

		// Message
		message := lv.highlightSearch(entry.Message)
		parts = append(parts, message)

		line := strings.Join(parts, " ")
		if len(line) > lv.Width-4 {
			line = line[:lv.Width-7] + "..."
		}
		lines = append(lines, line)
	}

	// Pad if needed
	for len(lines) < visibleHeight {
		lines = append(lines, "")
	}

	// Status bar
	statusStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#666666")).
		Italic(true)

	autoScrollIndicator := ""
	if lv.AutoScroll {
		autoScrollIndicator = " [AUTO]"
	}

	status := statusStyle.Render(
		"Lines: " + string(rune('0'+len(lv.Entries)%10)) +
			" | Filter: " + LogLevelLabels[lv.Filter] +
			autoScrollIndicator,
	)
	lines = append(lines, status)

	content := lipgloss.JoinVertical(lipgloss.Left, lines...)
	return lv.style.Container.Width(lv.Width).Render(content)
}

// FileExplorer provides a file explorer widget
type FileExplorer struct {
	ID         string
	Width      int
	Height     int
	RootPath   string
	CurrentDir string
	Entries    []FileEntry
	Selected   int
	ShowHidden bool
	scrollY    int
	style      FileExplorerStyle
}

// FileEntry represents a file or directory
type FileEntry struct {
	Name    string
	Path    string
	IsDir   bool
	Size    int64
	ModTime string
	Icon    string
}

// FileExplorerStyle holds styling
type FileExplorerStyle struct {
	Container lipgloss.Style
	Dir       lipgloss.Style
	File      lipgloss.Style
	Selected  lipgloss.Style
	Path      lipgloss.Style
	Size      lipgloss.Style
}

// DefaultFileExplorerStyle returns default styling
func DefaultFileExplorerStyle() FileExplorerStyle {
	return FileExplorerStyle{
		Container: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7B2CBF")),
		Dir: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E040FB")).
			Bold(true),
		File: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")),
		Selected: lipgloss.NewStyle().
			Background(lipgloss.Color("#9D4EDD")).
			Foreground(lipgloss.Color("#FFFFFF")).
			Bold(true),
		Path: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7B2CBF")).
			Bold(true),
		Size: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")),
	}
}

// NewFileExplorer creates a file explorer
func NewFileExplorer(id string, width, height int) *FileExplorer {
	return &FileExplorer{
		ID:       id,
		Width:    width,
		Height:   height,
		RootPath: ".",
		Entries:  []FileEntry{},
		style:    DefaultFileExplorerStyle(),
	}
}

// SetPath sets the current directory
func (fe *FileExplorer) SetPath(path string) *FileExplorer {
	fe.CurrentDir = path
	return fe
}

// AddEntry adds a file entry (for demo/mock purposes)
func (fe *FileExplorer) AddEntry(name string, isDir bool, size int64, icon string) *FileExplorer {
	fe.Entries = append(fe.Entries, FileEntry{
		Name:  name,
		IsDir: isDir,
		Size:  size,
		Icon:  icon,
	})
	return fe
}

// View renders the file explorer
func (fe *FileExplorer) View() string {
	var lines []string

	// Path header
	pathLine := fe.style.Path.Render("📁 " + fe.CurrentDir)
	lines = append(lines, pathLine)
	lines = append(lines, "")

	// Entries
	visibleHeight := fe.Height - 4
	endIdx := fe.scrollY + visibleHeight
	if endIdx > len(fe.Entries) {
		endIdx = len(fe.Entries)
	}

	for i := fe.scrollY; i < endIdx; i++ {
		entry := fe.Entries[i]
		style := fe.style.File
		if entry.IsDir {
			style = fe.style.Dir
		}
		if i == fe.Selected {
			style = fe.style.Selected
		}

		icon := entry.Icon
		if icon == "" {
			if entry.IsDir {
				icon = "📁"
			} else {
				icon = "📄"
			}
		}

		line := style.Render(icon + " " + entry.Name)
		lines = append(lines, line)
	}

	// Pad
	for len(lines) < fe.Height-2 {
		lines = append(lines, "")
	}

	content := lipgloss.JoinVertical(lipgloss.Left, lines...)
	return fe.style.Container.Width(fe.Width).Height(fe.Height).Render(content)
}

