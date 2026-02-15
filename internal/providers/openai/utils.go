package openai

import "github.com/atozi-ai/gateway/internal/domain/llm"

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

	return openAIChatRequest{
		Model:          req.Model,
		Messages:       req.Messages,
		Temperature:    options.Temperature,
		TopP:           options.TopP,
		MaxTokens:      options.MaxTokens,
		Stop:           options.Stop,
		ResponseFormat: responseFormat,
		Verbosity:      options.Verbosity,
	}
}
