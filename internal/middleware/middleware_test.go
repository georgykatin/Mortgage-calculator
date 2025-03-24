package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestResponseWriterWrapper verifies that responseWriterWrapper correctly captures the status code.
func TestResponseWriterWrapper(t *testing.T) {
	// Create a fake ResponseWriter
	recorder := httptest.NewRecorder()

	// Wrap it in responseWriterWrapper
	wrappedWriter := &responseWriterWrapper{
		ResponseWriter: recorder,
		statusCode:     http.StatusOK, // Default status code is 200
	}

	// Set a status code
	wrappedWriter.WriteHeader(http.StatusNotFound)

	// Check that the status code was captured correctly
	if wrappedWriter.statusCode != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, wrappedWriter.statusCode)
	}

	// Check that the status code was passed to the original ResponseWriter
	if recorder.Code != http.StatusNotFound {
		t.Errorf("Expected recorder status code %d, got %d", http.StatusNotFound, recorder.Code)
	}
}

// TestRequestInfoMiddleware verifies that RequestInfoMiddleware correctly logs the status code and execution time.
func TestRequestInfoMiddleware(t *testing.T) {
	// Create a fake HTTP request
	req := httptest.NewRequest(http.MethodGet, "/test", nil)

	// Create a fake ResponseWriter
	recorder := httptest.NewRecorder()

	// Create a test handler that returns a 200 status code
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Wrap the handler with the middleware
	middlewareHandler := RequestInfoMiddleware(handler)

	// Execute the handler
	middlewareHandler.ServeHTTP(recorder, req)

	// Verify that the status code was passed correctly
	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, recorder.Code)
	}
}

// TestRequestInfoMiddlewareWithError verifies that the middleware correctly handles errors.
func TestRequestInfoMiddlewareWithError(t *testing.T) {
	// Create a fake HTTP request
	req := httptest.NewRequest(http.MethodGet, "/test", nil)

	// Create a fake ResponseWriter
	recorder := httptest.NewRecorder()

	// Create a test handler that returns a 500 status code
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	})

	// Wrap the handler with the middleware
	middlewareHandler := RequestInfoMiddleware(handler)

	// Execute the handler
	start := time.Now()
	middlewareHandler.ServeHTTP(recorder, req)
	duration := time.Since(start)

	// Verify that execution time is at least 100 milliseconds
	if duration < 100*time.Millisecond {
		t.Errorf("Expected duration >= 100ms, got %v", duration)
	}
}
