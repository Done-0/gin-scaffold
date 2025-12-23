# Prompt 书写规范

> 基于 `internal/ai/internal/prompter` 和 `internal/utils/template` 的实际实现

---

## 一、JSON 结构

```json
{
  "name": "string",
  "description": "string (可选)",
  "variables": {
    "key": "说明文本"
  },
  "messages": [
    {
      "role": "system|user|assistant",
      "content": "文本内容，支持模板语法"
    }
  ]
}
```

---

## 二、类型定义

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

## 三、字段约束

| 字段 | 类型 | 必填 | 约束 |
|-----|------|-----|------|
| `name` | string | 是 | 不能为空 |
| `description` | string | 否 | - |
| `variables` | map[string]string | 否 | - |
| `messages` | []Message | 是 | 至少一条 |

---

## 四、模板语法

### 1. 变量

```go
{{.variable}}
```

### 2. 条件

```go
{{if .condition}}...{{else}}...{{end}}
{{if gt .value 10}}...{{end}}
{{if eq .status "active"}}...{{end}}
```

### 3. 循环

```go
{{range $index, $item := .list}}
  {{$index}} - {{$item.field}}
{{end}}
```

### 4. 内置函数

```go
{{add 1 2}}                    // 返回 3
{{unixToTime 1706140800}}      // 返回 "2025年01月24日 15时30分"
```

---

## 五、Prompter 接口

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

行为：
- `path` 参数是相对于 prompts 目录的文件路径（不含 `.json` 后缀）
- `vars == nil` → 返回原始模板
- `vars != nil` → 替换变量后返回
- JSON 内的 `name` 字段仅作为描述性元数据，不影响查找

示例：
```go
// configs/prompts/test.json
raw, _ := p.GetTemplate(ctx, "test", nil)

// configs/prompts/character/character_generator.json
tmpl, _ := p.GetTemplate(ctx, "character/character_generator", nil)

// configs/prompts/stories/horror/elevator_game.json
tmpl, _ := p.GetTemplate(ctx, "stories/horror/elevator_game", nil)

// 替换变量
vars := map[string]any{"name": "World"}
rendered, _ := p.GetTemplate(ctx, "test", &vars)
```

### 2. ListTemplates

```go
func (p *prompter) ListTemplates(ctx context.Context, prefix string) ([]string, error)
```

行为：
- `prefix == ""` → 列出所有模板
- `prefix != ""` → 只列出指定目录下的模板
- 返回相对路径，不含 `.json` 扩展名

示例：
```go
// 列出所有模板
names, _ := p.ListTemplates(ctx, "")
// ["character_generator", "stories/midnight_store", "stories/horror/elevator_game"]

// 只列出 stories/horror 目录
names, _ := p.ListTemplates(ctx, "stories/horror")
// ["stories/horror/elevator_game"]
```

### 3. CreateTemplate

```go
func (p *prompter) CreateTemplate(ctx context.Context, path string, tmpl *Template) error
```

约束：
- `path` 不能为空
- `tmpl.Messages` 至少一条
- 文件不能已存在
- 自动创建中间目录
- `tmpl.Name` 仅作为元数据，可随意命名

示例：
```go
tmpl := &prompt.Template{
    Name:     "午夜便利店",  // 显示名称，随意命名
    Messages: []prompt.Message{{Role: "system", Content: "..."}},
}
p.CreateTemplate(ctx, "stories/midnight_store", tmpl) // 路径决定存储位置
```

### 4. UpdateTemplate

```go
func (p *prompter) UpdateTemplate(ctx context.Context, path string, tmpl *Template) error
```

约束：
- `path` 不能为空
- 文件必须存在
- 原地更新，不支持移动（如需移动，先删除再创建）

示例：
```go
tmpl := &prompt.Template{
    Name:     "更新后的名称",  // 仅更新元数据
    Messages: []prompt.Message{{Role: "system", Content: "Updated content"}},
}
p.UpdateTemplate(ctx, "stories/midnight_store", tmpl)
```

### 5. DeleteTemplate

```go
func (p *prompter) DeleteTemplate(ctx context.Context, path string) error
```

行为：
- 如果是文件：删除 `{name}.json`
- 如果是目录：删除整个目录及其所有内容

示例：
```go
// 删除单个文件
p.DeleteTemplate(ctx, "stories/midnight_store")

// 删除整个目录
p.DeleteTemplate(ctx, "stories") // 删除 stories/ 及其所有子内容
```

---

## 六、使用示例

```go
package main

import (
    "context"

    "github.com/Done-0/gin-scaffold/internal/ai/internal/prompt"
)

func main() {
    p := prompt.New()
    ctx := context.Background()

    // 加载并替换变量
    vars := map[string]any{
        "symbol": "BTC/USDT",
        "price": 66500.00,
    }

    tmpl, err := p.GetTemplate(ctx, "trading_analyzer", &vars)
    if err != nil {
        panic(err)
    }

    // 使用替换后的内容
    fmt.Println(tmpl.Messages[0].Content)
}
```

---

## 七、文件结构

```
configs/prompts/
├── character_generator.json     # 根目录模板
├── stories/                     # 子目录
│   ├── midnight_store.json
│   └── horror/                  # 嵌套子目录
│       └── elevator_game.json
└── ...
```

- 位置：由配置 `AI.Prompt.Dir` 指定，默认 `configs/prompts/`
- 格式：`{path}.json`
- 命名：小写字母+下划线，支持多级目录

---

## 八、实际示例

```json
{
  "name": "stories/example",
  "description": "示例提示词",
  "variables": {
    "user_name": "用户姓名",
    "user_age": "用户年龄"
  },
  "messages": [
    {
      "role": "system",
      "content": "你是AI助手。当前用户：{{.user_name}}，年龄：{{.user_age}}岁。"
    },
    {
      "role": "user",
      "content": "{{.user_message}}"
    }
  ]
}
```
