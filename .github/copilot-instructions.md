---
description: '[ossify] Instructions for writing Go code following idiomatic Go practices and community standards'
applyTo: '**/*.go,**/go.mod,**/go.sum'
---
<!-- See also: https://raw.githubusercontent.com/github/awesome-copilot/refs/heads/main/instructions/go.instructions.md -->
# Copilot Instructions for ossify

This document provides guidelines for GitHub Copilot when working on this repository.

## Project Overview

ossify is a CLI tool written in Go that helps apply conventions and structure to open-source projects. It uses the Cobra library for CLI commands and follows standard Go project conventions.

## Go Guidelines

### General

- Use the latest stable version of Go (currently 1.25)
- Follow the official [Go Code Review Comments](https://go.dev/wiki/CodeReviewComments) style guide
- Use `gofmt` and `golangci-lint` for formatting and linting
- Prefer explicit code over implicit logic
- Always handle errors explicitly; do not use `_` to ignore errors except in deferred cleanup functions
- Always run `go mod tidy` after making changes to `go.mod` and `go.sum`

### Testing

- **Prefer table-driven tests** for testing multiple cases
- Test files must end in `_test.go`
- Keep tests in the same package (white-box testing)
- Use `_test` package suffix for black-box testing
- Tests should test actual behavior, not reproduce implementation logic
- Use `t.Helper()` in test helper functions
- Use `t.Fatalf()` for fatal errors, `t.Errorf()` for non-fatal assertions
- Use `defer func() { _ = cleanup() }()` pattern for cleanup in tests to satisfy errcheck linter
- Create temporary directories with `os.MkdirTemp()` and clean up with deferred `os.RemoveAll()`
- Name tests descriptively using `Test_functionName_scenario`
- Consider using `testify` or similar libraries when they add value, but don't over-complicate simple tests

Example table-driven test:
```go
func TestSomething(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {"valid input", "foo", "bar", false},
        {"empty input", "", "", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := Something(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("Something() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if got != tt.want {
                t.Errorf("Something() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

### Code Style

- Use `switch` statements with explicit cases rather than if-else chains when checking against multiple values of the same variable
- Do not remove `switch` cases or `default` branches without explicit instruction; prefer keeping explicit cases for clarity
- Avoid Yoda conditions (e.g., use `x == "value"` not `"value" == x`)
- Use tagged switch statements when checking the same variable multiple times (satisfies staticcheck QF1003)
- Keep functions focused and small; extract helpers when logic becomes complex

### Error Handling

- Return errors rather than calling `os.Exit()` in library code
- Use `fmt.Errorf("context: %w", err)` for error wrapping
- Use `errors.Join` where possible to combine multiple errors
- Use `errors.New` for simple static errors
- Use `fmt.Errorf` for dynamic errors with formatted messages
- Create custom error types for domain-specific errors
- Export error variables for sentinel errors (e.g., `var ErrNotFound = errors.New("not found")`)
- Use `errors.Is` and `errors.As` for error checking instead of direct comparison
- Check all error return values (errcheck linter)
- In CLI commands, use `cobra.CheckErr` with a wrapped error for consistent error handling:
  ```go
  cobra.CheckErr(fmt.Errorf("loading config: %w", err))
  ```

### Error Propagation

- Add context when propagating errors up the stack
- Don't log and return errors (choose one)
- Handle errors at the appropriate level
- Consider using structured errors for better debugging

### Naming Conventions

- Use MixedCaps or mixedCaps, not underscores
- Acronyms should be all caps (e.g., `HTTP`, `ID`, `URL`)
- Interface names should describe behavior, often ending in `-er` (e.g., `Reader`, `Writer`)
- Package names should be lowercase, single-word, and not plural

### Code Documentation

- Prioritize self-documenting code through clear naming and structure
- All exported functions, types, and constants must have clear concise explanations
- Start documentation with the symbol name
- Use `// Deprecated:` comments for deprecated items
- Write documentation in English by default
- Use examples in documentation when helpful
- Keep documentation close to code
- Update documentation when code changes
- Avoid emoji in documentation and comments

### README and Documentation Files

- Include clear setup instructions
- Document dependencies and requirements
- Provide usage examples
- Document configuration options
- Include troubleshooting section

### Dependencies

- Use Go modules for dependency management
- Keep dependencies minimal and well-maintained
- Prefer standard library over external packages when reasonable

## CLI Commands (Cobra)

- Register commands in `init()` functions
- Use persistent flags for options that apply to subcommands
- Make mutually exclusive flags explicit and validate at runtime with clear error messages
- Provide both short and long descriptions for commands
- Include usage examples in command long descriptions

## Project Structure

```
cmd/           # CLI commands (Cobra)
internal/      # Private application code
  config/      # Configuration management
  model/       # Domain models
  licenses/    # License data and operations
  util/        # Utility functions
```

## Git Commit Messages

- Use conventional commit format: `type(scope): description`
- Types: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`
- Keep the subject line under 72 characters
- Use the body for detailed explanations when needed
- For multi-line commit messages, use `git commit -F <file>` instead of `git commit -m`


## Performance Optimization

### Memory Management

- Minimize allocations in hot paths
- Reuse objects when possible (consider `sync.Pool`)
- Use value receivers for small structs
- Preallocate slices when size is known
- Avoid unnecessary string conversions

### I/O: Readers and Buffers

- Most `io.Reader` streams are consumable once; reading advances state. Do not assume a reader can be re-read without special handling
- If you must read data multiple times, buffer it once and recreate readers on demand:
    - Use `io.ReadAll` (or a limited read) to obtain `[]byte`, then create fresh readers via `bytes.NewReader(buf)` or `bytes.NewBuffer(buf)` for each reuse
    - For strings, use `strings.NewReader(s)`; you can `Seek(0, io.SeekStart)` on `*bytes.Reader` to rewind
- For HTTP requests, do not reuse a consumed `req.Body`. Instead:
    - Keep the original payload as `[]byte` and set `req.Body = io.NopCloser(bytes.NewReader(buf))` before each send
    - Prefer configuring `req.GetBody` so the transport can recreate the body for redirects/retries: `req.GetBody = func() (io.ReadCloser, error) { return io.NopCloser(bytes.NewReader(buf)), nil }`
- To duplicate a stream while reading, use `io.TeeReader` (copy to a buffer while passing through) or write to multiple sinks with `io.MultiWriter`
- Reusing buffered readers: call `(*bufio.Reader).Reset(r)` to attach to a new underlying reader; do not expect it to “rewind” unless the source supports seeking
- For large payloads, avoid unbounded buffering; consider streaming, `io.LimitReader`, or on-disk temporary storage to control memory

- Use `io.Pipe` to stream without buffering the whole payload:
    - Write to `*io.PipeWriter` in a separate goroutine while the reader consumes
    - Always close the writer; use `CloseWithError(err)` on failures
    - `io.Pipe` is for streaming, not rewinding or making readers reusable

- **Warning:** When using `io.Pipe` (especially with multipart writers), all writes must be performed in strict, sequential order. Do not write concurrently or out of order—multipart boundaries and chunk order must be preserved. Out-of-order or parallel writes can corrupt the stream and result in errors.

- Streaming multipart/form-data with `io.Pipe`:
    - `pr, pw := io.Pipe()`; `mw := multipart.NewWriter(pw)`; use `pr` as the HTTP request body
    - Set `Content-Type` to `mw.FormDataContentType()`
    - In a goroutine: write all parts to `mw` in the correct order; on error `pw.CloseWithError(err)`; on success `mw.Close()` then `pw.Close()`
    - Do not store request/in-flight form state on a long-lived client; build per call
    - Streamed bodies are not rewindable; for retries/redirects, buffer small payloads or provide `GetBody`