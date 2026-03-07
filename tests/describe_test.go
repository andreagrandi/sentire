package tests

import (
	"encoding/json"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestDescribeAllCommands(t *testing.T) {
	binary := buildSentire(t)
	stdout, _, exitCode := runSentire(t, binary, "describe")

	if exitCode != 0 {
		t.Fatalf("Expected exit code 0, got %d", exitCode)
	}

	var result struct {
		Commands []struct {
			Name         string   `json:"name"`
			Description  string   `json:"description"`
			OutputFields []string `json:"output_fields"`
		} `json:"commands"`
	}
	if err := json.Unmarshal([]byte(stdout), &result); err != nil {
		t.Fatalf("Invalid JSON output: %v\nOutput: %s", err, stdout)
	}

	if len(result.Commands) == 0 {
		t.Fatal("Expected at least one command in describe output")
	}

	found := make(map[string]bool)
	for _, cmd := range result.Commands {
		found[cmd.Name] = true
	}

	expectedCommands := []string{
		"events list-issues",
		"events get-issue",
		"events get-event",
		"inspect",
		"projects list",
		"projects get",
		"org list-projects",
		"org stats",
	}

	for _, name := range expectedCommands {
		if !found[name] {
			t.Errorf("Expected command %q in describe output", name)
		}
	}
}

func TestDescribeSpecificCommand(t *testing.T) {
	binary := buildSentire(t)
	stdout, _, exitCode := runSentire(t, binary, "describe", "events", "list-issues")

	if exitCode != 0 {
		t.Fatalf("Expected exit code 0, got %d", exitCode)
	}

	var result struct {
		Name         string   `json:"name"`
		Description  string   `json:"description"`
		OutputFields []string `json:"output_fields"`
		Args         []struct {
			Name     string `json:"name"`
			Required bool   `json:"required"`
		} `json:"args"`
	}
	if err := json.Unmarshal([]byte(stdout), &result); err != nil {
		t.Fatalf("Invalid JSON output: %v\nOutput: %s", err, stdout)
	}

	if result.Name != "events list-issues" {
		t.Errorf("Expected name 'events list-issues', got %q", result.Name)
	}

	if len(result.Args) == 0 {
		t.Error("Expected at least one arg")
	}

	if len(result.OutputFields) == 0 {
		t.Error("Expected output_fields to be populated")
	}

	fieldSet := make(map[string]bool)
	for _, f := range result.OutputFields {
		fieldSet[f] = true
	}
	for _, expected := range []string{"id", "title", "status", "lastSeen"} {
		if !fieldSet[expected] {
			t.Errorf("Expected output field %q", expected)
		}
	}
}

func TestDescribeUnknownCommand(t *testing.T) {
	binary := buildSentire(t)

	cmd := exec.Command(binary, "describe", "nonexistent")
	cmd.Env = []string{"PATH=" + os.Getenv("PATH"), "SENTRY_API_TOKEN=test"}
	var stderr strings.Builder
	cmd.Stderr = &stderr
	err := cmd.Run()

	exitCode := 0
	if exitErr, ok := err.(*exec.ExitError); ok {
		exitCode = exitErr.ExitCode()
	}

	if exitCode != 4 {
		t.Errorf("Expected exit code 4 for unknown command, got %d", exitCode)
	}
}
