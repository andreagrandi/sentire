package cli

import (
	"fmt"

	"github.com/andreagrandi/sentire/internal/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Long:  "Display version, build time, and build information for sentire",
	Run: func(cmd *cobra.Command, args []string) {
		detailed, _ := cmd.Flags().GetBool("detailed")
		if detailed {
			fmt.Println(version.GetFullVersionInfo())
		} else {
			fmt.Println(version.GetVersionInfo())
		}
	},
}

func init() {
	versionCmd.Flags().BoolP("detailed", "d", false, "Show detailed version information")
	rootCmd.AddCommand(versionCmd)
}
