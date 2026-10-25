package mcp

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

type MCPClient struct {
	client *client.Client
}

func NewSSEClient(ctx context.Context, endpoint string, headers map[string]string) (*MCPClient, error) {
	c, err := client.NewSSEMCPClient(endpoint, client.WithHeaders(headers))
	if err != nil {
		return nil, fmt.Errorf("create client: %w", err)
	}

	if err := c.Start(ctx); err != nil {
		return nil, fmt.Errorf("start: %w", err)
	}

	_, err = c.Initialize(ctx, mcp.InitializeRequest{
		Params: mcp.InitializeParams{
			ProtocolVersion: mcp.LATEST_PROTOCOL_VERSION,
			ClientInfo: mcp.Implementation{
				Name:    "sumo-mcp-client",
				Version: "0.1.0",
			},
		},
	})
	if err != nil {
		_ = c.Close()
		return nil, fmt.Errorf("initialize: %w", err)
	}

	return &MCPClient{client: c}, nil
}

func (m *MCPClient) ListTools(ctx context.Context) ([]mcp.Tool, error) {
	res, err := m.client.ListTools(ctx, mcp.ListToolsRequest{})
	if err != nil {
		return nil, fmt.Errorf("list tools: %w", err)
	}
	return res.Tools, nil
}

func (m *MCPClient) Close() error {
	return m.client.Close()
}