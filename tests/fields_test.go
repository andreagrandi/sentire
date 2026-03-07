package tests

import (
	"bytes"
	"encoding/json"
	"sentire/internal/cli/formatter"
	"sentire/pkg/models"
	"testing"
	"time"

	"github.com/spf13/cobra"
)

func createTestCommandWithFields(format, fields string) *cobra.Command {
	cmd := &cobra.Command{Use: "test"}
	cmd.Flags().StringP("format", "f", format, "Test format flag")
	cmd.Flags().String("fields", fields, "Test fields flag")
	return cmd
}

func TestFieldsFilteringSingleObject(t *testing.T) {
	issue := &models.Issue{
		ID:      "123",
		ShortID: "PROJ-1",
		Title:   "Test Issue",
		Status:  "unresolved",
		Level:   "error",
	}

	var buf bytes.Buffer
	cmd := createTestCommandWithFields("json", "id,title,status")
	f, err := formatter.NewFormatter(cmd, &buf)
	if err != nil {
		t.Fatalf("Failed to create formatter: %v", err)
	}

	if err := f.FormatIssue(issue); err != nil {
		t.Fatalf("Failed to format issue: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("Failed to parse output: %v", err)
	}

	if len(result) != 3 {
		t.Errorf("Expected 3 fields, got %d: %v", len(result), result)
	}
	if result["id"] != "123" {
		t.Errorf("Expected id=123, got %v", result["id"])
	}
	if result["title"] != "Test Issue" {
		t.Errorf("Expected title=Test Issue, got %v", result["title"])
	}
	if result["status"] != "unresolved" {
		t.Errorf("Expected status=unresolved, got %v", result["status"])
	}
}

func TestFieldsFilteringSlice(t *testing.T) {
	events := []models.Event{
		{ID: "e1", EventID: "evt1", Title: "First", DateCreated: time.Now()},
		{ID: "e2", EventID: "evt2", Title: "Second", DateCreated: time.Now()},
	}

	var buf bytes.Buffer
	cmd := createTestCommandWithFields("json", "id,title")
	f, err := formatter.NewFormatter(cmd, &buf)
	if err != nil {
		t.Fatalf("Failed to create formatter: %v", err)
	}

	if err := f.FormatEvents(events); err != nil {
		t.Fatalf("Failed to format events: %v", err)
	}

	var result []map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("Failed to parse output: %v", err)
	}

	if len(result) != 2 {
		t.Fatalf("Expected 2 items, got %d", len(result))
	}

	for i, item := range result {
		if len(item) != 2 {
			t.Errorf("Item %d: expected 2 fields, got %d: %v", i, len(item), item)
		}
		if _, ok := item["id"]; !ok {
			t.Errorf("Item %d: missing 'id' field", i)
		}
		if _, ok := item["title"]; !ok {
			t.Errorf("Item %d: missing 'title' field", i)
		}
	}
}

func TestNoFieldsFilteringReturnsAll(t *testing.T) {
	issue := &models.Issue{
		ID:      "123",
		ShortID: "PROJ-1",
		Title:   "Test Issue",
		Status:  "unresolved",
		Level:   "error",
	}

	var buf bytes.Buffer
	cmd := createTestCommandWithFields("json", "")
	f, err := formatter.NewFormatter(cmd, &buf)
	if err != nil {
		t.Fatalf("Failed to create formatter: %v", err)
	}

	if err := f.FormatIssue(issue); err != nil {
		t.Fatalf("Failed to format issue: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("Failed to parse output: %v", err)
	}

	if len(result) < 5 {
		t.Errorf("Expected many fields without filter, got %d", len(result))
	}
}
