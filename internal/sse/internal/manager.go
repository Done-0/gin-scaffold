// Package internal provides SSE manager implementation
// Author: Done-0
// Created: 2025-08-31
package internal

import (
	"net/http"

	"github.com/gin-contrib/sse"
	"github.com/gin-gonic/gin"

	"github.com/Done-0/gin-scaffold/configs"
)

type Manager struct{}

func NewManager(*configs.Config) *Manager {
	return &Manager{}
}

func (m *Manager) StreamToClient(c *gin.Context, events <-chan *sse.Event) error {
	c.Status(http.StatusOK)
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	for event := range events {
		c.Render(-1, *event)
		c.Writer.Flush()
	}
	return nil
}
