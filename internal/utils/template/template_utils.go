// Package template provides template variable substitution utilities
// Author: Done-0
// Created: 2025-08-31
package template

import (
	"bytes"
	"fmt"
	"text/template"
	"time"
)

// Message represents a single message in a conversation
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Replace parses and executes a Go template with the given variables and custom functions.
func Replace(text string, vars map[string]any) (string, error) {
	if text == "" {
		return "", nil
	}

	funcMap := template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
		"unixToTime": func(unixTime int64) string {
			return time.Unix(unixTime, 0).Format("2006年01月02日 15时04分")
		},
	}

	tmpl, err := template.New("prompt").Funcs(funcMap).Parse(text)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, vars); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}
