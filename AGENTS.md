# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Sentire is a CLI tool for the Sentry API written in Go, providing comprehensive access to Sentry's debugging data including complete stack traces, contexts, and event details. The tool features both 1:1 API mapping commands and user-friendly custom commands like `inspect` for parsing Sentry URLs.

## Build and Development Commands

```bash
# Building
make build                    # Standard build
make build-release            # Optimized build with stripped symbols
go build -o sentire ./cmd/sentire  # Direct go build

# Testing
make test                     # Run all tests
make test-coverage           # Run tests with HTML coverage report
go test ./tests/ -v          # Run tests with verbose output
go test ./tests/inspect_test.go ./tests/events_test.go -v  # Run specific test files

# Development
make fmt                     # Format all Go code
make lint                    # Run golangci-lint (requires installation)
make deps                    # Download and tidy dependencies
make clean                   # Remove build artifacts

# Cross-compilation
make build-all               # Build for Linux, macOS (Intel/ARM), and Windows
```

## Code Architecture

### Layer Structure
- **`cmd/sentire/main.go`**: Entry point that calls CLI executor
- **`internal/cli/`**: Cobra-based CLI commands and user interface
- **`internal/api/`**: Sentry API method implementations (1:1 mapping)
- **`internal/client/`**: HTTP client with auth, rate limiting, and pagination
- **`pkg/models/`**: Comprehensive data models matching Sentry's API responses
- **`tests/`**: Test suite with mock servers and integration tests

### Key Design Patterns

**CLI Command Registration**: Each command group (events, org, projects, inspect) has its own file in `internal/cli/` and registers with the root command via `init()` functions.

**API Client Architecture**: The `internal/client/client.go` provides a base HTTP client that handles:
- Bearer token authentication via `SENTRY_API_TOKEN` env var
- Rate limit tracking from Sentry's response headers
- Cursor-based pagination parsing from Link headers
- Proper error handling with meaningful messages

**Data Models**: The `pkg/models/` contain comprehensive structs that capture ALL available data from Sentry APIs:
- **Event model**: Complete debugging data including stack traces, breadcrumbs, contexts, exceptions
- **Issue model**: Enhanced with priority, substatus, culprit, ownership, and detailed metadata
- **Project model**: All capability flags, insights flags, and configuration options
- **Organization models**: Detailed statistics with category breakdowns

### Special Features

**Inspect Command**: Custom command (`internal/cli/inspect.go`) that parses Sentry URLs using regex to extract organization and issue ID, then automatically fetches the "recommended" event with full debugging context.

**Pagination Handling**: Automatic pagination support with `--all` flag that continues fetching until all results are retrieved using cursor-based pagination.

**Comprehensive Event Data**: Unlike many Sentry tools, this captures complete event data including stack frames with context lines, variable values, breadcrumbs trail, and all debugging contexts.

## Authentication and Configuration

Set `SENTRY_API_TOKEN` environment variable before using any commands. The client validates this on startup and provides clear error messages if missing.

## Testing Strategy

**Mock Server Testing**: Tests use `httptest.NewServer` to create mock Sentry API endpoints, allowing comprehensive testing without live API calls.

**Integration Testing**: Tests verify complete request/response cycles including URL construction, header parsing, and JSON marshaling/unmarshaling.

**Error Case Testing**: Comprehensive error handling tests for invalid URLs, missing authentication, API errors, and malformed responses.

## Model Field Mapping

When adding new API endpoints or updating existing ones, ensure data models capture ALL available fields from Sentry's API documentation. The project prioritizes completeness over simplicity - users should get access to all debugging data available through the API.

Critical model components:
- **Event.Entries**: Contains stack traces, exceptions, breadcrumbs
- **Event.Contexts**: Browser, OS, runtime, device information  
- **Issue metadata**: Now uses `interface{}` to capture complex nested structures
- **Project capability flags**: All `has*` and `hasInsights*` boolean flags

## Output and User Experience

All commands output JSON by default with proper indentation. The `inspect` command demonstrates the user-friendly approach - parse URLs that users commonly encounter (from Slack notifications, emails) and provide immediate access to debugging data.

Commands follow the pattern: `./sentire <resource> <action> [args] [flags]` with the exception of the custom `inspect` command which prioritizes ease of use over API consistency.

## New release

When asked to create a new release please refer to @.opencode/agent/new-release.md
