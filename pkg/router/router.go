// Package router 提供应用程序路由注册功能
// 创建者：Done-0
// 创建时间：2025-07-01
package router

import (
	"github.com/gin-gonic/gin"

	"github.com/Done-0/gin-scaffold/pkg/router/routes"
)

// New @title		Metaphysics API
// @version			1.0
// @description		This is the API documentation for Metaphysics.
// @host			localhost:8080
// @BasePath		/
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description 输入格式: Bearer {token}
// New 函数用于注册应用程序的路由
// 参数：
//   - app: Gin 实例
func New(app *gin.Engine) {
	// 创建多版本 API 路由组
	api1 := app.Group("/api/v1")
	api2 := app.Group("/api/v2")

	// 注册测试相关的路由
	routes.RegisterTestRoutes(api1, api2)
}
