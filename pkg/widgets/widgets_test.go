// Package widgets provides tests for the widget library
package widgets

import (
	"testing"

	"github.com/makeatui/makeatui/pkg/widgets/advanced"
	"github.com/makeatui/makeatui/pkg/widgets/display"
	"github.com/makeatui/makeatui/pkg/widgets/effects"
	"github.com/makeatui/makeatui/pkg/widgets/feedback"
	"github.com/makeatui/makeatui/pkg/widgets/input"
	"github.com/makeatui/makeatui/pkg/widgets/layout"
	"github.com/makeatui/makeatui/pkg/widgets/mouse"
	"github.com/makeatui/makeatui/pkg/widgets/navigation"
)

// Test Layout Components
func TestGrid(t *testing.T) {
	g := layout.NewGrid(3, 3)
	if g == nil {
		t.Fatal("Grid should not be nil")
	}

	g.SetGap(1)
	g.SetCell(0, 0, "cell1")
	g.SetCell(1, 1, "cell2")

	view := g.Render()
	if view == "" {
		t.Error("Grid Render should not be empty")
	}
}

func TestFlex(t *testing.T) {
	f := layout.NewFlex(layout.FlexRow)
	if f == nil {
		t.Fatal("Flex should not be nil")
	}

	item := layout.NewFlexItem("content")
	if item == nil {
		t.Fatal("FlexItem should not be nil")
	}

	f.AddItem(item)
	f.Add("content2")
}

func TestSplit(t *testing.T) {
	s := layout.NewSplit(layout.SplitHorizontal, 0.5)
	if s == nil {
		t.Fatal("Split should not be nil")
	}

	s.SetPanes("left", "right")
	s.SetSize(80, 24)

	view := s.Render()
	if view == "" {
		t.Error("Split Render should not be empty")
	}
}

// Test Input Components
func TestTextInput(t *testing.T) {
	ti := input.NewTextInput("test-input")
	if ti == nil {
		t.Fatal("TextInput should not be nil")
	}

	ti.SetLabel("Name").SetPlaceholder("Enter name")
	ti.SetRequired(true)

	if !ti.Required {
		t.Error("TextInput should be required")
	}
}

func TestSelect(t *testing.T) {
	s := input.NewSelect("test-select")
	if s == nil {
		t.Fatal("Select should not be nil")
	}

	s.AddOption("1", "Option 1")
	s.AddOption("2", "Option 2")

	if len(s.Options) != 2 {
		t.Errorf("Expected 2 options, got %d", len(s.Options))
	}
}

func TestSlider(t *testing.T) {
	s := input.NewSlider("test-slider", 0, 100)
	if s == nil {
		t.Fatal("Slider should not be nil")
	}

	s.SetValue(50).SetStep(5)

	if s.Value != 50 {
		t.Errorf("Expected value 50, got %f", s.Value)
	}
}

// Test Display Components
func TestTable(t *testing.T) {
	tbl := display.NewTable("test-table")
	if tbl == nil {
		t.Fatal("Table should not be nil")
	}

	tbl.AddColumn("name", "Name", 20)
	tbl.AddColumn("value", "Value", 20)
	tbl.AddRow(map[string]string{"name": "Row1", "value": "Data1"})

	if len(tbl.Rows) != 1 {
		t.Errorf("Expected 1 row, got %d", len(tbl.Rows))
	}
}

func TestTree(t *testing.T) {
	tree := display.NewTree("test-tree")
	if tree == nil {
		t.Fatal("Tree should not be nil")
	}

	root := display.NewTreeNode("root", "Root")
	tree.SetRoot(root)
}

// Test Feedback Components
func TestToast(t *testing.T) {
	toast := feedback.NewToast("test-toast", "Message", feedback.ToastInfo)
	if toast == nil {
		t.Fatal("Toast should not be nil")
	}

	view := toast.View()
	if view == "" {
		t.Error("Toast View should not be empty")
	}
}

// Test Effects
func TestMatrixRain(t *testing.T) {
	m := effects.NewMatrixRain(40, 20)
	if m == nil {
		t.Fatal("MatrixRain should not be nil")
	}

	view := m.View()
	if view == "" {
		t.Error("MatrixRain View should not be empty")
	}
}

func TestTypingEffect(t *testing.T) {
	te := effects.NewTypingEffect("Hello World")
	if te == nil {
		t.Fatal("TypingEffect should not be nil")
	}
}

// Test Mouse Support
func TestZoneManager(t *testing.T) {
	zm := mouse.NewZoneManager()
	if zm == nil {
		t.Fatal("ZoneManager should not be nil")
	}

	zone := &mouse.Zone{
		ID:     "test-zone",
		X:      0,
		Y:      0,
		Width:  10,
		Height: 5,
	}
	zm.Register(zone)

	found := zm.HitTest(5, 2)
	if found == nil || found.ID != "test-zone" {
		t.Error("Should find zone at coordinates")
	}
}

// Test Navigation Components
func TestMenu(t *testing.T) {
	m := navigation.NewMenu("test-menu")
	if m == nil {
		t.Fatal("Menu should not be nil")
	}

	m.AddItem("item1", "Item 1", "", "")
	m.AddItem("item2", "Item 2", "", "")

	if len(m.Items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(m.Items))
	}
}

func TestBreadcrumb(t *testing.T) {
	b := navigation.NewBreadcrumb("test-breadcrumb")
	if b == nil {
		t.Fatal("Breadcrumb should not be nil")
	}

	b.AddItem("home", "Home", "", "/")
	b.AddItem("products", "Products", "", "/products")

	view := b.View()
	if view == "" {
		t.Error("Breadcrumb View should not be empty")
	}
}

// Test Advanced Components
func TestFileManager(t *testing.T) {
	fm := advanced.NewFileManager("test-fm", 60, 20)
	if fm == nil {
		t.Fatal("FileManager should not be nil")
	}

	// Should start in current directory
	if fm.CurrentDir == "" {
		t.Error("FileManager should have a current directory")
	}

	view := fm.View()
	if view == "" {
		t.Error("FileManager View should not be empty")
	}
}

func TestWindow(t *testing.T) {
	w := advanced.NewWindow("test-window", "Test Window", 10, 5, 40, 20)
	if w == nil {
		t.Fatal("Window should not be nil")
	}

	w.SetContent("Window content")

	if w.X != 10 || w.Y != 5 {
		t.Error("Window position should be set correctly")
	}
}

func TestWindowManager(t *testing.T) {
	wm := advanced.NewWindowManager("wm", 80, 24)
	if wm == nil {
		t.Fatal("WindowManager should not be nil")
	}

	w := advanced.NewWindow("win1", "Window 1", 0, 0, 30, 15)
	wm.AddWindow(w)

	if len(wm.Windows) != 1 {
		t.Errorf("Expected 1 window, got %d", len(wm.Windows))
	}
}

