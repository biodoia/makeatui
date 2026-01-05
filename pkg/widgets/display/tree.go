// Package display - Tree component (inspired by Textual/Ratatui)
package display

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/makeatui/makeatui/pkg/widgets/mouse"
)

// TreeNode represents a node in the tree
type TreeNode struct {
	ID       string
	Label    string
	Icon     string
	Data     interface{}
	Children []*TreeNode
	Expanded bool
	Selected bool
	parent   *TreeNode
}

// TreeStyle holds tree styling
type TreeStyle struct {
	Node        lipgloss.Style
	NodeSel     lipgloss.Style
	NodeHover   lipgloss.Style
	Icon        lipgloss.Style
	IconExpand  lipgloss.Style
	Connector   lipgloss.Style
	Guide       GuideStyle
}

// GuideStyle defines tree guide characters
type GuideStyle struct {
	Vertical   string
	Horizontal string
	Corner     string
	Tee        string
	Expanded   string
	Collapsed  string
}

// DefaultGuideStyle returns default guide characters
func DefaultGuideStyle() GuideStyle {
	return GuideStyle{
		Vertical:   "│",
		Horizontal: "─",
		Corner:     "└",
		Tee:        "├",
		Expanded:   "▼",
		Collapsed:  "▶",
	}
}

// RoundGuideStyle returns rounded guide style
func RoundGuideStyle() GuideStyle {
	return GuideStyle{
		Vertical:   "│",
		Horizontal: "─",
		Corner:     "╰",
		Tee:        "├",
		Expanded:   "▾",
		Collapsed:  "▸",
	}
}

// DefaultTreeStyle returns default styling
func DefaultTreeStyle() TreeStyle {
	return TreeStyle{
		Node: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")),
		NodeSel: lipgloss.NewStyle().
			Background(lipgloss.Color("#9D4EDD")).
			Foreground(lipgloss.Color("#FFFFFF")).
			Bold(true),
		NodeHover: lipgloss.NewStyle().
			Background(lipgloss.Color("#3C096C")).
			Foreground(lipgloss.Color("#FFFFFF")),
		Icon: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E040FB")),
		IconExpand: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7B2CBF")),
		Connector: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#5A189A")),
		Guide: DefaultGuideStyle(),
	}
}

// Tree provides a tree view component
type Tree struct {
	ID          string
	Root        *TreeNode
	ShowRoot    bool
	Focused     bool
	cursor      int
	flatNodes   []*TreeNode // flattened visible nodes
	style       TreeStyle
	zoneManager *mouse.ZoneManager
}

// NewTree creates a new tree
func NewTree(id string) *Tree {
	return &Tree{
		ID:          id,
		ShowRoot:    true,
		style:       DefaultTreeStyle(),
		zoneManager: mouse.NewZoneManager(),
	}
}

// SetRoot sets the root node
func (t *Tree) SetRoot(node *TreeNode) *Tree {
	t.Root = node
	t.flatten()
	return t
}

// NewTreeNode creates a new tree node
func NewTreeNode(id, label string) *TreeNode {
	return &TreeNode{
		ID:       id,
		Label:    label,
		Children: []*TreeNode{},
	}
}

// AddChild adds a child node
func (n *TreeNode) AddChild(child *TreeNode) *TreeNode {
	child.parent = n
	n.Children = append(n.Children, child)
	return n
}

// SetIcon sets the node icon
func (n *TreeNode) SetIcon(icon string) *TreeNode {
	n.Icon = icon
	return n
}

// SetExpanded sets expanded state
func (n *TreeNode) SetExpanded(expanded bool) *TreeNode {
	n.Expanded = expanded
	return n
}

// Toggle toggles expansion
func (n *TreeNode) Toggle() *TreeNode {
	n.Expanded = !n.Expanded
	return n
}

// IsLeaf returns true if node has no children
func (n *TreeNode) IsLeaf() bool {
	return len(n.Children) == 0
}

// Depth returns node depth
func (n *TreeNode) Depth() int {
	depth := 0
	current := n.parent
	for current != nil {
		depth++
		current = current.parent
	}
	return depth
}

// flatten creates a flat list of visible nodes
func (t *Tree) flatten() {
	t.flatNodes = []*TreeNode{}
	if t.Root == nil {
		return
	}
	t.flattenNode(t.Root, t.ShowRoot)
}

func (t *Tree) flattenNode(node *TreeNode, include bool) {
	if include {
		t.flatNodes = append(t.flatNodes, node)
	}

	if node.Expanded || !include {
		for _, child := range node.Children {
			t.flattenNode(child, true)
		}
	}
}

// GetSelectedNode returns the selected node
func (t *Tree) GetSelectedNode() *TreeNode {
	if t.cursor >= 0 && t.cursor < len(t.flatNodes) {
		return t.flatNodes[t.cursor]
	}
	return nil
}

// ExpandAll expands all nodes
func (t *Tree) ExpandAll() *Tree {
	if t.Root != nil {
		t.expandAllNode(t.Root)
		t.flatten()
	}
	return t
}

func (t *Tree) expandAllNode(node *TreeNode) {
	node.Expanded = true
	for _, child := range node.Children {
		t.expandAllNode(child)
	}
}

// CollapseAll collapses all nodes
func (t *Tree) CollapseAll() *Tree {
	if t.Root != nil {
		t.collapseAllNode(t.Root)
		t.flatten()
	}
	return t
}

func (t *Tree) collapseAllNode(node *TreeNode) {
	node.Expanded = false
	for _, child := range node.Children {
		t.collapseAllNode(child)
	}
}

