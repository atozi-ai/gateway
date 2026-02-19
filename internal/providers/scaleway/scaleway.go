package scaleway

import (
	"context"

	"github.com/atozi-ai/gateway/internal/domain/llm"
	"github.com/atozi-ai/gateway/internal/providers/openai_compat"
)

const baseURL = "https://api.scaleway.ai/v1"

// Provider implements llm.Provider for the Scaleway AI API.
type Provider struct {
	client *openaicompat.Client
}

// New creates a Scaleway AI provider.
func New() *Provider {
	return &Provider{
		client: openaicompat.NewClient(openaicompat.Config{
			BaseURL: baseURL,
			APIKey:  "",
		}),
	}
}

func (p *Provider) Name() string { return "scaleway" }

func (p *Provider) Chat(ctx context.Context, req llm.ChatRequest) (*llm.ChatResponse, error) {
	return p.client.Chat(ctx, req)
}

func (p *Provider) ChatStream(ctx context.Context, req llm.ChatRequest, callback func(*llm.StreamChunk) error) error {
	return p.client.ChatStream(ctx, req, callback)
}
