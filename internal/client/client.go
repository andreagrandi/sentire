package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/andreagrandi/sentire/internal/config"
	"github.com/andreagrandi/sentire/internal/redact"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	BaseURL   = "https://sentry.io/api/0"
	UserAgent = "sentire/1.0.0"

	// DefaultTimeout bounds a single API request, including reading the
	// response body. It applies to every command unless the caller's
	// context carries an earlier deadline.
	DefaultTimeout = 30 * time.Second
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
	Limit               int
	Remaining           int
	Reset               time.Time
	ConcurrentLimit     int
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

// APIError represents an error returned by the Sentry API
type APIError struct {
	Message    string
	StatusCode int
}

func (e *APIError) Error() string {
	return e.Message
}

// RequestError represents a transport-level request failure. Timeout and
// Canceled distinguish the two cases the CLI reports specially; both are
// false for other transport failures (DNS, connection refused, etc.).
type RequestError struct {
	Message  string
	Timeout  bool
	Canceled bool
}

func (e *RequestError) Error() string {
	return e.Message
}

// requestErrorIfContext classifies err as a timeout or cancellation when it
// originates from a context deadline, a context cancellation, or any other
// network timeout. It returns nil when err is none of these.
func requestErrorIfContext(err error) *RequestError {
	switch {
	case errors.Is(err, context.Canceled):
		return &RequestError{Message: "request canceled", Canceled: true}
	case errors.Is(err, context.DeadlineExceeded):
		return &RequestError{Message: "request timed out", Timeout: true}
	}

	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return &RequestError{Message: "request timed out", Timeout: true}
	}

	return nil
}

// classifyRequestError converts a transport error from the HTTP client into a
// RequestError, separating timeouts and cancellations from other failures.
func classifyRequestError(err error, token string) *RequestError {
	if reqErr := requestErrorIfContext(err); reqErr != nil {
		return reqErr
	}
	// Transport errors come from net/http and may include the request URL;
	// strip the token defensively in case it ever surfaces in a wrapped error.
	return &RequestError{Message: fmt.Sprintf("http request failed: %s", redact.Secret(err.Error(), token))}
}

// NewClient creates a new Sentry API client
func NewClient() (*Client, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	return &Client{
		BaseURL: BaseURL,
		HTTPClient: &http.Client{
			Timeout: DefaultTimeout,
		},
		Token:     cfg.SentryAPIToken,
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
		return nil, classifyRequestError(err, c.Token)
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
		return response, &APIError{
			Message:    fmt.Sprintf("API request failed with status %d: %s", resp.StatusCode, redact.Secret(string(body), c.Token)),
			StatusCode: resp.StatusCode,
		}
	}

	return response, nil
}

// Get performs a context-aware GET request. The context propagates
// cancellation and any deadline from the caller into the HTTP request.
func (c *Client) Get(ctx context.Context, endpoint string, params url.Values) (*Response, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	fullURL := c.BaseURL + endpoint
	if params != nil {
		fullURL += "?" + params.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	return c.Do(req)
}

// DecodeJSON decodes JSON response into the provided interface
func (c *Client) DecodeJSON(resp *Response, v interface{}) error {
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		// A timeout or cancellation while streaming the body surfaces here
		// rather than from Do; report it as the request error it is.
		if reqErr := requestErrorIfContext(err); reqErr != nil {
			return reqErr
		}
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
