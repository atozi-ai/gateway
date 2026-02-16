package openai

import (
	"encoding/json"

	"github.com/atozi-ai/gateway/internal/domain/llm"
)

func toOpenAIRequest(req llm.ChatRequest) openAIChatRequest {
	options := req.Options

	var responseFormat *openAIResponseFormat
	if options.ResponseFormat != nil {
		responseFormat = &openAIResponseFormat{
			Type: options.ResponseFormat.Type,
		}

		if options.ResponseFormat.Type == "json_schema" && len(options.ResponseFormat.Schema) > 0 {
			responseFormat.JSONSchema = &openAIJSONSchema{
				Name:   "response",
				Schema: options.ResponseFormat.Schema,
				Strict: true,
			}
		}
	}

	// Convert tools
	var tools []openAITool
	if len(options.Tools) > 0 {
		tools = make([]openAITool, len(options.Tools))
		for i, tool := range options.Tools {
			tools[i] = openAITool{
				Type: tool.Type,
			}
			if tool.Function != nil {
				tools[i].Function = &openAIFunctionTool{
					Name:        tool.Function.Name,
					Description: tool.Function.Description,
					Parameters:  tool.Function.Parameters,
				}
			}
		}
	}

	// Convert stream options
	var streamOptions *openAIStreamOptions
	if options.StreamOptions != nil {
		streamOptions = &openAIStreamOptions{
			IncludeUsage: options.StreamOptions.IncludeUsage,
		}
	}

	// Handle tool_choice - can be string or object
	var toolChoice interface{}
	if options.ToolChoice != nil {
		// Try to unmarshal as string first, then as ToolChoice object
		if toolChoiceStr, ok := options.ToolChoice.(string); ok {
			toolChoice = toolChoiceStr
		} else {
			// Try to marshal/unmarshal as ToolChoice
			if toolChoiceObj, ok := options.ToolChoice.(llm.ToolChoice); ok {
				toolChoice = toolChoiceObj
			} else {
				// Fallback: marshal to JSON and unmarshal
				jsonBytes, err := json.Marshal(options.ToolChoice)
				if err == nil {
					var tc interface{}
					if err := json.Unmarshal(jsonBytes, &tc); err == nil {
						toolChoice = tc
					}
				}
			}
		}
	}

	return openAIChatRequest{
		Model:             req.Model,
		Messages:          req.Messages,
		FrequencyPenalty:  options.FrequencyPenalty,
		LogitBias:         options.LogitBias,
		Logprobs:          options.Logprobs,
		TopLogprobs:       options.TopLogprobs,
		MaxTokens:         options.MaxTokens,
		N:                 options.N,
		PresencePenalty:   options.PresencePenalty,
		ResponseFormat:    responseFormat,
		Seed:              options.Seed,
		Stop:              options.Stop,
		Stream:            options.Stream,
		StreamOptions:     streamOptions,
		Temperature:       options.Temperature,
		ToolChoice:        toolChoice,
		Tools:             tools,
		TopP:              options.TopP,
		User:              options.User,
		ParallelToolCalls: options.ParallelToolCalls,
		Verbosity:         options.Verbosity,
	}
}
