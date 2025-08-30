package api

import (
	"fmt"
	"net/url"
	"sentire/internal/client"
	"sentire/pkg/models"
)

// OrganizationsAPI provides methods for interacting with Sentry Organizations API
type OrganizationsAPI struct {
	client *client.Client
}

// NewOrganizationsAPI creates a new Organizations API client
func NewOrganizationsAPI(client *client.Client) *OrganizationsAPI {
	return &OrganizationsAPI{client: client}
}

// ListProjectsOptions contains options for listing organization projects
type ListProjectsOptions struct {
	Cursor string
}

// ListProjects retrieves projects for an organization
func (o *OrganizationsAPI) ListProjects(orgSlug string, opts *ListProjectsOptions) ([]models.Project, *client.PaginationInfo, error) {
	endpoint := fmt.Sprintf("/organizations/%s/projects/", orgSlug)

	params := url.Values{}
	if opts != nil && opts.Cursor != "" {
		params.Set("cursor", opts.Cursor)
	}

	resp, err := o.client.Get(endpoint, params)
	if err != nil {
		return nil, nil, err
	}

	var projects []models.Project
	if err := o.client.DecodeJSON(resp, &projects); err != nil {
		return nil, nil, err
	}

	return projects, resp.Pagination, nil
}

// GetStatsOptions contains options for retrieving organization statistics
type GetStatsOptions struct {
	Field       string   // Required: "sum(quantity)" or "sum(times_seen)"
	StatsPeriod string   // Time range (e.g., "1d", "7d")
	Interval    string   // Time series resolution
	Start       string   // Start time (ISO-8601)
	End         string   // End time (ISO-8601)
	Project     []string // Project ID filters
	Category    []string // Event category filters
	Outcome     []string // Event outcome filters
	Reason      []string // Event reason filters
	Download    bool     // Download as CSV
}

// GetStats retrieves event statistics for an organization
func (o *OrganizationsAPI) GetStats(orgSlug string, opts *GetStatsOptions) (*models.OrganizationStats, error) {
	if opts == nil || opts.Field == "" {
		return nil, fmt.Errorf("field parameter is required")
	}

	endpoint := fmt.Sprintf("/organizations/%s/stats-summary/", orgSlug)

	params := url.Values{}
	params.Set("field", opts.Field)

	if opts.StatsPeriod != "" {
		params.Set("statsPeriod", opts.StatsPeriod)
	}
	if opts.Interval != "" {
		params.Set("interval", opts.Interval)
	}
	if opts.Start != "" {
		params.Set("start", opts.Start)
	}
	if opts.End != "" {
		params.Set("end", opts.End)
	}
	for _, proj := range opts.Project {
		params.Add("project", proj)
	}
	for _, cat := range opts.Category {
		params.Add("category", cat)
	}
	for _, outcome := range opts.Outcome {
		params.Add("outcome", outcome)
	}
	for _, reason := range opts.Reason {
		params.Add("reason", reason)
	}
	if opts.Download {
		params.Set("download", "true")
	}

	resp, err := o.client.Get(endpoint, params)
	if err != nil {
		return nil, err
	}

	var stats models.OrganizationStats
	if err := o.client.DecodeJSON(resp, &stats); err != nil {
		return nil, err
	}

	return &stats, nil
}
