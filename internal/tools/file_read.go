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
	return "This is used to read the content of the file"
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

func (t *FileRead) InputSchema() map[string]any {
	schema := make(map[string]any)
	schema["path"] = map[string]any{
		"type":        "string",
		"description": "This is the path to the file where the file is located.",
	}

	return schema
}

var _ Tool = &FileRead{}
