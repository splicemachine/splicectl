package cmd

import (
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Args:  cobra.MinimumNArgs(1),
	Short: "List resources in the Splice Machine workspace Cluster",
	Long: `EXAMPLES
	splicectl list workspace`,
	Run: func(cmd *cobra.Command, args []string) {},
}

func init() {
	RootCmd.AddCommand(listCmd)
}
