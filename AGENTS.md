# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Session Start Workflow

Before making code or documentation changes in this repo:

1. Switch back to `master`.
2. Pull the latest changes with `git pull --ff-only`.
3. Create a new branch with a short, descriptive name related to the feature being added or the bug being fixed.
4. Make the requested changes on that branch.

Do not start work from an old feature branch unless the user explicitly asks to continue that branch.

## Adding a Ticket, Issue, or Bug

When the user asks to add a "ticket", "issue", or "bug" to this repo, all of the
steps below are required — creating the GitHub issue alone is not enough.

**1. Create the issue** with the repo label and an area label:

```bash
gh issue create --repo andreagrandi/sentire \
  --title "<concise title>" \
  --body "<description>" \
  --label "sentire" \
  --label "area:<area>"
```

- Always apply the `sentire` label — every work item in this repo carries it.
- Add the matching `area:<area>` label when one exists: `area:feature`,
  `area:ux`, `area:agent`, `area:docs`, `area:security`, `area:testing`,
  `area:release`. The `Reliability` and `Packaging` areas have no label; for
  those, set only the project Area field (step 3).
- Add a type label when it fits: `bug` (bug report), `enhancement` (feature
  request), or `documentation` (docs-only work).

**2. Add the issue to the "CLI Tools" project** and capture the item ID
(https://github.com/users/andreagrandi/projects/1):

```bash
ITEM_ID=$(gh project item-add 1 --owner andreagrandi \
  --url <issue-url> --format json --jq .id)
```

**3. Set Priority, Area, and Status** on the project item. The project and
field IDs are identical for every repo on this board:

- Project ID: `PVT_kwHOAAm1584BYDlZ`
- Priority — field `PVTSSF_lAHOAAm1584BYDlZzhTMDck`:
  High `ed3787e3`, Medium `3e3ea407`, Low `994234f4`
- Area — field `PVTSSF_lAHOAAm1584BYDlZzhTMDco`:
  Reliability `6595432d`, Packaging `6895c50a`, UX `2bc024bb`,
  Testing `0d5bc016`, Feature `6390f97d`, Agent `3a2d6f7e`,
  Docs `f5c50514`, Security `062d12a3`, Release `b344aeab`
- Status — field `PVTSSF_lAHOAAm1584BYDlZzhTMDNQ`:
  Todo `f75ad846`, In Progress `47fc9ee4`, Done `98236657`

```bash
# Priority — always set it
gh project item-edit --id "$ITEM_ID" --project-id PVT_kwHOAAm1584BYDlZ \
  --field-id PVTSSF_lAHOAAm1584BYDlZzhTMDck \
  --single-select-option-id <priority-option-id>

# Area — match the area:<area> label from step 1
gh project item-edit --id "$ITEM_ID" --project-id PVT_kwHOAAm1584BYDlZ \
  --field-id PVTSSF_lAHOAAm1584BYDlZzhTMDco \
  --single-select-option-id <area-option-id>

# Status — new tickets start as Todo
gh project item-edit --id "$ITEM_ID" --project-id PVT_kwHOAAm1584BYDlZ \
  --field-id PVTSSF_lAHOAAm1584BYDlZzhTMDNQ \
  --single-select-option-id f75ad846
```

If the user does not state a priority or area, ask before creating the issue.
Follow the conventions of existing project issues — do not invent new labels or
fields.

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

## Changelog

`CHANGELOG.md` is the canonical record of user-visible changes. Update it in
the same change that introduces the behaviour — not as a follow-up.

### When to add an entry

Add an entry when a change affects any of:

- CLI surface (new/renamed/removed commands, flags, arguments, defaults)
- Observable behaviour (output shape, exit codes, error messages, env vars)
- Packaging or installation (release artifacts, Homebrew tap, `go install` path)
- Configuration file format or supported locations
- Dependency bumps that ship in the binary or change the minimum Go version
- Agent-facing contracts: `describe` schema, `context` output, CONTEXT.md guidance

Skip the changelog only for changes with no user-visible effect: internal
refactors, test-only changes, CI tweaks that do not ship, formatting, or
edits to private agent instructions (AGENTS.md / .opencode/).

### How to add an entry

- Add the entry under `## [Unreleased]` at the top of `CHANGELOG.md`. Do not
  invent a new version heading — the release workflow promotes `[Unreleased]`
  when the version is bumped.
- Use the existing subsections: `Added`, `Changed`, `Fixed`, `Technical`.
  Create a subsection only if it is missing and your change needs it.
- Keep entries to a single line each. Lead with the user-facing effect, not
  the implementation detail. Reference the command or flag explicitly when
  relevant (e.g. ``Add `--all` flag to `events list-issues```).
- Group related changes into one entry rather than splitting per-commit.

Before opening a PR, check that `CHANGELOG.md` has an `[Unreleased]` entry
covering your change, or be ready to explain in the PR why no entry is
needed.

## New release

When asked to create a new release please refer to @.opencode/agent/new-release.md
