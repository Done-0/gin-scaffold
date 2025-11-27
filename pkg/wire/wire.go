//go:build wireinject

// Package wire provides Wire dependency injection definitions
// Author: Done-0
// Created: 2025-09-25
package wire

import (
	"github.com/google/wire"

	"github.com/Done-0/gin-scaffold/configs"
	"github.com/Done-0/gin-scaffold/internal/ai"
	"github.com/Done-0/gin-scaffold/internal/db"
	"github.com/Done-0/gin-scaffold/internal/i18n"
	"github.com/Done-0/gin-scaffold/internal/logger"
	"github.com/Done-0/gin-scaffold/internal/sse"

	// "github.com/Done-0/gin-scaffold/internal/queue"

	"github.com/Done-0/gin-scaffold/internal/redis"
	"github.com/Done-0/gin-scaffold/pkg/serve/controller"
)

// Container holds all application dependencies
type Container struct {
	Config *configs.Config

	// Infrastructure
	AIManager       *ai.AIManager
	DatabaseManager db.DatabaseManager
	RedisManager    redis.RedisManager
	LoggerManager   logger.LoggerManager
	I18nManager     i18n.I18nManager
	SSEManager      sse.SSEManager
	// QueueProducer   queue.Producer

	// Controllers
	TestController *controller.TestController

	// Services

	// mappers
}

// NewContainer initializes the complete application container using Wire
func NewContainer(config *configs.Config) (*Container, error) {
	panic(wire.Build(
		AllProviders,
		wire.Struct(new(Container), "*"),
	))
}
