// Package errno provides system-level error code definitions
// Author: Done-0
// Created: 2025-09-25
package errno

import (
	"github.com/Done-0/gin-scaffold/internal/utils/errorx/code"
)

// System-level error codes: 10000 ~ 19999
// Used: 10001-10008
// Next available: 10009
const (
	ErrInternalServer     = 10001 // Internal server error
	ErrInvalidParams      = 10002 // Parameter validation failed
	ErrUnauthorized       = 10003 // Authentication failed
	ErrForbidden          = 10004 // Insufficient permissions
	ErrResourceNotFound   = 10005 // Resource not found
	ErrResourceConflict   = 10006 // Resource conflict
	ErrTooManyRequests    = 10007 // Request rate limit exceeded
	ErrServiceUnavailable = 10008 // Service unavailable
)

func init() {
	code.Register(ErrInternalServer, "internal server error: {{.msg}}")
	code.Register(ErrInvalidParams, "invalid parameter: {{.msg}}")
	code.Register(ErrUnauthorized, "unauthorized access: {{.msg}}")
	code.Register(ErrForbidden, "permission denied: {{.resource}}")
	code.Register(ErrResourceNotFound, "{{.resource}} not found: {{.id}}")
	code.Register(ErrResourceConflict, "{{.resource}} already exists: {{.id}}")
	code.Register(ErrTooManyRequests, "too many requests: {{.limit}} per {{.period}}")
	code.Register(ErrServiceUnavailable, "service unavailable: {{.service}}")
}
