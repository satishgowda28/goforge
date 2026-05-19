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
	return "Write the content to the file"
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
func (t *FileWrite) InputSchema() map[string]any {
	schema := make(map[string]any)
	schema["path"] = map[string]any{
		"type":        "string",
		"description": "This is the path to the file where the file is located. Do not include file content here",
	}
	schema["content"] = map[string]any{
		"type":        "string",
		"description": "This is the content of the file that need to be written and saved",
	}

	return schema
}

var _ Tool = &FileWrite{}
