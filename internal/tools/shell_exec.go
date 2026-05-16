package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"slices"

	"github.com/satishgowda28/goforge/pkg/config"
)

type ShellExec struct{}

func (t *ShellExec) Name() string {
	return "shell_exec"
}

func (t *ShellExec) Description() string {
	return "Description of the tool"
}

func (t *ShellExec) Execute(ctx context.Context, input string) (string, error) {
	type ioType struct {
		Command string   `json:"command"`
		Args    []string `json:"args"`
	}
	toolName := t.Name()

	var ioArg ioType
	if err := json.Unmarshal([]byte(input), &ioArg); err != nil {
		return "", fmt.Errorf("Error in decoding the input in tool %s , Error : %w", toolName, err)
	}
	cmdAllowed := config.Get().Tools.ShellAllowedCommands
	if !(slices.Contains(cmdAllowed, ioArg.Command)) {
		return "", fmt.Errorf("This is shell command is not allowed %s", ioArg.Command)
	}
	cmd := exec.CommandContext(ctx, ioArg.Command, ioArg.Args...)

	var stdOut bytes.Buffer
	var stdError bytes.Buffer
	cmd.Stdout = &stdOut
	cmd.Stderr = &stdError

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("Erro while runnug the cmd: %s, Error: %w", ioArg.Command, err)
	}

	return stdOut.String(), nil
}

var _ Tool = &ShellExec{}
