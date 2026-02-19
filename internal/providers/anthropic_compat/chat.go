package anthropic_compat

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/atozi-ai/gateway/internal/domain/llm"
	"github.com/atozi-ai/gateway/internal/platform/logger"
)

const (
	anthropicVersion = "2023-06-01"
	maxResponseBytes = 10 * 1024 * 1024 // 10 MB
)

func (c *Client) Chat(ctx context.Context, req llm.ChatRequest) (*llm.ChatResponse, error) {
	body := toRequest(req)

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, llm.NewInternalError(fmt.Sprintf("failed to marshal request: %v", err))
	}

	log := logger.FromContext(ctx)
	log.Info().Str("model", req.Model).Msg("Using Anthropic Model")

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint(), bytes.NewReader(jsonBody))
	if err != nil {
		return nil, llm.NewInternalError(fmt.Sprintf("failed to create request: %v", err))
	}
	c.setHeaders(httpReq, req.APIKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, llm.NewProviderError(503, fmt.Sprintf("failed to execute request: %v", err), "service_unavailable", "request_failed")
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(io.LimitReader(resp.Body, 10*1024*1024))
	if err != nil {
		return nil, llm.NewInternalError(fmt.Sprintf("failed to read response: %v", err))
	}

	if err := checkError(resp.StatusCode, respBody); err != nil {
		return nil, err
	}

	var raw messageResponse
	if err := json.Unmarshal(respBody, &raw); err != nil {
		return nil, llm.NewInternalError(fmt.Sprintf("failed to unmarshal response: %v", err))
	}

	content := ""
	for _, block := range raw.Content {
		if block.Type == contentTypeText {
			content += block.Text
		}
	}

	return &llm.ChatResponse{
		ID:      raw.ID,
		Model:   raw.Model,
		Content: content,
		Raw:     respBody,
	}, nil
}

func (c *Client) ChatStream(ctx context.Context, req llm.ChatRequest, callback func(*llm.StreamChunk) error) error {
	stream := true
	req.Options.Stream = &stream

	body := toRequest(req)

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return llm.NewInternalError(fmt.Sprintf("failed to marshal request: %v", err))
	}

	log := logger.FromContext(ctx)
	log.Info().Str("model", req.Model).Msg("Using Anthropic Model for streaming")

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint(), bytes.NewReader(jsonBody))
	if err != nil {
		return llm.NewInternalError(fmt.Sprintf("failed to create request: %v", err))
	}
	c.setHeaders(httpReq, req.APIKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return llm.NewProviderError(503, fmt.Sprintf("failed to execute request: %v", err), "service_unavailable", "request_failed")
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, readErr := io.ReadAll(io.LimitReader(resp.Body, 10*1024*1024))
		if readErr != nil {
			return llm.NewInternalError(fmt.Sprintf("failed to read error response: %v", readErr))
		}
		if err := checkError(resp.StatusCode, respBody); err != nil {
			return err
		}
	}

	return readSSEStream(ctx, resp.Body, callback)
}

func checkError(statusCode int, body []byte) error {
	if statusCode >= 200 && statusCode < 300 {
		return nil
	}

	var apiErr errorResponse
	if err := json.Unmarshal(body, &apiErr); err == nil && apiErr.Message != "" {
		return &llm.ProviderError{
			StatusCode: statusCode,
			Message:    apiErr.Message,
			Type:       "api_error",
			Raw:        body,
		}
	}

	return &llm.ProviderError{
		StatusCode: statusCode,
		Message:    fmt.Sprintf("API returned status %d", statusCode),
		Raw:        body,
	}
}

func (c *Client) endpoint() string {
	baseURL := c.cfg.BaseURL
	if baseURL == "" {
		baseURL = "https://api.anthropic.com"
	}
	return baseURL + "/v1/messages"
}

func readSSEStream(ctx context.Context, r io.Reader, callback func(*llm.StreamChunk) error) error {
	const maxScanTokenSize = 64 * 1024
	scanner := bufio.NewScanner(r)
	buf := make([]byte, 0, bufio.MaxScanTokenSize)
	scanner.Buffer(buf, maxScanTokenSize)

	log := logger.FromContext(ctx)

	messageID := ""
	model := ""
	var contentIndex int

	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		line := scanner.Text()

		if strings.TrimSpace(line) == "" {
			continue
		}

		if !strings.HasPrefix(line, "data: ") {
			continue
		}

		data := strings.TrimPrefix(line, "data: ")

		if strings.TrimSpace(data) == "[DONE]" {
			chunk := &llm.StreamChunk{
				ID:      messageID,
				Model:   model,
				Choices: []llm.StreamChoice{},
			}
			finishReason := "stop"
			chunk.Choices = append(chunk.Choices, llm.StreamChoice{
				Index:        contentIndex,
				FinishReason: &finishReason,
			})
			return callback(chunk)
		}

		var event map[string]interface{}
		if err := json.Unmarshal([]byte(data), &event); err != nil {
			log.Warn().Err(err).Str("data", data).Msg("Failed to parse SSE event")
			continue
		}

		eventType, _ := event["type"].(string)

		switch eventType {
		case "message_start":
			if msg, ok := event["message"].(map[string]interface{}); ok {
				if id, ok := msg["id"].(string); ok {
					messageID = id
				}
				if m, ok := msg["model"].(string); ok {
					model = m
				}
			}

			chunk := &llm.StreamChunk{
				ID:      messageID,
				Model:   model,
				Choices: []llm.StreamChoice{},
			}
			role := "assistant"
			chunk.Choices = append(chunk.Choices, llm.StreamChoice{
				Index: 0,
				Delta: llm.StreamDelta{
					Role: &role,
				},
			})
			if err := callback(chunk); err != nil {
				return err
			}

		case "content_block_start":
			if cb, ok := event["content_block"].(map[string]interface{}); ok {
				if idx, ok := event["index"].(float64); ok {
					contentIndex = int(idx)
				}
				_ = cb
			}

		case "content_block_delta":
			if delta, ok := event["delta"].(map[string]interface{}); ok {
				if text, ok := delta["text"].(string); ok {
					chunk := &llm.StreamChunk{
						ID:      messageID,
						Model:   model,
						Choices: []llm.StreamChoice{},
					}
					chunk.Choices = append(chunk.Choices, llm.StreamChoice{
						Index: contentIndex,
						Delta: llm.StreamDelta{
							Content: &text,
						},
					})
					if err := callback(chunk); err != nil {
						return err
					}
				}
			}

		case "content_block_stop":
			contentIndex++

		case "message_delta":
			var usage *llm.Usage
			if u, ok := event["usage"].(map[string]interface{}); ok {
				if outputTokens, ok := u["output_tokens"].(float64); ok {
					usage = &llm.Usage{
						CompletionTokens: int(outputTokens),
					}
				}
			}

			var finishReason *string
			if sr, ok := event["stop_reason"].(string); ok {
				finishReason = &sr
			}

			chunk := &llm.StreamChunk{
				ID:      messageID,
				Model:   model,
				Usage:   usage,
				Choices: []llm.StreamChoice{},
			}
			if finishReason != nil {
				chunk.Choices = append(chunk.Choices, llm.StreamChoice{
					Index:        contentIndex,
					FinishReason: finishReason,
				})
			}
			if err := callback(chunk); err != nil {
				return err
			}

		case "message_stop":
			return nil

		case "ping":
		}
	}

	if err := scanner.Err(); err != nil {
		return llm.NewInternalError(fmt.Sprintf("failed to read stream: %v", err))
	}

	return nil
}
