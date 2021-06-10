package cmd

import (
	"github.com/spf13/cobra"
)

var restartCmd = &cobra.Command{
	Use:   "restart",
	Args:  cobra.MinimumNArgs(1),
	Short: "Restart components of the Splice DB Cluster",
	Long: `EXAMPLES
	splicectl restart workspace --database-name <dbname>`,
	Run: func(cmd *cobra.Command, args []string) {},
}

func init() {
	RootCmd.AddCommand(restartCmd)
}
