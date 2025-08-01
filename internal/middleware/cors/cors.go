// Package cors 提供CORS跨域中间件配置
// 创建者：Done-0
// 创建时间：2025-07-01
package cors

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// New 创建CORS中间件
// 返回值：
//
//	gin.HandlerFunc: CORS中间件
func New() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins: []string{"*"}, // 生产环境应指定具体域名，如: []string{"https://example.com"}
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders: []string{
			"Origin",           // 请求源
			"Content-Type",     // 内容类型
			"Accept",           // 接受类型
			"Authorization",    // 认证头
			"Refresh-Token",    // 刷新令牌头
			"X-Requested-With", // AJAX标识
		},
		ExposeHeaders: []string{
			"Content-Length",
			"Authorization",
			"Refresh-Token",
		},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}
