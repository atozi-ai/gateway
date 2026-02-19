package baseten

import (
	"context"

	"github.com/atozi-ai/gateway/internal/domain/llm"
	"github.com/atozi-ai/gateway/internal/providers/openai_compat"
)

const baseURL = "https://model-<model-id>.api.baseten.co/production/predict"

// Provider implements llm.Provider for the Baseten API.
// Note: Baseten uses model-specific URLs, but we use a placeholder here.
type Provider struct {
	client *openaicompat.Client
}

// New creates a Baseten provider.
func New() *Provider {
	return &Provider{
		client: openaicompat.NewClient(openaicompat.Config{
			BaseURL: "https://api.baseten.co",
			APIKey:  "",
		}),
	}
}

func (p *Provider) Name() string { return "baseten" }

func (p *Provider) Chat(ctx context.Context, req llm.ChatRequest) (*llm.ChatResponse, error) {
	return p.client.Chat(ctx, req)
}

func (p *Provider) ChatStream(ctx context.Context, req llm.ChatRequest, callback func(*llm.StreamChunk) error) error {
	return p.client.ChatStream(ctx, req, callback)
}
