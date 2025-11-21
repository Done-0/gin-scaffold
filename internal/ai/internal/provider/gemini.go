// Package provider implements AI provider interfaces
// Author: Done-0
// Created: 2025-08-31
package provider

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"sync/atomic"
	"time"

	"golang.org/x/time/rate"
	"google.golang.org/genai"

	"github.com/Done-0/gin-scaffold/configs"

	rateUtil "github.com/Done-0/gin-scaffold/internal/utils/rate"
)

type geminiProvider struct {
	config       *configs.ProviderInstanceConfig
	rateLimiter  *rate.Limiter
	keyCounter   *uint64
	modelCounter *uint64
}

func NewGemini(config *configs.ProviderInstanceConfig, keyCounter *uint64, modelCounter *uint64) (Provider, error) {
	rateLimit, burst, err := rateUtil.ParseLimit(config.RateLimit)
	if err != nil {
		return nil, err
	}
	limiter := rate.NewLimiter(rateLimit, burst)

	return &geminiProvider{
		config:       config,
		rateLimiter:  limiter,
		keyCounter:   keyCounter,
		modelCounter: modelCounter,
	}, nil
}

func (p *geminiProvider) Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	model := req.Model
	if model == "" {
		modelIndex := atomic.AddUint64(p.modelCounter, 1) - 1
		model = p.config.Models[modelIndex%uint64(len(p.config.Models))]
	}

	keyIndex := atomic.AddUint64(p.keyCounter, 1) - 1
	apiKey := p.config.Keys[keyIndex%uint64(len(p.config.Keys))]

	if err := p.rateLimiter.Wait(ctx); err != nil {
		return nil, err
	}

	clientConfig := &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	}

	client, err := genai.NewClient(ctx, clientConfig)
	if err != nil {
		return nil, err
	}

	var prompt string
	for _, msg := range req.Messages {
		prompt += msg.Content + "\n"
	}

	var resp *genai.GenerateContentResponse
	for attempt := 0; attempt <= p.config.MaxRetries; attempt++ {
		resp, err = client.Models.GenerateContent(
			ctx,
			model,
			genai.Text(prompt),
			&genai.GenerateContentConfig{
				Temperature:     &p.config.Temperature,
				MaxOutputTokens: int32(p.config.MaxTokens),
				TopP:            &p.config.TopP,
				TopK:            func() *float32 { v := float32(p.config.TopK); return &v }(),
				ThinkingConfig: &genai.ThinkingConfig{
					IncludeThoughts: true,
				},
			},
		)
		if err == nil {
			break
		}
		if attempt < p.config.MaxRetries {
			time.Sleep(time.Duration(attempt+1) * time.Second)
		}
	}
	if err != nil {
		return nil, err
	}

	var content, reasoningContent string
	if len(resp.Candidates) > 0 {
		candidate := resp.Candidates[0]
		if candidate.Content != nil {
			for _, part := range candidate.Content.Parts {
				if len(part.Text) == 0 {
					continue
				}
				if part.Thought {
					reasoningContent += part.Text
				} else {
					content += part.Text
				}
			}
		}
	}

	now := time.Now()
	bytes := make([]byte, 16)
	rand.Read(bytes)
	chatResp := &ChatResponse{
		ID:      "chatcmpl-" + hex.EncodeToString(bytes),
		Object:  "chat.completion",
		Created: now.Unix(),
		Model:   model,
		Choices: []Choice{{
			Index: 0,
			Message: Message{
				Role:             "assistant",
				Content:          content,
				ReasoningContent: reasoningContent,
			},
			FinishReason: string(resp.Candidates[0].FinishReason),
		}},
		SystemFingerprint: "",
		Provider:          "gemini",
	}

	if resp.Candidates[0].FinishReason != "" {
		chatResp.Usage = Usage{
			PromptTokens:     int(resp.UsageMetadata.PromptTokenCount),
			CompletionTokens: int(resp.UsageMetadata.CandidatesTokenCount),
			TotalTokens:      int(resp.UsageMetadata.TotalTokenCount),
		}
	}

	return chatResp, nil
}

