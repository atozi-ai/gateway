package azure

import (
	"context"

	"github.com/atozi-ai/gateway/internal/domain/llm"
	openaicompat "github.com/atozi-ai/gateway/internal/providers/openai_compat"
)

// Provider implements llm.Provider for the OpenAI API.
type Provider struct {
	apiKey   string
	endpoint string
}

func New(apiKey string, endpoint string) *Provider {
	return &Provider{
		apiKey:   apiKey,
		endpoint: endpoint,
	}
}

func (p *Provider) Name() string { return "azure" }

func (p *Provider) Chat(ctx context.Context, req llm.ChatRequest) (*llm.ChatResponse, error) {
	client := p.getClient(req)
	return client.Chat(ctx, req)
}

func (p *Provider) ChatStream(ctx context.Context, req llm.ChatRequest, callback func(*llm.StreamChunk) error) error {
	client := p.getClient(req)
	return client.ChatStream(ctx, req, callback)
}

func (p *Provider) getClient(req llm.ChatRequest) *openaicompat.Client {
	endpoint := p.endpoint
	if req.Options.AzureEndpoint != nil && *req.Options.AzureEndpoint != "" {
		endpoint = *req.Options.AzureEndpoint
	}

	if endpoint == "" {
		panic("Azure endpoint is required. Provide azure_endpoint in request options")
	}

	return openaicompat.NewClient(openaicompat.Config{
		BaseURL: endpoint,
		APIKey:  p.apiKey,
		Headers: map[string]string{
			"api-key": p.apiKey,
		},
	})
}
