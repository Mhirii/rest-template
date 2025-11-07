# AGENTS.md

## Build, Lint, and Test Commands
- **Build:** `make build` (binary in `bin/`)
- **Run:** `make run` or `cd cmd/rest-template && go run main.go`
- **Test all:** `make test` or `go test ./...`
- **Test single package:** `go test ./internal/handlers` (or any package path)
- **Test single file:** `go test -run TestFuncName ./internal/handlers/example.handler.go`
- **Integration tests:** Run `test_endpoints.sh` (requires running server)
- **Tidy modules:** `make tidy`

## Code Style Guidelines
- **Imports:** Group stdlib, third-party, and local packages; use Go's import formatting.
- **Formatting:** Always run `gofmt` or `go fmt ./...` before committing.
- **Types:** Prefer explicit types; use structs for request/response models.
- **Naming:** Use `CamelCase` for types, `mixedCaps` for variables/functions; acronyms are capitalized (e.g., `API`, `ID`).
- **Error Handling:** Return errors, use context-aware logging (`zerolog.Ctx(ctx)`), and prefer Huma error helpers (e.g., `huma.Error404NotFound`).
- **Comments:** Use Go doc comments for exported types/functions; keep endpoint comments concise.
- **Concurrency:** Use sync primitives (e.g., `sync.RWMutex`) for shared state.
- **File Structure:** Handlers in `internal/handlers`, DTOs in `internal/dto`, services in `internal/service`.
- **Tests:** Place unit tests in the same package; integration/e2e tests in `test/` or scripts.

## Huma
This project uses Huma as the API framework, All official HUMA documentation can be found under docs/huma

---
Agents should follow these conventions for all code contributions in this repository.
