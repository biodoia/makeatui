// Package display - Tree Update and View methods
package display

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/makeatui/makeatui/pkg/widgets/mouse"
)

// GetZone returns the mouse zone
func (t *Tree) GetZone(x, y, width, height int) *mouse.Zone {
	return &mouse.Zone{
		ID:     t.ID,
		X:      x,
		Y:      y,
		Width:  width,
		Height: height,
	}
}

// Update handles messages
func (t *Tree) Update(msg tea.Msg) (*Tree, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			t.cursor--
			if t.cursor < 0 {
				t.cursor = len(t.flatNodes) - 1
			}
		case "down", "j":
			t.cursor++
			if t.cursor >= len(t.flatNodes) {
				t.cursor = 0
			}
		case "enter", " ", "right", "l":
			if node := t.GetSelectedNode(); node != nil && !node.IsLeaf() {
				node.Toggle()
				t.flatten()
			}
		case "left", "h":
			if node := t.GetSelectedNode(); node != nil {
				if node.Expanded && !node.IsLeaf() {
					node.Expanded = false
					t.flatten()
				} else if node.parent != nil {
					// Move to parent
					for i, n := range t.flatNodes {
						if n == node.parent {
							t.cursor = i
							break
						}
					}
				}
			}
		case "e":
			t.ExpandAll()
		case "c":
			t.CollapseAll()
		}
	}

	return t, nil
}

// View renders the tree
func (t *Tree) View() string {
	if t.Root == nil || len(t.flatNodes) == 0 {
		return ""
	}

	var lines []string
	for i, node := range t.flatNodes {
		line := t.renderNode(node, i)
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

// renderNode renders a single node
func (t *Tree) renderNode(node *TreeNode, index int) string {
	depth := node.Depth()
	if !t.ShowRoot {
		depth--
	}

	// Build indent and guides
	var prefix strings.Builder
	current := node
	guides := make([]string, depth)

	for i := depth - 1; i >= 0; i-- {
		parent := current.parent
		if parent != nil {
			isLast := current == parent.Children[len(parent.Children)-1]
			if i == depth-1 {
				if isLast {
					guides[i] = t.style.Guide.Corner + t.style.Guide.Horizontal
				} else {
					guides[i] = t.style.Guide.Tee + t.style.Guide.Horizontal
				}
			} else {
				if isLast {
					guides[i] = " "
				} else {
					guides[i] = t.style.Guide.Vertical
				}
			}
		}
		current = parent
	}

	for _, g := range guides {
		prefix.WriteString(t.style.Connector.Render(g))
	}

	// Expand/collapse indicator
	if !node.IsLeaf() {
		if node.Expanded {
			prefix.WriteString(t.style.IconExpand.Render(t.style.Guide.Expanded + " "))
		} else {
			prefix.WriteString(t.style.IconExpand.Render(t.style.Guide.Collapsed + " "))
		}
	} else {
		prefix.WriteString("  ")
	}

	// Icon
	if node.Icon != "" {
		prefix.WriteString(t.style.Icon.Render(node.Icon + " "))
	}

	// Label
	labelStyle := t.style.Node
	if index == t.cursor && t.Focused {
		labelStyle = t.style.NodeSel
	}
	if node.Selected {
		labelStyle = t.style.NodeSel
	}

	return prefix.String() + labelStyle.Render(node.Label)
}

// SetStyle sets the tree style
func (t *Tree) SetStyle(style TreeStyle) *Tree {
	t.style = style
	return t
}

// SetShowRoot sets whether to show root node
func (t *Tree) SetShowRoot(show bool) *Tree {
	t.ShowRoot = show
	t.flatten()
	return t
}

// FindNode finds a node by ID
func (t *Tree) FindNode(id string) *TreeNode {
	return t.findNodeRecursive(t.Root, id)
}

func (t *Tree) findNodeRecursive(node *TreeNode, id string) *TreeNode {
	if node == nil {
		return nil
	}
	if node.ID == id {
		return node
	}
	for _, child := range node.Children {
		if found := t.findNodeRecursive(child, id); found != nil {
			return found
		}
	}
	return nil
}

