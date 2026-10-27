package claude

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Call sends the conversation and tools to Claude and returns the full response.
func Call(apiKey string, messages []any, tools []any) (map[string]any, error) {
	body, err := json.Marshal(map[string]any{
		"model":      "claude-haiku-5-5",
		"max_tokens": 2048,
		"tools":      tools,
		"messages":   messages,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://api.anthropic.com/v1/messages", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")
	req.Header.Set("content-type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("claude returned %d: %s", resp.StatusCode, raw)
	}

	var out map[string]any
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, err
	}
	return out, nil
}
