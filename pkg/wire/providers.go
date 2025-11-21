// Package wire provides dependency injection configuration using Google Wire
// Author: Done-0
// Created: 2025-09-25
package wire

import (
	"github.com/google/wire"

	"github.com/Done-0/gin-scaffold/internal/ai"
	"github.com/Done-0/gin-scaffold/internal/db"
	"github.com/Done-0/gin-scaffold/internal/i18n"
	"github.com/Done-0/gin-scaffold/internal/logger"
	"github.com/Done-0/gin-scaffold/internal/sse"

	// "github.com/Done-0/gin-scaffold/internal/queue"

	"github.com/Done-0/gin-scaffold/internal/redis"
	"github.com/Done-0/gin-scaffold/pkg/serve/controller"
	"github.com/Done-0/gin-scaffold/pkg/serve/service/impl"
)

// InfrastructureProviders provides infrastructure layer dependencies
var InfrastructureProviders = wire.NewSet(
	ai.New,
	db.New,
	logger.New,
	i18n.New,
	// queue.NewProducer,
	redis.New,
	sse.New,
)

// MapperProviders provides data access layer dependencies
var MapperProviders = wire.NewSet()

// ServiceProviders provides business logic layer dependencies
var ServiceProviders = wire.NewSet(
	impl.NewTestService,
)

// ControllerProviders provides controller layer dependencies
var ControllerProviders = wire.NewSet(
	controller.NewTestController,
)

// AllProviders combines all provider sets in dependency order
var AllProviders = wire.NewSet(
	InfrastructureProviders,
	MapperProviders,
	ServiceProviders,
	ControllerProviders,
)
