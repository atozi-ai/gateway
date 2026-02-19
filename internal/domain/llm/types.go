package llm

import (
	"encoding/json"
	"fmt"
)

type ChatRequest struct {
	Model    string
	Messages []Message
	Options  ChatOptions
	APIKey   string // API key to use for this request (overrides provider's default)
}

type ChatResponse struct {
	ID      string
	Model   string
	Content string
	Raw     json.RawMessage // Raw response from provider
}

type Role string

const (
	RoleSystem    Role = "system"
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
	RoleTool      Role = "tool"
)

type Message struct {
	Role    Role   `json:"role"`
	Content string `json:"content"`
}

type Verbosity string

const (
	LowVerbosity    Verbosity = "low"
	MediumVerbosity Verbosity = "medium"
	HighVerbosity   Verbosity = "high"
)

type ChatOptions struct {
	Temperature *float32
	MaxTokens   *int
	TopP        *float32
	Stop        []string
	Verbosity   *Verbosity

	ResponseFormat *ResponseFormat

	FrequencyPenalty *float32
	PresencePenalty  *float32
	LogitBias        map[string]int
	Logprobs         *bool
	TopLogprobs      *int
	N                *int
	Seed             *int
	User             *string

	Tools             []Tool
	ToolChoice        interface{} // Can be "none", "auto", or ToolChoice object
	ParallelToolCalls *bool
	ToolResolution    *ToolResolution

	Stream        *bool
	StreamOptions *StreamOptions

	// AWS credentials for Bedrock
	AWSAccessKeyID     *string
	AWSSecretAccessKey *string
	AWSRegion          *string

	// GCP credentials for Vertex AI
	GCPProjectID *string
	GCPLocation  *string
}

type ResponseFormat struct {
	Type   string
	Schema json.RawMessage
}

type Tool struct {
	Type     string        `json:"type"`
	Function *FunctionTool `json:"function,omitempty"`
}

type FunctionTool struct {
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	Parameters  json.RawMessage `json:"parameters,omitempty"`
}

type ToolChoice struct {
	Type     string `json:"type"` // "function"
	Function *struct {
		Name string `json:"name"`
	} `json:"function,omitempty"`
}

type ToolResolution struct {
	Type string `json:"type"` // "auto" or "required"
}

type StreamOptions struct {
	IncludeUsage       *bool `json:"include_usage,omitempty"`
	IncludeAccumulated *bool `json:"include_accumulated,omitempty"` // Include accumulated content in each chunk
}

// ProviderError represents an error from a provider with details
type ProviderError struct {
	StatusCode int             `json:"statusCode"`
	Message    string          `json:"message"`
	Type       string          `json:"type,omitempty"`
	Code       string          `json:"code,omitempty"`
	Param      string          `json:"param,omitempty"`
	Raw        json.RawMessage `json:"raw,omitempty"` // Raw error response from provider
}

func (e *ProviderError) Error() string {
	if e.Code != "" {
		return fmt.Sprintf("provider error [%d]: %s (type: %s, code: %s)", e.StatusCode, e.Message, e.Type, e.Code)
	}
	return fmt.Sprintf("provider error [%d]: %s", e.StatusCode, e.Message)
}

// NewProviderError creates a new ProviderError with the given parameters
func NewProviderError(statusCode int, message, errorType, code string) *ProviderError {
	return &ProviderError{
		StatusCode: statusCode,
		Message:    message,
		Type:       errorType,
		Code:       code,
	}
}

// NewValidationError creates a ProviderError for validation failures
func NewValidationError(message, code string) *ProviderError {
	return NewProviderError(400, message, "invalid_request_error", code)
}

// NewUnauthorizedError creates a ProviderError for authentication failures
func NewUnauthorizedError(message string) *ProviderError {
	return NewProviderError(401, message, "authentication_error", "unauthorized")
}

// NewInternalError creates a ProviderError for internal server errors
func NewInternalError(message string) *ProviderError {
	return NewProviderError(500, message, "internal_error", "internal_server_error")
}
