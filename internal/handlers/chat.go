package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/atozi-ai/gateway/internal/domain/llm"
	"github.com/atozi-ai/gateway/internal/platform/logger"
	"github.com/atozi-ai/gateway/internal/providers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
)

type ChatHandler struct{}

func NewChatHandler() *ChatHandler {
	return &ChatHandler{}
}

type ChatRequestPayload struct {
	Model    string              `json:"model"`
	Endpoint string              `json:"endpoint,omitempty"`
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

	Stream             *bool                 `json:"stream,omitempty"`
	StreamOptions      *StreamOptionsPayload `json:"streamOptions,omitempty"`
	Raw                *bool                 `json:"raw,omitempty"`
	IncludeAccumulated *bool                 `json:"includeAccumulated,omitempty"`
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
	IncludeUsage       *bool `json:"includeUsage,omitempty"`
	IncludeAccumulated *bool `json:"includeAccumulated,omitempty"`
}

type ResponseFormatPayload struct {
	Type   string          `json:"type"`
	Schema json.RawMessage `json:"schema,omitempty"`
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
	Content           string          `json:"content,omitempty"`
	Parsed            json.RawMessage `json:"parsed,omitempty"`
	Raw               json.RawMessage `json:"raw,omitempty"`
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
	AccumulatedContent *string           `json:"accumulatedContent,omitempty"`
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

func writeError(w http.ResponseWriter, ctx context.Context, err error) {
	var pe *llm.ProviderError
	if providerErr, ok := err.(*llm.ProviderError); ok {
		pe = providerErr
	} else {
		pe = llm.NewInternalError("Internal server error")
	}

	setRequestIDHeader(w, ctx)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(pe.StatusCode)

	resp := map[string]interface{}{
		"error": map[string]interface{}{
			"message":    pe.Message,
			"type":       pe.Type,
			"code":       pe.Code,
			"param":      pe.Param,
			"statusCode": pe.StatusCode,
		},
	}
	if len(pe.Raw) > 0 {
		resp["raw"] = json.RawMessage(pe.Raw)
	}

	json.NewEncoder(w).Encode(resp)
}

func setRequestIDHeader(w http.ResponseWriter, ctx context.Context) {
	if reqID := middleware.GetReqID(ctx); reqID != "" {
		w.Header().Set("X-Request-Id", reqID)
	}
}

func (h *ChatHandler) Chat(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(r.Context())

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		writeError(w, r.Context(), llm.NewUnauthorizedError("missing Authorization header"))
		return
	}

	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		writeError(w, r.Context(), llm.NewValidationError("invalid Authorization header format", "invalid_auth_format"))
		return
	}
	apiKey := strings.TrimPrefix(authHeader, bearerPrefix)
	if apiKey == "" {
		writeError(w, r.Context(), llm.NewUnauthorizedError("missing API key in Authorization header"))
		return
	}

	var payload ChatRequestPayload

	r.Body = http.MaxBytesReader(w, r.Body, 10*1024*1024)

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Error().Err(err).Msg("Failed to decode request")
		writeError(w, r.Context(), llm.NewValidationError("Invalid request body", "invalid_json"))
		return
	}

	if payload.Model == "" {
		writeError(w, r.Context(), llm.NewValidationError("model is required", "missing_model"))
		return
	}

	if len(payload.Messages) == 0 {
		writeError(w, r.Context(), llm.NewValidationError("messages are required", "missing_messages"))
		return
	}

	if len(payload.Messages) > 1000 {
		writeError(w, r.Context(), llm.NewValidationError("too many messages (max 1000)", "too_many_messages"))
		return
	}

	req := llm.ChatRequest{
		Model:    payload.Model,
		Messages: payload.Messages,
		APIKey:   apiKey,
	}

	streamQueryParam := r.URL.Query().Get("stream")
	if streamQueryParam != "" && (payload.Options == nil || payload.Options.Stream == nil) {
		if payload.Options == nil {
			payload.Options = &ChatOptionsPayload{}
		}
		streamValue := streamQueryParam == "true" || streamQueryParam == "1"
		payload.Options.Stream = &streamValue
	}

	includeRaw := false
	if payload.Options != nil && payload.Options.Raw != nil {
		includeRaw = *payload.Options.Raw
	} else {
		rawQueryParam := r.URL.Query().Get("raw")
		includeRaw = rawQueryParam == "true" || rawQueryParam == "1"
	}

	includeAccumulated := false
	if payload.Options != nil && payload.Options.IncludeAccumulated != nil {
		includeAccumulated = *payload.Options.IncludeAccumulated
	} else {
		accumulatedQueryParam := r.URL.Query().Get("includeAccumulated")
		includeAccumulated = accumulatedQueryParam == "true" || accumulatedQueryParam == "1"
	}

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

		if payload.Options.ResponseFormat != nil {
			req.Options.ResponseFormat = &llm.ResponseFormat{
				Type:   payload.Options.ResponseFormat.Type,
				Schema: payload.Options.ResponseFormat.Schema,
			}
		}

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

		if payload.Options.ToolChoice != nil {
			req.Options.ToolChoice = payload.Options.ToolChoice
		}

		if payload.Options.ToolResolution != nil {
			req.Options.ToolResolution = &llm.ToolResolution{
				Type: payload.Options.ToolResolution.Type,
			}
		}

		if payload.Options.StreamOptions != nil {
			req.Options.StreamOptions = &llm.StreamOptions{
				IncludeUsage:       payload.Options.StreamOptions.IncludeUsage,
				IncludeAccumulated: payload.Options.StreamOptions.IncludeAccumulated,
			}
		}
	}

	provider, model, err := providers.Get(req.Model, apiKey, payload.Endpoint)
	if err != nil {
		log.Error().Err(err).Str("model", req.Model).Msg("Invalid provider/model")
		writeError(w, r.Context(), err)
		return
	}
	req.Model = model

	log.Info().
		Str("provider", provider.Name()).
		Str("model", req.Model).
		Bool("structured", req.Options.ResponseFormat != nil).
		Bool("stream", req.Options.Stream != nil && *req.Options.Stream).
		Msg("Processing chat request")

	isStreaming := req.Options.Stream != nil && *req.Options.Stream

	var ctx context.Context
	var cancel context.CancelFunc

	if isStreaming {
		ctx, cancel = h.createIdleTimeoutContext(r.Context(), 180*time.Second)
	} else {
		ctx, cancel = context.WithTimeout(r.Context(), 180*time.Second)
	}
	defer cancel()

	if isStreaming {
		h.handleStreamingChat(w, ctx, provider, req, log, includeRaw, includeAccumulated)
		return
	}

	resp, err := provider.Chat(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("Chat request failed")
		writeError(w, r.Context(), err)
		return
	}

	type rawResponse struct {
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

	var parsedResponse rawResponse
	if len(resp.Raw) > 0 {
		if err := json.Unmarshal(resp.Raw, &parsedResponse); err != nil {
			log.Warn().Err(err).Msg("Failed to parse raw response, using basic fields")
			parsedResponse.ID = resp.ID
			parsedResponse.Model = resp.Model
			parsedResponse.Object = "chat.completion"
		}
	} else {
		parsedResponse.ID = resp.ID
		parsedResponse.Model = resp.Model
		parsedResponse.Object = "chat.completion"
	}

	choices := make([]ChoicePayload, len(parsedResponse.Choices))
	for i, choice := range parsedResponse.Choices {
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

	var usage *UsagePayload
	if parsedResponse.Usage != nil {
		var promptDetails *TokensDetailsPayload
		if parsedResponse.Usage.PromptTokensDetails != nil {
			promptDetails = &TokensDetailsPayload{
				CachedTokens:             parsedResponse.Usage.PromptTokensDetails.CachedTokens,
				AudioTokens:              parsedResponse.Usage.PromptTokensDetails.AudioTokens,
				ReasoningTokens:          parsedResponse.Usage.PromptTokensDetails.ReasoningTokens,
				AcceptedPredictionTokens: parsedResponse.Usage.PromptTokensDetails.AcceptedPredictionTokens,
				RejectedPredictionTokens: parsedResponse.Usage.PromptTokensDetails.RejectedPredictionTokens,
			}
		}
		var completionDetails *TokensDetailsPayload
		if parsedResponse.Usage.CompletionTokensDetails != nil {
			completionDetails = &TokensDetailsPayload{
				CachedTokens:             parsedResponse.Usage.CompletionTokensDetails.CachedTokens,
				AudioTokens:              parsedResponse.Usage.CompletionTokensDetails.AudioTokens,
				ReasoningTokens:          parsedResponse.Usage.CompletionTokensDetails.ReasoningTokens,
				AcceptedPredictionTokens: parsedResponse.Usage.CompletionTokensDetails.AcceptedPredictionTokens,
				RejectedPredictionTokens: parsedResponse.Usage.CompletionTokensDetails.RejectedPredictionTokens,
			}
		}
		usage = &UsagePayload{
			PromptTokens:            parsedResponse.Usage.PromptTokens,
			CompletionTokens:        parsedResponse.Usage.CompletionTokens,
			TotalTokens:             parsedResponse.Usage.TotalTokens,
			PromptTokensDetails:     promptDetails,
			CompletionTokensDetails: completionDetails,
		}
	}

	response := ChatResponsePayload{
		ID:                parsedResponse.ID,
		Object:            parsedResponse.Object,
		Created:           parsedResponse.Created,
		Model:             parsedResponse.Model,
		SystemFingerprint: parsedResponse.SystemFingerprint,
		Choices:           choices,
		Usage:             usage,
		ServiceTier:       parsedResponse.ServiceTier,
	}

	if includeAccumulated {
		response.Content = resp.Content
	}

	if includeRaw && len(resp.Raw) > 0 {
		response.Raw = resp.Raw
	}

	if req.Options.ResponseFormat != nil && resp.Content != "" {
		var parsed json.RawMessage
		if err := json.Unmarshal([]byte(resp.Content), &parsed); err == nil {
			response.Parsed = parsed
		} else {
			log.Warn().Err(err).Msg("Failed to parse structured response as JSON")
		}
	}

	w.Header().Set("Content-Type", "application/json")
	setRequestIDHeader(w, r.Context())
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Error().Err(err).Msg("Failed to encode response")
		return
	}
}

