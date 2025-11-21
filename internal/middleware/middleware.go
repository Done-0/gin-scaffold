// Package middleware provides common middleware functionality
// Author: Done-0
// Created: 2025-09-25
package middleware

import (
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"

	"github.com/Done-0/gin-scaffold/configs"
	"github.com/Done-0/gin-scaffold/internal/middleware/cors"
)

// New registers all middleware to the Gin engine
func New(r *gin.Engine, config *configs.Config) {
	// Recovery middleware (should be first)
	r.Use(gin.Recovery())

	// Request ID middleware
	r.Use(requestid.New())

	// CORS middleware
	r.Use(cors.New(config))

	// Logger middleware
	r.Use(gin.Logger())
}
