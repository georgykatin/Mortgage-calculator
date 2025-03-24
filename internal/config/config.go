// Package config provides functionality to load and parse the application's configuration
// from a YAML file. It defines a Config structure that maps to the configuration file
// and includes a function to load the configuration and return it as a Config object.
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config represents the application's configuration structure. It contains settings for various
// parts of the application, such as the server configuration.
type Config struct {
	// Server contains configuration settings related to the server, such as the port number.
	Server struct {
		// Port is the port number on which the server will listen for incoming requests.
		Port int `yaml:"port"`
	} `yaml:"server"`
}

// LoadConfig loads the configuration from the specified YAML file. It reads the file, unmarshals
// the content into a Config structure, and returns the Config object or an error if something goes wrong.
func LoadConfig(filename string) (*Config, error) {
	// Define the base directory for configuration files
	basePath := "./internal/config/"

	// Convert filename to an absolute path
	absFilename, err := filepath.Abs(filepath.Clean(filename))
	if err != nil {
		return nil, fmt.Errorf("failed to normalize filename: %w", err)
	}

	// Ensure the file is inside the base directory
	// If the file is already in the base directory, use it directly
	if !strings.HasPrefix(absFilename, basePath) {
		// If the file is outside basePath, join it with basePath
		absFilename = filepath.Join(basePath, filename)
	}

	// Clean the file path to avoid duplicate basePath entries
	absFilename = filepath.Clean(absFilename)

	// Read the configuration file
	data, err := os.ReadFile(absFilename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filename, err)
	}

	// Create an instance of Config to store the parsed data
	var config Config

	// Unmarshal YAML data into the Config structure
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	// Return the populated Config object
	return &config, nil
}
