// Package internal provides error status handling internal implementation
// Author: Done-0
// Created: 2025-09-25
package internal

import (
	"errors"
	"fmt"
	"strings"
)

// StatusError status error interface
type StatusError interface {
	error
	Code() int32 // Get status code
}

// statusError status error implementation
type statusError struct {
	code    int32             // Error code
	message string            // Error message
	extra   map[string]string // Extra information
	params  map[string]any    // Template parameters
}

// withStatus error wrapper with status
type withStatus struct {
	status *statusError // Status error
	cause  error        // Original error
	stack  string       // Stack information
}

// Code gets status code
func (w *statusError) Code() int32 {
	return w.code
}

// Msg gets error message
func (w *statusError) Msg() string {
	return w.message
}

// Error error string representation
func (w *statusError) Error() string {
	return fmt.Sprintf("code=%d message=%s", w.code, w.message)
}

// Extra gets extra information
func (w *statusError) Extra() map[string]string {
	return w.extra
}

// Params gets template parameters
func (w *statusError) Params() map[string]any {
	return w.params
}

// Code gets status code
func (w *withStatus) Code() int32 {
	return w.status.Code()
}

// Msg gets error message
func (w *withStatus) Msg() string {
	return w.status.Msg()
}

// Extra gets extra information
func (w *withStatus) Extra() map[string]string {
	return w.status.extra
}

// Params gets template parameters
func (w *withStatus) Params() map[string]any {
	return w.status.params
}

// Unwrap supports Go errors.Unwrap()
func (w *withStatus) Unwrap() error {
	return w.cause
}

// Is supports Go errors.Is()
func (w *withStatus) Is(target error) bool {
	var ws StatusError
	if errors.As(target, &ws) && w.status.Code() == ws.Code() {
		return true
	}
	return false
}

// As supports Go errors.As()
func (w *withStatus) As(target any) bool {
	return errors.As(w.status, target)
}

// StackTrace gets stack trace information
func (w *withStatus) StackTrace() string {
	return w.stack
}

// Error error string representation
func (w *withStatus) Error() string {
	b := strings.Builder{}
	b.WriteString(w.status.Error())

	if w.cause != nil {
		b.WriteString("\n")
		b.WriteString(fmt.Sprintf("cause=%s", w.cause))
	}

	if w.stack != "" {
		b.WriteString("\n")
		b.WriteString(fmt.Sprintf("stack=%s", w.stack))
	}

	return b.String()
}

// Option configuration option
type Option func(ws *withStatus)

// Param creates parameter replacement option
func Param(k, v string) Option {
	return func(ws *withStatus) {
		if ws == nil || ws.status == nil {
			return
		}
		if ws.status.params == nil {
			ws.status.params = make(map[string]any)
		}
		ws.status.params[k] = v
		ws.status.message = strings.ReplaceAll(ws.status.message, fmt.Sprintf("{{.%s}}", k), v)
	}
}

// Extra creates extra information option
func Extra(k, v string) Option {
	return func(ws *withStatus) {
		if ws == nil || ws.status == nil {
			return
		}
		if ws.status.extra == nil {
			ws.status.extra = make(map[string]string)
		}
		ws.status.extra[k] = v
	}
}

// NewByCode creates new error based on error code
func NewByCode(code int32, options ...Option) error {
	ws := &withStatus{
		status: getStatusByCode(code),
		cause:  nil,
		stack:  stack(),
	}

	for _, opt := range options {
		opt(ws)
	}

	return ws
}

// WrapByCode wraps existing error with error code
func WrapByCode(err error, code int32, options ...Option) error {
	if err == nil {
		return nil
	}

	ws := &withStatus{
		status: getStatusByCode(code),
		cause:  err,
	}

	for _, opt := range options {
		opt(ws)
	}

	// skip if stack has already exist
	var stackTracer StackTracer
	if errors.As(err, &stackTracer) {
		return ws
	}

	ws.stack = stack()

	return ws
}

// getStatusByCode gets status error based on error code
func getStatusByCode(code int32) *statusError {
	codeDefinition, ok := CodeDefinitions[code]
	if ok {
		// predefined err code
		return &statusError{
			code:    code,
			message: codeDefinition.Message,
			extra:   make(map[string]string),
			params:  make(map[string]any),
		}
	}

	return &statusError{
		code:    code,
		message: DefaultErrorMsg,
		extra:   make(map[string]string),
		params:  make(map[string]any),
	}
}
