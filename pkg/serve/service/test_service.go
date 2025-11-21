// Package service provides test service interfaces
// Author: Done-0
// Created: 2025-09-25
package service

import (
	"github.com/gin-contrib/sse"
	"github.com/gin-gonic/gin"

	"github.com/Done-0/gin-scaffold/pkg/serve/controller/dto"
	"github.com/Done-0/gin-scaffold/pkg/vo"
)

// TestService test service interface
type TestService interface {
	TestPing(c *gin.Context) (*vo.TestPingResponse, error)
	TestHello(c *gin.Context) (*vo.TestHelloResponse, error)
	TestLogger(c *gin.Context) (*vo.TestLoggerResponse, error)
	TestRedis(c *gin.Context, req *dto.TestRedisRequest) (*vo.TestRedisResponse, error)
	TestSuccess(c *gin.Context) (*vo.TestSuccessResponse, error)
	TestError(c *gin.Context) (*vo.TestErrorResponse, error)
	TestLong(c *gin.Context, req *dto.TestLongRequest) (*vo.TestLongResponse, error)
	TestI18n(c *gin.Context) (*vo.TestI18nResponse, error)
	TestStream(c *gin.Context, req *dto.TestStreamRequest) (<-chan *sse.Event, error)
}
