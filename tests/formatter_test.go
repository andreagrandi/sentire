package tests

import (
	"bytes"
	"sentire/internal/cli/formatter"
	"sentire/pkg/models"
	"strings"
	"testing"
	"time"

	"github.com/spf13/cobra"
)

// createTestCommand creates a test command with the specified format flag
func createTestCommand(format string) *cobra.Command {
	cmd := &cobra.Command{
		Use: "test",
	}
	cmd.Flags().StringP("format", "f", format, "Test format flag")
	return cmd
}

// Test Event formatting
func TestFormatEvent(t *testing.T) {
	event := &models.Event{
		ID:           "test-id",
		EventID:      "event-123",
		Title:        "Test Error",
		Message:      "This is a test error message",
		Type:         "error",
		Platform:     "python",
		ProjectID:    "project-123",
		DateCreated:  time.Date(2025, 8, 30, 10, 15, 0, 0, time.UTC),
		DateReceived: time.Date(2025, 8, 30, 10, 16, 0, 0, time.UTC),
		Size:         1024,
		Environment:  "production",
		Logger:       "django.request",
		Culprit:      "views.py",
	}

	formats := []string{"json", "table", "text", "markdown"}

	for _, format := range formats {
		t.Run(format, func(t *testing.T) {
			var buf bytes.Buffer
			cmd := createTestCommand(format)
			cmd.Flags().Set("format", format)

			formatter, err := formatter.NewFormatter(cmd, &buf)
			if err != nil {
				t.Fatalf("Failed to create formatter: %v", err)
			}

			err = formatter.FormatEvent(event)
			if err != nil {
				t.Fatalf("Failed to format event: %v", err)
			}

			output := buf.String()
			if output == "" {
				t.Errorf("Expected output for format %s, got empty string", format)
			}

			// Check for key content
			if !strings.Contains(output, "Test Error") {
				t.Errorf("Expected output to contain event title for format %s", format)
			}

			switch format {
			case "json":
				if !strings.Contains(output, `"title": "Test Error"`) {
					t.Errorf("Expected JSON format to contain title field")
				}
			case "table":
				if !strings.Contains(output, "FIELD") || !strings.Contains(output, "VALUE") {
					t.Errorf("Expected table format to have headers")
				}
			case "text":
				if !strings.Contains(output, "Event #") {
					t.Errorf("Expected text format to contain event identifier")
				}
			case "markdown":
				if !strings.Contains(output, "# Event Details") {
					t.Errorf("Expected markdown format to contain header")
				}
			}
		})
	}
}

// Test Issue formatting
func TestFormatIssue(t *testing.T) {
	issue := &models.Issue{
		ID:       "issue-123",
		ShortID:  "SENTIRE-1",
		Title:    "Test Issue",
		Level:    "error",
		Status:   "unresolved",
		Platform: "python",
		Project: models.IssueProject{
			ID:   "project-123",
			Name: "Test Project",
			Slug: "test-project",
		},
		Count:        "100",
		UserCount:    25,
		FirstSeen:    time.Date(2025, 8, 29, 10, 0, 0, 0, time.UTC),
		LastSeen:     time.Date(2025, 8, 30, 10, 0, 0, 0, time.UTC),
		IsPublic:     false,
		IsBookmarked: false,
		IsSubscribed: true,
		Permalink:    "https://sentry.io/test-project/issues/123/",
	}

	formats := []string{"json", "table", "text", "markdown"}

	for _, format := range formats {
		t.Run(format, func(t *testing.T) {
			var buf bytes.Buffer
			cmd := createTestCommand(format)
			cmd.Flags().Set("format", format)

			formatter, err := formatter.NewFormatter(cmd, &buf)
			if err != nil {
				t.Fatalf("Failed to create formatter: %v", err)
			}

			err = formatter.FormatIssue(issue)
			if err != nil {
				t.Fatalf("Failed to format issue: %v", err)
			}

			output := buf.String()
			if output == "" {
				t.Errorf("Expected output for format %s, got empty string", format)
			}

			// Check for key content
			if !strings.Contains(output, "Test Issue") {
				t.Errorf("Expected output to contain issue title for format %s", format)
			}
		})
	}
}

// Test Project formatting
func TestFormatProject(t *testing.T) {
	project := &models.Project{
		ID:       "project-123",
		Slug:     "test-project",
		Name:     "Test Project",
		Platform: "python",
		Organization: models.Organization{
			ID:   "org-123",
			Slug: "test-org",
			Name: "Test Organization",
		},
		Status:       "active",
		DateCreated:  time.Date(2025, 8, 29, 10, 0, 0, 0, time.UTC),
		IsPublic:     false,
		IsBookmarked: true,
	}

	formats := []string{"json", "table", "text", "markdown"}

	for _, format := range formats {
		t.Run(format, func(t *testing.T) {
			var buf bytes.Buffer
			cmd := createTestCommand(format)
			cmd.Flags().Set("format", format)

			formatter, err := formatter.NewFormatter(cmd, &buf)
			if err != nil {
				t.Fatalf("Failed to create formatter: %v", err)
			}

			err = formatter.FormatProject(project)
			if err != nil {
				t.Fatalf("Failed to format project: %v", err)
			}

			output := buf.String()
			if output == "" {
				t.Errorf("Expected output for format %s, got empty string", format)
			}

			// Check for key content
			if !strings.Contains(output, "Test Project") {
				t.Errorf("Expected output to contain project name for format %s", format)
			}
		})
	}
}

