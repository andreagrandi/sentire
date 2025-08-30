# Sentire

A simple and user-friendly command-line interface for the Sentry API, written in Go.

## Overview

Sentire provides an intuitive CLI for interacting with Sentry's API. It covers the essential Sentry operations including managing events, issues, projects, and organizations.

## Installation

### Prerequisites

- Go 1.25 or later
- A Sentry API token with appropriate permissions

### Building from source

```bash
git clone <repository-url>
cd sentire
go build -o sentire ./cmd/sentire
```

## Configuration

Before using sentire, you must set your Sentry API token as an environment variable:

```bash
export SENTRY_API_TOKEN=your_sentry_api_token_here
```

You can obtain an API token from your Sentry organization settings under "Auth Tokens".

## Usage

### Basic Commands

```bash
# Show help
sentire --help

# List all your projects
sentire projects list

# Get a specific project
sentire projects get <organization> <project>

# List organization projects
sentire org list-projects <organization>

# Get organization statistics
sentire org stats <organization> --field="sum(quantity)"
```

### Events and Issues

```bash
# List events for a project
sentire events list-project <organization> <project> --period=24h

# List events for an issue
sentire events list-issue <organization> <issue-id> --full

# List issues for an organization
sentire events list-issues <organization> --query="is:unresolved"

# Get a specific event
sentire events get-event <organization> <project> <event-id>

# Get a specific issue
sentire events get-issue <organization> <issue-id>

# Get an issue event (latest, oldest, recommended, or specific ID)
sentire events get-issue-event <organization> <issue-id> latest
```

### Command Options

Most list commands support these common options:

- `--all`: Fetch all pages of results (default: single page)
- `--format json`: Output format (currently supports JSON)
- `--verbose`: Enable verbose output

#### Time-based Filtering

Many commands support time filtering options:

- `--period <period>`: Relative time period (e.g., "1h", "24h", "7d", "30d")
- `--start <iso-date>`: Start time in ISO-8601 format
- `--end <iso-date>`: End time in ISO-8601 format

#### Environment and Project Filtering

- `--environment <env>`: Filter by environment (can be specified multiple times)
- `--project <id>`: Filter by project ID (can be specified multiple times)

## Examples

### List recent high-priority issues

```bash
sentire events list-issues my-org --query="is:unresolved issue.priority:[high,medium]" --period=7d
```

### Get all events for a project in the last 24 hours

```bash
sentire events list-project my-org my-project --period=24h --all --full
```

### Get organization statistics for the last week

```bash
sentire org stats my-org --field="sum(quantity)" --period=7d --project=123 --project=456
```

### List issues in production environment

```bash
sentire events list-issues my-org --environment=production --period=24h
```

## API Coverage

Sentire currently supports the following Sentry API endpoints:

### Events
- ✅ List project events (`/projects/{org}/{project}/events/`)
- ✅ List issue events (`/organizations/{org}/issues/{issue}/events/`)  
- ✅ List organization issues (`/organizations/{org}/issues/`)
- ✅ Get project event (`/projects/{org}/{project}/events/{event}/`)
- ✅ Get issue (`/organizations/{org}/issues/{issue}/`)
- ✅ Get issue event (`/organizations/{org}/issues/{issue}/events/{event}/`)

### Organizations
- ✅ List organization projects (`/organizations/{org}/projects/`)
- ✅ Get organization statistics (`/organizations/{org}/stats-summary/`)

### Projects
- ✅ List all projects (`/projects/`)
- ✅ Get project (`/projects/{org}/{project}/`)

## Rate Limiting

Sentire automatically handles Sentry's rate limiting by:

- Tracking rate limit headers from API responses
- Displaying current rate limit status in verbose mode
- Implementing proper error handling for rate limit exceeded scenarios

## Error Handling

The CLI provides clear error messages for common scenarios:

- Missing or invalid API token
- API authentication failures
- Resource not found errors
- Network connectivity issues
- Rate limiting

## Testing

Run the test suite:

```bash
go test ./tests/ -v
```

The tests include:
- Unit tests for the HTTP client
- Integration tests for all API methods
- Mock server tests for CLI commands

## Contributing

1. Fork the repository
2. Create a feature branch
3. Add tests for your changes
4. Ensure all tests pass
5. Submit a pull request

## License

[Add your license information here]

## Roadmap

Future enhancements may include:

- Table output format for better readability
- Configuration file support
- Additional Sentry API endpoints
- Webhook management
- Release management
- Performance monitoring queries
- Export functionality (CSV, JSON files)
- Interactive mode for complex queries