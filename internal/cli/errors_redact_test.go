package cli

import (
	"bytes"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/andreagrandi/sentire/internal/redact"
)

// TestWriteErrorOutputRedactsToken verifies the CLI's final error writer
// strips the SENTRY_API_TOKEN value regardless of which upstream layer
// produced the error.
func TestWriteErrorOutputRedactsToken(t *testing.T) {
	const token = "cli-secret-token"

	os.Setenv("SENTRY_API_TOKEN", token)
	defer os.Unsetenv("SENTRY_API_TOKEN")

	cases := []struct {
		name   string
		format string
		err    error
	}{
		{
			name:   "json CLIError",
			format: "json",
			err:    NewAPIError("API request failed: token leaked " + token),
		},
		{
			name:   "ndjson CLIError",
			format: "ndjson",
			err:    NewAPIError("API request failed: token leaked " + token),
		},
		{
			name:   "text plain error",
			format: "text",
			err:    errors.New("unexpected failure containing " + token),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			writeErrorOutput(&buf, tc.err, tc.format)

			out := buf.String()
			if strings.Contains(out, token) {
				t.Errorf("Token leaked in %s output: %s", tc.format, out)
			}
			if !strings.Contains(out, redact.Placeholder) {
				t.Errorf("Expected placeholder %q in %s output, got: %s", redact.Placeholder, tc.format, out)
			}
		})
	}
}
