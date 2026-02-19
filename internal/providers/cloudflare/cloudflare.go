package cloudflare

import (
	"context"

	"github.com/atozi-ai/gateway/internal/domain/llm"
	"github.com/atozi-ai/gateway/internal/providers/openai_compat"
)

const baseURL = "https://api.cloudflare.com/client/v4/accounts/{account_id}/ai/v1"

// Provider implements llm.Provider for the Cloudflare Workers AI API.
// Note: Requires account_id in the URL.
type Provider struct {
	client *openaicompat.Client
}

// New creates a Cloudflare Workers AI provider.
func New() *Provider {
	return &Provider{
		client: openaicompat.NewClient(openaicompat.Config{
			BaseURL: "https://api.cloudflare.com/client/v4/accounts/-/ai/v1",
			APIKey:  "",
		}),
	}
}

func (p *Provider) Name() string { return "cloudflare" }

func (p *Provider) Chat(ctx context.Context, req llm.ChatRequest) (*llm.ChatResponse, error) {
	return p.client.Chat(ctx, req)
}

func (p *Provider) ChatStream(ctx context.Context, req llm.ChatRequest, callback func(*llm.StreamChunk) error) error {
	return p.client.ChatStream(ctx, req, callback)
}
