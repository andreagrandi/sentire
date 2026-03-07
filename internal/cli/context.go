package cli

import (
	_ "embed"
	"fmt"

	"github.com/spf13/cobra"
)

//go:embed context_content.md
var contextContent string

var contextCmd = &cobra.Command{
	Use:   "context",
	Short: "Print agent-friendly context and usage guide",
	Long:  "Output the CONTEXT.md file with guidance for AI agents on how to use sentire effectively",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Print(contextContent)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(contextCmd)
}
