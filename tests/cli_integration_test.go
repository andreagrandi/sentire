package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/andreagrandi/sentire/pkg/models"
)

// decodeCLIJSON unmarshals stdout produced by a sentire command, failing the
// test with the raw output when the payload is not valid JSON.
func decodeCLIJSON(t *testing.T, output string, v interface{}) {
	t.Helper()
	if err := json.Unmarshal([]byte(output), v); err != nil {
		t.Fatalf("CLI output is not valid JSON: %v\noutput:\n%s", err, output)
	}
}

// sampleIssues returns issues with recent timestamps so relative-time output is
// stable across formats.
func sampleIssues() []models.Issue {
	now := time.Now()
	return []models.Issue{
		{
			ID:        "issue-1",
			ShortID:   "WEB-1",
			Title:     "TypeError in checkout flow",
			Level:     "error",
			Status:    "unresolved",
			Priority:  "high",
			Project:   models.IssueProject{Name: "Web App", Slug: "web-app"},
			Count:     "42",
			UserCount: 17,
			FirstSeen: now.Add(-6 * 24 * time.Hour),
			LastSeen:  now.Add(-2 * time.Hour),
		},
		{
			ID:        "issue-2",
			ShortID:   "API-7",
			Title:     "Timeout calling payments service",
			Level:     "warning",
			Status:    "unresolved",
			Priority:  "medium",
			Project:   models.IssueProject{Name: "API", Slug: "api"},
			Count:     "9",
			UserCount: 4,
			FirstSeen: now.Add(-3 * 24 * time.Hour),
			LastSeen:  now.Add(-30 * time.Minute),
		},
	}
}

// sampleEvent returns an event with debugging details for inspect/get flows.
func sampleEvent() models.Event {
	return models.Event{
		ID:          "event-abc",
		EventID:     "a1b2c3d4e5f60718293a4b5c6d7e8f90",
		Title:       "TypeError in checkout flow",
		Message:     "Cannot read property 'total' of undefined",
		Type:        "error",
		Platform:    "javascript",
		ProjectID:   "project-1",
		Environment: "production",
		Culprit:     "checkout.js",
		DateCreated: time.Now().Add(-2 * time.Hour),
	}
}

// --- Projects -------------------------------------------------------------

func TestCLIProjectsList(t *testing.T) {
	projects := []models.Project{
		{ID: "p1", Slug: "web-app", Name: "Web App", Platform: "javascript"},
		{ID: "p2", Slug: "api", Name: "API Service", Platform: "python"},
	}

	server := mockServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/projects/" {
			t.Errorf("expected path /projects/, got %s", r.URL.Path)
		}
		if auth := r.Header.Get("Authorization"); auth != "Bearer test-token" {
			t.Errorf("expected bearer test-token, got %q", auth)
		}
		writeJSON(t, w, projects)
	})

	t.Run("json", func(t *testing.T) {
		res := runCLI(t, server.URL, "projects", "list")
		if res.exitCode != 0 {
			t.Fatalf("expected exit 0, got %d (stderr: %s)", res.exitCode, res.stderr)
		}
		var got []models.Project
		decodeCLIJSON(t, res.stdout, &got)
		if len(got) != 2 || got[0].Slug != "web-app" || got[1].Slug != "api" {
			t.Errorf("unexpected projects: %+v", got)
		}
	})

	t.Run("table", func(t *testing.T) {
		res := runCLI(t, server.URL, "projects", "list", "--format", "table")
		if res.exitCode != 0 {
			t.Fatalf("expected exit 0, got %d (stderr: %s)", res.exitCode, res.stderr)
		}
		if !strings.Contains(res.stdout, "Web App") || !strings.Contains(res.stdout, "API Service") {
			t.Errorf("table output missing project names:\n%s", res.stdout)
		}
	})
}

