// Package vo provides test-related value object definitions
// Author: Done-0
// Created: 2025-09-25
package vo

// TestPingResponse ping test response
type TestPingResponse struct {
	Time    string `json:"time"`
	Message string `json:"message"`
}

// TestHelloResponse hello test response
type TestHelloResponse struct {
	Version string `json:"version"`
	Message string `json:"message"`
}

// TestLoggerResponse logger test response
type TestLoggerResponse struct {
	Level   string `json:"level"`
	Message string `json:"message"`
}

// TestRedisResponse redis test response
type TestRedisResponse struct {
	Key     string `json:"key"`
	Value   string `json:"value"`
	TTL     int    `json:"ttl"`
	Message string `json:"message"`
}

// TestSuccessResponse success test response
type TestSuccessResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// TestErrorResponse error test response
type TestErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// TestLongResponse long request test response
type TestLongResponse struct {
	Duration int    `json:"duration"`
	Message  string `json:"message"`
}

// TestI18nResponse i18n test response
type TestI18nResponse struct {
	Message string `json:"message"`
}
