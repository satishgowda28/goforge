package agent

import (
	"context"

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

func (r *Runner) Run(ctx context.Context, task string) (RunResult, error) {
	return RunResult{}, nil
}
