// Package internal provides Redis manager implementation
// Author: Done-0
// Created: 2025-09-25
package internal

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/Done-0/gin-scaffold/configs"
)

// Manager represents a Redis manager with dependency injection
type Manager struct {
	config *configs.Config
	client *redis.Client
}

// NewManager creates a new Redis manager instance
func NewManager(config *configs.Config) (*Manager, error) {
	return &Manager{
		config: config,
	}, nil
}

// Client returns the Redis client instance
func (m *Manager) Client() *redis.Client {
	return m.client
}

// Initialize sets up the Redis connection
func (m *Manager) Initialize() error {
	db, _ := strconv.Atoi(m.config.RedisConfig.RedisDB)
	m.client = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", m.config.RedisConfig.RedisHost, m.config.RedisConfig.RedisPort),
		Password:     m.config.RedisConfig.RedisPassword,                             // Database password, default empty string
		DB:           db,                                                             // Database index
		DialTimeout:  time.Duration(m.config.RedisConfig.DialTimeout) * time.Second,  // Connection timeout
		ReadTimeout:  time.Duration(m.config.RedisConfig.ReadTimeout) * time.Second,  // Read timeout
		WriteTimeout: time.Duration(m.config.RedisConfig.WriteTimeout) * time.Second, // Write timeout
		PoolSize:     m.config.RedisConfig.PoolSize,                                  // Maximum connection pool size
		MinIdleConns: m.config.RedisConfig.MinIdleConns,                              // Minimum idle connections
	})

	if err := m.client.Ping(context.Background()).Err(); err != nil {
		return fmt.Errorf("failed to ping Redis: %w", err)
	}

	log.Println("Redis connected successfully")
	return nil
}

// Close closes the Redis connection
func (m *Manager) Close() error {
	if m.client == nil {
		return nil
	}

	if err := m.client.Close(); err != nil {
		return fmt.Errorf("failed to close Redis connection: %w", err)
	}

	log.Println("Redis connection closed successfully")
	return nil
}
