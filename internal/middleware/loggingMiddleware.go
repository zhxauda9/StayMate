package middleware

import (
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

// LoggingMiddleware logs HTTP requests and responses.
func LoggingMiddleware(next http.Handler, logger *zerolog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Log incoming request details
		logger.Info().Str("method", r.Method).Str("Endpoint", r.URL.Path).Str("IP", r.RemoteAddr).Msg("Request received")

		// Capture the response status code by wrapping the ResponseWriter
		responseWriter := &LoggingResponseWriter{ResponseWriter: w}
		next.ServeHTTP(responseWriter, r)

		// Log request processing time and response status
		if responseWriter.status == 0 {
			responseWriter.status = 200
		}
		logger.Info().Str("method", r.Method).Str("path", r.URL.Path).Int("status", responseWriter.status).Any("duration", time.Since(start)).Msg("Request processed")
	})
}

// Custom ResponseWriter to check status code. This done for logging
type LoggingResponseWriter struct {
	http.ResponseWriter
	status int
}

// Overriding WriteHeader function.
func (w *LoggingResponseWriter) WriteHeader(statusCode int) {
	w.status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}
