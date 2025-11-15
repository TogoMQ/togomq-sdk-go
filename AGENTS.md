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

### Public API Methods

The Client exposes the following public methods:

- **Pub** - Publishes messages via streaming channel
- **PubBatch** - Publishes a batch of messages at once
- **Sub** - Subscribes to messages from topics (supports wildcards)
- **CountMessages** - Counts messages in a topic (supports wildcards)
- **Close** - Closes the client connection

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

# Verify configuration is valid
golangci-lint config verify

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

**Important**: Always verify the configuration is valid before committing with `golangci-lint config verify`

### Pre-Commit Checklist

Before committing code, ensure:

```bash
# 1. Format code
gofmt -s -w .

# 2. Verify linter configuration (if .golangci.yml was modified)
golangci-lint config verify

# 3. Run linter
golangci-lint run

# 4. Run tests with race detection
go test -race ./...

# 5. Verify no errors in all steps above
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

## Release Process

The project uses **automated releases** with conventional commits. Releases are created automatically when code is merged to the `main` branch.

### Conventional Commits

Use conventional commit messages to control versioning:

**Commit Message Format:**
```
<type>(<scope>): <subject>

<body>

<footer>
```

**Types and Version Bumps:**

- `feat:` or `feature:` - New feature → **Minor version bump** (v1.0.0 → v1.1.0)
- `fix:` or `bugfix:` - Bug fix → **Patch version bump** (v1.0.0 → v1.0.1)
- `BREAKING CHANGE:` - Breaking change → **Major version bump** (v1.0.0 → v2.0.0)
- `docs:`, `test:`, `chore:`, `style:`, `refactor:` - No version bump (unless specified)

**Examples:**

```bash
# Minor version bump (new feature)
git commit -m "feat: add retry mechanism for failed publishes"

# Patch version bump (bug fix)
git commit -m "fix: correct race condition in subscription handler"

# Major version bump (breaking change)
git commit -m "feat: redesign client API

BREAKING CHANGE: removed WithQueue method, use WithTopic instead"

# Multiple changes
git commit -m "feat: add message batching

- Implement batch send optimization
- Add configurable batch size
- Update documentation"

# No version bump
git commit -m "docs: update installation instructions"
git commit -m "test: add integration tests"
git commit -m "chore: update dependencies"
```

### Release Workflow

1. **Develop**: Make changes in a feature branch
2. **Commit**: Use conventional commit messages
3. **Create PR**: Open pull request to `main`
4. **Review & Merge**: After approval, merge to `main`
5. **Automatic Release**: GitHub Actions automatically:
   - Runs all tests and linters
   - Analyzes commit messages
   - Creates new version tag
   - Generates changelog
   - Creates GitHub release
   - Comments on associated PRs

### Version Guidelines

Follow semantic versioning (semver):

**Major (v1.0.0 → v2.0.0)**: Breaking API changes
- Removing public methods or types
- Changing method signatures
- Removing configuration options
- Changing default behavior that breaks existing code

**Minor (v1.0.0 → v1.1.0)**: New features, backward compatible
- Adding new methods or types
- Adding new configuration options
- Adding new functionality without breaking existing code

**Patch (v1.0.0 → v1.0.1)**: Bug fixes, backward compatible
- Fixing bugs
- Performance improvements
- Documentation updates
- Internal refactoring

### Pre-Release Checklist

Before merging to `main`:

- [ ] All tests pass locally and in CI
- [ ] golangci-lint shows no errors
- [ ] Code is properly formatted
- [ ] README.md is up to date
- [ ] Examples are tested and working
- [ ] Commit messages follow conventional format
- [ ] Breaking changes are clearly documented in commit message

### Go Module Usage

After release, users can install specific versions:

```bash
# Install specific version
go get github.com/TogoMQ/togomq-sdk-go@v1.2.3

# Install latest version
go get github.com/TogoMQ/togomq-sdk-go@latest

# List available versions
go list -m -versions github.com/TogoMQ/togomq-sdk-go
```

### First Release

To create the first release (v0.1.0), manually create and push a tag:

```bash
git tag -a v0.1.0 -m "Initial release"
git push origin v0.1.0
```

Or create a release through GitHub UI. After the first tag exists, all subsequent releases will be automatic.

### Monitoring Releases

- Check the **Actions** tab in GitHub to see release workflow status
- View **Releases** page to see published versions
- Automated comments will be added to PRs when they're released

### Troubleshooting Releases

**No release created after merge:**
- Check commit messages follow conventional format
- Verify workflow ran successfully in Actions tab
- Ensure commits have types that trigger version bumps (feat, fix, BREAKING CHANGE)

**Wrong version number:**
- Review commit messages for correct type prefixes
- Check for BREAKING CHANGE footer for major bumps
- Verify conventional commit format

**Tests failed:**
- Release is blocked if tests fail
- Fix issues and push again
- Workflow will retry on next push

## References

- [TogoMQ Documentation](https://togomq.io/docs)
- [gRPC Go Documentation](https://grpc.io/docs/languages/go/)
- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
