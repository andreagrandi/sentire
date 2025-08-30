package formatter

import (
	"encoding/json"
	"io"
	"sentire/pkg/models"
)

// JSONFormatter outputs data in JSON format (default behavior)
type JSONFormatter struct {
	writer io.Writer
}

// NewJSONFormatter creates a new JSON formatter
func NewJSONFormatter(writer io.Writer) *JSONFormatter {
	return &JSONFormatter{writer: writer}
}

// FormatEvent formats a single event as JSON
func (f *JSONFormatter) FormatEvent(event *models.Event) error {
	return f.FormatGeneric(event)
}

// FormatEvents formats multiple events as JSON
func (f *JSONFormatter) FormatEvents(events []models.Event) error {
	return f.FormatGeneric(events)
}

// FormatIssue formats a single issue as JSON
func (f *JSONFormatter) FormatIssue(issue *models.Issue) error {
	return f.FormatGeneric(issue)
}

// FormatIssues formats multiple issues as JSON
func (f *JSONFormatter) FormatIssues(issues []models.Issue) error {
	return f.FormatGeneric(issues)
}

// FormatProject formats a single project as JSON
func (f *JSONFormatter) FormatProject(project *models.Project) error {
	return f.FormatGeneric(project)
}

// FormatProjects formats multiple projects as JSON
func (f *JSONFormatter) FormatProjects(projects []models.Project) error {
	return f.FormatGeneric(projects)
}

// FormatOrgStats formats organization stats as JSON
func (f *JSONFormatter) FormatOrgStats(stats *models.OrganizationStats) error {
	return f.FormatGeneric(stats)
}

// FormatGeneric formats any data as JSON
func (f *JSONFormatter) FormatGeneric(data interface{}) error {
	encoder := json.NewEncoder(f.writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}