type idleTimeoutContext struct {
	context.Context
	mu           sync.RWMutex
	lastActivity time.Time
	timeout      time.Duration
	cancel       func()
	once         sync.Once
	doneCh       chan struct{}
}

func (i *idleTimeoutContext) Done() <-chan struct{} {
	i.once.Do(func() {
		i.doneCh = make(chan struct{})
		go func() {
			ticker := time.NewTicker(10 * time.Second)
			defer ticker.Stop()

			for {
				select {
				case <-i.Context.Done():
					close(i.doneCh)
					return
				case <-ticker.C:
					i.mu.RLock()
					since := time.Since(i.lastActivity)
					i.mu.RUnlock()

					if since > i.timeout {
						i.cancel()
						close(i.doneCh)
						return
					}
				}
			}
		}()
	})
	return i.doneCh
}

func (i *idleTimeoutContext) Err() error {
	if i.doneCh == nil {
		return nil
	}
	return i.Context.Err()
}

func (i *idleTimeoutContext) RecordActivity() {
	i.mu.Lock()
	i.lastActivity = time.Now()
	i.mu.Unlock()
}

func (h *ChatHandler) createIdleTimeoutContext(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	baseCtx, cancel := context.WithCancel(ctx)
	return &idleTimeoutContext{
		Context:      baseCtx,
		lastActivity: time.Now(),
		timeout:      timeout,
		cancel:       cancel,
	}, cancel
}

