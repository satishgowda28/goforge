package agent

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/satishgowda28/goforge/internal/llm"
	"github.com/satishgowda28/goforge/internal/tools"
)

type Runner struct {
	LLM    llm.LLMProvider
	Tools  *tools.Registry
	Memory Memory
	Config AgentConfig
}

type RunResult struct {
	FinalAnswer string
	Steps       []Step
}

func NewRunner(llm llm.LLMProvider, tools *tools.Registry, memory Memory, config AgentConfig) *Runner {
	return &Runner{
		LLM:    llm,
		Tools:  tools,
		Memory: memory,
		Config: config,
	}
}

/* type ToolsDefination struct {
	Name        string
	Description string
	InputSchema map[string]any
} */

func (r *Runner) Run(ctx context.Context, task string) (RunResult, error) {

	maxStep := r.Config.MaxSteps
	if maxStep == 0 {
		maxStep = 10
	}
	messages := llm.Message{
		Role: "user",
		Content: []llm.ContentBlock{
			{Type: "text", Text: task},
		},
	}
	var allTools []llm.ToolsDefination
	var prompt llm.CompletionRequest
	prompt.SystemPrompt = r.Config.SystemPrompt
	prompt.MaxTokens = 4000
	prompt.Messages = append(prompt.Messages, messages)
	for _, toolName := range r.Config.Tools {
		tool, ok := r.Tools.Get(toolName)
		if ok {
			toolDefination := llm.ToolsDefination{
				Name:        tool.Name(),
				Description: tool.Description(),
				InputSchema: tool.InputSchema(),
			}
			allTools = append(allTools, toolDefination)
		}
	}
	prompt.Tools = allTools

	for range maxStep {
		res, err := r.LLM.Complete(ctx, prompt)
		if err != nil {
			return RunResult{}, fmt.Errorf("%w", err)
		}
		// this is specific to anthorpic we need to make global
		if res.StopReason == "tool_use" {
			assistantMsg := llm.Message{
				Role:    "assistant",
				Content: res.Content,
			}
			prompt.Messages = append(prompt.Messages, assistantMsg)
			for _, block := range res.Content {
				switch block.Type {
				case "tool_use":
					tool, ok := r.Tools.Get(block.Name)
					if !ok {
						fmt.Printf("tool not found: %s, skipping\n", block.Name)
						continue
					}
					input, err := json.Marshal(block.Input)
					if err != nil {
						return RunResult{}, fmt.Errorf("error in Serializing input tool: %s, error:%w", block.Name, err)
					}
					toolResult, err := tool.Execute(ctx, string(input))
					if err != nil {
						return RunResult{}, fmt.Errorf("error in executing tool: %s, error:%w", block.Name, err)
					}
					toolsResponse := llm.Message{
						Role: "user",
						Content: []llm.ContentBlock{
							{
								Type:      "tool_result",
								ToolUseID: block.ID,
								Content:   toolResult,
							},
						},
					}
					prompt.Messages = append(prompt.Messages, toolsResponse)
					r.Memory.Add(Step{
						Type:      "action",
						ToolName:  block.Name,
						ToolUseID: block.ID,
						Input:     string(input),
					})
					r.Memory.Add(Step{
						Type:      "observation",
						ToolUseID: block.ID,
						Output:    toolResult,
					})
				case "text":
					r.Memory.Add(Step{
						Type:    "thought",
						Thought: block.Text,
					})
				}
			}
		}
		if res.StopReason == "end_turn" {
			for _, block := range res.Content {
				if block.Type == "text" {
					return RunResult{FinalAnswer: block.Text, Steps: r.Memory.Recent(maxStep)}, nil
				}
			}
		}

	}
	return RunResult{}, fmt.Errorf("Max step exceded")
}
