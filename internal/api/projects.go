package api

import (
	"fmt"
	"net/url"
	"sentire/internal/client"
	"sentire/pkg/models"
)

// ProjectsAPI provides methods for interacting with Sentry Projects API
type ProjectsAPI struct {
	client *client.Client
}

// NewProjectsAPI creates a new Projects API client
func NewProjectsAPI(client *client.Client) *ProjectsAPI {
	return &ProjectsAPI{client: client}
}

// ListAllProjectsOptions contains options for listing all projects
type ListAllProjectsOptions struct {
	Cursor string
}

// ListProjects retrieves all projects the user has access to
func (p *ProjectsAPI) ListProjects(opts *ListAllProjectsOptions) ([]models.Project, *client.PaginationInfo, error) {
	endpoint := "/projects/"

	params := url.Values{}
	if opts != nil && opts.Cursor != "" {
		params.Set("cursor", opts.Cursor)
	}

	resp, err := p.client.Get(endpoint, params)
	if err != nil {
		return nil, nil, err
	}

	var projects []models.Project
	if err := p.client.DecodeJSON(resp, &projects); err != nil {
		return nil, nil, err
	}

	return projects, resp.Pagination, nil
}

// GetProject retrieves a specific project
func (p *ProjectsAPI) GetProject(orgSlug, projectSlug string) (*models.Project, error) {
	endpoint := fmt.Sprintf("/projects/%s/%s/", orgSlug, projectSlug)

	resp, err := p.client.Get(endpoint, nil)
	if err != nil {
		return nil, err
	}

	var project models.Project
	if err := p.client.DecodeJSON(resp, &project); err != nil {
		return nil, err
	}

	return &project, nil
}
