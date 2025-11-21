// Package vo provides common value objects
// Author: Done-0
// Created: 2025-09-25
package vo

import (
	"time"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"

	"github.com/Done-0/gin-scaffold/internal/types/errno"
	"github.com/Done-0/gin-scaffold/internal/utils/errorx"
)

// Result common API response structure
type Result struct {
	Code      int    `json:"code,omitempty"`
	Msg       string `json:"msg,omitempty"`
	Data      any    `json:"data,omitempty"`
	RequestId string `json:"requestId"`
	TimeStamp int64  `json:"timeStamp"`
}

// Success successful response
func Success(c *gin.Context, data any, messageKey ...string) Result {
	if errData, ok := data.(error); ok {
		data = errData.Error()
	}
	return Result{
		Data:      data,
		RequestId: requestid.Get(c),
		TimeStamp: time.Now().Unix(),
	}
}

// Fail creates error response
func Fail(c *gin.Context, data any, err error) Result {
	if errData, ok := data.(error); ok {
		data = errData.Error()
	}

	var code int
	var message string

	switch e := err.(type) {
	case errorx.StatusError:
		code = int(e.Code())
		message = err.Error()
	default:
		code = errno.ErrInternalServer
		message = "Internal server error"
	}

	return Result{
		Code:      code,
		Msg:       message,
		Data:      data,
		RequestId: requestid.Get(c),
		TimeStamp: time.Now().Unix(),
	}
}
