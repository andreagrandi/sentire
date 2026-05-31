package tests

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/andreagrandi/sentire/internal/client"
	"github.com/andreagrandi/sentire/internal/redact"
)

func TestRedactSecret(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		secret string
		want   string
	}{
		{
			name:   "replaces bare token",
			input:  "auth=abc123",
			secret: "abc123",
			want:   "auth=" + redact.Placeholder,
		},
		{
			name:   "replaces multiple occurrences",
			input:  "abc123 and abc123 again",
			secret: "abc123",
			want:   redact.Placeholder + " and " + redact.Placeholder + " again",
		},
		{
			name:   "redacts Bearer-prefixed token via substring match",
			input:  "Authorization: Bearer abc123",
			secret: "abc123",
			want:   "Authorization: Bearer " + redact.Placeholder,
		},
		{
			name:   "empty secret is a no-op",
			input:  "anything",
			secret: "",
			want:   "anything",
		},
		{
			name:   "no occurrence leaves string untouched",
			input:  "nothing to see",
			secret: "abc123",
			want:   "nothing to see",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := redact.Secret(tc.input, tc.secret)
			if got != tc.want {
				t.Errorf("Secret(%q, %q) = %q, want %q", tc.input, tc.secret, got, tc.want)
			}
		})
	}
}

// TestClientAPIErrorRedactsTokenInBody guards against the case where the
// upstream API echoes the Sentry token in an error response body.
func TestClientAPIErrorRedactsTokenInBody(t *testing.T) {
	const token = "super-secret-token-value"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error": "invalid token: ` + token + `"}`))
	}))
	defer server.Close()

	os.Setenv("SENTRY_API_TOKEN", token)
	defer os.Unsetenv("SENTRY_API_TOKEN")

	c, err := client.NewClient()
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	c.BaseURL = server.URL

	_, err = c.Get(context.Background(), "/test", nil)
	if err == nil {
		t.Fatal("Expected error for 401 response")
	}

	if strings.Contains(err.Error(), token) {
		t.Errorf("Expected token to be redacted in error message, got: %s", err.Error())
	}
	if !strings.Contains(err.Error(), redact.Placeholder) {
		t.Errorf("Expected error message to contain placeholder %q, got: %s", redact.Placeholder, err.Error())
	}
}

// TestClientTransportErrorRedactsToken simulates a transport failure where the
// token has somehow ended up in the URL, and verifies it is stripped from the
// wrapped error message that net/http surfaces.
func TestClientTransportErrorRedactsToken(t *testing.T) {
	const token = "transport-error-token"

	os.Setenv("SENTRY_API_TOKEN", token)
	defer os.Unsetenv("SENTRY_API_TOKEN")

	c, err := client.NewClient()
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	c.BaseURL = "http://invalid.localhost.invalid/" + token

	_, err = c.Get(context.Background(), "/test", nil)
	if err == nil {
		t.Fatal("Expected transport error for unreachable host")
	}

	if strings.Contains(err.Error(), token) {
		t.Errorf("Expected token to be redacted in transport error, got: %s", err.Error())
	}
}
