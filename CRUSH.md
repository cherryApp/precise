# Crush Development Guide

## Build/Test/Lint Commands

- **Build**: `go build .` or `task build`
- **Run**: `go run .` or `task dev` (with profiling enabled)
- **Install**: `task install` or `go install -v .`
- **Test**: `task test` or `go test ./...` (run single test: `go test ./internal/llm/prompt -run TestGetContextFromPaths`)
- **Test with Race**: `go test -race ./...`
- **Test with Coverage**: `go test -cover ./...`
- **Update Golden Files**: `go test ./... -update` (regenerates .golden files when test output changes)
  - Update specific package: `go test ./internal/tui/components/core -update`
- **Benchmarks**: `go test -bench=. ./...`
- **Lint**: `task lint-fix` (uses golangci-lint with gofumpt/goimports)
- **Format**: `task fmt` (gofumpt -w .)
- **Generate Schema**: `task schema` (generates JSON schema for configuration)
- **Profiling**: `task profile:cpu`, `task profile:heap`, `task profile:allocs`

## Code Style Guidelines

- **Go Version**: 1.25.0 (minimum requirement)
- **Imports**: Group stdlib first, then external packages, then internal packages
- **Formatting**: Use gofumpt (stricter than gofmt), enabled in golangci-lint
- **Naming**: Standard Go conventions - PascalCase for exported, camelCase for unexported
- **Types**: Prefer explicit types, use type aliases for clarity (e.g., `type AgentName string`)
- **Error handling**: Return errors explicitly, use `fmt.Errorf` for wrapping with `%w` verb
- **Context**: Always pass `context.Context` as first parameter for operations
- **Interfaces**: Define interfaces in consuming packages, keep them small and focused
- **Structs**: Use struct embedding for composition, group related fields
- **Constants**: Use typed constants with iota for enums, group in const blocks
- **Testing**: Use testify's `require` package, parallel tests with `t.Parallel()`,
  `t.SetEnv()` to set environment variables. Always use `t.Tempdir()` when in
  need of a temporary directory (no cleanup needed).
- **JSON tags**: Use snake_case for JSON field names
- **File permissions**: Use octal notation (0o755, 0o644) for file permissions
- **Comments**: End comments in periods unless comments are at the end of the line.
- **Logging**: Use `slog` for structured logging, avoid `log` package

## Testing with Mock Providers

When writing tests that involve provider configurations, use the mock providers to avoid API calls:

```go
func TestYourFunction(t *testing.T) {
    // Enable mock providers for testing
    originalUseMock := config.UseMockProviders
    config.UseMockProviders = true
    defer func() {
        config.UseMockProviders = originalUseMock
        config.ResetProviders()
    }()

    // Reset providers to ensure fresh mock data
    config.ResetProviders()

    // Your test code here - providers will now return mock data
    providers := config.Providers()
    // ... test logic
}
```

## Formatting

- ALWAYS format any Go code you write.
  - First, try `gofumpt -w .`.
  - If `gofumpt` is not available, use `goimports`.
  - If `goimports` is not available, use `gofmt`.
  - You can also use `task fmt` to run `gofumpt -w .` on the entire project,
    as long as `gofumpt` is on the `PATH`.