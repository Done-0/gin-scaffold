// Package middleware 提供中间件集成和初始化功能
// 创建者：Done-0
// 创建时间：2025-07-01
package middleware

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"

	"github.com/Done-0/gin-scaffold/internal/middleware/cors"
	"github.com/Done-0/gin-scaffold/internal/middleware/recovery"
	"github.com/Done-0/gin-scaffold/internal/middleware/secure"
)

// New 初始化并注册所有中间件
// 参数：
//
//	app: Gin 实例
func New(app *gin.Engine) {
	// 全局请求 ID 中间件
	app.Use(requestid.New())

	// 日志中间件
	app.Use(gin.Logger())

	// COR 跨域中间件
	app.Use(cors.New())

	// Gzip 压缩中间件
	app.Use(gzip.Gzip(gzip.DefaultCompression))

	// 安全中间件
	app.Use(secure.New())

	// Recovery 中间件
	app.Use(recovery.New())
}
