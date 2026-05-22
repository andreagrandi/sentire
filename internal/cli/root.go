package cli

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/andreagrandi/sentire/internal/version"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "sentire",
	Short:   "A command-line tool for the Sentry API",
	Version: version.Version,
	Long: `Sentire is a simple and user-friendly command-line interface for interacting with the Sentry API.
It allows you to query events, issues, projects, and organizations directly from your terminal.

Before using sentire, make sure to set your Sentry API token:
  export SENTRY_API_TOKEN=your_token_here`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute runs the root command. The command context is canceled on SIGINT
// (Ctrl+C) or SIGTERM so in-flight API requests stop promptly.
func Execute() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		format, _ := rootCmd.PersistentFlags().GetString("format")
		writeErrorOutput(os.Stderr, err, format)
		os.Exit(exitCodeFromError(err))
	}
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringP("format", "f", "json", "Output format: json, ndjson, table, text, markdown")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Verbose output")
	rootCmd.PersistentFlags().String("fields", "", "Comma-separated list of fields to include in JSON output")
}
