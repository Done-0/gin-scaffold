// Package internal provides database migration functionality
// Author: Done-0
// Created: 2025-08-24
package internal

import (
	"fmt"
	"log"

	"github.com/Done-0/gin-scaffold/internal/model"
)

// migrate performs database auto migration
func (m *Manager) migrate() error {
	err := m.db.AutoMigrate(
		model.GetAllModels()...,
	)
	if err != nil {
		return fmt.Errorf("failed to auto migrate database: %w", err)
	}

	log.Println("Database auto migration succeeded")
	return nil
}
