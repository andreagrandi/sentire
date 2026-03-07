package tests

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"sentire/internal/cli/formatter"
	"sentire/internal/client"
	"sentire/internal/config"
	"strings"
	"testing"
)

func TestCLIErrorJSON(t *testing.T) {
	type CLIError struct {
		Message  string `json:"error"`
		Code     string `json:"code"`
		ExitCode int    `json:"-"`
	}

	err := CLIError{
		Message:  "SENTRY_API_TOKEN not set",
		Code:     "auth_missing",
		ExitCode: 2,
	}

	b, marshalErr := json.Marshal(err)
	if marshalErr != nil {
		t.Fatalf("Failed to marshal CLIError: %v", marshalErr)
	}

	var result map[string]interface{}
	json.Unmarshal(b, &result)

	if result["error"] != "SENTRY_API_TOKEN not set" {
		t.Errorf("Expected error field, got %v", result["error"])
	}
	if result["code"] != "auth_missing" {
		t.Errorf("Expected code field, got %v", result["code"])
	}
	if _, ok := result["ExitCode"]; ok {
		t.Error("ExitCode should not be in JSON output")
	}
}

func TestAuthErrorType(t *testing.T) {
	err := &config.AuthError{Message: "token missing"}
	if err.Error() != "token missing" {
		t.Errorf("Expected 'token missing', got %q", err.Error())
	}
}

func TestAPIErrorType(t *testing.T) {
	err := &client.APIError{Message: "not found", StatusCode: 404}
	if err.Error() != "not found" {
		t.Errorf("Expected 'not found', got %q", err.Error())
	}
	if err.StatusCode != 404 {
		t.Errorf("Expected status 404, got %d", err.StatusCode)
	}
}

func TestFormatErrorType(t *testing.T) {
	err := &formatter.FormatError{Message: "unsupported format: xml"}
	if err.Error() != "unsupported format: xml" {
		t.Errorf("Expected error message, got %q", err.Error())
	}
}

func TestUnsupportedFormatReturnsFormatError(t *testing.T) {
	cmd := createTestCommand("xml")
	cmd.Flags().Set("format", "xml")

	var buf bytes.Buffer
	_, err := formatter.NewFormatter(cmd, &buf)
	if err == nil {
		t.Fatal("Expected error for unsupported format")
	}

	_, ok := err.(*formatter.FormatError)
	if !ok {
		t.Errorf("Expected *formatter.FormatError, got %T", err)
	}
	if !strings.Contains(err.Error(), "xml") {
		t.Errorf("Expected error to mention 'xml', got %q", err.Error())
	}
}

func buildSentire(t *testing.T) string {
	t.Helper()
	binary := t.TempDir() + "/sentire"
	cmd := exec.Command("go", "build", "-o", binary, "./cmd/sentire")
	cmd.Dir = ".."
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to build sentire: %v\n%s", err, out)
	}
	return binary
}

func runSentire(t *testing.T, binary string, args ...string) (string, string, int) {
	t.Helper()
	cmd := exec.Command(binary, args...)
	cmd.Env = append(os.Environ(), "SENTRY_API_TOKEN=test-token")
	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	exitCode := 0
	if exitErr, ok := err.(*exec.ExitError); ok {
		exitCode = exitErr.ExitCode()
	} else if err != nil {
		t.Fatalf("Failed to run sentire: %v", err)
	}
	return stdout.String(), stderr.String(), exitCode
}

func TestExitCodeAuthError(t *testing.T) {
	binary := buildSentire(t)

	cmd := exec.Command(binary, "events", "list-issues", "my-org")
	cmd.Env = []string{"PATH=" + os.Getenv("PATH"), "HOME=/tmp/nonexistent"}
	var stderr strings.Builder
	cmd.Stderr = &stderr
	err := cmd.Run()

	exitCode := 0
	if exitErr, ok := err.(*exec.ExitError); ok {
		exitCode = exitErr.ExitCode()
	}

	if exitCode != 2 {
		t.Errorf("Expected exit code 2 for auth error, got %d\nstderr: %s", exitCode, stderr.String())
	}
	if !strings.Contains(stderr.String(), "auth_missing") {
		t.Errorf("Expected stderr to contain 'auth_missing', got %q", stderr.String())
	}
}

func TestExitCodeAPIError(t *testing.T) {
	binary := buildSentire(t)

	// Use a valid slug so validation passes, then the API call will fail
	_, stderr, exitCode := runSentire(t, binary, "events", "list-issues", "my-org")

	if exitCode != 3 {
		t.Errorf("Expected exit code 3 for API error, got %d\nstderr: %s", exitCode, stderr)
	}
}
