package provider

import (
	"sync/atomic"
	"testing"

	"github.com/Done-0/gin-scaffold/configs"
)

func TestNewOpenAI(t *testing.T) {
	tests := []struct {
		name    string
		config  *configs.ProviderInstanceConfig
		wantErr bool
	}{
		{
			name: "valid config",
			config: &configs.ProviderInstanceConfig{
				Enabled:    true,
				BaseURL:    "https://api.openai.com/v1",
				Keys:       []string{"test-key"},
				Models:     []string{"gpt-3.5-turbo"},
				Timeout:    30,
				MaxRetries: 3,
				RateLimit:  "60/min",
			},
			wantErr: false,
		},
		{
			name: "invalid rate limit",
			config: &configs.ProviderInstanceConfig{
				Enabled:    true,
				BaseURL:    "https://api.openai.com/v1",
				Keys:       []string{"test-key"},
				Models:     []string{"gpt-3.5-turbo"},
				Timeout:    30,
				MaxRetries: 3,
				RateLimit:  "invalid",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keyCounter := new(uint64)
			modelCounter := new(uint64)
			provider, err := NewOpenAI(tt.config, keyCounter, modelCounter)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewOpenAI() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && provider == nil {
				t.Error("NewOpenAI() returned nil provider")
			}
			if tt.wantErr && provider != nil {
				t.Error("NewOpenAI() should return nil on error")
			}
		})
	}
}

func TestNewGemini(t *testing.T) {
	tests := []struct {
		name    string
		config  *configs.ProviderInstanceConfig
		wantErr bool
	}{
		{
			name: "valid config",
			config: &configs.ProviderInstanceConfig{
				Enabled:    true,
				BaseURL:    "https://generativelanguage.googleapis.com/v1beta",
				Keys:       []string{"test-key"},
				Models:     []string{"gemini-pro"},
				Timeout:    30,
				MaxRetries: 3,
				RateLimit:  "30/min",
			},
			wantErr: false,
		},
		{
			name: "invalid rate limit",
			config: &configs.ProviderInstanceConfig{
				Enabled:    true,
				BaseURL:    "https://generativelanguage.googleapis.com/v1beta",
				Keys:       []string{"test-key"},
				Models:     []string{"gemini-pro"},
				Timeout:    30,
				MaxRetries: 3,
				RateLimit:  "invalid",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keyCounter := new(uint64)
			modelCounter := new(uint64)
			provider, err := NewGemini(tt.config, keyCounter, modelCounter)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewGemini() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && provider == nil {
				t.Error("NewGemini() returned nil provider")
			}
			if tt.wantErr && provider != nil {
				t.Error("NewGemini() should return nil on error")
			}
		})
	}
}

func TestProviderRateLimiter(t *testing.T) {
	config := &configs.ProviderInstanceConfig{
		Enabled:    true,
		BaseURL:    "https://api.openai.com/v1",
		Keys:       []string{"test-key"},
		Models:     []string{"gpt-3.5-turbo"},
		Timeout:    30,
		MaxRetries: 3,
		RateLimit:  "2/s",
	}

	keyCounter := new(uint64)
	modelCounter := new(uint64)
	provider, err := NewOpenAI(config, keyCounter, modelCounter)
	if err != nil {
		t.Fatalf("NewOpenAI() failed: %v", err)
	}

	openaiProvider, ok := provider.(*openAIProvider)
	if !ok {
		t.Fatal("provider is not *openAIProvider")
	}

	if openaiProvider.rateLimiter == nil {
		t.Fatal("rate limiter is nil")
	}

	if !openaiProvider.rateLimiter.Allow() {
		t.Error("rate limiter should allow first request")
	}
}

func TestKeyRoundRobin(t *testing.T) {
	config := &configs.ProviderInstanceConfig{
		Enabled:    true,
		BaseURL:    "https://api.openai.com/v1",
		Keys:       []string{"key-1", "key-2", "key-3"},
		Models:     []string{"gpt-3.5-turbo"},
		Timeout:    30,
		MaxRetries: 3,
		RateLimit:  "100/s",
	}

	keyCounter := new(uint64)
	modelCounter := new(uint64)
	expected := []string{"key-1", "key-2", "key-3", "key-1", "key-2", "key-3"}

	for i, want := range expected {
		provider, err := NewOpenAI(config, keyCounter, modelCounter)
		if err != nil {
			t.Fatalf("NewOpenAI() failed: %v", err)
		}

		p := provider.(*openAIProvider)
		keyIndex := atomic.AddUint64(p.keyCounter, 1) - 1
		got := config.Keys[keyIndex%uint64(len(config.Keys))]

		if got != want {
			t.Errorf("request %d: got %s, want %s", i, got, want)
		}
	}

	if *keyCounter != uint64(len(expected)) {
		t.Errorf("counter = %d, want %d", *keyCounter, len(expected))
	}
}

