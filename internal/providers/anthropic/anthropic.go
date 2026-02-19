package anthropic

import (
	"context"

	"github.com/atozi-ai/gateway/internal/domain/llm"
	"github.com/atozi-ai/gateway/internal/providers/anthropic_compat"
)

const defaultMaxTokens = 4096

type Provider struct {
	client *anthropic_compat.Client
}

func New(apiKey string) *Provider {
	return &Provider{
		client: anthropic_compat.NewClient(anthropic_compat.Config{
			APIKey: apiKey,
		}),
	}
}

func NewWithBaseURL(apiKey string, baseURL string) *Provider {
	return &Provider{
		client: anthropic_compat.NewClient(anthropic_compat.Config{
			BaseURL: baseURL,
			APIKey:  apiKey,
		}),
	}
}

func (p *Provider) Name() string { return "anthropic" }

func (p *Provider) Chat(ctx context.Context, req llm.ChatRequest) (*llm.ChatResponse, error) {
	return p.client.Chat(ctx, req)
}

func (p *Provider) ChatStream(ctx context.Context, req llm.ChatRequest, callback func(*llm.StreamChunk) error) error {
	return p.client.ChatStream(ctx, req, callback)
}

func (p *Provider) GetDefaultMaxTokens() int {
	return defaultMaxTokens
}