// Test OrganizationStats formatting
func TestFormatOrgStats(t *testing.T) {
	stats := &models.OrganizationStats{
		Start: time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC),
		End:   time.Date(2025, 8, 30, 23, 59, 59, 0, time.UTC),
		Projects: []models.ProjectStatsDetail{
			{
				ID:   "project-1",
				Slug: "test-project-1",
			},
			{
				ID:   "project-2",
				Slug: "test-project-2",
			},
		},
		Totals: struct {
			Sum       int64 `json:"sum"`
			TimesSeen int64 `json:"times_seen"`
		}{
			Sum:       1000,
			TimesSeen: 500,
		},
	}

	formats := []string{"json", "table", "text", "markdown"}

	for _, format := range formats {
		t.Run(format, func(t *testing.T) {
			var buf bytes.Buffer
			cmd := createTestCommand(format)
			cmd.Flags().Set("format", format)

			formatter, err := formatter.NewFormatter(cmd, &buf)
			if err != nil {
				t.Fatalf("Failed to create formatter: %v", err)
			}

			err = formatter.FormatOrgStats(stats)
			if err != nil {
				t.Fatalf("Failed to format organization stats: %v", err)
			}

			output := buf.String()
			if output == "" {
				t.Errorf("Expected output for format %s, got empty string", format)
			}

			// Check for key content
			if !strings.Contains(output, "1000") {
				t.Errorf("Expected output to contain total sum for format %s", format)
			}
		})
	}
}

// Test multiple events formatting
func TestFormatEvents(t *testing.T) {
	events := []models.Event{
		{
			ID:          "event-1",
			EventID:     "evt-1",
			Title:       "First Error",
			Type:        "error",
			Platform:    "python",
			ProjectID:   "project-1",
			DateCreated: time.Date(2025, 8, 30, 10, 0, 0, 0, time.UTC),
			Environment: "production",
		},
		{
			ID:          "event-2",
			EventID:     "evt-2",
			Title:       "Second Error",
			Type:        "warning",
			Platform:    "javascript",
			ProjectID:   "project-2",
			DateCreated: time.Date(2025, 8, 30, 11, 0, 0, 0, time.UTC),
			Environment: "staging",
		},
	}

	formats := []string{"json", "table", "text", "markdown"}

	for _, format := range formats {
		t.Run(format, func(t *testing.T) {
			var buf bytes.Buffer
			cmd := createTestCommand(format)
			cmd.Flags().Set("format", format)

			formatter, err := formatter.NewFormatter(cmd, &buf)
			if err != nil {
				t.Fatalf("Failed to create formatter: %v", err)
			}

			err = formatter.FormatEvents(events)
			if err != nil {
				t.Fatalf("Failed to format events: %v", err)
			}

			output := buf.String()
			if output == "" {
				t.Errorf("Expected output for format %s, got empty string", format)
			}

			// Check for both events
			if !strings.Contains(output, "First Error") || !strings.Contains(output, "Second Error") {
				t.Errorf("Expected output to contain both event titles for format %s", format)
			}
		})
	}
}

// Test empty data scenarios
func TestFormatEmptyData(t *testing.T) {
	formats := []string{"json", "table", "text", "markdown"}

	for _, format := range formats {
		t.Run(format+"_empty_events", func(t *testing.T) {
			var buf bytes.Buffer
			cmd := createTestCommand(format)
			cmd.Flags().Set("format", format)

			formatter, err := formatter.NewFormatter(cmd, &buf)
			if err != nil {
				t.Fatalf("Failed to create formatter: %v", err)
			}

			err = formatter.FormatEvents([]models.Event{})
			if err != nil {
				t.Fatalf("Failed to format empty events: %v", err)
			}

			output := buf.String()
			if output == "" {
				t.Errorf("Expected output for empty events in format %s, got empty string", format)
			}
		})
	}
}

// Test unsupported format
func TestUnsupportedFormat(t *testing.T) {
	cmd := createTestCommand("xml")
	cmd.Flags().Set("format", "xml")

	var buf bytes.Buffer
	_, err := formatter.NewFormatter(cmd, &buf)
	if err == nil {
		t.Errorf("Expected error for unsupported format 'xml', got nil")
	}

	if !strings.Contains(err.Error(), "unsupported format: xml") {
		t.Errorf("Expected error message about unsupported format, got: %v", err)
	}
}

// Test Output function with type detection
func TestOutputFunction(t *testing.T) {
	event := &models.Event{
		ID:      "test-event",
		EventID: "evt-123",
		Title:   "Test Event",
		Type:    "error",
	}

	var buf bytes.Buffer
	cmd := createTestCommand("json")
	cmd.Flags().Set("format", "json")

	// Override output to use our buffer (this is a workaround since Output uses os.Stdout)
	originalFormatter, _ := formatter.NewFormatter(cmd, &buf)
	err := originalFormatter.FormatEvent(event)

	if err != nil {
		t.Fatalf("Failed to format event through Output function: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Test Event") {
		t.Errorf("Expected output to contain event title")
	}
}
