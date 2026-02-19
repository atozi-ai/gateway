package anthropic_compat

import (
	"github.com/atozi-ai/gateway/internal/domain/llm"
)

func toRequest(req llm.ChatRequest) messageRequest {
	opts := req.Options

	maxTokens := defaultMaxTokens
	if opts.MaxTokens != nil && *opts.MaxTokens > 0 {
		maxTokens = *opts.MaxTokens
	}

	var systemBlocks []contentBlock
	var userMessages []message

	for _, msg := range req.Messages {
		switch msg.Role {
		case llm.RoleSystem:
			systemBlocks = append(systemBlocks, contentBlock{
				Type: contentTypeText,
				Text: msg.Content,
			})
		case llm.RoleUser:
			userMessages = append(userMessages, message{
				Role: roleUser,
				Content: []contentBlock{
					{
						Type: contentTypeText,
						Text: msg.Content,
					},
				},
			})
		case llm.RoleAssistant:
			userMessages = append(userMessages, message{
				Role: roleAssistant,
				Content: []contentBlock{
					{
						Type: contentTypeText,
						Text: msg.Content,
					},
				},
			})
		case llm.RoleTool:
			userMessages = append(userMessages, message{
				Role: roleUser,
				Content: []contentBlock{
					{
						Type:      contentTypeToolResult,
						Content:   msg.Content,
						ToolUseID: "",
					},
				},
			})
		}
	}

	var tools []tool
	if len(opts.Tools) > 0 {
		tools = make([]tool, len(opts.Tools))
		for i, t := range opts.Tools {
			tools[i] = tool{
				Name:        t.Function.Name,
				Description: t.Function.Description,
				Parameters:  t.Function.Parameters,
			}
		}
	}

	var toolChoice *toolChoice
	if opts.ToolChoice != nil {
		tc := resolveToolChoice(opts.ToolChoice)
		if tc != nil {
			toolChoice = tc
		}
	}

	var outputCfg *outputConfig
	if opts.ResponseFormat != nil {
		if opts.ResponseFormat.Type == "json_object" || opts.ResponseFormat.Type == "json_schema" {
			outputCfg = &outputConfig{
				Format: &outputFormat{
					Type:   "json_schema",
					Schema: opts.ResponseFormat.Schema,
				},
			}
		}
	}

	return messageRequest{
		Model:         req.Model,
		MaxTokens:     maxTokens,
		Messages:      userMessages,
		System:        systemBlocks,
		Temperature:   opts.Temperature,
		TopP:          opts.TopP,
		Tools:         tools,
		ToolChoice:    toolChoice,
		Stream:        opts.Stream,
		StopSequences: opts.Stop,
		OutputConfig:  outputCfg,
	}
}

func resolveToolChoice(raw interface{}) *toolChoice {
	if raw == nil {
		return nil
	}
	if s, ok := raw.(string); ok {
		if s == "none" || s == "auto" {
			return &toolChoice{Type: s}
		}
	}
	if tc, ok := raw.(llm.ToolChoice); ok {
		if tc.Function != nil {
			return &toolChoice{
				Type: "tool",
				Name: tc.Function.Name,
			}
		}
	}
	return nil
}
