package create

import (
	"github.com/spf13/cobra"
	c "github.com/splicemachine/splicectl/cmd"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Args:  cobra.MinimumNArgs(1),
	Short: "Create cluster resources",
	Long: `EXAMPLES
	splicectl create workspace --database-name splicedb --dnsprefix splicedb`,
	Run: func(cmd *cobra.Command, args []string) {},
}

func init() {
	c.RootCmd.AddCommand(createCmd)
}
