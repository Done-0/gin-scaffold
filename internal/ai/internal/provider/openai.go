// Package provider implements AI provider interfaces
// Author: Done-0
// Created: 2025-08-31
package provider

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"io"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/sashabaranov/go-openai"
	"golang.org/x/time/rate"

	"github.com/Done-0/gin-scaffold/configs"

	rateUtil "github.com/Done-0/gin-scaffold/internal/utils/rate"
)

type openAIProvider struct {
	config       *configs.ProviderInstanceConfig
	rateLimiter  *rate.Limiter
	keyCounter   *uint64
	modelCounter *uint64
}

func NewOpenAI(config *configs.ProviderInstanceConfig, keyCounter *uint64, modelCounter *uint64) (Provider, error) {
	rateLimit, burst, err := rateUtil.ParseLimit(config.RateLimit)
	if err != nil {
		return nil, err
	}
	limiter := rate.NewLimiter(rateLimit, burst)

	return &openAIProvider{
		config:       config,
		rateLimiter:  limiter,
		keyCounter:   keyCounter,
		modelCounter: modelCounter,
	}, nil
}

func (p *openAIProvider) Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	model := req.Model
	if model == "" {
		modelIndex := atomic.AddUint64(p.modelCounter, 1) - 1
		model = p.config.Models[modelIndex%uint64(len(p.config.Models))]
	}

	keyIndex := atomic.AddUint64(p.keyCounter, 1) - 1
	apiKey := p.config.Keys[keyIndex%uint64(len(p.config.Keys))]

	config := openai.DefaultConfig(apiKey)
	config.BaseURL = p.config.BaseURL
	config.HTTPClient = &http.Client{
		Timeout: time.Duration(p.config.Timeout) * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 20,
			MaxConnsPerHost:     100,
			IdleConnTimeout:     90 * time.Second,
			DisableKeepAlives:   false,
		},
	}
	client := openai.NewClientWithConfig(config)

	messages := make([]openai.ChatCompletionMessage, len(req.Messages))
	for i, msg := range req.Messages {
		messages[i] = openai.ChatCompletionMessage{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	if err := p.rateLimiter.Wait(ctx); err != nil {
		return nil, err
	}

	request := openai.ChatCompletionRequest{
		Model:       model,
		Messages:    messages,
		MaxTokens:   p.config.MaxTokens,
		Temperature: p.config.Temperature,
		TopP:        p.config.TopP,
	}

	var resp openai.ChatCompletionResponse
	var err error
	for attempt := 0; attempt <= p.config.MaxRetries; attempt++ {
		resp, err = client.CreateChatCompletion(ctx, request)
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

	bytes := make([]byte, 16)
	rand.Read(bytes)
	chatResp := &ChatResponse{
		ID:      "chatcmpl-" + hex.EncodeToString(bytes),
		Object:  resp.Object,
		Created: resp.Created,
		Model:   resp.Model,
		Choices: []Choice{{
			Index: resp.Choices[0].Index,
			Message: Message{
				Role:             resp.Choices[0].Message.Role,
				Content:          resp.Choices[0].Message.Content,
				ReasoningContent: resp.Choices[0].Message.ReasoningContent,
			},
			FinishReason: string(resp.Choices[0].FinishReason),
		}},
		SystemFingerprint: resp.SystemFingerprint,
		Provider:          "openai",
	}

	if string(resp.Choices[0].FinishReason) != "" {
		chatResp.Usage = Usage{
			PromptTokens:     resp.Usage.PromptTokens,
			CompletionTokens: resp.Usage.CompletionTokens,
			TotalTokens:      resp.Usage.TotalTokens,
		}
	}

	return chatResp, nil
}

func (p *openAIProvider) ChatStream(ctx context.Context, req *ChatRequest) (<-chan *ChatStreamResponse, error) {
	model := req.Model
	if model == "" {
		modelIndex := atomic.AddUint64(p.modelCounter, 1) - 1
		model = p.config.Models[modelIndex%uint64(len(p.config.Models))]
	}

	keyIndex := atomic.AddUint64(p.keyCounter, 1) - 1
	apiKey := p.config.Keys[keyIndex%uint64(len(p.config.Keys))]

	config := openai.DefaultConfig(apiKey)
	config.BaseURL = p.config.BaseURL
	config.HTTPClient = &http.Client{
		Timeout: time.Duration(p.config.Timeout) * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 20,
			MaxConnsPerHost:     100,
			IdleConnTimeout:     90 * time.Second,
			DisableKeepAlives:   false,
		},
	}
	client := openai.NewClientWithConfig(config)

	if err := p.rateLimiter.Wait(ctx); err != nil {
		return nil, err
	}

	messages := make([]openai.ChatCompletionMessage, len(req.Messages))
	for i, msg := range req.Messages {
		messages[i] = openai.ChatCompletionMessage{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	request := openai.ChatCompletionRequest{
		Model:       model,
		Messages:    messages,
		Stream:      true,
		MaxTokens:   p.config.MaxTokens,
		Temperature: float32(p.config.Temperature),
		TopP:        float32(p.config.TopP),
	}

	var stream *openai.ChatCompletionStream
	var err error
	for attempt := 0; attempt <= p.config.MaxRetries; attempt++ {
		stream, err = client.CreateChatCompletionStream(ctx, request)
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

	ch := make(chan *ChatStreamResponse)
	go func() {
		defer close(ch)
		defer stream.Close()

		for {
			response, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					break
				}
				break
			}

			bytes := make([]byte, 16)
			rand.Read(bytes)
			streamResp := &ChatStreamResponse{
				ID:                "chatcmpl-" + hex.EncodeToString(bytes),
				Object:            response.Object,
				Created:           response.Created,
				Model:             response.Model,
				SystemFingerprint: response.SystemFingerprint,
				Provider:          "openai",
			}

			if len(response.Choices) > 0 {
				streamResp.Choices = []StreamChoice{{
					Index: response.Choices[0].Index,
					Delta: MessageDelta{
						Role:             response.Choices[0].Delta.Role,
						Content:          response.Choices[0].Delta.Content,
						ReasoningContent: response.Choices[0].Delta.ReasoningContent,
					},
					FinishReason: string(response.Choices[0].FinishReason),
				}}
			}

			if response.Usage != nil {
				streamResp.Usage = &Usage{
					PromptTokens:     response.Usage.PromptTokens,
					CompletionTokens: response.Usage.CompletionTokens,
					TotalTokens:      response.Usage.TotalTokens,
				}
			}

			ch <- streamResp
		}
	}()
	return ch, nil
}
