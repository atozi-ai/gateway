package zai

import (
	"context"

	"github.com/atozi-ai/gateway/internal/domain/llm"
	"github.com/atozi-ai/gateway/internal/providers/openai_compat"
)

const baseURL = "https://api.z.ai/api/paas/v4"

// Provider implements llm.Provider for the Z.ai API.
type Provider struct {
	client *openaicompat.Client
}

// New creates an Z.ai provider with the given API key.
func New(apiKey string) *Provider {
	return &Provider{
		client: openaicompat.NewClient(openaicompat.Config{
			BaseURL: baseURL,
			APIKey:  apiKey,
		}),
	}
}

func (p *Provider) Name() string { return "zai" }

func (p *Provider) Chat(ctx context.Context, req llm.ChatRequest) (*llm.ChatResponse, error) {
	return p.client.Chat(ctx, req)
}

func (p *Provider) ChatStream(ctx context.Context, req llm.ChatRequest, callback func(*llm.StreamChunk) error) error {
	return p.client.ChatStream(ctx, req, callback)
}
