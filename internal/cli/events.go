package cli

import (
	"encoding/json"
	"os"
	"sentire/internal/api"
	"sentire/internal/client"

	"github.com/spf13/cobra"
)

var eventsCmd = &cobra.Command{
	Use:   "events",
	Short: "Manage Sentry events and issues",
	Long:  "Commands for interacting with Sentry events and issues",
}

var listProjectEventsCmd = &cobra.Command{
	Use:   "list-project <organization> <project>",
	Short: "List events for a project",
	Long:  "Retrieve a list of events for a specific project",
	Args:  cobra.ExactArgs(2),
	RunE:  runListProjectEvents,
}

var listIssueEventsCmd = &cobra.Command{
	Use:   "list-issue <organization> <issue-id>",
	Short: "List events for an issue",
	Long:  "Retrieve a list of events for a specific issue",
	Args:  cobra.ExactArgs(2),
	RunE:  runListIssueEvents,
}

var listIssuesCmd = &cobra.Command{
	Use:   "list-issues <organization>",
	Short: "List issues for an organization",
	Long:  "Retrieve a list of issues for a specific organization",
	Args:  cobra.ExactArgs(1),
	RunE:  runListIssues,
}

var getEventCmd = &cobra.Command{
	Use:   "get-event <organization> <project> <event-id>",
	Short: "Get a specific event",
	Long:  "Retrieve details for a specific event in a project",
	Args:  cobra.ExactArgs(3),
	RunE:  runGetEvent,
}

var getIssueCmd = &cobra.Command{
	Use:   "get-issue <organization> <issue-id>",
	Short: "Get a specific issue",
	Long:  "Retrieve details for a specific issue",
	Args:  cobra.ExactArgs(2),
	RunE:  runGetIssue,
}

var getIssueEventCmd = &cobra.Command{
	Use:   "get-issue-event <organization> <issue-id> <event-id>",
	Short: "Get a specific event for an issue",
	Long:  "Retrieve a specific event associated with an issue. Event ID can be 'latest', 'oldest', 'recommended', or a specific event ID",
	Args:  cobra.ExactArgs(3),
	RunE:  runGetIssueEvent,
}

func init() {
	rootCmd.AddCommand(eventsCmd)
	
	eventsCmd.AddCommand(listProjectEventsCmd)
	eventsCmd.AddCommand(listIssueEventsCmd)
	eventsCmd.AddCommand(listIssuesCmd)
	eventsCmd.AddCommand(getEventCmd)
	eventsCmd.AddCommand(getIssueCmd)
	eventsCmd.AddCommand(getIssueEventCmd)

	// Flags for list-project command
	listProjectEventsCmd.Flags().String("period", "", "Time period (e.g., '24h', '7d')")
	listProjectEventsCmd.Flags().String("start", "", "Start time (ISO-8601)")
	listProjectEventsCmd.Flags().String("end", "", "End time (ISO-8601)")
	listProjectEventsCmd.Flags().Bool("full", false, "Include full event body")
	listProjectEventsCmd.Flags().Bool("sample", false, "Return events in pseudo-random order")
	listProjectEventsCmd.Flags().Bool("all", false, "Fetch all pages")

	// Flags for list-issue command
	listIssueEventsCmd.Flags().String("period", "", "Time period (e.g., '24h', '7d')")
	listIssueEventsCmd.Flags().String("start", "", "Start time (ISO-8601)")
	listIssueEventsCmd.Flags().String("end", "", "End time (ISO-8601)")
	listIssueEventsCmd.Flags().StringSlice("environment", nil, "Filter by environments")
	listIssueEventsCmd.Flags().Bool("full", false, "Include full event body")
	listIssueEventsCmd.Flags().Bool("sample", false, "Return events in pseudo-random order")
	listIssueEventsCmd.Flags().String("query", "", "Search query")
	listIssueEventsCmd.Flags().Bool("all", false, "Fetch all pages")

	// Flags for list-issues command
	listIssuesCmd.Flags().StringSlice("environment", nil, "Filter by environments")
	listIssuesCmd.Flags().StringSlice("project", nil, "Filter by project IDs")
	listIssuesCmd.Flags().String("period", "", "Time period (e.g., '24h', '7d')")
	listIssuesCmd.Flags().String("start", "", "Start time (ISO-8601)")
	listIssuesCmd.Flags().String("end", "", "End time (ISO-8601)")
	listIssuesCmd.Flags().String("query", "is:unresolved issue.priority:[high,medium]", "Search/filter query")
	listIssuesCmd.Flags().String("sort", "", "Sort order (date, freq, inbox)")
	listIssuesCmd.Flags().Int("limit", 0, "Maximum number of results")
	listIssuesCmd.Flags().Bool("all", false, "Fetch all pages")

	// Flags for get-issue-event command
	getIssueEventCmd.Flags().StringSlice("environment", nil, "Filter by environments")
}

