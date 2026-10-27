package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"

	"sumo-mcp-client/cmd/internal/auth"
	"sumo-mcp-client/cmd/internal/chat"
	"sumo-mcp-client/cmd/internal/config"
	"sumo-mcp-client/cmd/internal/sumo"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	token, expiry, err := auth.FetchToken(cfg.TokenURL, cfg.ClientID, cfg.ClientSecret)
	if err != nil {
		log.Fatal(err)
	}

	c, err := sumo.Connect(ctx, cfg.MCPURL, token)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { c.Close() }()

	list, err := c.ListTools(ctx, mcp.ListToolsRequest{})
	if err != nil {
		log.Fatal("list tools: ", err)
	}
	defs := chat.ToolDefs(list.Tools)
	fmt.Printf("Connected. %d Sumo tools available. Type a question, or 'exit'.\n", len(defs))

	var history []any
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("\n> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if line == "exit" {
			return
		}

		// Refresh the token and reconnect if it expires within a minute.
		if time.Until(expiry) < time.Minute {
			token, expiry, err = auth.FetchToken(cfg.TokenURL, cfg.ClientID, cfg.ClientSecret)
			if err != nil {
				fmt.Println("token refresh failed:", err)
				continue
			}
			c.Close()
			c, err = sumo.Connect(ctx, cfg.MCPURL, token)
			if err != nil {
				fmt.Println("reconnect failed:", err)
				return
			}
		}

		history, err = chat.Ask(ctx, c, cfg.AnthropicKey, defs, history, line)
		if err != nil {
			fmt.Println("error:", err)
		}
	}
}
