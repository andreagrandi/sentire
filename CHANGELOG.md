# Changelog

All notable changes to Sentire will be documented in this file.

## [Unreleased]

### Added
- `SENTRY_API_BASE_URL` environment variable to target a self-hosted Sentry instance
- Agent workflow recipes in `context` output and `CONTEXT.md` (triage, URL inspection, release reporting)
- README section pointing agents at `sentire context` and `sentire describe`
- Test coverage locking the `describe` output schema and the agent context guide
- Contributor and agent guidance for keeping `CHANGELOG.md` updated, plus a PR template prompting for an `[Unreleased]` entry
- `docs/workflows.md` with recipe-style triage, inspection, and reporting workflows plus a troubleshooting checklist for auth, permissions, filters, and empty results
- `Ctrl+C` (SIGINT) and SIGTERM cancel in-flight API requests, reported with a `canceled` error code

### Changed
- Request timeouts surface a clear `timeout` error code instead of a generic transport failure
- Issue lists (`table`, `text`, `markdown`) now show a `Priority` column and relative `Last Seen` times for faster incident triage
- Event lists (`table`, `text`, `markdown`) show a relative `Created` time; single issue and event views append the relative time to absolute timestamps

### Fixed
- Redact `SENTRY_API_TOKEN` from error and verbose output

### Technical
- Use full GitHub module path (`github.com/andreagrandi/sentire`) for public `go install`
- Bump `github.com/spf13/cobra` from 1.8.1 to 1.10.2
- Lower `go.mod` minimum to Go 1.24 and add Go 1.24/1.25 to the CI build matrix
- API requests use context-aware HTTP request creation, propagating cancellation and deadlines from CLI commands into the client
- Add `make vet` target and CI step running `go vet ./...`

## [0.3.0] - 2026-03-07

### Added
- `context` command and CONTEXT.md agent skill file
- `describe` command for schema introspection
- `--format ndjson` for newline-delimited JSON output
- `--fields` flag for JSON output filtering
- Input validation for slugs, IDs, and URLs
- Structured errors with typed codes and distinct exit codes

### Fixed
- TestNewClient failing when config file exists

### Technical
- Bump Go from 1.25 to 1.26.1
- Bump GitHub Actions dependencies (actions/cache, actions/checkout, actions/setup-go)

## [0.2.0] - 2025-01-03

### Added
- Configuration file support for storing API tokens and settings
- Version information command with build details
- Automated CI/CD pipeline with GitHub Actions
- Dependabot integration for dependency updates
- Human-readable output formats (table, text, markdown)
- Logo and improved documentation
- Inspect command for parsing Sentry URLs

### Changed
- Improved README documentation with examples
- Enhanced code formatting and linting

### Technical
- Added comprehensive test suite
- Implemented proper GitHub Actions workflow
- Added release automation configuration

## [0.1.0] - Initial Release

### Added
- Basic CLI structure for Sentry API access
- Core API client with authentication
- Event and issue listing functionality
- Organization and project management commands
- JSON output format