func TestKeyRoundRobinGemini(t *testing.T) {
	config := &configs.ProviderInstanceConfig{
		Enabled:    true,
		BaseURL:    "https://generativelanguage.googleapis.com",
		Keys:       []string{"key-1", "key-2", "key-3"},
		Models:     []string{"gemini-pro"},
		Timeout:    30,
		MaxRetries: 3,
		RateLimit:  "100/s",
	}

	keyCounter := new(uint64)
	modelCounter := new(uint64)
	expected := []string{"key-1", "key-2", "key-3", "key-1", "key-2", "key-3"}

	for i, want := range expected {
		provider, err := NewGemini(config, keyCounter, modelCounter)
		if err != nil {
			t.Fatalf("NewGemini() failed: %v", err)
		}

		p := provider.(*geminiProvider)
		keyIndex := atomic.AddUint64(p.keyCounter, 1) - 1
		got := config.Keys[keyIndex%uint64(len(config.Keys))]

		if got != want {
			t.Errorf("request %d: got %s, want %s", i, got, want)
		}
	}

	if *keyCounter != uint64(len(expected)) {
		t.Errorf("counter = %d, want %d", *keyCounter, len(expected))
	}
}

func TestDynamicKeyUpdate(t *testing.T) {
	config := &configs.ProviderInstanceConfig{
		Enabled:    true,
		BaseURL:    "https://api.openai.com/v1",
		Keys:       []string{"key-1", "key-2", "key-3"},
		Models:     []string{"gpt-3.5-turbo"},
		Timeout:    30,
		MaxRetries: 3,
		RateLimit:  "100/s",
	}

	keyCounter := new(uint64)
	modelCounter := new(uint64)

	for i := 0; i < 3; i++ {
		provider, _ := NewOpenAI(config, keyCounter, modelCounter)
		p := provider.(*openAIProvider)
		keyIndex := atomic.AddUint64(p.keyCounter, 1) - 1
		_ = config.Keys[keyIndex%uint64(len(config.Keys))]
	}

	config.Keys = []string{"new-key-1", "new-key-2"}

	expected := []string{"new-key-2", "new-key-1", "new-key-2"}
	for i, want := range expected {
		provider, err := NewOpenAI(config, keyCounter, modelCounter)
		if err != nil {
			t.Fatalf("NewOpenAI() failed: %v", err)
		}

		p := provider.(*openAIProvider)
		keyIndex := atomic.AddUint64(p.keyCounter, 1) - 1
		got := config.Keys[keyIndex%uint64(len(config.Keys))]

		if got != want {
			t.Errorf("after config update, request %d: got %s, want %s", i, got, want)
		}
	}
}

