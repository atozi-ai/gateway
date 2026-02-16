package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/atozi-ai/gateway/internal/domain/llm"
	"github.com/atozi-ai/gateway/internal/platform/logger"
	"github.com/atozi-ai/gateway/internal/providers"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
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
	Temperature *float32       `json:"temperature,omitempty"`
	MaxTokens   *int           `json:"maxTokens,omitempty"`
	TopP        *float32       `json:"topP,omitempty"`
	Stop        []string       `json:"stop,omitempty"`
	Verbosity   *llm.Verbosity `json:"verbosity,omitempty"`

	ResponseFormat *ResponseFormatPayload `json:"responseFormat,omitempty"`

	FrequencyPenalty *float32       `json:"frequencyPenalty,omitempty"`
	PresencePenalty  *float32       `json:"presencePenalty,omitempty"`
	LogitBias        map[string]int `json:"logitBias,omitempty"`
	Logprobs         *bool          `json:"logprobs,omitempty"`
	TopLogprobs      *int           `json:"topLogprobs,omitempty"`
	N                *int           `json:"n,omitempty"`
	Seed             *int           `json:"seed,omitempty"`
	User             *string        `json:"user,omitempty"`

	Tools             []ToolPayload          `json:"tools,omitempty"`
	ToolChoice        interface{}            `json:"toolChoice,omitempty"`
	ParallelToolCalls *bool                  `json:"parallelToolCalls,omitempty"`
	ToolResolution    *ToolResolutionPayload `json:"toolResolution,omitempty"`

	Stream        *bool                 `json:"stream,omitempty"`
	StreamOptions *StreamOptionsPayload `json:"streamOptions,omitempty"`
	Raw                *bool `json:"raw,omitempty"`                 // Include raw provider response
	IncludeAccumulated *bool `json:"includeAccumulated,omitempty"`   // Include top-level accumulated content field
}

type ToolPayload struct {
	Type     string           `json:"type"`
	Function *FunctionPayload `json:"function,omitempty"`
}

type FunctionPayload struct {
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	Parameters  json.RawMessage `json:"parameters,omitempty"`
}

type ToolResolutionPayload struct {
	Type string `json:"type"`
}

type StreamOptionsPayload struct {
	IncludeUsage        *bool `json:"includeUsage,omitempty"`
	IncludeAccumulated *bool `json:"includeAccumulated,omitempty"` // Include accumulated content in each chunk
}

type ResponseFormatPayload struct {
	Type   string          `json:"type"`             // "json_schema" or "json_object"
	Schema json.RawMessage `json:"schema,omitempty"` // JSON schema for structured output
}

type ChatResponsePayload struct {
	ID                string          `json:"id"`
	Object            string          `json:"object"`
	Created           int64           `json:"created"`
	Model             string          `json:"model"`
	SystemFingerprint *string         `json:"systemFingerprint,omitempty"`
	Choices           []ChoicePayload `json:"choices"`
	Usage             *UsagePayload   `json:"usage,omitempty"`
	ServiceTier       *string         `json:"serviceTier,omitempty"`
	Content           string          `json:"content,omitempty"`          // Convenience field for first choice content (accumulated)
	Parsed            json.RawMessage `json:"parsed,omitempty"` // Parsed structured response if schema was provided
	Raw               json.RawMessage `json:"raw,omitempty"`    // Raw response from provider
}

type ChoicePayload struct {
	Index        int              `json:"index"`
	Message      MessagePayload   `json:"message"`
	FinishReason string           `json:"finishReason"`
	Logprobs     *LogprobsPayload `json:"logprobs,omitempty"`
}

type MessagePayload struct {
	Role               string            `json:"role"`
	Content            string            `json:"content"`
	AccumulatedContent *string           `json:"accumulatedContent,omitempty"` // Full accumulated content from all chunks
	Refusal            *string           `json:"refusal,omitempty"`
	Annotations        []interface{}     `json:"annotations,omitempty"`
	ToolCalls          []ToolCallPayload `json:"toolCalls,omitempty"`
}

