package llm

import "context"

type CompletionRequest struct {
	SystemPrompt string    // agent's persona/instructions
	Messages     []Message // conversation history
	MaxTokens    int
}

type ContentBlock struct {
	Type string
	Text string
}

type Message struct {
	Role    string
	Content []ContentBlock
}

type CompletionResponse struct {
	Content    []ContentBlock
	StopReason string
}

type StreamChunk struct {
	Delta string // text fragment
	Done  bool   // true on final chunk
	Err   error
}

type LLMProvider interface {
	Complete(ctx context.Context, req CompletionRequest) (CompletionResponse, error)

	Stream(ctx context.Context, req CompletionRequest) (<-chan StreamChunk, error)
}
