// Package app contains the main application logic for initializing the application,
// loading configuration settings, setting up dependencies, and starting the server.
// This is the entry point of the application, where all the required components
// are initialized and the server is started.
package app

import (
	"log"
	"sber/internal/cache"
	"sber/internal/config"
	"sber/internal/handlers"
	"sber/internal/server"
)

// Run is the main function for running the application. It loads the configuration from the specified YML file,
// initializes the storage system, creates handler instances, and starts the server with the configured handlers.
func Run() {
	// Load the application configuration from the YML file
	cfg, err := config.LoadConfig("config.yml")
	if err != nil {
		log.Fatalf("Failed to load config from yml file with err: %v", err)
	}

	// Initialize the cache storage system
	storage := cache.New()

	// Create the handlers using the initialized storage
	h := handlers.NewHandlers(storage)

	// Start the server with the configured handlers and loaded configuration
	server.New(h, cfg)
}