func TestCLIProjectsGet(t *testing.T) {
	server := mockServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/projects/test-org/web-app/" {
			t.Errorf("expected path /projects/test-org/web-app/, got %s", r.URL.Path)
		}
		writeJSON(t, w, models.Project{ID: "p1", Slug: "web-app", Name: "Web App", Platform: "javascript"})
	})

	res := runCLI(t, server.URL, "projects", "get", "test-org", "web-app")
	if res.exitCode != 0 {
		t.Fatalf("expected exit 0, got %d (stderr: %s)", res.exitCode, res.stderr)
	}
	var got models.Project
	decodeCLIJSON(t, res.stdout, &got)
	if got.Slug != "web-app" {
		t.Errorf("expected slug web-app, got %s", got.Slug)
	}
}

func TestCLIProjectsGetInvalidSlug(t *testing.T) {
	// Validation rejects the slug before any HTTP request is attempted.
	res := runCLI(t, "http://127.0.0.1:0", "projects", "get", "Test_Org", "web-app")
	if res.exitCode != 4 {
		t.Fatalf("expected exit 4 for invalid slug, got %d (stderr: %s)", res.exitCode, res.stderr)
	}
	assertErrorCode(t, res.stderr, "invalid_input")
}

// --- Organizations --------------------------------------------------------

func TestCLIOrgListProjects(t *testing.T) {
	server := mockServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/organizations/test-org/projects/" {
			t.Errorf("expected org projects path, got %s", r.URL.Path)
		}
		writeJSON(t, w, []models.Project{{ID: "p1", Slug: "web-app", Name: "Web App"}})
	})

	res := runCLI(t, server.URL, "org", "list-projects", "test-org")
	if res.exitCode != 0 {
		t.Fatalf("expected exit 0, got %d (stderr: %s)", res.exitCode, res.stderr)
	}
	var got []models.Project
	decodeCLIJSON(t, res.stdout, &got)
	if len(got) != 1 || got[0].Slug != "web-app" {
		t.Errorf("unexpected projects: %+v", got)
	}
}

func TestCLIOrgStats(t *testing.T) {
	stats := models.OrganizationStats{
		Start:    time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC),
		End:      time.Date(2026, 5, 22, 0, 0, 0, 0, time.UTC),
		Projects: []models.ProjectStatsDetail{{ID: "p1", Slug: "web-app"}},
	}
	stats.Totals.Sum = 1000
	stats.Totals.TimesSeen = 500

	server := mockServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/organizations/test-org/stats-summary/" {
			t.Errorf("expected stats-summary path, got %s", r.URL.Path)
		}
		if field := r.URL.Query().Get("field"); field != "sum(quantity)" {
			t.Errorf("expected default field sum(quantity), got %q", field)
		}
		writeJSON(t, w, stats)
	})

	t.Run("json", func(t *testing.T) {
		res := runCLI(t, server.URL, "org", "stats", "test-org")
		if res.exitCode != 0 {
			t.Fatalf("expected exit 0, got %d (stderr: %s)", res.exitCode, res.stderr)
		}
		var got models.OrganizationStats
		decodeCLIJSON(t, res.stdout, &got)
		if got.Totals.Sum != 1000 {
			t.Errorf("expected total sum 1000, got %d", got.Totals.Sum)
		}
	})

	t.Run("text", func(t *testing.T) {
		res := runCLI(t, server.URL, "org", "stats", "test-org", "--format", "text")
		if res.exitCode != 0 {
			t.Fatalf("expected exit 0, got %d (stderr: %s)", res.exitCode, res.stderr)
		}
		if !strings.Contains(res.stdout, "1000") {
			t.Errorf("text output missing total sum:\n%s", res.stdout)
		}
	})
}

// --- Events & Issues ------------------------------------------------------

