package middleware

import (
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

// LoggingMiddleware logs each HTTP request with method, path, status, duration, and trace/span IDs if present.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := &responseWriter{ResponseWriter: w, status: 200}

		// DEBUG: Log at start of middleware
		zerolog.Ctx(r.Context()).Info().Msg("[DEBUG] LoggingMiddleware entered")

		next.ServeHTTP(ww, r)

		duration := time.Since(start)
		logger := zerolog.Ctx(r.Context())
		logger.Info().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Int("status", ww.status).
			Dur("duration", duration).
			Msg("request completed")

		// DEBUG: Log after request log
		zerolog.Ctx(r.Context()).Info().Msg("[DEBUG] LoggingMiddleware after request log")
	})
}

// responseWriter wraps http.ResponseWriter to capture status code
// for logging purposes.
type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}
