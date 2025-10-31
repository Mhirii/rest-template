# Go Huma API Template

A reusable template for building backend APIs in Go using the [Huma](https://huma.rocks/) framework. This template provides a clean project structure, a minimal Huma setup, and best practices for scalable Go API development.

## Purpose

- Kickstart new Go API projects with Huma
- Encourage clean architecture and modular code
- Provide a working example endpoint and CLI

## Quick Start

```sh
# Clone this template
# git clone https://github.com/yourusername/rest-template.git
cd rest-template

# Initialize your own Go module (optional)
go mod tidy

# Run the API server
cd cmd/rest-template
GO111MODULE=on go run main.go
```

Visit [http://localhost:8888/greeting/world](http://localhost:8888/greeting/world) to see the example endpoint.

## Project Structure

```
rest-template/
├── cmd/rest-template/      # Main application entrypoint (main.go)
├── internal/api/           # API route handlers and logic
├── internal/config/        # Configuration loading and management
├── internal/service/       # Business logic/services
├── pkg/                    # Exported reusable packages (if any)
├── test/                   # Integration or e2e tests
├── go.mod
├── go.sum
├── .gitignore
├── README.md
├── LICENSE
├── Makefile
└── todo.md
```

## Documentation

- All official Huma documentation is available in the `docs/huma` directory.

---

Feel free to customize this template for your own needs!
