package ai21

import (
	"context"

	"github.com/atozi-ai/gateway/internal/domain/llm"
	"github.com/atozi-ai/gateway/internal/providers/openai_compat"
)

const baseURL = "https://api.ai21.com/v1"

// Provider implements llm.Provider for the AI21 Labs API.
type Provider struct {
	client *openaicompat.Client
}

// New creates an AI21 Labs provider.
func New() *Provider {
	return &Provider{
		client: openaicompat.NewClient(openaicompat.Config{
			BaseURL: baseURL,
			APIKey:  "",
		}),
	}
}

func (p *Provider) Name() string { return "ai21" }

func (p *Provider) Chat(ctx context.Context, req llm.ChatRequest) (*llm.ChatResponse, error) {
	return p.client.Chat(ctx, req)
}

func (p *Provider) ChatStream(ctx context.Context, req llm.ChatRequest, callback func(*llm.StreamChunk) error) error {
	return p.client.ChatStream(ctx, req, callback)
}
