// Package cmd provides application startup and runtime entry point
// Author: Done-0
// Created: 2025-09-25
package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/Done-0/gin-scaffold/configs"
	"github.com/Done-0/gin-scaffold/internal/middleware"
	"github.com/Done-0/gin-scaffold/pkg/router"
	"github.com/Done-0/gin-scaffold/pkg/wire"
)

// Start starts the HTTP server with graceful shutdown handling
func Start() {
	if err := configs.New(); err != nil {
		log.Fatalf("Failed to initialize config: %v", err)
	}

	cfgs, err := configs.GetConfig()
	if err != nil {
		log.Fatalf("Failed to get config: %v", err)
	}

	container, err := wire.NewContainer(cfgs)
	if err != nil {
		log.Fatalf("Failed to initialize container: %v", err)
	}

	// Initialize managers
	if err := container.LoggerManager.Initialize(); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer container.LoggerManager.Close()

	if err := container.DatabaseManager.Initialize(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer container.DatabaseManager.Close()

	if err := container.RedisManager.Initialize(); err != nil {
		log.Fatalf("Failed to initialize Redis: %v", err)
	}
	defer container.RedisManager.Close()

	// MQ consumers
	go func() {
		// TODO: Add specific consumer startup logic here
	}()

	// Set Gin mode based on environment
	env := os.Getenv("ENV")
	switch env {
	case "prod", "production":
		gin.SetMode(gin.ReleaseMode)
	case "dev", "development":
		gin.SetMode(gin.DebugMode)
	default:
		gin.SetMode(gin.DebugMode)
	}

	// Create Gin engine
	r := gin.New()
	middleware.New(r, cfgs)
	router.New(r, container)

	// Create HTTP server
	serverAddr := fmt.Sprintf("%s:%s", cfgs.AppConfig.AppHost, cfgs.AppConfig.AppPort)
	srv := &http.Server{
		Addr:    serverAddr,
		Handler: r,
	}

	// Start server in goroutine
	go func() {
		log.Printf("⇨ Gin server starting on %s", serverAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
