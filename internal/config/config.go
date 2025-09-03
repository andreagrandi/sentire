package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config holds the application configuration
type Config struct {
	SentryAPIToken string `json:"sentry_api_token"`
}

// LoadConfig loads configuration from environment variables or config file
// Environment variables take precedence over config file values
func LoadConfig() (*Config, error) {
	config := &Config{}

	// First, try to load from environment variable
	if token := os.Getenv("SENTRY_API_TOKEN"); token != "" {
		config.SentryAPIToken = token
		return config, nil
	}

	// If env var is not set, try to load from config file
	configPath, err := getConfigPath()
	if err != nil {
		return nil, fmt.Errorf("failed to determine config path: %w", err)
	}

	if err := loadFromFile(configPath, config); err != nil {
		// If config file doesn't exist or has issues, return error about missing token
		return nil, fmt.Errorf("SENTRY_API_TOKEN environment variable is required (or configure ~/.config/sentire/config.json)")
	}

	// Validate that we have a token
	if config.SentryAPIToken == "" {
		return nil, fmt.Errorf("SENTRY_API_TOKEN is required in config file")
	}

	return config, nil
}

// getConfigPath returns the path to the config file
func getConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	return filepath.Join(homeDir, ".config", "sentire", "config.json"), nil
}

// loadFromFile loads configuration from a JSON file
func loadFromFile(path string, config *Config) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(config); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	return nil
}

// SaveConfig saves the current configuration to the config file
// This is mainly for future use when we add more configuration options
func SaveConfig(config *Config) error {
	configPath, err := getConfigPath()
	if err != nil {
		return fmt.Errorf("failed to determine config path: %w", err)
	}

	// Create directory if it doesn't exist
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	file, err := os.Create(configPath)
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(config); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
