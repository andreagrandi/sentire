package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"sentire/internal/api"
	"sentire/internal/client"
	"sentire/pkg/models"
	"strings"
	"testing"
	"time"
)

func setupTestClient(handler http.HandlerFunc) (*client.Client, *httptest.Server) {
	server := httptest.NewServer(handler)
	
	os.Setenv("SENTRY_API_TOKEN", "test-token")
	c, _ := client.NewClient()
	c.BaseURL = server.URL
	
	return c, server
}

func TestListProjectEvents(t *testing.T) {
	expectedEvents := []models.Event{
		{
			ID:          "event1",
			EventID:     "abc123",
			Title:       "Test Error",
			Platform:    "python",
			DateCreated: time.Now(),
		},
		{
			ID:          "event2",
			EventID:     "def456",
			Title:       "Another Error",
			Platform:    "javascript",
			DateCreated: time.Now(),
		},
	}

	c, server := setupTestClient(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "/projects/test-org/test-project/events/") {
			t.Errorf("Expected path to contain '/projects/test-org/test-project/events/', got %s", r.URL.Path)
		}
		
		// Check query parameters
		if r.URL.Query().Get("statsPeriod") != "24h" {
			t.Errorf("Expected statsPeriod=24h, got %s", r.URL.Query().Get("statsPeriod"))
		}
		
		if r.URL.Query().Get("full") != "true" {
			t.Errorf("Expected full=true, got %s", r.URL.Query().Get("full"))
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expectedEvents)
	})
	defer server.Close()
	defer os.Unsetenv("SENTRY_API_TOKEN")

	eventsAPI := api.NewEventsAPI(c)
	
	opts := &api.ListProjectEventsOptions{
		StatsPeriod: "24h",
		Full:        true,
	}
	
	events, pagination, err := eventsAPI.ListProjectEvents("test-org", "test-project", opts)
	if err != nil {
		t.Fatalf("ListProjectEvents failed: %v", err)
	}
	
	if len(events) != 2 {
		t.Errorf("Expected 2 events, got %d", len(events))
	}
	
	if events[0].ID != "event1" {
		t.Errorf("Expected first event ID 'event1', got %s", events[0].ID)
	}
	
	if pagination == nil {
		t.Error("Expected pagination info, got nil")
	}
}

func TestListIssueEvents(t *testing.T) {
	expectedEvents := []models.Event{
		{
			ID:      "event1",
			EventID: "abc123",
			Title:   "Issue Event",
		},
	}

	c, server := setupTestClient(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "/organizations/test-org/issues/123/events/") {
			t.Errorf("Expected path to contain '/organizations/test-org/issues/123/events/', got %s", r.URL.Path)
		}
		
		environments := r.URL.Query()["environment"]
		if len(environments) != 2 || environments[0] != "prod" || environments[1] != "staging" {
			t.Errorf("Expected environments [prod, staging], got %v", environments)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expectedEvents)
	})
	defer server.Close()
	defer os.Unsetenv("SENTRY_API_TOKEN")

	eventsAPI := api.NewEventsAPI(c)
	
	opts := &api.ListIssueEventsOptions{
		Environment: []string{"prod", "staging"},
		Query:       "test query",
	}
	
	events, _, err := eventsAPI.ListIssueEvents("test-org", "123", opts)
	if err != nil {
		t.Fatalf("ListIssueEvents failed: %v", err)
	}
	
	if len(events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(events))
	}
	
	if events[0].ID != "event1" {
		t.Errorf("Expected event ID 'event1', got %s", events[0].ID)
	}
}

