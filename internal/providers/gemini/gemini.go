package gemini

import (
	"context"

	"github.com/atozi-ai/gateway/internal/domain/llm"
	"github.com/atozi-ai/gateway/internal/providers/gemini_compat"
)

type Provider struct {
	client *gemini_compat.Client
}

func New() *Provider {
	return &Provider{
		client: gemini_compat.NewClient(gemini_compat.Config{
			APIKey: "",
		}),
	}
}

func NewWithBaseURL(apiKey string, baseURL string) *Provider {
	return &Provider{
		client: gemini_compat.NewClient(gemini_compat.Config{
			BaseURL: baseURL,
			APIKey:  apiKey,
		}),
	}
}

func (p *Provider) Name() string { return "gemini" }

func (p *Provider) Chat(ctx context.Context, req llm.ChatRequest) (*llm.ChatResponse, error) {
	return p.client.Chat(ctx, req)
}

func (p *Provider) ChatStream(ctx context.Context, req llm.ChatRequest, callback func(*llm.StreamChunk) error) error {
	return p.client.ChatStream(ctx, req, callback)
}
