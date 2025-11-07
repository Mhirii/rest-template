package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/danielgtaylor/huma/v2/humacli"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/mhirii/rest-template/internal/config"
	"github.com/mhirii/rest-template/internal/handlers"
	"github.com/mhirii/rest-template/internal/logging"
	"github.com/mhirii/rest-template/internal/observability"
	"github.com/mhirii/rest-template/internal/service"
)

// Options for the CLI. Pass `--port` or set the `SERVICE_PORT` env var.
type Options struct {
	Port int `help:"Port to listen on" short:"p" default:"8888"`
}

// GreetingOutput represents the greeting operation response.
type GreetingOutput struct {
	Body struct {
		Message string `json:"message" example:"Hello, world!" doc:"Greeting message"`
	}
}

func main() {
	// CLI flags (must be defined before config.Load)
	flag.Int("port", 8888, "Port to listen on")
	flag.String("config", "", "Path to config file")
	flag.Parse()

	// Logger configuration
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
	if err != nil {
		panic(fmt.Sprintf("failed to open log file: %v", err))
	}
	multi := io.MultiWriter(os.Stdout, logFile)
	logger := zerolog.New(multi).With().Timestamp().Logger()
	logging.InitLogger(logger)

	config.Load()

	// --- OpenTelemetry Tracing Setup ---
	ctx := context.Background()
	shutdown, err := observability.SetupOTelTracing(ctx, "rest-template")
	if err != nil {
		panic(fmt.Sprintf("failed to set up OpenTelemetry: %v", err))
	}
	defer shutdown(context.Background())
	// -----------------------------------

	l := logging.L()
	l.Info().Msg("Starting API server")
	zerolog.Ctx(context.Background()).Info().Msg("Test log from zerolog.Ctx(context.Background())")
	cli := humacli.New(func(hooks humacli.Hooks, options *Options) {
		// Use config.Get() for config values
		cfg := config.Get()
		_ = cfg
		router := chi.NewMux()

		api := humachi.New(router, huma.DefaultConfig("Huma API Template", "1.0.0"))

		// Register example REST endpoints
		store := service.NewExampleStore()
		handlers.RegisterExampleRoutes(api, store)

		huma.Get(api, "/greeting/{name}", func(ctx context.Context, input *struct {
			Name string `path:"name" maxLength:"30" example:"world" doc:"Name to greet"`
		},
		) (*GreetingOutput, error) {
			resp := &GreetingOutput{}
			resp.Body.Message = fmt.Sprintf("Hello, %s!", input.Name)
			return resp, nil
		})

		hooks.OnStart(func() {
			l.Info().Int("port", cfg.Port).Msg("API server listening")
			// Wrap router with otelhttp for tracing
			importOtelHttp := func() {} // dummy to ensure import
			_ = importOtelHttp
			wrapped := otelhttp.NewHandler(router, "http.server")
			// Attach logger to context after OTel handler
			wrappedWithLogger := logging.RequestLoggingHandler(wrapped)
			http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), wrappedWithLogger)
		})
	})
	cli.Run()
}
