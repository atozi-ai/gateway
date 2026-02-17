package openaicompat

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

const maxResponseBytes = 10 * 1024 * 1024 // 10 MB

// Chat sends a non-streaming chat completions request and returns the parsed response.
func (c *Client) Chat(ctx context.Context, req llm.ChatRequest) (*llm.ChatResponse, error) {
	body := toRequest(req)

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, llm.NewInternalError(fmt.Sprintf("failed to marshal request: %v", err))
	}

	log := logger.FromContext(ctx)
	log.Info().Str("model", req.Model).Msg("Using Model")

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint(), bytes.NewReader(jsonBody))
	if err != nil {
		return nil, llm.NewInternalError(fmt.Sprintf("failed to create request: %v", err))
	}
	c.setHeaders(httpReq)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, llm.NewProviderError(503, fmt.Sprintf("failed to execute request: %v", err), "service_unavailable", "request_failed")
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseBytes))
	if err != nil {
		return nil, llm.NewInternalError(fmt.Sprintf("failed to read response: %v", err))
	}

	if err := checkError(resp.StatusCode, respBody); err != nil {
		return nil, err
	}

	var raw chatResponse
	if err := json.Unmarshal(respBody, &raw); err != nil {
		return nil, llm.NewInternalError(fmt.Sprintf("failed to unmarshal response: %v", err))
	}

	content := ""
	if len(raw.Choices) > 0 {
		content = raw.Choices[0].Message.Content
	}

	return &llm.ChatResponse{
		ID:      raw.ID,
		Model:   raw.Model,
		Content: content,
		Raw:     respBody,
	}, nil
}

// ChatStream sends a streaming chat completions request and invokes callback for each chunk.
// If the callback returns an error, streaming stops and that error is returned.
func (c *Client) ChatStream(
	ctx context.Context,
	req llm.ChatRequest,
	callback func(*llm.StreamChunk) error,
) error {
	// Ensure stream flag is set.
	stream := true
	req.Options.Stream = &stream

	body := toRequest(req)

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return llm.NewInternalError(fmt.Sprintf("failed to marshal request: %v", err))
	}

	log := logger.FromContext(ctx)
	log.Info().Str("model", req.Model).Msg("Using Model for streaming")

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint(), bytes.NewReader(jsonBody))
	if err != nil {
		return llm.NewInternalError(fmt.Sprintf("failed to create request: %v", err))
	}
	c.setHeaders(httpReq)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return llm.NewProviderError(503, fmt.Sprintf("failed to execute request: %v", err), "service_unavailable", "request_failed")
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, readErr := io.ReadAll(io.LimitReader(resp.Body, maxResponseBytes))
		if readErr != nil {
			return llm.NewInternalError(fmt.Sprintf("failed to read error response: %v", readErr))
		}
		if err := checkError(resp.StatusCode, respBody); err != nil {
			return err
		}
	}

	return readSSEStream(ctx, resp.Body, callback)
}

// endpoint returns the full chat completions URL.
func (c *Client) endpoint() string {
	baseURL := strings.TrimRight(c.cfg.BaseURL, "/")
	if strings.Contains(baseURL, "openai.azure.com") {
		return baseURL
	}
	return baseURL + "/chat/completions"
}

// checkError inspects the status code and tries to parse an API error.
func checkError(statusCode int, body []byte) error {
	if statusCode >= 200 && statusCode < 300 {
		return nil
	}

	var apiErr errorResponse
	if err := json.Unmarshal(body, &apiErr); err == nil && apiErr.Error.Message != "" {
		return &llm.ProviderError{
			StatusCode: statusCode,
			Message:    apiErr.Error.Message,
			Type:       apiErr.Error.Type,
			Code:       apiErr.Error.Code,
			Param:      apiErr.Error.Param,
			Raw:        body,
		}
	}

	return &llm.ProviderError{
		StatusCode: statusCode,
		Message:    fmt.Sprintf("API returned status %d", statusCode),
		Raw:        body,
	}
}

// readSSEStream reads an SSE stream and dispatches parsed chunks to callback.
func readSSEStream(
	ctx context.Context,
	r io.Reader,
	callback func(*llm.StreamChunk) error,
) error {
	// Set buffer size to handle large chunks efficiently (64KB)
	const maxScanTokenSize = 64 * 1024
	scanner := bufio.NewScanner(r)
	buf := make([]byte, 0, bufio.MaxScanTokenSize)
	scanner.Buffer(buf, maxScanTokenSize)

	log := logger.FromContext(ctx)

	for scanner.Scan() {
		// Check context cancellation
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
			return nil
		}

		chunk, err := parseStreamChunk([]byte(data))
		if err != nil {
			// Log parse errors but continue processing stream
			log.Warn().Err(err).Str("data", data).Msg("Failed to parse stream chunk")
			continue
		}

		if err := callback(chunk); err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return llm.NewInternalError(fmt.Sprintf("failed to read stream: %v", err))
	}

	return nil
}

// parseStreamChunk converts raw JSON into a domain StreamChunk.
func parseStreamChunk(data []byte) (*llm.StreamChunk, error) {
	var raw streamChunk
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	chunk := &llm.StreamChunk{
		ID:      raw.ID,
		Object:  raw.Object,
		Created: raw.Created,
		Model:   raw.Model,
		Raw:     data,
	}

	chunk.Choices = make([]llm.StreamChoice, len(raw.Choices))
	for i, c := range raw.Choices {
		sc := llm.StreamChoice{
			Index: c.Index,
		}

		if c.Delta.Role != nil {
			sc.Delta.Role = c.Delta.Role
		}
		if c.Delta.Content != nil {
			sc.Delta.Content = c.Delta.Content
		}
		if len(c.Delta.ToolCalls) > 0 {
			sc.Delta.ToolCalls = make([]interface{}, len(c.Delta.ToolCalls))
			copy(sc.Delta.ToolCalls, c.Delta.ToolCalls)
		}
		if c.FinishReason != nil {
			sc.FinishReason = c.FinishReason
		}
		if c.Logprobs != nil {
			sc.Logprobs = c.Logprobs
		}

		chunk.Choices[i] = sc
	}

	if raw.Usage != nil {
		chunk.Usage = &llm.Usage{
			PromptTokens:     raw.Usage.PromptTokens,
			CompletionTokens: raw.Usage.CompletionTokens,
			TotalTokens:      raw.Usage.TotalTokens,
		}
	}

	return chunk, nil
}
