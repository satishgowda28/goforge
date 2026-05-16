package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
)

type FileWrite struct{}

func (t *FileWrite) Name() string {
	return "file_write"
}

func (t *FileWrite) Description() string {
	return "Description of the tool"
}

func (t *FileWrite) Execute(ctx context.Context, input string) (string, error) {
	type ioType struct {
		Path    string `json:"path"`
		Content string `json:"content"`
	}
	toolName := t.Name()
	var ioArg ioType
	if err := json.Unmarshal([]byte(input), &ioArg); err != nil {
		return "", fmt.Errorf("Error in decoding the input in tool %s , Error : %w", toolName, err)
	}

	err := os.WriteFile(ioArg.Path, []byte(ioArg.Content), 0644)
	if err != nil {
		return "", fmt.Errorf("Error in writing the file tool: %s, path:%s , Error : %w", toolName, ioArg.Path, err)
	}

	return string("Success"), nil
}

var _ Tool = &FileWrite{}
