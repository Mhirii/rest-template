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
```
Visit [http://localhost:8888/greeting/world](http://localhost:8888/greeting/world) to see the example endpoint.

## Configuration
- The server port can be set via CLI flag `--port` or environment variable `SERVICE_PORT` (default: 8888).
- For more configuration, extend the `Options` struct in `main.go`, a config file loader is planned.

## Code Structure
- `cmd/rest-template/`: Main application entrypoint (`main.go`)
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
- API documentation is available at `http://localhost:8888/docs` when the server is running.

---
Feel free to customize this template for your own needs!
