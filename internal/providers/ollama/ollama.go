package ollama

import (
	"context"

	"github.com/atozi-ai/gateway/internal/domain/llm"
	"github.com/atozi-ai/gateway/internal/providers/openai_compat"
)

const baseURL = "http://localhost:11434/v1"

// Provider implements llm.Provider for the Ollama API (OpenAI-compatible).
// Ollama runs locally, so the baseURL assumes localhost.
// Users can override by setting OLLAMA_HOST environment variable.
type Provider struct {
	client *openaicompat.Client
}

// New creates an Ollama provider.
func New() *Provider {
	return &Provider{
		client: openaicompat.NewClient(openaicompat.Config{
			BaseURL: baseURL,
			APIKey:  "ollama", // Ollama doesn't require API key but OpenAI client needs one
		}),
	}
}

func (p *Provider) Name() string { return "ollama" }

func (p *Provider) Chat(ctx context.Context, req llm.ChatRequest) (*llm.ChatResponse, error) {
	return p.client.Chat(ctx, req)
}

func (p *Provider) ChatStream(ctx context.Context, req llm.ChatRequest, callback func(*llm.StreamChunk) error) error {
	return p.client.ChatStream(ctx, req, callback)
}
