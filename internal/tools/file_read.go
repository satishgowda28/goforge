package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
)

type FileRead struct{}

func (t *FileRead) Name() string {
	return "file_read"
}

func (t *FileRead) Description() string {
	return "Description of the tool"
}

func (t *FileRead) Execute(ctx context.Context, input string) (string, error) {
	type ioType struct {
		Path string `json:"path"`
	}
	toolName := t.Name()
	var ioArg ioType
	if err := json.Unmarshal([]byte(input), &ioArg); err != nil {
		return "", fmt.Errorf("Error in decoding the input in tool %s , Error : %w", toolName, err)
	}

	data, err := os.ReadFile(ioArg.Path)
	if err != nil {
		return "", fmt.Errorf("Error in reading the file tool: %s, path:%s , Error : %w", toolName, ioArg.Path, err)
	}

	return string(data), nil
}

var _ Tool = &FileRead{}
