package anthropic_compat

import (
	"encoding/json"
)

const defaultMaxTokens = 4096

type messageRole string

const (
	roleUser      messageRole = "user"
	roleAssistant messageRole = "assistant"
	roleSystem    messageRole = "system"
)

type contentBlockType string

const (
	contentTypeText       contentBlockType = "text"
	contentTypeToolUse    contentBlockType = "tool_use"
	contentTypeToolResult contentBlockType = "tool_result"
)

type contentBlock struct {
	Type      contentBlockType `json:"type,omitempty"`
	Text      string           `json:"text,omitempty"`
	ID        string           `json:"id,omitempty"`
	Name      string           `json:"name,omitempty"`
	Input     json.RawMessage  `json:"input,omitempty"`
	ToolUseID string           `json:"tool_use_id,omitempty"`
	Content   string           `json:"content,omitempty"`
}

type message struct {
	Role    messageRole    `json:"role"`
	Content []contentBlock `json:"content"`
}

type tool struct {
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	Parameters  json.RawMessage `json:"parameters,omitempty"`
}

type messageRequest struct {
	Model         string         `json:"model"`
	MaxTokens     int            `json:"max_tokens"`
	Messages      []message      `json:"messages"`
	System        []contentBlock `json:"system,omitempty"`
	Temperature   *float32       `json:"temperature,omitempty"`
	TopP          *float32       `json:"top_p,omitempty"`
	Tools         []tool         `json:"tools,omitempty"`
	ToolChoice    *toolChoice    `json:"tool_choice,omitempty"`
	Stream        *bool          `json:"stream,omitempty"`
	StopSequences []string       `json:"stop_sequences,omitempty"`
}

type toolChoice struct {
	Type string `json:"type"`
	Name string `json:"name,omitempty"`
}

type messageResponse struct {
	ID           string         `json:"id"`
	Type         string         `json:"type"`
	Role         messageRole    `json:"role"`
	Content      []contentBlock `json:"content"`
	Model        string         `json:"model"`
	StopReason   *string        `json:"stop_reason,omitempty"`
	StopSequence *string        `json:"stop_sequence,omitempty"`
	Usage        messageUsage   `json:"usage"`
}

type messageUsage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

type errorResponse struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}
