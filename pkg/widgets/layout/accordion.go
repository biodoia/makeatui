// Package layout - Accordion/Collapsible component (inspired by Textual)
package layout

import (
	"github.com/charmbracelet/lipgloss"
)

// AccordionItem represents a collapsible section
type AccordionItem struct {
	Title     string
	Content   string
	Expanded  bool
	Icon      AccordionIcon
}

// AccordionIcon defines expand/collapse icons
type AccordionIcon struct {
	Expanded  string
	Collapsed string
}

// DefaultAccordionIcon returns default icons
func DefaultAccordionIcon() AccordionIcon {
	return AccordionIcon{
		Expanded:  "▼",
		Collapsed: "▶",
	}
}

// Accordion implements an accordion/collapsible container
type Accordion struct {
	Items         []AccordionItem
	MultiExpand   bool // allow multiple expanded
	Animated      bool
	headerStyle   lipgloss.Style
	contentStyle  lipgloss.Style
	expandedStyle lipgloss.Style
	icon          AccordionIcon
}

// NewAccordion creates a new accordion
func NewAccordion() *Accordion {
	return &Accordion{
		Items:       []AccordionItem{},
		MultiExpand: false,
		icon:        DefaultAccordionIcon(),
		headerStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#E040FB")).
			Padding(0, 1),
		contentStyle: lipgloss.NewStyle().
			Padding(0, 2).
			Foreground(lipgloss.Color("#FFFFFF")),
		expandedStyle: lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#9D4EDD")),
	}
}

// AddItem adds an accordion item
func (a *Accordion) AddItem(title, content string) *Accordion {
	a.Items = append(a.Items, AccordionItem{
		Title:    title,
		Content:  content,
		Expanded: false,
		Icon:     a.icon,
	})
	return a
}

// Toggle toggles an item by index
func (a *Accordion) Toggle(index int) *Accordion {
	if index < 0 || index >= len(a.Items) {
		return a
	}

	if !a.MultiExpand {
		// Collapse all others
		for i := range a.Items {
			if i != index {
				a.Items[i].Expanded = false
			}
		}
	}

	a.Items[index].Expanded = !a.Items[index].Expanded
	return a
}

// ExpandAll expands all items
func (a *Accordion) ExpandAll() *Accordion {
	for i := range a.Items {
		a.Items[i].Expanded = true
	}
	return a
}

// CollapseAll collapses all items
func (a *Accordion) CollapseAll() *Accordion {
	for i := range a.Items {
		a.Items[i].Expanded = false
	}
	return a
}

// SetHeaderStyle sets the header style
func (a *Accordion) SetHeaderStyle(style lipgloss.Style) *Accordion {
	a.headerStyle = style
	return a
}

// SetContentStyle sets the content style  
func (a *Accordion) SetContentStyle(style lipgloss.Style) *Accordion {
	a.contentStyle = style
	return a
}

// SetIcon sets custom icons
func (a *Accordion) SetIcon(expanded, collapsed string) *Accordion {
	a.icon = AccordionIcon{Expanded: expanded, Collapsed: collapsed}
	for i := range a.Items {
		a.Items[i].Icon = a.icon
	}
	return a
}

// Render renders the accordion
func (a *Accordion) Render() string {
	if len(a.Items) == 0 {
		return ""
	}

	var sections []string
	for _, item := range a.Items {
		icon := item.Icon.Collapsed
		if item.Expanded {
			icon = item.Icon.Expanded
		}

		header := a.headerStyle.Render(icon + " " + item.Title)

		if item.Expanded {
			content := a.contentStyle.Render(item.Content)
			section := lipgloss.JoinVertical(lipgloss.Left, header, content)
			sections = append(sections, a.expandedStyle.Render(section))
		} else {
			sections = append(sections, header)
		}
	}

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// Collapsible is a single collapsible element
type Collapsible struct {
	Title    string
	Content  string
	Expanded bool
	style    lipgloss.Style
	icon     AccordionIcon
}

// NewCollapsible creates a new collapsible
func NewCollapsible(title, content string) *Collapsible {
	return &Collapsible{
		Title:    title,
		Content:  content,
		Expanded: false,
		icon:     DefaultAccordionIcon(),
		style: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7B2CBF")),
	}
}

// Toggle toggles the collapsible
func (c *Collapsible) Toggle() *Collapsible {
	c.Expanded = !c.Expanded
	return c
}

// Render renders the collapsible
func (c *Collapsible) Render() string {
	icon := c.icon.Collapsed
	if c.Expanded {
		icon = c.icon.Expanded
	}

	header := lipgloss.NewStyle().Bold(true).Render(icon + " " + c.Title)

	if c.Expanded {
		content := lipgloss.NewStyle().Padding(0, 2).Render(c.Content)
		return c.style.Render(lipgloss.JoinVertical(lipgloss.Left, header, content))
	}
	return c.style.Render(header)
}

