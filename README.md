# Sentire

<img src="sentire_logo.png" width="50%" alt="Sentire">

A simple and user-friendly command-line interface for the Sentry API, written in Go.

## Overview

Sentire provides an intuitive CLI for interacting with Sentry's API. It covers the essential Sentry operations including managing events, issues, projects, and organizations.

**Key Features:**
- Multiple output formats: JSON, table, text, and markdown
- Comprehensive API coverage for debugging workflows
- Built-in pagination and rate limiting
- Human-readable output for terminal usage
- Machine-readable JSON for scripting

## Installation

### Prerequisites

- A Sentry API token with appropriate permissions

### Using Homebrew (macOS)

```bash
brew install andreagrandi/tap/sentire
```

### Using Go install

```bash
go install github.com/andreagrandi/sentire/cmd/sentire@latest
```

### Download pre-built binaries

Download the latest release for your platform from the [releases page](https://github.com/andreagrandi/sentire/releases).

### Building from source

Prerequisites: Go 1.25 or later

```bash
git clone <repository-url>
cd sentire
go build -o sentire ./cmd/sentire
```

## Development Setup

### Git Hooks

To prevent formatting issues and ensure code quality, install the Git hooks:

```bash
# Install pre-commit hook for automatic Go formatting
./scripts/install-hooks.sh
```

The pre-commit hook will automatically format Go files using `gofmt -s` before each commit, ensuring consistency with CI requirements.

To format all Go files manually:
```bash
make fmt
```

To bypass the hook temporarily (not recommended):
```bash
git commit --no-verify
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

### URL Inspection

Sentire includes a special `inspect` command that can parse Sentry URLs directly:

```bash
# Inspect a Sentry issue URL and get detailed event information
sentire inspect "https://my-org.sentry.io/issues/123456789/"

# Get inspection results in table format
sentire inspect "https://my-org.sentry.io/issues/123456789/" --format table

# Get inspection results in markdown for documentation
sentire inspect "https://my-org.sentry.io/issues/123456789/" --format markdown
```

This command automatically extracts the organization and issue ID from the URL and fetches the most relevant debugging information.

### Command Options

Most list commands support these common options:

- `--all`: Fetch all pages of results (default: single page)
- `--format <format>`: Output format (json, table, text, markdown) - default: json
- `--verbose`: Enable verbose output

#### Output Formats

Sentire supports multiple output formats to suit different use cases:

- **`json`** (default): Machine-readable JSON format, ideal for scripting and automation
- **`table`**: Human-readable table format with borders, perfect for terminal viewing
- **`text`**: Clean plain text format, great for simple parsing and readability
- **`markdown`**: Documentation-friendly markdown format, useful for reports and documentation

**Format Examples:**
```bash
# Default JSON output
sentire events list-issues my-org

# Human-readable table format
sentire events list-issues my-org --format table

# Simple text format
sentire projects list --format text

# Markdown for documentation
sentire org stats my-org --format markdown
```

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
# JSON output (default)
sentire events list-issues my-org --query="is:unresolved issue.priority:[high,medium]" --period=7d

# Table format for better readability
sentire events list-issues my-org --query="is:unresolved issue.priority:[high,medium]" --period=7d --format table
```

### Get all events for a project in the last 24 hours

```bash
# Get all events with table format
sentire events list-project my-org my-project --period=24h --all --full --format table
```

### Get organization statistics for the last week

```bash
# Statistics in markdown format for reports
sentire org stats my-org --field="sum(quantity)" --period=7d --project=123 --project=456 --format markdown
```

### List issues in production environment

```bash
# Production issues in text format
sentire events list-issues my-org --environment=production --period=24h --format text
```

### Output Format Comparison

**Table format** (great for terminal viewing):
```
┌──────┬─────────────────────┬───────┬──────────┬─────────┬──────────────────┐
│  ID  │        TITLE        │ LEVEL │ STATUS   │  COUNT  │   LAST SEEN      │
├──────┼─────────────────────┼───────┼──────────┼─────────┼──────────────────┤
│ 1234 │ TypeError in login  │ error │ unresov. │   45    │ 2025-08-30 10:15 │
│ 1235 │ API timeout         │ warn  │ resolved │   12    │ 2025-08-30 09:30 │
└──────┴─────────────────────┴───────┴──────────┴─────────┴──────────────────┘
```

**Text format** (simple and scriptable):
```
Issues (2 total):

1. Issue #1234
   Title: TypeError in login component
   Level: error | Status: unresolved | Count: 45
   Project: web-app | Users: 23
   Last Seen: 2025-08-30 10:15

2. Issue #1235
   Title: API timeout on user endpoint
   Level: warning | Status: resolved | Count: 12
   Project: api-service | Users: 8
   Last Seen: 2025-08-30 09:30
```

**Markdown format** (documentation-ready):
```markdown
# Issues (2 total)

| ID | Title | Level | Status | Count | Users | Last Seen | Project |
|----|-------|-------|--------|-------|-------|-----------|---------|
| 1234 | TypeError in login... | error | unresolved | 45 | 23 | 08-30 10:15 | web-app |
| 1235 | API timeout on use... | warning | resolved | 12 | 8 | 08-30 09:30 | api-service |
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

## License

Licensed under the MIT License. See [LICENSE](LICENSE) for details.

## Roadmap

Future enhancements may include:

- ✅ **Multiple output formats** (JSON, table, text, markdown) - **COMPLETED**
- Configuration file support
- Additional Sentry API endpoints
- Webhook management
- Release management
- Performance monitoring queries
- Export functionality (CSV, JSON files)
- Interactive mode for complex queries
- Custom output format templates
- Shell auto-completion support