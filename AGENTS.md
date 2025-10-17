# AGENTS.md - Development Guidelines for AI Assistants

This document provides guidance for AI coding agents and developers working on the TogoMQ SDK for Go.

## Project Overview

This is the official Go SDK for TogoMQ, a modern message queue service. The SDK provides a clean, idiomatic Go API for publishing and subscribing to messages via gRPC.

## Architecture

### Core Components

1. **Client (`client.go`)** - Main entry point for SDK users
   - Manages gRPC connection lifecycle
   - Implements Pub and Sub methods
   - Handles authentication and metadata

2. **Configuration (`config.go`)** - Configuration management
   - Flexible option pattern for configuration
   - Validation of configuration parameters
   - Default values for common use cases

3. **Messages (`message.go`)** - Message types and conversions
   - User-friendly message structures
   - Conversion to/from gRPC protobuf types
   - Builder pattern for message construction

4. **Logger (`logger.go`)** - Logging infrastructure
   - Configurable log levels
   - Structured logging for debugging

5. **Errors (`errors.go`)** - Error handling
   - Custom error types with error codes
   - gRPC error wrapping and translation
   - User-friendly error messages

### Dependencies

- **github.com/TogoMQ/togomq-grpc-go** - Auto-generated gRPC protobuf definitions
- **google.golang.org/grpc** - gRPC client library
- Standard library only for everything else

## Development Principles

### 1. API Design

- **Simplicity First**: The API should be simple and intuitive for common use cases
- **Builder Pattern**: Use fluent builder pattern for complex configurations
- **Context Awareness**: All network operations should accept `context.Context`
- **Channel-Based Streaming**: Use Go channels for streaming operations
- **Idiomatic Go**: Follow Go best practices and conventions

### 2. Error Handling

- Always return errors, never panic
- Wrap gRPC errors with meaningful context
- Provide error codes for programmatic handling
- Include detailed error messages for debugging

### 3. Testing

- Unit tests for all public APIs
- No mocking of the TogoMQ gRPC library (test with real types)
- Table-driven tests for multiple scenarios
- Race detection enabled in CI

### 4. Logging

- Configurable log levels
- Log at appropriate levels:
  - Debug: Detailed operation info
  - Info: Important state changes
  - Warn: Recoverable issues
  - Error: Operation failures
- Never log sensitive data (tokens, message content)

## Adding New Features

### Adding a New Configuration Option

1. Add field to `Config` struct in `config.go`
2. Update `DefaultConfig()` with default value
3. Add `With*` option function
4. Update `Validate()` if needed
5. Add tests in `config_test.go`
6. Update README.md configuration section

### Adding a New Client Method

1. Add method to `Client` struct in `client.go`
2. Accept `context.Context` as first parameter
3. Add authentication metadata via `contextWithAuth()`
4. Implement proper error handling with `WrapGRPCError()`
5. Add logging at appropriate levels
6. Create tests
7. Add usage examples to README.md

### Adding Message Fields

1. Update `Message` struct in `message.go`
2. Add builder method if needed (e.g., `WithFieldName()`)
3. Update conversion functions (`toPubRequest()`, `fromSubResponse()`)
4. Add tests in `message_test.go`
5. Document in README.md

## Code Quality Standards

### Required Checks

1. **Format**: Code must pass `gofmt -s`
2. **Lint**: Code must pass `golangci-lint run` with no errors
3. **Tests**: All tests must pass with race detection enabled
4. **Coverage**: Maintain or improve test coverage

**Important**: Always run `golangci-lint run` locally before committing to catch issues early.

### Best Practices

- Use meaningful variable names
- Add comments for exported types and functions
- Keep functions focused and small
- Avoid global state
- Prefer composition over inheritance
- Use interfaces where appropriate

## CI/CD Pipeline

The project uses GitHub Actions for continuous integration:

1. **Format Check**: Ensures code is properly formatted
2. **Linting**: Runs golangci-lint for code quality
3. **Tests**: Runs all tests with race detection and coverage

All checks run on every push to any branch.

## Common Tasks

### Running Tests Locally

```bash
# Run all tests
go test ./...

# Run with race detection
go test -race ./...

# Run with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Format Code

```bash
# Format all files
gofmt -s -w .

# Check formatting
gofmt -s -l .
```

### Lint Code

```bash
# Install golangci-lint (if not already installed)
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linter (uses .golangci.yml configuration)
golangci-lint run

# Run linter with auto-fix for some issues
golangci-lint run --fix

# Run linter on specific files or directories
golangci-lint run ./path/to/code

# Run linter with verbose output
golangci-lint run -v
```

**Note**: The project includes a `.golangci.yml` configuration file that:
- Enables standard linters (gofmt, govet, errcheck, staticcheck, unused, gosimple, ineffassign, goimports, misspell, revive)
- Excludes deprecated linters (structcheck, varcheck, deadcode - replaced by unused)
- Allows `TogoMQError` type name (even though it stutters with package name)
- Excludes errcheck warnings in test files and examples

### Pre-Commit Checklist

Before committing code, ensure:

```bash
# 1. Format code
gofmt -s -w .

# 2. Run linter
golangci-lint run

# 3. Run tests with race detection
go test -race ./...

# 4. Verify no errors in all three steps above
```

## Future Enhancements

### Potential Features to Add

1. **Retry Logic**
   - Automatic retry on transient failures
   - Configurable retry policies
   - Exponential backoff

2. **Connection Pooling**
   - Pool of gRPC connections for high throughput
   - Load balancing across connections

3. **Metrics and Monitoring**
   - Prometheus metrics integration
   - Message throughput tracking
   - Error rate monitoring

4. **Advanced Message Features**
   - Message compression
   - Message encryption
   - Priority queues
   - Dead letter queues

5. **Testing Utilities**
   - Mock server for testing
   - Test helpers for common scenarios
   - Benchmark suite

6. **Additional Methods**
   - CountMessages (already in gRPC API)
   - HealthCheck (already in gRPC API)
   - Message acknowledgement
   - Batch operations

## Breaking Changes

When making breaking changes:

1. Increment major version (following semver)
2. Document migration path in CHANGELOG
3. Consider deprecation period if possible
4. Update all examples and documentation

## Questions and Support

For questions about development:

1. Check existing issues on GitHub
2. Review the TogoMQ gRPC proto definitions
3. Consult the main TogoMQ documentation
4. Open a discussion on GitHub

## Version History

- **v0.1.0** (Initial Release)
  - Basic Pub/Sub functionality
  - Configuration management
  - Error handling
  - Comprehensive logging
  - Full test coverage
  - CI/CD pipeline

## References

- [TogoMQ Documentation](https://togomq.io/docs)
- [gRPC Go Documentation](https://grpc.io/docs/languages/go/)
- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
