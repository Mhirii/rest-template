package logging

import (
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

var logger zerolog.Logger

func InitLogger(l zerolog.Logger) {
	logger = l
}

func L() zerolog.Logger {
	return logger
}

// RequestLoggingHandler attaches the logger to the context and logs the request.
func RequestLoggingHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := &responseWriter{ResponseWriter: w, status: 200}
		ctx := L().WithContext(r.Context())
		next.ServeHTTP(ww, r.WithContext(ctx))
		duration := time.Since(start)
		zerolog.Ctx(ctx).Info().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Int("status", ww.status).
			Dur("duration", duration).
			Str("agent", r.UserAgent()).
			Msg("request completed")
	})
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}
