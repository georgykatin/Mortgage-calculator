// Package server sets up and runs the HTTP server for the mortgage calculation service.
//
// It creates a new HTTP server, initializes the necessary handlers, and manages graceful shutdown.
// The server handles incoming requests for mortgage calculation and cache management.
// It also listens for system interrupts to initiate a clean shutdown of the server.
//
// Functions:
//   - New: Initializes the server with provided handlers and configuration, starts it, and manages graceful shutdown.
//   - initHandlers: Sets up the HTTP request handlers and applies middleware.
package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sber/internal/config"
	"sber/internal/handlers"
	"sber/internal/middleware"
	"syscall"
	"time"
)

// New initializes the HTTP server with the provided handlers and configuration.
// It listens for system interrupts to shut down gracefully.
func New(h *handlers.Handlers, cfg *config.Config) {
	// Create a new HTTP server with the specified configuration and timeouts
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.Server.Port), // Set the port for the server
		Handler:           initHandlers(h),                     // Initialize handlers
		ReadHeaderTimeout: 5 * time.Second,                     // Timeout for reading headers
		WriteTimeout:      10 * time.Second,                    // Timeout for writing the response
		ReadTimeout:       10 * time.Second,                    // Timeout for reading the request body
	}

	// Channel to receive signals for stopping the server (e.g., SIGTERM or SIGINT)
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	// Start the server in a goroutine for asynchronous request handling
	go func() {
		log.Printf("Server started on port %d", cfg.Server.Port)
		// Start the server, log and terminate if an error occurs (except for server closure)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Failed to start server with error: %v", err)
		}
	}()

	// Wait for a signal to stop the server
	<-stopChan
	log.Println("Shutting down server")

	// Create a context with a timeout for shutting down the server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt to gracefully shut down the server
	if err := srv.Shutdown(ctx); err != nil {
		// Log the error and terminate the program
		log.Println("Server shutdown error:", err)
		return // Return to indicate failure
	}
	// Log successful server stop
	log.Println("Server stopped")
}

// initHandlers initializes the HTTP handlers for the service and applies middleware.
func initHandlers(h *handlers.Handlers) http.Handler {
	// Create a new router to handle incoming requests
	r := http.NewServeMux()

	// Register handlers for specific routes
	r.HandleFunc("/execute", h.Execute) // Handler for the /execute route
	r.HandleFunc("/cache", h.Cache)     // Handler for the /cache route

	// Apply middleware to log request information
	return middleware.RequestInfoMiddleware(r)
}
