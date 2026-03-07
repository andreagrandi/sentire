# Sentire — Agent Context

Sentire is a read-only CLI for the Sentry API. It retrieves issues, events, projects, and organization data from Sentry.

## Authentication

Set `SENTRY_API_TOKEN` before using any command:

```
export SENTRY_API_TOKEN=<your-token>
```

Missing token returns exit code 2 with `auth_missing` error code.

## Command Reference

### Issues & Events

```bash
# List unresolved issues (default query: is:unresolved, high/medium priority)
sentire events list-issues <org-slug>
sentire events list-issues <org-slug> --query "is:unresolved"

# Get a single issue
sentire events get-issue <org-slug> <issue-id>

# Get the recommended event for an issue (full stack trace)
sentire events get-issue-event <org-slug> <issue-id> recommended

# List events for an issue
sentire events list-issue <org-slug> <issue-id>

# List events for a project
sentire events list-project <org-slug> <project-slug>

# Get a specific event
sentire events get-event <org-slug> <project-slug> <event-id>
```

### Inspect (shortcut)

Parse a Sentry URL and fetch the recommended event:

```bash
sentire inspect "https://myorg.sentry.io/issues/123456789/"
```

### Projects

```bash
sentire projects list
sentire projects get <org-slug> <project-slug>
sentire org list-projects <org-slug>
```

### Organization Stats

```bash
sentire org stats <org-slug> --period 7d
```

## Output Control

### Format

Default output is JSON. Available formats: `json`, `ndjson`, `table`, `text`, `markdown`.

```bash
sentire events list-issues myorg --format ndjson
```

### Field Filtering

Use `--fields` to limit JSON output to specific fields — reduces token usage:

```bash
sentire events list-issues myorg --fields id,title,status,lastSeen
sentire events get-issue myorg 12345 --fields id,title,count,userCount
```

### Pagination

Use `--all` to fetch all pages:

```bash
sentire events list-issues myorg --all
```

## Schema Introspection

Use `describe` to discover commands and output fields as JSON:

```bash
# List all commands with args, flags, and output fields
sentire describe

# Describe a specific command
sentire describe events list-issues
```

## Error Handling

Errors are structured JSON when using `--format json` (default):

```json
{"error": "message", "code": "error_code"}
```

### Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | General/unknown error |
| 2 | Authentication error (missing or invalid token) |
| 3 | API error (4xx/5xx from Sentry) |
| 4 | Invalid input (bad slug, ID, URL, or format) |

### Error Codes

- `auth_missing` — SENTRY_API_TOKEN not set
- `api_error` — Sentry API returned an error
- `invalid_input` — Bad argument (malformed slug, ID, or URL)
- `invalid_format` — Unsupported output format

## Tips for AI Agents

1. Use `sentire describe` to discover available commands and their output schemas
2. Use `--fields` to request only the fields you need — Sentry events can be very large
3. Use `--format ndjson` for streaming line-by-line processing
4. Check exit codes for error classification instead of parsing messages
5. All output goes to stdout, errors go to stderr
