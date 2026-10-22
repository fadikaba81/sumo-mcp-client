package mcp

import ( 
	"context"
	"fmt"
	"github.com/mark3labs/mcp-go/client"
)

type MCPClient struct {
	client *client.Client
}

func NewSSEClient(endpoint string) (*MCPClient, error) {
	c, err := client.NewClient(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	return &MCPClient{
		client: c
	}, nil
}

func (m *MCPClient) ListTools(cxt context.Context) {

	tools, err := m.client.ListTools(cxt)

	if err != nil {
		fmt.Printf("failed to list tools: %v\n", err)
		return
	}
	fmt.Printf("tools: %v\n", tools)

	for _, tool := range tools.Tools{
		fmt.Printf("tool: %v\n", tool.name)
	}

}