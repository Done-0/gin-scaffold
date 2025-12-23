// Package prompter provides dynamic prompt loading and management
// Author: Done-0
// Created: 2025-08-31
package prompter

import (
	"context"

	"github.com/Done-0/gin-scaffold/internal/utils/template"
)

type Prompter interface {
	GetTemplate(ctx context.Context, path string, vars *map[string]any) (*Template, error)
	ListTemplates(ctx context.Context, prefix string) ([]string, error)
	CreateTemplate(ctx context.Context, path string, tmpl *Template) error
	UpdateTemplate(ctx context.Context, path string, tmpl *Template) error
	DeleteTemplate(ctx context.Context, path string) error
}

type Template struct {
	Name        string            `json:"name"`
	Description string            `json:"description,omitempty"`
	Variables   map[string]string `json:"variables,omitempty"`
	Messages    []Message         `json:"messages"`
}

type Message = template.Message
