package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadConfig_Success(t *testing.T) {
	// Define the path to the ./internal/config directory
	baseDir := "./internal/config"
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		t.Fatalf("failed to create config directory: %v", err)
	}

	// Path to the temporary configuration file
	tempFile := filepath.Join(baseDir, "test_config.yaml")

	// Create content for the configuration file
	content := []byte(`
server:
  port: 8080
`)

	// Create a file with the provided content
	err := os.WriteFile(tempFile, content, 0644)
	if err != nil {
		t.Fatalf("failed to create temp config file: %v", err)
	}

	// Remove the file after the test completes
	defer os.Remove(tempFile)

	// Call LoadConfig with the created file
	cfg, err := LoadConfig("test_config.yaml")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify that the data was loaded correctly
	if cfg.Server.Port != 8080 {
		t.Errorf("expected port 8080, got %d", cfg.Server.Port)
	}
}

func TestLoadConfig_FileNotFound(t *testing.T) {
	// Attempt to load a non-existent file
	_, err := LoadConfig("non_existing_config.yaml")
	if err == nil {
		t.Fatal("expected error but got nil")
	}

	// Check that the error is related to the missing file
	if !strings.Contains(err.Error(), "failed to read file") || !strings.Contains(err.Error(), "non_existing_config.yaml") {
		t.Errorf("expected 'file not found' error, got: %v", err)
	}
}
