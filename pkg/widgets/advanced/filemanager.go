// Package advanced - File Manager widget (inspired by superfile and lf)
package advanced

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/makeatui/makeatui/pkg/widgets/mouse"
)

// FileEntry represents a file or directory
type FileEntry struct {
	Name     string
	Path     string
	IsDir    bool
	Size     int64
	ModTime  time.Time
	Mode     os.FileMode
	Selected bool
	Icon     string
}

// FileManager provides file browsing functionality
type FileManager struct {
	ID           string
	Width        int
	Height       int
	CurrentDir   string
	Entries      []FileEntry
	Cursor       int
	ScrollY      int
	ShowHidden   bool
	ShowDetails  bool
	MultiSelect  bool
	SelectedSet  map[string]bool
	SortBy       FileSortBy
	SortReverse  bool
	Filter       string
	Focused      bool
	Clipboard    []string
	ClipboardOp  ClipboardOp
	History      []string
	HistoryIdx   int
	style        FileManagerStyle
	zoneManager  *mouse.ZoneManager
}

// FileSortBy defines sort options
type FileSortBy int

const (
	SortByName FileSortBy = iota
	SortBySize
	SortByTime
	SortByType
)

// ClipboardOp defines clipboard operation
type ClipboardOp int

const (
	ClipboardNone ClipboardOp = iota
	ClipboardCopy
	ClipboardCut
)

// FileManagerStyle holds styling
type FileManagerStyle struct {
	Container    lipgloss.Style
	Header       lipgloss.Style
	PathBar      lipgloss.Style
	Item         lipgloss.Style
	ItemSelected lipgloss.Style
	ItemCursor   lipgloss.Style
	ItemDir      lipgloss.Style
	ItemExec     lipgloss.Style
	ItemHidden   lipgloss.Style
	Details      lipgloss.Style
	StatusBar    lipgloss.Style
	Scrollbar    lipgloss.Style
}

// DefaultFileManagerStyle returns default styling
func DefaultFileManagerStyle() FileManagerStyle {
	return FileManagerStyle{
		Container: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7B2CBF")),
		Header: lipgloss.NewStyle().
			Background(lipgloss.Color("#1A0533")).
			Foreground(lipgloss.Color("#E040FB")).
			Bold(true).
			Padding(0, 1),
		PathBar: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9D4EDD")).
			Padding(0, 1),
		Item: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Padding(0, 1),
		ItemSelected: lipgloss.NewStyle().
			Background(lipgloss.Color("#5A189A")).
			Foreground(lipgloss.Color("#FFFFFF")).
			Padding(0, 1),
		ItemCursor: lipgloss.NewStyle().
			Background(lipgloss.Color("#9D4EDD")).
			Foreground(lipgloss.Color("#FFFFFF")).
			Bold(true).
			Padding(0, 1),
		ItemDir: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00BFFF")).
			Bold(true).
			Padding(0, 1),
		ItemExec: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6BCB77")).
			Padding(0, 1),
		ItemHidden: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")).
			Padding(0, 1),
		Details: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888")),
		StatusBar: lipgloss.NewStyle().
			Background(lipgloss.Color("#0D0221")).
			Foreground(lipgloss.Color("#888888")).
			Padding(0, 1),
		Scrollbar: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#5A189A")),
	}
}

// NewFileManager creates a file manager
func NewFileManager(id string, width, height int) *FileManager {
	fm := &FileManager{
		ID:          id,
		Width:       width,
		Height:      height,
		ShowHidden:  false,
		ShowDetails: true,
		SortBy:      SortByName,
		SelectedSet: make(map[string]bool),
		History:     []string{},
		style:       DefaultFileManagerStyle(),
		zoneManager: mouse.NewZoneManager(),
	}

	// Start in current directory
	dir, err := os.Getwd()
	if err != nil {
		dir = "/"
	}
	_ = fm.SetDirectory(dir)

	return fm
}

// readDirectory reads directory contents
func (fm *FileManager) readDirectory(path string) ([]FileEntry, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var files []FileEntry

	// Add parent directory entry
	if path != "/" {
		files = append(files, FileEntry{
			Name:  "..",
			Path:  filepath.Dir(path),
			IsDir: true,
			Icon:  "📁",
		})
	}

	for _, entry := range entries {
		// Skip hidden files if not showing
		if !fm.ShowHidden && strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		fe := FileEntry{
			Name:    entry.Name(),
			Path:    filepath.Join(path, entry.Name()),
			IsDir:   entry.IsDir(),
			Size:    info.Size(),
			ModTime: info.ModTime(),
			Mode:    info.Mode(),
			Icon:    fm.getIcon(entry.Name(), entry.IsDir(), info.Mode()),
		}

		files = append(files, fe)
	}

	// Sort entries
	fm.sortEntries(files)

	return files, nil
}

// getIcon returns icon for file type
func (fm *FileManager) getIcon(name string, isDir bool, mode os.FileMode) string {
	if isDir {
		// Special directories
		switch strings.ToLower(name) {
		case "documents", "documenti":
			return "📄"
		case "downloads", "scaricati":
			return "📥"
		case "music", "musica":
			return "🎵"
		case "pictures", "immagini":
			return "🖼️"
		case "videos", "video":
			return "🎬"
		case ".git":
			return "🔧"
		case "node_modules":
			return "📦"
		default:
			return "📁"
		}
	}

	// Check if executable
	if mode&0111 != 0 {
		return "⚙️"
	}

	// By extension
	ext := strings.ToLower(filepath.Ext(name))
	switch ext {
	case ".go":
		return "🐹"
	case ".py":
		return "🐍"
	case ".js", ".ts":
		return "📜"
	case ".rs":
		return "🦀"
	case ".md":
		return "📝"
	case ".json", ".yaml", ".yml", ".toml":
		return "⚙️"
	case ".jpg", ".jpeg", ".png", ".gif", ".svg", ".webp":
		return "🖼️"
	case ".mp3", ".wav", ".flac", ".ogg":
		return "🎵"
	case ".mp4", ".mkv", ".avi", ".mov":
		return "🎬"
	case ".zip", ".tar", ".gz", ".rar", ".7z":
		return "📦"
	case ".pdf":
		return "📕"
	case ".html", ".css":
		return "🌐"
	case ".sh", ".bash":
		return "🖥️"
	case ".txt":
		return "📄"
	default:
		return "📄"
	}
}

// sortEntries sorts file entries
func (fm *FileManager) sortEntries(entries []FileEntry) {
	// Keep .. at top
	sort.SliceStable(entries, func(i, j int) bool {
		if entries[i].Name == ".." {
			return true
		}
		if entries[j].Name == ".." {
			return false
		}

		// Directories first
		if entries[i].IsDir != entries[j].IsDir {
			return entries[i].IsDir
		}

		var result bool
		switch fm.SortBy {
		case SortByName:
			result = strings.ToLower(entries[i].Name) < strings.ToLower(entries[j].Name)
		case SortBySize:
			result = entries[i].Size < entries[j].Size
		case SortByTime:
			result = entries[i].ModTime.After(entries[j].ModTime)
		case SortByType:
			extI := filepath.Ext(entries[i].Name)
			extJ := filepath.Ext(entries[j].Name)
			if extI == extJ {
				result = strings.ToLower(entries[i].Name) < strings.ToLower(entries[j].Name)
			} else {
				result = extI < extJ
			}
		default:
			result = strings.ToLower(entries[i].Name) < strings.ToLower(entries[j].Name)
		}

		if fm.SortReverse {
			return !result
		}
		return result
	})
}
