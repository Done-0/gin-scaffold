// Package internal provides i18n management functionality
// Author: Done-0
// Created: 2025-08-24
package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"

	"github.com/Done-0/gin-scaffold/internal/types/consts"
)

// Manager provides i18n functionality
type Manager struct {
	bundle *i18n.Bundle
}

// NewManager creates a new i18n manager instance
func NewManager() *Manager {
	return &Manager{}
}

// Bundle returns the i18n bundle instance
func (m *Manager) Bundle() *i18n.Bundle {
	return m.bundle
}

// Initialize sets up i18n system
func (m *Manager) Initialize() error {
	m.bundle = i18n.NewBundle(language.Chinese)
	m.bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	entries, err := os.ReadDir(consts.I18nConfigPath)
	if err != nil {
		return fmt.Errorf("failed to read i18n directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}
		if _, err := m.bundle.LoadMessageFile(filepath.Join(consts.I18nConfigPath, entry.Name())); err != nil {
			return fmt.Errorf("failed to load %s: %w", entry.Name(), err)
		}
	}

	log.Println("i18n system initialized successfully")
	return nil
}

// Close closes the i18n system
func (m *Manager) Close() error {
	m.bundle = nil
	log.Println("i18n closed successfully")
	return nil
}
