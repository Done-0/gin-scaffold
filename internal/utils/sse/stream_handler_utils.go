// Package sse provides Server-Sent Events streaming utilities
// Author: Done-0
// Created: 2025-08-31
package sse

import (
	"context"
	"encoding/json"

	"github.com/gin-gonic/gin"

	"github.com/Done-0/gin-scaffold/internal/sse"
)

type Handler func(ctx context.Context, ch chan<- *sse.Event)

// Stream processes data using a custom handler function
func Stream(c *gin.Context, handler Handler, manager sse.Manager) error {
	ch := make(chan *sse.Event, 100)

	go func() {
		defer close(ch)
		handler(context.Background(), ch)
	}()

	return manager.StreamToClient(c, ch)
}

// Send is a helper to emit a single event
func Send(ch chan<- *sse.Event, eventType string, data any) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}
	ch <- &sse.Event{
		Event: eventType,
		Data:  string(payload),
	}
	return nil
}
