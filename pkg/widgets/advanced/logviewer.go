// Package advanced - Log viewer widget
package advanced

import (
	"regexp"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/makeatui/makeatui/pkg/widgets/mouse"
)

// LogLevel represents log severity
type LogLevel int

const (
	LogTrace LogLevel = iota
	LogDebug
	LogInfo
	LogWarn
	LogError
	LogFatal
)

// LogEntry represents a log entry
type LogEntry struct {
	Timestamp time.Time
	Level     LogLevel
	Message   string
	Source    string
}

// LogViewer provides a log viewing widget
type LogViewer struct {
	ID          string
	Width       int
	Height      int
	Entries     []LogEntry
	Filter      LogLevel
	Search      string
	AutoScroll  bool
	ShowTime    bool
	ShowLevel   bool
	ShowSource  bool
	scrollY     int
	style       LogViewerStyle
	zoneManager *mouse.ZoneManager
}

// LogViewerStyle holds styling
type LogViewerStyle struct {
	Container lipgloss.Style
	Timestamp lipgloss.Style
	Source    lipgloss.Style
	Trace     lipgloss.Style
	Debug     lipgloss.Style
	Info      lipgloss.Style
	Warn      lipgloss.Style
	Error     lipgloss.Style
	Fatal     lipgloss.Style
	Highlight lipgloss.Style
}

// DefaultLogViewerStyle returns default styling
func DefaultLogViewerStyle() LogViewerStyle {
	return LogViewerStyle{
		Container: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7B2CBF")),
		Timestamp: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")),
		Source: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9D4EDD")),
		Trace: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")),
		Debug: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#4CC9F0")),
		Info: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6BCB77")),
		Warn: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFD93D")),
		Error: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF6B6B")),
		Fatal: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Bold(true),
		Highlight: lipgloss.NewStyle().
			Background(lipgloss.Color("#FFD93D")).
			Foreground(lipgloss.Color("#000000")),
	}
}

// LogLevelLabels maps levels to strings
var LogLevelLabels = map[LogLevel]string{
	LogTrace: "TRACE",
	LogDebug: "DEBUG",
	LogInfo:  "INFO",
	LogWarn:  "WARN",
	LogError: "ERROR",
	LogFatal: "FATAL",
}

// NewLogViewer creates a log viewer
func NewLogViewer(id string, width, height int) *LogViewer {
	return &LogViewer{
		ID:          id,
		Width:       width,
		Height:      height,
		Entries:     []LogEntry{},
		Filter:      LogTrace,
		AutoScroll:  true,
		ShowTime:    true,
		ShowLevel:   true,
		ShowSource:  false,
		style:       DefaultLogViewerStyle(),
		zoneManager: mouse.NewZoneManager(),
	}
}

// AddEntry adds a log entry
func (lv *LogViewer) AddEntry(level LogLevel, message, source string) *LogViewer {
	lv.Entries = append(lv.Entries, LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   message,
		Source:    source,
	})
	if lv.AutoScroll {
		lv.scrollToBottom()
	}
	return lv
}

// Log convenience methods
func (lv *LogViewer) Trace(msg string) *LogViewer { return lv.AddEntry(LogTrace, msg, "") }
func (lv *LogViewer) Debug(msg string) *LogViewer { return lv.AddEntry(LogDebug, msg, "") }
func (lv *LogViewer) Info(msg string) *LogViewer  { return lv.AddEntry(LogInfo, msg, "") }
func (lv *LogViewer) Warn(msg string) *LogViewer  { return lv.AddEntry(LogWarn, msg, "") }
func (lv *LogViewer) Error(msg string) *LogViewer { return lv.AddEntry(LogError, msg, "") }
func (lv *LogViewer) Fatal(msg string) *LogViewer { return lv.AddEntry(LogFatal, msg, "") }

// SetFilter sets minimum log level
func (lv *LogViewer) SetFilter(level LogLevel) *LogViewer {
	lv.Filter = level
	return lv
}

// SetSearch sets search query
func (lv *LogViewer) SetSearch(query string) *LogViewer {
	lv.Search = query
	return lv
}

// Clear clears all entries
func (lv *LogViewer) Clear() *LogViewer {
	lv.Entries = []LogEntry{}
	lv.scrollY = 0
	return lv
}

// scrollToBottom scrolls to latest entry
func (lv *LogViewer) scrollToBottom() {
	filtered := lv.getFilteredEntries()
	if len(filtered) > lv.Height-2 {
		lv.scrollY = len(filtered) - (lv.Height - 2)
	}
}

// getFilteredEntries returns entries matching filter and search
func (lv *LogViewer) getFilteredEntries() []LogEntry {
	var filtered []LogEntry
	for _, entry := range lv.Entries {
		if entry.Level < lv.Filter {
			continue
		}
		if lv.Search != "" {
			if !strings.Contains(strings.ToLower(entry.Message), strings.ToLower(lv.Search)) {
				continue
			}
		}
		filtered = append(filtered, entry)
	}
	return filtered
}

// getLevelStyle returns style for log level
func (lv *LogViewer) getLevelStyle(level LogLevel) lipgloss.Style {
	switch level {
	case LogTrace:
		return lv.style.Trace
	case LogDebug:
		return lv.style.Debug
	case LogInfo:
		return lv.style.Info
	case LogWarn:
		return lv.style.Warn
	case LogError:
		return lv.style.Error
	case LogFatal:
		return lv.style.Fatal
	}
	return lv.style.Info
}

// highlightSearch highlights search matches in text
func (lv *LogViewer) highlightSearch(text string) string {
	if lv.Search == "" {
		return text
	}
	re := regexp.MustCompile("(?i)" + regexp.QuoteMeta(lv.Search))
	return re.ReplaceAllStringFunc(text, func(match string) string {
		return lv.style.Highlight.Render(match)
	})
}

// Update handles messages
func (lv *LogViewer) Update(msg tea.Msg) (*LogViewer, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			lv.scrollY--
			if lv.scrollY < 0 {
				lv.scrollY = 0
			}
			lv.AutoScroll = false
		case "down", "j":
			filtered := lv.getFilteredEntries()
			maxScroll := len(filtered) - (lv.Height - 2)
			if maxScroll < 0 {
				maxScroll = 0
			}
			lv.scrollY++
			if lv.scrollY > maxScroll {
				lv.scrollY = maxScroll
			}
		case "G":
			lv.scrollToBottom()
			lv.AutoScroll = true
		case "g":
			lv.scrollY = 0
			lv.AutoScroll = false
		case "c":
			lv.Clear()
		}

	case mouse.ScrollMsg:
		if msg.Direction == mouse.ScrollUp {
			lv.scrollY--
			if lv.scrollY < 0 {
				lv.scrollY = 0
			}
			lv.AutoScroll = false
		} else {
			filtered := lv.getFilteredEntries()
			maxScroll := len(filtered) - (lv.Height - 2)
			if maxScroll < 0 {
				maxScroll = 0
			}
			lv.scrollY++
			if lv.scrollY > maxScroll {
				lv.scrollY = maxScroll
			}
		}
	}

	return lv, nil
}

