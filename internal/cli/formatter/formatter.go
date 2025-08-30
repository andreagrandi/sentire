package formatter

import (
	"fmt"
	"io"
	"os"
	"sentire/pkg/models"

	"github.com/spf13/cobra"
)

// Formatter defines the interface for different output formats
type Formatter interface {
	FormatEvent(event *models.Event) error
	FormatEvents(events []models.Event) error
	FormatIssue(issue *models.Issue) error
	FormatIssues(issues []models.Issue) error
	FormatProject(project *models.Project) error
	FormatProjects(projects []models.Project) error
	FormatOrgStats(stats *models.OrganizationStats) error
	FormatGeneric(data interface{}) error
}

// NewFormatter creates a formatter based on the format flag
func NewFormatter(cmd *cobra.Command, writer io.Writer) (Formatter, error) {
	format, _ := cmd.Flags().GetString("format")
	if writer == nil {
		writer = os.Stdout
	}

	switch format {
	case "json":
		return NewJSONFormatter(writer), nil
	case "table":
		return NewTableFormatter(writer), nil
	case "text":
		return NewTextFormatter(writer), nil
	case "markdown":
		return NewMarkdownFormatter(writer), nil
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
}

// Output is the main output function that replaces outputJSON
func Output(cmd *cobra.Command, data interface{}) error {
	formatter, err := NewFormatter(cmd, nil)
	if err != nil {
		return err
	}

	// Try to cast to specific types first, fallback to generic
	switch v := data.(type) {
	case *models.Event:
		return formatter.FormatEvent(v)
	case []models.Event:
		return formatter.FormatEvents(v)
	case *models.Issue:
		return formatter.FormatIssue(v)
	case []models.Issue:
		return formatter.FormatIssues(v)
	case *models.Project:
		return formatter.FormatProject(v)
	case []models.Project:
		return formatter.FormatProjects(v)
	case *models.OrganizationStats:
		return formatter.FormatOrgStats(v)
	case []interface{}:
		// Handle mixed type slices (common in current code)
		return formatter.FormatGeneric(v)
	default:
		return formatter.FormatGeneric(v)
	}
}
