// Package router provides application route registration functionality
// Author: Done-0
// Created: 2025-09-25
package router

import (
	"github.com/gin-gonic/gin"

	"github.com/Done-0/gin-scaffold/pkg/router/routes"
	"github.com/Done-0/gin-scaffold/pkg/wire"
)

// New registers application routes
func New(r *gin.Engine, container *wire.Container) {
	// Create API v1 route group
	v1 := r.Group("/api/v1")

	// Create API v2 route group
	v2 := r.Group("/api/v2")

	// Register routes by modules
	routes.RegisterTestRoutes(container, v1, v2)
}
