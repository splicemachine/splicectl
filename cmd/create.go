package cmd

import (
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Args:  cobra.MinimumNArgs(1),
	Short: "Create cluster resources",
	Long: `EXAMPLES
	splicectl create splice-database --database-name splicedb --dnsprefix splicedb`,
	Run: func(cmd *cobra.Command, args []string) {},
}

func init() {
	rootCmd.AddCommand(createCmd)
}
