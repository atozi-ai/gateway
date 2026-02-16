package openai

import (
	"encoding/json"

	"github.com/atozi-ai/gateway/internal/domain/llm"
)

type openAIChatResponse struct {
	ID                string                 `json:"id"`
	Object            string                 `json:"object"`
	Created           int64                  `json:"created"`
	Model             string                 `json:"model"`
	SystemFingerprint *string                `json:"system_fingerprint,omitempty"`
	Choices           []openAIChoice         `json:"choices"`
	Usage             *openAIUsage           `json:"usage,omitempty"`
	ServiceTier       *string                `json:"service_tier,omitempty"`
}

type openAIChoice struct {
	Index        int                    `json:"index"`
	Message      openAIMessage          `json:"message"`
	FinishReason string                 `json:"finish_reason"`
	Logprobs     *openAILogprobs        `json:"logprobs,omitempty"`
}

type openAIMessage struct {
	Role         string                 `json:"role"`
	Content      string                 `json:"content"`
	Refusal      *string                `json:"refusal,omitempty"`
	Annotations  []interface{}          `json:"annotations,omitempty"`
	ToolCalls    []openAIToolCall       `json:"tool_calls,omitempty"`
}

type openAIToolCall struct {
	ID       string         `json:"id"`
	Type     string         `json:"type"`
	Function openAIFunction `json:"function"`
}

type openAIFunction struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type openAILogprobs struct {
	Content []openAILogprobContent `json:"content,omitempty"`
}

type openAILogprobContent struct {
	Token       string              `json:"token"`
	Logprob     float64             `json:"logprob"`
	Bytes       []int               `json:"bytes,omitempty"`
	TopLogprobs []openAITopLogprob  `json:"top_logprobs,omitempty"`
}

type openAITopLogprob struct {
	Token   string `json:"token"`
	Logprob float64 `json:"logprob"`
	Bytes   []int   `json:"bytes,omitempty"`
}

type openAIUsage struct {
	PromptTokens         int                    `json:"prompt_tokens"`
	CompletionTokens     int                    `json:"completion_tokens"`
	TotalTokens          int                    `json:"total_tokens"`
	PromptTokensDetails  *openAITokensDetails   `json:"prompt_tokens_details,omitempty"`
	CompletionTokensDetails *openAITokensDetails `json:"completion_tokens_details,omitempty"`
}

type openAITokensDetails struct {
	CachedTokens            *int `json:"cached_tokens,omitempty"`
	AudioTokens             *int `json:"audio_tokens,omitempty"`
	ReasoningTokens         *int `json:"reasoning_tokens,omitempty"`
	AcceptedPredictionTokens *int `json:"accepted_prediction_tokens,omitempty"`
	RejectedPredictionTokens *int `json:"rejected_prediction_tokens,omitempty"`
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
	Model             string                    `json:"model"`
	Messages          []llm.Message            `json:"messages"`
	FrequencyPenalty  *float32                  `json:"frequency_penalty,omitempty"`
	LogitBias         map[string]int            `json:"logit_bias,omitempty"`
	Logprobs          *bool                     `json:"logprobs,omitempty"`
	TopLogprobs       *int                      `json:"top_logprobs,omitempty"`
	MaxTokens         *int                      `json:"max_tokens,omitempty"`
	N                 *int                      `json:"n,omitempty"`
	PresencePenalty   *float32                  `json:"presence_penalty,omitempty"`
	ResponseFormat    *openAIResponseFormat     `json:"response_format,omitempty"`
	Seed              *int                      `json:"seed,omitempty"`
	Stop              []string                  `json:"stop,omitempty"`
	Stream            *bool                     `json:"stream,omitempty"`
	StreamOptions     *openAIStreamOptions      `json:"stream_options,omitempty"`
	Temperature       *float32                  `json:"temperature,omitempty"`
	ToolChoice        interface{}               `json:"tool_choice,omitempty"`
	Tools             []openAITool              `json:"tools,omitempty"`
	TopP              *float32                  `json:"top_p,omitempty"`
	User              *string                   `json:"user,omitempty"`
	ParallelToolCalls *bool                     `json:"parallel_tool_calls,omitempty"`
	Verbosity         *llm.Verbosity            `json:"verbosity,omitempty"`
}

type openAITool struct {
	Type     string              `json:"type"`
	Function *openAIFunctionTool `json:"function,omitempty"`
}

type openAIFunctionTool struct {
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	Parameters  json.RawMessage `json:"parameters,omitempty"`
}

type openAIStreamOptions struct {
	IncludeUsage *bool `json:"include_usage,omitempty"`
}

type openAIErrorResponse struct {
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    string `json:"code"`
		Param   string `json:"param,omitempty"`
	} `json:"error"`
}
