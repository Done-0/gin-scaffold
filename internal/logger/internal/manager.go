// Package internal provides logger management functionality
// Author: Done-0
// Created: 2025-09-25
package internal

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"time"

	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"

	"github.com/Done-0/gin-scaffold/configs"

	rotateLogs "github.com/lestrrat-go/file-rotatelogs"
)

// Manager represents a logger manager with dependency injection
type Manager struct {
	config  *configs.Config
	logger  *logrus.Logger
	logFile io.Closer
}

// NewManager creates a new logger manager instance
func NewManager(config *configs.Config) (*Manager, error) {
	return &Manager{
		config: config,
	}, nil
}

// Logger returns the logger instance
func (m *Manager) Logger() *logrus.Logger {
	return m.logger
}

// Initialize sets up the logger system
func (m *Manager) Initialize() error {
	logFilePath := m.config.LogConfig.LogFilePath
	logFileName := m.config.LogConfig.LogFileName
	fileName := path.Join(logFilePath, logFileName)
	_ = os.MkdirAll(logFilePath, 0755)

	// Initialize logger
	formatter := &logrus.JSONFormatter{TimestampFormat: "2006-01-02 15:04:05"}
	m.logger = logrus.New()
	m.logger.SetFormatter(formatter)

	// Set log level
	if logLevel, err := logrus.ParseLevel(m.config.LogConfig.LogLevel); err == nil {
		m.logger.SetLevel(logLevel)
	} else {
		m.logger.SetLevel(logrus.InfoLevel)
	}

	// Configure log rotation
	writer, err := rotateLogs.New(
		path.Join(logFilePath, "%Y%m%d.log"),
		rotateLogs.WithLinkName(fileName),
		rotateLogs.WithMaxAge(time.Duration(m.config.LogConfig.LogMaxAge)*24*time.Hour),
		rotateLogs.WithRotationTime(24*time.Hour),
	)

	if err != nil {
		log.Printf("Failed to initialize log file rotation: %v, using standard output", err)
		fileHandle, fileErr := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0755)

		if fileErr != nil {
			log.Printf("Failed to create log file: %v, using standard output", fileErr)
			m.logger.SetOutput(os.Stdout)
			m.logFile = nil
		} else {
			m.logger.SetOutput(io.MultiWriter(os.Stdout, fileHandle))
			m.logFile = fileHandle
		}
	} else {
		allLevels := []logrus.Level{
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
			logrus.WarnLevel,
			logrus.InfoLevel,
			logrus.DebugLevel,
			logrus.TraceLevel,
		}

		writeMap := make(lfshook.WriterMap, len(allLevels))
		for _, level := range allLevels {
			writeMap[level] = writer
		}

		m.logger.AddHook(lfshook.NewHook(writeMap, formatter))
		m.logger.SetOutput(os.Stdout)
		m.logFile = writer
	}

	log.Println("Logger system initialized successfully")
	return nil
}

// Close closes the logger system
func (m *Manager) Close() error {
	if m.logger == nil {
		return nil
	}

	m.logger.ReplaceHooks(make(logrus.LevelHooks))

	if m.logFile != nil {
		if err := m.logFile.Close(); err != nil {
			return fmt.Errorf("failed to close log file: %w", err)
		}
	}

	m.logger = nil
	m.logFile = nil

	log.Println("Logger system closed successfully")
	return nil
}
