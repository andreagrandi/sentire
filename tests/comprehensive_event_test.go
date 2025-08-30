package tests

import (
	"encoding/json"
	"net/http"
	"os"
	"sentire/internal/api"
	"sentire/pkg/models"
	"testing"
	"time"
)

func TestCompleteEventData(t *testing.T) {
	expectedEvent := models.Event{
		ID:           "complete-event-123",
		EventID:      "abc123def456",
		ProjectID:    "project-456",
		GroupID:      "group-789",
		Title:        "Complete Test Event",
		Message:      "Complete error message",
		Platform:     "python",
		Type:         "error",
		DateCreated:  time.Now(),
		DateReceived: time.Now(),
		Size:         12345,
		Culprit:      "test.function",
		Logger:       "test.logger",

		// Stack traces and entries
		Entries: []models.Entry{
			{
				Type: "exception",
				Data: map[string]interface{}{
					"values": []map[string]interface{}{
						{
							"type":  "ValueError",
							"value": "Test exception",
							"stacktrace": map[string]interface{}{
								"frames": []map[string]interface{}{
									{
										"filename":    "test.py",
										"function":    "test_function",
										"lineNo":      42,
										"inApp":       true,
										"contextLine": "raise ValueError('Test exception')",
									},
								},
							},
						},
					},
				},
			},
		},

		// Exception with stack trace
		Exception: &models.Exception{
			Values: []models.ExceptionValue{
				{
					Type:  "ValueError",
					Value: "Test exception value",
					Stacktrace: &models.Stacktrace{
						Frames: []models.StackFrame{
							{
								Filename:    "test.py",
								Function:    "test_function",
								LineNo:      &[]int{42}[0],
								ContextLine: "raise ValueError('Test exception')",
								InApp:       &[]bool{true}[0],
							},
						},
						HasSystemFrames: true,
					},
				},
			},
		},

		// Breadcrumbs
		Breadcrumbs: &models.Breadcrumbs{
			Values: []models.Breadcrumb{
				{
					Timestamp: time.Now(),
					Type:      "navigation",
					Category:  "http",
					Message:   "Test breadcrumb",
					Level:     "info",
				},
			},
		},

		// Request information
		Request: &models.Request{
			URL:    "https://example.com/api/test",
			Method: "POST",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		},

		// Contexts
		Contexts: &models.Contexts{
			Browser: &models.BrowserContext{
				Name:    "Chrome",
				Version: "91.0",
				Type:    "browser",
			},
			OS: &models.OSContext{
				Name:    "Ubuntu",
				Version: "20.04",
				Type:    "os",
			},
			Runtime: &models.RuntimeContext{
				Name:    "CPython",
				Version: "3.9.0",
				Type:    "runtime",
			},
		},

		// Tags
		Tags: []models.EventTag{
			{Key: "environment", Value: "test"},
			{Key: "level", Value: "error"},
		},

		// User
		User: &models.EventUser{
			ID:       "user123",
			Username: "testuser",
			Email:    "test@example.com",
		},

		// Extra data
		Extra: map[string]interface{}{
			"custom_field": "custom_value",
		},

		// Metadata
		Metadata: map[string]interface{}{
			"title": "Test Event",
		},

		// Release info
		Release: &models.EventRelease{
			Version: "1.0.0",
		},

		Environment: "test",

		// SDK
		SDK: &models.EventSDK{
			Name:    "sentry.python",
			Version: "1.0.0",
		},
	}

	c, server := setupTestClient(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expectedEvent)
	})
	defer server.Close()
	defer os.Unsetenv("SENTRY_API_TOKEN")

	eventsAPI := api.NewEventsAPI(c)

	event, err := eventsAPI.GetProjectEvent("test-org", "test-project", "complete-event-123")
	if err != nil {
		t.Fatalf("GetProjectEvent failed: %v", err)
	}

	// Verify basic fields
	if event.ID != "complete-event-123" {
		t.Errorf("Expected event ID 'complete-event-123', got %s", event.ID)
	}

	if event.Culprit != "test.function" {
		t.Errorf("Expected culprit 'test.function', got %s", event.Culprit)
	}

	if event.Size != 12345 {
		t.Errorf("Expected size 12345, got %d", event.Size)
	}

	// Verify entries are captured
	if len(event.Entries) == 0 {
		t.Error("Expected entries to be captured")
	}

	// Verify exception data
	if event.Exception == nil {
		t.Error("Expected exception data to be captured")
	} else if len(event.Exception.Values) == 0 {
		t.Error("Expected exception values to be captured")
	} else {
		exVal := event.Exception.Values[0]
		if exVal.Type != "ValueError" {
			t.Errorf("Expected exception type 'ValueError', got %s", exVal.Type)
		}

		if exVal.Stacktrace == nil {
			t.Error("Expected stacktrace to be captured")
		} else if len(exVal.Stacktrace.Frames) == 0 {
			t.Error("Expected stack frames to be captured")
		}
	}

	// Verify breadcrumbs
	if event.Breadcrumbs == nil {
		t.Error("Expected breadcrumbs to be captured")
	} else if len(event.Breadcrumbs.Values) == 0 {
		t.Error("Expected breadcrumb values to be captured")
	}

	// Verify request data
	if event.Request == nil {
		t.Error("Expected request data to be captured")
	} else {
		if event.Request.URL != "https://example.com/api/test" {
			t.Errorf("Expected request URL 'https://example.com/api/test', got %s", event.Request.URL)
		}
	}

	// Verify contexts
	if event.Contexts == nil {
		t.Error("Expected contexts to be captured")
	} else {
		if event.Contexts.Browser == nil {
			t.Error("Expected browser context to be captured")
		} else if event.Contexts.Browser.Name != "Chrome" {
			t.Errorf("Expected browser name 'Chrome', got %s", event.Contexts.Browser.Name)
		}

		if event.Contexts.Runtime == nil {
			t.Error("Expected runtime context to be captured")
		} else if event.Contexts.Runtime.Name != "CPython" {
			t.Errorf("Expected runtime name 'CPython', got %s", event.Contexts.Runtime.Name)
		}
	}

	// Verify tags
	if len(event.Tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(event.Tags))
	}

	// Verify user data
	if event.User == nil {
		t.Error("Expected user data to be captured")
	} else if event.User.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got %s", event.User.Username)
	}

	// Verify extra data
	if event.Extra == nil {
		t.Error("Expected extra data to be captured")
	} else if event.Extra["custom_field"] != "custom_value" {
		t.Errorf("Expected custom_field 'custom_value', got %v", event.Extra["custom_field"])
	}

	// Verify metadata
	if event.Metadata == nil {
		t.Error("Expected metadata to be captured")
	}

	// Verify release
	if event.Release == nil {
		t.Error("Expected release data to be captured")
	} else if event.Release.Version != "1.0.0" {
		t.Errorf("Expected release version '1.0.0', got %s", event.Release.Version)
	}

	// Verify SDK
	if event.SDK == nil {
		t.Error("Expected SDK data to be captured")
	} else if event.SDK.Name != "sentry.python" {
		t.Errorf("Expected SDK name 'sentry.python', got %s", event.SDK.Name)
	}
}
