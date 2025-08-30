package formatter

import (
	"fmt"
	"io"
	"reflect"
	"sentire/pkg/models"
	"strconv"

	"github.com/olekukonko/tablewriter"
)

// TableFormatter outputs data in table format
type TableFormatter struct {
	writer io.Writer
}

// NewTableFormatter creates a new table formatter
func NewTableFormatter(writer io.Writer) *TableFormatter {
	return &TableFormatter{writer: writer}
}

// FormatEvent formats a single event as a table
func (f *TableFormatter) FormatEvent(event *models.Event) error {
	table := tablewriter.NewWriter(f.writer)
	table.Header("Field", "Value")

	rows := [][]string{
		{"ID", event.ID},
		{"Event ID", event.EventID},
		{"Title", event.Title},
		{"Message", event.Message},
		{"Platform", event.Platform},
		{"Type", event.Type},
		{"Project ID", event.ProjectID},
		{"Date Created", event.DateCreated.Format("2006-01-02 15:04:05")},
		{"Date Received", event.DateReceived.Format("2006-01-02 15:04:05")},
		{"Size", strconv.FormatInt(event.Size, 10)},
	}

	if event.GroupID != "" {
		rows = append(rows, []string{"Group ID", event.GroupID})
	}
	if event.Logger != "" {
		rows = append(rows, []string{"Logger", event.Logger})
	}
	if event.Culprit != "" {
		rows = append(rows, []string{"Culprit", event.Culprit})
	}
	if event.Environment != "" {
		rows = append(rows, []string{"Environment", event.Environment})
	}

	for _, row := range rows {
		err := table.Append(row)
		if err != nil {
			return err
		}
	}

	table.Render()
	return nil
}

// FormatEvents formats multiple events as a table
func (f *TableFormatter) FormatEvents(events []models.Event) error {
	if len(events) == 0 {
		fmt.Fprintf(f.writer, "No events found\n")
		return nil
	}

	table := tablewriter.NewWriter(f.writer)
	table.Header("ID", "Title", "Type", "Platform", "Project ID", "Date Created", "Environment")

	for _, event := range events {
		row := []string{
			event.EventID,
			truncateString(event.Title, 30),
			event.Type,
			event.Platform,
			event.ProjectID,
			event.DateCreated.Format("2006-01-02 15:04"),
			event.Environment,
		}
		err := table.Append(row)
		if err != nil {
			return err
		}
	}

	table.Render()
	return nil
}

// FormatIssue formats a single issue as a table
func (f *TableFormatter) FormatIssue(issue *models.Issue) error {
	table := tablewriter.NewWriter(f.writer)
	table.Header("Field", "Value")

	rows := [][]string{
		{"ID", issue.ID},
		{"Short ID", issue.ShortID},
		{"Title", issue.Title},
		{"Level", issue.Level},
		{"Status", issue.Status},
		{"Platform", issue.Platform},
		{"Project", fmt.Sprintf("%s (%s)", issue.Project.Name, issue.Project.Slug)},
		{"Count", issue.Count},
		{"User Count", strconv.Itoa(issue.UserCount)},
		{"First Seen", issue.FirstSeen.Format("2006-01-02 15:04:05")},
		{"Last Seen", issue.LastSeen.Format("2006-01-02 15:04:05")},
		{"Is Public", strconv.FormatBool(issue.IsPublic)},
		{"Is Bookmarked", strconv.FormatBool(issue.IsBookmarked)},
		{"Is Subscribed", strconv.FormatBool(issue.IsSubscribed)},
	}

	if issue.Substatus != "" {
		rows = append(rows, []string{"Substatus", issue.Substatus})
	}
	if issue.Priority != "" {
		rows = append(rows, []string{"Priority", issue.Priority})
	}
	if issue.Culprit != "" {
		rows = append(rows, []string{"Culprit", issue.Culprit})
	}
	if issue.Logger != "" {
		rows = append(rows, []string{"Logger", issue.Logger})
	}

	for _, row := range rows {
		err := table.Append(row)
		if err != nil {
			return err
		}
	}

	table.Render()
	return nil
}

// FormatIssues formats multiple issues as a table
func (f *TableFormatter) FormatIssues(issues []models.Issue) error {
	if len(issues) == 0 {
		fmt.Fprintf(f.writer, "No issues found\n")
		return nil
	}

	table := tablewriter.NewWriter(f.writer)
	table.Header("ID", "Title", "Level", "Status", "Count", "User Count", "Last Seen", "Project")

	for _, issue := range issues {
		row := []string{
			issue.ShortID,
			truncateString(issue.Title, 30),
			issue.Level,
			issue.Status,
			issue.Count,
			strconv.Itoa(issue.UserCount),
			issue.LastSeen.Format("01-02 15:04"),
			issue.Project.Slug,
		}
		err := table.Append(row)
		if err != nil {
			return err
		}
	}

	table.Render()
	return nil
}

