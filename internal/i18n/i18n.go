// Package i18n provides internationalization management functionality
// Author: Done-0
// Created: 2025-08-24
package i18n

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"

	"github.com/Done-0/gin-scaffold/internal/i18n/internal"
)

// I18nManager defines the interface for i18n management
type I18nManager interface {
	Bundle() *i18n.Bundle
	Initialize() error
	Close() error
}

// New creates a new i18n manager instance
func New() I18nManager {
	return internal.NewManager()
}
