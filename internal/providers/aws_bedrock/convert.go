package aws_bedrock

import (
	"github.com/atozi-ai/gateway/internal/domain/llm"
)

type ConverseRequest struct {
	Model           string           `json:"model"`
	Messages        []BedrockMessage `json:"messages"`
	System          []SystemContent  `json:"system,omitempty"`
	InferenceConfig InferenceConfig  `json:"inferenceConfig,omitempty"`
}

type BedrockMessage struct {
	Role    string         `json:"role"`
	Content []ContentBlock `json:"content"`
}

type ContentBlock struct {
	Text string `json:"text,omitempty"`
}

type SystemContent struct {
	Text string `json:"text"`
}

type InferenceConfig struct {
	MaxTokens     *int     `json:"maxTokens,omitempty"`
	Temperature   *float64 `json:"temperature,omitempty"`
	TopP          *float64 `json:"topP,omitempty"`
	StopSequences []string `json:"stopSequences,omitempty"`
}

type ConverseResponse struct {
	ID         string `json:"id"`
	Model      string `json:"model"`
	StopReason string `json:"stopReason"`
	Output     Output `json:"output"`
	Usage      Usage  `json:"usage"`
}

type Output struct {
	Message Message `json:"message"`
}

type Message struct {
	Role    string                 `json:"role"`
	Content []ContentBlockResponse `json:"content"`
}

type ContentBlockResponse struct {
	Text string `json:"text"`
}

type Usage struct {
	InputTokens  int `json:"inputTokens"`
	OutputTokens int `json:"outputTokens"`
	TotalTokens  int `json:"totalTokens"`
}

type ConverseStreamChunk struct {
	Type       string `json:"type,omitempty"`
	Delta      *Delta `json:"delta,omitempty"`
	StopReason string `json:"stopReason,omitempty"`
}

type Delta struct {
	Text string `json:"text,omitempty"`
}

func convertToBedrockRequest(req llm.ChatRequest) ConverseRequest {
	messages := make([]BedrockMessage, 0, len(req.Messages))
	for _, msg := range req.Messages {
		content := ContentBlock{Text: msg.Content}
		messages = append(messages, BedrockMessage{
			Role:    string(msg.Role),
			Content: []ContentBlock{content},
		})
	}

	system := make([]SystemContent, 0)
	for _, msg := range req.Messages {
		if msg.Role == llm.RoleSystem {
			system = append(system, SystemContent{Text: msg.Content})
		}
	}

	infConfig := InferenceConfig{}
	if req.Options.Temperature != nil {
		t := float64(*req.Options.Temperature)
		infConfig.Temperature = &t
	}
	if req.Options.TopP != nil {
		p := float64(*req.Options.TopP)
		infConfig.TopP = &p
	}
	if req.Options.MaxTokens != nil {
		infConfig.MaxTokens = req.Options.MaxTokens
	}
	if len(req.Options.Stop) > 0 {
		infConfig.StopSequences = req.Options.Stop
	}

	return ConverseRequest{
		Messages:        messages,
		System:          system,
		InferenceConfig: infConfig,
	}
}

func convertFromBedrockResponse(resp ConverseResponse) *llm.ChatResponse {
	content := ""
	if len(resp.Output.Message.Content) > 0 {
		content = resp.Output.Message.Content[0].Text
	}

	return &llm.ChatResponse{
		ID:      resp.ID,
		Model:   resp.Model,
		Content: content,
	}
}
