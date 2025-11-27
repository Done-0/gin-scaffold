// Package vo provides common value objects
// Author: Done-0
// Created: 2025-09-25
package vo

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"

	"github.com/Done-0/gin-scaffold/internal/types/errno"
	"github.com/Done-0/gin-scaffold/internal/utils/errorx"
	"github.com/Done-0/gin-scaffold/internal/utils/i18n"
)

// Error information
type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Result common API response structure
type Result struct {
	Error     *Error `json:"error,omitempty"`
	Data      any    `json:"data,omitempty"`
	RequestId string `json:"requestId"`
	TimeStamp int64  `json:"timeStamp"`
}

// Success successful response
func Success(c *gin.Context, data any) Result {
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

	var code, message string

	switch e := err.(type) {
	case errorx.StatusError:
		code = strconv.Itoa(int(e.Code()))
		params := e.Params()

		if len(params) == 0 {
			message = i18n.T(c, code)
		} else {
			args := make([]string, len(params)*2)
			i := 0
			for k, v := range params {
				args[i] = k
				if s, ok := v.(string); ok {
					args[i+1] = s
				} else {
					args[i+1] = fmt.Sprintf("%v", v)
				}
				i += 2
			}
			message = i18n.T(c, code, args...)
		}
	default:
		code = strconv.Itoa(errno.ErrInternalServer)
		message = i18n.T(c, code, "msg", err.Error())
	}

	return Result{
		Error:     &Error{Code: code, Message: message},
		Data:      data,
		RequestId: requestid.Get(c),
		TimeStamp: time.Now().Unix(),
	}
}
