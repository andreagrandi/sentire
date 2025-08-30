package tests

import (
	"encoding/json"
	"net/http"
	"os"
	"sentire/internal/api"
	"sentire/pkg/models"
	"strings"
	"testing"
	"time"
)

func TestListAllProjects(t *testing.T) {
	expectedProjects := []models.Project{
		{
			ID:   "1",
			Slug: "project1",
			Name: "Project One",
			Organization: models.Organization{
				ID:   "org1",
				Slug: "org1",
				Name: "Organization One",
			},
			DateCreated: time.Now(),
			Platform:    "python",
		},
		{
			ID:   "2",
			Slug: "project2",
			Name: "Project Two",
			Organization: models.Organization{
				ID:   "org2",
				Slug: "org2",
				Name: "Organization Two",
			},
			DateCreated: time.Now(),
			Platform:    "javascript",
		},
	}

	c, server := setupTestClient(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/projects/" {
			t.Errorf("Expected path '/projects/', got %s", r.URL.Path)
		}
		
		cursor := r.URL.Query().Get("cursor")
		if cursor != "" && cursor != "test-cursor" {
			t.Errorf("Expected cursor 'test-cursor' or empty, got %s", cursor)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expectedProjects)
	})
	defer server.Close()
	defer os.Unsetenv("SENTRY_API_TOKEN")

	projectsAPI := api.NewProjectsAPI(c)
	
	opts := &api.ListAllProjectsOptions{
		Cursor: "test-cursor",
	}
	
	projects, pagination, err := projectsAPI.ListProjects(opts)
	if err != nil {
		t.Fatalf("ListProjects failed: %v", err)
	}
	
	if len(projects) != 2 {
		t.Errorf("Expected 2 projects, got %d", len(projects))
	}
	
	if projects[0].Slug != "project1" {
		t.Errorf("Expected first project slug 'project1', got %s", projects[0].Slug)
	}
	
	if projects[0].Platform != "python" {
		t.Errorf("Expected first project platform 'python', got %s", projects[0].Platform)
	}
	
	if projects[1].Name != "Project Two" {
		t.Errorf("Expected second project name 'Project Two', got %s", projects[1].Name)
	}
	
	if projects[1].Organization.Slug != "org2" {
		t.Errorf("Expected second project org slug 'org2', got %s", projects[1].Organization.Slug)
	}
	
	if pagination == nil {
		t.Error("Expected pagination info, got nil")
	}
}

func TestListAllProjectsWithNilOptions(t *testing.T) {
	expectedProjects := []models.Project{
		{
			ID:   "1",
			Slug: "project1",
			Name: "Project One",
		},
	}

	c, server := setupTestClient(func(w http.ResponseWriter, r *http.Request) {
		// Should not have cursor parameter
		cursor := r.URL.Query().Get("cursor")
		if cursor != "" {
			t.Errorf("Expected no cursor parameter, got %s", cursor)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expectedProjects)
	})
	defer server.Close()
	defer os.Unsetenv("SENTRY_API_TOKEN")

	projectsAPI := api.NewProjectsAPI(c)
	
	projects, _, err := projectsAPI.ListProjects(nil)
	if err != nil {
		t.Fatalf("ListProjects failed: %v", err)
	}
	
	if len(projects) != 1 {
		t.Errorf("Expected 1 project, got %d", len(projects))
	}
}

func TestGetProject(t *testing.T) {
	now := time.Now()
	expectedProject := models.Project{
		ID:          "project123",
		Slug:        "test-project",
		Name:        "Test Project",
		IsPublic:    false,
		IsBookmarked: true,
		Color:       "#3fbf3f",
		DateCreated: now,
		Platform:    "python",
		Platforms:   []string{"python", "javascript"},
		HasAccess:   true,
		Features:    []string{"releases", "alert-filters"},
		Status:      "active",
		Organization: models.Organization{
			ID:   "org123",
			Slug: "test-org",
			Name: "Test Organization",
		},
		Teams: []models.ProjectTeam{
			{
				ID:   "team1",
				Slug: "backend",
				Name: "Backend Team",
			},
		},
	}

	c, server := setupTestClient(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/projects/test-org/test-project/"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expectedProject)
	})
	defer server.Close()
	defer os.Unsetenv("SENTRY_API_TOKEN")

	projectsAPI := api.NewProjectsAPI(c)
	
	project, err := projectsAPI.GetProject("test-org", "test-project")
	if err != nil {
		t.Fatalf("GetProject failed: %v", err)
	}
	
	if project.ID != "project123" {
		t.Errorf("Expected project ID 'project123', got %s", project.ID)
	}
	
	if project.Slug != "test-project" {
		t.Errorf("Expected project slug 'test-project', got %s", project.Slug)
	}
	
	if project.Name != "Test Project" {
		t.Errorf("Expected project name 'Test Project', got %s", project.Name)
	}
	
	if !project.IsBookmarked {
		t.Error("Expected project to be bookmarked")
	}
	
	if project.Platform != "python" {
		t.Errorf("Expected platform 'python', got %s", project.Platform)
	}
	
	if len(project.Platforms) != 2 {
		t.Errorf("Expected 2 platforms, got %d", len(project.Platforms))
	}
	
	if len(project.Features) != 2 {
		t.Errorf("Expected 2 features, got %d", len(project.Features))
	}
	
	if project.Organization.Slug != "test-org" {
		t.Errorf("Expected organization slug 'test-org', got %s", project.Organization.Slug)
	}
	
	if len(project.Teams) != 1 {
		t.Errorf("Expected 1 team, got %d", len(project.Teams))
	}
	
	if project.Teams[0].Slug != "backend" {
		t.Errorf("Expected team slug 'backend', got %s", project.Teams[0].Slug)
	}
}

func TestGetProjectNotFound(t *testing.T) {
	c, server := setupTestClient(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"detail": "The requested resource does not exist"}`))
	})
	defer server.Close()
	defer os.Unsetenv("SENTRY_API_TOKEN")

	projectsAPI := api.NewProjectsAPI(c)
	
	_, err := projectsAPI.GetProject("nonexistent-org", "nonexistent-project")
	if err == nil {
		t.Error("Expected error for 404 response")
	}
	
	if !strings.Contains(err.Error(), "404") {
		t.Errorf("Expected error to contain status code 404, got: %v", err)
	}
}