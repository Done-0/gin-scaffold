// Package internal provides database connection functionality
// Author: Done-0
// Created: 2025-08-24
package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// connectToDB establishes a connection to the specified database
func (m *Manager) connectToDB(dbName string) (*gorm.DB, error) {
	dialector, err := m.getDialector(dbName)
	if err != nil {
		return nil, fmt.Errorf("failed to get database dialector: %w", err)
	}

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get SQL database instance: %w", err)
	}

	sqlDB.SetMaxOpenConns(50)                 // Maximum number of open connections
	sqlDB.SetMaxIdleConns(20)                 // Maximum number of idle connections
	sqlDB.SetConnMaxLifetime(5 * time.Minute) // Maximum connection lifetime
	sqlDB.SetConnMaxIdleTime(1 * time.Minute) // Maximum idle time before closing

	return db, nil
}

// getDialector returns the appropriate database dialector based on configuration
func (m *Manager) getDialector(dbName string) (gorm.Dialector, error) {
	switch m.config.DBConfig.DBDialect {
	case DialectPostgres:
		return m.getPostgresDialector(dbName), nil
	case DialectMySQL:
		return m.getMySQLDialector(dbName), nil
	case DialectSQLite:
		return m.getSQLiteDialector(dbName)
	default:
		return nil, fmt.Errorf("unsupported database dialect: %s", m.config.DBConfig.DBDialect)
	}
}

// getSQLiteDialector returns SQLite dialector and ensures directory exists
func (m *Manager) getSQLiteDialector(dbName string) (gorm.Dialector, error) {
	if err := os.MkdirAll(m.config.DBConfig.DBPath, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create SQLite database directory: %w", err)
	}

	dbPath := filepath.Join(m.config.DBConfig.DBPath, dbName+".db")
	return sqlite.Open(dbPath), nil
}

// getPostgresDialector returns PostgreSQL dialector
func (m *Manager) getPostgresDialector(dbName string) gorm.Dialector {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai client_encoding=UTF8",
		m.config.DBConfig.DBHost, m.config.DBConfig.DBUser, m.config.DBConfig.DBPassword, dbName, m.config.DBConfig.DBPort,
	)
	return postgres.Open(dsn)
}

// getMySQLDialector returns MySQL dialector
func (m *Manager) getMySQLDialector(dbName string) gorm.Dialector {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		m.config.DBConfig.DBUser, m.config.DBConfig.DBPassword, m.config.DBConfig.DBHost, m.config.DBConfig.DBPort, dbName,
	)
	return mysql.Open(dsn)
}
