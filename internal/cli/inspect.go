package cli

import (
	"fmt"
	"net/url"
	"regexp"
	"sentire/internal/api"
	"sentire/internal/cli/formatter"
	"sentire/internal/client"

	"github.com/spf13/cobra"
)

var inspectCmd = &cobra.Command{
	Use:   "inspect <url>",
	Short: "Inspect a Sentry issue from its URL",
	Long:  "Parse a Sentry issue URL and display the recommended event with full debugging details",
	Args:  cobra.ExactArgs(1),
	RunE:  runInspect,
}

func init() {
	rootCmd.AddCommand(inspectCmd)
}

// SentryURLParts contains extracted parts from a Sentry URL
type SentryURLParts struct {
	Organization string
	IssueID      string
}

// parseSentryURL extracts organization and issue ID from a Sentry URL
func parseSentryURL(rawURL string) (*SentryURLParts, error) {
	// Parse the URL
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL format: %w", err)
	}

	// Extract organization from subdomain
	// Expected format: https://orgname.sentry.io/...
	host := parsedURL.Host
	sentryRegex := regexp.MustCompile(`^([^.]+)\.sentry\.io$`)
	matches := sentryRegex.FindStringSubmatch(host)
	if len(matches) < 2 {
		return nil, fmt.Errorf("invalid Sentry URL: expected format https://orgname.sentry.io/...")
	}
	organization := matches[1]

	// Extract issue ID from path
	// Expected format: /issues/123456789/
	path := parsedURL.Path
	issueRegex := regexp.MustCompile(`/issues/(\d+)/?`)
	issueMatches := issueRegex.FindStringSubmatch(path)
	if len(issueMatches) < 2 {
		return nil, fmt.Errorf("invalid issue URL: expected format /issues/<issue_id>/")
	}
	issueID := issueMatches[1]

	return &SentryURLParts{
		Organization: organization,
		IssueID:      issueID,
	}, nil
}

func runInspect(cmd *cobra.Command, args []string) error {
	sentryURL := args[0]

	// Parse the URL to extract organization and issue ID
	parts, err := parseSentryURL(sentryURL)
	if err != nil {
		return fmt.Errorf("failed to parse Sentry URL: %w", err)
	}

	// Create API client
	c, err := client.NewClient()
	if err != nil {
		return err
	}

	eventsAPI := api.NewEventsAPI(c)

	// Get the recommended event for the issue
	event, err := eventsAPI.GetIssueEvent(parts.Organization, parts.IssueID, "recommended", nil)
	if err != nil {
		return fmt.Errorf("failed to retrieve issue event: %w", err)
	}

	// Output the event data
	return formatter.Output(cmd, event)
}
