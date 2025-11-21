// Package internal provides database management functionality
// Author: Done-0
// Created: 2025-08-24
package internal

import (
	"fmt"
	"log"

	"gorm.io/gorm"

	"github.com/Done-0/gin-scaffold/configs"
)

// Database dialect constants
const (
	DialectPostgres = "postgres" // PostgreSQL database
	DialectSQLite   = "sqlite"   // SQLite database
	DialectMySQL    = "mysql"    // MySQL database
)

// Manager represents a database manager with dependency injection
type Manager struct {
	config *configs.Config
	db     *gorm.DB
}

// NewManager creates a new database manager instance
func NewManager(config *configs.Config) *Manager {
	return &Manager{
		config: config,
	}
}

// DB returns the database instance
func (m *Manager) DB() *gorm.DB {
	return m.db
}

// Initialize sets up the database connection and performs migrations
func (m *Manager) Initialize() error {
	if err := m.setupDatabase(); err != nil {
		return fmt.Errorf("failed to setup database: %w", err)
	}

	if err := m.migrate(); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("Database initialized successfully")
	return nil
}

// Close closes the database connection
func (m *Manager) Close() error {
	if m.db == nil {
		return nil
	}

	sqlDB, err := m.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get SQL database instance: %w", err)
	}

	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}

	log.Println("Database closed successfully")
	return nil
}
