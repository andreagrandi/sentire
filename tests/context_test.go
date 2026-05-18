package tests

import (
	"strings"
	"testing"
)

func TestContextCommandPrintsAgentGuide(t *testing.T) {
	binary := buildSentire(t)
	stdout, stderr, exitCode := runSentire(t, binary, "context")

	if exitCode != 0 {
		t.Fatalf("Expected exit code 0, got %d\nstderr: %s", exitCode, stderr)
	}

	required := []string{
		"# Sentire — Agent Context",
		"## Agent Workflows",
		"### Workflow 1: Triage unresolved issues for an organization",
		"### Workflow 2: Inspect a Sentry URL pasted from Slack or email",
		"### Workflow 3: Build a release/environment report",
		"## Schema Introspection",
		"## Error Handling",
		"### Query Filters",
	}
	for _, marker := range required {
		if !strings.Contains(stdout, marker) {
			t.Errorf("Expected context output to contain %q", marker)
		}
	}
}
