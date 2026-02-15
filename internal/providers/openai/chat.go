package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/atozi-ai/gateway/internal/domain/llm"
	"github.com/atozi-ai/gateway/internal/platform/logger"
)

type OpenAI struct {
	client *http.Client
	apikey string
}

func New() llm.Provider {
	apikey := os.Getenv("OPENAI_API_KEY")
	if apikey == "" {
		logger.Log.Fatal().Msg("OPENAI_API_KEY environment variable is required")
	}

	return &OpenAI{
		client: newHTTPClient(),
		apikey: apikey,
	}
}

func (o *OpenAI) Name() string {
	return "openai"
}

func (o *OpenAI) Chat(
	ctx context.Context,
	req llm.ChatRequest,
) (*llm.ChatResponse, error) {
	body := toOpenAIRequest(req)

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	log := logger.FromContext(ctx)
	log.Info().Str("model", req.Model).Msg("Using Model")

	httpReq, err := http.NewRequestWithContext(
		ctx,
		"POST",
		"https://api.openai.com/v1/chat/completions",
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+o.apikey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := o.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(io.LimitReader(resp.Body, 10*1024*1024))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var apiError struct {
			Error struct {
				Message string `json:"message"`
				Type    string `json:"type"`
				Code    string `json:"code"`
			} `json:"error"`
		}
		if err := json.Unmarshal(bodyBytes, &apiError); err == nil && apiError.Error.Message != "" {
			return nil, fmt.Errorf("OpenAI API error: %s (type: %s, code: %s)",
				apiError.Error.Message, apiError.Error.Type, apiError.Error.Code)
		}
		return nil, fmt.Errorf("OpenAI API returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var raw openAIChatResponse
	if err := json.Unmarshal(bodyBytes, &raw); err != nil {
		return nil, err
	}

	content := ""
	if len(raw.Choices) > 0 {
		content = raw.Choices[0].Message.Content
	}

	return &llm.ChatResponse{
		ID:      raw.ID,
		Model:   raw.Model,
		Content: content,
	}, nil
}
