// Package file provides file loading and saving utilities
// Author: Done-0
// Created: 2025-08-31
package file

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// LoadJSONFile loads and validates a JSON file
func LoadJSONFile(filePath string, target any) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	if err := json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("invalid JSON format: %w", err)
	}

	return nil
}

// SaveJSONFile saves data to a JSON file with proper formatting
func SaveJSONFile(filePath string, data any) error {
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// GetFileNameWithoutExt returns filename without extension
func GetFileNameWithoutExt(filePath string) string {
	return strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))
}
