package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/danielgtaylor/huma/v2/humacli"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"

	"github.com/mhirii/rest-template/internal/handlers"
	"github.com/mhirii/rest-template/internal/logging"
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
	// Logger configuration
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(fmt.Sprintf("failed to open log file: %v", err))
	}
	logging.InitLogger(logging.LoggerConfig{
		Level:      zerolog.InfoLevel,
		Format:     "console", // or "json"
		Writers:    []io.Writer{os.Stdout, logFile},
		TimeFormat: "15:04:05",
	})

	l := logging.L()
	l.Info().Msg("Starting API server")
	cli := humacli.New(func(hooks humacli.Hooks, options *Options) {
		router := chi.NewMux()
		api := humachi.New(router, huma.DefaultConfig("Huma API Template", "1.0.0"))

		// Register example REST endpoints
		store := service.NewExampleStore()
		handlers.RegisterExampleRoutes(api, store)

		huma.Get(api, "/greeting/{name}", func(ctx context.Context, input *struct {
			Name string `path:"name" maxLength:"30" example:"world" doc:"Name to greet"`
		}) (*GreetingOutput, error) {
			resp := &GreetingOutput{}
			resp.Body.Message = fmt.Sprintf("Hello, %s!", input.Name)
			return resp, nil
		})

		hooks.OnStart(func() {
			l.Info().Int("port", options.Port).Msg("API server listening")
			http.ListenAndServe(fmt.Sprintf(":%d", options.Port), router)
		})
	})
	cli.Run()
}
