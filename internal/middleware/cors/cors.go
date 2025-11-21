// Package cors provides Gin middleware for CORS management
// Author: Done-0
// Created: 2025-09-25
package cors

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/Done-0/gin-scaffold/configs"
)

// New creates a Gin middleware for CORS management
func New(config *configs.Config) gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     config.AppConfig.CORSConfig.AllowOrigins,
		AllowMethods:     config.AppConfig.CORSConfig.AllowMethods,
		AllowHeaders:     config.AppConfig.CORSConfig.AllowHeaders,
		ExposeHeaders:    config.AppConfig.CORSConfig.ExposeHeaders,
		AllowCredentials: config.AppConfig.CORSConfig.AllowCredentials,
		MaxAge:           time.Duration(config.AppConfig.CORSConfig.MaxAge) * time.Hour,
	})
}
