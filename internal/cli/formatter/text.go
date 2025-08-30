package formatter

import (
	"fmt"
	"io"
	"reflect"
	"sentire/pkg/models"
	"strings"
)

// TextFormatter outputs data in plain text format
type TextFormatter struct {
	writer io.Writer
}

// NewTextFormatter creates a new text formatter
func NewTextFormatter(writer io.Writer) *TextFormatter {
	return &TextFormatter{writer: writer}
}

// FormatEvent formats a single event as text
func (f *TextFormatter) FormatEvent(event *models.Event) error {
	fmt.Fprintf(f.writer, "Event #%s\n", event.ID)
	fmt.Fprintf(f.writer, "Event ID: %s\n", event.EventID)
	fmt.Fprintf(f.writer, "Title: %s\n", event.Title)
	fmt.Fprintf(f.writer, "Message: %s\n", event.Message)
	fmt.Fprintf(f.writer, "Type: %s\n", event.Type)
	fmt.Fprintf(f.writer, "Platform: %s\n", event.Platform)
	fmt.Fprintf(f.writer, "Project ID: %s\n", event.ProjectID)
	fmt.Fprintf(f.writer, "Date Created: %s\n", event.DateCreated.Format("2006-01-02 15:04:05"))
	fmt.Fprintf(f.writer, "Date Received: %s\n", event.DateReceived.Format("2006-01-02 15:04:05"))
	fmt.Fprintf(f.writer, "Size: %d bytes\n", event.Size)

	if event.GroupID != "" {
		fmt.Fprintf(f.writer, "Group ID: %s\n", event.GroupID)
	}
	if event.Logger != "" {
		fmt.Fprintf(f.writer, "Logger: %s\n", event.Logger)
	}
	if event.Culprit != "" {
		fmt.Fprintf(f.writer, "Culprit: %s\n", event.Culprit)
	}
	if event.Environment != "" {
		fmt.Fprintf(f.writer, "Environment: %s\n", event.Environment)
	}

	// Add stack trace information if available
	if len(event.Entries) > 0 {
		fmt.Fprintf(f.writer, "\nEntries:\n")
		for i, entry := range event.Entries {
			if i > 2 { // Limit to first 3 entries for readability
				break
			}
			fmt.Fprintf(f.writer, "  %d. Type: %s\n", i+1, entry.Type)
		}
	}

	fmt.Fprintf(f.writer, "\n")
	return nil
}

// FormatEvents formats multiple events as text
func (f *TextFormatter) FormatEvents(events []models.Event) error {
	if len(events) == 0 {
		fmt.Fprintf(f.writer, "No events found\n")
		return nil
	}

	fmt.Fprintf(f.writer, "Events (%d total):\n\n", len(events))

	for i, event := range events {
		fmt.Fprintf(f.writer, "%d. Event #%s\n", i+1, event.EventID)
		fmt.Fprintf(f.writer, "   Title: %s\n", event.Title)
		fmt.Fprintf(f.writer, "   Type: %s | Platform: %s | Project ID: %s\n",
			event.Type, event.Platform, event.ProjectID)
		fmt.Fprintf(f.writer, "   Date: %s | Environment: %s\n",
			event.DateCreated.Format("2006-01-02 15:04"), event.Environment)
		fmt.Fprintf(f.writer, "\n")
	}

	return nil
}

