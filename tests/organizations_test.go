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

func TestListOrgProjects(t *testing.T) {
	expectedProjects := []models.Project{
		{
			ID:   "1",
			Slug: "test-project",
			Name: "Test Project",
			Organization: models.Organization{
				ID:   "org1",
				Slug: "test-org",
				Name: "Test Org",
			},
		},
		{
			ID:   "2",
			Slug: "another-project",
			Name: "Another Project",
			Organization: models.Organization{
				ID:   "org1",
				Slug: "test-org",
				Name: "Test Org",
			},
		},
	}

	c, server := setupTestClient(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "/organizations/test-org/projects/") {
			t.Errorf("Expected path to contain '/organizations/test-org/projects/', got %s", r.URL.Path)
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

	orgAPI := api.NewOrganizationsAPI(c)
	
	opts := &api.ListProjectsOptions{
		Cursor: "test-cursor",
	}
	
	projects, pagination, err := orgAPI.ListProjects("test-org", opts)
	if err != nil {
		t.Fatalf("ListProjects failed: %v", err)
	}
	
	if len(projects) != 2 {
		t.Errorf("Expected 2 projects, got %d", len(projects))
	}
	
	if projects[0].Slug != "test-project" {
		t.Errorf("Expected first project slug 'test-project', got %s", projects[0].Slug)
	}
	
	if projects[1].Name != "Another Project" {
		t.Errorf("Expected second project name 'Another Project', got %s", projects[1].Name)
	}
	
	if pagination == nil {
		t.Error("Expected pagination info, got nil")
	}
}

func TestGetOrgStats(t *testing.T) {
	expectedStats := models.OrganizationStats{
		Start: time.Now().Add(-24 * time.Hour),
		End:   time.Now(),
		Projects: []models.ProjectStatsDetail{
			{
				ID: "1",
				Slug: "project1",
				Stats: []models.CategoryStats{
					{
						Category: "error",
						Outcomes: models.StatsOutcomes{
							Accepted: 1000,
							Filtered: 0,
						},
						Totals: models.StatsTotals{
							Sum: 1000,
							TimesSeen: 500,
						},
					},
				},
			},
			{
				ID: "2", 
				Slug: "project2",
				Stats: []models.CategoryStats{
					{
						Category: "transaction",
						Outcomes: models.StatsOutcomes{
							Accepted: 2000,
							Filtered: 0,
						},
						Totals: models.StatsTotals{
							Sum: 2000,
							TimesSeen: 750,
						},
					},
				},
			},
		},
		Totals: struct {
			Sum       int64 `json:"sum"`
			TimesSeen int64 `json:"times_seen"`
		}{
			Sum:       3000,
			TimesSeen: 1250,
		},
	}

	c, server := setupTestClient(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "/organizations/test-org/stats-summary/") {
			t.Errorf("Expected path to contain '/organizations/test-org/stats-summary/', got %s", r.URL.Path)
		}
		
		field := r.URL.Query().Get("field")
		if field != "sum(quantity)" {
			t.Errorf("Expected field 'sum(quantity)', got %s", field)
		}
		
		period := r.URL.Query().Get("statsPeriod")
		if period != "7d" {
			t.Errorf("Expected statsPeriod '7d', got %s", period)
		}
		
		projects := r.URL.Query()["project"]
		if len(projects) != 2 || projects[0] != "1" || projects[1] != "2" {
			t.Errorf("Expected projects [1, 2], got %v", projects)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expectedStats)
	})
	defer server.Close()
	defer os.Unsetenv("SENTRY_API_TOKEN")

	orgAPI := api.NewOrganizationsAPI(c)
	
	opts := &api.GetStatsOptions{
		Field:       "sum(quantity)",
		StatsPeriod: "7d",
		Project:     []string{"1", "2"},
	}
	
	stats, err := orgAPI.GetStats("test-org", opts)
	if err != nil {
		t.Fatalf("GetStats failed: %v", err)
	}
	
	if len(stats.Projects) != 2 {
		t.Errorf("Expected 2 projects in stats, got %d", len(stats.Projects))
	}
	
	if stats.Projects[0].ID != "1" {
		t.Errorf("Expected first project ID '1', got %s", stats.Projects[0].ID)
	}
	
	if len(stats.Projects[0].Stats) > 0 && stats.Projects[0].Stats[0].Totals.Sum != 1000 {
		t.Errorf("Expected first project quantity 1000, got %d", stats.Projects[0].Stats[0].Totals.Sum)
	}
	
	if stats.Totals.Sum != 3000 {
		t.Errorf("Expected total sum 3000, got %d", stats.Totals.Sum)
	}
}

func TestGetStatsWithoutField(t *testing.T) {
	c, server := setupTestClient(func(w http.ResponseWriter, r *http.Request) {
		// Should not be called
		t.Error("Request should not be made without field parameter")
	})
	defer server.Close()
	defer os.Unsetenv("SENTRY_API_TOKEN")

	orgAPI := api.NewOrganizationsAPI(c)
	
	// Test with nil options
	_, err := orgAPI.GetStats("test-org", nil)
	if err == nil {
		t.Error("Expected error when opts is nil")
	}
	
	// Test with empty field
	opts := &api.GetStatsOptions{}
	_, err = orgAPI.GetStats("test-org", opts)
	if err == nil {
		t.Error("Expected error when field is empty")
	}
	
	if !strings.Contains(err.Error(), "field parameter is required") {
		t.Errorf("Expected error message to contain 'field parameter is required', got: %v", err)
	}
}