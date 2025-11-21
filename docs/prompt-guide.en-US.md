# Prompt Writing Guide

> Based on actual implementation in `internal/ai/internal/prompt` and `internal/utils/template`

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
| `name` | string | Yes | Cannot be empty (manager.go:85) |
| `description` | string | No | - |
| `variables` | map[string]string | No | - |
| `messages` | []Message | Yes | At least one (manager.go:89) |

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

Definition: `template.go:25-32`

---

## V. Manager Interface

```go
type Manager interface {
    GetTemplate(ctx, name, vars) (*Template, error)
    ListTemplates(ctx) ([]string, error)
    CreateTemplate(ctx, template) error
    UpdateTemplate(ctx, name, template) error
    DeleteTemplate(ctx, name) error
}
```

### 1. GetTemplate

```go
func (m *manager) GetTemplate(ctx context.Context, name string, vars *map[string]any) (*Template, error)
```

Behavior (manager.go:24-59):
- `vars == nil` → returns raw template (lines 36-38)
- `vars != nil` → returns template with replaced variables (lines 47-56)

Example:
```go
// Raw template
raw, _ := mgr.GetTemplate(ctx, "test", nil)

// With variable replacement
vars := map[string]any{"name": "World"}
rendered, _ := mgr.GetTemplate(ctx, "test", &vars)
```

### 2. CreateTemplate

Constraints (manager.go:79-98):
- `template.Name` cannot be empty (line 85)
- `template.Messages` must have at least one (line 89)
- File cannot already exist (line 94)

### 3. UpdateTemplate

Constraints (manager.go:100-134):
- `template.Name` cannot be empty (line 106)
- Original file must exist (line 111)
- When renaming, new name cannot already exist (line 117)

### 4. ListTemplates

Behavior (manager.go:61-77):
- Lists `*.json` file names (line 67)
- Returns list without extensions (line 74)

### 5. DeleteTemplate

Behavior (manager.go:136-142):
- Deletes `.json` file with specified name

---

## VI. Usage Example

```go
package main

import (
    "context"
    
    "github.com/Done-0/gin-scaffold/internal/ai/internal/prompt"
)

func main() {
    mgr := prompt.New()
    ctx := context.Background()
    
    // Load and replace variables
    vars := map[string]any{
        "symbol": "BTC/USDT",
        "price": 66500.00,
    }
    
    tmpl, err := mgr.GetTemplate(ctx, "trading_analyzer", &vars)
    if err != nil {
        panic(err)
    }
    
    // Use replaced content
    fmt.Println(tmpl.Messages[0].Content)
}
```

---

## VII. File Naming

- Location: `configs/prompts/`
- Format: `{name}.json`
- Naming: lowercase + underscore
- Example: `bazi_analyzer.json`

---

## VIII. Real Example

```json
{
  "name": "example",
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

