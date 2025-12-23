// Package prompter provides dynamic prompt loading and management
// Author: Done-0
// Created: 2025-08-31
package prompter

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Done-0/gin-scaffold/configs"
	"github.com/Done-0/gin-scaffold/internal/utils/file"
	"github.com/Done-0/gin-scaffold/internal/utils/template"
)

type prompter struct{}

func New() Prompter {
	return &prompter{}
}

func (p *prompter) GetTemplate(ctx context.Context, path string, vars *map[string]any) (*Template, error) {
	cfg, err := configs.GetConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get config: %w", err)
	}

	var tmpl Template
	if err := file.LoadJSONFile(filepath.Join(cfg.AI.Prompt.Dir, path+".json"), &tmpl); err != nil {
		return nil, fmt.Errorf("failed to load template '%s': %w", path, err)
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
		result.Messages[i] = Message{Role: msg.Role, Content: content}
	}
	return result, nil
}

func (p *prompter) ListTemplates(ctx context.Context, prefix string) ([]string, error) {
	cfg, err := configs.GetConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get config: %w", err)
	}

	baseDir := cfg.AI.Prompt.Dir
	searchDir := baseDir
	if prefix != "" {
		searchDir = filepath.Join(baseDir, prefix)
	}

	var names []string
	err = filepath.Walk(searchDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".json") {
			return err
		}
		relPath, _ := filepath.Rel(baseDir, path)
		names = append(names, filepath.ToSlash(strings.TrimSuffix(relPath, ".json")))
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list templates: %w", err)
	}
	return names, nil
}

func (p *prompter) CreateTemplate(ctx context.Context, path string, tmpl *Template) error {
	cfg, err := configs.GetConfig()
	if err != nil {
		return fmt.Errorf("failed to get config: %w", err)
	}

	if path == "" {
		return fmt.Errorf("path cannot be empty")
	}
	if len(tmpl.Messages) == 0 {
		return fmt.Errorf("template must have at least one message")
	}

	filePath := filepath.Join(cfg.AI.Prompt.Dir, path+".json")
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	if _, err := os.Stat(filePath); err == nil {
		return fmt.Errorf("template '%s' already exists", path)
	}
	return file.SaveJSONFile(filePath, tmpl)
}

func (p *prompter) UpdateTemplate(ctx context.Context, path string, tmpl *Template) error {
	cfg, err := configs.GetConfig()
	if err != nil {
		return fmt.Errorf("failed to get config: %w", err)
	}

	if path == "" {
		return fmt.Errorf("path cannot be empty")
	}

	filePath := filepath.Join(cfg.AI.Prompt.Dir, path+".json")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("template '%s' does not exist", path)
	}

	return file.SaveJSONFile(filePath, tmpl)
}

func (p *prompter) DeleteTemplate(ctx context.Context, path string) error {
	cfg, err := configs.GetConfig()
	if err != nil {
		return fmt.Errorf("failed to get config: %w", err)
	}

	fullPath := filepath.Join(cfg.AI.Prompt.Dir, path)

	if info, err := os.Stat(fullPath + ".json"); err == nil && !info.IsDir() {
		return os.Remove(fullPath + ".json")
	}

	if info, err := os.Stat(fullPath); err == nil && info.IsDir() {
		return os.RemoveAll(fullPath)
	}

	return fmt.Errorf("template '%s' not found", path)
}
