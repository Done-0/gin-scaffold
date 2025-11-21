// Package internal provides AI service internal implementation
// Author: Done-0
// Created: 2025-08-31
package internal

import (
	"github.com/Done-0/gin-scaffold/configs"
	"github.com/Done-0/gin-scaffold/internal/ai/internal/prompt"
	"github.com/Done-0/gin-scaffold/internal/ai/internal/provider"
)

type Manager struct {
	provider.Provider
	prompt.Manager
}

// New creates a new AI provider manager with dynamic prompt loading
func New(config *configs.Config) (*Manager, error) {
	return &Manager{
		Provider: provider.New(),
		Manager:  prompt.New(),
	}, nil
}
