package chat

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"

	"sumo-mcp-client/cmd/internal/claude"
	"sumo-mcp-client/cmd/internal/sumo"
)

// ToolDefs converts Sumo's MCP tools into Claude's tool format.
func ToolDefs(tools []mcp.Tool) []any {
	var defs []any
	for _, t := range tools {
		if t.Name == "x_amz_bedrock_agentcore_search" {
			continue
		}
		raw, _ := json.Marshal(t.InputSchema)
		var schema any
		_ = json.Unmarshal(raw, &schema)
		defs = append(defs, map[string]any{
			"name":         t.Name,
			"description":  t.Description,
			"input_schema": schema,
		})
	}
	return defs
}

// Ask runs one question through Claude, executing tool calls until Claude answers in text.
func Ask(ctx context.Context, c *client.Client, apiKey string, tools []any, history []any, question string) ([]any, error) {
	history = append(history, map[string]any{"role": "user", "content": question})

	for round := 0; round < 10; round++ {
		resp, err := claude.Call(apiKey, history, tools)
		if err != nil {
			return history, err
		}

		content, _ := resp["content"].([]any)
		history = append(history, map[string]any{"role": "assistant", "content": content})

		if resp["stop_reason"] != "tool_use" {
			for _, b := range content {
				blk, _ := b.(map[string]any)
				if blk["type"] == "text" {
					fmt.Println(blk["text"])
				}
			}
			return history, nil
		}

		var results []any
		for _, b := range content {
			blk, _ := b.(map[string]any)
			if blk["type"] != "tool_use" {
				continue
			}
			name, _ := blk["name"].(string)
			args, _ := blk["input"].(map[string]any)
			fmt.Printf("  [calling %s]\n", name)

			results = append(results, map[string]any{
				"type":        "tool_result",
				"tool_use_id": blk["id"],
				"content":     sumo.RunTool(ctx, c, name, args),
			})
		}
		history = append(history, map[string]any{"role": "user", "content": results})
	}

	return history, fmt.Errorf("stopped after 10 rounds without an answer")
}
