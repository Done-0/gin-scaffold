// Package vo 提供视图对象定义和响应结果包装
// 创建者：Done-0
// 创建时间：2025-07-01
package vo

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"

	bizErr "github.com/Done-0/gin-scaffold/internal/error"
)

// Result 通用 API 响应结果结构体
type Result struct {
	*bizErr.Err     // 错误信息
	Data        any `json:"data"`      // 响应数据
	RequestId   any `json:"requestId"` // 请求ID
	TimeStamp   any `json:"timeStamp"` // 响应时间戳
}

// Success 成功返回
// 参数：
//
//	c: Gin 上下文
//	data: 响应数据
//
// 返回值：
//
//	Result: 成功响应结果
func Success(c *gin.Context, data any) Result {
	return Result{
		Err:       nil,
		Data:      data,
		RequestId: c.GetHeader("X-Request-ID"),
		TimeStamp: time.Now().Unix(),
	}
}

// Fail 失败返回
// 参数：
//
//	c: Gin 上下文
//	data: 错误相关数据
//	err: 错误对象
//
// 返回值：
//
//	Result: 失败响应结果
func Fail(c *gin.Context, data any, err error) Result {
	var newBizErr *bizErr.Err
	if ok := errors.As(err, &newBizErr); ok {
		return Result{
			Err:       newBizErr,
			Data:      data,
			RequestId: c.GetHeader("X-Request-ID"),
			TimeStamp: time.Now().Unix(),
		}
	}

	return Result{
		Err:       bizErr.New(bizErr.SYSTEM_ERROR),
		Data:      data,
		RequestId: c.GetHeader("X-Request-ID"),
		TimeStamp: time.Now().Unix(),
	}
}
