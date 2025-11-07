# Go Huma API Template

A reusable template for building backend APIs in Go using the [Huma](https://huma.rocks/) framework. This template provides a clean project structure, a minimal Huma setup, and best practices for scalable Go API development.

## Purpose
- Kickstart new Go API projects with Huma
- Encourage clean architecture and modular code
- Provide a working example endpoint and CLI

## Quick Start
```sh
# Clone this template
git clone https://github.com/Mhirii/rest-template.git
cd rest-template

# Initialize your own Go module (update module name in go.mod)
go mod tidy

# Change the module name in go.mod to your own repo
# Update import paths in all files from 'github.com/mhirii/rest-template' to your module path
# Update README, LICENSE, and other metadata as needed

# Run the API server
cd cmd/rest-template
GO111MODULE=on go run main.go

# Or use the Makefile from the project root:
make run ARGS="--port=8000 --config=/path/to/configdir"
```
Visit [http://localhost:8888/greeting/world](http://localhost:8888/greeting/world) to see the example endpoint.

You can pass any supported CLI flags using the ARGS variable with make, e.g.:
```
make run ARGS="--port=8080 --log_level=debug"
```
This will start the server on port 8080 with debug logging.

## Configuration
- All configuration options are defined once in `internal/config/config.go` in the `Config` struct, with defaults set in `defaultConfig`.
- You can set config values via:
  - **CLI flags** (highest priority, e.g. `--port=8000 --log_level=debug --config=/path/to/configdir`)
  - **YAML config file** (`config.yaml` in the directory specified by `--config` or `CONFIG_PATH` env var, or current directory by default)
  - **Environment variables** (e.g. `SERVICE_PORT`, `LOG_LEVEL`)
  - **Defaults** (only set once in code)
- To add new config options, simply add a field to the `Config` struct and update `defaultConfig`.
- All flags, env vars, and YAML keys are bound automatically using struct tags and Viper.
- Example config file path usage:
  - `go run main.go --config=config.yaml --port=8000` (pass a file)
  - `go run main.go --config=/path/to/configdir --port=8000` (pass a directory)
  - Or set `CONFIG_PATH` env var before running.
- The config loader will automatically detect if you pass a file or a directory for `--config` and load the config accordingly.

## Code Structure
- `cmd/rest-template/`: Main application entrypoint (`main.go`)
- `internal/config/`: Centralized config loader and struct
- `internal/handlers/`: HTTP route registration and handler logic
- `internal/dto/`: Data Transfer Objects for request/response models
- `internal/service/`: Business logic and in-memory data store
- `internal/logging/`: Global logger setup and request logging middleware
- `internal/observability/`: OpenTelemetry tracing setup
- `internal/middleware/`: Additional middleware (e.g., logging)
- `docs/huma/`: Official Huma documentation and API docs

## What to Change After Cloning
- Update `go.mod` module name and all import paths
- Change project metadata (README, LICENSE, etc.)
- Replace example endpoints and DTOs with your own
- Add your own business logic and configuration

## Documentation
- All official Huma documentation is available in the `docs/huma` directory.
- API documentation is available at `http://<address>:<port>/docs` when the server is running.

---
Feel free to customize this template for your own needs!
