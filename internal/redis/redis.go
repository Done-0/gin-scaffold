// Package redis provides Redis connection and management functionality
// Author: Done-0
// Created: 2025-09-25
package redis

import (
	"github.com/redis/go-redis/v9"

	"github.com/Done-0/gin-scaffold/configs"
	"github.com/Done-0/gin-scaffold/internal/redis/internal"
)

// RedisManager defines the interface for Redis management operations
type RedisManager interface {
	Client() *redis.Client
	Initialize() error
	Close() error
}

// New creates a new Redis manager instance
func New(config *configs.Config) (RedisManager, error) {
	return internal.NewManager(config)
}
