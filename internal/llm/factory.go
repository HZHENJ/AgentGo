package llm

import (
	"sync"
	"context"
	"fmt"
	"agentgo/pkg/conf"
)


const (
	TypeOpenAI = "openai"
	TypeOllama = "ollama"
	TypeDeepseek = "deepseek"
)

type ModelBuilder func(ctx context.Context, cfg *conf.LLMConfig) (Model, error)

// global registry for model builders
var (
	mu sync.RWMutex
	registry = make(map[string]ModelBuilder)
)

// Register registers a model builder for a given mode type.
func Register(modelType string, builder ModelBuilder) {
	mu.Lock()
	defer mu.Unlock()

	if builder == nil {
		panic("llm: Register builder is nil")
	}

	if _, dup := registry[modelType]; dup {
		panic("llm: Register called twice for type " + modelType)
	}
	registry[modelType] = builder
}

// CreateModel creates a model instance based on the given model type and configuration.
func CreateModel(ctx context.Context, modelType string, cfg *conf.LLMConfig) (Model, error) {
	mu.RLock()
	builder, ok := registry[modelType]
	mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("llm: unsupported model type: %q (forgot to register?)", modelType)
	}
	return builder(ctx, cfg)
}