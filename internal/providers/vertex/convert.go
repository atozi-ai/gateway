package vertex

import (
	"github.com/atozi-ai/gateway/internal/domain/llm"
)

type VertexRequest struct {
	Model       string          `json:"model"`
	Messages    []VertexMessage `json:"messages"`
	Temperature *float32        `json:"temperature,omitempty"`
	MaxTokens   *int            `json:"max_tokens,omitempty"`
	TopP        *float32        `json:"top_p,omitempty"`
	Stop        []string        `json:"stop,omitempty"`
	Stream      bool            `json:"stream,omitempty"`
}

type VertexMessage struct {
	Role    string          `json:"role"`
	Content []VertexContent `json:"content"`
}

type VertexContent struct {
	Type     string `json:"type"`
	Text     string `json:"text,omitempty"`
	ImageURL string `json:"image_url,omitempty"`
}

type VertexResponse struct {
	ID      string         `json:"id"`
	Model   string         `json:"model"`
	Choices []VertexChoice `json:"choices"`
	Usage   VertexUsage    `json:"usage"`
}

type VertexChoice struct {
	Index        int           `json:"index"`
	Message      VertexMessage `json:"message"`
	FinishReason string        `json:"finish_reason"`
}

type VertexUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

func convertToVertexRequest(req llm.ChatRequest) VertexRequest {
	messages := make([]VertexMessage, 0, len(req.Messages))
	for _, msg := range req.Messages {
		content := VertexContent{
			Type: "text",
			Text: msg.Content,
		}
		messages = append(messages, VertexMessage{
			Role:    string(msg.Role),
			Content: []VertexContent{content},
		})
	}

	vertexReq := VertexRequest{
		Model:    req.Model,
		Messages: messages,
	}

	if req.Options.Temperature != nil {
		vertexReq.Temperature = req.Options.Temperature
	}
	if req.Options.MaxTokens != nil {
		vertexReq.MaxTokens = req.Options.MaxTokens
	}
	if req.Options.TopP != nil {
		vertexReq.TopP = req.Options.TopP
	}
	if len(req.Options.Stop) > 0 {
		vertexReq.Stop = req.Options.Stop
	}

	return vertexReq
}

func convertFromVertexResponse(resp VertexResponse) *llm.ChatResponse {
	content := ""

	if len(resp.Choices) > 0 {
		choice := resp.Choices[0]
		if len(choice.Message.Content) > 0 {
			content = choice.Message.Content[0].Text
		}
	}

	return &llm.ChatResponse{
		ID:      resp.ID,
		Model:   resp.Model,
		Content: content,
	}
}
