// Package dto provides test-related data transfer object definitions
// Author: Done-0
// Created: 2025-09-25
package dto

// TestRedisRequest redis test request
type TestRedisRequest struct {
	Key   string `json:"key" validate:"required"`
	Value string `json:"value" validate:"required"`
	TTL   int    `json:"ttl" validate:"omitempty,min=1"`
}

// TestLongRequest long request test
type TestLongRequest struct {
	Duration int `json:"duration" validate:"omitempty,min=1,max=10"`
}

// TestStreamRequest SSE 流式测试请求体
type TestStreamRequest struct {
	Name string `json:"name" validate:"required"`
}
