// Package feedback provides feedback components (Toast, Modal, Notification)
package feedback

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ToastType defines toast severity
type ToastType int

const (
	ToastInfo ToastType = iota
	ToastSuccess
	ToastWarning
	ToastError
)

// Toast represents a temporary notification
type Toast struct {
	ID        string
	Message   string
	Type      ToastType
	Duration  time.Duration
	Visible   bool
	Icon      string
	CreatedAt time.Time
	style     ToastStyle
}

// ToastStyle holds toast styling
type ToastStyle struct {
	Info    lipgloss.Style
	Success lipgloss.Style
	Warning lipgloss.Style
	Error   lipgloss.Style
}

// DefaultToastStyle returns default styling
func DefaultToastStyle() ToastStyle {
	base := lipgloss.NewStyle().
		Padding(0, 2).
		Border(lipgloss.RoundedBorder()).
		Bold(true)

	return ToastStyle{
		Info: base.Copy().
			BorderForeground(lipgloss.Color("#4CC9F0")).
			Foreground(lipgloss.Color("#4CC9F0")),
		Success: base.Copy().
			BorderForeground(lipgloss.Color("#6BCB77")).
			Foreground(lipgloss.Color("#6BCB77")),
		Warning: base.Copy().
			BorderForeground(lipgloss.Color("#FFD93D")).
			Foreground(lipgloss.Color("#FFD93D")),
		Error: base.Copy().
			BorderForeground(lipgloss.Color("#FF6B6B")).
			Foreground(lipgloss.Color("#FF6B6B")),
	}
}

// ToastIcons for different types
var ToastIcons = map[ToastType]string{
	ToastInfo:    "ℹ",
	ToastSuccess: "✓",
	ToastWarning: "⚠",
	ToastError:   "✗",
}

// NewToast creates a new toast
func NewToast(id, message string, toastType ToastType) *Toast {
	return &Toast{
		ID:        id,
		Message:   message,
		Type:      toastType,
		Duration:  3 * time.Second,
		Visible:   true,
		Icon:      ToastIcons[toastType],
		CreatedAt: time.Now(),
		style:     DefaultToastStyle(),
	}
}

// SetDuration sets toast duration
func (t *Toast) SetDuration(d time.Duration) *Toast {
	t.Duration = d
	return t
}

// SetIcon sets custom icon
func (t *Toast) SetIcon(icon string) *Toast {
	t.Icon = icon
	return t
}

// Show makes the toast visible
func (t *Toast) Show() *Toast {
	t.Visible = true
	t.CreatedAt = time.Now()
	return t
}

// Hide hides the toast
func (t *Toast) Hide() *Toast {
	t.Visible = false
	return t
}

// IsExpired checks if toast should be hidden
func (t *Toast) IsExpired() bool {
	return time.Since(t.CreatedAt) > t.Duration
}

// View renders the toast
func (t *Toast) View() string {
	if !t.Visible {
		return ""
	}

	var style lipgloss.Style
	switch t.Type {
	case ToastInfo:
		style = t.style.Info
	case ToastSuccess:
		style = t.style.Success
	case ToastWarning:
		style = t.style.Warning
	case ToastError:
		style = t.style.Error
	}

	content := t.Icon + " " + t.Message
	return style.Render(content)
}

// ToastManager manages multiple toasts
type ToastManager struct {
	Toasts   []*Toast
	MaxShow  int
	Position string // "top", "bottom", "top-right", etc.
}

// NewToastManager creates a toast manager
func NewToastManager() *ToastManager {
	return &ToastManager{
		Toasts:   []*Toast{},
		MaxShow:  5,
		Position: "top-right",
	}
}

// Add adds a toast
func (tm *ToastManager) Add(toast *Toast) *ToastManager {
	tm.Toasts = append(tm.Toasts, toast)
	return tm
}

// Info adds an info toast
func (tm *ToastManager) Info(message string) *ToastManager {
	return tm.Add(NewToast("", message, ToastInfo))
}

// Success adds a success toast
func (tm *ToastManager) Success(message string) *ToastManager {
	return tm.Add(NewToast("", message, ToastSuccess))
}

// Warning adds a warning toast
func (tm *ToastManager) Warning(message string) *ToastManager {
	return tm.Add(NewToast("", message, ToastWarning))
}

// Error adds an error toast
func (tm *ToastManager) Error(message string) *ToastManager {
	return tm.Add(NewToast("", message, ToastError))
}

// Cleanup removes expired toasts
func (tm *ToastManager) Cleanup() {
	var active []*Toast
	for _, t := range tm.Toasts {
		if !t.IsExpired() && t.Visible {
			active = append(active, t)
		}
	}
	tm.Toasts = active
}

// View renders all visible toasts
func (tm *ToastManager) View() string {
	tm.Cleanup()

	var views []string
	count := 0
	for _, t := range tm.Toasts {
		if count >= tm.MaxShow {
			break
		}
		if t.Visible {
			views = append(views, t.View())
			count++
		}
	}

	return lipgloss.JoinVertical(lipgloss.Right, views...)
}

// TickCmd returns a command to auto-hide expired toasts
func (tm *ToastManager) TickCmd() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return ToastTickMsg{}
	})
}

// ToastTickMsg is sent for toast timing
type ToastTickMsg struct{}

