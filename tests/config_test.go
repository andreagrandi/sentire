package tests

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sentire/internal/config"
	"testing"
)

func TestLoadConfigFromEnvironmentVariable(t *testing.T) {
	// Clean up any existing env var
	os.Unsetenv("SENTRY_API_TOKEN")

	// Set environment variable
	expectedToken := "test-env-token"
	os.Setenv("SENTRY_API_TOKEN", expectedToken)
	defer os.Unsetenv("SENTRY_API_TOKEN")

	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if cfg.SentryAPIToken != expectedToken {
		t.Errorf("Expected token %s, got %s", expectedToken, cfg.SentryAPIToken)
	}
}

func TestLoadConfigFromFile(t *testing.T) {
	// Clean up any existing env var
	os.Unsetenv("SENTRY_API_TOKEN")

	// Create a temporary config file
	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, ".config", "sentire")
	err := os.MkdirAll(configDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create config directory: %v", err)
	}

	configPath := filepath.Join(configDir, "config.json")
	expectedToken := "test-file-token"

	configData := map[string]string{
		"sentry_api_token": expectedToken,
	}

	file, err := os.Create(configPath)
	if err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(configData); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// Temporarily change home directory to our temp dir
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if cfg.SentryAPIToken != expectedToken {
		t.Errorf("Expected token %s, got %s", expectedToken, cfg.SentryAPIToken)
	}
}

func TestEnvironmentVariableTakesPrecedence(t *testing.T) {
	// Create a temporary config file
	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, ".config", "sentire")
	err := os.MkdirAll(configDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create config directory: %v", err)
	}

	configPath := filepath.Join(configDir, "config.json")
	fileToken := "test-file-token"

	configData := map[string]string{
		"sentry_api_token": fileToken,
	}

	file, err := os.Create(configPath)
	if err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(configData); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// Temporarily change home directory to our temp dir
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	// Set environment variable (should take precedence)
	envToken := "test-env-token"
	os.Setenv("SENTRY_API_TOKEN", envToken)
	defer os.Unsetenv("SENTRY_API_TOKEN")

	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Environment variable should take precedence
	if cfg.SentryAPIToken != envToken {
		t.Errorf("Expected env token %s to take precedence, got %s", envToken, cfg.SentryAPIToken)
	}
}

func TestLoadConfigNoTokenAvailable(t *testing.T) {
	// Clean up any existing env var
	os.Unsetenv("SENTRY_API_TOKEN")

	// Create a temporary directory with no config file
	tempDir := t.TempDir()

	// Temporarily change home directory to our temp dir (no config file exists)
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	_, err := config.LoadConfig()
	if err == nil {
		t.Error("Expected error when no token is available")
	}

	expectedErrorMsg := "SENTRY_API_TOKEN environment variable is required"
	if !containsString(err.Error(), expectedErrorMsg) {
		t.Errorf("Expected error message to contain '%s', got '%s'", expectedErrorMsg, err.Error())
	}
}

func TestLoadConfigMalformedJSON(t *testing.T) {
	// Clean up any existing env var
	os.Unsetenv("SENTRY_API_TOKEN")

	// Create a temporary config file with malformed JSON
	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, ".config", "sentire")
	err := os.MkdirAll(configDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create config directory: %v", err)
	}

	configPath := filepath.Join(configDir, "config.json")

	// Write malformed JSON
	file, err := os.Create(configPath)
	if err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}
	defer file.Close()

	_, err = file.WriteString(`{"sentry_api_token": "test-token"`) // Missing closing brace
	if err != nil {
		t.Fatalf("Failed to write malformed JSON: %v", err)
	}

	// Temporarily change home directory to our temp dir
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	_, err = config.LoadConfig()
	if err == nil {
		t.Error("Expected error when config file has malformed JSON")
	}
}

func TestLoadConfigEmptyTokenInFile(t *testing.T) {
	// Clean up any existing env var
	os.Unsetenv("SENTRY_API_TOKEN")

	// Create a temporary config file with empty token
	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, ".config", "sentire")
	err := os.MkdirAll(configDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create config directory: %v", err)
	}

	configPath := filepath.Join(configDir, "config.json")

	configData := map[string]string{
		"sentry_api_token": "", // Empty token
	}

	file, err := os.Create(configPath)
	if err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(configData); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// Temporarily change home directory to our temp dir
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	_, err = config.LoadConfig()
	if err == nil {
		t.Error("Expected error when config file has empty token")
	}

	expectedErrorMsg := "SENTRY_API_TOKEN is required in config file"
	if !containsString(err.Error(), expectedErrorMsg) {
		t.Errorf("Expected error message to contain '%s', got '%s'", expectedErrorMsg, err.Error())
	}
}

func TestSaveConfig(t *testing.T) {
	// Create a temporary directory
	tempDir := t.TempDir()

	// Temporarily change home directory to our temp dir
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	// Create config to save
	cfg := &config.Config{
		SentryAPIToken: "test-save-token",
	}

	// Save config
	err := config.SaveConfig(cfg)
	if err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Verify config file was created
	configPath := filepath.Join(tempDir, ".config", "sentire", "config.json")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("Config file was not created")
	}

	// Load and verify the saved config
	savedCfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load saved config: %v", err)
	}

	if savedCfg.SentryAPIToken != cfg.SentryAPIToken {
		t.Errorf("Expected saved token %s, got %s", cfg.SentryAPIToken, savedCfg.SentryAPIToken)
	}
}

// Helper function to check if a string contains a substring
func containsString(str, substr string) bool {
	return len(str) >= len(substr) &&
		(str == substr ||
			len(str) > len(substr) &&
				(str[:len(substr)] == substr ||
					str[len(str)-len(substr):] == substr ||
					findSubstring(str, substr)))
}

func findSubstring(str, substr string) bool {
	for i := 0; i <= len(str)-len(substr); i++ {
		if str[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
