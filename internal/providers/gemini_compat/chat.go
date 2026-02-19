package gemini_compat

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
	maxResponseBytes = 10 * 1024 * 1024
)

func (c *Client) Chat(ctx context.Context, req llm.ChatRequest) (*llm.ChatResponse, error) {
	body := toRequest(req)

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, llm.NewInternalError(fmt.Sprintf("failed to marshal request: %v", err))
	}

	log := logger.FromContext(ctx)
	log.Info().Str("model", req.Model).Msg("Using Gemini Model")

	model := strings.TrimPrefix(req.Model, "gemini/")
	url := c.endpoint(model)

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, llm.NewInternalError(fmt.Sprintf("failed to create request: %v", err))
	}
	c.setHeaders(httpReq, req.APIKey)

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

	var raw GenerateContentResponse
	if err := json.Unmarshal(respBody, &raw); err != nil {
		return nil, llm.NewInternalError(fmt.Sprintf("failed to unmarshal response: %v", err))
	}

	chatResp := toChatResponse(raw, model)
	chatResp.Raw = respBody

	return chatResp, nil
}

func (c *Client) ChatStream(ctx context.Context, req llm.ChatRequest, callback func(*llm.StreamChunk) error) error {
	stream := true
	req.Options.Stream = &stream

	body := toRequest(req)
	body.GenerationConfig.StopSequences = nil

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return llm.NewInternalError(fmt.Sprintf("failed to marshal request: %v", err))
	}

	log := logger.FromContext(ctx)
	log.Info().Str("model", req.Model).Msg("Using Gemini Model for streaming")

	model := strings.TrimPrefix(req.Model, "gemini/")
	url := c.streamEndpoint(model)

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(jsonBody))
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

func (c *Client) endpoint(model string) string {
	baseURL := c.cfg.BaseURL
	if baseURL == "" {
		baseURL = "https://generativelanguage.googleapis.com"
	}
	return fmt.Sprintf("%s/v1beta/models/%s:generateContent", baseURL, model)
}

func (c *Client) streamEndpoint(model string) string {
	baseURL := c.cfg.BaseURL
	if baseURL == "" {
		baseURL = "https://generativelanguage.googleapis.com"
	}
	return fmt.Sprintf("%s/v1beta/models/%s:streamGenerateContent", baseURL, model)
}

func checkError(statusCode int, body []byte) error {
	if statusCode >= 200 && statusCode < 300 {
		return nil
	}

	var apiErr ErrorResponse
	if err := json.Unmarshal(body, &apiErr); err == nil && apiErr.Error.Message != "" {
		return &llm.ProviderError{
			StatusCode: statusCode,
			Message:    apiErr.Error.Message,
			Type:       "api_error",
			Code:       apiErr.Error.Status,
			Raw:        body,
		}
	}

	return &llm.ProviderError{
		StatusCode: statusCode,
		Message:    fmt.Sprintf("API returned status %d", statusCode),
		Raw:        body,
	}
}

func readSSEStream(ctx context.Context, r io.Reader, callback func(*llm.StreamChunk) error) error {
	const maxScanTokenSize = 64 * 1024
	scanner := bufio.NewScanner(r)
	buf := make([]byte, 0, bufio.MaxScanTokenSize)
	scanner.Buffer(buf, maxScanTokenSize)

	log := logger.FromContext(ctx)

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
			return nil
		}

		var chunk map[string]interface{}
		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			log.Warn().Err(err).Str("data", data).Msg("Failed to parse SSE chunk")
			continue
		}

		candidates, ok := chunk["candidates"].([]interface{})
		if !ok || len(candidates) == 0 {
			continue
		}

		candidate, ok := candidates[0].(map[string]interface{})
		if !ok {
			continue
		}

		content, ok := candidate["content"].(map[string]interface{})
		if !ok {
			continue
		}

		parts, ok := content["parts"].([]interface{})
		if !ok || len(parts) == 0 {
			continue
		}

		part, ok := parts[0].(map[string]interface{})
		if !ok {
			continue
		}

		text, ok := part["text"].(string)
		if !ok || text == "" {
			continue
		}

		streamChunk := &llm.StreamChunk{
			ID:      "",
			Model:   "",
			Choices: []llm.StreamChoice{},
		}

		textPtr := text
		streamChunk.Choices = append(streamChunk.Choices, llm.StreamChoice{
			Index: 0,
			Delta: llm.StreamDelta{
				Content: &textPtr,
			},
		})

		if err := callback(streamChunk); err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return llm.NewInternalError(fmt.Sprintf("failed to read stream: %v", err))
	}

	return nil
}
