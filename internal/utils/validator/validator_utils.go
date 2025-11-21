// Package validator provides parameter validation utilities
// Author: Done-0
// Created: 2025-08-24
package validator

import "github.com/go-playground/validator/v10"

// ValidErrRes validation error result struct
type ValidErrRes struct {
	Error bool   // Whether error exists
	Field string // Error field name
	Tag   string // Error tag
	Value any    // Error value
}

// NewValidator global validator instance
var NewValidator = validator.New()

// Validate parameter validator
func Validate(data any) []ValidErrRes {
	var Errors []ValidErrRes
	errs := NewValidator.Struct(data)
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			var el ValidErrRes
			el.Error = true
			el.Field = err.Field()
			el.Tag = err.Tag()
			el.Value = err.Value()

			Errors = append(Errors, el)
		}
	}
	return Errors
}
