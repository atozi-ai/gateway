package openai

import (
	"encoding/json"

	"github.com/atozi-ai/gateway/internal/domain/llm"
)

type openAIChatResponse struct {
	ID      string `json:"id"`
	Model   string `json:"model"`
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type openAIResponseFormat struct {
	Type       string            `json:"type"`
	JSONSchema *openAIJSONSchema `json:"json_schema,omitempty"`
}

type openAIJSONSchema struct {
	Name   string          `json:"name"`
	Schema json.RawMessage `json:"schema"`
	Strict bool            `json:"strict,omitempty"`
}

type openAIChatRequest struct {
	Model          string                `json:"model"`
	Messages       []llm.Message         `json:"messages"`
	Temperature    *float32              `json:"temperature,omitempty"`
	TopP           *float32              `json:"top_p,omitempty"`
	MaxTokens      *int                  `json:"max_completion_tokens,omitempty"`
	Stop           []string              `json:"stop,omitempty"`
	Verbosity      *llm.Verbosity        `json:"verbosity,omitempty"`
	ResponseFormat *openAIResponseFormat `json:"response_format,omitempty"`
}
