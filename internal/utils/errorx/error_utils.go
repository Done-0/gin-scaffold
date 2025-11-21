// Package errorx provides error handling utilities with status codes and stack traces
// Author: Done-0
// Created: 2025-09-25
package errorx

import (
	"github.com/Done-0/gin-scaffold/internal/utils/errorx/internal"
)

// StatusError error interface with status code
type StatusError interface {
	error
	Code() int32              // Get error code
	Msg() string              // Get error message
	Extra() map[string]string // Get extra information
	Params() map[string]any   // Get template parameters
}

// Option StatusError configuration option
type Option = internal.Option

// KV creates key-value parameter option
func KV(k, v string) Option {
	return internal.Param(k, v)
}

// New creates new error based on error code
func New(code int32, options ...Option) error {
	return internal.NewByCode(code, options...)
}

// Register registers error code definition
func Register(code int32, msg string, opts ...internal.RegisterOption) {
	internal.Register(code, msg, opts...)
}

// RegisterOption registration option type
type RegisterOption = internal.RegisterOption
