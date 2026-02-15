package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/atozi-ai/gateway/internal/domain/llm"
	"github.com/atozi-ai/gateway/internal/platform/logger"
	"github.com/atozi-ai/gateway/internal/providers"
	"github.com/go-chi/chi/v5"
)

type ChatHandler struct{}

func NewChatHandler() *ChatHandler {
	return &ChatHandler{}
}

type ChatRequestPayload struct {
	Model    string              `json:"model"`
	Messages []llm.Message       `json:"messages"`
	Options  *ChatOptionsPayload `json:"options,omitempty"`
}

type ChatOptionsPayload struct {
	Temperature    *float32               `json:"temperature,omitempty"`
	MaxTokens      *int                   `json:"maxTokens,omitempty"`
	TopP           *float32               `json:"topP,omitempty"`
	Stop           []string               `json:"stop,omitempty"`
	Verbosity      *llm.Verbosity         `json:"verbosity,omitempty"`
	ResponseFormat *ResponseFormatPayload `json:"responseFormat,omitempty"`
}

type ResponseFormatPayload struct {
	Type   string          `json:"type"`             // "json_schema" or "json_object"
	Schema json.RawMessage `json:"schema,omitempty"` // JSON schema for structured output
}

type ChatResponsePayload struct {
	ID      string          `json:"id"`
	Model   string          `json:"model"`
	Content string          `json:"content"`
	Parsed  json.RawMessage `json:"parsed,omitempty"` // Parsed structured response if schema was provided
}

func (h *ChatHandler) Chat(w http.ResponseWriter, r *http.Request) {
	// Get logger with request ID from context
	log := logger.FromContext(r.Context())
	
	var payload ChatRequestPayload

	r.Body = http.MaxBytesReader(w, r.Body, 10*1024*1024)

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Error().Err(err).Msg("Failed to decode request")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if payload.Model == "" {
		http.Error(w, "model is required", http.StatusBadRequest)
		return
	}

	if len(payload.Messages) == 0 {
		http.Error(w, "messages are required", http.StatusBadRequest)
		return
	}

	if len(payload.Messages) > 1000 {
		http.Error(w, "too many messages (max 1000)", http.StatusBadRequest)
		return
	}

	// Create chat request
	req := llm.ChatRequest{
		Model:    payload.Model,
		Messages: payload.Messages,
	}

	// Convert options if provided
	if payload.Options != nil {
		req.Options = llm.ChatOptions{
			Temperature: payload.Options.Temperature,
			MaxTokens:   payload.Options.MaxTokens,
			TopP:        payload.Options.TopP,
			Stop:        payload.Options.Stop,
			Verbosity:   payload.Options.Verbosity,
		}

		// Handle response format for structured output
		if payload.Options.ResponseFormat != nil {
			req.Options.ResponseFormat = &llm.ResponseFormat{
				Type:   payload.Options.ResponseFormat.Type,
				Schema: payload.Options.ResponseFormat.Schema,
			}
		}
	}

	// Get provider and call chat
	provider := providers.Get(req.Model)
	log.Info().
		Str("provider", provider.Name()).
		Str("model", req.Model).
		Bool("structured", req.Options.ResponseFormat != nil).
		Msg("Processing chat request")

	ctx, cancel := context.WithTimeout(r.Context(), 120*time.Second)
	defer cancel()
	resp, err := provider.Chat(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("Chat request failed")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Build response
	response := ChatResponsePayload{
		ID:      resp.ID,
		Model:   resp.Model,
		Content: resp.Content,
	}

	// If structured response was requested, parse the JSON content
	if req.Options.ResponseFormat != nil && resp.Content != "" {
		// Try to parse the content as JSON
		var parsed json.RawMessage
		if err := json.Unmarshal([]byte(resp.Content), &parsed); err == nil {
			response.Parsed = parsed
		} else {
			// If parsing fails, log but don't fail the request
			log.Warn().Err(err).Msg("Failed to parse structured response as JSON")
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Error().Err(err).Msg("Failed to encode response")
		// Response may have been partially written, can't use http.Error
		return
	}
}

func (h *ChatHandler) RegisterRoutes(r chi.Router) {
	r.Post("/chat", h.Chat)
}
