package anthropic

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/satishgowda28/goforge/internal/llm"
	"github.com/satishgowda28/goforge/pkg/config"
)

type Client struct {
	client  anthropic.Client
	model   string
	timeout int
}

func NewClient() *Client {
	// config
	llmCfg := config.Get().LLM
	// client instance
	client := anthropic.NewClient(option.WithAPIKey(llmCfg.APIKey), option.WithRequestTimeout(time.Duration(llmCfg.TimeoutSeconds)*time.Second))
	return &Client{
		client: client,
		model:  llmCfg.Model,
	}
}

func (c *Client) Complete(ctx context.Context, req llm.CompletionRequest) (llm.CompletionResponse, error) {
	var messages []anthropic.MessageParam
	for _, msg := range req.Messages {
		var blocks []anthropic.ContentBlockParamUnion
		for _, block := range msg.Content {
			switch block.Type {
			case "text":
				blocks = append(blocks, anthropic.NewTextBlock(block.Text))
			case "tool_use":
				blocks = append(blocks, anthropic.NewToolUseBlock(block.ID, block.Input, block.Name))
			case "tool_result":
				blocks = append(blocks, anthropic.NewToolResultBlock(block.ToolUseID, block.Content, false))
			}
		}
		switch msg.Role {
		case "user":
			messages = append(messages, anthropic.NewUserMessage(blocks...))
		case "assistant":
			messages = append(messages, anthropic.NewAssistantMessage(blocks...))
		}
	}

	var tools []anthropic.ToolUnionParam
	for _, tool := range req.Tools {
		var required []string
		for key := range tool.InputSchema {
			required = append(required, key)
		}
		t := anthropic.ToolUnionParam{
			OfTool: &anthropic.ToolParam{
				Name:        tool.Name,
				Description: anthropic.String(tool.Description),
				InputSchema: anthropic.ToolInputSchemaParam{
					Properties: tool.InputSchema,
					Required:   required,
				},
			},
		}
		tools = append(tools, t)
	}

	message, err := c.client.Messages.New(ctx, anthropic.MessageNewParams{
		MaxTokens: int64(req.MaxTokens),
		System: []anthropic.TextBlockParam{
			{Text: req.SystemPrompt},
		},
		Messages: messages,
		Tools:    tools,
		Model:    c.model,
	})
	if err != nil {
		return llm.CompletionResponse{}, fmt.Errorf("client error sending message %w", err)
	}
	var block []llm.ContentBlock
	for _, b := range message.Content {
		switch b.Type {
		case "text":
			block = append(block, llm.ContentBlock{
				Type: b.Type,
				Text: b.AsText().Text,
			})
		case "tool_use":
			t := b.AsToolUse()
			var io map[string]any
			if err := json.Unmarshal(t.Input, &io); err != nil {
				return llm.CompletionResponse{}, fmt.Errorf("%w", err)
			}
			block = append(block, llm.ContentBlock{
				Type:  string(t.Type),
				ID:    t.ID,
				Name:  t.Name,
				Input: io,
			})
		}

	}
	return llm.CompletionResponse{
		Content:    block,
		StopReason: string(message.StopReason),
	}, nil
}

func (c *Client) Stream(ctx context.Context, req llm.CompletionRequest) (<-chan llm.StreamChunk, error) {

	return make(<-chan llm.StreamChunk), nil
}

var _ llm.LLMProvider = &Client{}
