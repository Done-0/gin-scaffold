// Package prompt provides dynamic prompt loading and management
// Author: Done-0
// Created: 2025-08-31
package prompt

import (
	"context"

	"github.com/Done-0/gin-scaffold/internal/utils/template"
)

type Manager interface {
	GetTemplate(ctx context.Context, name string, vars *map[string]any) (*Template, error)
	ListTemplates(ctx context.Context) ([]string, error)
	CreateTemplate(ctx context.Context, template *Template) error
	UpdateTemplate(ctx context.Context, name string, template *Template) error
	DeleteTemplate(ctx context.Context, name string) error
}

type Template struct {
	Name        string            `json:"name"`
	Description string            `json:"description,omitempty"`
	Variables   map[string]string `json:"variables,omitempty"`
	Messages    []Message         `json:"messages"`
}

type Message = template.Message
