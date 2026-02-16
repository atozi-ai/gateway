package llm

import "context"

type StreamChunk struct {
	ID      string
	Object  string
	Created int64
	Model   string
	Choices []StreamChoice
	Usage   *Usage
	Raw     []byte // Raw JSON chunk
}

type StreamChoice struct {
	Index        int
	Delta        StreamDelta
	FinishReason *string
	Logprobs     interface{}
}

type StreamDelta struct {
	Role      *string
	Content   *string
	ToolCalls []interface{}
}

type Usage struct {
	PromptTokens     int
	CompletionTokens int
	TotalTokens      int
}

type Provider interface {
	Name() string

	Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error)
	
	// ChatStream streams chat responses. The callback function is called for each chunk.
	// If the callback returns an error, streaming is stopped and that error is returned.
	ChatStream(ctx context.Context, req ChatRequest, callback func(*StreamChunk) error) error
}
