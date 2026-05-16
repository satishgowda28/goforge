package anthropic

import (
	"context"
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
			blocks = append(blocks, anthropic.NewTextBlock(block.Text))
		}
		switch msg.Role {
		case "user":
			messages = append(messages, anthropic.NewUserMessage(blocks...))
		case "assistant":
			messages = append(messages, anthropic.NewAssistantMessage(blocks...))
		}
	}
	message, err := c.client.Messages.New(ctx, anthropic.MessageNewParams{
		MaxTokens: int64(req.MaxTokens),
		System: []anthropic.TextBlockParam{
			{Text: req.SystemPrompt},
		},
		Messages: messages,
		Model:    c.model,
	})
	if err != nil {
		return llm.CompletionResponse{}, fmt.Errorf("client error sending message %w", err)
	}
	var block []llm.ContentBlock
	for _, b := range message.Content {
		block = append(block, llm.ContentBlock{
			Type: b.Type,
			Text: b.AsText().Text,
		})
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
