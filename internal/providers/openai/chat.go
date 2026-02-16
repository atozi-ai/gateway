package openai

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

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
		var apiError openAIErrorResponse
		if err := json.Unmarshal(bodyBytes, &apiError); err == nil && apiError.Error.Message != "" {
			return nil, &llm.ProviderError{
				StatusCode: resp.StatusCode,
				Message:    apiError.Error.Message,
				Type:       apiError.Error.Type,
				Code:       apiError.Error.Code,
				Param:      apiError.Error.Param,
				Raw:        bodyBytes,
			}
		}
		// If we can't parse the error, still return a ProviderError with raw response
		return nil, &llm.ProviderError{
			StatusCode: resp.StatusCode,
			Message:    fmt.Sprintf("OpenAI API returned status %d", resp.StatusCode),
			Raw:        bodyBytes,
		}
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
		Raw:     bodyBytes, // Include raw response
	}, nil
}

func (o *OpenAI) ChatStream(
	ctx context.Context,
	req llm.ChatRequest,
	callback func(*llm.StreamChunk) error,
) error {
	// Ensure stream is set to true
	stream := true
	req.Options.Stream = &stream

	body := toOpenAIRequest(req)

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	log := logger.FromContext(ctx)
	log.Info().Str("model", req.Model).Msg("Using Model for streaming")

	httpReq, err := http.NewRequestWithContext(
		ctx,
		"POST",
		"https://api.openai.com/v1/chat/completions",
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+o.apikey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := o.client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(io.LimitReader(resp.Body, 10*1024*1024))
		var apiError openAIErrorResponse
		if err := json.Unmarshal(bodyBytes, &apiError); err == nil && apiError.Error.Message != "" {
			return &llm.ProviderError{
				StatusCode: resp.StatusCode,
				Message:    apiError.Error.Message,
				Type:       apiError.Error.Type,
				Code:       apiError.Error.Code,
				Param:      apiError.Error.Param,
				Raw:        bodyBytes,
			}
		}
		return &llm.ProviderError{
			StatusCode: resp.StatusCode,
			Message:    fmt.Sprintf("OpenAI API returned status %d", resp.StatusCode),
			Raw:        bodyBytes,
		}
	}

	// Read SSE stream
	scanner := bufio.NewScanner(resp.Body)

	for scanner.Scan() {
		line := scanner.Text()

		// Skip empty lines
		if strings.TrimSpace(line) == "" {
			continue
		}

		// Check if this is a data line
		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")

			// Check for [DONE] marker
			if strings.TrimSpace(data) == "[DONE]" {
				return nil
			}

			// Parse JSON chunk
			var chunk openAIStreamChunk
			if err := json.Unmarshal([]byte(data), &chunk); err != nil {
				log.Warn().Err(err).Str("data", data).Msg("Failed to parse stream chunk")
				continue
			}

			// Convert to domain StreamChunk
			streamChunk := &llm.StreamChunk{
				ID:      chunk.ID,
				Object:  chunk.Object,
				Created: chunk.Created,
				Model:   chunk.Model,
				Raw:     []byte(data),
			}

			// Convert choices
			streamChunk.Choices = make([]llm.StreamChoice, len(chunk.Choices))
			for i, choice := range chunk.Choices {
				streamChoice := llm.StreamChoice{
					Index: choice.Index,
				}

				if choice.Delta.Role != nil {
					streamChoice.Delta.Role = choice.Delta.Role
				}
				if choice.Delta.Content != nil {
					streamChoice.Delta.Content = choice.Delta.Content
				}
				if len(choice.Delta.ToolCalls) > 0 {
					streamChoice.Delta.ToolCalls = make([]interface{}, len(choice.Delta.ToolCalls))
					for j, tc := range choice.Delta.ToolCalls {
						streamChoice.Delta.ToolCalls[j] = tc
					}
				}

				if choice.FinishReason != nil {
					streamChoice.FinishReason = choice.FinishReason
				}
				if choice.Logprobs != nil {
					streamChoice.Logprobs = choice.Logprobs
				}

				streamChunk.Choices[i] = streamChoice
			}

			// Convert usage if present
			if chunk.Usage != nil {
				streamChunk.Usage = &llm.Usage{
					PromptTokens:     chunk.Usage.PromptTokens,
					CompletionTokens: chunk.Usage.CompletionTokens,
					TotalTokens:       chunk.Usage.TotalTokens,
				}
			}

			// Call callback
			if err := callback(streamChunk); err != nil {
				return err
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to read stream: %w", err)
	}

	return nil
}
