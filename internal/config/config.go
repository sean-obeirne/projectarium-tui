package config

import (
	"os"
)

// Config holds the application configuration
type Config struct {
	APIBaseURL string
}

// Load loads configuration from environment variables
func Load() *Config {
	apiURL := os.Getenv("PROJECTARIUM_API_URL")
	if apiURL == "" {
		// Default to localhost
		apiURL = "http://localhost:8080/api"
	}

	return &Config{
		APIBaseURL: apiURL,
	}
}
