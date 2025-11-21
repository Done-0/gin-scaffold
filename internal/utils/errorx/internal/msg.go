// Package internal provides error message wrapping internal implementation
// Author: Done-0
// Created: 2025-09-25
package internal

import (
	"fmt"
)

// withMessage error wrapper with message
type withMessage struct {
	cause error  // Cause error
	msg   string // Message
}

// Unwrap returns the wrapped original error
func (w *withMessage) Unwrap() error {
	return w.cause
}

// Error error string representation
func (w *withMessage) Error() string {
	return fmt.Sprintf("%s\ncause=%s", w.msg, w.cause.Error())
}

// wrapf wraps error with formatted message (internal function)
func wrapf(err error, format string, args ...any) error {
	if err == nil {
		return nil
	}
	err = &withMessage{
		cause: err,
		msg:   fmt.Sprintf(format, args...),
	}

	return err
}

// Wrapf wraps error with formatted message and adds stack information
func Wrapf(err error, format string, args ...any) error {
	return withStackTraceIfNotExists(wrapf(err, format, args...))
}