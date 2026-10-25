package config

import (
	"os"
)

type Config struct {
	ANTHROPIC_API_KEY string
	SUMO_ACCESS_ID             string
	SUMO_ACCESS_KEY            string
	SUMO_MCP_URL               string
}

func LoadConfig() Config {
	cfg := Config{
		ANTHROPIC_API_KEY: os.Getenv("ANTHROPIC_API_KEY"),
		SUMO_ACCESS_ID:             os.Getenv("SUMO_ACCESS_ID"),
		SUMO_ACCESS_KEY:            os.Getenv("SUMO_ACCESS_KEY"),
		SUMO_MCP_URL:               os.Getenv("SUMO_MCP_URL"),
	}

	return cfg
}
