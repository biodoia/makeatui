# MakeaTUI 🎨

> A visual TUI design framework for AI agents, built with the Charm ecosystem.

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev)
[![License](https://img.shields.io/badge/License-MIT-purple.svg)](LICENSE)

MakeaTUI enables AI agents to design beautiful terminal user interfaces through a visual builder, programmatic API, or natural language descriptions. Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea), [Lip Gloss](https://github.com/charmbracelet/lipgloss), and the entire Charm ecosystem.

## ✨ Features

- **🎨 Visual Designer** - Interactive TUI builder with sidebar, canvas, and properties panel
- **🤖 AI Agent API** - Programmatic interface for automated TUI generation
- **📝 Natural Language** - Describe your interface, get working code
- **🎭 Template Engine** - Pre-built templates for common use cases
- **🌈 Ultraviolet Theme** - Stunning purple/magenta color palette
- **📹 VHS Integration** - Record and replay TUI sessions
- **🔌 MCP Server** - Model Context Protocol for AI integration
- **⚡ Code Generation** - Export to working Bubble Tea applications

## 🚀 Quick Start

```bash
# Clone and build
git clone https://github.com/biodoia/makeatui.git
cd makeatui
go build -o makeatui .

# Run the visual designer
./makeatui

# Or start the MCP server
./makeatui serve --port 8080
```

## 📦 Installation

```bash
go get github.com/makeatui/makeatui
```

## 🤖 Agent API Usage

```go
package main

import (
    "fmt"
    "github.com/makeatui/makeatui/pkg/agent"
)

func main() {
    // Create a new API instance
    api := agent.NewAPI("My Dashboard")
    
    // Add components
    api.AddBox("header", "Dashboard", 0, 0, 80, 3)
    api.AddBox("sidebar", "Menu", 0, 4, 20, 17)
    api.AddButton("save", "Save", 2, 6)
    api.AddProgress("cpu", 0.65, 25, 6, 30)
    
    // Export to Go code
    code := api.Export()
    fmt.Println(code)
}
```

## 🧠 AI Generation

```go
package main

import (
    "github.com/makeatui/makeatui/pkg/ai"
)

func main() {
    agent := ai.NewTUIAgent()
    
    // Generate from description
    api, _ := agent.GenerateFromDescription(
        "A system monitoring dashboard with CPU, memory, and disk metrics",
    )
    
    // Export the generated design
    code := api.Export()
}
```

## 🎭 Templates

```go
import "github.com/makeatui/makeatui/pkg/templates"

engine := templates.NewTemplateEngine()

// List available templates
for _, t := range engine.List() {
    fmt.Printf("%s: %s\n", t.Name, t.Description)
}

// Apply a template
api, _ := engine.Apply("dashboard-metrics")
```

**Built-in Templates:**
- `dashboard-basic` - Header, sidebar, main content
- `dashboard-metrics` - Metrics cards with progress bars
- `form-simple` - Registration form with inputs
- `form-wizard` - Multi-step wizard
- `list-browser` - File browser with preview
- `chat-interface` - Chat with message history
- `monitor-system` - System monitoring dashboard
- `settings-panel` - Settings with categories

## 🔌 MCP Server

Start the MCP server for AI agent integration:

```bash
./makeatui serve --port 8080
```

### Available Tools

| Tool | Description |
|------|-------------|
| `makeatui_create_session` | Create a design session |
| `makeatui_add_box` | Add a box component |
| `makeatui_add_text` | Add text |
| `makeatui_add_button` | Add a button |
| `makeatui_generate` | Generate from description |
| `makeatui_export` | Export as Go code or JSON |
| `makeatui_apply_template` | Apply a template |

### Client Example

```go
import "github.com/makeatui/makeatui/pkg/mcp"

client := mcp.NewClient("http://localhost:8080")
client.CreateSession("My Project")
client.AddBox("main", "Hello", 0, 0, 40, 10)
code, _ := client.Export("go")
```

## 📹 VHS Recording

Generate demo GIFs with VHS:

```go
import "github.com/makeatui/makeatui/pkg/vhs"

tape := vhs.CreateDemoTape("demo.gif")
tape.Run()
```

## 🎨 Theming

```go
import "github.com/makeatui/makeatui/internal/ui/styles"

// Use the Ultraviolet theme
theme := styles.Ultraviolet
styles := styles.NewStyles(theme)

// Or the Neon theme
theme := styles.Neon
```

## 📁 Project Structure

```
makeatui/
├── cmd/makeatui/       # CLI entry point
├── internal/
│   ├── app/            # Bubble Tea application
│   ├── codegen/        # Go code generator
│   └── ui/
│       ├── animation/  # Harmonica animations
│       ├── canvas/     # Design canvas
│       ├── components/ # UI components
│       ├── markdown/   # Glamour rendering
│       └── styles/     # Themes
├── pkg/
│   ├── agent/          # Agent API
│   ├── ai/             # AI TUI agent
│   ├── mcp/            # MCP server/client
│   ├── schema/         # Component schemas
│   ├── scripting/      # Gum scripting
│   ├── templates/      # Template engine
│   └── vhs/            # VHS integration
└── examples/
```

## ⌨️ Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `Tab` | Switch focus |
| `↑/↓/←/→` | Navigate |
| `Enter` | Select/Confirm |
| `d` | Delete component |
| `u` | Undo |
| `r` | Redo |
| `e` | Export |
| `?` | Help |
| `q` | Quit |

## 🛠️ Built With

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Styling
- [Bubbles](https://github.com/charmbracelet/bubbles) - Components
- [Glamour](https://github.com/charmbracelet/glamour) - Markdown
- [Harmonica](https://github.com/charmbracelet/harmonica) - Animations
- [Log](https://github.com/charmbracelet/log) - Logging

## 📄 License

MIT License - see [LICENSE](LICENSE) for details.

---

Made with 💜 for AI agents who deserve beautiful TUIs.

