package tests

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"sentire/internal/client"
	"strings"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	// Test with missing token
	os.Unsetenv("SENTRY_API_TOKEN")
	_, err := client.NewClient()
	if err == nil {
		t.Error("Expected error when SENTRY_API_TOKEN is not set")
	}

	// Test with valid token
	os.Setenv("SENTRY_API_TOKEN", "test-token")
	defer os.Unsetenv("SENTRY_API_TOKEN")
	
	c, err := client.NewClient()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	if c.Token != "test-token" {
		t.Errorf("Expected token 'test-token', got %s", c.Token)
	}
	
	if c.BaseURL != client.BaseURL {
		t.Errorf("Expected base URL %s, got %s", client.BaseURL, c.BaseURL)
	}
}

func TestClientDo(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check headers
		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-token" {
			t.Errorf("Expected Authorization header 'Bearer test-token', got %s", auth)
		}
		
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("Expected Content-Type 'application/json', got %s", contentType)
		}
		
		userAgent := r.Header.Get("User-Agent")
		if userAgent != client.UserAgent {
			t.Errorf("Expected User-Agent %s, got %s", client.UserAgent, userAgent)
		}
		
		// Set rate limit headers
		w.Header().Set("X-Sentry-Rate-Limit-Limit", "100")
		w.Header().Set("X-Sentry-Rate-Limit-Remaining", "99")
		w.Header().Set("X-Sentry-Rate-Limit-Reset", "1234567890")
		
		// Set Link header for pagination
		w.Header().Set("Link", `<https://sentry.io/api/0/test?cursor=next123>; rel="next"; results="true"`)
		
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"test": "data"}`))
	}))
	defer server.Close()

	os.Setenv("SENTRY_API_TOKEN", "test-token")
	defer os.Unsetenv("SENTRY_API_TOKEN")
	
	c, err := client.NewClient()
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	
	// Create request
	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	
	resp, err := c.Do(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()
	
	// Check rate limit parsing
	if c.RateLimit.Limit != 100 {
		t.Errorf("Expected rate limit 100, got %d", c.RateLimit.Limit)
	}
	
	if c.RateLimit.Remaining != 99 {
		t.Errorf("Expected rate limit remaining 99, got %d", c.RateLimit.Remaining)
	}
	
	expectedTime := time.Unix(1234567890, 0)
	if !c.RateLimit.Reset.Equal(expectedTime) {
		t.Errorf("Expected reset time %v, got %v", expectedTime, c.RateLimit.Reset)
	}
	
	// Check pagination parsing
	if resp.Pagination == nil {
		t.Error("Expected pagination info, got nil")
	} else {
		if resp.Pagination.NextCursor != "next123" {
			t.Errorf("Expected next cursor 'next123', got %s", resp.Pagination.NextCursor)
		}
		if !resp.Pagination.HasNext {
			t.Error("Expected HasNext to be true")
		}
	}
}

func TestClientGet(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		
		// Check query parameters
		if r.URL.Query().Get("test") != "value" {
			t.Errorf("Expected query parameter test=value, got %s", r.URL.Query().Get("test"))
		}
		
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true}`))
	}))
	defer server.Close()

	os.Setenv("SENTRY_API_TOKEN", "test-token")
	defer os.Unsetenv("SENTRY_API_TOKEN")
	
	c, err := client.NewClient()
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	
	// Override base URL to use test server
	c.BaseURL = server.URL
	
	params := url.Values{}
	params.Set("test", "value")
	
	resp, err := c.Get("/test", params)
	if err != nil {
		t.Fatalf("GET request failed: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestClientError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Bad request"}`))
	}))
	defer server.Close()

	os.Setenv("SENTRY_API_TOKEN", "test-token")
	defer os.Unsetenv("SENTRY_API_TOKEN")
	
	c, err := client.NewClient()
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	
	c.BaseURL = server.URL
	
	_, err = c.Get("/test", nil)
	if err == nil {
		t.Error("Expected error for 400 response")
	}
	
	if !strings.Contains(err.Error(), "400") {
		t.Errorf("Expected error to contain status code 400, got: %v", err)
	}
}

func TestDecodeJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"name": "test", "value": 42}`))
	}))
	defer server.Close()

	os.Setenv("SENTRY_API_TOKEN", "test-token")
	defer os.Unsetenv("SENTRY_API_TOKEN")
	
	c, err := client.NewClient()
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	
	c.BaseURL = server.URL
	
	resp, err := c.Get("/test", nil)
	if err != nil {
		t.Fatalf("GET request failed: %v", err)
	}
	
	var result map[string]interface{}
	err = c.DecodeJSON(resp, &result)
	if err != nil {
		t.Fatalf("Failed to decode JSON: %v", err)
	}
	
	if result["name"] != "test" {
		t.Errorf("Expected name 'test', got %v", result["name"])
	}
	
	if result["value"] != float64(42) {
		t.Errorf("Expected value 42, got %v", result["value"])
	}
}