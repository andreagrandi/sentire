package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	BaseURL     = "https://sentry.io/api/0"
	UserAgent   = "sentire/1.0.0"
	TokenEnvVar = "SENTRY_API_TOKEN"
)

// Client represents the Sentry API client
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	Token      string
	RateLimit  *RateLimiter
}

// RateLimiter tracks rate limit information
type RateLimiter struct {
	Limit           int
	Remaining       int
	Reset           time.Time
	ConcurrentLimit int
	ConcurrentRemaining int
}

// PaginationInfo contains pagination metadata
type PaginationInfo struct {
	NextCursor string
	PrevCursor string
	HasNext    bool
	HasPrev    bool
}

// Response wraps HTTP responses with pagination info
type Response struct {
	*http.Response
	Pagination *PaginationInfo
}

// NewClient creates a new Sentry API client
func NewClient() (*Client, error) {
	token := os.Getenv(TokenEnvVar)
	if token == "" {
		return nil, fmt.Errorf("environment variable %s is required", TokenEnvVar)
	}

	return &Client{
		BaseURL: BaseURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		Token:     token,
		RateLimit: &RateLimiter{},
	}, nil
}

// Do executes an HTTP request and returns the response
func (c *Client) Do(req *http.Request) (*Response, error) {
	// Set required headers
	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", UserAgent)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}

	// Parse rate limit headers
	c.parseRateLimitHeaders(resp)

	// Parse pagination from Link header
	pagination := c.parseLinkHeader(resp.Header.Get("Link"))

	response := &Response{
		Response:   resp,
		Pagination: pagination,
	}

	// Handle HTTP errors
	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		return response, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return response, nil
}

// Get performs a GET request
func (c *Client) Get(endpoint string, params url.Values) (*Response, error) {
	fullURL := c.BaseURL + endpoint
	if params != nil {
		fullURL += "?" + params.Encode()
	}

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	return c.Do(req)
}

// DecodeJSON decodes JSON response into the provided interface
func (c *Client) DecodeJSON(resp *Response, v interface{}) error {
	defer resp.Body.Close()
	
	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		return fmt.Errorf("failed to decode JSON response: %w", err)
	}
	
	return nil
}

// parseRateLimitHeaders extracts rate limiting information from response headers
func (c *Client) parseRateLimitHeaders(resp *http.Response) {
	if limit := resp.Header.Get("X-Sentry-Rate-Limit-Limit"); limit != "" {
		if val, err := strconv.Atoi(limit); err == nil {
			c.RateLimit.Limit = val
		}
	}

	if remaining := resp.Header.Get("X-Sentry-Rate-Limit-Remaining"); remaining != "" {
		if val, err := strconv.Atoi(remaining); err == nil {
			c.RateLimit.Remaining = val
		}
	}

	if reset := resp.Header.Get("X-Sentry-Rate-Limit-Reset"); reset != "" {
		if val, err := strconv.ParseInt(reset, 10, 64); err == nil {
			c.RateLimit.Reset = time.Unix(val, 0)
		}
	}

	if concurrentLimit := resp.Header.Get("X-Sentry-Rate-Limit-ConcurrentLimit"); concurrentLimit != "" {
		if val, err := strconv.Atoi(concurrentLimit); err == nil {
			c.RateLimit.ConcurrentLimit = val
		}
	}

	if concurrentRemaining := resp.Header.Get("X-Sentry-Rate-Limit-ConcurrentRemaining"); concurrentRemaining != "" {
		if val, err := strconv.Atoi(concurrentRemaining); err == nil {
			c.RateLimit.ConcurrentRemaining = val
		}
	}
}

// parseLinkHeader parses the Link header for pagination information
func (c *Client) parseLinkHeader(linkHeader string) *PaginationInfo {
	info := &PaginationInfo{}
	
	if linkHeader == "" {
		return info
	}

	links := strings.Split(linkHeader, ",")
	for _, link := range links {
		link = strings.TrimSpace(link)
		parts := strings.Split(link, ";")
		if len(parts) < 2 {
			continue
		}

		urlPart := strings.Trim(strings.TrimSpace(parts[0]), "<>")
		
		// Extract cursor from URL
		if u, err := url.Parse(urlPart); err == nil {
			cursor := u.Query().Get("cursor")
			
			// Check all parts for rel and results attributes
			var isNext, isPrev, hasResults bool
			for i := 1; i < len(parts); i++ {
				part := strings.TrimSpace(parts[i])
				if strings.Contains(part, `rel="next"`) {
					isNext = true
				} else if strings.Contains(part, `rel="previous"`) {
					isPrev = true
				}
				if strings.Contains(part, `results="true"`) {
					hasResults = true
				}
			}
			
			if isNext {
				info.NextCursor = cursor
				info.HasNext = hasResults
			} else if isPrev {
				info.PrevCursor = cursor
				info.HasPrev = true
			}
		}
	}

	return info
}