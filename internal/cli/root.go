package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sentire",
	Short: "A command-line tool for the Sentry API",
	Long: `Sentire is a simple and user-friendly command-line interface for interacting with the Sentry API.
It allows you to query events, issues, projects, and organizations directly from your terminal.

Before using sentire, make sure to set your Sentry API token:
  export SENTRY_API_TOKEN=your_token_here`,
	SilenceUsage: true,
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringP("format", "f", "json", "Output format: json, table")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Verbose output")
}