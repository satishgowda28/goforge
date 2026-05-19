package llm

import "context"

type CompletionRequest struct {
	SystemPrompt string    // agent's persona/instructions
	Messages     []Message // conversation history
	Tools        []ToolsDefination
	MaxTokens    int
}

type ToolsDefination struct {
	Name        string
	Description string
	InputSchema map[string]any
}

type ContentBlock struct {
	Type      string
	Text      string         // for "text" type
	ID        string         // for "tool_use" type
	Name      string         // for "tool_use" type
	Input     map[string]any // for "tool_use" type (parsed JSON)
	ToolUseID string         // for "tool_result" type
	Content   string         // for "tool_result" type
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