// FormatIssue formats a single issue as text
func (f *TextFormatter) FormatIssue(issue *models.Issue) error {
	fmt.Fprintf(f.writer, "Issue #%s (%s)\n", issue.ID, issue.ShortID)
	fmt.Fprintf(f.writer, "Title: %s\n", issue.Title)
	fmt.Fprintf(f.writer, "Level: %s\n", issue.Level)
	fmt.Fprintf(f.writer, "Status: %s", issue.Status)

	if issue.Substatus != "" {
		fmt.Fprintf(f.writer, " (%s)", issue.Substatus)
	}
	fmt.Fprintf(f.writer, "\n")

	if issue.Priority != "" {
		fmt.Fprintf(f.writer, "Priority: %s\n", issue.Priority)
	}

	fmt.Fprintf(f.writer, "Platform: %s\n", issue.Platform)
	fmt.Fprintf(f.writer, "Project: %s (%s)\n", issue.Project.Name, issue.Project.Slug)
	fmt.Fprintf(f.writer, "Count: %s\n", issue.Count)
	fmt.Fprintf(f.writer, "User Count: %d\n", issue.UserCount)
	fmt.Fprintf(f.writer, "First Seen: %s\n", issue.FirstSeen.Format("2006-01-02 15:04:05"))
	fmt.Fprintf(f.writer, "Last Seen: %s\n", issue.LastSeen.Format("2006-01-02 15:04:05"))

	if issue.Culprit != "" {
		fmt.Fprintf(f.writer, "Culprit: %s\n", issue.Culprit)
	}

	if issue.Logger != "" {
		fmt.Fprintf(f.writer, "Logger: %s\n", issue.Logger)
	}

	fmt.Fprintf(f.writer, "Public: %v | Bookmarked: %v | Subscribed: %v\n",
		issue.IsPublic, issue.IsBookmarked, issue.IsSubscribed)

	if issue.Permalink != "" {
		fmt.Fprintf(f.writer, "Permalink: %s\n", issue.Permalink)
	}

	fmt.Fprintf(f.writer, "\n")
	return nil
}

// FormatIssues formats multiple issues as text
func (f *TextFormatter) FormatIssues(issues []models.Issue) error {
	if len(issues) == 0 {
		fmt.Fprintf(f.writer, "No issues found\n")
		return nil
	}

	fmt.Fprintf(f.writer, "Issues (%d total):\n\n", len(issues))

	for i, issue := range issues {
		fmt.Fprintf(f.writer, "%d. Issue #%s\n", i+1, issue.ShortID)
		fmt.Fprintf(f.writer, "   Title: %s\n", issue.Title)
		fmt.Fprintf(f.writer, "   Level: %s | Status: %s | Count: %s\n",
			issue.Level, issue.Status, issue.Count)
		fmt.Fprintf(f.writer, "   Project: %s | Users: %d\n",
			issue.Project.Slug, issue.UserCount)
		fmt.Fprintf(f.writer, "   Last Seen: %s\n",
			issue.LastSeen.Format("2006-01-02 15:04"))
		fmt.Fprintf(f.writer, "\n")
	}

	return nil
}

// FormatProject formats a single project as text
func (f *TextFormatter) FormatProject(project *models.Project) error {
	fmt.Fprintf(f.writer, "Project: %s\n", project.Name)
	fmt.Fprintf(f.writer, "ID: %s\n", project.ID)
	fmt.Fprintf(f.writer, "Slug: %s\n", project.Slug)
	fmt.Fprintf(f.writer, "Platform: %s\n", project.Platform)
	fmt.Fprintf(f.writer, "Organization: %s (%s)\n", project.Organization.Name, project.Organization.Slug)
	fmt.Fprintf(f.writer, "Status: %s\n", project.Status)
	fmt.Fprintf(f.writer, "Date Created: %s\n", project.DateCreated.Format("2006-01-02 15:04:05"))
	fmt.Fprintf(f.writer, "Public: %v | Bookmarked: %v\n", project.IsPublic, project.IsBookmarked)

	fmt.Fprintf(f.writer, "\n")
	return nil
}

// FormatProjects formats multiple projects as text
func (f *TextFormatter) FormatProjects(projects []models.Project) error {
	if len(projects) == 0 {
		fmt.Fprintf(f.writer, "No projects found\n")
		return nil
	}

	fmt.Fprintf(f.writer, "Projects (%d total):\n\n", len(projects))

	for i, project := range projects {
		fmt.Fprintf(f.writer, "%d. %s (%s)\n", i+1, project.Name, project.Slug)
		fmt.Fprintf(f.writer, "   Platform: %s | Organization: %s\n",
			project.Platform, project.Organization.Slug)
		fmt.Fprintf(f.writer, "   Status: %s | Created: %s\n",
			project.Status, project.DateCreated.Format("2006-01-02"))
		fmt.Fprintf(f.writer, "\n")
	}

	return nil
}