// FormatProject formats a single project as a table
func (f *TableFormatter) FormatProject(project *models.Project) error {
	table := tablewriter.NewWriter(f.writer)
	table.Header("Field", "Value")

	rows := [][]string{
		{"ID", project.ID},
		{"Slug", project.Slug},
		{"Name", project.Name},
		{"Platform", project.Platform},
		{"Organization", fmt.Sprintf("%s (%s)", project.Organization.Name, project.Organization.Slug)},
		{"Date Created", project.DateCreated.Format("2006-01-02 15:04:05")},
		{"Status", project.Status},
		{"Is Public", strconv.FormatBool(project.IsPublic)},
		{"Is Bookmarked", strconv.FormatBool(project.IsBookmarked)},
	}

	for _, row := range rows {
		err := table.Append(row)
		if err != nil {
			return err
		}
	}

	table.Render()
	return nil
}

// FormatProjects formats multiple projects as a table
func (f *TableFormatter) FormatProjects(projects []models.Project) error {
	if len(projects) == 0 {
		fmt.Fprintf(f.writer, "No projects found\n")
		return nil
	}

	table := tablewriter.NewWriter(f.writer)
	table.Header("Slug", "Name", "Platform", "Organization", "Status", "Date Created")

	for _, project := range projects {
		row := []string{
			project.Slug,
			truncateString(project.Name, 25),
			project.Platform,
			project.Organization.Slug,
			project.Status,
			project.DateCreated.Format("2006-01-02"),
		}
		err := table.Append(row)
		if err != nil {
			return err
		}
	}

	table.Render()
	return nil
}

// FormatOrgStats formats organization stats as a table
func (f *TableFormatter) FormatOrgStats(stats *models.OrganizationStats) error {
	table := tablewriter.NewWriter(f.writer)
	table.Header("Metric", "Value")

	rows := [][]string{
		{"Start Time", stats.Start.Format("2006-01-02 15:04:05")},
		{"End Time", stats.End.Format("2006-01-02 15:04:05")},
		{"Total Sum", strconv.FormatInt(stats.Totals.Sum, 10)},
		{"Times Seen", strconv.FormatInt(stats.Totals.TimesSeen, 10)},
	}

	for _, row := range rows {
		err := table.Append(row)
		if err != nil {
			return err
		}
	}

	table.Render()
	return nil
}

// FormatGeneric formats any data as a table by reflecting on its structure
func (f *TableFormatter) FormatGeneric(data interface{}) error {
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
			// Fallback to simple key-value table for unknown types
			return f.formatUnknownSlice(data)
		}
	}

	// Handle single values
	return f.formatSingleValue(data)
}

// formatUnknownSlice formats a slice of unknown type
func (f *TableFormatter) formatUnknownSlice(data interface{}) error {
	v := reflect.ValueOf(data)
	
	table := tablewriter.NewWriter(f.writer)
	table.Header("Index", "Value")

	for i := 0; i < v.Len(); i++ {
		elem := v.Index(i)
		row := []string{
			strconv.Itoa(i),
			fmt.Sprintf("%v", elem.Interface()),
		}
		err := table.Append(row)
		if err != nil {
			return err
		}
	}

	table.Render()
	return nil
}

// formatSingleValue formats a single value as a key-value table
func (f *TableFormatter) formatSingleValue(data interface{}) error {
	table := tablewriter.NewWriter(f.writer)
	table.Header("Field", "Value")

	v := reflect.ValueOf(data)
	t := reflect.TypeOf(data)

	// Handle pointers
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			table.Append([]string{"Value", "<nil>"})
			table.Render()
			return nil
		}
		v = v.Elem()
		t = t.Elem()
	}

	// Handle structs
	if v.Kind() == reflect.Struct {
		for i := 0; i < v.NumField(); i++ {
			field := t.Field(i)
			value := v.Field(i)
			
			if !value.CanInterface() {
				continue
			}

			fieldValue := fmt.Sprintf("%v", value.Interface())
			table.Append([]string{field.Name, fieldValue})
		}
	} else {
		// Simple value
		table.Append([]string{"Value", fmt.Sprintf("%v", data)})
	}

	table.Render()
	return nil
}

