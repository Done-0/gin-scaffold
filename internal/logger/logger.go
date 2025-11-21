// Package logger provides application logging functionality initialization and configuration
// Author: Done-0
// Created: 2025-09-25
package logger

import (
	"github.com/sirupsen/logrus"

	"github.com/Done-0/gin-scaffold/configs"
	"github.com/Done-0/gin-scaffold/internal/logger/internal"
)

// LoggerManager defines the interface for logger management operations
type LoggerManager interface {
	Logger() *logrus.Logger
	Initialize() error
	Close() error
}

// New creates a new logger manager instance
func New(config *configs.Config) (LoggerManager, error) {
	return internal.NewManager(config)
}
