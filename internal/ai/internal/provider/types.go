// Package provider implements AI provider interfaces
// Author: Done-0
// Created: 2025-08-31
package provider

import "context"

type Provider interface {
	Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error)
	ChatStream(ctx context.Context, req *ChatRequest) (<-chan *ChatStreamResponse, error)
}

// Request types
type ChatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Temperature float64   `json:"temperature,omitempty"`
}

type Message struct {
	Role             string `json:"role"`
	Content          string `json:"content"`
	ReasoningContent string `json:"reasoning_content,omitempty"`
}

// Response types
type ChatResponse struct {
	ID                string   `json:"id"`
	Object            string   `json:"object"`
	Created           int64    `json:"created"`
	Model             string   `json:"model"`
	Choices           []Choice `json:"choices"`
	Usage             Usage    `json:"usage"`
	SystemFingerprint string   `json:"system_fingerprint,omitempty"`
	Provider          string   `json:"provider"`
}

type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// Stream response types
type ChatStreamResponse struct {
	ID                string         `json:"id"`
	Object            string         `json:"object"`
	Created           int64          `json:"created"`
	Model             string         `json:"model"`
	Choices           []StreamChoice `json:"choices"`
	SystemFingerprint string         `json:"system_fingerprint,omitempty"`
	Usage             *Usage         `json:"usage,omitempty"`
	Provider          string         `json:"provider"`
}

type StreamChoice struct {
	Index        int          `json:"index"`
	Delta        MessageDelta `json:"delta"`
	FinishReason string       `json:"finish_reason"`
}

type MessageDelta struct {
	Role             string `json:"role,omitempty"`
	Content          string `json:"content,omitempty"`
	ReasoningContent string `json:"reasoning_content,omitempty"`
}
