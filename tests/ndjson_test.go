package tests

import (
	"bytes"
	"encoding/json"
	"sentire/internal/cli/formatter"
	"sentire/pkg/models"
	"strings"
	"testing"
	"time"
)

func TestNDJSONSingleObject(t *testing.T) {
	issue := &models.Issue{
		ID:    "123",
		Title: "Test Issue",
	}

	var buf bytes.Buffer
	cmd := createTestCommandWithFields("ndjson", "")
	f, err := formatter.NewFormatter(cmd, &buf)
	if err != nil {
		t.Fatalf("Failed to create formatter: %v", err)
	}

	if err := f.FormatIssue(issue); err != nil {
		t.Fatalf("Failed to format issue: %v", err)
	}

	output := strings.TrimSpace(buf.String())
	lines := strings.Split(output, "\n")
	if len(lines) != 1 {
		t.Errorf("Expected 1 line, got %d", len(lines))
	}

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(lines[0]), &result); err != nil {
		t.Fatalf("Line is not valid JSON: %v", err)
	}
	if result["id"] != "123" {
		t.Errorf("Expected id=123, got %v", result["id"])
	}
}

func TestNDJSONMultipleObjects(t *testing.T) {
	events := []models.Event{
		{ID: "e1", Title: "First", DateCreated: time.Now()},
		{ID: "e2", Title: "Second", DateCreated: time.Now()},
		{ID: "e3", Title: "Third", DateCreated: time.Now()},
	}

	var buf bytes.Buffer
	cmd := createTestCommandWithFields("ndjson", "")
	f, err := formatter.NewFormatter(cmd, &buf)
	if err != nil {
		t.Fatalf("Failed to create formatter: %v", err)
	}

	if err := f.FormatEvents(events); err != nil {
		t.Fatalf("Failed to format events: %v", err)
	}

	output := strings.TrimSpace(buf.String())
	lines := strings.Split(output, "\n")
	if len(lines) != 3 {
		t.Fatalf("Expected 3 lines, got %d", len(lines))
	}

	for i, line := range lines {
		var result map[string]interface{}
		if err := json.Unmarshal([]byte(line), &result); err != nil {
			t.Errorf("Line %d is not valid JSON: %v", i, err)
		}
	}

	if strings.Contains(buf.String(), "[") || strings.Contains(buf.String(), "]") {
		t.Error("NDJSON output should not contain array brackets")
	}
}

func TestNDJSONWithFieldsFilter(t *testing.T) {
	events := []models.Event{
		{ID: "e1", EventID: "evt1", Title: "First", DateCreated: time.Now()},
		{ID: "e2", EventID: "evt2", Title: "Second", DateCreated: time.Now()},
	}

	var buf bytes.Buffer
	cmd := createTestCommandWithFields("ndjson", "id,title")
	f, err := formatter.NewFormatter(cmd, &buf)
	if err != nil {
		t.Fatalf("Failed to create formatter: %v", err)
	}

	if err := f.FormatEvents(events); err != nil {
		t.Fatalf("Failed to format events: %v", err)
	}

	output := strings.TrimSpace(buf.String())
	lines := strings.Split(output, "\n")
	if len(lines) != 2 {
		t.Fatalf("Expected 2 lines, got %d", len(lines))
	}

	for i, line := range lines {
		var result map[string]interface{}
		if err := json.Unmarshal([]byte(line), &result); err != nil {
			t.Fatalf("Line %d: invalid JSON: %v", i, err)
		}
		if len(result) != 2 {
			t.Errorf("Line %d: expected 2 fields, got %d: %v", i, len(result), result)
		}
	}
}
