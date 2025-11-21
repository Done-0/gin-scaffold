// Package model provides database model definitions and management
// Author: Done-0
// Created: 2025-08-24
package model

import (
	"github.com/Done-0/gin-scaffold/internal/model/user"
)

// GetAllModels gets and registers all models for database migration
func GetAllModels() []any {
	return []any{
		&user.User{}, // User model
	}
}
