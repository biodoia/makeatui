// Package app provides the main application model
package app

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/makeatui/makeatui/internal/ui/canvas"
	"github.com/makeatui/makeatui/internal/ui/styles"
	"github.com/makeatui/makeatui/pkg/schema"
)

// FocusArea represents which area of the UI is focused
type FocusArea int

const (
	FocusSidebar FocusArea = iota
	FocusCanvas
	FocusProperties
)

// Model is the main application model
type Model struct {
	width      int
	height     int
	theme      styles.Theme
	styles     styles.Styles
	canvas     *canvas.Canvas
	focus      FocusArea
	components []ComponentItem
	selected   int
	showHelp   bool
	quitting   bool
	projectName string
}

// ComponentItem represents a component in the sidebar
type ComponentItem struct {
	Type schema.ComponentType
	Name string
	Icon string
}

// New creates a new application model
func New() Model {
	theme := styles.Ultraviolet
	s := styles.NewStyles(theme)

	componentList := []ComponentItem{
		{Type: schema.TypeBox, Name: "Box", Icon: "□"},
		{Type: schema.TypeText, Name: "Text", Icon: "T"},
		{Type: schema.TypeButton, Name: "Button", Icon: "◉"},
		{Type: schema.TypeInput, Name: "Input", Icon: "▭"},
		{Type: schema.TypeList, Name: "List", Icon: "☰"},
		{Type: schema.TypeTable, Name: "Table", Icon: "▦"},
		{Type: schema.TypeProgress, Name: "Progress", Icon: "█"},
		{Type: schema.TypeSpinner, Name: "Spinner", Icon: "◌"},
		{Type: schema.TypeViewport, Name: "Viewport", Icon: "◱"},
		{Type: schema.TypeTabs, Name: "Tabs", Icon: "⊟"},
	}

	return Model{
		theme:       theme,
		styles:      s,
		canvas:      canvas.New(60, 20, theme),
		focus:       FocusSidebar,
		components:  componentList,
		selected:    0,
		projectName: "Untitled Project",
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return nil
}

// KeyMap defines keyboard shortcuts
type KeyMap struct {
	Up       key.Binding
	Down     key.Binding
	Left     key.Binding
	Right    key.Binding
	Enter    key.Binding
	Tab      key.Binding
	Delete   key.Binding
	Help     key.Binding
	Quit     key.Binding
	Export   key.Binding
	MoveMod  key.Binding
}

var keys = KeyMap{
	Up:       key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "up")),
	Down:     key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "down")),
	Left:     key.NewBinding(key.WithKeys("left", "h"), key.WithHelp("←/h", "left")),
	Right:    key.NewBinding(key.WithKeys("right", "l"), key.WithHelp("→/l", "right")),
	Enter:    key.NewBinding(key.WithKeys("enter", " "), key.WithHelp("enter/space", "add component")),
	Tab:      key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "switch focus")),
	Delete:   key.NewBinding(key.WithKeys("d", "delete"), key.WithHelp("d/del", "delete")),
	Help:     key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "help")),
	Quit:     key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
	Export:   key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "export")),
	MoveMod:  key.NewBinding(key.WithKeys("m"), key.WithHelp("m", "move mode")),
}

// Update handles messages
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Resize canvas
		canvasWidth := m.width - 30 - 4 // sidebar width + padding
		canvasHeight := m.height - 6     // toolbar + statusbar
		m.canvas = canvas.New(canvasWidth, canvasHeight, m.theme)

	case tea.KeyMsg:
		if key.Matches(msg, keys.Quit) {
			m.quitting = true
			return m, tea.Quit
		}

		if key.Matches(msg, keys.Help) {
			m.showHelp = !m.showHelp
			return m, nil
		}

		if key.Matches(msg, keys.Tab) {
			m.focus = (m.focus + 1) % 3
			return m, nil
		}

		switch m.focus {
		case FocusSidebar:
			return m.updateSidebar(msg)
		case FocusCanvas:
			return m.updateCanvas(msg)
		}
	}

	return m, nil
}

func (m Model) updateSidebar(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, keys.Up):
		if m.selected > 0 {
			m.selected--
		}
	case key.Matches(msg, keys.Down):
		if m.selected < len(m.components)-1 {
			m.selected++
		}
	case key.Matches(msg, keys.Enter):
		comp := m.components[m.selected]
		newComp := schema.NewComponent(comp.Type, comp.Name)
		newComp.Text = comp.Name
		m.canvas.AddComponent(newComp)
		m.focus = FocusCanvas
	}
	return m, nil
}

func (m Model) updateCanvas(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, keys.Up):
		if m.canvas.Mode == canvas.ModeMove && m.canvas.Selected >= 0 {
			m.canvas.MoveSelected(0, -1)
		} else {
			m.canvas.CursorY--
			if m.canvas.CursorY < 0 {
				m.canvas.CursorY = 0
			}
		}
	case key.Matches(msg, keys.Down):
		if m.canvas.Mode == canvas.ModeMove && m.canvas.Selected >= 0 {
			m.canvas.MoveSelected(0, 1)
		} else {
			m.canvas.CursorY++
			if m.canvas.CursorY >= m.canvas.Height {
				m.canvas.CursorY = m.canvas.Height - 1
			}
		}
	case key.Matches(msg, keys.Left):
		if m.canvas.Mode == canvas.ModeMove && m.canvas.Selected >= 0 {
			m.canvas.MoveSelected(-1, 0)
		} else {
			m.canvas.CursorX--
			if m.canvas.CursorX < 0 {
				m.canvas.CursorX = 0
			}
		}
	case key.Matches(msg, keys.Right):
		if m.canvas.Mode == canvas.ModeMove && m.canvas.Selected >= 0 {
			m.canvas.MoveSelected(1, 0)
		} else {
			m.canvas.CursorX++
			if m.canvas.CursorX >= m.canvas.Width {
				m.canvas.CursorX = m.canvas.Width - 1
			}
		}
	case key.Matches(msg, keys.Delete):
		m.canvas.RemoveSelected()
	case key.Matches(msg, keys.MoveMod):
		if m.canvas.Mode == canvas.ModeMove {
			m.canvas.Mode = canvas.ModeSelect
		} else {
			m.canvas.Mode = canvas.ModeMove
		}
	case key.Matches(msg, keys.Enter):
		// Select component at cursor position
		m.selectComponentAtCursor()
	}
	return m, nil
}

// selectComponentAtCursor selects the component at the current cursor position
func (m *Model) selectComponentAtCursor() {
	for i, comp := range m.canvas.Components {
		if m.canvas.CursorX >= comp.Position.X &&
			m.canvas.CursorX < comp.Position.X+comp.Size.Width &&
			m.canvas.CursorY >= comp.Position.Y &&
			m.canvas.CursorY < comp.Position.Y+comp.Size.Height {
			m.canvas.Selected = i
			return
		}
	}
	m.canvas.Selected = -1
}

// View renders the UI - continued in view.go
func (m Model) View() string {
	if m.quitting {
		return lipgloss.NewStyle().Foreground(m.theme.Primary).Render("\n  💜 Thanks for using MakeaTUI!\n\n")
	}
	return m.renderView()
}

