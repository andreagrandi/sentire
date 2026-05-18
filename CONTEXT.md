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

For agent pipelines, prefer `json` (default) or `ndjson`. The `table`, `text`,
and `markdown` formats are intended for humans and may change formatting
between releases — do not parse them.

### Field Filtering

Use `--fields` to limit JSON output to specific fields — reduces token usage:

```bash
sentire events list-issues myorg --fields id,title,status,lastSeen
sentire events get-issue myorg 12345 --fields id,title,count,userCount
```

Field names match the top-level JSON keys returned by Sentry. Run
`sentire describe <command>` to see the supported `output_fields` for each
command before composing a `--fields` list — unknown names are silently
dropped, so check first.

### Pagination

Use `--all` to fetch all pages:

```bash
sentire events list-issues myorg --all
```

`--all` follows cursor-based pagination until exhausted, which can return
large amounts of data. For exploration, omit `--all` and use `--limit` plus
`--fields` first; only enable `--all` when you genuinely need every page.

### Query Filters

The `--query` flag on `events list-issues` and `events list-issue` accepts
Sentry's search syntax. Common, safe building blocks:

- `is:unresolved`, `is:resolved`, `is:ignored`
- `issue.priority:[high,medium]`
- `level:error`, `level:warning`
- `environment:production`
- `firstSeen:-24h`, `lastSeen:-7d`
- `assigned:me`, `assigned:#team-slug`

Combine terms with spaces (implicit AND). Always quote the full value when
it contains spaces or brackets so the shell does not split it.

## Schema Introspection

Use `describe` to discover commands and output fields as JSON:

```bash
# List all commands with args, flags, and output fields
sentire describe

# Describe a specific command
sentire describe events list-issues
```

The `describe` output schema is stable — each command entry exposes `name`,
`description`, `args`, `flags`, and `output_fields`. Agents can rely on these
keys when generating subsequent invocations.

## Agent Workflows

The following recipes show end-to-end flows an agent can run without any
human input. Each step uses only stable JSON output and exit codes for
control flow.

### Workflow 1: Triage unresolved issues for an organization

Goal: produce a short, ranked list of issues that deserve attention this
week.

```bash
# 1. Pull the top unresolved high/medium priority issues from the last 7 days,
#    keep only the fields needed for ranking.
sentire events list-issues myorg \
  --query "is:unresolved issue.priority:[high,medium]" \
  --period 7d \
  --fields id,shortId,title,level,count,userCount,lastSeen,permalink \
  --format ndjson

# 2. For any issue worth a closer look, fetch the full issue record.
sentire events get-issue myorg <issue-id> \
  --fields id,shortId,title,culprit,metadata,firstSeen,lastSeen,count,userCount,status,assignedTo
```

Pick candidates from step 1 by sorting on `userCount` then `count`. Step 2
returns assignment and timing context that is not present in the list view.

### Workflow 2: Inspect a Sentry URL pasted from Slack or email

Goal: get from a notification link to a complete stack trace in one call.

```bash
# 1. Resolve the URL straight to the recommended event (full debugging payload).
sentire inspect "https://myorg.sentry.io/issues/123456789/" \
  --fields id,eventID,title,culprit,platform,dateCreated,tags,entries,contexts

# 2. If the recommended event is not enough, fall back to the issue's other
#    events ordered by recency.
sentire events list-issue myorg 123456789 \
  --fields id,eventID,dateCreated,message,tags \
  --format ndjson
```

`entries` contains the exception(s), stack frames, and breadcrumbs. `contexts`
contains browser/OS/runtime/device data. Both are large — request them only
when you actually need them.

### Workflow 3: Build a release/environment report

Goal: summarise recent error volume for a specific project and environment.

```bash
# 1. Confirm the project exists and capture its slug + platform.
sentire projects get myorg my-project \
  --fields slug,name,platform,status

# 2. Pull issues affecting production in the last 24h.
sentire events list-issues myorg \
  --query "is:unresolved environment:production" \
  --period 24h \
  --project <project-id> \
  --fields id,shortId,title,level,count,userCount,lastSeen \
  --format ndjson

# 3. Pull organisation-wide usage stats for the same window for context.
sentire org stats myorg --period 24h
```

Use `--format ndjson` so the agent can stream-process one issue per line
without holding the full array in memory.

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

1. Start with `sentire describe` to discover commands and their output schemas.
2. Use `--fields` to request only the fields you need — Sentry events can be very large.
3. Use `--format ndjson` for streaming line-by-line processing.
4. Prefer JSON/NDJSON for parsing; `table`/`text`/`markdown` are for humans.
5. Check exit codes for error classification instead of parsing messages.
6. All output goes to stdout, errors go to stderr.
7. The configured token is automatically redacted from error and verbose output, but tokens passed inside arguments (e.g. embedded in a URL) cannot be redacted — keep them out of arguments.