func (p *geminiProvider) ChatStream(ctx context.Context, req *ChatRequest) (<-chan *ChatStreamResponse, error) {
	model := req.Model
	if model == "" {
		modelIndex := atomic.AddUint64(p.modelCounter, 1) - 1
		model = p.config.Models[modelIndex%uint64(len(p.config.Models))]
	}

	keyIndex := atomic.AddUint64(p.keyCounter, 1) - 1
	apiKey := p.config.Keys[keyIndex%uint64(len(p.config.Keys))]

	if err := p.rateLimiter.Wait(ctx); err != nil {
		return nil, err
	}

	clientConfig := &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	}

	contents := make([]*genai.Content, len(req.Messages))
	for i, msg := range req.Messages {
		role := genai.RoleUser
		if msg.Role == "assistant" {
			role = genai.RoleModel
		}
		contents[i] = &genai.Content{
			Parts: []*genai.Part{{Text: msg.Content}},
			Role:  role,
		}
	}

	var client *genai.Client
	var err error
	for attempt := 0; attempt <= p.config.MaxRetries; attempt++ {
		client, err = genai.NewClient(ctx, clientConfig)
		if err == nil {
			break
		}
		if attempt < p.config.MaxRetries {
			time.Sleep(time.Duration(attempt+1) * time.Second)
		}
	}
	if err != nil {
		return nil, err
	}

	stream := client.Models.GenerateContentStream(
		ctx,
		model,
		contents,
		&genai.GenerateContentConfig{
			Temperature:     &p.config.Temperature,
			MaxOutputTokens: int32(p.config.MaxTokens),
			TopP:            &p.config.TopP,
			TopK:            func() *float32 { v := float32(p.config.TopK); return &v }(),
			ThinkingConfig: &genai.ThinkingConfig{
				IncludeThoughts: true,
			},
		},
	)

	ch := make(chan *ChatStreamResponse)
	go func() {
		defer close(ch)

		if stream == nil {
			return
		}

		for chunk := range stream {
			if chunk == nil {
				continue
			}

			if len(chunk.Candidates) == 0 {
				continue
			}

			candidate := chunk.Candidates[0]
			if candidate.Content == nil || len(candidate.Content.Parts) == 0 {
				continue
			}

			var normalContent, thinkingContent string

			for _, part := range candidate.Content.Parts {
				if part.Text == "" {
					continue
				}

				if part.Thought {
					thinkingContent += part.Text
				} else {
					normalContent += part.Text
				}
			}

			now := time.Now()

			bytes := make([]byte, 16)
			rand.Read(bytes)
			streamResp := &ChatStreamResponse{
				ID:      "chatcmpl-" + hex.EncodeToString(bytes),
				Object:  "chat.completion.chunk",
				Created: now.Unix(),
				Model:   model,
				Choices: []StreamChoice{{
					Index: 0,
					Delta: MessageDelta{
						Role:             "assistant",
						Content:          normalContent,
						ReasoningContent: thinkingContent,
					},
					FinishReason: string(candidate.FinishReason),
				}},
				SystemFingerprint: "",
				Provider:          "gemini",
			}

			if candidate.FinishReason != "" {
				streamResp.Usage = &Usage{
					PromptTokens:     int(chunk.UsageMetadata.PromptTokenCount),
					CompletionTokens: int(chunk.UsageMetadata.CandidatesTokenCount),
					TotalTokens:      int(chunk.UsageMetadata.TotalTokenCount),
				}
			}

			select {
			case ch <- streamResp:
			case <-ctx.Done():
				return
			}

			if candidate.FinishReason == genai.FinishReasonStop || candidate.FinishReason == genai.FinishReasonMaxTokens {
				break
			}
		}
	}()
	return ch, nil
}
