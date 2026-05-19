package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type HttpFetch struct{}

func (t *HttpFetch) Name() string {
	return "http_fetch"
}

func (t *HttpFetch) Description() string {
	return "Make a get request to get an appropriate data"
}

func (t *HttpFetch) Execute(ctx context.Context, input string) (string, error) {
	type ioType struct {
		URL string `json:"url"`
	}
	toolName := t.Name()
	var ioArg ioType
	if err := json.Unmarshal([]byte(input), &ioArg); err != nil {
		return "", fmt.Errorf("Error in decoding the input in tool %s , Error : %w", toolName, err)
	}

	client := &http.Client{}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ioArg.URL, nil)
	if err != nil {
		return "", fmt.Errorf("Err in creating request: %w", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error in response tool:%s, url: %s, error: %w", toolName, ioArg.URL, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed response tool:%s, url: %s", toolName, ioArg.URL)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error in reading resp tool:%s, url: %s, error: %w", toolName, ioArg.URL, err)
	}

	return string(body), nil
}
func (t *HttpFetch) InputSchema() map[string]any {
	schema := make(map[string]any)
	schema["url"] = map[string]any{
		"type":        "string",
		"description": "This is the URL used to make get request",
	}

	return schema
}

var _ Tool = &HttpFetch{}
