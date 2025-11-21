# Prompt 书写规范

> 基于 `internal/ai/internal/prompt` 和 `internal/utils/template` 的实际实现

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
| `name` | string | 是 | 不能为空（manager.go:85） |
| `description` | string | 否 | - |
| `variables` | map[string]string | 否 | - |
| `messages` | []Message | 是 | 至少一条（manager.go:89） |

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

定义位置：`template.go:25-32`

---

## 五、Manager 接口

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

行为（manager.go:24-59）：
- `vars == nil` → 返回原始模板（第36-38行）
- `vars != nil` → 替换变量后返回（第47-56行）

示例：
```go
// 原始模板
raw, _ := mgr.GetTemplate(ctx, "test", nil)

// 替换变量
vars := map[string]any{"name": "World"}
rendered, _ := mgr.GetTemplate(ctx, "test", &vars)
```

### 2. CreateTemplate

约束（manager.go:79-98）：
- `template.Name` 不能为空（第85行）
- `template.Messages` 至少一条（第89行）
- 文件不能已存在（第94行）

### 3. UpdateTemplate

约束（manager.go:100-134）：
- `template.Name` 不能为空（第106行）
- 原文件必须存在（第111行）
- 重命名时新名称不能已存在（第117行）

### 4. ListTemplates

行为（manager.go:61-77）：
- 列出 `*.json` 文件名（第67行）
- 返回不含扩展名的列表（第74行）

### 5. DeleteTemplate

行为（manager.go:136-142）：
- 删除指定名称的 `.json` 文件

---

## 六、使用示例

```go
package main

import (
    "context"
    
    "github.com/Done-0/gin-scaffold/internal/ai/internal/prompt"
)

func main() {
    mgr := prompt.New()
    ctx := context.Background()
    
    // 加载并替换变量
    vars := map[string]any{
        "symbol": "BTC/USDT",
        "price": 66500.00,
    }
    
    tmpl, err := mgr.GetTemplate(ctx, "trading_analyzer", &vars)
    if err != nil {
        panic(err)
    }
    
    // 使用替换后的内容
    fmt.Println(tmpl.Messages[0].Content)
}
```

---

## 七、文件命名

- 位置：`configs/prompts/`
- 格式：`{name}.json`
- 命名：小写字母+下划线
- 示例：`bazi_analyzer.json`

---

## 八、实际示例

```json
{
  "name": "example",
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

