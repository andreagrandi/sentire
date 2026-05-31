# Sentire Workflows and Troubleshooting

Practical recipes for using Sentire against the Sentry API, plus the
short list of problems that account for most failed first runs.

The README covers installation, configuration, and a flag reference.
[`CONTEXT.md`](../CONTEXT.md) is the same material aimed at AI agents
(stable JSON, no narration). This guide is for humans who want a
copy-paste starting point.

Every example assumes `SENTRY_API_TOKEN` is set. Replace `myorg`,
`my-project`, and issue/event IDs with your own values.

## Table of Contents

- [Triage recipes](#triage-recipes)
- [Inspection recipes](#inspection-recipes)
- [Reporting recipes](#reporting-recipes)
- [Output format examples](#output-format-examples)
- [Field filtering examples](#field-filtering-examples)
- [Troubleshooting](#troubleshooting)

## Triage recipes

### List unresolved issues worth looking at today

`events list-issues` defaults to `is:unresolved
issue.priority:[high,medium]`. For a daily triage queue, narrow it to
the last 24 hours and view as a table:

```bash
sentire events list-issues myorg --period 24h --format table
```

The table shows status, priority, event count, user count, and a
relative `Last Seen` (e.g. `3h ago`) so a row at a glance tells you
how active the issue is.

### Focus on a single environment

Production-only triage during an incident:

```bash
sentire events list-issues myorg \
  --query "is:unresolved environment:production" \
  --period 24h --format table
```

`--environment` can also be passed as a top-level flag, but folding it
into `--query` keeps the entire filter in one place that is easy to
edit.

### Sort by frequency or recency

Sentry supports `date`, `freq`, and `inbox` sort orders:

```bash
# Loudest issues first
sentire events list-issues myorg --sort freq --period 7d --format table

# Newest first (the default)
sentire events list-issues myorg --sort date --period 7d --format table
```

### Drill into one issue

Once a row in the table list catches your eye, fetch the full record:

```bash
sentire events get-issue myorg 123456789 --format text
```

`text` keeps the absolute `lastSeen` timestamp and appends the relative
time (`2026-05-23T08:14:02Z (3h ago)`) so you can correlate against
deploys without leaving the terminal.

## Inspection recipes

### Inspect a URL from Slack or email

The `inspect` command parses a Sentry issue URL and fetches the
recommended event with the full debugging payload (stack trace,
breadcrumbs, contexts):

```bash
sentire inspect "https://myorg.sentry.io/issues/123456789/"
```

For a quick scan in the terminal:

```bash
sentire inspect "https://myorg.sentry.io/issues/123456789/" --format text
```

For pasting into a ticket or runbook:

```bash
sentire inspect "https://myorg.sentry.io/issues/123456789/" --format markdown
```

### Fetch a specific event

When you already have the event ID — for example from a customer
report — call `get-event` directly:

```bash
sentire events get-event myorg my-project abcdef0123456789
```

### Walk the events on an issue

If the recommended event is not the one you need, list every event
attached to the issue (newest first) and pick one by ID:

```bash
sentire events list-issue myorg 123456789 --format table
```

Then:

```bash
sentire events get-issue-event myorg 123456789 <event-id>
```

The `<event-id>` argument also accepts `latest`, `oldest`, or
`recommended`.

## Reporting recipes

### Weekly summary in Markdown

`markdown` output is designed to be pasted directly into release
notes, post-mortems, or a Notion page:

```bash
sentire events list-issues myorg \
  --query "is:unresolved" --period 7d \
  --format markdown > weekly-issues.md
```

The result is a header followed by a Markdown table with status,
priority, event/user counts, last seen, and project.

### Org-wide usage stats

For a quick check on event volume before a release:

```bash
sentire org stats myorg --period 7d --format markdown
```

Filter to one or more projects:

```bash
sentire org stats myorg --period 7d \
  --project 123 --project 456 --format markdown
```

### Project inventory

List every project the token can see, in plain text:

```bash
sentire projects list --format text
```

Or limit to one organization:

```bash
sentire org list-projects myorg --format table
```

### Capture a full snapshot in JSON

`json` is the default and is the right choice for anything that will
feed another tool:

```bash
sentire events list-issues myorg --period 7d --all > issues.json
```

`--all` follows cursor-based pagination until the API returns no more
pages. Omit it during exploration and add `--limit` while you tune
filters.

## Output format examples

The same query rendered in each format. Pick `json` or `ndjson` for
scripting, the other three for humans.

### JSON (default — stable, scriptable)

```bash
sentire events list-issues myorg --period 24h
```

```json
[
  {
    "id": "1",
    "shortId": "SENTIRE-1",
    "title": "TypeError in login component",
    "level": "error",
    "status": "unresolved",
    "priority": "high",
    "count": "45",
    "userCount": 23,
    "lastSeen": "2026-05-23T09:14:02Z"
  }
]
```

### NDJSON (one issue per line — stream-friendly)

```bash
sentire events list-issues myorg --period 24h --format ndjson
```

```text
{"id":"1","shortId":"SENTIRE-1","title":"TypeError in login component",...}
{"id":"2","shortId":"SENTIRE-2","title":"API timeout on user endpoint",...}
```

NDJSON pairs well with `jq -c`:

```bash
sentire events list-issues myorg --period 24h --format ndjson \
  | jq -c 'select(.userCount > 10) | {shortId, title, userCount}'
```

### Table (best for terminal scanning)

```bash
sentire events list-issues myorg --period 24h --format table
```

```text
┌───────────┬──────────────────────────────┬─────────┬────────────┬──────────┬────────┬───────┬───────────┬─────────────┐
│    ID     │            TITLE             │  LEVEL  │   STATUS   │ PRIORITY │ EVENTS │ USERS │ LAST SEEN │   PROJECT   │
├───────────┼──────────────────────────────┼─────────┼────────────┼──────────┼────────┼───────┼───────────┼─────────────┤
│ SENTIRE-1 │ TypeError in login component │ error   │ unresolved │ high     │ 45     │ 23    │ 3h ago    │ web-app     │
│ SENTIRE-2 │ API timeout on user endpoint │ warning │ resolved   │ medium   │ 12     │ 8     │ 2d ago    │ api-service │
└───────────┴──────────────────────────────┴─────────┴────────────┴──────────┴────────┴───────┴───────────┴─────────────┘
```

### Text (compact, line-oriented)

```bash
sentire events list-issues myorg --period 24h --format text
```

```text
Issues (2 total):

1. TypeError in login component
   ID: SENTIRE-1 | Status: unresolved | Level: error | Priority: high
   Events: 45 | Users: 23 | Last seen: 3h ago
   Project: web-app

2. API timeout on user endpoint
   ID: SENTIRE-2 | Status: resolved | Level: warning | Priority: medium
   Events: 12 | Users: 8 | Last seen: 2d ago
   Project: api-service
```

### Markdown (pastes into reports and tickets)

```bash
sentire events list-issues myorg --period 24h --format markdown
```

```markdown
# Issues (2 total)

| ID | Title | Level | Status | Priority | Events | Users | Last Seen | Project |
|----|-------|-------|--------|----------|--------|-------|-----------|---------|
| SENTIRE-1 | TypeError in login component | error | unresolved | high | 45 | 23 | 3h ago | web-app |
| SENTIRE-2 | API timeout on user endpoint | warning | resolved | medium | 12 | 8 | 2d ago | api-service |
```

> The `table`, `text`, and `markdown` formats are tuned for humans and
> may change between releases. Parse `json` or `ndjson` instead.

## Field filtering examples

Sentry's `event` and `issue` payloads can be very large. `--fields`
trims the JSON to just the keys you want:

```bash
# Compact issue list — perfect for piping into jq or a spreadsheet
sentire events list-issues myorg --period 7d \
  --fields id,shortId,title,status,priority,count,userCount,lastSeen

# Event details without the heavy entries/contexts payload
sentire events get-issue-event myorg 123456789 recommended \
  --fields id,eventID,title,platform,dateCreated,tags

# Only fetch entries when you actually want the stack trace
sentire events get-issue-event myorg 123456789 recommended \
  --fields id,eventID,entries,contexts
```

Find the supported keys for any command with:

```bash
sentire describe events list-issues
```

The `output_fields` array in the response is the source of truth.
Unknown names are silently dropped, so check `describe` first when a
field you expect does not show up.

## Troubleshooting

### Authentication

**`auth_missing` (exit code 2) — `SENTRY_API_TOKEN is required`**

The token is not set in the environment and no config file is
present. Fix one of:

```bash
export SENTRY_API_TOKEN=your_token_here
# or
mkdir -p ~/.config/sentire
printf '{"sentry_api_token":"your_token_here"}\n' \
  > ~/.config/sentire/config.json
chmod 600 ~/.config/sentire/config.json
```

The environment variable wins when both are set, so a temporary
`SENTRY_API_TOKEN=… sentire …` override is fine.

**`api_error` with HTTP 401 — token rejected by Sentry**

The token is set but Sentry is refusing it. Common causes:

- The token has been revoked or rotated. Generate a new one in
  Sentry under *Settings → Auth Tokens*.
- You pasted the token with surrounding whitespace or quotes.
  `echo "$SENTRY_API_TOKEN" | wc -c` should match the original
  token length plus 1 (for the trailing newline `echo` adds).
- You are using a user auth token with a self-hosted Sentry but
  forgot to set `SENTRY_API_BASE_URL` — the token is being sent to
  `sentry.io` instead.

**Token leaking into shell history**

Prefix the export with a space when `HISTCONTROL` includes
`ignorespace`, or source the token from a credential manager
(`security`, `pass`, `1Password CLI`, etc.). Sentire redacts the
configured token from its own error and verbose output, but cannot
redact tokens embedded inside arguments — keep them out of URLs and
flag values.

### Permissions

**`api_error` with HTTP 403 — `You do not have permission to perform this action`**

The token authenticated but lacks the scopes the endpoint needs.
Check that the token has at least:

- `org:read` for `org list-projects`, `org stats`, and `events list-issues`
- `project:read` for `projects list` and `projects get`
- `event:read` for `events list-project`, `events get-event`, and
  `events get-issue-event`

Regenerate the token with the missing scopes — Sentry does not allow
adding scopes to an existing token.

**`api_error` with HTTP 404 — `The requested resource does not exist`**

Usually one of the slugs is wrong. Double-check:

- Organization slug (the segment after `https://` in your Sentry URL,
  e.g. `myorg` in `myorg.sentry.io`).
- Project slug (visible in *Settings → Projects*, not the display
  name).
- Issue ID — the numeric ID, not the `SENTIRE-123` short ID.

Run `sentire org list-projects myorg` to confirm the org has the
project under the slug you expect.

### Filters and queries

**`events list-issues` returns `[]`**

An empty array is a valid response — it means the filter matched
nothing. Things to try, in order:

1. Drop `--query` entirely:

   ```bash
   sentire events list-issues myorg --period 7d
   ```

   The default query is `is:unresolved
   issue.priority:[high,medium]`. If issues only have `low`
   priority, the default hides them.

2. Widen the time window:

   ```bash
   sentire events list-issues myorg --period 30d
   ```

3. Remove environment and project filters one at a time.

4. Run the same search inside the Sentry web UI. If the web UI also
   shows zero results, the filter — not Sentire — is the cause.

**`invalid_input` (exit code 4) — bad slug, ID, or URL**

The CLI rejects malformed inputs before calling the API. Typical
mistakes:

- Passing a Sentry web URL where a slug is expected (use the
  organization slug, not the full URL — `inspect` is the only
  command that takes a URL).
- Using a `SENTIRE-123` short ID where a numeric issue ID is
  required (`get-issue`, `list-issue`, `get-issue-event` all take
  the numeric ID).
- Trailing slashes or whitespace in arguments.

**Query with brackets gets split by the shell**

Always quote queries that contain spaces or brackets:

```bash
# Correct
sentire events list-issues myorg \
  --query "is:unresolved issue.priority:[high,medium]"

# Wrong — the shell expands `[high,medium]` as a glob
sentire events list-issues myorg \
  --query is:unresolved issue.priority:[high,medium]
```

### Pagination and large responses

**A single page is not enough**

Add `--all` to follow cursors until exhausted:

```bash
sentire events list-issues myorg --period 7d --all > issues.json
```

`--all` can return a lot of data on busy orgs. Combine with
`--fields` to keep the payload manageable.

**Responses are too large to scroll**

Pipe to `jq` for filtering or to a file for inspection:

```bash
sentire events list-issues myorg --period 7d \
  | jq -r '.[] | "\(.shortId)\t\(.title)"'
```

### Rate limiting

**`api_error` mentioning rate limits**

Sentire reads Sentry's rate-limit headers and surfaces a clear
error when the bucket is exhausted. Wait for the window to reset
(usually 60 seconds) and rerun, or drop `--all` to do less work per
invocation. Use `--verbose` to see the current rate-limit budget
on each successful call.

### Networking and self-hosted Sentry

**`timeout` (exit code 3) — request did not complete in time**

The connection or the API itself is slow. Retry once; if it
persists, check Sentry's status page or your VPN/proxy. The CLI
honours `Ctrl+C` and reports a `canceled` code if you abort
mid-flight.

**Self-hosted Sentry**

Point `SENTRY_API_BASE_URL` at the API root of your instance:

```bash
export SENTRY_API_BASE_URL=https://sentry.example.com/api/0
```

Without this, requests go to `https://sentry.io/api/0` and will
fail with 401/404 depending on whether the token also exists on
sentry.io.

### Output formats

**`invalid_format` (exit code 4) — unsupported output format**

Only `json`, `ndjson`, `table`, `text`, and `markdown` are
supported. Check for typos (`tabel`, `md`, `csv` — none of these
work).

**`--fields` returns less than expected**

`--fields` only applies to JSON and NDJSON output. With `--format
table`, `--format text`, or `--format markdown`, the formatters
choose the columns and `--fields` is ignored. Unknown field names
in the comma-separated list are also silently dropped — run
`sentire describe <command>` and copy keys from the `output_fields`
array.

### Discovery

If you are unsure which command or flag covers a use case, the CLI
documents itself:

```bash
# Top-level commands
sentire --help

# Subcommand flags and arguments
sentire events list-issues --help

# Machine-readable schema for every command
sentire describe
sentire describe events get-issue
```

The agent-oriented version of this material, including the JSON
contract for `describe`, lives in [`CONTEXT.md`](../CONTEXT.md).
