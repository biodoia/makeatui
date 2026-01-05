# MakeaTUI MCP Integration

MakeaTUI implements the Model Context Protocol (MCP) for seamless AI agent integration.

## Starting the Server

```bash
./makeatui serve --port 8080
```

## API Endpoints

### Sessions

#### Create Session
```
POST /sessions
Content-Type: application/json

{"name": "My Project"}
```

Response:
```json
{"id": "session_1", "name": "My Project"}
```

#### List Sessions
```
GET /sessions
```

#### Get Session
```
GET /sessions/{id}
```

#### Delete Session
```
DELETE /sessions/{id}
```

### Tools

#### Add Component
```
POST /tools/add_component
Content-Type: application/json

{
  "session_id": "session_1",
  "type": "box",
  "name": "main",
  "text": "Hello",
  "x": 0,
  "y": 0,
  "width": 40,
  "height": 10
}
```

#### Move Component
```
POST /tools/move_component

{
  "session_id": "session_1",
  "component_id": "comp_xxx",
  "x": 10,
  "y": 5
}
```

#### Remove Component
```
POST /tools/remove_component

{
  "session_id": "session_1",
  "component_id": "comp_xxx"
}
```

#### Set Text
```
POST /tools/set_text

{
  "session_id": "session_1",
  "component_id": "comp_xxx",
  "text": "New text"
}
```

#### Generate from Description
```
POST /tools/generate

{
  "session_id": "session_1",
  "description": "A dashboard with metrics and charts"
}
```

#### Export
```
GET /tools/export?session_id=session_1&format=go
GET /tools/export?session_id=session_1&format=json
```

### Templates

#### List Templates
```
GET /templates
```

#### Apply Template
```
POST /templates/apply

{
  "session_id": "session_1",
  "template_name": "dashboard-basic"
}
```

### Resources

#### Get Canvas
```
GET /resources/canvas?session_id=session_1
```

#### Get Components
```
GET /resources/components?session_id=session_1
```

### Health
```
GET /health
```

## MCP Tool Schemas

For AI model integration, use the tool schemas:

```go
import "github.com/makeatui/makeatui/pkg/mcp"

schemas := mcp.GetToolSchemas()
```

## Claude/ChatGPT Integration Example

```python
# Example tool use with Claude
tools = [
    {
        "name": "makeatui_generate",
        "description": "Generate a TUI layout from description",
        "input_schema": {
            "type": "object",
            "properties": {
                "description": {
                    "type": "string",
                    "description": "Natural language description"
                }
            },
            "required": ["description"]
        }
    }
]
```

## WebSocket Support (Planned)

Real-time updates via WebSocket:

```
WS /ws/session/{id}
```

Events:
- `component_added`
- `component_moved`
- `component_removed`
- `canvas_updated`

