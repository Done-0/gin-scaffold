// Package internal provides error stack trace handling internal implementation
// Author: Done-0
// Created: 2025-09-25
package internal

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

// StackTracer stack trace interface
type StackTracer interface {
	StackTrace() string // Get stack trace information
}

// withStack error wrapper with stack information
type withStack struct {
	cause error  // Cause error
	stack string // Stack information
}

// Unwrap returns the wrapped original error
func (w *withStack) Unwrap() error {
	return w.cause
}

// StackTrace gets stack trace information
func (w *withStack) StackTrace() string {
	return w.stack
}

// Error error string representation
func (w *withStack) Error() string {
	return fmt.Sprintf("%s\nstack=%s", w.cause.Error(), w.stack)
}

// stack gets current call stack information
func stack() string {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(2, pcs[:])

	b := strings.Builder{}
	for i := 0; i < n; i++ {
		fn := runtime.FuncForPC(pcs[i])

		file, line := fn.FileLine(pcs[i])
		name := trimPathPrefix(fn.Name())
		b.WriteString(fmt.Sprintf("%s:%d %s\n", file, line, name))
	}

	return b.String()
}

// trimPathPrefix removes path prefix, keeping only function name
func trimPathPrefix(s string) string {
	i := strings.LastIndex(s, "/")
	s = s[i+1:]
	i = strings.Index(s, ".")
	return s[i+1:]
}

// withStackTraceIfNotExists adds stack information if error doesn't have it
func withStackTraceIfNotExists(err error) error {
	if err == nil {
		return nil
	}

	// skip if stack has already exist
	var stackTracer StackTracer
	if errors.As(err, &stackTracer) {
		return err
	}

	return &withStack{
		err,
		stack(),
	}
}
