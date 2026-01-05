// Package layout - Flexbox layout component
package layout

import (
	"github.com/charmbracelet/lipgloss"
)

// FlexDirection defines flex direction
type FlexDirection int

const (
	FlexRow FlexDirection = iota
	FlexColumn
	FlexRowReverse
	FlexColumnReverse
)

// FlexAlign defines alignment
type FlexAlign int

const (
	AlignStart FlexAlign = iota
	AlignCenter
	AlignEnd
	AlignStretch
	AlignBaseline
)

// FlexJustify defines content justification
type FlexJustify int

const (
	JustifyStart FlexJustify = iota
	JustifyCenter
	JustifyEnd
	JustifySpaceBetween
	JustifySpaceAround
	JustifySpaceEvenly
)

// FlexItem represents an item in a flex container
type FlexItem struct {
	Content string
	Grow    int
	Shrink  int
	Basis   int
	Style   lipgloss.Style
}

// NewFlexItem creates a new flex item
func NewFlexItem(content string) *FlexItem {
	return &FlexItem{
		Content: content,
		Grow:    0,
		Shrink:  1,
		Basis:   0,
	}
}

// SetGrow sets the grow factor
func (f *FlexItem) SetGrow(grow int) *FlexItem {
	f.Grow = grow
	return f
}

// Flex implements a CSS Flexbox-like layout
type Flex struct {
	Direction FlexDirection
	Align     FlexAlign
	Justify   FlexJustify
	Gap       int
	Wrap      bool
	items     []*FlexItem
	width     int
	height    int
	style     lipgloss.Style
}

// NewFlex creates a new flex container
func NewFlex(direction FlexDirection) *Flex {
	return &Flex{
		Direction: direction,
		Align:     AlignStart,
		Justify:   JustifyStart,
		Gap:       1,
		items:     []*FlexItem{},
		style:     lipgloss.NewStyle(),
	}
}

// AddItem adds an item to the flex container
func (f *Flex) AddItem(item *FlexItem) *Flex {
	f.items = append(f.items, item)
	return f
}

// Add adds content with default flex settings
func (f *Flex) Add(content string) *Flex {
	return f.AddItem(NewFlexItem(content))
}

// SetSize sets the container size
func (f *Flex) SetSize(width, height int) *Flex {
	f.width = width
	f.height = height
	return f
}

// SetStyle sets the container style
func (f *Flex) SetStyle(style lipgloss.Style) *Flex {
	f.style = style
	return f
}

// Render renders the flex container
func (f *Flex) Render() string {
	if len(f.items) == 0 {
		return ""
	}

	contents := make([]string, len(f.items))
	for i, item := range f.items {
		contents[i] = item.Content
	}

	var position lipgloss.Position
	switch f.Align {
	case AlignStart:
		position = lipgloss.Top
	case AlignCenter:
		position = lipgloss.Center
	case AlignEnd:
		position = lipgloss.Bottom
	default:
		position = lipgloss.Top
	}

	var result string
	switch f.Direction {
	case FlexRow:
		result = lipgloss.JoinHorizontal(position, interleaveGap(contents, f.Gap)...)
	case FlexColumn:
		result = lipgloss.JoinVertical(lipgloss.Position(f.Align), interleaveGapVertical(contents, f.Gap)...)
	case FlexRowReverse:
		reversed := reverseStrings(contents)
		result = lipgloss.JoinHorizontal(position, interleaveGap(reversed, f.Gap)...)
	case FlexColumnReverse:
		reversed := reverseStrings(contents)
		result = lipgloss.JoinVertical(lipgloss.Position(f.Align), interleaveGapVertical(reversed, f.Gap)...)
	}

	return f.style.Render(result)
}

func interleaveGap(items []string, gap int) []string {
	if len(items) <= 1 || gap <= 0 {
		return items
	}
	result := make([]string, 0, len(items)*2-1)
	spacer := lipgloss.NewStyle().Width(gap).Render("")
	for i, item := range items {
		result = append(result, item)
		if i < len(items)-1 {
			result = append(result, spacer)
		}
	}
	return result
}

func interleaveGapVertical(items []string, gap int) []string {
	if len(items) <= 1 || gap <= 0 {
		return items
	}
	result := make([]string, 0, len(items)*2-1)
	spacer := ""
	for range gap {
		spacer += "\n"
	}
	for i, item := range items {
		result = append(result, item)
		if i < len(items)-1 {
			result = append(result, spacer)
		}
	}
	return result
}

func reverseStrings(items []string) []string {
	result := make([]string, len(items))
	for i, item := range items {
		result[len(items)-1-i] = item
	}
	return result
}

