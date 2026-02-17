package providers

import (
	"context"
	"testing"

	"github.com/atozi-ai/gateway/internal/domain/llm"
)

func TestOpenAIChat(t *testing.T) {
	req := llm.ChatRequest{
		Model: "openai/gpt-4.1-mini",
		Messages: []llm.Message{
			{
				Role:    llm.RoleUser,
				Content: "Hi",
			},
		},
	}

	p, model, err := Get(req.Model, "test-api-key")
	if err != nil {
		t.Fatal(err)
	}
	req.Model = model

	resp, err := p.Chat(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.Content == "" {
		t.Fatal("empty response")
	}
}