func TestListIssues(t *testing.T) {
	expectedIssues := []models.Issue{
		{
			ID:      "issue1",
			ShortID: "TEST-1",
			Title:   "Test Issue",
			Level:   "error",
			Status:  "unresolved",
		},
		{
			ID:      "issue2",
			ShortID: "TEST-2",
			Title:   "Another Issue",
			Level:   "warning",
			Status:  "resolved",
		},
	}

	c, server := setupTestClient(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "/organizations/test-org/issues/") {
			t.Errorf("Expected path to contain '/organizations/test-org/issues/', got %s", r.URL.Path)
		}
		
		if r.URL.Query().Get("query") != "is:unresolved" {
			t.Errorf("Expected query 'is:unresolved', got %s", r.URL.Query().Get("query"))
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expectedIssues)
	})
	defer server.Close()
	defer os.Unsetenv("SENTRY_API_TOKEN")

	eventsAPI := api.NewEventsAPI(c)
	
	opts := &api.ListIssuesOptions{
		Query: "is:unresolved",
		Limit: 50,
	}
	
	issues, _, err := eventsAPI.ListIssues("test-org", opts)
	if err != nil {
		t.Fatalf("ListIssues failed: %v", err)
	}
	
	if len(issues) != 2 {
		t.Errorf("Expected 2 issues, got %d", len(issues))
	}
	
	if issues[0].ID != "issue1" {
		t.Errorf("Expected first issue ID 'issue1', got %s", issues[0].ID)
	}
}

func TestGetProjectEvent(t *testing.T) {
	expectedEvent := models.Event{
		ID:          "event123",
		EventID:     "abc123def456",
		Title:       "Test Event Detail",
		Platform:    "python",
		DateCreated: time.Now(),
	}

	c, server := setupTestClient(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/projects/test-org/test-project/events/event123/"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expectedEvent)
	})
	defer server.Close()
	defer os.Unsetenv("SENTRY_API_TOKEN")

	eventsAPI := api.NewEventsAPI(c)
	
	event, err := eventsAPI.GetProjectEvent("test-org", "test-project", "event123")
	if err != nil {
		t.Fatalf("GetProjectEvent failed: %v", err)
	}
	
	if event.ID != "event123" {
		t.Errorf("Expected event ID 'event123', got %s", event.ID)
	}
	
	if event.Title != "Test Event Detail" {
		t.Errorf("Expected title 'Test Event Detail', got %s", event.Title)
	}
}

func TestGetIssue(t *testing.T) {
	expectedIssue := models.Issue{
		ID:      "issue123",
		ShortID: "TEST-123",
		Title:   "Test Issue Detail",
		Level:   "error",
		Status:  "unresolved",
	}

	c, server := setupTestClient(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/organizations/test-org/issues/issue123/"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expectedIssue)
	})
	defer server.Close()
	defer os.Unsetenv("SENTRY_API_TOKEN")

	eventsAPI := api.NewEventsAPI(c)
	
	issue, err := eventsAPI.GetIssue("test-org", "issue123")
	if err != nil {
		t.Fatalf("GetIssue failed: %v", err)
	}
	
	if issue.ID != "issue123" {
		t.Errorf("Expected issue ID 'issue123', got %s", issue.ID)
	}
	
	if issue.ShortID != "TEST-123" {
		t.Errorf("Expected short ID 'TEST-123', got %s", issue.ShortID)
	}
}

func TestGetIssueEvent(t *testing.T) {
	expectedEvent := models.Event{
		ID:      "event123",
		EventID: "abc123",
		Title:   "Issue Event Detail",
	}

	c, server := setupTestClient(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/organizations/test-org/issues/issue123/events/latest/"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}
		
		environments := r.URL.Query()["environment"]
		if len(environments) != 1 || environments[0] != "production" {
			t.Errorf("Expected environment [production], got %v", environments)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expectedEvent)
	})
	defer server.Close()
	defer os.Unsetenv("SENTRY_API_TOKEN")

	eventsAPI := api.NewEventsAPI(c)
	
	opts := &api.GetIssueEventOptions{
		Environment: []string{"production"},
	}
	
	event, err := eventsAPI.GetIssueEvent("test-org", "issue123", "latest", opts)
	if err != nil {
		t.Fatalf("GetIssueEvent failed: %v", err)
	}
	
	if event.ID != "event123" {
		t.Errorf("Expected event ID 'event123', got %s", event.ID)
	}
}