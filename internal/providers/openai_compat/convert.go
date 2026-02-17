package openaicompat

import (
	"encoding/json"

	"github.com/atozi-ai/gateway/internal/domain/llm"
)

// toRequest converts the domain ChatRequest into the OpenAI-compatible wire format.
func toRequest(req llm.ChatRequest) chatRequest {
	opts := req.Options

	var rf *responseFormat
	if opts.ResponseFormat != nil {
		rf = &responseFormat{Type: opts.ResponseFormat.Type}
		if opts.ResponseFormat.Type == "json_schema" && len(opts.ResponseFormat.Schema) > 0 {
			rf.JSONSchema = &jsonSchema{
				Name:   "response",
				Schema: opts.ResponseFormat.Schema,
				Strict: true,
			}
		}
	}

	var tools []tool
	if len(opts.Tools) > 0 {
		tools = make([]tool, len(opts.Tools))
		for i, t := range opts.Tools {
			tools[i] = tool{Type: t.Type}
			if t.Function != nil {
				tools[i].Function = &functionTool{
					Name:        t.Function.Name,
					Description: t.Function.Description,
					Parameters:  t.Function.Parameters,
				}
			}
		}
	}

	var so *streamOptions
	if opts.StreamOptions != nil {
		so = &streamOptions{IncludeUsage: opts.StreamOptions.IncludeUsage}
	}

	tc := resolveToolChoice(opts.ToolChoice)

	return chatRequest{
		Model:             req.Model,
		Messages:          req.Messages,
		FrequencyPenalty:  opts.FrequencyPenalty,
		LogitBias:         opts.LogitBias,
		Logprobs:          opts.Logprobs,
		TopLogprobs:       opts.TopLogprobs,
		MaxTokens:         opts.MaxTokens,
		N:                 opts.N,
		PresencePenalty:   opts.PresencePenalty,
		ResponseFormat:    rf,
		Seed:              opts.Seed,
		Stop:              opts.Stop,
		Stream:            opts.Stream,
		StreamOptions:     so,
		Temperature:       opts.Temperature,
		ToolChoice:        tc,
		Tools:             tools,
		TopP:              opts.TopP,
		User:              opts.User,
		ParallelToolCalls: opts.ParallelToolCalls,
		Verbosity:         opts.Verbosity,
	}
}

// resolveToolChoice normalises the polymorphic ToolChoice field.
func resolveToolChoice(raw interface{}) interface{} {
	if raw == nil {
		return nil
	}
	if s, ok := raw.(string); ok {
		return s
	}
	if tc, ok := raw.(llm.ToolChoice); ok {
		return tc
	}
	// Fallback: round-trip through JSON.
	b, err := json.Marshal(raw)
	if err != nil {
		return nil
	}
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return nil
	}
	return v
}
