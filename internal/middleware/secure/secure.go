// Package secure 提供安全中间件配置
// 创建者：Done-0
// 创建时间：2025-07-01
package secure

import (
	"github.com/gin-contrib/secure"
	"github.com/gin-gonic/gin"
)

// New 创建安全中间件
// 返回值：
//
//	gin.HandlerFunc: 安全中间件
func New() gin.HandlerFunc {
	return secure.New(secure.Config{
		FrameDeny:          true, // 禁止嵌入 frame
		ContentTypeNosniff: true, // 禁止 MIME 类型嗅探
		BrowserXssFilter:   true, // 启用 XSS 过滤
		ContentSecurityPolicy: "default-src 'self'; " + // 默认只允许同源资源
			"script-src 'self' 'unsafe-inline' 'unsafe-eval'; " + // 脚本：同源 + 内联 + eval
			"style-src 'self' 'unsafe-inline' https://fonts.googleapis.com; " + // 样式：同源 + 内联 + Google字体
			"img-src 'self' data: https: http:; " + // 图片：同源 + data URL + 所有 HTTPS/HTTP
			"font-src 'self' https://fonts.gstatic.com https://fonts.googleapis.com; " + // 字体：同源 + Google字体
			"connect-src 'self' https: wss:; " + // 连接：同源 + HTTPS API + WebSocket
			"media-src 'self' https: data:; " + // 媒体：同源 + HTTPS + data URL
			"frame-src 'self' https:; " + // 框架：同源 + HTTPS（支持第三方内嵌如地图、视频）
			"child-src 'self' https:; " + // 子资源：同源 + HTTPS
			"worker-src 'self'; " + // Web Worker：同源
			"frame-ancestors 'self'; " + // 允许同源嵌入
			"base-uri 'self'; " + // 限制 base 标签
			"object-src 'none'; " + // 禁用危险的 object/embed
			"upgrade-insecure-requests", // 自动将 HTTP 请求升级为 HTTPS

		ReferrerPolicy: "strict-origin-when-cross-origin",
	})
}