func (h *ChatHandler) handleStreamingChat(
	w http.ResponseWriter,
	ctx context.Context,
	provider llm.Provider,
	req llm.ChatRequest,
	log zerolog.Logger,
	includeRaw bool,
	includeAccumulated bool,
) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")
	setRequestIDHeader(w, ctx)

	flusher, ok := w.(http.Flusher)
	if !ok {
		log.Error().Msg("Streaming not supported by response writer")
		writeError(w, ctx, llm.NewInternalError("Streaming not supported"))
		return
	}

	includeAccumulatedInMessage := req.Options.StreamOptions != nil &&
		req.Options.StreamOptions.IncludeAccumulated != nil &&
		*req.Options.StreamOptions.IncludeAccumulated

	accumulatedContent := make(map[int]string)

	idleCtx, ok := ctx.(*idleTimeoutContext)
	if ok {
		idleCtx.RecordActivity()
	}

	err := provider.ChatStream(ctx, req, func(chunk *llm.StreamChunk) error {
		if ok {
			idleCtx.RecordActivity()
		}

		choices := make([]ChoicePayload, len(chunk.Choices))
		for i, choice := range chunk.Choices {
			message := MessagePayload{
				Role:    "assistant",
				Content: "",
			}

			if choice.Delta.Role != nil {
				message.Role = *choice.Delta.Role
			}

			if choice.Delta.Content != nil {
				deltaContent := *choice.Delta.Content
				message.Content = deltaContent

				accumulatedContent[choice.Index] += deltaContent

				if includeAccumulatedInMessage {
					if accContent := accumulatedContent[choice.Index]; accContent != "" {
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
				TotalTokens:      chunk.Usage.TotalTokens,
			}
		}

		var content string
		if len(choices) > 0 {
			content = accumulatedContent[choices[0].Index]
		} else {
			content = accumulatedContent[0]
		}

		streamResponse := ChatResponsePayload{
			ID:                chunk.ID,
			Object:            chunk.Object,
			Created:           chunk.Created,
			Model:             chunk.Model,
			SystemFingerprint: chunk.SystemFingerprint,
			Choices:           choices,
			Usage:             usage,
			ServiceTier:       chunk.ServiceTier,
		}

		if includeAccumulated {
			streamResponse.Content = content
		}

		if includeRaw && len(chunk.Raw) > 0 {
			streamResponse.Raw = json.RawMessage(chunk.Raw)
		}

		jsonData, err := json.Marshal(streamResponse)
		if err != nil {
			log.Error().Err(err).Msg("Failed to marshal stream chunk")
			return llm.NewInternalError(fmt.Sprintf("failed to marshal stream chunk: %v", err))
		}

		_, err = fmt.Fprintf(w, "data: %s\n\n", jsonData)
		if err != nil {
			log.Error().Err(err).Msg("Failed to write stream chunk")
			return llm.NewInternalError(fmt.Sprintf("failed to write stream chunk: %v", err))
		}

		flusher.Flush()
		return nil
	})

	if err != nil {
		log.Error().Err(err).Msg("Streaming chat request failed")

		var pe *llm.ProviderError
		if providerErr, ok := err.(*llm.ProviderError); ok {
			pe = providerErr
		} else {
			pe = llm.NewInternalError("Internal server error")
		}

		errorResponse := map[string]interface{}{
			"error": map[string]interface{}{
				"message":    pe.Message,
				"type":       pe.Type,
				"code":       pe.Code,
				"param":      pe.Param,
				"statusCode": pe.StatusCode,
			},
		}
		if len(pe.Raw) > 0 {
			errorResponse["raw"] = json.RawMessage(pe.Raw)
		}

		jsonData, _ := json.Marshal(errorResponse)
		fmt.Fprintf(w, "data: %s\n\n", jsonData)
		flusher.Flush()
		return
	}

	fmt.Fprintf(w, "data: [DONE]\n\n")
	flusher.Flush()
}

func (h *ChatHandler) RegisterRoutes(r chi.Router) {
	r.Post("/chat/completions", h.Chat)
}
