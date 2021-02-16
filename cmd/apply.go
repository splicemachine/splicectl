package cmd

import (
	"github.com/spf13/cobra"
)

var applyCmd = &cobra.Command{
	Use:   "apply",
	Args:  cobra.MinimumNArgs(1),
	Short: "Apply configurations to various resources of the Splice Machine Database Cluster",
	Long: `EXAMPLES
	splicectl get system-settings > ~/tmp/system-settings.json
	# edit the file
	splicectl apply system-settings --file ~/tmp/system-settings.json
	`,
	Run: func(cmd *cobra.Command, args []string) {},
}

func init() {
	rootCmd.AddCommand(applyCmd)
}
