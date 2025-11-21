// Package prompt provides dynamic prompt loading and management
// Author: Done-0
// Created: 2025-08-31
package prompt

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Done-0/gin-scaffold/configs"
	"github.com/Done-0/gin-scaffold/internal/utils/file"
	"github.com/Done-0/gin-scaffold/internal/utils/template"
)

type manager struct{}

// New creates a new prompt manager
func New() Manager {
	return &manager{}
}

func (m *manager) GetTemplate(ctx context.Context, name string, vars *map[string]any) (*Template, error) {
	cfgs, err := configs.GetConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get config: %w", err)
	}

	filePath := filepath.Join(cfgs.AI.Prompt.Dir, name+".json")
	var tmpl Template
	if err := file.LoadJSONFile(filePath, &tmpl); err != nil {
		return nil, fmt.Errorf("failed to load template '%s': %w", name, err)
	}

	if vars == nil {
		return &tmpl, nil
	}

	result := &Template{
		Name:        tmpl.Name,
		Description: tmpl.Description,
		Variables:   tmpl.Variables,
		Messages:    make([]Message, len(tmpl.Messages)),
	}

	for i, msg := range tmpl.Messages {
		content, err := template.Replace(msg.Content, *vars)
		if err != nil {
			return nil, fmt.Errorf("failed to replace variables in message %d: %w", i, err)
		}
		result.Messages[i] = Message{
			Role:    msg.Role,
			Content: content,
		}
	}

	return result, nil
}

func (m *manager) ListTemplates(ctx context.Context) ([]string, error) {
	cfgs, err := configs.GetConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get config: %w", err)
	}

	files, err := filepath.Glob(filepath.Join(cfgs.AI.Prompt.Dir, "*.json"))
	if err != nil {
		return nil, fmt.Errorf("failed to list templates: %w", err)
	}

	var names []string
	for _, filePath := range files {
		names = append(names, file.GetFileNameWithoutExt(filePath))
	}
	return names, nil
}

func (m *manager) CreateTemplate(ctx context.Context, template *Template) error {
	cfgs, err := configs.GetConfig()
	if err != nil {
		return fmt.Errorf("failed to get config: %w", err)
	}

	if template.Name == "" {
		return fmt.Errorf("template name cannot be empty")
	}

	if len(template.Messages) == 0 {
		return fmt.Errorf("template must have at least one message")
	}

	filePath := filepath.Join(cfgs.AI.Prompt.Dir, template.Name+".json")
	if _, err := os.Stat(filePath); err == nil {
		return fmt.Errorf("template '%s' already exists", template.Name)
	}
	return file.SaveJSONFile(filePath, template)
}

func (m *manager) UpdateTemplate(ctx context.Context, name string, template *Template) error {
	cfgs, err := configs.GetConfig()
	if err != nil {
		return fmt.Errorf("failed to get config: %w", err)
	}

	if template.Name == "" {
		return fmt.Errorf("template name cannot be empty")
	}

	oldFilePath := filepath.Join(cfgs.AI.Prompt.Dir, name+".json")
	if _, err := os.Stat(oldFilePath); os.IsNotExist(err) {
		return fmt.Errorf("template '%s' does not exist", name)
	}

	if template.Name != name {
		newFilePath := filepath.Join(cfgs.AI.Prompt.Dir, template.Name+".json")
		if _, err := os.Stat(newFilePath); err == nil {
			return fmt.Errorf("template '%s' already exists", template.Name)
		}

		if err := file.SaveJSONFile(newFilePath, template); err != nil {
			return fmt.Errorf("failed to save renamed template: %w", err)
		}

		if err := os.Remove(oldFilePath); err != nil {
			os.Remove(newFilePath)
			return fmt.Errorf("failed to remove old template file: %w", err)
		}

		return nil
	}

	return file.SaveJSONFile(oldFilePath, template)
}

func (m *manager) DeleteTemplate(ctx context.Context, name string) error {
	cfgs, err := configs.GetConfig()
	if err != nil {
		return fmt.Errorf("failed to get config: %w", err)
	}
	return os.Remove(filepath.Join(cfgs.AI.Prompt.Dir, name+".json"))
}
