package cmd

import (
	"github.com/spf13/cobra"
)

var changelogCmd = &cobra.Command{
	Use:   "changelog",
	Short: "List the most recent changes on the changelog for splicectl.",
	Long: `EXAMPLES
	splicectl changelog
`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: output the changelog var...
	},
}

func init() {
	rootCmd.AddCommand(changelogCmd)
}
