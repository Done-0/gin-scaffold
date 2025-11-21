// Package internal provides error code registration internal implementation
// Author: Done-0
// Created: 2025-09-25
package internal

const (
	DefaultErrorMsg = "Service Internal Error" // Default error message
)

var (
	CodeDefinitions = make(map[int32]*CodeDefinition) // Error code definition mapping
)

// CodeDefinition error code definition
type CodeDefinition struct {
	Code    int32  // Error code
	Message string // Error message template
}

// RegisterOption registration option function
type RegisterOption func(definition *CodeDefinition)

// Register registers error code definition
func Register(code int32, msg string, opts ...RegisterOption) {
	definition := &CodeDefinition{
		Code:    code,
		Message: msg,
	}

	for _, opt := range opts {
		opt(definition)
	}

	CodeDefinitions[code] = definition
}
