package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	TokenURL     string
	ClientID     string
	ClientSecret string
	MCPURL       string
	AnthropicKey string
}

// Load reads the environment and reports every missing variable at once.
func Load() (Config, error) {
	cfg := Config{
		TokenURL:     os.Getenv("SUMO_TOKEN_URL"),
		ClientID:     os.Getenv("SUMO_CLIENT_ID"),
		ClientSecret: os.Getenv("SUMO_CLIENT_KEY"),
		MCPURL:       os.Getenv("SUMO_MCP_URL"),
		AnthropicKey: os.Getenv("ANTHROPIC_API_KEY"),
	}

	var missing []string
	if cfg.TokenURL == "" {
		missing = append(missing, "SUMO_TOKEN_URL")
	}
	if cfg.ClientID == "" {
		missing = append(missing, "SUMO_CLIENT_ID")
	}
	if cfg.ClientSecret == "" {
		missing = append(missing, "SUMO_CLIENT_KEY")
	}
	if cfg.MCPURL == "" {
		missing = append(missing, "SUMO_MCP_URL")
	}
	if cfg.AnthropicKey == "" {
		missing = append(missing, "ANTHROPIC_API_KEY")
	}
	if len(missing) > 0 {
		return Config{}, fmt.Errorf("missing environment variables: %s", strings.Join(missing, ", "))
	}
	return cfg, nil
}
