// Package base provides base model definitions and common database operations
// Author: Done-0
// Created: 2025-09-25
package base

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/Done-0/gin-scaffold/internal/utils/snowflake"
)

// Base contains common model fields
type Base struct {
	ID        int64   `gorm:"primaryKey;type:bigint" json:"id"`          // Primary key (snowflake)
	CreatedAt int64   `gorm:"type:bigint" json:"created_at"`             // Creation timestamp
	UpdatedAt int64   `gorm:"type:bigint" json:"updated_at"`             // Update timestamp
	Ext       JSONMap `gorm:"type:json" json:"ext"`                      // Extension fields
	Deleted   bool    `gorm:"type:boolean;default:false" json:"deleted"` // Soft delete flag
}

// JSONMap handles JSON type fields
type JSONMap map[string]any

// Scan implements sql.Scanner interface
func (j *JSONMap) Scan(value any) error {
	switch v := value.(type) {
	case nil:
		*j = make(JSONMap)
	case []byte:
		return json.Unmarshal(v, j)
	default:
		return errors.New("cannot scan into JSONMap")
	}
	return nil
}

// Value implements driver.Valuer interface
func (j JSONMap) Value() (driver.Value, error) {
	return json.Marshal(j)
}

// BeforeCreate implements GORM hook
func (m *Base) BeforeCreate(db *gorm.DB) error {
	if m.ID != 0 {
		return nil
	}

	now := time.Now().Unix()
	m.CreatedAt = now
	m.UpdatedAt = now

	var err error
	m.ID, err = snowflake.GenerateID()
	return err
}

// BeforeUpdate implements GORM hook
func (m *Base) BeforeUpdate(db *gorm.DB) error {
	m.UpdatedAt = time.Now().Unix()
	return nil
}
