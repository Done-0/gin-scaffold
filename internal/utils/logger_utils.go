// Package utils 提供日志记录工具
// 创建者：Done-0
// 创建时间：2025-07-01
package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/Done-0/gin-scaffold/internal/global"
)

const (
	BIZLOG = "Bizlog" // 业务日志键名
)

// BizLogger 业务日志记录器
// 参数：
//
//	c: Gin 上下文
//
// 返回值：
//
//	*logrus.Entry: 日志条目
func BizLogger(c *gin.Context) *logrus.Entry {
	if bizLog, ok := c.Get(BIZLOG); ok {
		if entry, ok := bizLog.(*logrus.Entry); ok {
			return entry
		}
	}

	return logrus.NewEntry(global.SysLog)
}
