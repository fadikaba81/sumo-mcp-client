package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
)

type tokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

func fetchToken(tokenURL, clientID, clientSecret string) (string, time.Time, error) {
	form := url.Values{}
	form.Set("grant_type", "client_credentials")
	form.Set("client_id", clientID)
	form.Set("client_secret", clientSecret)

	resp, err := http.PostForm(tokenURL, form)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("token request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", time.Time{}, fmt.Errorf("token endpoint returned %d: %s", resp.StatusCode, body)
	}

	var tr tokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tr); err != nil {
		return "", time.Time{}, fmt.Errorf("decode token: %w", err)
	}
	if tr.AccessToken == "" {
		return "", time.Time{}, fmt.Errorf("no access_token in response")
	}

	expiry := time.Now().Add(time.Duration(tr.ExpiresIn) * time.Second)
	return tr.AccessToken, expiry, nil
}

func requireEnv(name string) string {
	v := os.Getenv(name)
	if v == "" {
		log.Fatalf("%s is not set", name)
	}
	return v
}

func main() {
	tokenURL := requireEnv("SUMO_TOKEN_URL")
	clientID := requireEnv("SUMO_CLIENT_ID")
	clientSecret := requireEnv("SUMO_CLIENT_KEY")
	mcpURL := requireEnv("SUMO_MCP_URL")

	token, expiry, err := fetchToken(tokenURL, clientID, clientSecret)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("token ok, expires %s\n", expiry.Format(time.RFC3339))

	ctx := context.Background()

	c, err := client.NewStreamableHttpClient(mcpURL, transport.WithHTTPHeaders(map[string]string{
		"Authorization": "Bearer " + token,
	}))
	if err != nil {
		log.Fatal("create client: ", err)
	}
	defer c.Close()

	_, err = c.Initialize(ctx, mcp.InitializeRequest{
		Params: mcp.InitializeParams{
			ProtocolVersion: mcp.LATEST_PROTOCOL_VERSION,
			ClientInfo:      mcp.Implementation{Name: "sumo-chat", Version: "0.1.0"},
		},
	})
	if err != nil {
		log.Fatal("initialize: ", err)
	}

	list, err := c.ListTools(ctx, mcp.ListToolsRequest{})
	if err != nil {
		log.Fatal("list tools: ", err)
	}

	fmt.Printf("%d tools:\n", len(list.Tools))
	for _, t := range list.Tools {
		fmt.Printf("  %-40s %s\n", t.Name, strings.TrimSpace(t.Description))
	}
}
