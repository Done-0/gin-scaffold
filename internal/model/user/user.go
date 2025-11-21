// Package user provides user data model definitions
// Author: Done-0
// Created: 2025-09-25
package user

import "github.com/Done-0/gin-scaffold/internal/model/base"

// User represents user model
type User struct {
	base.Base
	Email    string `gorm:"type:varchar(64);unique;not null" json:"email"`    // Email, primary login method
	Password string `gorm:"type:varchar(255);not null" json:"password"`       // Encrypted password
	Nickname string `gorm:"type:varchar(64);unique;not null" json:"nickname"` // User nickname
	Avatar   string `gorm:"type:varchar(255);default:null" json:"avatar"`     // User avatar
	Role     string `gorm:"type:varchar(32);default:'user'" json:"role"`      // User role
}

// TableName specifies table name
func (User) TableName() string {
	return "users"
}