type ToolCallPayload struct {
	ID       string              `json:"id"`
	Type     string              `json:"type"`
	Function FunctionCallPayload `json:"function"`
}

type FunctionCallPayload struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type LogprobsPayload struct {
	Content []LogprobContentPayload `json:"content,omitempty"`
}

type LogprobContentPayload struct {
	Token       string              `json:"token"`
	Logprob     float64             `json:"logprob"`
	Bytes       []int               `json:"bytes,omitempty"`
	TopLogprobs []TopLogprobPayload `json:"topLogprobs,omitempty"`
}

type TopLogprobPayload struct {
	Token   string  `json:"token"`
	Logprob float64 `json:"logprob"`
	Bytes   []int   `json:"bytes,omitempty"`
}

type UsagePayload struct {
	PromptTokens            int                   `json:"promptTokens"`
	CompletionTokens        int                   `json:"completionTokens"`
	TotalTokens             int                   `json:"totalTokens"`
	PromptTokensDetails     *TokensDetailsPayload `json:"promptTokensDetails,omitempty"`
	CompletionTokensDetails *TokensDetailsPayload `json:"completionTokensDetails,omitempty"`
}

type TokensDetailsPayload struct {
	CachedTokens             *int `json:"cachedTokens,omitempty"`
	AudioTokens              *int `json:"audioTokens,omitempty"`
	ReasoningTokens          *int `json:"reasoningTokens,omitempty"`
	AcceptedPredictionTokens *int `json:"acceptedPredictionTokens,omitempty"`
	RejectedPredictionTokens *int `json:"rejectedPredictionTokens,omitempty"`
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

	// Check for stream parameter in query string as fallback
	streamQueryParam := r.URL.Query().Get("stream")
	if streamQueryParam != "" && (payload.Options == nil || payload.Options.Stream == nil) {
		// Initialize Options if nil
		if payload.Options == nil {
			payload.Options = &ChatOptionsPayload{}
		}
		// Parse stream query parameter
		streamValue := streamQueryParam == "true" || streamQueryParam == "1"
		payload.Options.Stream = &streamValue
	}

	// Check for raw parameter - from options first, then query string as fallback
	includeRaw := false
	if payload.Options != nil && payload.Options.Raw != nil {
		includeRaw = *payload.Options.Raw
	} else {
		rawQueryParam := r.URL.Query().Get("raw")
		includeRaw = rawQueryParam == "true" || rawQueryParam == "1"
	}

	// Check for includeAccumulated parameter - from options first, then query string as fallback
	includeAccumulated := false
	if payload.Options != nil && payload.Options.IncludeAccumulated != nil {
		includeAccumulated = *payload.Options.IncludeAccumulated
	} else {
		accumulatedQueryParam := r.URL.Query().Get("includeAccumulated")
		includeAccumulated = accumulatedQueryParam == "true" || accumulatedQueryParam == "1"
	}

	// Convert options if provided
	if payload.Options != nil {
		req.Options = llm.ChatOptions{
			Temperature:       payload.Options.Temperature,
			MaxTokens:         payload.Options.MaxTokens,
			TopP:              payload.Options.TopP,
			Stop:              payload.Options.Stop,
			Verbosity:         payload.Options.Verbosity,
			FrequencyPenalty:  payload.Options.FrequencyPenalty,
			PresencePenalty:   payload.Options.PresencePenalty,
			LogitBias:         payload.Options.LogitBias,
			Logprobs:          payload.Options.Logprobs,
			TopLogprobs:       payload.Options.TopLogprobs,
			N:                 payload.Options.N,
			Seed:              payload.Options.Seed,
			User:              payload.Options.User,
			Stream:            payload.Options.Stream,
			ParallelToolCalls: payload.Options.ParallelToolCalls,
		}

		// Handle response format for structured output
		if payload.Options.ResponseFormat != nil {
			req.Options.ResponseFormat = &llm.ResponseFormat{
				Type:   payload.Options.ResponseFormat.Type,
				Schema: payload.Options.ResponseFormat.Schema,
			}
		}

		// Handle tools
		if len(payload.Options.Tools) > 0 {
			req.Options.Tools = make([]llm.Tool, len(payload.Options.Tools))
			for i, tool := range payload.Options.Tools {
				req.Options.Tools[i] = llm.Tool{
					Type: tool.Type,
				}
				if tool.Function != nil {
					req.Options.Tools[i].Function = &llm.FunctionTool{
						Name:        tool.Function.Name,
						Description: tool.Function.Description,
						Parameters:  tool.Function.Parameters,
					}
				}
			}
		}

		// Handle tool choice
		if payload.Options.ToolChoice != nil {
			req.Options.ToolChoice = payload.Options.ToolChoice
		}

		// Handle tool resolution
		if payload.Options.ToolResolution != nil {
			req.Options.ToolResolution = &llm.ToolResolution{
				Type: payload.Options.ToolResolution.Type,
			}
		}

		// Handle stream options
		if payload.Options.StreamOptions != nil {
			req.Options.StreamOptions = &llm.StreamOptions{
				IncludeUsage:        payload.Options.StreamOptions.IncludeUsage,
				IncludeAccumulated:  payload.Options.StreamOptions.IncludeAccumulated,
			}
		}
	}

	// Get provider and call chat
	provider := providers.Get(req.Model)
	log.Info().
		Str("provider", provider.Name()).
		Str("model", req.Model).
		Bool("structured", req.Options.ResponseFormat != nil).
		Bool("stream", req.Options.Stream != nil && *req.Options.Stream).
		Msg("Processing chat request")

	ctx, cancel := context.WithTimeout(r.Context(), 120*time.Second)
	defer cancel()

	// Check if streaming is requested
	if req.Options.Stream != nil && *req.Options.Stream {
		h.handleStreamingChat(w, r, ctx, provider, req, log, includeRaw, includeAccumulated)
		return
	}

	resp, err := provider.Chat(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("Chat request failed")
		
		// Check if it's a ProviderError and return it to the user
		if providerErr, ok := err.(*llm.ProviderError); ok {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(providerErr.StatusCode)
			
			errorResponse := map[string]interface{}{
				"error": map[string]interface{}{
					"message":    providerErr.Message,
					"type":       providerErr.Type,
					"code":       providerErr.Code,
					"param":      providerErr.Param,
					"statusCode":  providerErr.StatusCode,
				},
			}
			
			// Include raw error if available
			if len(providerErr.Raw) > 0 {
				errorResponse["raw"] = json.RawMessage(providerErr.Raw)
			}
			
			if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
				log.Error().Err(err).Msg("Failed to encode error response")
			}
			return
		}
		
		// For other errors, return generic internal server error
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Parse raw response to extract all fields
	var rawResponse struct {
		ID                string  `json:"id"`
		Object            string  `json:"object"`
		Created           int64   `json:"created"`
		Model             string  `json:"model"`
		SystemFingerprint *string `json:"system_fingerprint,omitempty"`
		Choices           []struct {
			Index   int `json:"index"`
			Message struct {
				Role        string        `json:"role"`
				Content     string        `json:"content"`
				Refusal     *string       `json:"refusal,omitempty"`
				Annotations []interface{} `json:"annotations,omitempty"`
				ToolCalls   []struct {
					ID       string `json:"id"`
					Type     string `json:"type"`
					Function struct {
						Name      string `json:"name"`
						Arguments string `json:"arguments"`
					} `json:"function"`
				} `json:"tool_calls,omitempty"`
			} `json:"message"`
			FinishReason string `json:"finish_reason"`
			Logprobs     *struct {
				Content []struct {
					Token       string  `json:"token"`
					Logprob     float64 `json:"logprob"`
					Bytes       []int   `json:"bytes,omitempty"`
					TopLogprobs []struct {
						Token   string  `json:"token"`
						Logprob float64 `json:"logprob"`
						Bytes   []int   `json:"bytes,omitempty"`
					} `json:"top_logprobs,omitempty"`
				} `json:"content,omitempty"`
			} `json:"logprobs,omitempty"`
		} `json:"choices"`
		Usage *struct {
			PromptTokens        int `json:"prompt_tokens"`
			CompletionTokens    int `json:"completion_tokens"`
			TotalTokens         int `json:"total_tokens"`
			PromptTokensDetails *struct {
				CachedTokens             *int `json:"cached_tokens,omitempty"`
				AudioTokens              *int `json:"audio_tokens,omitempty"`
				ReasoningTokens          *int `json:"reasoning_tokens,omitempty"`
				AcceptedPredictionTokens *int `json:"accepted_prediction_tokens,omitempty"`
				RejectedPredictionTokens *int `json:"rejected_prediction_tokens,omitempty"`
			} `json:"prompt_tokens_details,omitempty"`
			CompletionTokensDetails *struct {
				CachedTokens             *int `json:"cached_tokens,omitempty"`
				AudioTokens              *int `json:"audio_tokens,omitempty"`
				ReasoningTokens          *int `json:"reasoning_tokens,omitempty"`
				AcceptedPredictionTokens *int `json:"accepted_prediction_tokens,omitempty"`
				RejectedPredictionTokens *int `json:"rejected_prediction_tokens,omitempty"`
			} `json:"completion_tokens_details,omitempty"`
		} `json:"usage,omitempty"`
		ServiceTier *string `json:"service_tier,omitempty"`
	}

	// Try to parse the raw response
	if len(resp.Raw) > 0 {
		if err := json.Unmarshal(resp.Raw, &rawResponse); err != nil {
			log.Warn().Err(err).Msg("Failed to parse raw response, using basic fields")
			// Fallback to basic response
			rawResponse.ID = resp.ID
			rawResponse.Model = resp.Model
			rawResponse.Object = "chat.completion"
		}
	} else {
		// Fallback if raw is empty
		rawResponse.ID = resp.ID
		rawResponse.Model = resp.Model
		rawResponse.Object = "chat.completion"
	}

	// Convert choices from raw response format to payload format
	choices := make([]ChoicePayload, len(rawResponse.Choices))
	for i, choice := range rawResponse.Choices {
		toolCalls := make([]ToolCallPayload, len(choice.Message.ToolCalls))
		for j, tc := range choice.Message.ToolCalls {
			toolCalls[j] = ToolCallPayload{
				ID:   tc.ID,
				Type: tc.Type,
				Function: FunctionCallPayload{
					Name:      tc.Function.Name,
					Arguments: tc.Function.Arguments,
				},
			}
		}

		var logprobs *LogprobsPayload
		if choice.Logprobs != nil {
			logprobContent := make([]LogprobContentPayload, len(choice.Logprobs.Content))
			for j, lp := range choice.Logprobs.Content {
				topLogprobs := make([]TopLogprobPayload, len(lp.TopLogprobs))
				for k, tlp := range lp.TopLogprobs {
					topLogprobs[k] = TopLogprobPayload{
						Token:   tlp.Token,
						Logprob: tlp.Logprob,
						Bytes:   tlp.Bytes,
					}
				}
				logprobContent[j] = LogprobContentPayload{
					Token:       lp.Token,
					Logprob:     lp.Logprob,
					Bytes:       lp.Bytes,
					TopLogprobs: topLogprobs,
				}
			}
			logprobs = &LogprobsPayload{
				Content: logprobContent,
			}
		}

		choices[i] = ChoicePayload{
			Index:        choice.Index,
			FinishReason: choice.FinishReason,
			Logprobs:     logprobs,
			Message: MessagePayload{
				Role:        choice.Message.Role,
				Content:     choice.Message.Content,
				Refusal:     choice.Message.Refusal,
				Annotations: choice.Message.Annotations,
				ToolCalls:   toolCalls,
			},
		}
	}

	// Convert usage
	var usage *UsagePayload
	if rawResponse.Usage != nil {
		var promptDetails *TokensDetailsPayload
		if rawResponse.Usage.PromptTokensDetails != nil {
			promptDetails = &TokensDetailsPayload{
				CachedTokens:             rawResponse.Usage.PromptTokensDetails.CachedTokens,
				AudioTokens:              rawResponse.Usage.PromptTokensDetails.AudioTokens,
				ReasoningTokens:          rawResponse.Usage.PromptTokensDetails.ReasoningTokens,
				AcceptedPredictionTokens: rawResponse.Usage.PromptTokensDetails.AcceptedPredictionTokens,
				RejectedPredictionTokens: rawResponse.Usage.PromptTokensDetails.RejectedPredictionTokens,
			}
		}
		var completionDetails *TokensDetailsPayload
		if rawResponse.Usage.CompletionTokensDetails != nil {
			completionDetails = &TokensDetailsPayload{
				CachedTokens:             rawResponse.Usage.CompletionTokensDetails.CachedTokens,
				AudioTokens:              rawResponse.Usage.CompletionTokensDetails.AudioTokens,
				ReasoningTokens:          rawResponse.Usage.CompletionTokensDetails.ReasoningTokens,
				AcceptedPredictionTokens: rawResponse.Usage.CompletionTokensDetails.AcceptedPredictionTokens,
				RejectedPredictionTokens: rawResponse.Usage.CompletionTokensDetails.RejectedPredictionTokens,
			}
		}
		usage = &UsagePayload{
			PromptTokens:            rawResponse.Usage.PromptTokens,
			CompletionTokens:        rawResponse.Usage.CompletionTokens,
			TotalTokens:             rawResponse.Usage.TotalTokens,
			PromptTokensDetails:     promptDetails,
			CompletionTokensDetails: completionDetails,
		}
	}

	// Build response
	response := ChatResponsePayload{
		ID:                rawResponse.ID,
		Object:            rawResponse.Object,
		Created:           rawResponse.Created,
		Model:             rawResponse.Model,
		SystemFingerprint: rawResponse.SystemFingerprint,
		Choices:           choices,
		Usage:             usage,
		ServiceTier:       rawResponse.ServiceTier,
	}
	
	// Include top-level content field (accumulated) if requested
	if includeAccumulated {
		response.Content = resp.Content // Convenience field
	}
	
	// Include raw provider response if requested
	if includeRaw && len(resp.Raw) > 0 {
		response.Raw = resp.Raw
	}

	// If structured response was requested, parse the content as JSON
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

func (h *ChatHandler) handleStreamingChat(
	w http.ResponseWriter,
	r *http.Request,
	ctx context.Context,
	provider llm.Provider,
	req llm.ChatRequest,
	log zerolog.Logger,
	includeRaw bool,
	includeAccumulated bool,
) {
	// Set up Server-Sent Events headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no") // Disable buffering in nginx

	// Create a flusher to enable streaming
	flusher, ok := w.(http.Flusher)
	if !ok {
		log.Error().Msg("Streaming not supported by response writer")
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	// Check if accumulated content is requested in message payload
	includeAccumulatedInMessage := req.Options.StreamOptions != nil && 
		req.Options.StreamOptions.IncludeAccumulated != nil && 
		*req.Options.StreamOptions.IncludeAccumulated

	// Always track accumulated content for top-level content field
	accumulatedContent := make(map[int]string)

	// Stream chunks from provider
	err := provider.ChatStream(ctx, req, func(chunk *llm.StreamChunk) error {
		// Convert chunk to response format
		choices := make([]ChoicePayload, len(chunk.Choices))
		for i, choice := range chunk.Choices {
			message := MessagePayload{
				Role:    "assistant",
				Content: "",
			}

			if choice.Delta.Role != nil {
				message.Role = *choice.Delta.Role
			}
			
			// Handle delta content
			if choice.Delta.Content != nil {
				deltaContent := *choice.Delta.Content
				message.Content = deltaContent // Delta content for this chunk only
				
				// Always accumulate content for top-level field
				accumulatedContent[choice.Index] += deltaContent
				
				// Set accumulated content in message if requested
				if includeAccumulatedInMessage {
					if accContent, exists := accumulatedContent[choice.Index]; exists && accContent != "" {
						message.AccumulatedContent = &accContent
					}
				}
			}

			choices[i] = ChoicePayload{
				Index:   choice.Index,
				Message: message,
			}

			if choice.FinishReason != nil {
				choices[i].FinishReason = *choice.FinishReason
			}
		}

		var usage *UsagePayload
		if chunk.Usage != nil {
			usage = &UsagePayload{
				PromptTokens:     chunk.Usage.PromptTokens,
				CompletionTokens: chunk.Usage.CompletionTokens,
				TotalTokens:       chunk.Usage.TotalTokens,
			}
		}

		// Get accumulated content for top-level content field
		// Use first choice index if available, otherwise use index 0 (most common case)
		var content string
		if len(choices) > 0 {
			if accContent, exists := accumulatedContent[choices[0].Index]; exists {
				content = accContent
			}
		} else {
			// If no choices in this chunk (e.g., final usage chunk), use accumulated content from index 0
			if accContent, exists := accumulatedContent[0]; exists {
				content = accContent
			}
		}

		streamResponse := ChatResponsePayload{
			ID:      chunk.ID,
			Object:  chunk.Object,
			Created: chunk.Created,
			Model:   chunk.Model,
			Choices: choices,
			Usage:   usage,
		}
		
		// Include top-level content field (accumulated) if requested
		if includeAccumulated {
			streamResponse.Content = content // Accumulated content from all previous chunks
		}
		
		// Include raw provider response if requested
		if includeRaw && len(chunk.Raw) > 0 {
			streamResponse.Raw = json.RawMessage(chunk.Raw)
		}

		// Marshal to JSON
		jsonData, err := json.Marshal(streamResponse)
		if err != nil {
			log.Error().Err(err).Msg("Failed to marshal stream chunk")
			return err
		}

		// Write SSE format: "data: {json}\n\n"
		_, err = fmt.Fprintf(w, "data: %s\n\n", jsonData)
		if err != nil {
			log.Error().Err(err).Msg("Failed to write stream chunk")
			return err
		}

		// Flush to send data immediately
		flusher.Flush()
		return nil
	})

	if err != nil {
		log.Error().Err(err).Msg("Streaming chat request failed")

		// Check if it's a ProviderError
		if providerErr, ok := err.(*llm.ProviderError); ok {
			errorResponse := map[string]interface{}{
				"error": map[string]interface{}{
					"message":    providerErr.Message,
					"type":       providerErr.Type,
					"code":       providerErr.Code,
					"param":      providerErr.Param,
					"statusCode":  providerErr.StatusCode,
				},
			}

			jsonData, _ := json.Marshal(errorResponse)
			fmt.Fprintf(w, "data: %s\n\n", jsonData)
			flusher.Flush()
			return
		}

		// For other errors, send error in SSE format
		errorResponse := map[string]interface{}{
			"error": map[string]interface{}{
				"message": "Internal server error",
				"type":    "internal_error",
			},
		}
		jsonData, _ := json.Marshal(errorResponse)
		fmt.Fprintf(w, "data: %s\n\n", jsonData)
		flusher.Flush()
		return
	}

	// Send [DONE] marker to indicate stream completion
	fmt.Fprintf(w, "data: [DONE]\n\n")
	flusher.Flush()
}

func (h *ChatHandler) RegisterRoutes(r chi.Router) {
	r.Post("/chat", h.Chat)
}
