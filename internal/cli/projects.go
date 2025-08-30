package cli

import (
	"sentire/internal/api"
	"sentire/internal/client"

	"github.com/spf13/cobra"
)

var projectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "Manage Sentry projects",
	Long:  "Commands for interacting with Sentry projects",
}

var listProjectsCmd = &cobra.Command{
	Use:   "list",
	Short: "List all your projects",
	Long:  "Retrieve a list of all projects you have access to",
	Args:  cobra.NoArgs,
	RunE:  runListProjects,
}

var getProjectCmd = &cobra.Command{
	Use:   "get <organization> <project>",
	Short: "Get a specific project",
	Long:  "Retrieve details for a specific project",
	Args:  cobra.ExactArgs(2),
	RunE:  runGetProject,
}

func init() {
	rootCmd.AddCommand(projectsCmd)
	
	projectsCmd.AddCommand(listProjectsCmd)
	projectsCmd.AddCommand(getProjectCmd)

	// Flags for list command
	listProjectsCmd.Flags().Bool("all", false, "Fetch all pages")
}

func runListProjects(cmd *cobra.Command, args []string) error {
	c, err := client.NewClient()
	if err != nil {
		return err
	}

	projectsAPI := api.NewProjectsAPI(c)

	fetchAll, _ := cmd.Flags().GetBool("all")
	
	var allProjects []interface{}
	cursor := ""

	for {
		opts := &api.ListAllProjectsOptions{}
		if cursor != "" {
			opts.Cursor = cursor
		}

		projects, pagination, err := projectsAPI.ListProjects(opts)
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

func runGetProject(cmd *cobra.Command, args []string) error {
	orgSlug, projectSlug := args[0], args[1]

	c, err := client.NewClient()
	if err != nil {
		return err
	}

	projectsAPI := api.NewProjectsAPI(c)
	project, err := projectsAPI.GetProject(orgSlug, projectSlug)
	if err != nil {
		return err
	}

	return outputJSON(project)
}