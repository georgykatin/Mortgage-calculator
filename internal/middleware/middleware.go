// Package middleware provides a middleware function for logging HTTP response status codes
// and measuring the duration of each request in nanoseconds. It helps in tracking the performance
// of the application by logging the status code and the processing time for each HTTP request.
package middleware

import (
	"log"
	"net/http"
	"time"
)

// responseWriterWrapper is a custom wrapper for the http.ResponseWriter that allows capturing the status code
// of the response. It is used in the RequestInfoMiddleware to track the response status code.
type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code and sends it to the client.
func (rw *responseWriterWrapper) WriteHeader(code int) {
	rw.statusCode = code                // Store the response status code
	rw.ResponseWriter.WriteHeader(code) // Send the response header to the client
}

// RequestInfoMiddleware is a middleware function that logs the status code and duration
// of each HTTP request in nanoseconds. It can be used for performance monitoring and
// debugging the response times of API endpoints.
func RequestInfoMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now() // Capture the start time of the request

		// Wrap the ResponseWriter to capture the status code before the handler processes the request
		wrappedWriter := &responseWriterWrapper{ResponseWriter: w, statusCode: http.StatusOK}

		// Pass the wrapped ResponseWriter to the next handler in the chain
		next.ServeHTTP(wrappedWriter, r)

		// Measure the duration of the request processing
		duration := time.Since(start)

		// Log the status code and request processing time in nanoseconds
		code := wrappedWriter.statusCode
		log.Printf("status_code: %v , duration: %v ns", code, duration.Nanoseconds())
	})
}
