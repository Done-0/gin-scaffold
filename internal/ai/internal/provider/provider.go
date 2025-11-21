// Package provider implements AI provider interfaces
// Author: Done-0
// Created: 2025-08-31
package provider

import (
	"context"
	"fmt"
	"log"
	"sync/atomic"

	"github.com/Done-0/gin-scaffold/configs"
)

type provider struct {
	instanceCounter uint64             // Round Robin selection
	keyCounters     map[string]*uint64 // key: "provider:instance", value: counter pointer
	modelCounters   map[string]*uint64 // key: "provider:instance", value: counter pointer
}

func New() Provider {
	return &provider{
		keyCounters:   make(map[string]*uint64),
		modelCounters: make(map[string]*uint64),
	}
}

func (p *provider) Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	client, err := p.getProvider()
	if err != nil {
		return nil, err
	}
	return client.Chat(ctx, req)
}

func (p *provider) ChatStream(ctx context.Context, req *ChatRequest) (<-chan *ChatStreamResponse, error) {
	client, err := p.getProvider()
	if err != nil {
		return nil, err
	}
	return client.ChatStream(ctx, req)
}

func (p *provider) getProvider() (Provider, error) {
	cfg, err := configs.GetConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get config: %w", err)
	}

	type providerInstance struct {
		name     string
		instance configs.ProviderInstanceConfig
	}

	var instances []providerInstance
	for name, prov := range cfg.AI.Providers {
		if !prov.Enabled {
			continue
		}

		for _, inst := range prov.Instances {
			if inst.Enabled {
				instances = append(instances, providerInstance{
					name:     name,
					instance: inst,
				})
			}
		}
	}

	// Round Robin selection
	selected := instances[atomic.AddUint64(&p.instanceCounter, 1)%uint64(len(instances))]

	log.Printf("Using %s provider, instance: %s", selected.name, selected.instance.Name)

	counterKey := fmt.Sprintf("%s:%s", selected.name, selected.instance.Name)

	keyCounter, exists := p.keyCounters[counterKey]
	if !exists {
		keyCounter = new(uint64)
		p.keyCounters[counterKey] = keyCounter
	}

	modelCounter, exists := p.modelCounters[counterKey]
	if !exists {
		modelCounter = new(uint64)
		p.modelCounters[counterKey] = modelCounter
	}

	switch selected.name {
	case "openai":
		return NewOpenAI(&selected.instance, keyCounter, modelCounter)
	case "gemini":
		return NewGemini(&selected.instance, keyCounter, modelCounter)
	}

	return nil, fmt.Errorf("unsupported provider: %s", selected.name)
}
