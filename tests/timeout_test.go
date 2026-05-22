package tests

import (
	"context"
	"errors"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/andreagrandi/sentire/internal/client"
)

// hangingHandler blocks until the client gives up on the request, so a test
// server never holds Close() open longer than the client itself waits.
func hangingHandler(w http.ResponseWriter, r *http.Request) {
	select {
	case <-r.Context().Done():
	case <-time.After(5 * time.Second):
	}
}

func TestNewClientDefaultTimeout(t *testing.T) {
	os.Setenv("SENTRY_API_TOKEN", "test-token")
	defer os.Unsetenv("SENTRY_API_TOKEN")

	c, err := client.NewClient()
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	if c.HTTPClient.Timeout != client.DefaultTimeout {
		t.Errorf("Expected HTTP client timeout %s, got %s", client.DefaultTimeout, c.HTTPClient.Timeout)
	}
}

func TestGetTimeoutViaContextDeadline(t *testing.T) {
	c, server := setupTestClient(hangingHandler)
	defer server.Close()
	defer os.Unsetenv("SENTRY_API_TOKEN")

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	_, err := c.Get(ctx, "/test", nil)
	if err == nil {
		t.Fatal("Expected timeout error, got nil")
	}

	var reqErr *client.RequestError
	if !errors.As(err, &reqErr) {
		t.Fatalf("Expected *client.RequestError, got %T: %v", err, err)
	}
	if !reqErr.Timeout {
		t.Errorf("Expected Timeout to be true, got %+v", reqErr)
	}
	if reqErr.Canceled {
		t.Errorf("Expected Canceled to be false for a deadline, got %+v", reqErr)
	}
}

func TestGetTimeoutViaClientTimeout(t *testing.T) {
	c, server := setupTestClient(hangingHandler)
	defer server.Close()
	defer os.Unsetenv("SENTRY_API_TOKEN")

	// The HTTP client's own timeout must surface as a request timeout too.
	c.HTTPClient.Timeout = 50 * time.Millisecond

	_, err := c.Get(context.Background(), "/test", nil)
	if err == nil {
		t.Fatal("Expected timeout error, got nil")
	}

	var reqErr *client.RequestError
	if !errors.As(err, &reqErr) {
		t.Fatalf("Expected *client.RequestError, got %T: %v", err, err)
	}
	if !reqErr.Timeout {
		t.Errorf("Expected Timeout to be true, got %+v", reqErr)
	}
}

func TestGetCanceledBeforeRequest(t *testing.T) {
	c, server := setupTestClient(hangingHandler)
	defer server.Close()
	defer os.Unsetenv("SENTRY_API_TOKEN")

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := c.Get(ctx, "/test", nil)
	if err == nil {
		t.Fatal("Expected cancellation error, got nil")
	}

	var reqErr *client.RequestError
	if !errors.As(err, &reqErr) {
		t.Fatalf("Expected *client.RequestError, got %T: %v", err, err)
	}
	if !reqErr.Canceled {
		t.Errorf("Expected Canceled to be true, got %+v", reqErr)
	}
	if reqErr.Timeout {
		t.Errorf("Expected Timeout to be false for a cancellation, got %+v", reqErr)
	}
}

func TestGetCanceledMidRequest(t *testing.T) {
	started := make(chan struct{})
	c, server := setupTestClient(func(w http.ResponseWriter, r *http.Request) {
		close(started)
		<-r.Context().Done()
	})
	defer server.Close()
	defer os.Unsetenv("SENTRY_API_TOKEN")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		<-started
		cancel()
	}()

	_, err := c.Get(ctx, "/test", nil)
	if err == nil {
		t.Fatal("Expected cancellation error, got nil")
	}

	var reqErr *client.RequestError
	if !errors.As(err, &reqErr) {
		t.Fatalf("Expected *client.RequestError, got %T: %v", err, err)
	}
	if !reqErr.Canceled {
		t.Errorf("Expected Canceled to be true, got %+v", reqErr)
	}
}

func TestDecodeJSONTimeout(t *testing.T) {
	c, server := setupTestClient(func(w http.ResponseWriter, r *http.Request) {
		// Send headers and a partial body, then stall so the timeout fires
		// while DecodeJSON is streaming the body.
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
		w.Write([]byte(`{"partial":`))
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
		select {
		case <-r.Context().Done():
		case <-time.After(5 * time.Second):
		}
	})
	defer server.Close()
	defer os.Unsetenv("SENTRY_API_TOKEN")

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	resp, err := c.Get(ctx, "/test", nil)
	if err != nil {
		t.Fatalf("Get failed before body read: %v", err)
	}

	var out map[string]interface{}
	err = c.DecodeJSON(resp, &out)
	if err == nil {
		t.Fatal("Expected timeout error while decoding body, got nil")
	}

	var reqErr *client.RequestError
	if !errors.As(err, &reqErr) {
		t.Fatalf("Expected *client.RequestError, got %T: %v", err, err)
	}
	if !reqErr.Timeout {
		t.Errorf("Expected Timeout to be true, got %+v", reqErr)
	}
}