func TestCLIEventsListIssues(t *testing.T) {
	issues := sampleIssues()

	server := mockServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/organizations/test-org/issues/" {
			t.Errorf("expected issues path, got %s", r.URL.Path)
		}
		writeJSON(t, w, issues)
	})

	t.Run("json", func(t *testing.T) {
		res := runCLI(t, server.URL, "events", "list-issues", "test-org")
		if res.exitCode != 0 {
			t.Fatalf("expected exit 0, got %d (stderr: %s)", res.exitCode, res.stderr)
		}
		var got []models.Issue
		decodeCLIJSON(t, res.stdout, &got)
		if len(got) != 2 || got[0].ShortID != "WEB-1" {
			t.Errorf("unexpected issues: %+v", got)
		}
	})

	for _, format := range []string{"table", "text", "markdown"} {
		t.Run(format, func(t *testing.T) {
			res := runCLI(t, server.URL, "events", "list-issues", "test-org", "--format", format)
			if res.exitCode != 0 {
				t.Fatalf("expected exit 0, got %d (stderr: %s)", res.exitCode, res.stderr)
			}
			if !strings.Contains(res.stdout, "TypeError in checkout flow") {
				t.Errorf("%s output missing issue title:\n%s", format, res.stdout)
			}
			if !strings.Contains(strings.ToLower(res.stdout), "priority") {
				t.Errorf("%s output missing priority column:\n%s", format, res.stdout)
			}
		})
	}
}

func TestCLIEventsListIssuesPagination(t *testing.T) {
	issues := sampleIssues()

	server := mockServer(t, func(w http.ResponseWriter, r *http.Request) {
		// The client only reads the cursor query param from the Link URL, so a
		// fixed absolute URL is enough to drive the next page.
		switch r.URL.Query().Get("cursor") {
		case "":
			w.Header().Set("Link",
				`<https://sentry.local/issues/?cursor=page2>; rel="next"; results="true"`)
			writeJSON(t, w, issues[:1])
		case "page2":
			writeJSON(t, w, issues[1:])
		default:
			t.Errorf("unexpected cursor %q", r.URL.Query().Get("cursor"))
		}
	})

	res := runCLI(t, server.URL, "events", "list-issues", "test-org", "--all")
	if res.exitCode != 0 {
		t.Fatalf("expected exit 0, got %d (stderr: %s)", res.exitCode, res.stderr)
	}
	var got []models.Issue
	decodeCLIJSON(t, res.stdout, &got)
	if len(got) != 2 {
		t.Errorf("expected 2 issues across pages, got %d", len(got))
	}
}

func TestCLIEventsGetIssue(t *testing.T) {
	server := mockServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/organizations/test-org/issues/123456/" {
			t.Errorf("expected single issue path, got %s", r.URL.Path)
		}
		writeJSON(t, w, models.Issue{ID: "123456", ShortID: "WEB-1", Title: "TypeError in checkout flow"})
	})

	res := runCLI(t, server.URL, "events", "get-issue", "test-org", "123456")
	if res.exitCode != 0 {
		t.Fatalf("expected exit 0, got %d (stderr: %s)", res.exitCode, res.stderr)
	}
	var got models.Issue
	decodeCLIJSON(t, res.stdout, &got)
	if got.ShortID != "WEB-1" {
		t.Errorf("expected short ID WEB-1, got %s", got.ShortID)
	}
}

func TestCLIEventsListProject(t *testing.T) {
	server := mockServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/projects/test-org/web-app/events/" {
			t.Errorf("expected project events path, got %s", r.URL.Path)
		}
		writeJSON(t, w, []models.Event{sampleEvent()})
	})

	res := runCLI(t, server.URL, "events", "list-project", "test-org", "web-app")
	if res.exitCode != 0 {
		t.Fatalf("expected exit 0, got %d (stderr: %s)", res.exitCode, res.stderr)
	}
	var got []models.Event
	decodeCLIJSON(t, res.stdout, &got)
	if len(got) != 1 || got[0].Title != "TypeError in checkout flow" {
		t.Errorf("unexpected events: %+v", got)
	}
}

