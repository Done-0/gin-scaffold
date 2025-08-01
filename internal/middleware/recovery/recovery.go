// Package recovery 提供全局异常恢复中间件
// 创建者：Done-0
// 创建时间：2025-07-01
package recovery

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"

	"github.com/Done-0/gin-scaffold/internal/global"
	"github.com/Done-0/gin-scaffold/pkg/vo"
)

// New 返回 Gin 框架的全局异常恢复中间件
// 返回值：
//
//	gin.HandlerFunc: 全局异常恢复中间件
func New() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				// 获取堆栈信息
				stackSize := 4096
				var stack []byte
				for {
					stack = make([]byte, stackSize)
					n := runtime.Stack(stack, false)
					if n < stackSize {
						stack = stack[:n]
						break
					}
					stackSize *= 2
				}

				// 记录日志
				global.SysLog.WithFields(map[string]any{
					"stack_trace": string(stack),
				}).Errorf("发生运行时异常: %v", rec)

				var err error
				if e, ok := rec.(error); ok {
					err = e
				} else {
					err = fmt.Errorf("%v", rec)
				}
				c.AbortWithStatusJSON(http.StatusInternalServerError, vo.Fail(c, nil, err))
			}
		}()
		c.Next()
	}
}
