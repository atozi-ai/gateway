package llm

import "encoding/json"

type ChatRequest struct {
	Model    string
	Messages []Message
	Options  ChatOptions
}

type ChatResponse struct {
	ID      string
	Model   string
	Content string
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
}

type ResponseFormat struct {
	Type   string
	Schema json.RawMessage
}