// FormatOrgStats formats organization stats as text
func (f *TextFormatter) FormatOrgStats(stats *models.OrganizationStats) error {
	fmt.Fprintf(f.writer, "Organization Statistics\n")
	fmt.Fprintf(f.writer, "======================\n\n")

	fmt.Fprintf(f.writer, "Period Start: %s\n", stats.Start.Format("2006-01-02 15:04:05"))
	fmt.Fprintf(f.writer, "Period End: %s\n", stats.End.Format("2006-01-02 15:04:05"))
	fmt.Fprintf(f.writer, "Total Sum: %d\n", stats.Totals.Sum)
	fmt.Fprintf(f.writer, "Times Seen: %d\n", stats.Totals.TimesSeen)

	if len(stats.Projects) > 0 {
		fmt.Fprintf(f.writer, "\nProjects (%d):\n", len(stats.Projects))
		for i, project := range stats.Projects {
			if i > 4 { // Limit for readability
				fmt.Fprintf(f.writer, "... and %d more projects\n", len(stats.Projects)-i)
				break
			}
			fmt.Fprintf(f.writer, "  %d. %s (%s)\n", i+1, project.Slug, project.ID)
		}
	}

	fmt.Fprintf(f.writer, "\n")
	return nil
}

// FormatGeneric formats any data as text
func (f *TextFormatter) FormatGeneric(data interface{}) error {
	v := reflect.ValueOf(data)

	// Handle slices/arrays
	if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
		if v.Len() == 0 {
			fmt.Fprintf(f.writer, "No data found\n")
			return nil
		}

		// For slice of interface{}, try to determine the type
		firstElem := v.Index(0).Interface()

		// Handle known types
		switch firstElem.(type) {
		case models.Event:
			events := make([]models.Event, v.Len())
			for i := 0; i < v.Len(); i++ {
				events[i] = v.Index(i).Interface().(models.Event)
			}
			return f.FormatEvents(events)
		case models.Issue:
			issues := make([]models.Issue, v.Len())
			for i := 0; i < v.Len(); i++ {
				issues[i] = v.Index(i).Interface().(models.Issue)
			}
			return f.FormatIssues(issues)
		case models.Project:
			projects := make([]models.Project, v.Len())
			for i := 0; i < v.Len(); i++ {
				projects[i] = v.Index(i).Interface().(models.Project)
			}
			return f.FormatProjects(projects)
		default:
			// Fallback to simple list for unknown types
			return f.formatUnknownSlice(data)
		}
	}

	// Handle single values
	return f.formatSingleValue(data)
}

// formatUnknownSlice formats a slice of unknown type
func (f *TextFormatter) formatUnknownSlice(data interface{}) error {
	v := reflect.ValueOf(data)

	fmt.Fprintf(f.writer, "Data (%d items):\n\n", v.Len())

	for i := 0; i < v.Len(); i++ {
		elem := v.Index(i)
		fmt.Fprintf(f.writer, "%d. %v\n", i+1, elem.Interface())
	}

	fmt.Fprintf(f.writer, "\n")
	return nil
}

// formatSingleValue formats a single value as text
func (f *TextFormatter) formatSingleValue(data interface{}) error {
	v := reflect.ValueOf(data)
	t := reflect.TypeOf(data)

	// Handle pointers
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			fmt.Fprintf(f.writer, "Value: <nil>\n")
			return nil
		}
		v = v.Elem()
		t = t.Elem()
	}

	// Handle structs
	if v.Kind() == reflect.Struct {
		typeName := t.Name()
		if typeName == "" {
			typeName = "Data"
		}
		fmt.Fprintf(f.writer, "%s:\n", typeName)
		fmt.Fprintf(f.writer, "%s\n", strings.Repeat("=", len(typeName)+1))

		for i := 0; i < v.NumField(); i++ {
			field := t.Field(i)
			value := v.Field(i)

			if !value.CanInterface() {
				continue
			}

			fmt.Fprintf(f.writer, "%s: %v\n", field.Name, value.Interface())
		}
	} else {
		// Simple value
		fmt.Fprintf(f.writer, "Value: %v\n", data)
	}

	fmt.Fprintf(f.writer, "\n")
	return nil
}
