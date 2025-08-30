package cli

import (
	"sentire/internal/api"
	"sentire/internal/client"

	"github.com/spf13/cobra"
)

var orgCmd = &cobra.Command{
	Use:   "org",
	Short: "Manage Sentry organizations",
	Long:  "Commands for interacting with Sentry organizations",
}

var listOrgProjectsCmd = &cobra.Command{
	Use:   "list-projects <organization>",
	Short: "List projects for an organization",
	Long:  "Retrieve a list of projects for a specific organization",
	Args:  cobra.ExactArgs(1),
	RunE:  runListOrgProjects,
}

var getOrgStatsCmd = &cobra.Command{
	Use:   "stats <organization>",
	Short: "Get organization statistics",
	Long:  "Retrieve event statistics for an organization",
	Args:  cobra.ExactArgs(1),
	RunE:  runGetOrgStats,
}

func init() {
	rootCmd.AddCommand(orgCmd)
	
	orgCmd.AddCommand(listOrgProjectsCmd)
	orgCmd.AddCommand(getOrgStatsCmd)

	// Flags for list-projects command
	listOrgProjectsCmd.Flags().Bool("all", false, "Fetch all pages")

	// Flags for stats command
	getOrgStatsCmd.Flags().String("field", "sum(quantity)", "Field to query: sum(quantity) or sum(times_seen)")
	getOrgStatsCmd.Flags().String("period", "", "Time period (e.g., '1d', '7d')")
	getOrgStatsCmd.Flags().String("interval", "", "Time series resolution")
	getOrgStatsCmd.Flags().String("start", "", "Start time (ISO-8601)")
	getOrgStatsCmd.Flags().String("end", "", "End time (ISO-8601)")
	getOrgStatsCmd.Flags().StringSlice("project", nil, "Filter by project IDs")
	getOrgStatsCmd.Flags().StringSlice("category", nil, "Filter by event categories")
	getOrgStatsCmd.Flags().StringSlice("outcome", nil, "Filter by event outcomes")
	getOrgStatsCmd.Flags().StringSlice("reason", nil, "Filter by event reasons")
	getOrgStatsCmd.Flags().Bool("download", false, "Download response as CSV")
}

func runListOrgProjects(cmd *cobra.Command, args []string) error {
	orgSlug := args[0]

	c, err := client.NewClient()
	if err != nil {
		return err
	}

	orgAPI := api.NewOrganizationsAPI(c)

	fetchAll, _ := cmd.Flags().GetBool("all")
	
	var allProjects []interface{}
	cursor := ""

	for {
		opts := &api.ListProjectsOptions{}
		if cursor != "" {
			opts.Cursor = cursor
		}

		projects, pagination, err := orgAPI.ListProjects(orgSlug, opts)
		if err != nil {
			return err
		}

		for _, project := range projects {
			allProjects = append(allProjects, project)
		}

		if !fetchAll || pagination == nil || !pagination.HasNext {
			break
		}
		cursor = pagination.NextCursor
	}

	return outputJSON(allProjects)
}

func runGetOrgStats(cmd *cobra.Command, args []string) error {
	orgSlug := args[0]

	c, err := client.NewClient()
	if err != nil {
		return err
	}

	orgAPI := api.NewOrganizationsAPI(c)

	opts := &api.GetStatsOptions{}
	if field, _ := cmd.Flags().GetString("field"); field != "" {
		opts.Field = field
	}
	if period, _ := cmd.Flags().GetString("period"); period != "" {
		opts.StatsPeriod = period
	}
	if interval, _ := cmd.Flags().GetString("interval"); interval != "" {
		opts.Interval = interval
	}
	if start, _ := cmd.Flags().GetString("start"); start != "" {
		opts.Start = start
	}
	if end, _ := cmd.Flags().GetString("end"); end != "" {
		opts.End = end
	}
	if projects, _ := cmd.Flags().GetStringSlice("project"); len(projects) > 0 {
		opts.Project = projects
	}
	if categories, _ := cmd.Flags().GetStringSlice("category"); len(categories) > 0 {
		opts.Category = categories
	}
	if outcomes, _ := cmd.Flags().GetStringSlice("outcome"); len(outcomes) > 0 {
		opts.Outcome = outcomes
	}
	if reasons, _ := cmd.Flags().GetStringSlice("reason"); len(reasons) > 0 {
		opts.Reason = reasons
	}
	if download, _ := cmd.Flags().GetBool("download"); download {
		opts.Download = true
	}

	stats, err := orgAPI.GetStats(orgSlug, opts)
	if err != nil {
		return err
	}

	return outputJSON(stats)
}