func TestCLIEventsGetEvent(t *testing.T) {
	eventID := "a1b2c3d4e5f60718293a4b5c6d7e8f90"
	server := mockServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/projects/test-org/web-app/events/"+eventID+"/" {
			t.Errorf("expected event detail path, got %s", r.URL.Path)
		}
		writeJSON(t, w, sampleEvent())
	})

	res := runCLI(t, server.URL, "events", "get-event", "test-org", "web-app", eventID)
	if res.exitCode != 0 {
		t.Fatalf("expected exit 0, got %d (stderr: %s)", res.exitCode, res.stderr)
	}
	var got models.Event
	decodeCLIJSON(t, res.stdout, &got)
	if got.EventID != eventID {
		t.Errorf("expected event ID %s, got %s", eventID, got.EventID)
	}
}

func TestCLIEventsGetIssueEvent(t *testing.T) {
	server := mockServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/organizations/test-org/issues/123456/events/latest/" {
			t.Errorf("expected issue event path, got %s", r.URL.Path)
		}
		writeJSON(t, w, sampleEvent())
	})

	res := runCLI(t, server.URL, "events", "get-issue-event", "test-org", "123456", "latest")
	if res.exitCode != 0 {
		t.Fatalf("expected exit 0, got %d (stderr: %s)", res.exitCode, res.stderr)
	}
	var got models.Event
	decodeCLIJSON(t, res.stdout, &got)
	if got.Title != "TypeError in checkout flow" {
		t.Errorf("unexpected event title: %s", got.Title)
	}
}

// --- Inspect --------------------------------------------------------------

func TestCLIInspect(t *testing.T) {
	server := mockServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/organizations/laterpay/issues/6796439331/events/recommended/" {
			t.Errorf("expected recommended event path, got %s", r.URL.Path)
		}
		writeJSON(t, w, sampleEvent())
	})

	url := "https://laterpay.sentry.io/issues/6796439331/?alert_type=issue"

	t.Run("json", func(t *testing.T) {
		res := runCLI(t, server.URL, "inspect", url)
		if res.exitCode != 0 {
			t.Fatalf("expected exit 0, got %d (stderr: %s)", res.exitCode, res.stderr)
		}
		var got models.Event
		decodeCLIJSON(t, res.stdout, &got)
		if got.Title != "TypeError in checkout flow" {
			t.Errorf("unexpected event title: %s", got.Title)
		}
	})

	t.Run("markdown", func(t *testing.T) {
		res := runCLI(t, server.URL, "inspect", url, "--format", "markdown")
		if res.exitCode != 0 {
			t.Fatalf("expected exit 0, got %d (stderr: %s)", res.exitCode, res.stderr)
		}
		if !strings.Contains(res.stdout, "# Event Details") {
			t.Errorf("markdown output missing event header:\n%s", res.stdout)
		}
	})
}

func TestCLIInspectInvalidURL(t *testing.T) {
	res := runCLI(t, "http://127.0.0.1:0", "inspect", "https://example.com/not-sentry")
	if res.exitCode != 4 {
		t.Fatalf("expected exit 4 for invalid URL, got %d (stderr: %s)", res.exitCode, res.stderr)
	}
	assertErrorCode(t, res.stderr, "invalid_input")
}

// --- Agent-facing commands ------------------------------------------------

func TestCLIContext(t *testing.T) {
	// context is self-contained and needs no API token.
	res := runCLIEnv(t, testEnv("", ""), "context")
	if res.exitCode != 0 {
		t.Fatalf("expected exit 0, got %d (stderr: %s)", res.exitCode, res.stderr)
	}
	if !strings.Contains(res.stdout, "# Sentire — Agent Context") {
		t.Errorf("context output missing guide heading:\n%s", res.stdout)
	}
}

