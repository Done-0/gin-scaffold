# Prompt Writing Guide

> Based on actual implementation in `internal/ai/internal/prompter` and `internal/utils/template`

---

## I. JSON Structure

```json
{
  "name": "string",
  "description": "string (optional)",
  "variables": {
    "key": "description text"
  },
  "messages": [
    {
      "role": "system|user|assistant",
      "content": "text content, supports template syntax"
    }
  ]
}
```

---

## II. Type Definitions

```go
// internal/ai/internal/prompt/types.go
type Template struct {
    Name        string            `json:"name"`
    Description string            `json:"description,omitempty"`
    Variables   map[string]string `json:"variables,omitempty"`
    Messages    []Message         `json:"messages"`
}

type Message = template.Message

// internal/utils/template/template.go
type Message struct {
    Role    string `json:"role"`
    Content string `json:"content"`
}
```

---

## III. Field Constraints

| Field | Type | Required | Constraint |
|-------|------|----------|------------|
| `name` | string | Yes | Cannot be empty |
| `description` | string | No | - |
| `variables` | map[string]string | No | - |
| `messages` | []Message | Yes | At least one |

---

## IV. Template Syntax

### 1. Variables

```go
{{.variable}}
```

### 2. Conditionals

```go
{{if .condition}}...{{else}}...{{end}}
{{if gt .value 10}}...{{end}}
{{if eq .status "active"}}...{{end}}
```

### 3. Loops

```go
{{range $index, $item := .list}}
  {{$index}} - {{$item.field}}
{{end}}
```

### 4. Built-in Functions

```go
{{add 1 2}}                    // returns 3
{{unixToTime 1706140800}}      // returns "2025年01月24日 15时30分"
```

---

## V. Prompter Interface

```go
type Prompter interface {
    GetTemplate(ctx, path, vars) (*Template, error)
    ListTemplates(ctx, prefix) ([]string, error)
    CreateTemplate(ctx, path, tmpl) error
    UpdateTemplate(ctx, path, tmpl) error
    DeleteTemplate(ctx, path) error
}
```

### 1. GetTemplate

```go
func (p *prompter) GetTemplate(ctx context.Context, path string, vars *map[string]any) (*Template, error)
```

Behavior:
- `path` parameter is the file path relative to prompts directory (without `.json` suffix)
- `vars == nil` → returns raw template
- `vars != nil` → returns template with replaced variables
- The `name` field inside JSON is purely descriptive metadata, does not affect lookup

Example:
```go
// configs/prompts/test.json
raw, _ := p.GetTemplate(ctx, "test", nil)

// configs/prompts/character/character_generator.json
tmpl, _ := p.GetTemplate(ctx, "character/character_generator", nil)

// configs/prompts/stories/horror/elevator_game.json
tmpl, _ := p.GetTemplate(ctx, "stories/horror/elevator_game", nil)

// With variable replacement
vars := map[string]any{"name": "World"}
rendered, _ := p.GetTemplate(ctx, "test", &vars)
```

### 2. ListTemplates

```go
func (p *prompter) ListTemplates(ctx context.Context, prefix string) ([]string, error)
```

Behavior:
- `prefix == ""` → lists all templates
- `prefix != ""` → lists only templates under specified directory
- Returns relative paths without `.json` extension

Example:
```go
// List all templates
names, _ := p.ListTemplates(ctx, "")
// ["character_generator", "stories/midnight_store", "stories/horror/elevator_game"]

// List only stories/horror directory
names, _ := p.ListTemplates(ctx, "stories/horror")
// ["stories/horror/elevator_game"]
```

### 3. CreateTemplate

```go
func (p *prompter) CreateTemplate(ctx context.Context, path string, tmpl *Template) error
```

Constraints:
- `path` cannot be empty
- `tmpl.Messages` must have at least one
- File cannot already exist
- Auto-creates intermediate directories
- `tmpl.Name` is metadata only, can be named freely

Example:
```go
tmpl := &prompt.Template{
    Name:     "Midnight Store",  // Display name, can be anything
    Messages: []prompt.Message{{Role: "system", Content: "..."}},
}
p.CreateTemplate(ctx, "stories/midnight_store", tmpl) // Path determines storage location
```

### 4. UpdateTemplate

```go
func (p *prompter) UpdateTemplate(ctx context.Context, path string, tmpl *Template) error
```

Constraints:
- `path` cannot be empty
- File must exist
- Updates in place, does not support moving (use delete + create to move)

Example:
```go
tmpl := &prompt.Template{
    Name:     "Updated Name",  // Only updates metadata
    Messages: []prompt.Message{{Role: "system", Content: "Updated content"}},
}
p.UpdateTemplate(ctx, "stories/midnight_store", tmpl)
```

### 5. DeleteTemplate

```go
func (p *prompter) DeleteTemplate(ctx context.Context, path string) error
```

Behavior:
- If file: deletes `{path}.json`
- If directory: deletes entire directory and all contents

Example:
```go
// Delete single file
p.DeleteTemplate(ctx, "stories/midnight_store")

// Delete entire directory
p.DeleteTemplate(ctx, "stories") // Deletes stories/ and all sub-contents
```

---

## VI. Usage Example

```go
package main

import (
    "context"

    "github.com/Done-0/gin-scaffold/internal/ai/internal/prompt"
)

func main() {
    p := prompt.New()
    ctx := context.Background()

    // Load and replace variables
    vars := map[string]any{
        "symbol": "BTC/USDT",
        "price": 66500.00,
    }

    tmpl, err := p.GetTemplate(ctx, "trading_analyzer", &vars)
    if err != nil {
        panic(err)
    }

    // Use replaced content
    fmt.Println(tmpl.Messages[0].Content)
}
```

---

## VII. File Structure

```
configs/prompts/
├── character_generator.json     # Root directory template
├── stories/                     # Subdirectory
│   ├── midnight_store.json
│   └── horror/                  # Nested subdirectory
│       └── elevator_game.json
└── ...
```

- Location: specified by config `AI.Prompt.Dir`, default `configs/prompts/`
- Format: `{path}.json`
- Naming: lowercase + underscore, supports multi-level directories

---

## VIII. Real Example

```json
{
  "name": "stories/example",
  "description": "Example prompt",
  "variables": {
    "user_name": "User name",
    "user_age": "User age"
  },
  "messages": [
    {
      "role": "system",
      "content": "You are an AI assistant. Current user: {{.user_name}}, age: {{.user_age}}."
    },
    {
      "role": "user",
      "content": "{{.user_message}}"
    }
  ]
}
```