// TestMultiInstanceMultiKeyMultiModel tests multiple instances with multiple keys and models
func TestMultiInstanceMultiKeyMultiModel(t *testing.T) {
	instance1 := &configs.ProviderInstanceConfig{
		Enabled:    true,
		Name:       "instance-01",
		BaseURL:    "https://api.openai.com/v1",
		Keys:       []string{"i1-key-A", "i1-key-B"},
		Models:     []string{"gpt-3.5", "gpt-4"},
		Timeout:    30,
		MaxRetries: 3,
		RateLimit:  "100/s",
	}

	instance2 := &configs.ProviderInstanceConfig{
		Enabled:    true,
		Name:       "instance-02",
		BaseURL:    "https://api.openai.com/v1",
		Keys:       []string{"i2-key-X", "i2-key-Y", "i2-key-Z"},
		Models:     []string{"gpt-4-turbo"},
		Timeout:    30,
		MaxRetries: 3,
		RateLimit:  "100/s",
	}

	keyCounter1 := new(uint64)
	modelCounter1 := new(uint64)
	keyCounter2 := new(uint64)
	modelCounter2 := new(uint64)

	type testCase struct {
		config        *configs.ProviderInstanceConfig
		keyCounter    *uint64
		modelCounter  *uint64
		expectedKey   string
		expectedModel string
	}

	tests := []testCase{
		{instance1, keyCounter1, modelCounter1, "i1-key-A", "gpt-3.5"},
		{instance2, keyCounter2, modelCounter2, "i2-key-X", "gpt-4-turbo"},
		{instance1, keyCounter1, modelCounter1, "i1-key-B", "gpt-4"},
		{instance2, keyCounter2, modelCounter2, "i2-key-Y", "gpt-4-turbo"},
		{instance1, keyCounter1, modelCounter1, "i1-key-A", "gpt-3.5"},
		{instance2, keyCounter2, modelCounter2, "i2-key-Z", "gpt-4-turbo"},
	}

	for i, tc := range tests {
		provider, err := NewOpenAI(tc.config, tc.keyCounter, tc.modelCounter)
		if err != nil {
			t.Fatalf("request %d: NewOpenAI() failed: %v", i, err)
		}

		p := provider.(*openAIProvider)

		keyIndex := atomic.AddUint64(p.keyCounter, 1) - 1
		gotKey := tc.config.Keys[keyIndex%uint64(len(tc.config.Keys))]

		modelIndex := atomic.AddUint64(p.modelCounter, 1) - 1
		gotModel := tc.config.Models[modelIndex%uint64(len(tc.config.Models))]

		if gotKey != tc.expectedKey {
			t.Errorf("request %d (%s): key = %s, want %s", i, tc.config.Name, gotKey, tc.expectedKey)
		}

		if gotModel != tc.expectedModel {
			t.Errorf("request %d (%s): model = %s, want %s", i, tc.config.Name, gotModel, tc.expectedModel)
		}
	}

	if *keyCounter1 != 3 {
		t.Errorf("instance1 key counter = %d, want 3", *keyCounter1)
	}
	if *modelCounter1 != 3 {
		t.Errorf("instance1 model counter = %d, want 3", *modelCounter1)
	}
	if *keyCounter2 != 3 {
		t.Errorf("instance2 key counter = %d, want 3", *keyCounter2)
	}
	if *modelCounter2 != 3 {
		t.Errorf("instance2 model counter = %d, want 3", *modelCounter2)
	}
}

// TestConfigShrinkage verifies behavior when config items are reduced
func TestConfigShrinkage(t *testing.T) {
	config := &configs.ProviderInstanceConfig{
		Enabled:    true,
		Name:       "test",
		BaseURL:    "https://api.openai.com/v1",
		Keys:       []string{"k1", "k2", "k3", "k4", "k5"},
		Models:     []string{"m1", "m2", "m3"},
		Timeout:    30,
		MaxRetries: 3,
		RateLimit:  "100/s",
	}

	keyCounter := new(uint64)
	modelCounter := new(uint64)

	for i := 0; i < 10; i++ {
		provider, _ := NewOpenAI(config, keyCounter, modelCounter)
		p := provider.(*openAIProvider)
		keyIndex := atomic.AddUint64(p.keyCounter, 1) - 1
		_ = config.Keys[keyIndex%uint64(len(config.Keys))]
	}

	if *keyCounter != 10 {
		t.Errorf("after 10 requests: keyCounter = %d, want 10", *keyCounter)
	}

	config.Keys = []string{"new-k1", "new-k2"}
	config.Models = []string{"new-m1"}

	expectedKeys := []string{"new-k1", "new-k2", "new-k1", "new-k2", "new-k1"}
	expectedModels := []string{"new-m1", "new-m1", "new-m1", "new-m1", "new-m1"}

	for i := 0; i < 5; i++ {
		provider, err := NewOpenAI(config, keyCounter, modelCounter)
		if err != nil {
			t.Fatalf("NewOpenAI() failed: %v", err)
		}

		p := provider.(*openAIProvider)

		keyIndex := atomic.AddUint64(p.keyCounter, 1) - 1
		gotKey := config.Keys[keyIndex%uint64(len(config.Keys))]

		modelIndex := atomic.AddUint64(p.modelCounter, 1) - 1
		gotModel := config.Models[modelIndex%uint64(len(config.Models))]

		if gotKey != expectedKeys[i] {
			t.Errorf("request %d after shrink: key = %s, want %s (counter=%d)", i, gotKey, expectedKeys[i], *keyCounter)
		}

		if gotModel != expectedModels[i] {
			t.Errorf("request %d after shrink: model = %s, want %s (counter=%d)", i, gotModel, expectedModels[i], *modelCounter)
		}
	}

	if *keyCounter != 15 {
		t.Errorf("final keyCounter = %d, want 15", *keyCounter)
	}
}
