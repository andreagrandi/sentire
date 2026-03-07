package tests

import (
	"strings"
	"testing"
)

func TestValidateOrgSlug(t *testing.T) {
	binary := buildSentire(t)

	tests := []struct {
		name         string
		args         []string
		wantExitCode int
		wantStderr   string
	}{
		{
			name:         "valid org slug",
			args:         []string{"events", "list-issues", "my-org"},
			wantExitCode: 3, // API error since we're hitting a fake token
		},
		{
			name:         "invalid org slug with uppercase",
			args:         []string{"events", "list-issues", "MyOrg"},
			wantExitCode: 4,
			wantStderr:   "invalid_input",
		},
		{
			name:         "invalid org slug with special chars",
			args:         []string{"events", "list-issues", "my@org!"},
			wantExitCode: 4,
			wantStderr:   "invalid_input",
		},
		{
			name:         "invalid issue ID non-numeric",
			args:         []string{"events", "get-issue", "my-org", "abc"},
			wantExitCode: 4,
			wantStderr:   "invalid_input",
		},
		{
			name:         "valid issue ID",
			args:         []string{"events", "get-issue", "my-org", "123456"},
			wantExitCode: 3, // API error
		},
		{
			name:         "invalid event ID",
			args:         []string{"events", "get-event", "my-org", "my-project", "not-hex"},
			wantExitCode: 4,
			wantStderr:   "invalid_input",
		},
		{
			name:         "valid special event ID",
			args:         []string{"events", "get-issue-event", "my-org", "12345", "recommended"},
			wantExitCode: 3, // API error
		},
		{
			name:         "inspect with non-sentry URL",
			args:         []string{"inspect", "https://example.com/issues/123"},
			wantExitCode: 4,
			wantStderr:   "invalid_input",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, stderr, exitCode := runSentire(t, binary, tt.args...)
			if exitCode != tt.wantExitCode {
				t.Errorf("exit code = %d, want %d\nstderr: %s", exitCode, tt.wantExitCode, stderr)
			}
			if tt.wantStderr != "" && !strings.Contains(stderr, tt.wantStderr) {
				t.Errorf("stderr = %q, want it to contain %q", stderr, tt.wantStderr)
			}
		})
	}
}
