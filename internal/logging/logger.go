package logging

import (
	"io"
	"os"
	"sync"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// LoggerConfig holds configuration for the logger.
type LoggerConfig struct {
	Level      zerolog.Level // Log level
	Format     string        // "json" or "console"
	Writers    []io.Writer   // Output writers (stdout, file, etc.)
	TimeFormat string        // Optional time format for console
}

var (
	logger zerolog.Logger
	once   sync.Once
)

// InitLogger initializes the global logger with the given config.
func InitLogger(cfg LoggerConfig) {
	once.Do(func() {
		var w io.Writer
		if len(cfg.Writers) == 0 {
			w = os.Stdout
		} else if len(cfg.Writers) == 1 {
			w = cfg.Writers[0]
		} else {
			w = io.MultiWriter(cfg.Writers...)
		}

		if cfg.Format == "console" {
			cw := zerolog.ConsoleWriter{Out: w, TimeFormat: cfg.TimeFormat}
			logger = zerolog.New(cw).With().Timestamp().Logger()
		} else {
			logger = zerolog.New(w).With().Timestamp().Logger()
		}

		zerolog.SetGlobalLevel(cfg.Level)
		log.Logger = logger
	})
}

// L returns the global logger instance.
func L() zerolog.Logger {
	return logger
}

// WithTrace returns a logger with traceid/spanid fields from context.
func WithTrace(ctx io.Writer, traceID, spanID string) zerolog.Logger {
	return logger.With().Str("traceid", traceID).Str("spanid", spanID).Logger()
}
