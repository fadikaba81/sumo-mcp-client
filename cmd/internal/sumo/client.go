package sumo

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
)

// Connect creates an MCP client with the bearer token and initializes it.
func Connect(ctx context.Context, mcpURL, token string) (*client.Client, error) {
	c, err := client.NewStreamableHttpClient(mcpURL, transport.WithHTTPHeaders(map[string]string{
		"Authorization": "Bearer " + token,
	}))
	if err != nil {
		return nil, fmt.Errorf("create client: %w", err)
	}

	_, err = c.Initialize(ctx, mcp.InitializeRequest{
		Params: mcp.InitializeParams{
			ProtocolVersion: mcp.LATEST_PROTOCOL_VERSION,
			ClientInfo:      mcp.Implementation{Name: "sumo-chat", Version: "0.1.0"},
		},
	})
	if err != nil {
		c.Close()
		return nil, fmt.Errorf("initialize: %w", err)
	}
	return c, nil
}

// RunTool calls a Sumo tool and returns its content as text, truncated if long.
func RunTool(ctx context.Context, c *client.Client, name string, args map[string]any) string {
	callCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	res, err := c.CallTool(callCtx, mcp.CallToolRequest{
		Params: mcp.CallToolParams{Name: name, Arguments: args},
	})
	if err != nil {
		return "error: " + err.Error()
	}

	out, _ := json.Marshal(res.Content)
	text := string(out)
	if len(text) > 20000 {
		text = text[:20000] + "...[truncated]"
	}
	return text
}
