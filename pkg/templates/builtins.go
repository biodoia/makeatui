// Package templates - Built-in template definitions
package templates

// DashboardBasic creates a basic dashboard template
func DashboardBasic() *Template {
	return NewTemplateBuilder(
		"dashboard-basic",
		"A basic dashboard with header, sidebar, and main content area",
		CategoryDashboard,
	).
		AddBox("header", "Dashboard", 0, 0, 80, 3).
		AddBox("sidebar", "Menu", 0, 4, 20, 17).
		AddBox("main", "", 21, 4, 59, 17).
		AddText("title", "Welcome to Dashboard", 23, 6).
		SetPreview(`
┌──────────────────────────────────────────────────────────────────────────────┐
│ Dashboard                                                                    │
├───────────────────┬──────────────────────────────────────────────────────────┤
│ Menu              │                                                          │
│                   │  Welcome to Dashboard                                    │
│ > Home            │                                                          │
│   Settings        │                                                          │
│   Users           │                                                          │
│   Reports         │                                                          │
│                   │                                                          │
└───────────────────┴──────────────────────────────────────────────────────────┘
		`).
		Build()
}

// DashboardMetrics creates a metrics dashboard template
func DashboardMetrics() *Template {
	return NewTemplateBuilder(
		"dashboard-metrics",
		"A dashboard with metric cards and charts",
		CategoryDashboard,
	).
		AddBox("header", "System Metrics", 0, 0, 80, 3).
		AddBox("metric1", "CPU", 0, 4, 19, 7).
		AddProgress("cpu-bar", 0.65, 2, 8, 15).
		AddBox("metric2", "Memory", 20, 4, 19, 7).
		AddProgress("mem-bar", 0.42, 22, 8, 15).
		AddBox("metric3", "Disk", 40, 4, 19, 7).
		AddProgress("disk-bar", 0.78, 42, 8, 15).
		AddBox("metric4", "Network", 60, 4, 19, 7).
		AddProgress("net-bar", 0.23, 62, 8, 15).
		AddBox("chart", "Activity", 0, 12, 80, 9).
		SetPreview(`
┌──────────────────────────────────────────────────────────────────────────────┐
│ System Metrics                                                               │
├──────────────────┬──────────────────┬──────────────────┬─────────────────────┤
│ CPU              │ Memory           │ Disk             │ Network             │
│ ████████░░ 65%   │ ████░░░░░░ 42%   │ ████████░░ 78%   │ ██░░░░░░░░ 23%      │
├──────────────────┴──────────────────┴──────────────────┴─────────────────────┤
│ Activity                                                                     │
│     ▄▄ ▄▄    ▄▄        ▄▄ ▄▄                                                │
│   ▄▄██▄██▄▄▄▄██▄▄    ▄▄██▄██▄▄                                              │
└──────────────────────────────────────────────────────────────────────────────┘
		`).
		Build()
}

// FormSimple creates a simple form template
func FormSimple() *Template {
	return NewTemplateBuilder(
		"form-simple",
		"A simple form with input fields",
		CategoryForm,
	).
		AddBox("form", "User Registration", 10, 2, 60, 18).
		AddText("label1", "Username:", 12, 5).
		AddText("label2", "Email:", 12, 8).
		AddText("label3", "Password:", 12, 11).
		AddButton("submit", "[ Submit ]", 30, 15).
		AddButton("cancel", "[ Cancel ]", 45, 15).
		SetPreview(`
          ┌──────────────────────────────────────────────────────────┐
          │ User Registration                                        │
          │                                                          │
          │  Username:   ┌────────────────────────────────────┐      │
          │              └────────────────────────────────────┘      │
          │  Email:      ┌────────────────────────────────────┐      │
          │              └────────────────────────────────────┘      │
          │  Password:   ┌────────────────────────────────────┐      │
          │              └────────────────────────────────────┘      │
          │                                                          │
          │                    [ Submit ]   [ Cancel ]               │
          └──────────────────────────────────────────────────────────┘
		`).
		Build()
}

// FormMultiStep creates a multi-step wizard form
func FormMultiStep() *Template {
	return NewTemplateBuilder(
		"form-wizard",
		"A multi-step wizard form",
		CategoryWizard,
	).
		AddBox("container", "Setup Wizard", 5, 1, 70, 20).
		AddText("step", "Step 1 of 3: Basic Info", 8, 3).
		AddProgress("progress", 0.33, 8, 5, 64).
		AddBox("content", "", 8, 7, 64, 10).
		AddButton("prev", "[ Previous ]", 8, 18).
		AddButton("next", "[ Next ]", 55, 18).
		Build()
}

// ListBrowser creates a file/item browser template
func ListBrowser() *Template {
	return NewTemplateBuilder(
		"list-browser",
		"A list browser with preview panel",
		CategoryList,
	).
		AddBox("list", "Files", 0, 0, 30, 21).
		AddList("filelist", []string{"README.md", "main.go", "go.mod"}, 1, 2, 28, 18).
		AddBox("preview", "Preview", 31, 0, 49, 21).
		AddText("preview-content", "Select a file to preview", 33, 3).
		Build()
}

// ChatInterface creates a chat-style interface
func ChatInterface() *Template {
	return NewTemplateBuilder(
		"chat-interface",
		"A chat interface with message history and input",
		CategoryChat,
	).
		AddBox("messages", "Chat", 0, 0, 80, 18).
		AddBox("input", "", 0, 19, 80, 3).
		AddText("prompt", "> ", 2, 20).
		Build()
}

// MonitorSystem creates a system monitor template
func MonitorSystem() *Template {
	return NewTemplateBuilder(
		"monitor-system",
		"A system monitoring dashboard",
		CategoryMonitor,
	).
		AddBox("processes", "Processes", 0, 0, 50, 12).
		AddBox("logs", "Logs", 51, 0, 29, 12).
		AddBox("resources", "Resources", 0, 13, 80, 8).
		AddProgress("cpu", 0.45, 2, 15, 25).
		AddProgress("mem", 0.62, 2, 17, 25).
		AddProgress("disk", 0.33, 2, 19, 25).
		Build()
}

// SettingsPanel creates a settings panel template
func SettingsPanel() *Template {
	return NewTemplateBuilder(
		"settings-panel",
		"A settings panel with categories",
		CategorySettings,
	).
		AddBox("categories", "Settings", 0, 0, 25, 21).
		AddList("cat-list", []string{"General", "Appearance", "Keyboard", "Advanced"}, 1, 2, 23, 18).
		AddBox("options", "General Settings", 26, 0, 54, 21).
		AddText("opt1", "[ ] Dark Mode", 28, 3).
		AddText("opt2", "[ ] Auto Save", 28, 5).
		AddText("opt3", "[ ] Notifications", 28, 7).
		Build()
}

