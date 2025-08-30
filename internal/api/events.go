package api

import (
	"fmt"
	"net/url"
	"sentire/internal/client"
	"sentire/pkg/models"
)

// EventsAPI provides methods for interacting with Sentry Events API
type EventsAPI struct {
	client *client.Client
}

// NewEventsAPI creates a new Events API client
func NewEventsAPI(client *client.Client) *EventsAPI {
	return &EventsAPI{client: client}
}

// ListProjectEventsOptions contains options for listing project events
type ListProjectEventsOptions struct {
	StatsPeriod string
	Start       string
	End         string
	Full        bool
	Sample      bool
	Cursor      string
}

// ListProjectEvents retrieves events for a specific project
func (e *EventsAPI) ListProjectEvents(orgSlug, projectSlug string, opts *ListProjectEventsOptions) ([]models.Event, *client.PaginationInfo, error) {
	endpoint := fmt.Sprintf("/projects/%s/%s/events/", orgSlug, projectSlug)
	
	params := url.Values{}
	if opts != nil {
		if opts.StatsPeriod != "" {
			params.Set("statsPeriod", opts.StatsPeriod)
		}
		if opts.Start != "" {
			params.Set("start", opts.Start)
		}
		if opts.End != "" {
			params.Set("end", opts.End)
		}
		if opts.Full {
			params.Set("full", "true")
		}
		if opts.Sample {
			params.Set("sample", "true")
		}
		if opts.Cursor != "" {
			params.Set("cursor", opts.Cursor)
		}
	}

	resp, err := e.client.Get(endpoint, params)
	if err != nil {
		return nil, nil, err
	}

	var events []models.Event
	if err := e.client.DecodeJSON(resp, &events); err != nil {
		return nil, nil, err
	}

	return events, resp.Pagination, nil
}

// ListIssueEventsOptions contains options for listing issue events
type ListIssueEventsOptions struct {
	Start       string
	End         string
	StatsPeriod string
	Environment []string
	Full        bool
	Sample      bool
	Query       string
	Cursor      string
}

// ListIssueEvents retrieves events for a specific issue
func (e *EventsAPI) ListIssueEvents(orgSlug string, issueID string, opts *ListIssueEventsOptions) ([]models.Event, *client.PaginationInfo, error) {
	endpoint := fmt.Sprintf("/organizations/%s/issues/%s/events/", orgSlug, issueID)
	
	params := url.Values{}
	if opts != nil {
		if opts.Start != "" {
			params.Set("start", opts.Start)
		}
		if opts.End != "" {
			params.Set("end", opts.End)
		}
		if opts.StatsPeriod != "" {
			params.Set("statsPeriod", opts.StatsPeriod)
		}
		for _, env := range opts.Environment {
			params.Add("environment", env)
		}
		if opts.Full {
			params.Set("full", "true")
		}
		if opts.Sample {
			params.Set("sample", "true")
		}
		if opts.Query != "" {
			params.Set("query", opts.Query)
		}
		if opts.Cursor != "" {
			params.Set("cursor", opts.Cursor)
		}
	}

	resp, err := e.client.Get(endpoint, params)
	if err != nil {
		return nil, nil, err
	}

	var events []models.Event
	if err := e.client.DecodeJSON(resp, &events); err != nil {
		return nil, nil, err
	}

	return events, resp.Pagination, nil
}

// ListIssuesOptions contains options for listing organization issues
type ListIssuesOptions struct {
	Environment []string
	Project     []string
	StatsPeriod string
	Start       string
	End         string
	Query       string
	Sort        string
	Limit       int
	Cursor      string
}

// ListIssues retrieves issues for an organization
func (e *EventsAPI) ListIssues(orgSlug string, opts *ListIssuesOptions) ([]models.Issue, *client.PaginationInfo, error) {
	endpoint := fmt.Sprintf("/organizations/%s/issues/", orgSlug)
	
	params := url.Values{}
	if opts != nil {
		for _, env := range opts.Environment {
			params.Add("environment", env)
		}
		for _, proj := range opts.Project {
			params.Add("project", proj)
		}
		if opts.StatsPeriod != "" {
			params.Set("statsPeriod", opts.StatsPeriod)
		}
		if opts.Start != "" {
			params.Set("start", opts.Start)
		}
		if opts.End != "" {
			params.Set("end", opts.End)
		}
		if opts.Query != "" {
			params.Set("query", opts.Query)
		}
		if opts.Sort != "" {
			params.Set("sort", opts.Sort)
		}
		if opts.Limit > 0 {
			params.Set("limit", fmt.Sprintf("%d", opts.Limit))
		}
		if opts.Cursor != "" {
			params.Set("cursor", opts.Cursor)
		}
	}

	resp, err := e.client.Get(endpoint, params)
	if err != nil {
		return nil, nil, err
	}

	var issues []models.Issue
	if err := e.client.DecodeJSON(resp, &issues); err != nil {
		return nil, nil, err
	}

	return issues, resp.Pagination, nil
}

// GetProjectEvent retrieves a specific event for a project
func (e *EventsAPI) GetProjectEvent(orgSlug, projectSlug, eventID string) (*models.Event, error) {
	endpoint := fmt.Sprintf("/projects/%s/%s/events/%s/", orgSlug, projectSlug, eventID)
	
	resp, err := e.client.Get(endpoint, nil)
	if err != nil {
		return nil, err
	}

	var event models.Event
	if err := e.client.DecodeJSON(resp, &event); err != nil {
		return nil, err
	}

	return &event, nil
}

// GetIssue retrieves a specific issue
func (e *EventsAPI) GetIssue(orgSlug, issueID string) (*models.Issue, error) {
	endpoint := fmt.Sprintf("/organizations/%s/issues/%s/", orgSlug, issueID)
	
	resp, err := e.client.Get(endpoint, nil)
	if err != nil {
		return nil, err
	}

	var issue models.Issue
	if err := e.client.DecodeJSON(resp, &issue); err != nil {
		return nil, err
	}

	return &issue, nil
}

// GetIssueEventOptions contains options for retrieving issue events
type GetIssueEventOptions struct {
	Environment []string
}

// GetIssueEvent retrieves a specific event for an issue
func (e *EventsAPI) GetIssueEvent(orgSlug, issueID, eventID string, opts *GetIssueEventOptions) (*models.Event, error) {
	endpoint := fmt.Sprintf("/organizations/%s/issues/%s/events/%s/", orgSlug, issueID, eventID)
	
	params := url.Values{}
	if opts != nil {
		for _, env := range opts.Environment {
			params.Add("environment", env)
		}
	}

	resp, err := e.client.Get(endpoint, params)
	if err != nil {
		return nil, err
	}

	var event models.Event
	if err := e.client.DecodeJSON(resp, &event); err != nil {
		return nil, err
	}

	return &event, nil
}