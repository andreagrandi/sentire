package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// sentireBinary is the path to the sentire binary built once for the whole test
// package by TestMain. The CLI integration tests run it as a subprocess so they
// exercise argument parsing, output formatting, and exit codes end to end.
var sentireBinary string

// emptyHomeDir is an empty directory used as HOME for CLI subprocesses so a real
// ~/.config/sentire/config.json on the developer's machine never leaks into a
// test run.
var emptyHomeDir string

// TestMain builds the sentire binary into a temporary directory before running
// the package tests and removes it afterwards.
func TestMain(m *testing.M) {
	dir, err := os.MkdirTemp("", "sentire-cli-tests")
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to create temp dir:", err)
		os.Exit(1)
	}

	sentireBinary = filepath.Join(dir, "sentire")
	emptyHomeDir = filepath.Join(dir, "home")
	if err := os.Mkdir(emptyHomeDir, 0o755); err != nil {
		fmt.Fprintln(os.Stderr, "failed to create home dir:", err)
		os.RemoveAll(dir)
		os.Exit(1)
	}

	build := exec.Command("go", "build", "-o", sentireBinary, "github.com/andreagrandi/sentire/cmd/sentire")
	build.Stderr = os.Stderr
	if err := build.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "failed to build sentire binary:", err)
		os.RemoveAll(dir)
		os.Exit(1)
	}

	code := m.Run()
	os.RemoveAll(dir)
	os.Exit(code)
}

// cliResult captures the observable result of a sentire subprocess run.
type cliResult struct {
	stdout   string
	stderr   string
	exitCode int
}

// testEnv builds a hermetic environment for a sentire subprocess. It drops any
// SENTRY_API_* and HOME values inherited from the caller, points HOME at an
// empty directory, and adds the token and base URL only when non-empty.
func testEnv(baseURL, token string) []string {
	var env []string
	for _, e := range os.Environ() {
		if strings.HasPrefix(e, "SENTRY_API_TOKEN=") ||
			strings.HasPrefix(e, "SENTRY_API_BASE_URL=") ||
			strings.HasPrefix(e, "HOME=") {
			continue
		}
		env = append(env, e)
	}
	env = append(env, "HOME="+emptyHomeDir)
	if token != "" {
		env = append(env, "SENTRY_API_TOKEN="+token)
	}
	if baseURL != "" {
		env = append(env, "SENTRY_API_BASE_URL="+baseURL)
	}
	return env
}

// runCLIEnv runs the sentire binary with the given environment and arguments.
func runCLIEnv(t *testing.T, env []string, args ...string) cliResult {
	t.Helper()

	cmd := exec.Command(sentireBinary, args...)
	cmd.Env = env

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	result := cliResult{}
	if err := cmd.Run(); err != nil {
		exitErr, ok := err.(*exec.ExitError)
		if !ok {
			t.Fatalf("failed to run sentire %v: %v", args, err)
		}
		result.exitCode = exitErr.ExitCode()
	}
	result.stdout = stdout.String()
	result.stderr = stderr.String()
	return result
}

// runCLI runs the sentire binary against a mock server with a valid test token.
func runCLI(t *testing.T, baseURL string, args ...string) cliResult {
	t.Helper()
	return runCLIEnv(t, testEnv(baseURL, "test-token"), args...)
}

// mockServer starts an httptest server with handler and registers its shutdown.
func mockServer(t *testing.T, handler http.HandlerFunc) *httptest.Server {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	return server
}

// writeJSON writes v as a JSON response body, failing the test on encoding errors.
func writeJSON(t *testing.T, w http.ResponseWriter, v interface{}) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		t.Errorf("failed to encode mock response: %v", err)
	}
}
