package tests

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/andreagrandi/sentire/internal/cli/formatter"
	"github.com/andreagrandi/sentire/pkg/models"
)

// triageIssues returns a fixed set of issues with recent timestamps so that
// relative-time output is predictable.
func triageIssues() []models.Issue {
	now := time.Now()
	return []models.Issue{
		{
			ID:        "issue-1",
			ShortID:   "SENTIRE-1",
			Title:     "TypeError in login component",
			Level:     "error",
			Status:    "unresolved",
			Priority:  "high",
			Project:   models.IssueProject{Name: "Web App", Slug: "web-app"},
			Count:     "45",
			UserCount: 23,
			FirstSeen: now.Add(-5 * 24 * time.Hour),
			LastSeen:  now.Add(-3 * time.Hour),
		},
		{
			ID:        "issue-2",
			ShortID:   "SENTIRE-2",
			Title:     "API timeout on user endpoint",
			Level:     "warning",
			Status:    "resolved",
			Project:   models.IssueProject{Name: "API Service", Slug: "api-service"},
			Count:     "12",
			UserCount: 8,
			FirstSeen: now.Add(-10 * 24 * time.Hour),
			LastSeen:  now.Add(-2 * 24 * time.Hour),
		},
	}
}

func renderIssues(t *testing.T, format string, issues []models.Issue) string {
	t.Helper()
	var buf bytes.Buffer
	cmd := createTestCommand(format)
	cmd.Flags().Set("format", format)

	f, err := formatter.NewFormatter(cmd, &buf)
	if err != nil {
		t.Fatalf("Failed to create %s formatter: %v", format, err)
	}
	if err := f.FormatIssues(issues); err != nil {
		t.Fatalf("Failed to format issues as %s: %v", format, err)
	}
	return buf.String()
}

// Issue lists must surface priority and relative last-seen times for triage.
func TestFormatIssuesTriageFields(t *testing.T) {
	issues := triageIssues()

	for _, format := range []string{"table", "text", "markdown"} {
		t.Run(format, func(t *testing.T) {
			output := renderIssues(t, format, issues)

			if !strings.Contains(strings.ToLower(output), "priority") {
				t.Errorf("%s issue list missing priority column", format)
			}
			if !strings.Contains(output, "high") {
				t.Errorf("%s issue list missing priority value", format)
			}
			if !strings.Contains(output, "3h ago") {
				t.Errorf("%s issue list missing relative last-seen time", format)
			}
			if !strings.Contains(output, "2d ago") {
				t.Errorf("%s issue list missing relative last-seen for second issue", format)
			}
		})
	}
}

// A missing priority should render as a dash in the markdown list rather than
// an empty, misaligned cell.
func TestFormatIssuesEmptyPriorityRendersDash(t *testing.T) {
	output := renderIssues(t, "markdown", triageIssues())
	if !strings.Contains(output, "| - |") {
		t.Errorf("markdown issue list should render '| - |' for missing priority")
	}
}

// Single-issue views keep absolute timestamps but add a relative suffix.
func TestFormatIssueTimestampHasRelativeSuffix(t *testing.T) {
	issue := triageIssues()[0]

	for _, format := range []string{"table", "text", "markdown"} {
		t.Run(format, func(t *testing.T) {
			var buf bytes.Buffer
			cmd := createTestCommand(format)
			cmd.Flags().Set("format", format)

			f, err := formatter.NewFormatter(cmd, &buf)
			if err != nil {
				t.Fatalf("Failed to create %s formatter: %v", format, err)
			}
			if err := f.FormatIssue(&issue); err != nil {
				t.Fatalf("Failed to format issue as %s: %v", format, err)
			}

			output := buf.String()
			if !strings.Contains(output, "(3h ago)") {
				t.Errorf("%s issue view missing relative last-seen suffix", format)
			}
		})
	}
}

// Event lists surface a relative creation time for quick scanning.
func TestFormatEventsRelativeCreated(t *testing.T) {
	events := []models.Event{
		{
			EventID:     "evt-1",
			Title:       "First Error",
			Type:        "error",
			Platform:    "python",
			ProjectID:   "project-1",
			Environment: "production",
			DateCreated: time.Now().Add(-4 * time.Hour),
		},
	}

	for _, format := range []string{"table", "text", "markdown"} {
		t.Run(format, func(t *testing.T) {
			var buf bytes.Buffer
			cmd := createTestCommand(format)
			cmd.Flags().Set("format", format)

			f, err := formatter.NewFormatter(cmd, &buf)
			if err != nil {
				t.Fatalf("Failed to create %s formatter: %v", format, err)
			}
			if err := f.FormatEvents(events); err != nil {
				t.Fatalf("Failed to format events as %s: %v", format, err)
			}

			output := buf.String()
			if !strings.Contains(output, "4h ago") {
				t.Errorf("%s event list missing relative creation time", format)
			}
		})
	}
}
