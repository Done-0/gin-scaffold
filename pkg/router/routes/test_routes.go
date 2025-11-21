// Package routes provides route registration functionality
// Author: Done-0
// Created: 2025-09-25
package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/Done-0/gin-scaffold/pkg/wire"
)

// RegisterTestRoutes registers test module routes
func RegisterTestRoutes(container *wire.Container, v1, v2 *gin.RouterGroup) {
	// V1 routes
	test := v1.Group("/test")
	{
		test.GET("/testPing", container.TestController.TestPing)
		test.GET("/testHello", container.TestController.TestHello)
		test.GET("/testLogger", container.TestController.TestLogger)
		test.POST("/testRedis", container.TestController.TestRedis)
		test.GET("/testSuccessRes", container.TestController.TestSuccess)
		test.GET("/testErrRes", container.TestController.TestError)
		test.GET("/testErrorMiddleware", container.TestController.TestErrorMiddleware)
		test.GET("/testI18n", container.TestController.TestI18n)
		test.POST("/testStream", container.TestController.TestStream)
	}

	// V2 routes
	v2Test := v2.Group("/test")
	{
		v2Test.POST("/testLongReq", container.TestController.TestLong)
	}
}
