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
	// Try config file first
	configFile := os.Getenv("HOME") + "/.config/pj-tui.env"
	if _, err := os.Stat(configFile); err == nil {
		// Read config file
		if data, err := os.ReadFile(configFile); err == nil {
			// Simple parsing of KEY=VALUE format
			for _, line := range splitLines(string(data)) {
				if len(line) > 0 && line[0] != '#' {
					if idx := findChar(line, '='); idx != -1 {
						key := line[:idx]
						value := line[idx+1:]
						if key == "PROJECTARIUM_API_URL" {
							os.Setenv(key, value)
						}
					}
				}
			}
		}
	}

	apiURL := os.Getenv("PROJECTARIUM_API_URL")
	if apiURL == "" {
		// Default to localhost
		apiURL = "http://localhost:8888/api"
	}

	return &Config{
		APIBaseURL: apiURL,
	}
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}

func findChar(s string, c byte) int {
	for i := 0; i < len(s); i++ {
		if s[i] == c {
			return i
		}
	}
	return -1
}