func TestCLIDescribe(t *testing.T) {
	res := runCLIEnv(t, testEnv("", ""), "describe")
	if res.exitCode != 0 {
		t.Fatalf("expected exit 0, got %d (stderr: %s)", res.exitCode, res.stderr)
	}

	var got struct {
		Commands []struct {
			Name string `json:"name"`
		} `json:"commands"`
	}
	decodeCLIJSON(t, res.stdout, &got)

	found := false
	for _, c := range got.Commands {
		if c.Name == "events list-issues" {
			found = true
		}
	}
	if !found {
		t.Errorf("describe output missing 'events list-issues' command")
	}
}

func TestCLIDescribeCommand(t *testing.T) {
	res := runCLIEnv(t, testEnv("", ""), "describe", "events", "list-issues")
	if res.exitCode != 0 {
		t.Fatalf("expected exit 0, got %d (stderr: %s)", res.exitCode, res.stderr)
	}

	var got struct {
		Name         string   `json:"name"`
		OutputFields []string `json:"output_fields"`
	}
	decodeCLIJSON(t, res.stdout, &got)
	if got.Name != "events list-issues" {
		t.Errorf("expected command name 'events list-issues', got %q", got.Name)
	}
	if len(got.OutputFields) == 0 {
		t.Errorf("expected non-empty output_fields for events list-issues")
	}
}

// --- Error handling -------------------------------------------------------

func TestCLIAPIErrorJSON(t *testing.T) {
	server := mockServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, `{"detail":"internal server error"}`)
	})

	res := runCLI(t, server.URL, "events", "list-issues", "test-org")
	if res.exitCode != 3 {
		t.Fatalf("expected exit 3 for API error, got %d (stderr: %s)", res.exitCode, res.stderr)
	}
	assertErrorCode(t, res.stderr, "api_error")
}

func TestCLIAPIErrorHumanReadable(t *testing.T) {
	server := mockServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, `{"detail":"not found"}`)
	})

	res := runCLI(t, server.URL, "events", "list-issues", "test-org", "--format", "table")
	if res.exitCode != 3 {
		t.Fatalf("expected exit 3 for API error, got %d (stderr: %s)", res.exitCode, res.stderr)
	}
	if !strings.HasPrefix(res.stderr, "Error:") {
		t.Errorf("expected human-readable error prefix, got: %s", res.stderr)
	}
}

func TestCLIMissingToken(t *testing.T) {
	// No token and an empty HOME means no config file is found either.
	res := runCLIEnv(t, testEnv("", ""), "events", "list-issues", "test-org")
	if res.exitCode != 2 {
		t.Fatalf("expected exit 2 for missing token, got %d (stderr: %s)", res.exitCode, res.stderr)
	}
	assertErrorCode(t, res.stderr, "auth_missing")
}

func TestCLIUnsupportedFormat(t *testing.T) {
	server := mockServer(t, func(w http.ResponseWriter, r *http.Request) {
		writeJSON(t, w, []models.Issue{})
	})

	res := runCLI(t, server.URL, "events", "list-issues", "test-org", "--format", "xml")
	if res.exitCode != 4 {
		t.Fatalf("expected exit 4 for unsupported format, got %d (stderr: %s)", res.exitCode, res.stderr)
	}
	if !strings.Contains(res.stderr, "unsupported format: xml") {
		t.Errorf("expected unsupported format message, got: %s", res.stderr)
	}
}

// assertErrorCode checks that stderr carries a structured JSON error with the
// expected machine-readable code.
func assertErrorCode(t *testing.T, stderr, wantCode string) {
	t.Helper()
	var got struct {
		Error string `json:"error"`
		Code  string `json:"code"`
	}
	if err := json.Unmarshal([]byte(stderr), &got); err != nil {
		t.Fatalf("error output is not JSON: %v\nstderr:\n%s", err, stderr)
	}
	if got.Code != wantCode {
		t.Errorf("expected error code %q, got %q (stderr: %s)", wantCode, got.Code, stderr)
	}
}
