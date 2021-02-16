package cmd

import (
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Args:  cobra.MinimumNArgs(1),
	Short: "Get various resources of the Splice Machiene Database Cluster",
	Long: `EXAMPLES
	splicectl get default-cr > ~/tmp/default-cr.json
	splicectl get system-settings > ~/tmp/system-settings.json
	splicectl get database-status --database-name "test" `,
	Run: func(cmd *cobra.Command, args []string) {},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
