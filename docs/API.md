# MakeaTUI API Reference

## Agent API

The Agent API provides a simple interface for AI agents to create TUI designs programmatically.

### Creating an API Instance

```go
import "github.com/makeatui/makeatui/pkg/agent"

api := agent.NewAPI("Project Name")
```

### Adding Components

#### AddBox

Add a box container with optional title.

```go
id := api.AddBox(name, title, x, y, width, height)
```

| Parameter | Type | Description |
|-----------|------|-------------|
| name | string | Component identifier |
| title | string | Box title (displayed in border) |
| x, y | int | Position on canvas |
| width, height | int | Dimensions |

#### AddText

Add a text label.

```go
id := api.AddText(name, content, x, y)
```

#### AddButton

Add a clickable button.

```go
id := api.AddButton(name, label, x, y)
```

#### AddList

Add a selectable list.

```go
id := api.AddList(name, items, x, y, width, height)
```

| Parameter | Type | Description |
|-----------|------|-------------|
| items | []string | List items |

#### AddProgress

Add a progress bar.

```go
id := api.AddProgress(name, value, x, y, width)
```

| Parameter | Type | Description |
|-----------|------|-------------|
| value | float64 | Progress value (0.0 - 1.0) |

### Modifying Components

#### Move

Move a component to a new position.

```go
api.Move(componentID, newX, newY)
```

#### Resize

Resize a component.

```go
api.Resize(componentID, newWidth, newHeight)
```

#### SetText

Update component text.

```go
api.SetText(componentID, newText)
```

#### Delete

Remove a component.

```go
api.Delete(componentID)
```

### History

#### Undo

Revert the last action.

```go
success := api.Undo()
```

#### Redo

Reapply the last undone action.

```go
success := api.Redo()
```

### Export

#### Export (Go Code)

Generate working Bubble Tea code.

```go
code := api.Export()
```

#### ExportJSON

Export canvas as JSON.

```go
jsonStr, err := api.ExportJSON()
```

### Other Methods

```go
api.Clear()                    // Remove all components
api.SetTheme("ultraviolet")    // Set theme
api.GetCanvas()                // Get canvas state
api.ListComponents()           // List all components
```

---

## AI Agent

The TUI Agent can generate interfaces from natural language.

```go
import "github.com/makeatui/makeatui/pkg/ai"

agent := ai.NewTUIAgent()
```

### Generate from Description

```go
api, err := agent.GenerateFromDescription("a chat interface with message history")
```

### Generate from Use Case

```go
api, err := agent.GenerateFromUseCase(ai.UseCaseDashboard, "My Dashboard")
```

**Use Case Types:**
- `UseCaseDashboard` - Metrics and charts
- `UseCaseForm` - Input forms
- `UseCaseBrowser` - List with preview
- `UseCaseChat` - Chat interface
- `UseCaseMonitor` - System monitoring
- `UseCaseWizard` - Multi-step wizard
- `UseCaseSettings` - Settings panel

### Get Suggestions

```go
suggestions := agent.GetSuggestions()
for _, s := range suggestions {
    fmt.Printf("%s: %s\n", s.Type, s.Description)
}
```

---

## Template Engine

```go
import "github.com/makeatui/makeatui/pkg/templates"

engine := templates.NewTemplateEngine()
```

### List Templates

```go
templates := engine.List()
```

### Apply Template

```go
api, err := engine.Apply("dashboard-basic")
```

### Create Custom Template

```go
template := templates.NewTemplateBuilder(
    "my-template",
    "My custom template",
    templates.CategoryDashboard,
).
    AddBox("header", "Title", 0, 0, 80, 3).
    AddText("content", "Hello", 5, 5).
    Build()

engine.Register(template)
```

---

## MCP Client

```go
import "github.com/makeatui/makeatui/pkg/mcp"

client := mcp.NewClient("http://localhost:8080")
```

### Session Management

```go
sessionID, _ := client.CreateSession("My Project")
```

### Component Operations

```go
client.AddBox("main", "Hello", 0, 0, 40, 10)
client.AddText("label", "Text", 5, 5)
client.Move(componentID, 10, 10)
client.Remove(componentID)
```

### AI Generation

```go
client.Generate("a dashboard with metrics")
```

### Export

```go
goCode, _ := client.Export("go")
jsonData, _ := client.Export("json")
```

### Templates

```go
templates, _ := client.ListTemplates()
client.ApplyTemplate("dashboard-basic")
```

