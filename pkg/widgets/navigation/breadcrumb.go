// Package navigation - Breadcrumb and Stepper components
package navigation

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/makeatui/makeatui/pkg/widgets/mouse"
)

// BreadcrumbItem represents a breadcrumb item
type BreadcrumbItem struct {
	ID    string
	Label string
	Icon  string
	Path  string
}

// Breadcrumb provides breadcrumb navigation
type Breadcrumb struct {
	ID        string
	Items     []*BreadcrumbItem
	Separator string
	Clickable bool
	style     BreadcrumbStyle
	zoneManager *mouse.ZoneManager
}

// BreadcrumbStyle holds breadcrumb styling
type BreadcrumbStyle struct {
	Item       lipgloss.Style
	ItemActive lipgloss.Style
	Separator  lipgloss.Style
	Icon       lipgloss.Style
}

// DefaultBreadcrumbStyle returns default styling
func DefaultBreadcrumbStyle() BreadcrumbStyle {
	return BreadcrumbStyle{
		Item: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7B2CBF")),
		ItemActive: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E040FB")).
			Bold(true),
		Separator: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#5A189A")),
		Icon: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9D4EDD")).
			MarginRight(1),
	}
}

// NewBreadcrumb creates a breadcrumb
func NewBreadcrumb(id string) *Breadcrumb {
	return &Breadcrumb{
		ID:          id,
		Items:       []*BreadcrumbItem{},
		Separator:   " / ",
		Clickable:   true,
		style:       DefaultBreadcrumbStyle(),
		zoneManager: mouse.NewZoneManager(),
	}
}

// AddItem adds a breadcrumb item
func (b *Breadcrumb) AddItem(id, label, icon, path string) *Breadcrumb {
	b.Items = append(b.Items, &BreadcrumbItem{
		ID:    id,
		Label: label,
		Icon:  icon,
		Path:  path,
	})
	return b
}

// SetSeparator sets the separator
func (b *Breadcrumb) SetSeparator(sep string) *Breadcrumb {
	b.Separator = sep
	return b
}

// SetPath sets breadcrumb from path string
func (b *Breadcrumb) SetPath(path, separator string) *Breadcrumb {
	parts := strings.Split(path, separator)
	b.Items = []*BreadcrumbItem{}
	currentPath := ""
	for i, part := range parts {
		if part == "" {
			continue
		}
		if i > 0 {
			currentPath += separator
		}
		currentPath += part
		b.Items = append(b.Items, &BreadcrumbItem{
			ID:    part,
			Label: part,
			Path:  currentPath,
		})
	}
	return b
}

// View renders the breadcrumb
func (b *Breadcrumb) View() string {
	var parts []string

	for i, item := range b.Items {
		isLast := i == len(b.Items)-1

		// Icon
		itemView := ""
		if item.Icon != "" {
			itemView = b.style.Icon.Render(item.Icon)
		}

		// Label
		style := b.style.Item
		if isLast {
			style = b.style.ItemActive
		}
		itemView += style.Render(item.Label)

		parts = append(parts, itemView)

		// Separator
		if !isLast {
			parts = append(parts, b.style.Separator.Render(b.Separator))
		}
	}

	return lipgloss.JoinHorizontal(lipgloss.Center, parts...)
}

// Stepper provides a step indicator (wizard-style)
type Stepper struct {
	ID        string
	Steps     []StepItem
	Current   int
	Completed []bool
	Clickable bool
	Vertical  bool
	style     StepperStyle
}

// StepItem represents a step
type StepItem struct {
	Label       string
	Description string
	Icon        string
}

// StepperStyle holds stepper styling
type StepperStyle struct {
	Step        lipgloss.Style
	StepActive  lipgloss.Style
	StepDone    lipgloss.Style
	StepPending lipgloss.Style
	Label       lipgloss.Style
	LabelActive lipgloss.Style
	Desc        lipgloss.Style
	Connector   lipgloss.Style
}

// DefaultStepperStyle returns default styling
func DefaultStepperStyle() StepperStyle {
	return StepperStyle{
		Step: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")).
			Width(3).
			Align(lipgloss.Center),
		StepActive: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E040FB")).
			Background(lipgloss.Color("#3C096C")).
			Bold(true).
			Width(3).
			Align(lipgloss.Center),
		StepDone: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6BCB77")).
			Width(3).
			Align(lipgloss.Center),
		StepPending: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#5A189A")).
			Width(3).
			Align(lipgloss.Center),
		Label: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")),
		LabelActive: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E040FB")).
			Bold(true),
		Desc: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")).
			Italic(true),
		Connector: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#5A189A")),
	}
}

// NewStepper creates a stepper
func NewStepper(id string) *Stepper {
	return &Stepper{
		ID:        id,
		Steps:     []StepItem{},
		Completed: []bool{},
		style:     DefaultStepperStyle(),
	}
}

// AddStep adds a step
func (s *Stepper) AddStep(label, description, icon string) *Stepper {
	s.Steps = append(s.Steps, StepItem{
		Label:       label,
		Description: description,
		Icon:        icon,
	})
	s.Completed = append(s.Completed, false)
	return s
}

// SetCurrent sets current step
func (s *Stepper) SetCurrent(step int) *Stepper {
	if step >= 0 && step < len(s.Steps) {
		s.Current = step
	}
	return s
}

// Next moves to next step
func (s *Stepper) Next() *Stepper {
	if s.Current < len(s.Steps)-1 {
		s.Completed[s.Current] = true
		s.Current++
	}
	return s
}

// Prev moves to previous step
func (s *Stepper) Prev() *Stepper {
	if s.Current > 0 {
		s.Current--
	}
	return s
}

// Complete marks current step as completed
func (s *Stepper) Complete() *Stepper {
	if s.Current >= 0 && s.Current < len(s.Steps) {
		s.Completed[s.Current] = true
	}
	return s
}

// Update handles messages
func (s *Stepper) Update(msg tea.Msg) (*Stepper, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "left", "h":
			s.Prev()
		case "right", "l":
			s.Next()
		case "enter", " ":
			s.Complete()
			s.Next()
		}
	}
	return s, nil
}

