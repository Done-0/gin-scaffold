// Package impl provides test service implementation
// Author: Done-0
// Created: 2025-09-25
package impl

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/gin-contrib/sse"
	"github.com/gin-gonic/gin"

	"github.com/Done-0/gin-scaffold/internal/ai"
	"github.com/Done-0/gin-scaffold/internal/logger"
	"github.com/Done-0/gin-scaffold/internal/redis"
	"github.com/Done-0/gin-scaffold/internal/types/errno"
	"github.com/Done-0/gin-scaffold/pkg/serve/controller/dto"
	"github.com/Done-0/gin-scaffold/pkg/serve/service"
	"github.com/Done-0/gin-scaffold/pkg/vo"
)

// TestServiceImpl test service implementation
type TestServiceImpl struct {
	loggerManager logger.LoggerManager
	redisManager  redis.RedisManager
	aiManager     *ai.AIManager
}

// NewTestService creates test service implementation
func NewTestService(loggerManager logger.LoggerManager, redisManager redis.RedisManager, aiManager *ai.AIManager) service.TestService {
	return &TestServiceImpl{
		loggerManager: loggerManager,
		redisManager:  redisManager,
		aiManager:     aiManager,
	}
}

// TestPing handles ping test
func (ts *TestServiceImpl) TestPing(c *gin.Context) (*vo.TestPingResponse, error) {
	return &vo.TestPingResponse{
		Message: "Pong successfully!",
		Time:    time.Now().Format(time.RFC3339),
	}, nil
}

// TestHello handles hello test
func (ts *TestServiceImpl) TestHello(c *gin.Context) (*vo.TestHelloResponse, error) {
	return &vo.TestHelloResponse{
		Message: "Hello, gin-scaffold! 🎉!",
		Version: "1.0.0",
	}, nil
}

// TestLogger handles logger test
func (ts *TestServiceImpl) TestLogger(c *gin.Context) (*vo.TestLoggerResponse, error) {
	logger := ts.loggerManager.Logger()
	logger.Info("Test logger endpoint called")

	return &vo.TestLoggerResponse{
		Message: "Log test succeeded!",
		Level:   "info",
	}, nil
}

// TestRedis handles redis test
func (ts *TestServiceImpl) TestRedis(c *gin.Context, req *dto.TestRedisRequest) (*vo.TestRedisResponse, error) {
	client := ts.redisManager.Client()

	ttl := time.Duration(req.TTL) * time.Second

	err := client.Set(c.Request.Context(), req.Key, req.Value, ttl).Err()
	if err != nil {
		return nil, err
	}

	val, err := client.Get(c.Request.Context(), req.Key).Result()
	if err != nil {
		return nil, err
	}

	if val != req.Value {
		return nil, errors.New("redis test failed: value mismatch")
	}

	return &vo.TestRedisResponse{
		Message: "Cache functionality test completed!",
		Key:     req.Key,
		Value:   val,
		TTL:     int(ttl.Seconds()),
	}, nil
}

// TestSuccess handles success test
func (ts *TestServiceImpl) TestSuccess(c *gin.Context) (*vo.TestSuccessResponse, error) {
	return &vo.TestSuccessResponse{
		Message: "Successful response validation passed!",
		Status:  "success",
	}, nil
}

// TestError handles error test
func (ts *TestServiceImpl) TestError(c *gin.Context) (*vo.TestErrorResponse, error) {
	return &vo.TestErrorResponse{
		Message: "Server exception",
		Code:    errno.ErrInternalServer,
	}, nil
}

// TestLong handles long request test
func (ts *TestServiceImpl) TestLong(c *gin.Context, req *dto.TestLongRequest) (*vo.TestLongResponse, error) {
	duration := time.Duration(req.Duration) * time.Second
	time.Sleep(duration)

	return &vo.TestLongResponse{
		Message:  "Simulated long-running request completed!",
		Duration: int(duration.Seconds()),
	}, nil
}

// TestI18n handles i18n test
func (ts *TestServiceImpl) TestI18n(c *gin.Context) (*vo.TestI18nResponse, error) {
	return &vo.TestI18nResponse{
		Message: "i18n test succeeded!",
	}, nil
}

// TestStream handles SSE related test
func (ts *TestServiceImpl) TestStream(c *gin.Context, req *dto.TestStreamRequest) (<-chan *sse.Event, error) {
	vars := map[string]any{
		"user_name":    req.Name,
		"greet_time":   time.Now().Format("2006-01-02 15:04:05"),
		"user_message": fmt.Sprintf("This is a message from %s", req.Name),
	}

	tmpl, err := ts.aiManager.GetTemplate(context.Background(), "example", &vars)
	if err != nil {
		ts.loggerManager.Logger().Errorf("failed to load prompt template 'example': %v", err)
		return nil, err
	}

	if len(tmpl.Messages) == 0 {
		return nil, errors.New("prompt template 'example' has no messages")
	}

	messages := make([]ai.Message, len(tmpl.Messages))
	for i, msg := range tmpl.Messages {
		messages[i] = ai.Message{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	events := make(chan *sse.Event, 100)

	go func() {
		defer close(events)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		heartbeatTicker := time.NewTicker(15 * time.Second)
		defer heartbeatTicker.Stop()

		stream, err := ts.aiManager.ChatStream(ctx, &ai.ChatRequest{
			Messages: messages,
		})
		if err != nil {
			ts.loggerManager.Logger().Errorf("failed to start AI chat stream: %v", err)
			return
		}

		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				case <-heartbeatTicker.C:
					events <- &sse.Event{Event: "heartbeat", Data: ""}
				}
			}
		}()

		for resp := range stream {
			if resp == nil {
				continue
			}

			payload, err := json.Marshal(resp)
			if err != nil {
				continue
			}

			events <- &sse.Event{Data: string(payload)}
		}
	}()

	return events, nil
}
