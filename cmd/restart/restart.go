package restart

import (
	"github.com/spf13/cobra"
	c "github.com/splicemachine/splicectl/cmd"
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
	c.RootCmd.AddCommand(restartCmd)
}
