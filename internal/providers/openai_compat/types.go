package openaicompat

import (
	"encoding/json"

	"github.com/atozi-ai/gateway/internal/domain/llm"
)

// --- Chat completions request ---

type chatRequest struct {
	Model             string              `json:"model"`
	Messages          []llm.Message       `json:"messages"`
	FrequencyPenalty  *float32            `json:"frequency_penalty,omitempty"`
	LogitBias         map[string]int      `json:"logit_bias,omitempty"`
	Logprobs          *bool               `json:"logprobs,omitempty"`
	TopLogprobs       *int                `json:"top_logprobs,omitempty"`
	MaxTokens         *int                `json:"max_tokens,omitempty"`
	N                 *int                `json:"n,omitempty"`
	PresencePenalty   *float32            `json:"presence_penalty,omitempty"`
	ResponseFormat    *responseFormat     `json:"response_format,omitempty"`
	Seed              *int                `json:"seed,omitempty"`
	Stop              []string            `json:"stop,omitempty"`
	Stream            *bool               `json:"stream,omitempty"`
	StreamOptions     *streamOptions      `json:"stream_options,omitempty"`
	Temperature       *float32            `json:"temperature,omitempty"`
	ToolChoice        interface{}         `json:"tool_choice,omitempty"`
	Tools             []tool              `json:"tools,omitempty"`
	TopP              *float32            `json:"top_p,omitempty"`
	User              *string             `json:"user,omitempty"`
	ParallelToolCalls *bool               `json:"parallel_tool_calls,omitempty"`
	Verbosity         *llm.Verbosity      `json:"verbosity,omitempty"`
}

type responseFormat struct {
	Type       string      `json:"type"`
	JSONSchema *jsonSchema `json:"json_schema,omitempty"`
}

type jsonSchema struct {
	Name   string          `json:"name"`
	Schema json.RawMessage `json:"schema"`
	Strict bool            `json:"strict,omitempty"`
}

type tool struct {
	Type     string        `json:"type"`
	Function *functionTool `json:"function,omitempty"`
}

type functionTool struct {
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	Parameters  json.RawMessage `json:"parameters,omitempty"`
}

type streamOptions struct {
	IncludeUsage *bool `json:"include_usage,omitempty"`
}

// --- Chat completions response ---

type chatResponse struct {
	ID                string    `json:"id"`
	Object            string    `json:"object"`
	Created           int64     `json:"created"`
	Model             string    `json:"model"`
	SystemFingerprint *string   `json:"system_fingerprint,omitempty"`
	Choices           []choice  `json:"choices"`
	Usage             *usage    `json:"usage,omitempty"`
	ServiceTier       *string   `json:"service_tier,omitempty"`
}

type choice struct {
	Index        int       `json:"index"`
	Message      message   `json:"message"`
	FinishReason string    `json:"finish_reason"`
	Logprobs     *logprobs `json:"logprobs,omitempty"`
}

type message struct {
	Role        string        `json:"role"`
	Content     string        `json:"content"`
	Refusal     *string       `json:"refusal,omitempty"`
	Annotations []interface{} `json:"annotations,omitempty"`
	ToolCalls   []toolCall    `json:"tool_calls,omitempty"`
}

type toolCall struct {
	ID       string   `json:"id"`
	Type     string   `json:"type"`
	Function function `json:"function"`
}

type function struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type logprobs struct {
	Content []logprobContent `json:"content,omitempty"`
}

type logprobContent struct {
	Token       string        `json:"token"`
	Logprob     float64       `json:"logprob"`
	Bytes       []int         `json:"bytes,omitempty"`
	TopLogprobs []topLogprob  `json:"top_logprobs,omitempty"`
}

type topLogprob struct {
	Token   string  `json:"token"`
	Logprob float64 `json:"logprob"`
	Bytes   []int   `json:"bytes,omitempty"`
}

type usage struct {
	PromptTokens            int            `json:"prompt_tokens"`
	CompletionTokens        int            `json:"completion_tokens"`
	TotalTokens             int            `json:"total_tokens"`
	PromptTokensDetails     *tokensDetails `json:"prompt_tokens_details,omitempty"`
	CompletionTokensDetails *tokensDetails `json:"completion_tokens_details,omitempty"`
}

type tokensDetails struct {
	CachedTokens             *int `json:"cached_tokens,omitempty"`
	AudioTokens              *int `json:"audio_tokens,omitempty"`
	ReasoningTokens          *int `json:"reasoning_tokens,omitempty"`
	AcceptedPredictionTokens *int `json:"accepted_prediction_tokens,omitempty"`
	RejectedPredictionTokens *int `json:"rejected_prediction_tokens,omitempty"`
}

// --- Error response ---

type errorResponse struct {
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    string `json:"code"`
		Param   string `json:"param,omitempty"`
	} `json:"error"`
}

// --- Streaming types ---

type streamChunk struct {
	ID      string         `json:"id"`
	Object  string         `json:"object"`
	Created int64          `json:"created"`
	Model   string         `json:"model"`
	Choices []streamChoice `json:"choices"`
	Usage   *streamUsage   `json:"usage,omitempty"`
}

type streamChoice struct {
	Index        int          `json:"index"`
	Delta        streamDelta  `json:"delta"`
	FinishReason *string      `json:"finish_reason,omitempty"`
	Logprobs     interface{}  `json:"logprobs,omitempty"`
}

type streamDelta struct {
	Role      *string       `json:"role,omitempty"`
	Content   *string       `json:"content,omitempty"`
	ToolCalls []interface{} `json:"tool_calls,omitempty"`
}

type streamUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}