func runListProjectEvents(cmd *cobra.Command, args []string) error {
	orgSlug, projectSlug := args[0], args[1]

	c, err := client.NewClient()
	if err != nil {
		return err
	}

	eventsAPI := api.NewEventsAPI(c)

	opts := &api.ListProjectEventsOptions{}
	if period, _ := cmd.Flags().GetString("period"); period != "" {
		opts.StatsPeriod = period
	}
	if start, _ := cmd.Flags().GetString("start"); start != "" {
		opts.Start = start
	}
	if end, _ := cmd.Flags().GetString("end"); end != "" {
		opts.End = end
	}
	if full, _ := cmd.Flags().GetBool("full"); full {
		opts.Full = true
	}
	if sample, _ := cmd.Flags().GetBool("sample"); sample {
		opts.Sample = true
	}

	fetchAll, _ := cmd.Flags().GetBool("all")
	
	var allEvents []interface{}
	cursor := ""

	for {
		if cursor != "" {
			opts.Cursor = cursor
		}

		events, pagination, err := eventsAPI.ListProjectEvents(orgSlug, projectSlug, opts)
		if err != nil {
			return err
		}

		for _, event := range events {
			allEvents = append(allEvents, event)
		}

		if !fetchAll || pagination == nil || !pagination.HasNext {
			break
		}
		cursor = pagination.NextCursor
	}

	return outputJSON(allEvents)
}

func runListIssueEvents(cmd *cobra.Command, args []string) error {
	orgSlug, issueID := args[0], args[1]

	c, err := client.NewClient()
	if err != nil {
		return err
	}

	eventsAPI := api.NewEventsAPI(c)

	opts := &api.ListIssueEventsOptions{}
	if period, _ := cmd.Flags().GetString("period"); period != "" {
		opts.StatsPeriod = period
	}
	if start, _ := cmd.Flags().GetString("start"); start != "" {
		opts.Start = start
	}
	if end, _ := cmd.Flags().GetString("end"); end != "" {
		opts.End = end
	}
	if environments, _ := cmd.Flags().GetStringSlice("environment"); len(environments) > 0 {
		opts.Environment = environments
	}
	if full, _ := cmd.Flags().GetBool("full"); full {
		opts.Full = true
	}
	if sample, _ := cmd.Flags().GetBool("sample"); sample {
		opts.Sample = true
	}
	if query, _ := cmd.Flags().GetString("query"); query != "" {
		opts.Query = query
	}

	fetchAll, _ := cmd.Flags().GetBool("all")
	
	var allEvents []interface{}
	cursor := ""

	for {
		if cursor != "" {
			opts.Cursor = cursor
		}

		events, pagination, err := eventsAPI.ListIssueEvents(orgSlug, issueID, opts)
		if err != nil {
			return err
		}

		for _, event := range events {
			allEvents = append(allEvents, event)
		}

		if !fetchAll || pagination == nil || !pagination.HasNext {
			break
		}
		cursor = pagination.NextCursor
	}

	return outputJSON(allEvents)
}

func runListIssues(cmd *cobra.Command, args []string) error {
	orgSlug := args[0]

	c, err := client.NewClient()
	if err != nil {
		return err
	}

	eventsAPI := api.NewEventsAPI(c)

	opts := &api.ListIssuesOptions{}
	if environments, _ := cmd.Flags().GetStringSlice("environment"); len(environments) > 0 {
		opts.Environment = environments
	}
	if projects, _ := cmd.Flags().GetStringSlice("project"); len(projects) > 0 {
		opts.Project = projects
	}
	if period, _ := cmd.Flags().GetString("period"); period != "" {
		opts.StatsPeriod = period
	}
	if start, _ := cmd.Flags().GetString("start"); start != "" {
		opts.Start = start
	}
	if end, _ := cmd.Flags().GetString("end"); end != "" {
		opts.End = end
	}
	if query, _ := cmd.Flags().GetString("query"); query != "" {
		opts.Query = query
	}
	if sort, _ := cmd.Flags().GetString("sort"); sort != "" {
		opts.Sort = sort
	}
	if limit, _ := cmd.Flags().GetInt("limit"); limit > 0 {
		opts.Limit = limit
	}

	fetchAll, _ := cmd.Flags().GetBool("all")
	
	var allIssues []interface{}
	cursor := ""

	for {
		if cursor != "" {
			opts.Cursor = cursor
		}

		issues, pagination, err := eventsAPI.ListIssues(orgSlug, opts)
		if err != nil {
			return err
		}

		for _, issue := range issues {
			allIssues = append(allIssues, issue)
		}

		if !fetchAll || pagination == nil || !pagination.HasNext {
			break
		}
		cursor = pagination.NextCursor
	}

	return outputJSON(allIssues)
}

func runGetEvent(cmd *cobra.Command, args []string) error {
	orgSlug, projectSlug, eventID := args[0], args[1], args[2]

	c, err := client.NewClient()
	if err != nil {
		return err
	}

	eventsAPI := api.NewEventsAPI(c)
	event, err := eventsAPI.GetProjectEvent(orgSlug, projectSlug, eventID)
	if err != nil {
		return err
	}

	return outputJSON(event)
}

func runGetIssue(cmd *cobra.Command, args []string) error {
	orgSlug, issueID := args[0], args[1]

	c, err := client.NewClient()
	if err != nil {
		return err
	}

	eventsAPI := api.NewEventsAPI(c)
	issue, err := eventsAPI.GetIssue(orgSlug, issueID)
	if err != nil {
		return err
	}

	return outputJSON(issue)
}

func runGetIssueEvent(cmd *cobra.Command, args []string) error {
	orgSlug, issueID, eventID := args[0], args[1], args[2]

	c, err := client.NewClient()
	if err != nil {
		return err
	}

	eventsAPI := api.NewEventsAPI(c)

	opts := &api.GetIssueEventOptions{}
	if environments, _ := cmd.Flags().GetStringSlice("environment"); len(environments) > 0 {
		opts.Environment = environments
	}

	event, err := eventsAPI.GetIssueEvent(orgSlug, issueID, eventID, opts)
	if err != nil {
		return err
	}

	return outputJSON(event)
}

func outputJSON(data interface{}) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}