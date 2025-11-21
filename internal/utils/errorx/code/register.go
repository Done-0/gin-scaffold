// Package code provides error code registration adapter
// Author: Done-0
// Created: 2025-09-25
package code

import (
	"github.com/Done-0/gin-scaffold/internal/utils/errorx"
)

// RegisterOptionFn registration option function type
type RegisterOptionFn = errorx.RegisterOption

// Register registers predefined error code information
func Register(code int32, msg string, opts ...RegisterOptionFn) {
	errorx.Register(code, msg, opts...)
}
