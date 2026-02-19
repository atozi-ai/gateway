package gemini_compat

import (
	"encoding/json"

	"github.com/atozi-ai/gateway/internal/domain/llm"
)

func toRequest(req llm.ChatRequest) GenerateContentRequest {
	opts := req.Options

	maxTokens := defaultMaxTokens
	if opts.MaxTokens != nil && *opts.MaxTokens > 0 {
		maxTokens = *opts.MaxTokens
	}

	var contents []Content
	var systemInstruction *Content

	for _, msg := range req.Messages {
		switch msg.Role {
		case llm.RoleSystem:
			systemInstruction = &Content{
				Role:  "system",
				Parts: []Part{{Text: msg.Content}},
			}
		case llm.RoleUser:
			contents = append(contents, Content{
				Role:  "user",
				Parts: []Part{{Text: msg.Content}},
			})
		case llm.RoleAssistant:
			contents = append(contents, Content{
				Role:  "model",
				Parts: []Part{{Text: msg.Content}},
			})
		}
	}

	var tools []Tool
	if len(opts.Tools) > 0 {
		tools = make([]Tool, len(opts.Tools))
		for i, t := range opts.Tools {
			tools[i] = Tool{
				FunctionDeclarations: []FunctionDeclaration{
					{
						Name:        t.Function.Name,
						Description: t.Function.Description,
						Parameters:  t.Function.Parameters,
					},
				},
			}
		}
	}

	genConfig := &GenerationConfig{
		MaxOutputTokens: &maxTokens,
		Temperature:     opts.Temperature,
		TopP:            opts.TopP,
		StopSequences:   opts.Stop,
	}

	return GenerateContentRequest{
		Contents:          contents,
		SystemInstruction: systemInstruction,
		GenerationConfig:  genConfig,
		Tools:             tools,
	}
}

func toChatResponse(resp GenerateContentResponse, model string) *llm.ChatResponse {
	var content string

	if len(resp.Candidates) > 0 {
		candidate := resp.Candidates[0]

		if len(candidate.Content.Parts) > 0 {
			content = candidate.Content.Parts[0].Text
		}
	}

	raw, _ := json.Marshal(resp)

	return &llm.ChatResponse{
		ID:      "",
		Model:   model,
		Content: content,
		Raw:     raw,
	}
}
