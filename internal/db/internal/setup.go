// Package internal provides database setup and initialization functionality
// Author: Done-0
// Created: 2025-08-24
package internal

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gorm.io/gorm"
)

// setupDatabase handles database initialization based on dialect
func (m *Manager) setupDatabase() error {
	dialect := m.config.DBConfig.DBDialect

	switch dialect {
	case DialectSQLite:
		if err := m.ensureSQLiteDBExists(); err != nil {
			return fmt.Errorf("SQLite database creation failed: %w", err)
		}
	case DialectPostgres, DialectMySQL:
		if err := m.setupSystemDatabase(); err != nil {
			return fmt.Errorf("system database setup failed: %w", err)
		}
	default:
		return fmt.Errorf("unsupported database dialect: %s", dialect)
	}

	// Connect to target database
	var err error
	m.db, err = m.connectToDB(m.config.DBConfig.DBName)
	if err != nil {
		return fmt.Errorf("failed to connect to database '%s': %w", m.config.DBConfig.DBName, err)
	}

	log.Printf("Database '%s' connected successfully", m.config.DBConfig.DBName)
	return nil
}

// setupSystemDatabase handles PostgreSQL and MySQL system database setup
func (m *Manager) setupSystemDatabase() error {
	systemDBName := m.getSystemDBName()
	systemDB, err := m.connectToDB(systemDBName)
	if err != nil {
		return fmt.Errorf("failed to connect to system database '%s': %w", systemDBName, err)
	}
	defer func() {
		if sqlDB, err := systemDB.DB(); err == nil {
			sqlDB.Close()
		}
	}()

	return m.ensureDBExists(systemDB)
}

// getSystemDBName returns the system database name for the current dialect
func (m *Manager) getSystemDBName() string {
	switch m.config.DBConfig.DBDialect {
	case DialectPostgres:
		return "postgres"
	case DialectMySQL:
		return "information_schema"
	default:
		return ""
	}
}

// ensureDBExists ensures the database exists, creates it if it doesn't exist
func (m *Manager) ensureDBExists(db *gorm.DB) error {
	switch m.config.DBConfig.DBDialect {
	case DialectPostgres:
		return m.ensurePostgresDBExists(db)
	case DialectMySQL:
		return m.ensureMySQLDBExists(db)
	case DialectSQLite:
		return m.ensureSQLiteDBExists()
	default:
		return fmt.Errorf("unsupported database dialect: %s", m.config.DBConfig.DBDialect)
	}
}

// ensureSQLiteDBExists ensures SQLite database exists, creates if it doesn't exist
func (m *Manager) ensureSQLiteDBExists() error {
	dbPath := filepath.Join(m.config.DBConfig.DBPath, m.config.DBConfig.DBName+".db")

	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		if err := os.MkdirAll(m.config.DBConfig.DBPath, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create SQLite database directory: %w", err)
		}

		file, err := os.Create(dbPath)
		if err != nil {
			return fmt.Errorf("failed to create SQLite database file: %w", err)
		}
		file.Close()

		log.Printf("SQLite database '%s' created successfully", dbPath)
	}

	return nil
}

// ensurePostgresDBExists ensures PostgreSQL database exists, creates if it doesn't exist
func (m *Manager) ensurePostgresDBExists(db *gorm.DB) error {
	var exists bool
	query := "SELECT EXISTS (SELECT 1 FROM pg_database WHERE datname = ?)"
	if err := db.Raw(query, m.config.DBConfig.DBName).Scan(&exists).Error; err != nil {
		return fmt.Errorf("failed to check if PostgreSQL database exists: %w", err)
	}

	if !exists {
		createQuery := fmt.Sprintf(`CREATE DATABASE "%s" OWNER "%s"`, m.config.DBConfig.DBName, m.config.DBConfig.DBUser)
		if err := db.Exec(createQuery).Error; err != nil {
			return fmt.Errorf("failed to create PostgreSQL database: %w", err)
		}
		log.Printf("PostgreSQL database '%s' created successfully", m.config.DBConfig.DBName)
	}

	return nil
}

// ensureMySQLDBExists ensures MySQL database exists, creates if it doesn't exist
func (m *Manager) ensureMySQLDBExists(db *gorm.DB) error {
	var exists bool
	query := "SELECT EXISTS (SELECT 1 FROM information_schema.schemata WHERE schema_name = ?)"
	if err := db.Raw(query, m.config.DBConfig.DBName).Scan(&exists).Error; err != nil {
		return fmt.Errorf("failed to check if MySQL database exists: %w", err)
	}

	if !exists {
		createQuery := fmt.Sprintf("CREATE DATABASE `%s`", m.config.DBConfig.DBName)
		if err := db.Exec(createQuery).Error; err != nil {
			return fmt.Errorf("failed to create MySQL database: %w", err)
		}
		log.Printf("MySQL database '%s' created successfully", m.config.DBConfig.DBName)
	}

	return nil
}
